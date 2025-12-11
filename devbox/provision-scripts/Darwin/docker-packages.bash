#!/usr/bin/env bash
set -e

if [[ ! $(brew list | grep -w docker) ]]; then
    brew install docker
    brew install docker-buildx
    brew install colima

    mkdir -p "$HOME/.docker"
    echo -e "{\n\t\"cliPluginsExtraDirs\": [ \"$HOMEBREW_PREFIX/lib/docker/cli-plugins\" ]\n}" > "$HOME/.docker/config.json"

    colima start
    brew services start colima
fi
