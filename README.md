# Required packages

The machine provisioner requires some basic packages.

For MacOs: Start  of with `xcode-select --install` to launch the installer for some base packages.

All other operating systems don't have any further requirements

# Usage

## Installation

### Install provisioner and devbox packages
```BASH
wget -qO- https://raw.githubusercontent.com/modell-aachen/modac-dev-machine-setup/refs/heads/main/nixpkgs/modac-dev-provisioner/share/modac-dev-provisioner/bin/install | bash; source ~/.bashrc
```

### Login to 1Password

* The 1Password app is provisioned. You should be able to login

### Provide needed secrets to 1Password for Harbor
1) in 1Password: New Item
1) Add Login
1) Change `Login` title to `Harbor`
1) set username to your modac email address
1) set password (https://harbor.modac.cloud -> Login -> user profile [top right corner] -> User Profile -> CLI secret)
1) enable cli integration

check that you can login to 1password:
```BASH
op vault list
```

### Provision your system

```BASH
machine provision
```


Back to QwikiContrib: [QwikiContrib](https://github.com/modell-aachen/QwikiContrib/)

## Updates
Update your `$(devbox global path)/devbox.json`:
```BASH
machine edit-config
```

Apply updates:
```BASH
machine provision
```

# FAQ

## I want to use a different 'REPOS_DIRECTORY' than '$HOME/qwiki-repos'

Add an ENV variable to "$(devbox global path)/devbox.json:
```BASH
machine edit-config
```
* add `"REPOS_DIRECTORY": "$HOME/path"` to the `env` object

## Docker setup fails

The docker task requires a system restart. If any succeeding docker-related tasks fail (kubectl, calico, etc.) try restarting the system
and the provisioning process (`machine provision`).

## Calico: ImageInspectError

The calico task can fail because of an [incompatibility issue](https://github.com/k3s-io/k3s/issues/9279):

Workaround:

```BASH
apt-cache policy docker-ce | grep 24.0.7
```

Use the printed version (e.g. 5:24.0.7-1~ubuntu.22.04~jammy) to ...

```BASH
sudo apt install docker-ce=*VERSION*
```

## CLAUDE.md: How to add personal preferences

The team's CLAUDE.md contains what we all can agree on. But everyone works differently. So feel free to enter your personal preferences to ~/.claude/personal-CLAUDE.md.
It will be included automatically and take precedence over the team defaults if conflicting.
