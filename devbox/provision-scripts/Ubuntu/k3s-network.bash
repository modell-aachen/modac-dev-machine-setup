#!/usr/bin/env bash

ip="172.25.5.1"
name="k3s-vr"
ip_address_line=$( nmcli connection show k3s-vr /dev/null 2>&1 | grep ipv4.addresses )

set -e

if [[ "$ip_address_line" == *ipv4.addresses* && "$ip_address_line" != *$ip* ]]; then
    nmcli connection delete $name
fi

if [[ "$ip_address_line" != ipv4.addresses*$ip* ]]; then
    nmcli connection add \
        type dummy \
        ifname $name \
        ipv4.addresses $ip \
        con-name $name
fi
