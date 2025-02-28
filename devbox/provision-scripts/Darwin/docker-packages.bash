#!/usr/bin/env bash
set -e

if [[ ! $(brew list | grep -w docker) ]]; then
    brew install docker
fi
