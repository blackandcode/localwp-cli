#!/usr/bin/env sh
set -eu

if ! command -v go >/dev/null 2>&1; then
  echo "localwp install: Go is required when installing from source." >&2
  echo "Install Go, or download a prebuilt binary from the GitHub Releases page." >&2
  exit 1
fi

SCRIPT_DIR=$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)
INSTALL_DIR="${LOCALWP_INSTALL_DIR:-$HOME/.local/bin}"
BINARY="$INSTALL_DIR/localwp"

mkdir -p "$INSTALL_DIR"
cd "$SCRIPT_DIR"
go build -trimpath -o "$BINARY" ./cmd/localwp
chmod +x "$BINARY"

echo "Installed localwp to:"
echo "  $BINARY"
echo
case ":$PATH:" in
  *":$INSTALL_DIR:"*)
    echo "Run:"
    echo "  localwp --version"
    echo "  localwp --sites"
    ;;
  *)
    echo "Add this directory to PATH in your shell profile:"
    echo "  export PATH=\"$INSTALL_DIR:\$PATH\""
    echo
    echo "Then open a new terminal and run:"
    echo "  localwp --version"
    ;;
esac
