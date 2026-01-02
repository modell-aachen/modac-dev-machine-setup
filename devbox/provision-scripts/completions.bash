#!/usr/bin/env bash
set -e

function install_completion() {
    local shell=$1
    local cmd=$2
    local version=$3
    local shell_path="$HOME/.${shell}rc"
    local completions_path="$HOME/.${shell}_completions"
    local cmd_completion_path="$completions_path/${cmd}_$version.sh"


    if [[ -f "$shell_path" && ! -f "$cmd_completion_path" ]]; then
        mkdir -p "$completions_path"
        rm -f "$completions_path/$cmd"*".sh"

        echo "Installing $cmd completion for $shell under $cmd_completion_path"
        "$cmd" completion "$shell" > "$cmd_completion_path"
    fi
}

for shell in bash zsh; do
    install_completion "$shell" flux 2_5_1
    install_completion "$shell" op 2_30_3
    install_completion "$shell" helm 3_17_3
    install_completion "$shell" devspace 6_3_15
    install_completion "$shell" kubectl 1_32_3
done

