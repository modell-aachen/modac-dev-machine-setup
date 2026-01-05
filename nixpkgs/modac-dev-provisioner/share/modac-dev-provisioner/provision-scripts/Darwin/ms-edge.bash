#!/usr/bin/env bash
set -e

if [[ ! $(brew list | grep -w microsoft-edge) ]]; then
    brew install microsoft-edge
fi