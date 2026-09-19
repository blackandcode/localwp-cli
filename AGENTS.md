# Agent instructions

This repository implements a cross-platform CLI companion for Local.

When changing the project:

- Preserve the public interface `localwp <WP-CLI arguments>` on Windows, macOS and Linux.
- Keep ordinary WP-CLI execution non-interactive.
- Preserve WP-CLI argument boundaries and the underlying process exit code.
- Do not bundle Local, PHP, WordPress, MySQL/MariaDB, or WP-CLI binaries.
- Treat Local runtime paths as implementation details that may change between Local releases.
- Keep platform-specific discovery isolated and testable.
- Do not make the generated Local Site Shell entry a requirement for normal one-shot WP-CLI commands.
- Keep `--shell` as an explicit interactive-only path.
- Update README, architecture docs, installation docs and `skills/localwp-cli/SKILL.md` when user-visible behavior changes.
- Run `go test ./...`, `go vet ./...`, `gofmt`, and native builds before merging.
- Preserve CI coverage for Windows, macOS and Linux.
- When runtime discovery changes, add fixture/unit coverage and smoke-test against a real Local installation when practical.

The reusable consumer-facing coding-agent skill is `skills/localwp-cli/SKILL.md`.
