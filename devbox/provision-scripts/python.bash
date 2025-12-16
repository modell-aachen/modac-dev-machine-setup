#!/usr/bin/env bash
set -e

source "$(dirname "$0")/helper.bash"

log_info "Installing Python 3 via pyenv"
pyenv install 3 -s
pyenv global 3

if [[ "$( pip3 config list | grep global.bread-system-packages )" != *true* ]]; then
    log_info "Configuring pip to break system packages"
    pip3 config set global.break-system-packages true
fi

log_info "Upgrading pip"
pip3 install --upgrade pip

function pyenv_init_hook() {
    local shell=$1
    local shell_path="$HOME/.${shell}rc"
    local cmd='eval "$(pyenv init - '"$shell"')"'

    if [[ -f "$shell_path" && -z "$( cat "$shell_path" | grep "$cmd" )" ]]; then
        log_info "Setting up pyenv init hook for $shell"

        echo "$cmd" >> "$shell_path"
    fi
}

for shell in bash zsh; do
    pyenv_init_hook "$shell"
done
