# FancyVerteiler

With FancyVerteiler, you can deploy Minecraft and Hytale plugins to multiple platforms at once via GitHub Actions or the standalone app.

## Features

- Configure multiple platforms in a single JSON configuration file.
- Automatically read version and changelog from files.
- Send notifications to a Discord channel via webhook.

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
    hangar_api_key: ${{ secrets.HANGAR_API_KEY }}
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
- `<platform>_api_key` is only required if you want to publish to <platform>.

Example json config:
```json
{
  "project_name": "FancyNpcs",
  "plugin_jar_path": "./plugins/fancynpcs/build/libs/FancyNpcs-%VERSION%.jar",
  "changelog_path": "./plugins/fancynpcs/CHANGELOG.md",
  "version_path": "./plugins/fancynpcs/VERSION",
  "fancyspaces": {
    "space_id": "fn",
    "platform": "paper",
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
  "hangar": {
    "author": "peter",
    "project_id": "EeyAn23L",
    "supported_versions": [ "1.21.10", "1.21.11" ],
    "channel": "release"
  },
  "curseforge": {
    "project_id": "123456",
    "type": "plugin",
    "game_versions": [ "1.21.10", "1.21.11" ],
    "release_type": "release"
  },
  "orbis": {
    "resource_id": "1234",
    "channel": "RELEASE",
    "hytale_version_ids": [ "cmj1x42ef001k4qz9r03ojrpe" ]
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

To automatically get the last commit SHA and message from git, you can add the following steps before the FancyVerteiler step:
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

You can download the latest version of the standalone app from [FancySpaces](http://fancyspaces.net/spaces/fancyverteiler).