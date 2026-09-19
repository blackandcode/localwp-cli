$ErrorActionPreference = 'Stop'

function Fail([string]$Message) {
    [Console]::Error.WriteLine("localwp install: $Message")
    exit 1
}

if (-not (Get-Command go -ErrorAction SilentlyContinue)) {
    Fail "Go is required when installing from source. Install Go, or download a prebuilt binary from the GitHub Releases page."
}

$repoRoot = Split-Path -Parent $MyInvocation.MyCommand.Path
$installDir = Join-Path $env:LOCALAPPDATA 'localwp-cli\bin'
$binary = Join-Path $installDir 'localwp.exe'

New-Item -ItemType Directory -Path $installDir -Force | Out-Null

Push-Location $repoRoot
try {
    & go build -trimpath -o $binary .\cmd\localwp
    if ($LASTEXITCODE -ne 0) { Fail "Go build failed." }
}
finally {
    Pop-Location
}

$userPath = [Environment]::GetEnvironmentVariable('Path', 'User')
if ($null -eq $userPath) { $userPath = '' }
$entries = @($userPath -split ';' | Where-Object { -not [string]::IsNullOrWhiteSpace($_) })
$present = $false
foreach ($entry in $entries) {
    if ($entry.TrimEnd('\').Equals($installDir.TrimEnd('\'), [System.StringComparison]::OrdinalIgnoreCase)) {
        $present = $true
        break
    }
}
if (-not $present) {
    $newPath = if ([string]::IsNullOrWhiteSpace($userPath)) { $installDir } else { $userPath.TrimEnd(';') + ';' + $installDir }
    [Environment]::SetEnvironmentVariable('Path', $newPath, 'User')
}

Write-Host "Installed localwp to:"
Write-Host "  $binary"
Write-Host ""
Write-Host "Open a NEW terminal, then run:"
Write-Host "  localwp --version"
Write-Host "  localwp --sites"
