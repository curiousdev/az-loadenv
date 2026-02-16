# Troubleshooting

## Authentication errors

### "Failed to authenticate"

`az-loadenv` could not obtain Azure credentials through any method in the DefaultAzureCredential chain.

**Fix:**
- Run `az login` for local development
- Set `AZURE_TENANT_ID`, `AZURE_CLIENT_ID`, and `AZURE_CLIENT_SECRET` for service principal auth
- Verify credentials with: `az account show`

### "Failed to detect subscription"

No subscription could be determined.

**Fix:**
- Set `AZURE_SUBSCRIPTION_ID` environment variable, or
- Run `az login` followed by `az account set --subscription "..."` to set a default

### "cannot read ~/.azure/azureProfile.json"

The Azure CLI profile file doesn't exist.

**Fix:** Run `az login` at least once to create the profile.

### "no default subscription found"

The profile exists but no subscription is marked as default. This can happen with multiple subscriptions.

**Fix:** Run `az account set --subscription "My Subscription"`

## App Service errors

### "Failed to list app settings"

The API call to fetch app settings failed.

**Common causes:**
- The `--app` or `--rg` values are incorrect
- The identity doesn't have permission to read the web app's settings
- The web app doesn't exist in the specified resource group

**Fix:**
- Verify the app exists: `az webapp show --name my-api --resource-group my-rg`
- Check permissions: the identity needs `Microsoft.Web/sites/config/list/action`

### "No app settings found"

The web app exists but has no application settings configured.

This is not an error — the tool exits cleanly with no output file.

## Key Vault errors

### "secret: error" in output

One or more Key Vault references could not be resolved. The error is logged to stderr and the original reference is written to the `.env` file.

**Common causes:**
- Missing vault access (no `Get` permission on secrets)
- Vault firewall blocking your IP
- Secret or version doesn't exist
- Network connectivity issues

**Debugging steps:**

```bash
# Test vault access directly
az keyvault secret show --vault-name my-vault --name my-secret

# Check your IP against vault firewall rules
az keyvault network-rule list --name my-vault

# Verify the secret URI from the app setting
az webapp config appsettings list --name my-api --resource-group my-rg
```

### Resolution is slow

Key Vault resolution uses up to 10 concurrent requests. If resolution is slow:
- Check network latency to `*.vault.azure.net`
- Vault throttling may be occurring — Azure Key Vault has [rate limits](https://learn.microsoft.com/en-us/azure/key-vault/general/service-limits) per vault

## Output file issues

### Permission denied writing .env

The current user doesn't have write permission in the output directory.

**Fix:** Specify a writable path with `-o`:
```bash
az-loadenv --app my-api --rg my-resource-group -o /tmp/.env
```

### Values look wrong with my .env parser

Some parsers don't handle quoted or escaped values correctly. Try raw mode:

```bash
az-loadenv --app my-api --rg my-resource-group --raw
```

See [Raw mode](raw-mode.md) for details.

## Timeout

All operations run under a 2-minute timeout. If you consistently hit this limit, check:
- Network connectivity to Azure APIs
- Number of Key Vault references (many secrets across multiple vaults can be slow)
- Azure service health at https://status.azure.com
