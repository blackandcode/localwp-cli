# Contributing

Contributions that improve Local compatibility, platform discovery, tests, diagnostics, documentation or agent workflows are welcome.

## Development requirements

- Go 1.24+
- Git
- Local for real integration testing

## Build and install from source

End users should use the [release installation guide](installation/README.md). For development, clone and build your checkout:

```sh
git clone https://github.com/blackandcode/localwp-cli.git
cd localwp-cli
go build -o localwp ./cmd/localwp
```

On Windows, use `go build -o localwp.exe ./cmd/localwp`.

### Windows

From the repository root:

```cmd
.\install-localwp.cmd
```

This builds the checkout, installs `%LOCALAPPDATA%\localwp-cli\bin\localwp.exe`, and adds the folder to your User PATH. Use the explicit `.\` prefix to select this repository's launcher. Close and reopen the terminal application and editor afterward.

### macOS and Linux

From the repository root:

```sh
sh ./install-localwp.sh
```

This builds into `~/.local/bin/localwp`. To choose a different directory:

```sh
LOCALWP_INSTALL_DIR="$HOME/bin" sh ./install-localwp.sh
```

If needed, add `export PATH="$HOME/.local/bin:$PATH"` to your shell startup file (for example `~/.zshrc` for zsh or `~/.bashrc` for interactive Bash) and open a new terminal. Use your chosen directory if you changed the default.

### Verify

```sh
localwp --version
localwp --sites
```

Start your development site in Local, move into its project directory, and run `localwp --doctor` followed by `localwp plugin list`.

The root `install-localwp.*` scripts build from source. The separate `installation/install-release.*` helpers download published binaries and do not require Go.

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
