# Architecture

`localwp-cli` is intentionally small and uses only the Go standard library.

## Execution flow

```text
localwp <wp-cli args>
        |
        v
resolve Local application data directory
        |
        v
read sites.json
        |
        v
select site
  --local-site
  current directory
  single running site
        |
        v
resolve Local runtime
  PHP binary
  site php.ini
  bundled wp-cli.phar
  database binaries/config
        |
        v
execute
PHP -c <site php.ini> wp-cli.phar --path=<site>/app/public <args...>
        |
        v
return WP-CLI exit code
```

## Platform data directories

```text
Windows   %APPDATA%\Local
macOS     ~/Library/Application Support/Local
Linux     $XDG_CONFIG_HOME/Local or ~/.config/Local
```

## Runtime discovery

Local can ship runtime components with the application and can also download Lightning Services into its user application-data directory.

The resolver searches both categories instead of assuming that every PHP/MySQL version is bundled in one fixed application path.

Platform-specific Local resource roots are isolated in `internal/localwp/runtime.go`.

## WP-CLI execution

Normal commands do not enter Local's interactive Site Shell.

Instead the executable directly starts the Local-managed PHP binary with:

```text
-c <site php.ini>
<Local wp-cli.phar>
--path=<WordPress root>
<user arguments>
```

This design preserves argument boundaries and process exit codes without nesting a persistent shell.

## Interactive shell

`localwp --shell` is a separate feature. It invokes Local's generated `ssh-entry/<site-id>.bat` or `.sh` entry when available.

The shell-entry file is deliberately not required for normal WP-CLI commands.

## Testability

Runtime lookup supports environment overrides so resolver behavior can be tested without hard-coding a developer machine:

```text
LOCALWP_DATA_DIR
LOCALWP_RESOURCES_DIR
LOCALWP_PHP_BINARY
LOCALWP_WP_CLI_PHAR
LOCALWP_MYSQL_BIN_DIR
```

The production default remains automatic Local discovery.

## Release installation and agent routing

The `installation/install-release.ps1` and `installation/install-release.sh` helpers download the latest published GitHub Release for a supported platform, verify its archive against that release's `SHA256SUMS.txt`, and install only the CLI in the user's bin directory. Windows updates User PATH; the POSIX helper updates Bash/zsh startup files. They do not bundle or install Local or its runtimes. Root-level `install-localwp.*` scripts remain source-build tools for contributors.

The consumer agent skill requires `localwp` only when the current environment has Local installed, the intended Local site is active, and the wrapper is on PATH. It allows a configured `wp` fallback when the wrapper is unavailable, including CI and remote environments. This is agent guidance; the CLI's site selection and execution behavior are unchanged.
