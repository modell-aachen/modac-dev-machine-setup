#!/usr/bin/env bash
set -euo pipefail

REPO_ROOT="$(git rev-parse --show-toplevel)"
PLUGIN_JSON="$REPO_ROOT/devbox/plugins/modac/plugin.json"

if [ ! -f "$PLUGIN_JSON" ]; then
  echo "ERROR: plugin.json not found at $PLUGIN_JSON" >&2
  exit 1
fi

LATEST_HASH="$(git rev-parse HEAD)"
CURRENT_HASH="$(grep -oP '(?<=&rev=)[a-f0-9]+' "$PLUGIN_JSON")"

if [ "$LATEST_HASH" = "$CURRENT_HASH" ]; then
  echo "Already up to date (${LATEST_HASH:0:7})"
  exit 0
fi

sed -i "s/&rev=${CURRENT_HASH}/\&rev=${LATEST_HASH}/" "$PLUGIN_JSON"
echo "Updated machine hash: ${CURRENT_HASH:0:7} -> ${LATEST_HASH:0:7}"
