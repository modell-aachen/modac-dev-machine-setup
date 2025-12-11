#!/usr/bin/env bash
set -e

brew install \
    bash \
    gettext \
    gnu-getopt \
    gpg \
    openssh \
    libfido2

brew install --cask \
    visual-studio-code \
    yubico-authenticator
