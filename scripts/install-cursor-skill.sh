#!/usr/bin/env sh
set -eu
if [ "$#" -ne 1 ]; then
  echo "Usage: $0 /path/to/project" >&2
  exit 2
fi
SCRIPT_DIR=$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)
REPO_ROOT=$(CDPATH= cd -- "$SCRIPT_DIR/.." && pwd)
DEST="$1/.cursor/skills/localwp-cli"
mkdir -p "$DEST"
cp "$REPO_ROOT/skills/localwp-cli/SKILL.md" "$DEST/SKILL.md"
echo "Installed LocalWP CLI skill to $DEST/SKILL.md"
