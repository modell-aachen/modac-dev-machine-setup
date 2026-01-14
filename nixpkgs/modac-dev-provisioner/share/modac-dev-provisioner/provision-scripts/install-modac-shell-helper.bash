#!/usr/bin/env bash
set -e

if [  -z "$( uv tool list | grep modac-shell-helper )" ]; then
    if [ -z "$(gh auth token)"  ]; then
        echo -e "You need to setup ${BI_RED}gh auth login$NC!"
    else
        export GIT_USERNAME="x-access-token"
        export GIT_PASSWORD="$(gh auth token)"
        export GIT_ASKPASS=/usr/bin/printf
        export GIT_TERMINAL_PROMPT=0
        uv tool install "git+https://github.com/modell-aachen/modac-shell-helper.git@main" --force
    fi
else
    uv tool upgrade modac-shell-helper
fi
