# Contributing

Contributions that improve Local compatibility, platform discovery, tests, diagnostics, documentation or agent workflows are welcome.

## Development requirements

- Go 1.24+
- Git
- Local for real integration testing

## Before opening a pull request

Run:

```sh
gofmt -w ./cmd ./internal
go test ./...
go vet ./...
go build ./cmd/localwp
```

Or:

```sh
make check
```

## Cross-platform changes

A change that touches path discovery, environment variables, process execution or installation should be considered for all three host families:

- Windows
- macOS
- Linux

GitHub Actions runs native tests on each platform and cross-builds all supported release architectures.

## Local compatibility

Local owns its internal runtime layout, and those details can change between releases. Prefer resilient lookup over one exact hard-coded versioned path.

When adding a new discovery rule:

1. keep existing fallbacks when safe;
2. add a test;
3. expose the resolved value through `localwp --doctor` when useful;
4. update architecture/install docs if behavior changes.

## Tests versus real Local

Automated tests should not require the Local desktop GUI. Use fixtures and environment overrides for deterministic test coverage.

For runtime-layout changes, additionally test on a real running Local site before release when possible.
