# FancyVerteiler

This action allows you to push version updates of Minecraft or Hytale plugins to multiple platforms at once.

## Features

Supported Minecraft plugin platforms:
- [FancySpaces](https://fancyspaces.net/)
- [Modrinth](https://modrinth.com/)
- [CurseForge](https://www.curseforge.com/)

Supported Hytale plugin platforms:
- [FancySpaces](https://fancyspaces.net/)
- [Orbis](https://orbis.place/)
- [Modtale](https://modtale.net/)
- [UnifiedHytale](https://www.unifiedhytale.com)
- [HytaHub](https://hytahub.com/)

Send notifications to a Discord channel via webhook.

## Usage

### GitHub Actions

Include the following in your GitHub Actions workflow:
```yml
- uses: fancyinnovations/fancyverteiler@main
  with:
    config_path: "plugins/fancynpcs/release_deployment_config.json"
    commit_sha: "see example for git integration below"
    commit_message: "see example for git integration below"
    fancyspaces_api_key: ${{ secrets.FANCYSPACES_API_KEY }}
    modrinth_api_key: ${{ secrets.MODRINTH_API_KEY }}
    curseforge_api_key: ${{ secrets.CURSEFORGE_API_KEY }}
    orbis_api_key: ${{ secrets.ORBIS_API_KEY }}
    modtale_api_key: ${{ secrets.MODTALE_API_KEY }}
    unifiedhytale_api_key: ${{ secrets.UNIFIEDHYTALE_API_KEY }}
    hytahub_api_key: ${{ secrets.HYTAHUB_API_KEY }}
    discord_webhook_url: ${{ secrets.DISCORD_WEBHOOK_URL }}
```

Inputs:
- `config_path` (required): Path to the JSON configuration file for FancyVerteiler.
- `commit_sha` (optional): The commit SHA to replace in the changelog.
- `commit_message` (optional): The commit message to replace in the changelog.
- `fancyspaces_api_key` is only required if you want to publish to FancySpaces.
- `modrinth_api_key` is only required if you want to publish to Modrinth.
- `curseforge_api_key` is only required if you want to publish to CurseForge.
- `orbis_api_key` is only required if you want to publish to Orbis.
- `modtale_api_key` is only required if you want to publish to Modtale.
- `unifiedhytale_api_key` is only required if you want to publish to UnifiedHytale.
- `hytahub_api_key` is only required if you want to publish to HytaHub.
- `discord_webhook_url` is only required if you want to send notifications to Discord.

Example config:
```json
{
  "project_name": "FancyNpcs",
  "plugin_jar_path": "../../../../plugins/fancynpcs/build/libs/FancyNpcs-%VERSION%.jar",
  "changelog_path": "../../../../plugins/fancynpcs/CHANGELOG.md",
  "version_path": "../../../../plugins/fancynpcs/VERSION",
  "fancyspaces": {
    "space_id": "fancyinnovations/fancynpcs",
    "platform": "minecraft_plugin",
    "channel": "release",
    "supported_versions": [ "1.21.10", "1.21.11" ]
  },
  "modrinth": {
    "project_id": "EeyAn23L",
    "supported_versions": [ "1.21.10", "1.21.11" ],
    "channel": "release",
    "loaders": [ "paper", "folia" ],
    "featured": true
  },
  "curseforge": {
    "project_id": "123456",
    "type": "plugin",
    "game_versions": [ "1.21.10", "1.21.11" ],
    "release_type": "release"
  },
  "orbis": {
    "resource_id": "1234",
    "is_pre_release": false,
    "hytale_version_ids": [ "2026.12.02.12312313" ]
  },
  "modtale": {
    "project_id": "abcdef123456",
    "channel": "RELEASE",
    "game_versions": [ "2026.12.02.12312313" ]
  },
  "unifiedhytale": {
    "project_id": "abcdef123456",
    "game_versions": [ "2026.12.02.12312313" ],
    "release_channel": "release"
  },
  "hytahub": {
    "slug": "mymode",
    "channel": "release"
  }
}
```

Full example with git integration:
```yml
      - name: Get last commit SHA and message
        id: last_commit
        run: |
          {
            echo "commit_sha=$(git rev-parse --short HEAD)"
            echo "commit_msg<<EOF"
            git log -1 --pretty=%B
            echo "EOF"
          } >> "$GITHUB_OUTPUT"

      - name: Deploy
        uses: fancyinnovations/fancyverteiler@main
        with:
          config_path: "/plugins/fancynpcs/release_deployment_config.json"
          commit_sha: ${{ steps.last_commit.outputs.commit_sha }}
          commit_message: ${{ steps.last_commit.outputs.commit_msg }}
          modrinth_api_key: ${{ secrets.MODRINTH_API_KEY }}
          curseforge_api_key: ${{ secrets.CURSEFORGE_API_KEY }}
          discord_webhook_url: ${{ secrets.DISCORD_WEBHOOK_URL }}

```

This will replace `%COMMIT_HASH%` and `%COMMIT_MESSAGE%` in the changelog with the actual commit hash and message.

### Standalone

You can also run FancyVerteiler as a standalone app.
Everything works the same way as in GitHub Actions, but you need to provide the inputs as environment variables.

Environment variables:
- `FV_CONFIG_PATH`
- `FV_DISCORD_WEBHOOK_URL`
- `FV_COMMIT_SHA`
- `FV_MESSAGE_SHA`
- `FV_{PLATFORM}_API_KEY` (example: `FV_FANCYSPACES_API_KEY`)