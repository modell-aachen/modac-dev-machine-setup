#!/usr/bin/env bash
set -e

sudo apt update
sudo apt install -y \
    curl \
    easy-rsa \
    htop \
    inotify-tools \
    net-tools \
    network-manager-openvpn \
    network-manager-openvpn-gnome \
    openvpn \
    restic \
    python-is-python3 \
    apt-transport-https \
    ca-certificates \
    gnupg \
    age
