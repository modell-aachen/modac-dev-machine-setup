#!/usr/bin/env bash
set -e

config_file="$( devbox global path)/devbox.json"
tmp_file="$( devbox global path)/tmp.json"

[[ -f "$tmp_file" ]] && rm "$tmp_file"

if [[ "$( jq .env_from $config_file )" == *"null"* ]]; then
    echo "Adding env_from to $config_file"
    jq  -c ". + {\"env_from\": \"$HOME/.env\"}" "$(devbox global path)/devbox.json" | jq . > "$tmp_file"
    mv "$tmp_file" "$config_file"
fi

if [[ $( cat "$HOME/.env" ) == *"export "* ]]; then
    echo "Removing export from $HOME/.env"
    sed -i 's/^export //' "$HOME/.env"
fi
