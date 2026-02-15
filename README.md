# az-loadenv

Export Azure Web App settings to a `.env` file with automatic Key Vault secret resolution.

`az-loadenv` fetches application settings from an Azure App Service web app, resolves any `@Microsoft.KeyVault(SecretUri=...)` references to their actual secret values, and writes everything to a `.env` file ready for local development.

## Features

- **Key Vault resolution** — automatically detects and resolves Key Vault references to plaintext values
- **Concurrent secret fetching** — resolves up to 10 Key Vault secrets in parallel
- **Atomic writes** — output is written to a temp file then renamed, so the `.env` file is never left in a partial state
- **Secure defaults** — output file is created with `0600` permissions (owner read/write only)
- **Smart quoting** — values with spaces, quotes, newlines, or shell metacharacters are automatically double-quoted and escaped
- **Deterministic output** — settings are sorted alphabetically by key for clean diffs
- **Auto subscription detection** — picks up the active subscription from `AZURE_SUBSCRIPTION_ID` or `~/.azure/azureProfile.json`
- **Cross-platform** — pre-built binaries for Linux, macOS, and Windows

## Quick install

**macOS / Linux:**
```bash
curl -fsSL https://curiousdev.github.io/az-loadenv/install.sh | bash
```

**Windows (PowerShell):**
```powershell
irm https://curiousdev.github.io/az-loadenv/install.ps1 | iex
```

## Installation

### Download a release

Pre-built binaries are available on the [Releases](https://github.com/curiousdev/az-loadenv/releases) page for:

| Platform | Architecture | Archive |
|---|---|---|
| Linux | x86_64 | `az-loadenv-linux-amd64.tar.gz` |
| Linux | ARM64 | `az-loadenv-linux-arm64.tar.gz` |
| macOS | Intel | `az-loadenv-darwin-amd64.tar.gz` |
| macOS | Apple Silicon | `az-loadenv-darwin-arm64.tar.gz` |
| Windows | x86_64 | `az-loadenv-windows-amd64.zip` |

```bash
# Example: macOS Apple Silicon
curl -L https://github.com/curiousdev/az-loadenv/releases/latest/download/az-loadenv-darwin-arm64.tar.gz | tar xz
sudo mv az-loadenv /usr/local/bin/
```

### Build from source

Requires Go 1.25+:

```bash
go install github.com/curiousdev/az-loadenv@latest
```

## Usage

```
az-loadenv --app <name> --rg <resource-group> [flags]
```

### Flags

| Flag | Description | Default |
|---|---|---|
| `--app` | Azure Web App name | *(required)* |
| `--rg` | Resource group name | *(required)* |
| `-o` | Output file path | `.env` |
| `--raw` | Write values without quoting or escaping | `false` |
| `--version` | Print version and exit | |

### Examples

```bash
# Write settings to .env (default)
az-loadenv --app my-api --rg my-resource-group

# Write settings to a custom file
az-loadenv --app my-api --rg my-resource-group -o .env.local

# Use with a specific subscription
AZURE_SUBSCRIPTION_ID=xxx az-loadenv --app my-api --rg my-resource-group

# Use with a service principal (CI/CD)
export AZURE_TENANT_ID=xxx AZURE_CLIENT_ID=xxx AZURE_CLIENT_SECRET=xxx
az-loadenv --app my-api --rg my-resource-group
```

### Output format

Settings are written as `KEY=VALUE`, one per line, sorted alphabetically:

```
API_URL=https://api.example.com
DB_CONNECTION="Server=db.example.com;Password=s3cret"
SIMPLE_FLAG=true
```

Values containing spaces, quotes, newlines, or other special characters are automatically double-quoted and escaped.

## Authentication

`az-loadenv` uses Azure's [DefaultAzureCredential](https://learn.microsoft.com/en-us/azure/developer/go/azure-sdk-authentication), which tries the following methods in order:

| Priority | Method | When to use |
|---|---|---|
| 1 | Environment variables | CI/CD pipelines, containers |
| 2 | Workload identity | Kubernetes, GitHub Actions |
| 3 | Managed identity | Azure VMs, App Service, Container Apps |
| 4 | Azure CLI | Local development (`az login`) |
| 5 | Azure Developer CLI | Local development (`azd auth login`) |

For local development, the simplest path is:

```bash
az login
az-loadenv --app my-api --rg my-resource-group
```

## Subscription detection

The Azure subscription is resolved automatically:

1. **`AZURE_SUBSCRIPTION_ID`** environment variable, if set
2. **Default subscription** from `~/.azure/azureProfile.json` (set by `az login` / `az account set`)

## Environment variables

| Variable | Purpose |
|---|---|
| `AZURE_SUBSCRIPTION_ID` | Override automatic subscription detection |
| `AZURE_CLIENT_ID` | Service principal authentication |
| `AZURE_CLIENT_SECRET` | Service principal authentication |
| `AZURE_TENANT_ID` | Service principal authentication |

## Key Vault resolution

Any app setting whose value matches the Azure Key Vault reference format is automatically resolved:

```
@Microsoft.KeyVault(SecretUri=https://my-vault.vault.azure.net/secrets/my-secret)
@Microsoft.KeyVault(SecretUri=https://my-vault.vault.azure.net/secrets/my-secret/version-id)
```

The authenticating identity must have `Get` permission on secrets in the referenced vault(s). Both versioned and unversioned secret URIs are supported, and secrets can span multiple vaults.

If a secret fails to resolve, `az-loadenv` logs the error to stderr and writes the original Key Vault reference to the output file so other settings are not blocked.

## Security considerations

- The output `.env` file is created with **`0600` permissions** (owner read/write only)
- Writes are **atomic** (temp file + rename) to prevent partial reads
- Secret values are never printed to stderr — only setting names are logged
- Add `.env` to your `.gitignore` to avoid committing secrets

## Building

```bash
# Development build
go build -o az-loadenv .

# Production build with version info
go build -trimpath -ldflags="-s -w -X main.version=1.0.0 -X main.build=1" -o az-loadenv .
```

## License

MIT
