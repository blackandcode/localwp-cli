---
name: localwp-cli
description: Run WP-CLI operations targeting Local WordPress (LocalWP, Local by Flywheel) on Windows, macOS, or Linux. Trigger for wp commands, plugin/theme management, core, options, cache, users, database export or search-replace, and custom aws-backend-api request/logs commands in Local projects. Detect an installed Local app and an active target site, use localwp when available on PATH, and use wp in configured environments without the wrapper such as CI/CD, remote containers, or headless servers.
---

# LocalWP CLI

Use this workflow for WP-CLI work targeting a Local WordPress development site. A WordPress repository alone does not establish that Local is installed or that the target site is active.

## Runtime detection and command choice

1. Check whether the wrapper is available on PATH:
   - Windows: `where.exe localwp` (use the executable name explicitly in PowerShell).
   - macOS/Linux: `which localwp` (or `command -v localwp` if `which` is unavailable).
2. Confirm Local is installed in the current execution environment and the intended project is registered with it. When the wrapper is available, run `localwp --sites`, match the site's path/name/domain to the requested project, and check its running state. Use `localwp --local-site "Site Name" --doctor` to inspect the selected runtime and running state. The mere presence of `localwp` on PATH is not proof of an active Local site.
3. **When Local is installed, the intended site is active, and the wrapper is available, always use `localwp <args>` for that site's WP-CLI operations.** Prefer running inside the target site tree; use `--local-site` when needed to avoid selecting a different running site.
4. **When `localwp` is unavailable, fall back to `wp <args>`** if WP-CLI is available and configured for the intended WordPress installation. This includes CI/CD, remote containers, and headless servers. Check `where.exe wp` on Windows or `which wp` / `command -v wp` on POSIX and confirm the target working directory or explicit `--path`. A host's Local installation does not make a container or remote server a Local environment. If neither command is usable, report the missing prerequisite; do not automatically install software.
5. For a known Local site that is stopped, start it within the authorized workflow or ask the user to start it, then recheck. Do not switch to `wp` merely to bypass a stopped site, a selection error, or a broken Local runtime. For non-Local targets, use their configured `wp` command even if `localwp` happens to be installed.

Keep ordinary WP-CLI execution non-interactive. Preserve argument boundaries and check the process exit code. Do not enter a Local Site Shell for ordinary commands.

## Examples

```text
localwp core version
localwp plugin list
localwp plugin list --format=json
localwp theme list --format=json
localwp option get siteurl
localwp cache flush
localwp plugin activate my-plugin
localwp aws-backend-api request --help
localwp aws-backend-api logs --help
```

For a site that registers the custom `aws-backend-api` command, inspect its help first, then run `localwp aws-backend-api request <documented arguments>` or `localwp aws-backend-api logs <documented arguments>`. These are plugin-provided commands, not built-in WordPress or wrapper commands; do not invent endpoint, payload, or log flags.

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
