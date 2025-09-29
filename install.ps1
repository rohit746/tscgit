#!/usr/bin/env pwsh

# tscgit installer for Windows
# Usage: iwr -useb https://raw.githubusercontent.com/rohit746/tscgit/main/install.ps1 | iex

$ErrorActionPreference = "Stop"

# Configuration
$repo = "rohit746/tscgit"
$installDir = "$env:LOCALAPPDATA\tscgit"
$binPath = "$installDir\tscgit.exe"

Write-Host "Installing tscgit..." -ForegroundColor Green

# Create install directory
if (-not (Test-Path $installDir)) {
    New-Item -ItemType Directory -Path $installDir -Force | Out-Null
}

# Detect architecture
$arch = if ([Environment]::Is64BitOperatingSystem) { "amd64" } else { "386" }
$os = "windows"

try {
    # Get latest release info
    Write-Host "Fetching latest release info..." -ForegroundColor Yellow
    $apiUrl = "https://api.github.com/repos/$repo/releases/latest"
    $release = Invoke-RestMethod -Uri $apiUrl -Headers @{"User-Agent" = "tscgit-installer"}
    
    $version = $release.tag_name
    Write-Host "Latest version: $version" -ForegroundColor Cyan
    
    # Find the Windows binary
    $asset = $release.assets | Where-Object { $_.name -like "*windows*$arch*" -and $_.name -like "*.zip" }
    if (-not $asset) {
        throw "Could not find Windows binary for architecture $arch"
    }
    
    $downloadUrl = $asset.browser_download_url
    $tempFile = "$env:TEMP\tscgit-$version.zip"
    $tempDir = "$env:TEMP\tscgit-$version"
    
    Write-Host "Downloading $($asset.name)..." -ForegroundColor Yellow
    Invoke-WebRequest -Uri $downloadUrl -OutFile $tempFile -UserAgent "tscgit-installer"
    
    # Extract and install
    Write-Host "Extracting..." -ForegroundColor Yellow
    if (Test-Path $tempDir) {
        Remove-Item $tempDir -Recurse -Force
    }
    Expand-Archive -Path $tempFile -DestinationPath $tempDir -Force
    
    # Find the binary in the extracted files
    $binaryPath = Get-ChildItem -Path $tempDir -Name "tscgit.exe" -Recurse | Select-Object -First 1
    if (-not $binaryPath) {
        throw "Could not find tscgit.exe in downloaded archive"
    }
    
    $fullBinaryPath = Join-Path $tempDir $binaryPath
    Copy-Item $fullBinaryPath $binPath -Force
    
    # Clean up
    Remove-Item $tempFile -Force -ErrorAction SilentlyContinue
    Remove-Item $tempDir -Recurse -Force -ErrorAction SilentlyContinue
    
    Write-Host "✓ tscgit installed to $binPath" -ForegroundColor Green
    
    # Add to PATH if not already there
    $userPath = [Environment]::GetEnvironmentVariable("PATH", "User")
    if ($userPath -notlike "*$installDir*") {
        Write-Host "Adding $installDir to PATH..." -ForegroundColor Yellow
        [Environment]::SetEnvironmentVariable("PATH", "$userPath;$installDir", "User")
        Write-Host "✓ Added to PATH (restart your terminal to use 'tscgit' directly)" -ForegroundColor Green
    }
    
    # Test installation
    Write-Host "`nTesting installation..." -ForegroundColor Yellow
    $version_output = & $binPath version 2>&1
    if ($LASTEXITCODE -eq 0) {
        Write-Host "✓ Installation successful!" -ForegroundColor Green
        Write-Host $version_output
        Write-Host "`nTo get started, restart your terminal and run:" -ForegroundColor Cyan
        Write-Host "  tscgit lessons" -ForegroundColor White
    } else {
        Write-Warning "Installation completed but binary test failed. You may need to restart your terminal."
    }
    
} catch {
    Write-Error "Installation failed: $($_.Exception.Message)"
    Write-Host "`nTry manual installation:" -ForegroundColor Yellow
    Write-Host "1. Download the latest release from: https://github.com/$repo/releases" -ForegroundColor White
    Write-Host "2. Extract the ZIP file" -ForegroundColor White  
    Write-Host "3. Copy tscgit.exe to a directory in your PATH" -ForegroundColor White
    exit 1
}