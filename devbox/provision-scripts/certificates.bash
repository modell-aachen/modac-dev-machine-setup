#!/usr/bin/env bash
set -e

source "$(dirname "$0")/helper.bash"

if [[ ! -f "$QWIKI_DEVELOPMENT_ROOT_CA/rootCA.pem" ]]; then
    log_info "Generating root CA in '$QWIKI_DEVELOPMENT_ROOT_CA'"
    CAROOT="$QWIKI_DEVELOPMENT_ROOT_CA" mkcert -install
fi


host="localhost"
location="$HOME/certs/$host"

if [[ ! -f "$location/$host.pem" ]]; then
    log_info "Generating certificate for '$host' in '$location'"
    mkdir -p "$location"

    pushd "$location" > /dev/null
    CAROOT="$QWIKI_DEVELOPMENT_ROOT_CA" mkcert "$host"
    popd > /dev/null
fi
