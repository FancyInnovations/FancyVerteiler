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
	"log/slog"
	"os"

	"github.com/OliverSchlueter/goutils/env"
	"github.com/OliverSchlueter/goutils/sloki"
)

const (
	configPathEnv        = "FV_CONFIG_PATH" // required
	discordWebhookUrlEnv = "FV_DISCORD_WEBHOOK_URL"
	githubRepoURLEnv     = "FV_GITHUB_REPO_URL"
	commitShaEnv         = "FV_COMMIT_SHA"
	commitMessageEnv     = "FV_MESSAGE_SHA"

	fancyspacesApiKeyEnv   = "FV_FANCYSPACES_API_KEY"
	modrinthApiKeyEnv      = "FV_MODRINTH_API_KEY"
	hangarApiKeyEnv        = "FV_HANGAR_API_KEY"
	orbisApiKeyEnv         = "FV_ORBIS_API_KEY"
	modtaleApiKeyEnv       = "FV_MODTALE_API_KEY"
	curseforgeApiKeyEnv    = "FV_CURSEFORGE_API_KEY"
	unifiedhytaleApiKeyEnv = "FV_UNIFIEDHytale_API_KEY"
	hytahubApiKeyEnv       = "FV_HYTAHUB_API_KEY"
)

func main() {
	configPath := env.MustGetStr(configPathEnv)

	discWebhookURL := os.Getenv(discordWebhookUrlEnv)

	slog.Info("Reading config", slog.String("path", configPath))

	cfg, err := config.ReadFromPath(configPath)
	if err != nil {
		slog.Error("Failed to read config", sloki.WrapError(err))
		return
	}

	slog.Info("Successfully read config", slog.String("project", cfg.ProjectName))

	githubRepoURL := env.MustGetStr(githubRepoURLEnv)
	sha := env.MustGetStr(commitShaEnv)
	message := env.MustGetStr(commitMessageEnv)
	gs := git.New(githubRepoURL, sha, message)

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
		disc := discord.New(gs)
		if err := disc.SendSuccessMessage(discWebhookURL, cfg); err != nil {
			slog.Error("Failed to send Discord success message", sloki.WrapError(err))
		} else {
			slog.Info("Successfully sent Discord success message")
		}
	}
}

func deployToFancySpaces(cfg *config.DeploymentConfig, gs *git.Service) {
	apiKey := env.MustGetStr(fancyspacesApiKeyEnv)

	slog.Info("Deploying to FancySpaces space", slog.String("space_id", cfg.FancySpaces.SpaceID))

	fs := fancyspaces.New(apiKey, gs)
	if err := fs.Deploy(cfg); err != nil {
		slog.Error("Failed to deploy to FancySpaces", sloki.WrapError(err))
		return
	}
	slog.Info("Successfully deployed to FancySpaces", slog.String("space_id", cfg.FancySpaces.SpaceID))
}

func deployToModrinth(cfg *config.DeploymentConfig, gs *git.Service) {
	apiKey := env.MustGetStr(modrinthApiKeyEnv)

	slog.Info("Deploying to Modrinth project", slog.String("project_id", cfg.Modrinth.ProjectID))

	mr := modrinth.New(apiKey, gs)
	if err := mr.Deploy(cfg); err != nil {
		slog.Error("Failed to deploy to Modrinth", sloki.WrapError(err))
		return
	}
	slog.Info("Successfully deployed to Modrinth", slog.String("project_id", cfg.Modrinth.ProjectID))
}

func deployToHangar(cfg *config.DeploymentConfig, gs *git.Service) {
	apiKey := env.MustGetStr(hangarApiKeyEnv)

	slog.Info("Deploying to Hangar project", slog.String("project_id", cfg.Hangar.ProjectID))

	hn := hangar.New(apiKey, gs)
	if err := hn.Deploy(cfg); err != nil {
		slog.Error("Failed to deploy to Hangar", sloki.WrapError(err))
		return
	}
	slog.Info("Successfully deployed to Hangar", slog.String("project_id", cfg.Hangar.ProjectID))
}

func deployToOrbis(cfg *config.DeploymentConfig, gs *git.Service) {
	apiKey := env.MustGetStr(orbisApiKeyEnv)

	slog.Info("Deploying to Orbis resource", slog.String("resource_id", cfg.Orbis.ResourceID))

	ob := orbis.New(apiKey, gs)
	if err := ob.Deploy(cfg); err != nil {
		slog.Error("Failed to deploy to Orbis", sloki.WrapError(err))
		return
	}
	slog.Info("Successfully deployed to Orbis", slog.String("resource_id", cfg.Orbis.ResourceID))
}

func deployToModtale(cfg *config.DeploymentConfig, gs *git.Service) {
	apiKey := env.MustGetStr(modtaleApiKeyEnv)

	slog.Info("Deploying to Modtale project", slog.String("project_id", cfg.Modtale.ProjectID))

	mt := modtale.New(apiKey, gs)
	if err := mt.Deploy(cfg); err != nil {
		slog.Error("Failed to deploy to Modtale", sloki.WrapError(err))
		return
	}
	slog.Info("Successfully deployed to Modtale", slog.String("project_id", cfg.Modtale.ProjectID))
}

func deployToCurseforge(cfg *config.DeploymentConfig, gs *git.Service) {
	apiKey := env.MustGetStr(curseforgeApiKeyEnv)

	slog.Info("Deploying to CurseForge project", slog.String("project_id", cfg.CurseForge.ProjectID))

	cf := curseforge.New(apiKey, gs)
	if err := cf.Deploy(cfg); err != nil {
		slog.Error("Failed to deploy to CurseForge", sloki.WrapError(err))
		return
	}
	slog.Info("Successfully deployed to CurseForge", slog.String("project_id", cfg.CurseForge.ProjectID))
}

func deployToUnifiedHytale(cfg *config.DeploymentConfig, gs *git.Service) {
	apiKey := env.MustGetStr(unifiedhytaleApiKeyEnv)

	slog.Info("Deploying to UnifiedHytale project", slog.String("project_id", cfg.UnifiedHytale.ProjectID))

	mt := unifiedhytale.New(apiKey, gs)
	if err := mt.Deploy(cfg); err != nil {
		slog.Error("Failed to deploy to UnifiedHytale", sloki.WrapError(err))
		return
	}
	slog.Info("Successfully deployed to UnifiedHytale", slog.String("project_id", cfg.UnifiedHytale.ProjectID))
}

func deployToHytahub(cfg *config.DeploymentConfig, gs *git.Service) {
	apiKey := env.MustGetStr(hytahubApiKeyEnv)

	slog.Info("Deploying to Hytahub channel", slog.String("slug", cfg.Hytahub.Slug))

	ht := hytahub.New(apiKey, gs)
	if err := ht.Deploy(cfg); err != nil {
		slog.Error("Failed to deploy to Hytahub", sloki.WrapError(err))
		return
	}
	slog.Info("Successfully deployed to Hytahub", slog.String("slug", cfg.Hytahub.Slug))
}
