#!/usr/bin/env bash
set -e

source "$(dirname "$0")/helper.bash"

if [ -f "$HOME/.local/bin/poetry" ] && [[ "$($HOME/.local/bin/poetry --version)" == *"$POETRY_VERSION"* ]]; then
    log_info "poetry in version $POETRY_VERSION is already installed"
else
    log_info "Installing poetry version $POETRY_VERSION"
    curl -sSL https://install.python-poetry.org | python3 - --version $POETRY_VERSION
fi
