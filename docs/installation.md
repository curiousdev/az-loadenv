# Installation

## Package managers

### Homebrew (macOS / Linux)

```bash
brew tap curiousdev/tap
brew install az-loadenv-cli
```

The `brew tap` only needs to be run once. After that, `brew install az-loadenv-cli` and `brew upgrade az-loadenv-cli` work directly.

### WinGet (Windows)

```powershell
winget install curiousdev.az-loadenv-cli
```

Updates automatically with `winget upgrade`.

## Install scripts

One-liner install scripts that detect your OS and architecture automatically.

**macOS / Linux:**
```bash
curl -fsSL https://curiousdev.github.io/az-loadenv/install.sh | bash
```

**Windows (PowerShell):**
```powershell
irm https://curiousdev.github.io/az-loadenv/install.ps1 | iex
```

The shell script installs to `/usr/local/bin` (prompts for sudo if needed). The PowerShell script installs to `%LOCALAPPDATA%\az-loadenv` and adds it to your user PATH.

## Download a pre-built binary

Pre-built binaries are published with every release on GitHub:

| Platform | Architecture | Archive |
|---|---|---|
| Linux | x86_64 | `az-loadenv-linux-amd64.tar.gz` |
| Linux | ARM64 | `az-loadenv-linux-arm64.tar.gz` |
| macOS | Intel | `az-loadenv-darwin-amd64.tar.gz` |
| macOS | Apple Silicon | `az-loadenv-darwin-arm64.tar.gz` |
| Windows | x86_64 | `az-loadenv-windows-amd64.zip` |

### macOS / Linux

```bash
# macOS Apple Silicon
curl -L https://github.com/curiousdev/az-loadenv/releases/latest/download/az-loadenv-darwin-arm64.tar.gz | tar xz
sudo mv az-loadenv /usr/local/bin/

# macOS Intel
curl -L https://github.com/curiousdev/az-loadenv/releases/latest/download/az-loadenv-darwin-amd64.tar.gz | tar xz
sudo mv az-loadenv /usr/local/bin/

# Linux x86_64
curl -L https://github.com/curiousdev/az-loadenv/releases/latest/download/az-loadenv-linux-amd64.tar.gz | tar xz
sudo mv az-loadenv /usr/local/bin/

# Linux ARM64
curl -L https://github.com/curiousdev/az-loadenv/releases/latest/download/az-loadenv-linux-arm64.tar.gz | tar xz
sudo mv az-loadenv /usr/local/bin/
```

### Windows

1. Download `az-loadenv-windows-amd64.zip` from the [latest release](https://github.com/curiousdev/az-loadenv/releases/latest)
2. Extract `az-loadenv.exe`
3. Move it to a directory in your `PATH` (e.g. `C:\Users\<you>\bin`)

## Install with Go

Requires Go 1.25 or later:

```bash
go install github.com/curiousdev/az-loadenv@latest
```

This places the binary in `$GOPATH/bin` (usually `~/go/bin`).

## Verify installation

```bash
az-loadenv --version
```

You should see output like:

```
az-loadenv 1.0.0+42
```
