package main

import (
	"FancyVerteiler/internal/config"
	"FancyVerteiler/internal/curseforge"
	"FancyVerteiler/internal/discord"
	"FancyVerteiler/internal/fancyspaces"
	"FancyVerteiler/internal/git"
	"FancyVerteiler/internal/hangar"
	"FancyVerteiler/internal/hytahub"
	"FancyVerteiler/internal/modrinth"
	"FancyVerteiler/internal/modtale"
	"FancyVerteiler/internal/orbis"
	"FancyVerteiler/internal/unifiedhytale"

	"github.com/sethvargo/go-githubactions"
)

func main() {
	configPath := githubactions.GetInput("config_path")
	if configPath == "" {
		githubactions.Fatalf("missing input 'config_path'")
	}

	discWebhookURL := githubactions.GetInput("discord_webhook_url")

	githubactions.Infof("Reading config: %s", configPath)

	config.BasePath = "/github/workspace"
	cfg, err := config.ReadFromPath(configPath)
	if err != nil {
		githubactions.Fatalf("failed to read config: %v", err)
	}

	githubactions.Infof("Successfully read config for project: %s", cfg.ProjectName)

	sha := githubactions.GetInput("commit_sha")
	if sha == "" {
		sha = "unknown"
	}
	message := githubactions.GetInput("commit_message")
	if message == "" {
		message = "unknown"
	}
	gs := git.New(sha, message)

	if cfg.FancySpaces != nil {
		deployToFancySpaces(cfg, gs)
	}
	if cfg.Modrinth != nil {
		deployToModrinth(cfg, gs)
	}
	if cfg.Hangar != nil {
		deployToHangar(cfg, gs)
	}
	if cfg.Orbis != nil {
		deployToOrbis(cfg, gs)
	}
	if cfg.Modtale != nil {
		deployToModtale(cfg, gs)
	}
	if cfg.CurseForge != nil {
		deployToCurseforge(cfg, gs)
	}
	if cfg.UnifiedHytale != nil {
		deployToUnifiedHytale(cfg, gs)
	}
	if cfg.Hytahub != nil {
		deployToHytahub(cfg, gs)
	}

	if discWebhookURL != "" {
		disc := discord.New()
		if err := disc.SendSuccessMessage(discWebhookURL, cfg); err != nil {
			githubactions.Errorf("Failed to send Discord success message: %v", err)
		} else {
			githubactions.Infof("Successfully sent Discord success message")
		}
	}
}

func deployToFancySpaces(cfg *config.DeploymentConfig, gs *git.Service) {
	apiKey := githubactions.GetInput("fancyspaces_api_key")
	if apiKey == "" {
		githubactions.Errorf("Missing input 'fancyspaces_api_key'")
		return
	}

	githubactions.Infof("Deploying to FancySpaces space: %s", cfg.FancySpaces.SpaceID)

	fs := fancyspaces.New(apiKey, gs)
	if err := fs.Deploy(cfg); err != nil {
		githubactions.Errorf("Failed to deploy to FancySpaces: %v", err)
		return
	}
	githubactions.Infof("Successfully deployed to FancySpaces space: %s", cfg.FancySpaces.SpaceID)
}

func deployToModrinth(cfg *config.DeploymentConfig, gs *git.Service) {
	apiKey := githubactions.GetInput("modrinth_api_key")
	if apiKey == "" {
		githubactions.Errorf("Missing input 'modrinth_api_key'")
		return
	}

	githubactions.Infof("Deploying to Modrinth project: %s", cfg.Modrinth.ProjectID)

	mr := modrinth.New(apiKey, gs)
	if err := mr.Deploy(cfg); err != nil {
		githubactions.Errorf("Failed to deploy to Modrinth: %v", err)
		return
	}
	githubactions.Infof("Successfully deployed to Modrinth project: %s", cfg.Modrinth.ProjectID)
}

