package modrinth

import (
	"FancyVerteiler/internal/config"
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
	hc     *http.Client
	apiKey string
}

func New(apiKey string) *Service {
	return &Service{
		hc:     &http.Client{},
		apiKey: apiKey,
	}
}

func (s *Service) Deploy(cfg *config.DeploymentConfig) error {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	data, err := dataJson(cfg)
	if err != nil {
		return err
	}

	_ = writer.WriteField("data", data)

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

	fileWriter, err := writer.CreateFormFile("pluginFile", filepath.Base(pluginJarPath))
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

	req, err := http.NewRequest("POST", "https://api.modrinth.com/v2/version", body)
	if err != nil {
		return err
	}

	// Set the correct Content-Type with boundary
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", s.apiKey)

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

func dataJson(cfg *config.DeploymentConfig) (string, error) {
	ver, err := cfg.Version()
	if err != nil {
		return "", err
	}
	cl, err := cfg.Changelog()
	if err != nil {
		return "", err
	}

	req := CreateVersionReq{
		Name:          ver,
		VersionNumber: ver,
		Changelog:     cl,
		Dependencies:  []ProjectDependency{},
		GameVersions:  cfg.Modrinth.SupportedVersions,
		VersionType:   cfg.Modrinth.Channel,
		Loaders:       cfg.Modrinth.Loaders,
		Featured:      cfg.Modrinth.Featured,
		Status:        "listed",
		ProjectID:     cfg.Modrinth.ProjectID,
		FileParts:     []string{"pluginFile"},
		PrimaryFile:   "pluginFile",
	}

	data, err := json.Marshal(req)
	if err != nil {
		return "", err
	}

	return string(data), nil
}
