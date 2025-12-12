#!/usr/bin/env bash
set -e

brew install \
    bash \
    gettext \
    gnu-getopt \
    gpg \
    openssh \
    libfido2 \
    openvpn \
    nmap

brew install --cask \
    visual-studio-code \
    yubico-authenticator \
    openvpn-connect \
    orbstack
