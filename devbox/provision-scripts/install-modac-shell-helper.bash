#!/usr/bin/env bash
set -e

dir=/tmp/modac-shell-helper
pyproject="pyproject.toml"

mkdir -p "$dir"
pushd "$dir"

curl -s -f -u "bot-ro:${NEXUS_BOT_TOKEN}" \
    "https://nexus.modac.cloud/repository/modac-shell-helper/latest/pyproject.toml" \
    --output "$pyproject"
filename="modac_shell_helper-$(poetry version -s)-py3-none-any.whl"
curl -s -f -u "bot-ro:${NEXUS_BOT_TOKEN}" \
    "https://nexus.modac.cloud/repository/modac-shell-helper/latest/modac-shell-helper.whl" \
    --output "$filename"
echo "Installing $filename"
pip3 install --force-reinstall --user "$filename"

popd
rm -r "$dir"
