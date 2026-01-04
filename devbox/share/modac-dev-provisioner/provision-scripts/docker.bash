#!/usr/bin/env bash
set -e

if [ ! -f $HOME/.docker_buildx_builder_created ]; then
    sudo docker buildx create --use
    touch $HOME/.docker_buildx_builder_created

    echo -e "Please ${BI_RED}logout and login again$NC to use docker without sudo"
fi
