package discord

import (
	"FancyVerteiler/internal/config"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"path/filepath"
	"strings"
)

type Service struct {
	hc *http.Client
}

func New() *Service {
	return &Service{
		hc: &http.Client{},
	}
}

func (s *Service) SendSuccessMessage(webhookURL string, cfg *config.DeploymentConfig) error {
	desc, err := buildDescription(cfg)
	if err != nil {
		return err
	}

	msg := Message{
		Content: "Plugin deployment completed!",
		Embeds: []Embed{
			{
				Title:       fmt.Sprintf("New version of %s deployed!", cfg.ProjectName),
				Description: desc,
				Color:       0x00FF00,
			},
		},
	}

	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	resp, err := s.hc.Post(webhookURL, "application/json", bytes.NewReader(data))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("failed to send Discord message, status code: %d, and failed to read body: %v", resp.StatusCode, err)
		}
		slog.Debug("Discord webhook response status", slog.String("body", string(body)), slog.Int("status_code", resp.StatusCode))

		return fmt.Errorf("failed to send Discord message, status code: %d", resp.StatusCode)
	}

	return nil
}

func buildDescription(cfg *config.DeploymentConfig) (string, error) {
	ver, err := cfg.Version()
	if err != nil {
		return "", err
	}

	desc := fmt.Sprintf("**Version:** %s", ver)

	if cfg.FancySpaces != nil || cfg.Modrinth != nil {
		var channel string
		if cfg.FancySpaces != nil {
			channel = cfg.FancySpaces.Channel
		} else if cfg.Modrinth != nil {
			channel = cfg.Modrinth.Channel
		}

		desc += fmt.Sprintf("\n**Channel:** %s", channel)
	}

	if cfg.FancySpaces != nil {
		fileName := filepath.Base(cfg.PluginJarPath)
		fileName = strings.ReplaceAll(fileName, "%VERSION%", ver)
		desc += fmt.Sprintf("\n**FancySpaces:** https://fancyspaces.net/api/v1/spaces/%s/versions/%s/files/%s", cfg.FancySpaces.SpaceID, ver, fileName)
	}

	if cfg.Modrinth != nil {
		desc += fmt.Sprintf("\n**Modrinth:** https://modrinth.com/plugin/%s/version/%s", cfg.ProjectName, ver)
	}

	return desc, nil
}
