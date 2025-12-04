package main

import (
	"FancyVerteiler/internal/config"
	"FancyVerteiler/internal/discord"
	"FancyVerteiler/internal/git"
	"FancyVerteiler/internal/modrinth"
	"FancyVerteiler/internal/modtale"

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

	gs := git.New()

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

	if discWebhookURL != "" {
		disc := discord.New()
		if err := disc.SendSuccessMessage(discWebhookURL, cfg); err != nil {
			githubactions.Fatalf("failed to send Discord success message: %v", err)
		}
		githubactions.Infof("Successfully sent Discord success message")
	}
}
