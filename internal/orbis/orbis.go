package orbis

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
	versionID, err := s.createVersion(cfg)
	if err != nil {
		return fmt.Errorf("failed to create version: %w", err)
	}

	if err := s.uploadFile(cfg, versionID); err != nil {
		return fmt.Errorf("failed to upload file: %w", err)
	}
	return nil
}

func (s *Service) createVersion(cfg *config.DeploymentConfig) (string, error) {
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
		VersionNumber:      ver,
		Name:               ver,
		Channel:            cfg.Orbis.Channel,
		CompatibleVersions: cfg.Orbis.SupportedVersions,
		Changelog:          cl,
	}

	data, err := json.Marshal(req)
	if err != nil {
		return "", err
	}

	reqBody, err := http.NewRequest("POST", "https://api.orbis.place.net/resources/"+cfg.Orbis.ResourceID+"/versions", strings.NewReader(string(data)))
	if err != nil {
		return "", err
	}
	reqBody.Header.Set("Content-Type", "application/json")
	reqBody.Header.Set("Authorization", "Bearer "+s.apiKey)
	reqBody.Header.Set("User-Agent", "FancyVerteiler (https://github.com/FancyInnovations/FancyVerteiler)")

	resp, err := s.hc.Do(reqBody)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var respVer Version
	if err := json.NewDecoder(resp.Body).Decode(&respVer); err != nil {
		return "", err
	}

	return respVer.ID, nil
}

func (s *Service) uploadFile(cfg *config.DeploymentConfig, versionID string) error {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

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

	req, err := http.NewRequest("POST", "https://api.orbis.place.net/resources/"+cfg.Orbis.ResourceID+"/versions/"+versionID+"/files", body)
	if err != nil {
		return err
	}

	// Set the correct Content-Type with boundary
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", "Bearer "+s.apiKey)

	resp, err := s.hc.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to create version: %s", string(respBody))
	}

	return nil
}
