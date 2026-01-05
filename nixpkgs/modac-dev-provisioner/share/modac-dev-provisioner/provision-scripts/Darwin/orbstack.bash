#!/usr/bin/env bash

if [[ "$(orbctl status)" == "Stopped" ]]; then
    echo "Starting OrbStack..."
    orbctl start

    echo "Login to OrbStack .."
    orbctl login
fi
