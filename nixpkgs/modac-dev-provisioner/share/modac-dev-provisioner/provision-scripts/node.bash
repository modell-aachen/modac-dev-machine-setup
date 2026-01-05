#!/usr/bin/env bash
set -e

devbox_path=$(devbox global path)
yarn_path="$devbox_path/.devbox/virtenv/nodejs/corepack-bin/yarn"

if [ ! -f "$yarn_path" ]; then
    corepack install -g yarn@latest
fi
