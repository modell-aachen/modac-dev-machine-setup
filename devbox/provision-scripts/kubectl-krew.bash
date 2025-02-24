#!/bin/bash
set -e

for plugin in "ctx" "ns" "konfig" "oidc-login"; do
    kubectl krew install $plugin
done
