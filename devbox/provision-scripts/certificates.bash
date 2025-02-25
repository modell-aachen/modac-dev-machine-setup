#!/usr/bin/env bash
set -e

if [[ ! -f "$QWIKI_DEVELOPMENT_ROOT_CA/rootCA.pem" ]]; then
    echo "Generating root CA in '$QWIKI_DEVELOPMENT_ROOT_CA'"
    CAROOT="$QWIKI_DEVELOPMENT_ROOT_CA" mkcert -install
fi


host="localhost"
location="$HOME/certs/$host"

if [[ ! -f "$location/$host.pem" ]]; then
    echo "Generating certificate for '$host' in '$location'"
    mkdir -p "$location"

    pushd "$location"
    CAROOT="$QWIKI_DEVELOPMENT_ROOT_CA" mkcert "$host"
    popd
fi
