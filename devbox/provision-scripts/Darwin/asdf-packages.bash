#!/usr/bin/env bash
set -e

if [[ ! $(brew list | grep -w asdf) ]]; then
    brew install asdf
fi
