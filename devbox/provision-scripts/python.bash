#!/usr/bin/env bash
set -e

pyenv install 3 -s
pyenv global 3

pip3 install --upgrade pip

. ./provision-scripts/helper.bash

for shell in bash zsh; do
    pyenv_envs "$shell"
    source_path "$shell"
done
