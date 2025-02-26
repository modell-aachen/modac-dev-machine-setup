#!/bin/bash
set -e

for plugin in "ctx" "ns" "konfig" "oidc-login"; do
    krew install $plugin
done
