#!/usr/bin/env bash
set -e

if [ -f "$HOME/.local/bin/poetry" ] && [[ "$($HOME/.local/bin/poetry --version)" == *"$POETRY_VERSION"* ]]; then
    echo "poetry in version $POETRY_VERSION is already installed"
else
    curl -sSL https://install.python-poetry.org | python3 - --version $POETRY_VERSION
fi
