package fancyspaces

import (
	"FancyVerteiler/internal/config"
	"FancyVerteiler/internal/git"
	"bytes"
	"encoding/json"
	"fmt"
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
	if err := s.createVersion(cfg); err != nil {
		return fmt.Errorf("failed to create version: %w", err)
	}

	if err := s.uploadFile(cfg); err != nil {
		return fmt.Errorf("failed to upload file: %w", err)
	}

	return nil
}

func (s *Service) createVersion(cfg *config.DeploymentConfig) error {
	ver, err := cfg.Version()
	if err != nil {
		return err
	}

	cl, err := cfg.Changelog()
	if err != nil {
		return err
	}
	cl = strings.ReplaceAll(cl, "%COMMIT_HASH%", s.git.CommitSHA())
	cl = strings.ReplaceAll(cl, "%COMMIT_MESSAGE%", s.git.CommitMessage())

	req := CreateVersionReq{
		Name:                      ver,
		Platform:                  cfg.FancySpaces.Platform,
		Channel:                   cfg.FancySpaces.Channel,
		Changelog:                 cl,
		SupportedPlatformVersions: cfg.FancySpaces.SupportedVersions,
	}

	data, err := json.Marshal(req)
	if err != nil {
		return err
	}

	reqBody, err := http.NewRequest("POST", "https://fancyspaces.net/api/v1/spaces/"+cfg.FancySpaces.SpaceID+"/versions", strings.NewReader(string(data)))
	if err != nil {
		return err
	}
	reqBody.Header.Set("Content-Type", "application/json")
	reqBody.Header.Set("Authorization", s.apiKey)
	reqBody.Header.Set("User-Agent", "FancyVerteiler (https://github.com/FancyInnovations/FancyVerteiler)")

	resp, err := s.hc.Do(reqBody)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}

func (s *Service) uploadFile(cfg *config.DeploymentConfig) error {
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

	pluginJarData, err := os.ReadFile(pluginJarPath)
	if err != nil {
		return err
	}

	pluginJarName := filepath.Base(pluginJarPath)

	url := fmt.Sprintf("https://fancyspaces.net/api/v1/spaces/%s/versions/%s/files/%s", cfg.FancySpaces.SpaceID, ver, pluginJarName)
	reqBody, err := http.NewRequest("POST", url, bytes.NewReader(pluginJarData))
	if err != nil {
		return err
	}
	reqBody.Header.Set("Authorization", s.apiKey)
	reqBody.Header.Set("User-Agent", "FancyVerteiler (https://github.com/FancyInnovations/FancyVerteiler)")

	resp, err := s.hc.Do(reqBody)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}
