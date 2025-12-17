package curseforge

import (
	"FancyVerteiler/internal/config"
	"FancyVerteiler/internal/git"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type Service struct {
	git    *git.Service
	hc     *http.Client
	apiKey string
}

func New(apiKey string, git *git.Service) *Service {
	return &Service{
		git:    git,
		hc:     &http.Client{},
		apiKey: apiKey,
	}
}

func (s *Service) Deploy(cfg *config.DeploymentConfig) error {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Add metadata JSON
	metadata, err := s.metadataJson(cfg)
	if err != nil {
		return err
	}

	_ = writer.WriteField("metadata", metadata)

	// Add the plugin file
	ver, err := cfg.Version()
	if err != nil {
		return err
	}

	pluginJarPath := config.BasePath + cfg.PluginJarPath
	pluginJarPath = strings.ReplaceAll(pluginJarPath, "%VERSION%", ver)
	file, err := os.Open(pluginJarPath)
	if err != nil {
		return err
	}
	defer file.Close()

	fileWriter, err := writer.CreateFormFile("file", filepath.Base(pluginJarPath))
	if err != nil {
		return err
	}

	_, err = io.Copy(fileWriter, file)
	if err != nil {
		return err
	}

	// Close the writer to finalize the multipart form
	err = writer.Close()
	if err != nil {
		return err
	}

	// Create the request
	url := fmt.Sprintf("https://minecraft.curseforge.com/api/projects/%s/upload-file", cfg.CurseForge.ProjectID)
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return err
	}

	// Set headers
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("X-Api-Token", s.apiKey)
	req.Header.Set("User-Agent", "FancyVerteiler (https://github.com/FancyInnovations/FancyVerteiler)")

	resp, err := s.hc.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to create version (status %d): %s", resp.StatusCode, string(respBody))
	}

	return nil
}

func (s *Service) metadataJson(cfg *config.DeploymentConfig) (string, error) {
	ver, err := cfg.Version()
	if err != nil {
		return "", err
	}

	cl, err := cfg.Changelog()
	if err != nil {
		return "", err
	}
	cl = strings.ReplaceAll(cl, "%COMMIT_HASH%", s.git.CommitSHA())
	cl = strings.ReplaceAll(cl, "%COMMIT_MESSAGE%", s.git.CommitMessage())

	// Convert game versions (supports both int IDs and string versions)
	gameVersionIDs := make([]int, 0, len(cfg.CurseForge.GameVersions)+1)
	
	// Get the appropriate loader ID based on project type
	projectType := cfg.CurseForge.Type
	if projectType == "" {
		projectType = "plugin" // Default to plugin for backward compatibility
	}
	
	loaderID, err := GetLoaderID(projectType, cfg.CurseForge.Loader)
	if err != nil {
		return "", err
	}
	
	// Add loader ID first
	gameVersionIDs = append(gameVersionIDs, loaderID)
	
	for _, v := range cfg.CurseForge.GameVersions {
		switch val := v.(type) {
		case float64: // JSON numbers are parsed as float64
			gameVersionIDs = append(gameVersionIDs, int(val))
		case int:
			gameVersionIDs = append(gameVersionIDs, val)
		case string:
			if id, ok := ConvertVersionString(val); ok {
				gameVersionIDs = append(gameVersionIDs, id)
			} else {
				return "", fmt.Errorf("unknown Minecraft version: %s", val)
			}
		default:
			return "", fmt.Errorf("invalid game version type: %T", v)
		}
	}

	req := CreateVersionReq{
		Changelog:     cl,
		ChangelogType: "markdown",
		DisplayName:   ver,
		GameVersions:  gameVersionIDs,
		ReleaseType:   cfg.CurseForge.ReleaseType,
		Relations:     nil, // Set to nil initially
	}

	// Only add relations if they are provided in config
	if cfg.CurseForge.Relations != nil && len(cfg.CurseForge.Relations.Projects) > 0 {
		relations := &CreateVersionRelations{
			Projects: make([]ProjectRelation, len(cfg.CurseForge.Relations.Projects)),
		}
		for i, proj := range cfg.CurseForge.Relations.Projects {
			relations.Projects[i] = ProjectRelation{
				Slug: proj.Slug,
				Type: proj.Type,
			}
		}
		req.Relations = relations
	}

	data, err := json.Marshal(req)
	if err != nil {
		return "", err
	}

	return string(data), nil
}