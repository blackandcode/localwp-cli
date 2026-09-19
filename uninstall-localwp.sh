#!/usr/bin/env sh
set -eu
INSTALL_DIR="${LOCALWP_INSTALL_DIR:-$HOME/.local/bin}"
BINARY="$INSTALL_DIR/localwp"
if [ -f "$BINARY" ]; then
  rm -f "$BINARY"
  echo "Removed $BINARY"
else
  echo "localwp is not installed at $BINARY"
fi
