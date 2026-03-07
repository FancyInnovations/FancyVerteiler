package discord

import (
	"FancyVerteiler/internal/config"
	"FancyVerteiler/internal/git"
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
	hc  *http.Client
	git *git.Service
}

func New(git *git.Service) *Service {
	return &Service{
		hc:  &http.Client{},
		git: git,
	}
}

func (s *Service) SendSuccessMessage(webhookURL string, cfg *config.DeploymentConfig) error {
	desc, err := s.buildDescription(cfg)
	if err != nil {
		return err
	}

	ver, err := cfg.Version()
	if err != nil {
		return err
	}

	msg := Message{
		Content: "New version of " + cfg.ProjectName + " published!",
		Embeds: []Embed{
			{
				Title:       fmt.Sprintf("%s v%s published!", cfg.ProjectName, ver),
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

func (s *Service) buildDescription(cfg *config.DeploymentConfig) (string, error) {
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

		desc += fmt.Sprintf("\n**Channel:** %s", strings.ToUpper(channel))
	}

	desc += fmt.Sprintf("\n**Commit ([%s](%s)):**", s.git.CommitSHA(), s.git.CommitURL())
	desc += fmt.Sprintf("\n```\n%s\n```", s.git.CommitMessage())

	desc += "\n"
	desc += "\n**Download Links:**"

	if cfg.FancySpaces != nil {
		fileName := filepath.Base(cfg.PluginJarPath)
		fileName = strings.ReplaceAll(fileName, "%VERSION%", ver)
		desc += fmt.Sprintf("\n- [FancySpaces](https://fancyspaces.net/spaces/%s/versions/%s)", cfg.FancySpaces.SpaceID, ver)
	}

	if cfg.Modrinth != nil {
		desc += fmt.Sprintf("\n- [Modrinth](https://modrinth.com/plugin/%s/version/%s)", cfg.ProjectName, ver)
	}

	if cfg.Hangar != nil {
		desc += fmt.Sprintf("\n- [Hangar](https://hangar.papermc.io/%s/%s/versions/%s)", cfg.Hangar.Author, cfg.Hangar.ProjectID, ver)
	}

	if cfg.CurseForge != nil {
		desc += fmt.Sprintf("\n- [CurseForge](https://www.curseforge.com/minecraft/bukkit-plugins/%s/files/all)", cfg.ProjectName)
	}

	return desc, nil
}