func deployToHangar(cfg *config.DeploymentConfig, gs *git.Service) {
	apiKey := githubactions.GetInput("hangar_api_key")
	if apiKey == "" {
		githubactions.Errorf("Missing input 'hangar_api_key'")
		return
	}

	githubactions.Infof("Deploying to Hangar project: %s", cfg.Hangar.ProjectID)

	hg := hangar.New(apiKey, gs)
	if err := hg.Deploy(cfg); err != nil {
		githubactions.Errorf("Failed to deploy to Hangar: %v", err)
		return
	}
	githubactions.Infof("Successfully deployed to Hangar project: %s", cfg.Hangar.ProjectID)
}

func deployToOrbis(cfg *config.DeploymentConfig, gs *git.Service) {
	apiKey := githubactions.GetInput("orbis_api_key")
	if apiKey == "" {
		githubactions.Errorf("Missing input 'orbis_api_key'")
		return
	}

	githubactions.Infof("Deploying to Orbis resource: %s", cfg.Orbis.ResourceID)

	ob := orbis.New(apiKey, gs)
	if err := ob.Deploy(cfg); err != nil {
		githubactions.Errorf("Failed to deploy to Orbis: %v", err)
		return
	}
	githubactions.Infof("Successfully deployed to Orbis resource: %s", cfg.Orbis.ResourceID)
}

func deployToModtale(cfg *config.DeploymentConfig, gs *git.Service) {
	apiKey := githubactions.GetInput("modtale_api_key")
	if apiKey == "" {
		githubactions.Errorf("Missing input 'modtale_api_key'")
		return
	}

	githubactions.Infof("Deploying to Modtale project: %s", cfg.Modtale.ProjectID)

	mt := modtale.New(apiKey, gs)
	if err := mt.Deploy(cfg); err != nil {
		githubactions.Errorf("Failed to deploy to Modtale: %v", err)
		return
	}
	githubactions.Infof("Successfully deployed to Modtale project: %s", cfg.Modtale.ProjectID)
}

func deployToCurseforge(cfg *config.DeploymentConfig, gs *git.Service) {
	apiKey := githubactions.GetInput("curseforge_api_key")
	if apiKey == "" {
		githubactions.Errorf("Missing input 'curseforge_api_key'")
		return
	}

	githubactions.Infof("Deploying to CurseForge project: %s", cfg.CurseForge.ProjectID)

	cf := curseforge.New(apiKey, gs)
	if err := cf.Deploy(cfg); err != nil {
		githubactions.Errorf("Failed to deploy to CurseForge: %v", err)
		return
	}
	githubactions.Infof("Successfully deployed to CurseForge project: %s", cfg.CurseForge.ProjectID)
}

func deployToUnifiedHytale(cfg *config.DeploymentConfig, gs *git.Service) {
	apiKey := githubactions.GetInput("unifiedhytale_api_key")
	if apiKey == "" {
		githubactions.Errorf("Missing input 'unifiedhytale_api_key'")
		return
	}

	githubactions.Infof("Deploying to Modtale project: %s", cfg.UnifiedHytale.ProjectID)

	mt := unifiedhytale.New(apiKey, gs)
	if err := mt.Deploy(cfg); err != nil {
		githubactions.Errorf("Failed to deploy to UnifiedHytale: %v", err)
		return
	}
	githubactions.Infof("Successfully deployed to UnifiedHytale project: %s", cfg.UnifiedHytale.ProjectID)
}

func deployToHytahub(cfg *config.DeploymentConfig, gs *git.Service) {
	apiKey := githubactions.GetInput("hytahub_api_key")
	if apiKey == "" {
		githubactions.Errorf("Missing input 'hytahub_api_key'")
		return
	}

	githubactions.Infof("Deploying to Hytahub project: %s", cfg.Hytahub.Slug)

	mt := hytahub.New(apiKey, gs)
	if err := mt.Deploy(cfg); err != nil {
		githubactions.Errorf("Failed to deploy to Hytahub: %v", err)
		return
	}
	githubactions.Infof("Successfully deployed to Hytahub project: %s", cfg.Hytahub.Slug)
}
