#!/usr/bin/env bash
set -e

source "$(dirname "$0")/helper.bash"

if gh auth status >/dev/null 2>&1; then
    log_info "GitHub CLI is already authenticated"
else
    log_info "Logging into GitHub CLI"
    gh auth login --web
fi