#!/usr/bin/env bash
set -e

source "$(dirname "$0")/../helper.bash"

if [[ ! $(brew list | grep -w docker) ]]; then
    log_info "Installing Docker packages via Homebrew"
    brew install \
        docker \
        docker-buildx \
        docker-compose \
        docker-completion

    log_info "Setting Docker context to orbstack"
    docker context use orbstack
fi

dockerConfigPath="$HOME/.docker/config.json"
if [[ ! -f "$dockerConfigPath" || -z "$( cat "$dockerConfigPath" | grep "cliPluginsExtraDirs" )" ]]; then
    log_info "Configuring Docker CLI plugins directory"
    mkdir -p "$HOME/.docker"
    echo -e "{\n\t\"cliPluginsExtraDirs\": [ \"$HOMEBREW_PREFIX/lib/docker/cli-plugins\" ]\n}" > "$dockerConfigPath"
fi