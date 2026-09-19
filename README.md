# localwp-cli

**Use WP-CLI with Local from the terminal you already have open. On Windows, macOS, and Linux.**

`localwp-cli` is a small command-line companion for [Local](https://localwp.com/) that automatically finds the right Local WordPress site, uses that site's Local-managed PHP/WP-CLI runtime, runs your command, and returns you to your normal terminal.

```sh
localwp plugin list
localwp core version
localwp cache flush
localwp option get siteurl
```

No manual Site Shell hopping. No guessing which PHP version belongs to the project. No persistent interactive shell for your coding agent to get stuck inside.

> [!IMPORTANT]
> **[Local](https://localwp.com/) is required.** Local is the WordPress development application this tool is built for. Install a current Local version for your operating system before using `localwp-cli`.
>
> Local publishes builds for Windows, macOS, and Debian-based Linux. See the official [Local installation documentation](https://localwp.com/help-docs/getting-started/installing-local/) for current supported OS versions and requirements.

[![CI](https://github.com/blackandcode/localwp-cli/actions/workflows/ci.yml/badge.svg)](https://github.com/blackandcode/localwp-cli/actions/workflows/ci.yml)
[![Platforms](https://img.shields.io/badge/platform-Windows%20%7C%20macOS%20%7C%20Linux-4c8bf5)](https://localwp.com/)
[![Local](https://img.shields.io/badge/requires-Local-7b5cff)](https://localwp.com/)
[![WP-CLI](https://img.shields.io/badge/powered%20by-WP--CLI-21759b)](https://wp-cli.org/)
[![License: MIT](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)

## The problem it solves

Local already gives each site a correctly configured Site Shell. That is useful interactively, but it becomes repetitive when your normal workflow already happens in Terminal, Windows Terminal, VS Code, Cursor, Claude Code, Codex, or another coding environment.

Without `localwp-cli`:

```text
Open Local
  -> select the site
  -> open Site Shell
  -> run wp ...
  -> exit
  -> return to your original terminal
```

With `localwp-cli`:

```sh
localwp plugin list
```

That is the entire workflow.

## Why use it?

- **One command on every supported OS** — the same `localwp <wp-cli arguments>` interface on Windows, macOS, and Linux.
- **Automatic site detection** — run it from anywhere inside a Local site directory tree.
- **Uses Local's runtime** — selects the PHP version, `php.ini`, WP-CLI and database tooling managed by Local for that site.
- **No global WP-CLI setup required** — the tool invokes Local's bundled `wp-cli.phar`.
- **No interactive shell required** — normal commands execute once and return immediately.
- **Agent friendly** — coding agents get a predictable command, output, and process exit code.
- **Explicit targeting when needed** — select a site by Local name, domain, or ID.
- **Built-in diagnostics** — `localwp --doctor` tells you exactly what runtime was resolved.
- **Real Site Shell still available** — `localwp --shell` remains available when you intentionally want an interactive Local shell.

## Quick start

Start the WordPress site in Local, then move into any directory under that Local site. For example:

```sh
cd ~/"Local Sites"/my-site/app/public/wp-content/plugins/my-plugin
```

On Windows the path may look like:

```cmd
cd "C:\Users\you\Local Sites\my-site\app\public\wp-content\plugins\my-plugin"
```

Now use normal WP-CLI arguments after `localwp`:

```sh
localwp plugin list
localwp plugin status
localwp plugin activate my-plugin
localwp theme list
localwp core version
localwp cache flush
localwp option get siteurl
localwp user list --format=json
localwp db export backup.sql
localwp search-replace "http://old.local" "https://new.local" --dry-run
```

Think of it as replacing:

```text
wp <command>
```

with:

```text
localwp <command>
```

when the WordPress project runs in Local.

## Requirements

For normal use:

- [Local](https://localwp.com/) installed.
- A WordPress site created/imported into Local.
- The target site started in Local for commands that require the database or running services.
- Windows, macOS, or Linux supported by your installed Local version.

For building or installing **from source**, Go 1.24+ is also required. End users installing a prebuilt GitHub Release binary do not need Go.

## Installation

### Option 1: GitHub Release binary

Every version tag is built by GitHub Actions for:

| OS | Architecture |
| --- | --- |
| Windows | x86-64 / amd64 |
| macOS | Intel / amd64 |
| macOS | Apple Silicon / arm64 |
| Linux | x86-64 / amd64 |
| Linux | arm64 |

Download the archive for your platform from **Releases**, extract `localwp` (`localwp.exe` on Windows), and place it in a directory on your `PATH`.

The release workflow also publishes `SHA256SUMS.txt` so downloaded binaries can be verified.

### Option 2: Install from source on Windows

Clone this repository and run:

```cmd
install.cmd
```

The installer builds the Go CLI and installs it to:

```text
%LOCALAPPDATA%\localwp-cli\bin\localwp.exe
```

It also adds that directory to your **User PATH**. Open a new terminal and verify:

```cmd
localwp --version
localwp --sites
```

### Option 3: Install from source on macOS or Linux

Clone the repository, then run:

```sh
./install.sh
```

By default the binary is installed to:

```text
~/.local/bin/localwp
```

If `~/.local/bin` is not already on `PATH`, the installer prints the exact line to add to your shell profile.

Verify:

```sh
localwp --version
localwp --sites
```

See [Installation](docs/INSTALLATION.md) for full setup and troubleshooting instructions.

## How site detection works

`localwp-cli` selects the target site in this order:

1. An explicit `--local-site` selector.
2. The Local site whose project path contains your current working directory.
3. If the current directory is outside every Local site, the only site that appears to be running.

If multiple sites are possible, `localwp-cli` refuses to guess.

Explicit selection:

```sh
localwp --local-site "My Site" plugin list
```

The selector may be a Local site name, domain, or site ID:

```sh
localwp --local-site "my-site.local" option get siteurl
```

## Useful wrapper commands

List Local sites:

```sh
localwp --sites
```

Inspect the selected site and resolved Local runtime:

```sh
localwp --doctor
```

Show help:

```sh
localwp --help
```

Open Local's interactive Site Shell intentionally:

```sh
localwp --shell
```

Normal `localwp <wp-cli command>` calls are one-shot commands. You do **not** need to run `exit` afterward.

## Cross-platform Local discovery

The CLI follows Local's platform-specific application data locations:

```text
Windows   %APPDATA%\Local
macOS     ~/Library/Application Support/Local
Linux     ~/.config/Local
```

It reads Local's `sites.json`, resolves the selected site's runtime configuration under Local's application data, then finds the corresponding Local-managed PHP and WP-CLI resources.

The exact Local binary locations remain implementation details; `localwp --doctor` is the supported way to inspect what was discovered on a machine.

## Built for AI coding agents

Interactive shells are inconvenient for autonomous agents. They introduce persistent state and require the agent to know when to leave the shell.

With `localwp-cli`, an agent can simply run:

```sh
localwp plugin list --format=json
localwp core version
localwp cache flush
```

Each command:

1. identifies the Local site;
2. resolves the site-specific Local runtime;
3. executes WP-CLI once;
4. returns stdout/stderr and the WP-CLI process exit code;
5. terminates normally.

This repository includes an agent skill at:

```text
skills/localwp-cli/SKILL.md
```

For Cursor, copy/install it into a project with:

**Windows**

```cmd
scripts\install-cursor-skill.cmd "C:\path\to\project"
```

**macOS / Linux**

```sh
./scripts/install-cursor-skill.sh /path/to/project
```

See [AI agent usage](docs/AI-AGENTS.md).

## CI/CD

This repository includes GitHub Actions for both continuous integration and release delivery.

### CI

Every push and pull request runs:

- Go unit tests on Windows;
- Go unit tests on macOS;
- Go unit tests on Linux;
- `go vet` on every OS;
- native builds on every OS;
- CLI smoke tests;
- `gofmt` verification;
- cross-compilation of all release targets.

### Releases

Pushing a tag such as:

```sh
git tag v1.0.0
git push origin v1.0.0
```

triggers the release workflow. It re-runs tests, builds every supported binary, creates `.zip`/`.tar.gz` packages, generates SHA-256 checksums, and publishes a GitHub Release.

See [CI/CD](docs/CI-CD.md) for details.

## Testing against Local itself

GitHub-hosted runners test the portable CLI logic, parsing, selection, runtime resolution behavior, argument preservation, native compilation, and cross-platform builds.

A hosted CI runner does **not** launch the Local desktop application and provision a real WordPress site. Real Local integration should therefore also be tested before important releases on at least one actual Local installation per platform when practical.

The separation is intentional: deterministic CI tests validate our code; real-application smoke testing validates compatibility with a particular Local release.

## Environment overrides

The following environment variables are primarily useful for diagnostics, development, and CI fixtures:

```text
LOCALWP_DATA_DIR
LOCALWP_RESOURCES_DIR
LOCALWP_PHP_BINARY
LOCALWP_WP_CLI_PHAR
LOCALWP_MYSQL_BIN_DIR
LOCALWP_SKIP_RUNNING_CHECK
```

Most users never need to set them.

## Development

Run the full local checks:

```sh
make check
```

or directly:

```sh
go test ./...
go vet ./...
go build ./cmd/localwp
```

See [Contributing](CONTRIBUTING.md).

## Project status

`localwp-cli` is an independent community utility. It is **not an official Local product** and is not affiliated with WP Engine or the Local team.

Local and its runtime layout can evolve. The project therefore keeps runtime discovery isolated, covered by tests, and visible through `localwp --doctor` so compatibility changes can be diagnosed quickly.

## License

MIT. See [LICENSE](LICENSE).
