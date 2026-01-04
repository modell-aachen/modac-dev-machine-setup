#!/usr/bin/env bash
set -e

sudo apt update
sudo apt install -y \
    python3-pip \
    easy-rsa \
    inotify-tools \
    net-tools \
    network-manager-openvpn \
    network-manager-openvpn-gnome \
    openvpn \
    python-is-python3 \
    apt-transport-https \
    ca-certificates \
    gnupg \
    libnss3-tools

# Install python build dependencies
sudo apt install -y \
    make \
    build-essential \
    libssl-dev \
    zlib1g-dev \
    libbz2-dev \
    libreadline-dev \
    libsqlite3-dev \
    curl \
    git \
    libncursesw5-dev  \
    xz-utils \
    tk-dev \
    libxml2-dev \
    libxmlsec1-dev \
    libffi-dev \
    liblzma-dev \
    libzstd-dev
