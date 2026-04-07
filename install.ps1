$ErrorActionPreference = "Stop"

$Version = (Get-Content VERSION).Trim()
$Commit  = (git rev-parse --short HEAD 2>$null) ?? "unknown"
$Ldflags = "-s -w -X curriculum/cmd.version=$Version -X curriculum/cmd.commit=$Commit"

Write-Host "Building cur $Version+$Commit..."
Push-Location src
go install -ldflags $Ldflags .
Pop-Location
Write-Host "Installed: $(Get-Command cur -ErrorAction SilentlyContinue | Select-Object -ExpandProperty Source)"
