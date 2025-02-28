#!/usr/bin/env bash
set -e

sudo apt update
sudo apt install -y \
    python3-pip \
    easy-rsa \
    htop \
    inotify-tools \
    net-tools \
    network-manager-openvpn \
    network-manager-openvpn-gnome \
    openvpn \
    python-is-python3 \
    apt-transport-https \
    ca-certificates \
    gnupg \
    age
