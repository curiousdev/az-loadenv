# Releases

## Versioning

`az-loadenv` uses [semantic versioning](https://semver.org/) with build metadata:

```
<major>.<minor>.<patch>+<build>
```

- **major** — breaking changes to CLI flags or behavior
- **minor** — new features, backwards-compatible
- **patch** — bug fixes
- **build** — auto-incrementing CI build number

Check the version of your installed binary:

```bash
az-loadenv --version
# az-loadenv 1.2.0+47
```

## Release process

Releases are automated through GitHub Actions:

### 1. Tag a version

```bash
git tag v1.0.0
git push origin v1.0.0
```

### 2. CI builds all targets

The CI workflow detects the `v*` tag and:
- Builds binaries for all 5 platform targets
- Injects the version and build number into the binary via ldflags
- Archives each binary (`.tar.gz` for Linux/macOS, `.zip` for Windows)

### 3. GitHub Release is created

A GitHub Release is published automatically with:
- All platform archives attached as downloadable assets
- Auto-generated release notes from commit history

## Release artifacts

Each release includes:

| File | Contents |
|---|---|
| `az-loadenv-linux-amd64.tar.gz` | `az-loadenv` binary |
| `az-loadenv-linux-arm64.tar.gz` | `az-loadenv` binary |
| `az-loadenv-darwin-amd64.tar.gz` | `az-loadenv` binary |
| `az-loadenv-darwin-arm64.tar.gz` | `az-loadenv` binary |
| `az-loadenv-windows-amd64.zip` | `az-loadenv.exe` binary |

## Pre-release versions

For pre-release testing, use a pre-release tag:

```bash
git tag v1.1.0-rc.1
git push origin v1.1.0-rc.1
```

This still triggers the release workflow. You can mark it as a pre-release in the GitHub UI after creation.
