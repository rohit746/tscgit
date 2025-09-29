#!/usr/bin/env pwsh
Set-StrictMode -Version Latest
$ErrorActionPreference = 'Stop'

function New-TemporaryDirectory {
    param(
        [string]$Prefix = 'tscgit-lessons-'
    )
    $tempPath = [System.IO.Path]::GetTempPath()
    $dir = Join-Path $tempPath ($Prefix + [System.IO.Path]::GetRandomFileName())
    New-Item -ItemType Directory -Path $dir | Out-Null
    return $dir
}

$ScriptDir = Split-Path -Parent $MyInvocation.MyCommand.Definition
$RootDir = Resolve-Path (Join-Path $ScriptDir '..')
$BinDir = Join-Path $RootDir 'bin'
if (-not (Test-Path $BinDir)) {
    New-Item -ItemType Directory -Path $BinDir | Out-Null
}

$TscgitBin = Join-Path $BinDir 'tscgit.exe'
Write-Host "Building tscgit CLI..."
go build -o $TscgitBin (Join-Path $RootDir 'cmd/tscgit')

$workDir = New-TemporaryDirectory
Write-Host "Using temp workspace: $workDir"
try {
    $gitConfig = Join-Path $workDir 'gitconfig'
    @"
[user]
    name = Test Student
    email = student@example.com
[init]
    defaultBranch = master
"@ | Set-Content -LiteralPath $gitConfig

    $env:GIT_CONFIG_GLOBAL = $gitConfig
    $env:GIT_AUTHOR_NAME = 'Test Student'
    $env:GIT_AUTHOR_EMAIL = 'student@example.com'
    $env:GIT_COMMITTER_NAME = $env:GIT_AUTHOR_NAME
    $env:GIT_COMMITTER_EMAIL = $env:GIT_AUTHOR_EMAIL

    $repoDir = Join-Path $workDir 'webflyx'
    New-Item -ItemType Directory -Path $repoDir | Out-Null
    Set-Location $repoDir

    git init | Out-Null
    '# Webflyx Practice' | Set-Content -LiteralPath 'README.md'

    & $TscgitBin run 0
    & $TscgitBin run 1
    & $TscgitBin run 2
    & $TscgitBin run 3

    '# contents' | Set-Content -LiteralPath 'contents.md'
    & $TscgitBin run 4a

    git add contents.md
    & $TscgitBin run 4b

    git commit -m 'A: add contents.md' | Out-Null
    & $TscgitBin run 5

    git cat-file commit HEAD | Set-Content -LiteralPath 'catfileout.txt'
    $blobId = (git rev-parse HEAD:contents.md).Trim()
    git cat-file blob $blobId | Set-Content -LiteralPath 'blobfile.txt'
    & $TscgitBin run 6a
    & $TscgitBin run 6b

    '# Titles' | Set-Content -LiteralPath 'titles.md'
    git add titles.md
    git commit -m 'B: add titles.md' | Out-Null
    & $TscgitBin run 7

    git config --global init.defaultBranch main
    git branch -M main
    & $TscgitBin run 8

    git switch -c add_classics | Out-Null
    & $TscgitBin run 9

    "title,year" | Set-Content -LiteralPath 'classics.csv'
    "Metropolis,1927" | Add-Content -LiteralPath 'classics.csv'
    git add classics.csv
    git commit -m 'C: add classics.csv' | Out-Null
    & $TscgitBin run 10

    git switch main | Out-Null
    '' | Add-Content -LiteralPath 'contents.md'
    '## Classics' | Add-Content -LiteralPath 'contents.md'
    git add contents.md
    git commit -m 'D: update contents.md' | Out-Null
    & $TscgitBin run 11a

    git merge add_classics --no-ff -m 'E: merge add_classics' | Out-Null
    & $TscgitBin run 11b

    git branch -D 'C:' *> $null
    git switch -c feature/lesson-branch | Out-Null
    '- Added new section' | Add-Content -LiteralPath 'contents.md'
    git add contents.md
    git commit -m '[branch] add lesson notes' | Out-Null

    & $TscgitBin verify init-basics
    & $TscgitBin verify branch-basics

    Write-Host "`nAll lessons and verifications completed successfully." -ForegroundColor Green
}
finally {
    Set-Location $RootDir
    if (Test-Path $workDir) {
        try {
            Remove-Item -Recurse -Force $workDir
        }
        catch {
            Start-Sleep -Seconds 1
            Remove-Item -Recurse -Force $workDir -ErrorAction SilentlyContinue
        }
    }
    Remove-Item Env:GIT_CONFIG_GLOBAL -ErrorAction SilentlyContinue
    Remove-Item Env:GIT_AUTHOR_NAME -ErrorAction SilentlyContinue
    Remove-Item Env:GIT_AUTHOR_EMAIL -ErrorAction SilentlyContinue
    Remove-Item Env:GIT_COMMITTER_NAME -ErrorAction SilentlyContinue
    Remove-Item Env:GIT_COMMITTER_EMAIL -ErrorAction SilentlyContinue
}
