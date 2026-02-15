$ErrorActionPreference = "Stop"

$repo = "curiousdev/az-loadenv"
$binary = "az-loadenv.exe"
$archive = "az-loadenv-windows-amd64.zip"
$installDir = "$env:LOCALAPPDATA\az-loadenv"
$url = "https://github.com/$repo/releases/latest/download/$archive"

Write-Host "Downloading $archive..."

$tmp = New-TemporaryFile | Rename-Item -NewName { $_.Name + ".zip" } -PassThru
Invoke-WebRequest -Uri $url -OutFile $tmp.FullName -UseBasicParsing

Write-Host "Extracting..."

$extractDir = Join-Path $env:TEMP "az-loadenv-install"
if (Test-Path $extractDir) { Remove-Item $extractDir -Recurse -Force }
Expand-Archive -Path $tmp.FullName -DestinationPath $extractDir
Remove-Item $tmp.FullName -Force

if (-not (Test-Path $installDir)) { New-Item -ItemType Directory -Path $installDir | Out-Null }
Move-Item -Path (Join-Path $extractDir $binary) -Destination (Join-Path $installDir $binary) -Force
Remove-Item $extractDir -Recurse -Force

# Add to PATH if not already present
$userPath = [Environment]::GetEnvironmentVariable("Path", "User")
if ($userPath -notlike "*$installDir*") {
    [Environment]::SetEnvironmentVariable("Path", "$userPath;$installDir", "User")
    $env:Path = "$env:Path;$installDir"
    Write-Host "Added $installDir to user PATH."
}

Write-Host "Installed to $installDir\$binary"
& "$installDir\$binary" --version
