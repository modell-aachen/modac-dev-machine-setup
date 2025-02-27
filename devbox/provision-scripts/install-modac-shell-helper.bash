#!/usr/bin/env bash
set -e

if [[ -z "$NEXUS_BOT_TOKEN" ]]; then
    echo -e "You need to set ${BI_RED}NEXUS_BOT_TOKEN$NC environment variable in $HOME/.env"
    exit 1
fi

if [[ "$( pip3 config list | grep global.bread-system-packages )" != *true* ]]; then
    echo "Configuring pip to break system packages"
    pip3 config set global.break-system-packages true
fi

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
pip3 install --break-system-packages --force-reinstall --user "$filename"

popd
rm -r "$dir"
