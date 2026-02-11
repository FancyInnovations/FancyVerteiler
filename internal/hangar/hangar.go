package hangar

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

func (s *Service) GetJWT() (string, error) {
	resp, err := http.Post("https://hangar.papermc.io/api/v1/authenticate?apiKey="+s.apiKey, "text/plain", nil)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("failed to authenticate: %s", string(respBody))
	}

	var authResp AuthenticateResp
	err = json.NewDecoder(resp.Body).Decode(&authResp)
	if err != nil {
		return "", err
	}

	return authResp.Token, nil
}

func (s *Service) Deploy(cfg *config.DeploymentConfig) error {
	jwt, err := s.GetJWT()
	if err != nil {
		return err
	}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	data, err := s.dataJson(cfg)
	if err != nil {
		return err
	}

	_ = writer.WriteField("versionUpload", data)

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

	fileWriter, err := writer.CreateFormFile("files", filepath.Base(pluginJarPath))
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

	req, err := http.NewRequest("POST", "https://hangar.papermc.io/api/v1/projects/"+cfg.Hangar.Author+"/"+cfg.Hangar.ProjectID+"/upload", body)
	if err != nil {
		return err
	}

	// Set the correct Content-Type with boundary
	req.Header.Set("Authorization", "HangarAuth "+jwt)
	req.Header.Set("User-Agent", "FancyVerteiler (https://github.com/FancyInnovations/FancyVerteiler)")

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

func (s *Service) dataJson(cfg *config.DeploymentConfig) (string, error) {
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

	req := VersionUploadReq{
		Version:            ver,
		PluginDependencies: map[Platform][]PluginDependency{},
		PlatformDependencies: map[Platform][]string{
			PlatformPaper: cfg.Hangar.SupportedVersions,
		},
		Description: cl,
		Files: []MultipartFileOrURL{
			{
				Platforms:   []Platform{PlatformPaper},
				ExternalURL: nil,
			},
		},
		Channel: cfg.Hangar.Channel,
	}

	data, err := json.Marshal(req)
	if err != nil {
		return "", err
	}

	return string(data), nil
}
