# Authentication

## DefaultAzureCredential

`az-loadenv` authenticates using Azure's [DefaultAzureCredential](https://learn.microsoft.com/en-us/azure/developer/go/azure-sdk-authentication) chain, which tries multiple authentication methods in sequence until one succeeds.

### Credential chain

| Priority | Method | Use case |
|---|---|---|
| 1 | Environment variables | CI/CD pipelines, containers |
| 2 | Workload identity | Kubernetes pods, GitHub Actions |
| 3 | Managed identity | Azure VMs, App Service, Container Apps |
| 4 | Azure CLI | Local development |
| 5 | Azure Developer CLI | Local development |

The first method that succeeds is used. No configuration is needed — just make sure one of these methods is available.

## Local development

The simplest approach is to use the Azure CLI:

```bash
# Login interactively
az login

# Run az-loadenv
az-loadenv --app my-api --rg my-resource-group
```

Or with the Azure Developer CLI:

```bash
azd auth login
az-loadenv --app my-api --rg my-resource-group
```

## Service principal

For CI/CD and automated environments, use a service principal by setting environment variables:

```bash
export AZURE_TENANT_ID=xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
export AZURE_CLIENT_ID=xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
export AZURE_CLIENT_SECRET=your-client-secret

az-loadenv --app my-api --rg my-resource-group
```

The service principal needs:
- **Reader** role (or more) on the App Service resource to list settings
- **Key Vault Secrets User** role (or `Get` secret permission) on any referenced Key Vaults

## Managed identity

On Azure compute resources (VMs, App Service, Container Apps), managed identity is used automatically. No environment variables needed — just assign the appropriate roles to the identity.

## GitHub Actions

Use the [azure/login](https://github.com/marketplace/actions/azure-login) action with a service principal or federated credentials (OIDC):

```yaml
- uses: azure/login@v2
  with:
    client-id: ${{ secrets.AZURE_CLIENT_ID }}
    tenant-id: ${{ secrets.AZURE_TENANT_ID }}
    subscription-id: ${{ secrets.AZURE_SUBSCRIPTION_ID }}

- run: az-loadenv --app my-api --rg my-resource-group
```

## Required permissions

### App Service

The authenticating identity needs permission to read app settings. The minimum built-in role is **Website Contributor** on the App Service resource, or a custom role with the `Microsoft.Web/sites/config/list/action` permission.

### Key Vault

For settings that reference Key Vault secrets, the identity also needs access to the vault(s):

- **RBAC model:** assign the **Key Vault Secrets User** role
- **Access policy model:** grant **Get** permission on secrets
