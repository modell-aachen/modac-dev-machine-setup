#!/usr/bin/env bash
set -e

source "$(dirname "$0")/../helper.bash"

if [ -n "$CONTAINER_ID" ]; then
    log_info "Running inside a distrobox, skipping k3s network setup"
    exit 0
fi

ip="172.25.5.1"
name="k3s-vr"
ip_address_line=$( nmcli connection show k3s-vr /dev/null 2>&1 | grep ipv4.addresses )


if [[ "$ip_address_line" == *ipv4.addresses* && "$ip_address_line" != *$ip* ]]; then
    log_info "Deleting existing k3s-vr connection with incorrect IP"
    nmcli connection delete $name
fi

if [[ "$ip_address_line" != ipv4.addresses*$ip* ]]; then
    log_info "Creating k3s-vr network connection with IP $ip"
    nmcli connection add \
        type dummy \
        ifname $name \
        ipv4.addresses $ip \
        ipv4.method manual \
        con-name $name
fi
