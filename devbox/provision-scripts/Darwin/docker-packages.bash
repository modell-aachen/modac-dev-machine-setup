#!/usr/bin/env bash
set -e

if [[ ! $(brew list | grep -w docker) ]]; then
    brew install docker
    brew install docker-buildx
    brew install colima
fi
