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

	req := CreateVersionReq{
		Changelog:     cl,
		ChangelogType: "markdown",
		DisplayName:   ver,
		GameVersions:  cfg.CurseForge.GameVersions,
		ReleaseType:   cfg.CurseForge.ReleaseType,
		Relations: CreateVersionRelations{
			Projects: []ProjectRelation{},
		},
	}

	data, err := json.Marshal(req)
	if err != nil {
		return "", err
	}

	return string(data), nil
}