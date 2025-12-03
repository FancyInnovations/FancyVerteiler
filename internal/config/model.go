package config

import (
	"encoding/json"
	"os"
)

var BasePath = "."

type DeploymentConfig struct {
	ProjectName string `json:"project_name"`

	PluginJarPath string `json:"plugin_jar_path"`
	pluginJar     []byte

	ChangelogPath string `json:"changelog_path"`
	changelog     string

	VersionPath string `json:"version_path"`
	version     string

	Modrinth *Modrinth `json:"modrinth,omitempty"`
	Modtale  *Modtale  `json:"modtale,omitempty"`
}

type Modrinth struct {
	ProjectID         string   `json:"project_id"`
	SupportedVersions []string `json:"supported_versions"`
	Channel           string   `json:"channel"`
	Loaders           []string `json:"loaders"`
	Featured          bool     `json:"featured"`
}

type Modtale struct {
	ProjectID    string   `json:"project_id"`
	GameVersions []string `json:"game_versions"`
}

func (d *DeploymentConfig) PluginJar() ([]byte, error) {
	if d.pluginJar != nil {
		return d.pluginJar, nil
	}

	data, err := os.ReadFile(BasePath + "/" + d.PluginJarPath)
	if err != nil {
		return nil, err
	}

	d.pluginJar = data

	return data, nil
}

func (d *DeploymentConfig) Version() (string, error) {
	if d.version != "" {
		return d.version, nil
	}

	data, err := os.ReadFile(BasePath + "/" + d.VersionPath)
	if err != nil {
		return "", err
	}

	d.version = string(data)

	return string(data), nil
}

func (d *DeploymentConfig) Changelog() (string, error) {
	if d.changelog != "" {
		return d.changelog, nil
	}

	data, err := os.ReadFile(BasePath + "/" + d.ChangelogPath)
	if err != nil {
		return "", err
	}

	d.changelog = string(data)

	return string(data), nil
}

func ReadFromPath(path string) (*DeploymentConfig, error) {
	data, err := os.ReadFile(BasePath + "/" + path)
	if err != nil {
		return nil, err
	}

	var config DeploymentConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}
	return &config, nil
}
