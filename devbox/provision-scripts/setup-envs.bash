#!/usr/bin/env bash
set -e

config_file="$( devbox global path)/devbox.json"
tmp_file="$( devbox global path)/tmp.json"

[[ -f "$tmp_file" ]] && rm "$tmp_file"

if [[ "$( jq .env_from $config_file )" != *"$HOME/.secrets/.env"* ]]; then
    echo "Adding env_from to $config_file"
    jq  -c ". + {\"env_from\": \"$HOME/.secrets/.env\"}" "$config_file" | jq . > "$tmp_file"
    mv "$tmp_file" "$config_file"
fi

if [[ "$( jq .op_secrets_tpl $config_file )" == "null" ]]; then
    echo "Adding op_secrets_tpl to $config_file"
    op_tpl=$( jq -cM '.op_secrets_tpl' "$PROVISIONER_DIRECTORY/devbox/templates/devbox.json" )
    jq  -c ". + {\"op_secrets_tpl\": $op_tpl}" "$config_file" | jq . >| "$tmp_file"
    mv "$tmp_file" "$config_file"
fi

if [[ $( cat "$HOME/.env" ) == *"export "* ]]; then
    echo "Removing export from $HOME/.env"
    sed -i 's/^export //' "$HOME/.env"
fi

mkdir -p $HOME/.secrets
jq -r '.op_secrets_tpl | to_entries | .[] | .key + "=" + (.value | @sh)' "$config_file" >| $HOME/.secrets/env.tpl
op inject --in-file $HOME/.secrets/env.tpl --out-file $HOME/.secrets/.env --force
