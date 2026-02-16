package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"net/url"
	"os"
	"os/signal"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/appservice/armappservice"
	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azsecrets"
	"golang.org/x/sync/errgroup"
)

var (
	version = "dev"
	build   = "0"
)

var kvRefPattern = regexp.MustCompile(`^@Microsoft\.KeyVault\(SecretUri=(https://.+)\)$`)

type setting struct {
	name   string
	value  string
	secret bool
	errMsg string
}

const helpText = `az-loadenv — Export Azure Web App settings to a .env file

Fetches application settings from an Azure App Service web app and writes
them to a .env file. Any settings that reference Azure Key Vault secrets
(@Microsoft.KeyVault(SecretUri=...)) are automatically resolved to their
actual values. Key Vault references are resolved concurrently (up to 10 at
a time) and the output file is written atomically to avoid partial writes.

Usage:
  az-loadenv --app <name> --rg <resource-group> [flags]

Flags:
  --app     string   Azure Web App name (required)
  --rg      string   Resource group name (required)
  -o        string   Output file path (default ".env")
  --raw              Write values without quoting or escaping
  --version          Print version and exit

Authentication:
  Uses DefaultAzureCredential, which tries the following in order:
    1. Environment variables (AZURE_CLIENT_ID, AZURE_CLIENT_SECRET, AZURE_TENANT_ID)
    2. Workload identity (Kubernetes, GitHub Actions)
    3. Managed identity (Azure VMs, App Service, Container Apps)
    4. Azure CLI (az login)
    5. Azure Developer CLI (azd auth login)

Subscription detection:
  The Azure subscription is detected automatically:
    1. AZURE_SUBSCRIPTION_ID environment variable (if set)
    2. Default subscription from ~/.azure/azureProfile.json (az login)

Output format:
  Settings are written in KEY=VALUE format, sorted alphabetically by key.
  Values containing spaces, quotes, newlines, or other special characters
  are automatically double-quoted and escaped. The output file is created
  with 0600 permissions (owner read/write only) to protect secrets.

Environment variables:
  AZURE_SUBSCRIPTION_ID       Override subscription detection
  AZURE_CLIENT_ID             Service principal authentication
  AZURE_CLIENT_SECRET         Service principal authentication
  AZURE_TENANT_ID             Service principal authentication

Examples:
  # Write settings to .env (default)
  az-loadenv --app my-api --rg my-resource-group

  # Write settings to a custom file
  az-loadenv --app my-api --rg my-resource-group -o .env.local

  # Use with a specific subscription
  AZURE_SUBSCRIPTION_ID=xxx az-loadenv --app my-api --rg my-resource-group

  # Use with a service principal (CI/CD)
  export AZURE_TENANT_ID=xxx AZURE_CLIENT_ID=xxx AZURE_CLIENT_SECRET=xxx
  az-loadenv --app my-api --rg my-resource-group

Documentation:
  https://curiousdev.github.io/az-loadenv/
`

func main() {
	flag.Usage = func() {
		fmt.Fprint(os.Stderr, helpText)
	}

	showVersion := flag.Bool("version", false, "")
	raw := flag.Bool("raw", false, "")
	app := flag.String("app", "", "")
	rg := flag.String("rg", "", "")
	output := flag.String("o", ".env", "")
	flag.Parse()

	if *showVersion {
		fmt.Printf("az-loadenv %s+%s\n", version, build)
		os.Exit(0)
	}

	if *app == "" || *rg == "" {
		flag.Usage()
		os.Exit(1)
	}

	subscriptionID, err := detectSubscription()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to detect subscription: %v\n", err)
		os.Exit(1)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	ctx, cancel := context.WithTimeout(ctx, 2*time.Minute)
	defer cancel()

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to authenticate: %v\n", err)
		os.Exit(1)
	}

	client, err := armappservice.NewWebAppsClient(subscriptionID, cred, nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create web apps client: %v\n", err)
		os.Exit(1)
	}

	settings, err := client.ListApplicationSettings(ctx, *rg, *app, nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to list app settings: %v\n", err)
		os.Exit(1)
	}

	if settings.Properties == nil {
		fmt.Fprintln(os.Stderr, "No app settings found")
		os.Exit(0)
	}

	// Collect and sort keys for deterministic output
	keys := make([]string, 0, len(settings.Properties))
	for k := range settings.Properties {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	entries := make([]setting, len(keys))
	secrets := 0
	for i, name := range keys {
		value := ""
		if settings.Properties[name] != nil {
			value = *settings.Properties[name]
		}
		entries[i] = setting{name: name, value: value}
		if kvRefPattern.MatchString(value) {
			entries[i].secret = true
			secrets++
		} else {
			fmt.Fprintf(os.Stderr, "  %s\n", name)
		}
	}

	// Resolve Key Vault references concurrently, printing each as it resolves
	if secrets > 0 {
		resolveSecrets(ctx, cred, entries)
	}

	// Atomic file write: write to temp file, then rename
	if err := atomicWriteEnv(*output, entries, *raw); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to write %s: %v\n", *output, err)
		os.Exit(1)
	}

	fmt.Fprintf(os.Stderr, "\nWrote %d settings (%d config, %d secrets) to %s\n",
		len(entries), len(entries)-secrets, secrets, *output)
}

