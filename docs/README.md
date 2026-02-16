# az-loadenv

> Export Azure Web App settings to a `.env` file with automatic Key Vault secret resolution.

## What is az-loadenv?

`az-loadenv` is a command-line tool that pulls application settings from an Azure App Service web app and writes them to a local `.env` file. It automatically detects and resolves `@Microsoft.KeyVault(SecretUri=...)` references, replacing them with the actual secret values so your local development environment matches production configuration.

## Why?

Azure App Service stores configuration as app settings. During local development, you typically need those same values in a `.env` file. Manually copying settings is tedious and error-prone, especially when some values are Key Vault references that need to be resolved to their actual secrets.

`az-loadenv` automates this in a single command:

```bash
az-loadenv --app my-api --rg my-resource-group
```

## Key features

- **Automatic Key Vault resolution** — detects `@Microsoft.KeyVault(SecretUri=...)` references and resolves them to plaintext values
- **Concurrent secret fetching** — resolves up to 10 secrets in parallel for fast execution
- **Atomic file writes** — uses temp file + rename so your `.env` is never left in a partial state
- **Secure by default** — output file gets `0600` permissions (owner read/write only)
- **Smart value formatting** — automatically quotes and escapes special characters, or use `--raw` to skip escaping
- **Deterministic output** — settings sorted alphabetically for clean, diffable output
- **Zero configuration auth** — uses Azure's DefaultAzureCredential chain, so it works with `az login`, service principals, managed identity, and more
- **Cross-platform** — pre-built binaries for Linux, macOS (Intel + Apple Silicon), and Windows

## Quick install

**Homebrew (macOS / Linux):**
```bash
brew install curiousdev/tap/az-loadenv-cli
```

**WinGet (Windows):**
```powershell
winget install curiousdev.az-loadenv-cli
```

**Shell script (macOS / Linux):**
```bash
curl -fsSL https://curiousdev.github.io/az-loadenv/install.sh | bash
```

**PowerShell script (Windows):**
```powershell
irm https://curiousdev.github.io/az-loadenv/install.ps1 | iex
```

## Quick example

```bash
# Login to Azure
az login

# Export settings to .env
az-loadenv --app my-api --rg my-resource-group

# Output:
#   APP_NAME
#   DATABASE_URL (secret)
#   REDIS_HOST
#   STORAGE_KEY (secret)
#
# Wrote 4 settings (2 config, 2 secrets) to .env
```

The resulting `.env` file:

```
APP_NAME=my-api
DATABASE_URL="Server=db.example.com;Password=s3cret value"
REDIS_HOST=redis.example.com
STORAGE_KEY=base64encodedkey==
```
