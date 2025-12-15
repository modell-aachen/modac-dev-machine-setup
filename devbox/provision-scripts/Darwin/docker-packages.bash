#!/usr/bin/env bash
set -e

if [[ ! $(brew list | grep -w docker) ]]; then
    brew install \
        docker \
        docker-buildx 

    mkdir -p "$HOME/.docker"
    echo -e "{\n\t\"cliPluginsExtraDirs\": [ \"$HOMEBREW_PREFIX/lib/docker/cli-plugins\" ]\n}" > "$HOME/.docker/config.json"

    docker context use orbstack 
fi
