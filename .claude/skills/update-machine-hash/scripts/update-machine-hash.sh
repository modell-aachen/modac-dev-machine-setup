#!/usr/bin/env bash
set -euo pipefail

REPO_ROOT="$(git rev-parse --show-toplevel)"
PLUGIN_JSONS=(
  "$REPO_ROOT/devbox/plugins/modac/plugin.json"
  "$REPO_ROOT/devbox/plugins/modac-service/plugin.json"
)

LATEST_HASH="$(git rev-parse HEAD)"

for plugin_json in "${PLUGIN_JSONS[@]}"; do
  if [ ! -f "$plugin_json" ]; then
    echo "ERROR: plugin.json not found at $plugin_json" >&2
    exit 1
  fi

  # Portable across macOS (BSD) and Linux (GNU): avoid grep -oP and sed -i quirks.
  # There are multiple &rev= refs (machine + google-cloud-sdk-gke); bump them all.
  # Hashes are read per file since the plugins can drift apart.
  CURRENT_HASH="$(grep -oE 'rev=[a-f0-9]{40}' "$plugin_json" | head -1 | cut -d= -f2)"

  if [ "$LATEST_HASH" = "$CURRENT_HASH" ]; then
    echo "Already up to date (${LATEST_HASH:0:7}): $plugin_json"
    continue
  fi

  tmp="$(mktemp)"
  sed "s/rev=${CURRENT_HASH}/rev=${LATEST_HASH}/g" "$plugin_json" > "$tmp"
  cat "$tmp" > "$plugin_json"
  rm -f "$tmp"
  echo "Updated machine hash: ${CURRENT_HASH:0:7} -> ${LATEST_HASH:0:7}: $plugin_json"
done
