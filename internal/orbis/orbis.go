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

	if err := s.updateChangelog(cfg, versionID); err != nil {
		return fmt.Errorf("failed to update changelog: %w", err)
	}

	versionFileID, err := s.uploadFile(cfg, versionID)
	if err != nil {
		return fmt.Errorf("failed to upload file: %w", err)
	}

	if err := s.setPrimaryVersionFile(cfg, versionID, versionFileID); err != nil {
		return fmt.Errorf("failed to set primary file: %w", err)
	}

	if err := s.submitForReview(cfg, versionID); err != nil {
		return fmt.Errorf("failed to submit for review: %w", err)
	}

	return nil
}

func (s *Service) createVersion(cfg *config.DeploymentConfig) (string, error) {
	ver, err := cfg.Version()
	if err != nil {
		return "", err
	}

	req := CreateVersionReq{
		VersionNumber:              ver,
		Name:                       ver,
		Channel:                    cfg.Orbis.Channel,
		CompatibleHytaleVersionIds: cfg.Orbis.CompatibleHytaleVersionIds,
	}

	data, err := json.Marshal(req)
	if err != nil {
		return "", err
	}

	reqBody, err := http.NewRequest("POST", "https://api.orbis.place/resources/"+cfg.Orbis.ResourceID+"/versions", strings.NewReader(string(data)))
	if err != nil {
		return "", err
	}
	reqBody.Header.Set("Content-Type", "application/json")
	reqBody.Header.Set("x-api-key", s.apiKey)
	reqBody.Header.Set("User-Agent", "FancyVerteiler (https://github.com/FancyInnovations/FancyVerteiler)")

	resp, err := s.hc.Do(reqBody)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return "", fmt.Errorf("failed to read body: %w", err)
		}

		return "", fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(body))
	}

	var respVer Version
	if err := json.NewDecoder(resp.Body).Decode(&respVer); err != nil {
		return "", err
	}

	return respVer.ID, nil
}

func (s *Service) updateChangelog(cfg *config.DeploymentConfig, versionID string) error {
	cl, err := cfg.Changelog()
	if err != nil {
		return err
	}
	cl = strings.ReplaceAll(cl, "%COMMIT_HASH%", s.git.CommitSHA())
	cl = strings.ReplaceAll(cl, "%COMMIT_MESSAGE%", s.git.CommitMessage())

	req := UpdateChangelogReq{
		Changelog: cl,
	}

	data, err := json.Marshal(req)
	if err != nil {
		return err
	}

	reqBody, err := http.NewRequest("PATCH", "https://api.orbis.place/resources/"+cfg.Orbis.ResourceID+"/versions/"+versionID+"/changelog", strings.NewReader(string(data)))
	if err != nil {
		return err
	}
	reqBody.Header.Set("Content-Type", "application/json")
	reqBody.Header.Set("x-api-key", s.apiKey)
	reqBody.Header.Set("User-Agent", "FancyVerteiler (https://github.com/FancyInnovations/FancyVerteiler)")

	resp, err := s.hc.Do(reqBody)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("failed to read body: %w", err)
		}

		return fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(body))
	}

	return nil
}

func (s *Service) uploadFile(cfg *config.DeploymentConfig, versionID string) (string, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	ver, err := cfg.Version()
	if err != nil {
		return "", err
	}

	pluginJarPath := config.BasePath + cfg.PluginJarPath
	pluginJarPath = strings.ReplaceAll(pluginJarPath, "%VERSION%", ver)
	file, err := os.Open(pluginJarPath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	fileWriter, err := writer.CreateFormFile("file", filepath.Base(pluginJarPath))
	if err != nil {
		return "", err
	}

	_, err = io.Copy(fileWriter, file)
	if err != nil {
		return "", err
	}

	// Close the writer to finalize the multipart form
	err = writer.Close()
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", "https://api.orbis.place/resources/"+cfg.Orbis.ResourceID+"/versions/"+versionID+"/files", body)
	if err != nil {
		return "", err
	}

	// Set the correct Content-Type with boundary
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("x-api-key", s.apiKey)
	req.Header.Set("User-Agent", "FancyVerteiler (https://github.com/FancyInnovations/FancyVerteiler)")

	resp, err := s.hc.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			return "", fmt.Errorf("failed to read body: %w", err)
		}

		return "", fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(respBody))
	}

	var respVerFile VersionFile
	if err := json.NewDecoder(resp.Body).Decode(&respVerFile); err != nil {
		return "", err
	}

	return respVerFile.ID, nil
}

func (s *Service) setPrimaryVersionFile(cfg *config.DeploymentConfig, versionId, fileId string) error {
	req := SetPrimaryFileReq{
		FileID: fileId,
	}

	data, err := json.Marshal(req)
	if err != nil {
		return err
	}

	reqBody, err := http.NewRequest("PATCH", "https://api.orbis.place/resources/"+cfg.Orbis.ResourceID+"/versions/"+versionId+"/files/primary", strings.NewReader(string(data)))
	if err != nil {
		return err
	}
	reqBody.Header.Set("Content-Type", "application/json")
	reqBody.Header.Set("x-api-key", s.apiKey)
	reqBody.Header.Set("User-Agent", "FancyVerteiler (https://github.com/FancyInnovations/FancyVerteiler)")

	resp, err := s.hc.Do(reqBody)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("failed to read body: %w", err)
		}

		return fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(body))
	}

	return nil
}

func (s *Service) submitForReview(cfg *config.DeploymentConfig, versionId string) error {
	reqBody, err := http.NewRequest("POST", "https://api.orbis.place/resources/"+cfg.Orbis.ResourceID+"/versions/"+versionId+"/submit", nil)
	if err != nil {
		return err
	}
	reqBody.Header.Set("x-api-key", s.apiKey)
	reqBody.Header.Set("User-Agent", "FancyVerteiler (https://github.com/FancyInnovations/FancyVerteiler)")

	resp, err := s.hc.Do(reqBody)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("failed to read body: %w", err)
		}

		return fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(body))
	}

	return nil
}
