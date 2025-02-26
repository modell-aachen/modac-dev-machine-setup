#!/usr/bin/env bash
set -e

if [[ "$( op whoami 2>&1 )" == *"[ERROR]"* ]]; then
    echo "Your are not logged in to 1Password CLI. Please login and try again."
    exit 1
else
    curl -s https://modell-aachen.github.io/k8s-kubeconfig-setup/kubeconfig-setup.sh | bash -s -- --merge
fi
