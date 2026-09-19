# AI coding agents

`localwp-cli` exists partly to make Local-hosted WordPress projects easier for coding agents to operate safely and predictably.

## Agent rule

When a WordPress project is running in Local, prefer:

```text
localwp <WP-CLI arguments>
```

over direct `wp` execution.

Examples:

```sh
localwp plugin list --format=json
localwp theme list --format=json
localwp option get siteurl
localwp core version
localwp cache flush
```

This rule is the same on Windows, macOS, and Linux.

## Why

The agent should not have to discover or permanently configure:

- which PHP version the Local site uses;
- which Local `php.ini` belongs to the site;
- where Local's WP-CLI is installed;
- which database CLI binaries belong to the site;
- which site shell needs to be opened;
- when an interactive shell needs to be exited.

`localwp` resolves those details per invocation.

## Non-interactive by default

Normal commands are one-shot operations:

```text
agent
  -> localwp plugin list --format=json
  -> Local site/runtime resolution
  -> WP-CLI
  -> stdout/stderr + exit code
  -> agent continues
```

An agent should **not** send `exit` after a normal `localwp` command.

## Site targeting

Prefer running from inside the Local project tree.

When that is not possible:

```sh
localwp --local-site "My Site" plugin list --format=json
```

If the site name is unknown:

```sh
localwp --sites
```

## Diagnostics

When a Local-related WP-CLI command fails:

```sh
localwp --doctor
```

Inspect the selected site, runtime paths and running state before changing project configuration.

## Structured output

Agents should prefer machine-readable WP-CLI output when supported:

```sh
localwp plugin list --format=json
localwp theme list --format=json
localwp user list --format=json
```

## Database safety

Use dry runs before mutations when WP-CLI supports them:

```sh
localwp search-replace "http://old.local" "https://new.local" --dry-run
```

Do not drop, reset, overwrite or import over a database unless the task explicitly requires it.

## Interactive Site Shell

`localwp --shell` should be exceptional for agents. It is intended for commands that genuinely require a persistent interactive Local shell.

For ordinary WP-CLI work, keep using one-shot `localwp ...` commands.

## Reusable skill

The repository contains:

```text
skills/localwp-cli/SKILL.md
```

Install it into a Cursor project:

Windows:

```cmd
scripts\install-cursor-skill.cmd "C:\path\to\project"
```

macOS/Linux:

```sh
./scripts/install-cursor-skill.sh /path/to/project
```
