#!/usr/bin/env bash
set -e

source "$(dirname "$0")/helper.bash"

if [[ "$( op signin 2>&1 )" == *"[ERROR]"* ]]; then
    log_error "You are not logged in to 1Password CLI. Please log into 1Password CLI and try again."
    exit 1
else
    log_info "Setting up Kubernetes cluster configuration"
    curl -s https://modell-aachen.github.io/k8s-kubeconfig-setup/kubeconfig-setup.sh | bash -s -- --merge
fi
