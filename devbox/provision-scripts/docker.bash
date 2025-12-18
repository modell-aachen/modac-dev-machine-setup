#!/usr/bin/env bash
set -e

source "$(dirname "$0")/helper.bash"

if [ ! -f $HOME/.docker_buildx_builder_created ]; then
    log_info "Creating docker buildx builder"
    sudo docker buildx create --use
    touch $HOME/.docker_buildx_builder_created

    log_warn "Please logout and login again to use docker without sudo"
fi
