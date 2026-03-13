#!/usr/bin/env bash
set -euo pipefail

REPO="Subodh8/ClaudeFlux"
BINARY="claudeflux"
INSTALL_DIR="${CLAUDEFLUX_INSTALL_DIR:-/usr/local/bin}"

OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)
case $ARCH in
  x86_64)  ARCH="amd64" ;;
  aarch64|arm64) ARCH="arm64" ;;
  *) echo "Unsupported architecture: $ARCH" && exit 1 ;;
esac

LATEST=$(curl -sSL "https://api.github.com/repos/$REPO/releases/latest" \
  | grep '"tag_name"' | sed 's/.*"tag_name": "\(.*\)".*/\1/')

ARCHIVE="${BINARY}-${OS}-${ARCH}.tar.gz"
URL="https://github.com/$REPO/releases/download/$LATEST/$ARCHIVE"

echo "→ Installing ClaudeFlux $LATEST ($OS/$ARCH)"

TMP=$(mktemp -d)
trap "rm -rf $TMP" EXIT

curl -sSL "$URL" -o "$TMP/$ARCHIVE"
tar -xzf "$TMP/$ARCHIVE" -C "$TMP"

chmod +x "$TMP/$BINARY"

"$TMP/$BINARY" version > /dev/null 2>&1 || {
  echo "✗ Binary verification failed."
  exit 1
}

if [ -w "$INSTALL_DIR" ]; then
  mv "$TMP/$BINARY" "$INSTALL_DIR/$BINARY"
else
  sudo mv "$TMP/$BINARY" "$INSTALL_DIR/$BINARY"
fi

echo ""
echo "✓ ClaudeFlux $LATEST installed to $INSTALL_DIR/$BINARY"
echo ""
echo "  Next steps:"
echo "  1. Install claude CLI: https://docs.anthropic.com/claude-code"
echo "  2. Run: claudeflux run workflow.yaml --dashboard"
echo "  3. Docs: https://claudeflux.dev/docs"
