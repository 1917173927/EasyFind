param(
    [ValidateSet("amd64", "arm64")]
    [string]$Arch = "amd64",
    [string]$OutputDir = "dist",
    [string]$BinaryName = "easyfind"
)

$ErrorActionPreference = "Stop"

$projectRoot = Split-Path -Parent $PSScriptRoot
Push-Location $projectRoot

$oldGoos = $env:GOOS
$oldGoarch = $env:GOARCH
$oldCgo = $env:CGO_ENABLED

try {
    if (-not (Test-Path $OutputDir)) {
        New-Item -Path $OutputDir -ItemType Directory | Out-Null
    }

    $env:CGO_ENABLED = "0"
    $env:GOOS = "linux"
    $env:GOARCH = $Arch

    $outputFile = Join-Path $OutputDir "$BinaryName-linux-$Arch"
    Write-Host "Building $outputFile ..."

    go build -trimpath -ldflags "-s -w" -o $outputFile ./cmd

    Write-Host "Build success: $outputFile"
}
finally {
    $env:GOOS = $oldGoos
    $env:GOARCH = $oldGoarch
    $env:CGO_ENABLED = $oldCgo
    Pop-Location
}
