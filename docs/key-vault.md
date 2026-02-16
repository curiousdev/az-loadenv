# Key Vault resolution

## Overview

Azure App Service supports [Key Vault references](https://learn.microsoft.com/en-us/azure/app-service/app-service-key-vault-references) in application settings. Instead of storing a secret value directly, a setting can reference a Key Vault secret:

```
@Microsoft.KeyVault(SecretUri=https://my-vault.vault.azure.net/secrets/my-secret)
```

When `az-loadenv` encounters a setting with this format, it automatically fetches the actual secret value from Key Vault and writes the resolved plaintext to the `.env` file.

## Supported formats

Both versioned and unversioned secret URIs are supported:

```
# Unversioned (latest version)
@Microsoft.KeyVault(SecretUri=https://my-vault.vault.azure.net/secrets/my-secret)

# Versioned (specific version)
@Microsoft.KeyVault(SecretUri=https://my-vault.vault.azure.net/secrets/my-secret/abc123def456)
```

## Multiple vaults

Settings can reference secrets from different Key Vaults. `az-loadenv` creates a client per vault and resolves all references regardless of which vault they point to.

```
# These settings can coexist — they reference different vaults
DB_PASSWORD=@Microsoft.KeyVault(SecretUri=https://prod-vault.vault.azure.net/secrets/db-pass)
API_KEY=@Microsoft.KeyVault(SecretUri=https://shared-vault.vault.azure.net/secrets/api-key)
```

## Concurrent resolution

Key Vault secrets are resolved concurrently with up to 10 parallel requests. This significantly reduces total execution time when an app has many secret references. Vault clients are reused across secrets from the same vault.

## Error handling

If a secret fails to resolve (permission denied, secret not found, network error), `az-loadenv`:

1. Logs the error to stderr: `SETTING_NAME (secret: error)`
2. Writes the original `@Microsoft.KeyVault(SecretUri=...)` value to the output file
3. Continues resolving other secrets — one failure does not block the rest

This partial-success approach ensures you still get a usable `.env` file even if some secrets are inaccessible.

## Required permissions

The identity used by `az-loadenv` must have permission to read secrets from the referenced vault(s):

### RBAC authorization (recommended)

Assign the **Key Vault Secrets User** role on the vault:

```bash
az role assignment create \
  --role "Key Vault Secrets User" \
  --assignee <principal-id> \
  --scope /subscriptions/<sub>/resourceGroups/<rg>/providers/Microsoft.KeyVault/vaults/<vault>
```

### Access policy authorization

If the vault uses access policies instead of RBAC, grant **Get** permission on secrets:

```bash
az keyvault set-policy \
  --name my-vault \
  --object-id <principal-id> \
  --secret-permissions get
```

## Debugging resolution failures

If secrets are not resolving, check:

1. **Identity has vault access** — verify with `az keyvault secret show --vault-name my-vault --name my-secret`
2. **Vault firewall** — if the vault has network restrictions, your client IP must be allowed
3. **Secret exists** — the referenced secret (and version, if specified) must exist in the vault
4. **DNS resolution** — the vault hostname (e.g., `my-vault.vault.azure.net`) must be resolvable from your environment
