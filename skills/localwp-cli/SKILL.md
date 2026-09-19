---
name: localwp-cli
description: Use WP-CLI safely and non-interactively against WordPress sites managed by Local on Windows, macOS, or Linux. Use when a project runs in Local and the agent needs WordPress inspection, plugin/theme operations, option/cache operations, database-aware WP-CLI commands, or Local runtime diagnostics.
---

# LocalWP CLI

Use `localwp` for WP-CLI operations when the WordPress development site is managed by [Local](https://localwp.com/).

This rule applies on Windows, macOS, and Linux.

## Core rule

Use:

```text
localwp <wp-cli arguments>
```

instead of invoking `wp` directly when the task targets the Local-managed WordPress project.

Do not enter an interactive Local Site Shell for ordinary WP-CLI work.

## Examples

```text
localwp core version
localwp plugin list
localwp plugin list --format=json
localwp theme list --format=json
localwp option get siteurl
localwp cache flush
localwp plugin activate my-plugin
```

Treat everything after `localwp` as normal WP-CLI arguments except the documented `localwp` wrapper options.

## Site selection

Prefer automatic current-directory detection. Run commands with the working directory somewhere inside the Local site tree.

If selection is ambiguous, use:

```text
localwp --local-site "Site Name" plugin list
```

The selector can be a Local site name, domain, or site ID.

Discover sites with:

```text
localwp --sites
```

## Diagnostics

If a command fails because the Local site or runtime cannot be resolved, run:

```text
localwp --doctor
```

Check:

- selected Local site;
- WordPress root;
- site running state;
- Local PHP binary/version;
- site-specific `php.ini`;
- Local bundled WP-CLI;
- database runtime information.

Start the target site in Local when a command needs the database or other site services.

## Non-interactive behavior

Each normal invocation is one-shot:

1. detect the site;
2. resolve Local's site-specific runtime;
3. execute WP-CLI once;
4. return stdout/stderr and WP-CLI's exit code;
5. terminate.

Do not send a separate `exit` command after normal `localwp` commands.

## Structured output

Prefer machine-readable output when available:

```text
localwp plugin list --format=json
localwp theme list --format=json
localwp user list --format=json
```

Always check the exit code before treating a command as successful.

## Database safety

Use dry runs before mutations when supported:

```text
localwp search-replace "http://old.local" "https://new.local" --dry-run
```

Do not drop, reset, overwrite, or destructively import a database unless the task explicitly requires it.

## Interactive shell

Only use:

```text
localwp --shell
```

when a persistent Local Site Shell is genuinely required.

Exit that interactive shell normally with `exit`.
