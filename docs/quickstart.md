# Quick start

## Prerequisites

- An Azure account with access to an App Service web app
- The Azure CLI installed and logged in, **or** service principal credentials

## Step 1: Authenticate

The simplest way is through the Azure CLI:

```bash
az login
```

If your web app's settings reference Key Vault secrets, your identity also needs `Get` permission on those secrets in the relevant vault(s).

## Step 2: Export settings

```bash
az-loadenv --app my-api --rg my-resource-group
```

This creates a `.env` file in the current directory with all the app's settings. Key Vault references are resolved automatically.

You'll see progress on stderr as each setting is processed:

```
  APP_NAME
  LOG_LEVEL
  DATABASE_URL (secret)
  REDIS_CONNECTION (secret)
  STORAGE_ACCOUNT

Wrote 5 settings (3 config, 2 secrets) to .env
```

## Step 3: Use the .env file

Most frameworks and tools load `.env` files automatically. For example:

```bash
# Node.js (with dotenv)
node -r dotenv/config app.js

# Python (with python-dotenv)
python app.py

# Docker
docker run --env-file .env myimage

# .NET (with DotNetEnv)
dotnet run
```

## Common options

```bash
# Write to a different file
az-loadenv --app my-api --rg my-resource-group -o .env.local

# Target a specific subscription
AZURE_SUBSCRIPTION_ID=xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx \
  az-loadenv --app my-api --rg my-resource-group

# Skip quoting/escaping for SDKs that handle raw values
az-loadenv --app my-api --rg my-resource-group --raw
```
