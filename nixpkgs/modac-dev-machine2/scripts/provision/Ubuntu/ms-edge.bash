#!/usr/bin/env bash
set -e

gpg_key_dir="/etc/apt/trusted.gpg.d/"
gpg_key_filename="microsoft-edge.gpg"

if [[ ! -f "$gpg_key_dir$gpg_key_filename" ]]; then
    curl https://packages.microsoft.com/keys/microsoft.asc | gpg --dearmor > "/tmp/$gpg_key_filename"
    sudo install -o root -g root -m 644 "/tmp/$gpg_key_filename" "$gpg_key_dir"
    sudo rm "/tmp/$gpg_key_filename"
fi

ms_edge_package_source="/etc/apt/sources.list.d/microsoft-edge-dev.list"
if [[ ! -f "$ms_edge_package_source" ]]; then
    sudo sh -c "echo \"deb [arch=amd64] https://packages.microsoft.com/repos/edge stable main\" > $ms_edge_package_source"

    sudo apt update
    sudo apt install -y microsoft-edge-stable
fi
