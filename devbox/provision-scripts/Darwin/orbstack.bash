#!/usr/bin/env bash

source "$(dirname "$0")/../helper.bash"

if [[ "$(orbctl status)" == "Stopped" ]]; then
    log_info "Starting OrbStack"
    orbctl start

    log_info "Logging in to OrbStack"
    orbctl login
fi
