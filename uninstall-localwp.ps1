$ErrorActionPreference = 'Stop'
$installDir = Join-Path $env:LOCALAPPDATA 'localwp-cli\bin'
$root = Split-Path -Parent $installDir

$userPath = [Environment]::GetEnvironmentVariable('Path', 'User')
if ($null -ne $userPath) {
    $entries = @($userPath -split ';' | Where-Object {
        -not [string]::IsNullOrWhiteSpace($_) -and
        -not $_.TrimEnd('\').Equals($installDir.TrimEnd('\'), [System.StringComparison]::OrdinalIgnoreCase)
    })
    [Environment]::SetEnvironmentVariable('Path', ($entries -join ';'), 'User')
}

if (Test-Path $root) {
    Remove-Item -LiteralPath $root -Recurse -Force
}

Write-Host "localwp removed. Open a new terminal for PATH changes to take effect."
