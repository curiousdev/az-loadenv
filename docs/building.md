# Building from source

## Prerequisites

- [Go](https://go.dev/dl/) 1.25 or later

## Development build

```bash
git clone https://github.com/curiousdev/az-loadenv.git
cd az-loadenv
go build -o az-loadenv .
```

## Production build

Strip debug information and embed version metadata:

```bash
go build -trimpath \
  -ldflags="-s -w -X main.version=1.0.0 -X main.build=1" \
  -o az-loadenv .
```

| Flag | Purpose |
|---|---|
| `-trimpath` | Remove local file paths from the binary |
| `-s` | Strip symbol table |
| `-w` | Strip DWARF debug info |
| `-X main.version=...` | Set the version string |
| `-X main.build=...` | Set the build number |

## Cross-compilation

Go supports cross-compilation natively. Set `GOOS` and `GOARCH`:

```bash
# Linux ARM64
GOOS=linux GOARCH=arm64 go build -trimpath -ldflags="-s -w" -o az-loadenv .

# Windows x86_64
GOOS=windows GOARCH=amd64 go build -trimpath -ldflags="-s -w" -o az-loadenv.exe .

# macOS Apple Silicon
GOOS=darwin GOARCH=arm64 go build -trimpath -ldflags="-s -w" -o az-loadenv .
```

No additional toolchains or dependencies are required — Go handles everything.

## Supported targets

| GOOS | GOARCH | Platform |
|---|---|---|
| `linux` | `amd64` | Linux x86_64 |
| `linux` | `arm64` | Linux ARM64 (AWS Graviton, etc.) |
| `darwin` | `amd64` | macOS Intel |
| `darwin` | `arm64` | macOS Apple Silicon |
| `windows` | `amd64` | Windows x86_64 |

## CI builds

The project includes a GitHub Actions workflow (`.github/workflows/ci.yml`) that builds all five targets on every push and creates GitHub Releases with archives on version tags. See [Releases](releases.md) for details.
