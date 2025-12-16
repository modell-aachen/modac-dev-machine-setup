#!/usr/bin/env bash
set -e

pyenv install 3 -s
pyenv global 3

if [[ "$( pip3 config list | grep global.bread-system-packages )" != *true* ]]; then
    echo "Configuring pip to break system packages"
    pip3 config set global.break-system-packages true
fi

pip3 install --upgrade pip

function pyenv_init_hook() {
    local shell=$1
    local shell_path="$HOME/.${shell}rc"
    local cmd='eval "$(pyenv init - '"$shell"')"'

    if [[ -f "$shell_path" && -z "$( cat "$shell_path" | grep "$cmd" )" ]]; then
        echo "Setting up pyenv init hook for $shell"

        echo "$cmd" >> "$shell_path"
    fi
}

for shell in bash zsh; do
    pyenv_init_hook "$shell"
done
