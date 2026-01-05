#!/usr/bin/env bash
set -e

if [[ ! $(brew list | grep -w docker) ]]; then
    brew install \
        docker \
        docker-buildx \
        docker-compose \
        docker-completion

    docker context use orbstack 
fi

dockerConfigPath="$HOME/.docker/config.json"
if [[ ! -f "$dockerConfigPath" || -z "$( cat "$dockerConfigPath" | grep "cliPluginsExtraDirs" )" ]]; then
    mkdir -p "$HOME/.docker"
    echo -e "{\n\t\"cliPluginsExtraDirs\": [ \"$HOMEBREW_PREFIX/lib/docker/cli-plugins\" ]\n}" > "$dockerConfigPath"
fi