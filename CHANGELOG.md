# Changelog

All notable changes to this project will be documented here.

## Unreleased

### Fixed

- Fixed macOS 26 GitHub Actions test crashes (`dyld: missing LC_UUID load command`) by raising the source-build minimum to Go 1.24 and running CI/release verification on Go 1.26.x.
- Disabled CGO in the native CI verification jobs for deterministic pure-Go test and build binaries.

## 1.0.0 - Initial release

### Added

- Cross-platform `localwp <wp-cli arguments>` CLI for Windows, macOS and Linux.
- Automatic Local site detection from the current directory.
- Explicit site selection by name, domain, or Local site ID.
- Local `sites.json` discovery for Windows, macOS and Linux.
- Local-managed PHP, site `php.ini`, WP-CLI and database runtime discovery.
- `localwp --sites` site listing.
- `localwp --doctor` runtime diagnostics.
- Optional `localwp --shell` interactive Local Site Shell access.
- Windows and POSIX source installers/uninstallers.
- Reusable AI coding-agent skill.
- GitHub Actions CI on Windows, macOS and Linux.
- Cross-compilation for Windows amd64, Linux amd64/arm64, and macOS amd64/arm64.
- Tag-driven GitHub Release workflow with SHA-256 checksums.
