#Requires -Version 5.1
$ErrorActionPreference = 'Stop'

$InstallDir = Join-Path $env:USERPROFILE '.curriculum\bin'
$RepoDir = Split-Path -Parent $MyInvocation.MyCommand.Definition
$Version = (Get-Content (Join-Path $RepoDir 'VERSION') -Raw).Trim()
$Commit = try { (git -C $RepoDir rev-parse --short HEAD 2>$null) } catch { 'unknown' }
if (-not $Commit) { $Commit = 'unknown' }
$Ldflags = "-s -w -X curriculum/cmd.version=$Version -X curriculum/cmd.commit=$Commit"

Write-Host "Building cur v${Version}+${Commit}..."
Push-Location (Join-Path $RepoDir 'src')
try {
    & go build -ldflags $Ldflags -o (Join-Path $RepoDir 'cur.exe') .
    if ($LASTEXITCODE -ne 0) { throw 'Build failed' }
} finally {
    Pop-Location
}

Write-Host "Installing to $InstallDir..."
if (-not (Test-Path $InstallDir)) {
    New-Item -ItemType Directory -Path $InstallDir -Force | Out-Null
}
Move-Item -Force (Join-Path $RepoDir 'cur.exe') (Join-Path $InstallDir 'cur.exe')

# Add to user-level PATH if not already present
$UserPath = [Environment]::GetEnvironmentVariable('Path', 'User')
if ($UserPath -split ';' | Where-Object { $_ -eq $InstallDir }) {
    Write-Host "PATH already contains $InstallDir"
} else {
    $NewPath = $InstallDir + ';' + $UserPath
    [Environment]::SetEnvironmentVariable('Path', $NewPath, 'User')
    $env:PATH = $InstallDir + ';' + $env:PATH
    Write-Host "Added $InstallDir to user PATH"
}

Write-Host 'Done. Restart your terminal, then run: cur --help'
