# CI/CD

The repository uses GitHub Actions to validate all supported host platforms and produce release binaries.

## Continuous integration

Workflow:

```text
.github/workflows/ci.yml
```

It runs for pushes, pull requests and manual dispatches.

### Native matrix

Tests execute natively on:

- `windows-latest`
- `macos-latest`
- `ubuntu-latest`

Each runner uses the latest Go 1.26 patch release with `CGO_ENABLED=0` and performs:

```sh
go test ./...
go vet ./...
go build -trimpath ./cmd/localwp
go run ./cmd/localwp --version
go run ./cmd/localwp --help
```

The project minimum is Go 1.24. Go 1.24 changed the macOS linker to emit a Mach-O `LC_UUID` load command by default; older Go toolchains can fail on newer macOS runners with `dyld: missing LC_UUID load command`. CI intentionally uses a newer supported Go line instead of deriving the CI toolchain from the minimum version in `go.mod`.

### Formatting

A separate job fails when committed Go files are not formatted by `gofmt`.

### Cross-build validation

CI also cross-compiles:

- Windows amd64
- Linux amd64
- Linux arm64
- macOS amd64
- macOS arm64

The resulting binaries are uploaded as CI artifacts.

## What CI tests

The automated suite tests deterministic code that does not require a GUI:

- Local `sites.json` parsing;
- site selection by directory/name/domain/ID;
- Local runtime path resolution;
- wrapper option parsing;
- preservation of WP-CLI arguments;
- environment construction;
- compilation and CLI startup on every host OS.

## Real Local integration testing

GitHub-hosted runners do not provision and launch the Local desktop application as part of CI.

Before important releases, a maintainer should additionally smoke-test against real Local installations when practical:

```sh
localwp --sites
localwp --doctor
localwp core version
localwp plugin list
localwp option get siteurl
```

This validates compatibility with the currently installed Local release in addition to deterministic automated tests.

## Release CD

Workflow:

```text
.github/workflows/release.yml
```

A pushed `v*` tag triggers release delivery.

Example:

```sh
git tag v1.0.0
git push origin v1.0.0
```

The workflow:

1. runs tests and vet on Windows, macOS and Linux;
2. builds all release architectures;
3. injects the tag version into the binary;
4. packages Windows as ZIP and Unix targets as tar.gz;
5. includes README and LICENSE in each package;
6. creates `SHA256SUMS.txt`;
7. publishes a GitHub Release using generated release notes.

A release can also be started manually with `workflow_dispatch` for an existing tag.
