#!/usr/bin/env bash
set -e

config_file="$( devbox global path)/devbox.json"
tmp_file="$( devbox global path)/tmp.json"

asdf plugin add erlang
asdf plugin add elixir

asdf_dir=$(asdf info | grep ASDF_DIR | sed 's/ASDF_DIR=//')

if ! grep -q "$asdf_dir" "$config_file" ; then
    jq -c ".env += {ASDF_DIR: \"$asdf_dir\"}" "$config_file" | jq . >| "$tmp_file"
    mv "$tmp_file" "$config_file"
fi
