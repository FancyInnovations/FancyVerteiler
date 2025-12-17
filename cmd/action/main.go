package main

import (
	"FancyVerteiler/internal/config"
	"FancyVerteiler/internal/curseforge"
	"FancyVerteiler/internal/discord"
	"FancyVerteiler/internal/fancyspaces"
	"FancyVerteiler/internal/git"
	"FancyVerteiler/internal/modrinth"
	"FancyVerteiler/internal/modtale"
	"FancyVerteiler/internal/orbis"

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
		apiKey := githubactions.GetInput("fancyspaces_api_key")
		if apiKey == "" {
			githubactions.Fatalf("missing input 'fancyspaces_api_key'")
		}

		githubactions.Infof("Deploying to FancySpaces space: %s", cfg.FancySpaces.SpaceID)

		fs := fancyspaces.New(apiKey, gs)
		if err := fs.Deploy(cfg); err != nil {
			githubactions.Fatalf("failed to deploy to FancySpaces: %v", err)
		}
		githubactions.Infof("Successfully deployed to FancySpaces space: %s", cfg.FancySpaces.SpaceID)
	}

	if cfg.Modrinth != nil {
		apiKey := githubactions.GetInput("modrinth_api_key")
		if apiKey == "" {
			githubactions.Fatalf("missing input 'modrinth_api_key'")
		}

		githubactions.Infof("Deploying to Modrinth project: %s", cfg.Modrinth.ProjectID)

		mr := modrinth.New(apiKey, gs)
		if err := mr.Deploy(cfg); err != nil {
			githubactions.Fatalf("failed to deploy to Modrinth: %v", err)
		}
		githubactions.Infof("Successfully deployed to Modrinth project: %s", cfg.Modrinth.ProjectID)
	}

	if cfg.Orbis != nil {
		apiKey := githubactions.GetInput("orbis_api_key")
		if apiKey == "" {
			githubactions.Fatalf("missing input 'orbis_api_key'")
		}

		githubactions.Infof("Deploying to Orbis resource: %s", cfg.Orbis.ResourceID)

		ob := orbis.New(apiKey, gs)
		if err := ob.Deploy(cfg); err != nil {
			githubactions.Fatalf("failed to deploy to Orbis: %v", err)
		}
		githubactions.Infof("Successfully deployed to Orbis resource: %s", cfg.Orbis.ResourceID)
	}

	if cfg.Modtale != nil {
		apiKey := githubactions.GetInput("modtale_api_key")
		if apiKey == "" {
			githubactions.Fatalf("missing input 'modtale_api_key'")
		}

		githubactions.Infof("Deploying to Modtale project: %s", cfg.Modtale.ProjectID)

		mt := modtale.New(apiKey, gs)
		if err := mt.Deploy(cfg); err != nil {
			githubactions.Fatalf("failed to deploy to Modtale: %v", err)
		}
		githubactions.Infof("Successfully deployed to Modtale project: %s", cfg.Modtale.ProjectID)
	}

	if cfg.CurseForge != nil {
		apiKey := githubactions.GetInput("curseforge_api_key")
		if apiKey == "" {
			githubactions.Fatalf("missing input 'curseforge_api_key'")
		}

		githubactions.Infof("Deploying to CurseForge project: %s", cfg.CurseForge.ProjectID)

		cf := curseforge.New(apiKey, gs)
		if err := cf.Deploy(cfg); err != nil {
			githubactions.Fatalf("failed to deploy to CurseForge: %v", err)
		}
		githubactions.Infof("Successfully deployed to CurseForge project: %s", cfg.CurseForge.ProjectID)
	}

	if discWebhookURL != "" {
		disc := discord.New()
		if err := disc.SendSuccessMessage(discWebhookURL, cfg); err != nil {
			githubactions.Fatalf("failed to send Discord success message: %v", err)
		}
		githubactions.Infof("Successfully sent Discord success message")
	}
}