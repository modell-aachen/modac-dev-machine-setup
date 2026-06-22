#!/usr/bin/env bash
set -euo pipefail

REPO_ROOT="$(git rev-parse --show-toplevel)"
PLUGIN_JSON="$REPO_ROOT/devbox/plugins/modac/plugin.json"

if [ ! -f "$PLUGIN_JSON" ]; then
  echo "ERROR: plugin.json not found at $PLUGIN_JSON" >&2
  exit 1
fi

LATEST_HASH="$(git rev-parse HEAD)"
# Portable across macOS (BSD) and Linux (GNU): avoid grep -oP and sed -i quirks.
# There are multiple &rev= refs (machine + google-cloud-sdk-gke); bump them all.
CURRENT_HASH="$(grep -oE 'rev=[a-f0-9]{40}' "$PLUGIN_JSON" | head -1 | cut -d= -f2)"

if [ "$LATEST_HASH" = "$CURRENT_HASH" ]; then
  echo "Already up to date (${LATEST_HASH:0:7})"
  exit 0
fi

tmp="$(mktemp)"
sed "s/rev=${CURRENT_HASH}/rev=${LATEST_HASH}/g" "$PLUGIN_JSON" > "$tmp"
cat "$tmp" > "$PLUGIN_JSON"
rm -f "$tmp"
echo "Updated machine hash: ${CURRENT_HASH:0:7} -> ${LATEST_HASH:0:7}"
