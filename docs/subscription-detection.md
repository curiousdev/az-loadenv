# Subscription detection

## How it works

`az-loadenv` determines which Azure subscription to use through a two-step process:

### 1. Environment variable (highest priority)

If `AZURE_SUBSCRIPTION_ID` is set, that value is used directly:

```bash
AZURE_SUBSCRIPTION_ID=xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx \
  az-loadenv --app my-api --rg my-resource-group
```

This is the recommended approach for CI/CD pipelines and automated environments.

### 2. Azure CLI profile (fallback)

If no environment variable is set, `az-loadenv` reads `~/.azure/azureProfile.json` — the profile file managed by the Azure CLI — and uses whichever subscription is marked as the default.

This file is created and updated by `az login` and `az account set`:

```bash
# Login
az login

# Set the active subscription (if you have multiple)
az account set --subscription "My Subscription"

# Now az-loadenv uses that subscription
az-loadenv --app my-api --rg my-resource-group
```

## Common scenarios

### Multiple subscriptions

If your Azure account has access to multiple subscriptions, make sure the correct one is active:

```bash
# List subscriptions
az account list --output table

# Set the desired one
az account set --subscription "Production"
```

Or override per-command with the environment variable:

```bash
AZURE_SUBSCRIPTION_ID=xxxxxxxx az-loadenv --app my-api --rg my-resource-group
```

### CI/CD pipelines

Always set `AZURE_SUBSCRIPTION_ID` explicitly in pipelines to avoid depending on profile state:

```yaml
env:
  AZURE_SUBSCRIPTION_ID: ${{ secrets.AZURE_SUBSCRIPTION_ID }}
```

## Error messages

| Error | Cause | Fix |
|---|---|---|
| `cannot read ~/.azure/azureProfile.json` | Azure CLI not installed or never logged in | Run `az login` |
| `no default subscription found` | Profile exists but no subscription is marked default | Run `az account set --subscription "..."` |