// resolveSecrets resolves Key Vault references concurrently with bounded parallelism,
// printing each secret name to stderr as it resolves.
func resolveSecrets(ctx context.Context, cred azcore.TokenCredential, entries []setting) {
	var mu sync.Mutex
	vaultClients := make(map[string]*azsecrets.Client)

	g, ctx := errgroup.WithContext(ctx)
	g.SetLimit(10)

	for i := range entries {
		if !entries[i].secret {
			continue
		}
		i := i
		g.Go(func() error {
			match := kvRefPattern.FindStringSubmatch(entries[i].value)
			if match == nil {
				return nil
			}
			resolved, err := resolveKeyVaultRef(ctx, cred, &mu, vaultClients, match[1])
			if err != nil {
				entries[i].errMsg = err.Error()
				fmt.Fprintf(os.Stderr, "  %s (secret: error)\n", entries[i].name)
			} else {
				entries[i].value = resolved
				fmt.Fprintf(os.Stderr, "  %s (secret)\n", entries[i].name)
			}
			return nil // don't fail the group; individual errors are tracked per-entry
		})
	}

	_ = g.Wait()
}

func resolveKeyVaultRef(ctx context.Context, cred azcore.TokenCredential, mu *sync.Mutex, clients map[string]*azsecrets.Client, secretURI string) (string, error) {
	u, err := url.Parse(secretURI)
	if err != nil {
		return "", fmt.Errorf("invalid secret URI: %w", err)
	}

	vaultURL := fmt.Sprintf("%s://%s", u.Scheme, u.Host)

	mu.Lock()
	client, ok := clients[vaultURL]
	if !ok {
		client, err = azsecrets.NewClient(vaultURL, cred, nil)
		if err != nil {
			mu.Unlock()
			return "", fmt.Errorf("creating keyvault client: %w", err)
		}
		clients[vaultURL] = client
	}
	mu.Unlock()

	// Path is /secrets/<name> or /secrets/<name>/<version>
	parts := strings.Split(strings.TrimPrefix(u.Path, "/"), "/")
	if len(parts) < 2 || parts[0] != "secrets" {
		return "", fmt.Errorf("unexpected secret URI path: %s", u.Path)
	}
	secretName := parts[1]
	version := ""
	if len(parts) >= 3 {
		version = parts[2]
	}

	resp, err := client.GetSecret(ctx, secretName, version, nil)
	if err != nil {
		return "", fmt.Errorf("getting secret %s: %w", secretName, err)
	}

	if resp.Value == nil {
		return "", fmt.Errorf("secret %s has nil value", secretName)
	}
	return *resp.Value, nil
}

// needsQuoting returns true if the value contains characters that require
// double-quoting in a .env file.
func needsQuoting(val string) bool {
	if len(val) == 0 {
		return false
	}
	if val[0] == ' ' || val[0] == '\t' || val[len(val)-1] == ' ' || val[len(val)-1] == '\t' {
		return true
	}
	return strings.ContainsAny(val, " \t\n\r\"'#`$\\")
}

// formatEnvValue formats a value for a .env file, quoting if necessary.
func formatEnvValue(val string) string {
	if !needsQuoting(val) {
		return val
	}
	escaped := strings.ReplaceAll(val, `\`, `\\`)
	escaped = strings.ReplaceAll(escaped, `"`, `\"`)
	escaped = strings.ReplaceAll(escaped, "\n", `\n`)
	escaped = strings.ReplaceAll(escaped, "\r", `\r`)
	return `"` + escaped + `"`
}

// atomicWriteEnv writes the .env file atomically via a temp file + rename.
func atomicWriteEnv(path string, entries []setting, raw bool) error {
	dir := filepath.Dir(path)
	tmp, err := os.CreateTemp(dir, ".env.tmp.*")
	if err != nil {
		return fmt.Errorf("creating temp file: %w", err)
	}
	tmpPath := tmp.Name()
	defer os.Remove(tmpPath) // clean up on any failure path

	var lines []string
	for _, e := range entries {
		val := e.value
		if !raw {
			val = formatEnvValue(val)
		}
		lines = append(lines, fmt.Sprintf("%s=%s", e.name, val))
	}
	content := strings.Join(lines, "\n") + "\n"

	if _, err := tmp.WriteString(content); err != nil {
		tmp.Close()
		return fmt.Errorf("writing temp file: %w", err)
	}
	if err := tmp.Close(); err != nil {
		return fmt.Errorf("closing temp file: %w", err)
	}

	if err := os.Chmod(tmpPath, 0600); err != nil {
		return fmt.Errorf("setting file permissions: %w", err)
	}

	if err := os.Rename(tmpPath, path); err != nil {
		return fmt.Errorf("renaming temp file: %w", err)
	}

	return nil
}

// azureProfile represents the structure of ~/.azure/azureProfile.json
type azureProfile struct {
	Subscriptions []struct {
		ID        string `json:"id"`
		IsDefault bool   `json:"isDefault"`
	} `json:"subscriptions"`
}

func detectSubscription() (string, error) {
	// Check environment variable first (CI/CD, containers)
	if id := os.Getenv("AZURE_SUBSCRIPTION_ID"); id != "" {
		return id, nil
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("cannot find home directory: %w", err)
	}

	profilePath := filepath.Join(homeDir, ".azure", "azureProfile.json")
	data, err := os.ReadFile(profilePath)
	if err != nil {
		return "", fmt.Errorf("cannot read %s: %w (run 'az login' first)", profilePath, err)
	}

	// Azure CLI sometimes writes a UTF-8 BOM; strip it
	data = stripBOM(data)

	var profile azureProfile
	if err := json.Unmarshal(data, &profile); err != nil {
		return "", fmt.Errorf("cannot parse %s: %w", profilePath, err)
	}

	for _, sub := range profile.Subscriptions {
		if sub.IsDefault {
			return sub.ID, nil
		}
	}

	return "", fmt.Errorf("no default subscription found in %s (run 'az account set')", profilePath)
}

func stripBOM(data []byte) []byte {
	if len(data) >= 3 && data[0] == 0xEF && data[1] == 0xBB && data[2] == 0xBF {
		return data[3:]
	}
	return data
}
