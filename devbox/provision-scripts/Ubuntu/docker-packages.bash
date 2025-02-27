#!/usr/bin/env bash
set -e

if [ ! -f /etc/apt/keyrings/docker.asc ]; then
    sudo apt-get install ca-certificates curl
    sudo install -m 0755 -d /etc/apt/keyrings
    sudo curl -fsSL https://download.docker.com/linux/ubuntu/gpg -o /etc/apt/keyrings/docker.asc
    sudo chmod a+r /etc/apt/keyrings/docker.asc
fi

if [ ! -f /etc/apt/sources.list.d/docker.list ]; then
    echo "deb [signed-by=/etc/apt/keyrings/docker.asc] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable" | sudo tee /etc/apt/sources.list.d/docker.list
    sudo apt update
fi

sudo apt-get install -y \
    docker-ce \
    docker-ce-cli \
    containerd.io \
    docker-buildx-plugin \
    docker-compose-plugin

BIRed='\033[1;91m'
NC='\033[0m'

if [[ $( id -nG "$USER" | grep -w docker ) != *docker* ]]; then
    sudo usermod -aG docker "$USER"
    echo -e "Please ${BIRed}logout and login again$NC to use docker without sudo"
    exit 1
fi
