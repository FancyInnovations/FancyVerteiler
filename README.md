# FancyVerteiler

This action allows you to push version updates of Minecraft or Hytale plugins to multiple platforms at once.

## Usage

Include the following in your GitHub Actions workflow:
```yml
- uses: fancyinnovations/fancyverteiler@main
  with:
    config_path: "plugins/fancynpcs/release_deployment_config.json"
    modrinth_api_key: ${{ secrets.MODRINTH_API_KEY }}
    modtale_api_key: ${{ secrets.MODTALE_API_KEY }}
    discord_webhook_url: ${{ secrets.DISCORD_WEBHOOK_URL }}
```

Inputs:
- `config_path` (required): Path to the JSON configuration file for FancyVerteiler.
- `modrinth_api_key` is only required if you want to publish to Modrinth.
- `modtale_api_key` is only required if you want to publish to Modtale.
- `discord_webhook_url` is only required if you want to send notifications to Discord.

Example config:
```json
{
  "project_name": "FancyNpcs",
  "plugin_jar_path": "../../../../plugins/fancynpcs/build/libs/FancyNpcs-%VERSION%.jar",
  "changelog_path": "../../../../plugins/fancynpcs/CHANGELOG.md",
  "version_path": "../../../../plugins/fancynpcs/VERSION",
  "modrinth": {
    "project_id": "EeyAn23L",
    "supported_versions": [ "1.21.10", "1.21.11" ],
    "channel": "release",
    "loaders": [ "paper", "folia" ],
    "featured": true
  },
  "modtale": {
    "project_id": "abcdef123456",
    "game_versions": [ "2026.12.02.12312313" ]
  }
}
```