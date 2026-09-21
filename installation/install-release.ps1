$ErrorActionPreference = 'Stop'

# Install only the published Windows x64 binary, without Go or administrator access.
if ($env:OS -ne 'Windows_NT' -or $env:PROCESSOR_ARCHITECTURE -eq 'ARM64' -or
    $env:PROCESSOR_ARCHITEW6432 -eq 'ARM64' -or -not [Environment]::Is64BitOperatingSystem) {
    throw 'This installer requires x64 Windows. See installation/README.md for supported platforms.'
}
[Net.ServicePointManager]::SecurityProtocol = [Net.SecurityProtocolType]::Tls12
$repo = 'https://api.github.com/repos/blackandcode/localwp-cli/releases/latest'
$release = Invoke-RestMethod -Uri $repo -Headers @{ 'User-Agent' = 'localwp-cli-installer' }
$tag = $release.tag_name
if ($tag -notmatch '^v[0-9A-Za-z._-]+$') { throw 'Unexpected release tag.' }
$name = "localwp_$($tag.Substring(1))_windows_amd64"
$archive = "$name.zip"
$base = "https://github.com/blackandcode/localwp-cli/releases/download/$tag"
$tempDir = Join-Path ([IO.Path]::GetTempPath()) ("localwp-install-" + [guid]::NewGuid())
New-Item -ItemType Directory -Path $tempDir | Out-Null
try {
    Write-Host "Downloading localwp $tag from GitHub Releases..."
    Invoke-WebRequest -UseBasicParsing -Uri "$base/$archive" -OutFile (Join-Path $tempDir $archive)
    Invoke-WebRequest -UseBasicParsing -Uri "$base/SHA256SUMS.txt" -OutFile (Join-Path $tempDir 'SHA256SUMS.txt')
    $lines = @(Get-Content (Join-Path $tempDir 'SHA256SUMS.txt') | Where-Object {
        $_ -match ('^[0-9a-fA-F]{64}\s+\*?' + [regex]::Escape($archive) + '$')
    })
    if ($lines.Count -ne 1) { throw 'Missing or ambiguous release checksum.' }
    $expected = ($lines[0] -split '\s+')[0]
    $actual = (Get-FileHash -Algorithm SHA256 -LiteralPath (Join-Path $tempDir $archive)).Hash
    if ($actual -ne $expected) { throw 'Checksum mismatch; installation stopped.' }
    Expand-Archive -LiteralPath (Join-Path $tempDir $archive) -DestinationPath $tempDir
    $installDir = Join-Path $env:LOCALAPPDATA 'localwp-cli\bin'
    New-Item -ItemType Directory -Path $installDir -Force | Out-Null
    Copy-Item -LiteralPath (Join-Path $tempDir "$name/localwp.exe") -Destination (Join-Path $installDir 'localwp.exe') -Force
    $userPath = [Environment]::GetEnvironmentVariable('Path', 'User')
    $entries = @($userPath -split ';' | ForEach-Object { $_.Trim().TrimEnd('\') })
    if ($entries -notcontains $installDir.TrimEnd('\')) {
        $newPath = if ([string]::IsNullOrWhiteSpace($userPath)) { $installDir } else { $userPath.TrimEnd(';') + ';' + $installDir }
        [Environment]::SetEnvironmentVariable('Path', $newPath, 'User')
    }
    Write-Host "Installed $tag to $installDir and added it to your User PATH."
    Write-Host 'Close and reopen your terminal application (and editor), then run: localwp --version'
}
finally {
    # Only remove this invocation's uniquely named temporary directory.
    $resolvedTemp = [IO.Path]::GetFullPath($tempDir)
    $tempRoot = [IO.Path]::GetFullPath([IO.Path]::GetTempPath()).TrimEnd('\') + '\'
    if ($resolvedTemp.StartsWith($tempRoot, [StringComparison]::OrdinalIgnoreCase) -and
        (Split-Path -Leaf $resolvedTemp) -match '^localwp-install-[0-9a-f-]{36}$') {
        Remove-Item -LiteralPath $resolvedTemp -Recurse -Force
    }
}
