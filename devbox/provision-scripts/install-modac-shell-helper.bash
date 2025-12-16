#!/usr/bin/env bash
set -e

source "$(dirname "$0")/helper.bash"

eval "$(devbox global shellenv --init-hook)"
if [[ -z "$NEXUS_BOT_TOKEN" ]]; then
    log_error "You need to set NEXUS_BOT_TOKEN environment variable in $HOME/.env"
    exit 1
fi

dir=/tmp/modac-shell-helper
pyproject="pyproject.toml"

mkdir -p "$dir"
pushd "$dir" > /dev/null

log_info "Downloading modac-shell-helper from Nexus"
curl -s -f -u "bot-ro:${NEXUS_BOT_TOKEN}" \
    "https://nexus.modac.cloud/repository/modac-shell-helper/latest/pyproject.toml" \
    --output "$pyproject"
filename="modac_shell_helper-$(poetry version -s)-py3-none-any.whl"
curl -s -f -u "bot-ro:${NEXUS_BOT_TOKEN}" \
    "https://nexus.modac.cloud/repository/modac-shell-helper/latest/modac-shell-helper.whl" \
    --output "$filename"
log_info "Installing $filename"
pip3 install --break-system-packages --force-reinstall --user "$filename"

popd > /dev/null
rm -r "$dir"
