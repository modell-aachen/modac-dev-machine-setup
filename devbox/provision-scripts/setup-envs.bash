#!/usr/bin/env bash
set -e

source "$(dirname "$0")/helper.bash"

config_file="$( devbox global path)/devbox.json"
tmp_file="$( devbox global path)/tmp.json"

[[ -f "$tmp_file" ]] && rm "$tmp_file"

if [[ "$( jq .env_from $config_file )" != *"$HOME/.secrets/.env"* ]]; then
    log_info "Adding env_from to $config_file"
    jq  -c ". + {\"env_from\": \"$HOME/.secrets/.env\"}" "$config_file" | jq . > "$tmp_file"
    mv "$tmp_file" "$config_file"
fi

log_info "Updating op_secrets_tpl in $config_file"
op_tpl=$( jq -cM '.op_secrets_tpl' "$PROVISIONER_DIRECTORY/devbox/templates/devbox.json" )
jq  -c ". + {\"op_secrets_tpl\": $op_tpl}" "$config_file" | jq . >| "$tmp_file"
mv "$tmp_file" "$config_file"

if [[ $( cat "$HOME/.env" ) == *"export "* ]]; then
    log_info "Removing export from $HOME/.env"
    sed -i 's/^export //' "$HOME/.env"
fi

log_info "Injecting secrets from 1Password"
mkdir -p $HOME/.secrets
jq -r '.op_secrets_tpl | to_entries | .[] | .key + "=" + (.value | @sh)' "$config_file" >| $HOME/.secrets/env.tpl
op inject --in-file $HOME/.secrets/env.tpl --out-file $HOME/.secrets/.env --force
