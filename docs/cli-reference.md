# CLI reference

## Synopsis

```
az-loadenv --app <name> --rg <resource-group> [flags]
```

## Flags

| Flag | Type | Required | Default | Description |
|---|---|---|---|---|
| `--app` | string | yes | | Azure Web App name |
| `--rg` | string | yes | | Resource group containing the web app |
| `-o` | string | no | `.env` | Output file path |
| `--raw` | bool | no | `false` | Write values without quoting or escaping |
| `--version` | bool | no | | Print version and exit |

## Exit codes

| Code | Meaning |
|---|---|
| `0` | Success |
| `1` | Error (authentication failure, missing arguments, API error, write failure) |

## Behavior

1. **Subscription detection** — resolves the Azure subscription from `AZURE_SUBSCRIPTION_ID` or `~/.azure/azureProfile.json`
2. **Authentication** — obtains credentials via Azure's DefaultAzureCredential chain
3. **Fetch settings** — calls the Azure App Service API to list all application settings
4. **Resolve secrets** — any value matching `@Microsoft.KeyVault(SecretUri=...)` is resolved concurrently (up to 10 in parallel)
5. **Write output** — settings are sorted alphabetically and written atomically to the output file with `0600` permissions

## Stderr output

All progress and diagnostic messages go to stderr, so the tool works cleanly in pipelines. The output file path and a summary are printed after completion:

```
  APP_NAME
  DATABASE_URL (secret)
  LOG_LEVEL
  REDIS_HOST (secret: error)

Wrote 4 settings (1 config, 3 secrets) to .env
```

Settings that fail to resolve are logged with `(secret: error)` and their original Key Vault reference value is written to the output file.

## Timeout

All operations run under a 2-minute timeout. If Azure API calls or Key Vault resolution take longer than this, the operation is cancelled. The tool also responds to `SIGINT` (Ctrl+C) for clean cancellation.
