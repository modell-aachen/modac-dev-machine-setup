#!/usr/bin/env bash
set -e

. ./provision-scripts/helper.bash

for shell in bash zsh; do
    install_completion "$shell" flux 2_5_1
    install_completion "$shell" op 2_30_3
    install_completion "$shell" helm 3_17_3
    install_completion "$shell" devspace 6_3_15
    install_completion "$shell" kubectl 1_32_3
done

