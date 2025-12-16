#!/usr/bin/env bash
set -e 

if gh auth status >/dev/null 2>&1; then
    echo "GitHub CLI is already authenticated."
else
    echo "Logging into GitHub CLI..."
    gh auth login --web
fi