#!/bin/sh

echo "Adding default route to $route_vpn_gateway with /0 mask..."
ip route add default via $route_vpn_gateway

echo "Removing /1 routes..."
ip route del 0.0.0.0/1 via $route_vpn_gateway
ip route del 128.0.0.0/1 via $route_vpn_gateway