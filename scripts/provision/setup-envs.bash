#!/usr/bin/env bash
set -e

base_path="$(cd "$(dirname "$0")" && pwd)/.."
config_file="$( devbox global path)/devbox.json"
tmp_file="$( devbox global path)/tmp.json"

[[ -f "$tmp_file" ]] && rm "$tmp_file"

if [[ "$( jq .env_from $config_file )" != *"$HOME/.secrets/.env"* ]]; then
    echo "Adding env_from to $config_file"
    jq  -c ". + {\"env_from\": \"$HOME/.secrets/.env\"}" "$config_file" | jq . > "$tmp_file"
    mv "$tmp_file" "$config_file"
fi

echo "Merging op_secrets_tpl in $config_file with template"

jq -c '
  . as $cfg
  | (input | .op_secrets_tpl) as $tpl
  | $cfg
  | .op_secrets_tpl = (
      ($cfg.op_secrets_tpl // {})      # existing map (maybe empty)
      + ( $tpl | to_entries            # template entries
          | map(
              if (.key | in($cfg.op_secrets_tpl // {}))
              then empty               # skip keys that already exist
              else { (.key): .value }  # add new keys
              end
            )
          | add // {}
        )
    )
' "$config_file" "$base_path/templates/devbox.json" \
  | jq . > "$tmp_file"

mv "$tmp_file" "$config_file"

if [[ $( cat "$HOME/.env" ) == *"export "* ]]; then
    echo "Removing export from $HOME/.env"
    sed -i 's/^export //' "$HOME/.env"
fi

mkdir -p "$HOME/.secrets"
jq -r '.op_secrets_tpl | to_entries | .[] | .key + "=" + (.value | @sh)' "$config_file" >| $HOME/.secrets/env.tpl
op inject \
    --in-file $HOME/.secrets/env.tpl --out-file $HOME/.secrets/.env --force
