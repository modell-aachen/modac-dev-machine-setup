#!/usr/bin/env bash
set -e

if [ ! -f $HOME/.docker_buildx_builder_created ]; then
    docker buildx create --use
    touch $HOME/.docker_buildx_builder_created
fi
