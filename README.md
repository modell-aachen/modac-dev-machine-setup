# Required packages

The machine provisioner requires some basic packages.

For MacOs: Start  of with `xcode-select --install` to launch the installer for some base packages.

All other operating systems don't have any further requirements

# Usage

## Installation

### Install provisioner and devbox packages

Linux
```BASH
wget -qO- https://raw.githubusercontent.com/modell-aachen/modac-dev-machine-setup/refs/heads/main/install | bash; source ~/.bashrc
```

MacOs
```BASH
curl -fsSL https://raw.githubusercontent.com/modell-aachen/modac-dev-machine-setup/refs/heads/main/install | "$SHELL"
source ~/."$(basename $SHELL)"rc
```

### Login to 1Password

* The 1Password app is provisioned. You should be able to login

### Provide needed secrets to 1Password

**Harbor**
1) in 1Password: New Item
1) Add Login
1) Change `Login` title to `Harbor`
1) set username to your modac email address
1) set password (https://harbor.modac.cloud -> Login -> user profile [top right corner] -> User Profile -> CLI secret)
1) enable cli integration (1Password app > ... > Settings > Developer > Command-Line Interface)

**GitHub token (Claude Code plugins)**
1) create a classic PAT at https://github.com/settings/tokens (*Generate new token (classic)*) with scopes `repo` and `read:org`
1) in 1Password (`Employee` vault): New Item → Password, title `Github Token`, paste the PAT in the `password` field

check that you can login to 1password:
```BASH
op vault list
```

### [OPTIONAL] restore devbox.json

If you want to restore other files than `devbox.json` you have to call the `backup restore` twice.
The first time, the `devbox.json` is restored.
The second command will restore the files defined in the restored `devbox.json`.

```BASH
eval "$(op signin)"
machine backup restore
machine backup restore
```

### Provision your system

```BASH
machine provision
```


Back to QwikiContrib: [QwikiContrib](https://github.com/modell-aachen/QwikiContrib/)

## Service machine (Windows/WSL)

Service machines get a minimal setup for Kubernetes access instead of the full
developer environment: a reduced devbox package set (kubectl, helm, krew, k9s,
gcloud with GKE auth plugin) and only the service-flagged provisioning modules
(`machine provision list-modules` marks them with `*`). 1Password is installed
CLI-only, without the desktop app.

### Windows preparation

1) Open an admin PowerShell and run `wsl --install -d Ubuntu`
1) Reboot and create your UNIX user when Ubuntu starts
1) Make sure systemd is enabled in `/etc/wsl.conf` (default on current Ubuntu images):
   ```
   [boot]
   systemd=true
   ```
   After changing it, run `wsl --shutdown` from Windows and start WSL again.

### Installation

Inside WSL (or on a plain Ubuntu machine — the service profile is not tied to WSL):
```BASH
wget -qO- https://raw.githubusercontent.com/modell-aachen/modac-dev-machine-setup/refs/heads/main/install | MACHINE_PROFILE=service bash; source ~/.bashrc
```

The profile is persisted in `~/.machine/profile`, so later runs of
`machine provision` select the service module set automatically. To switch an
existing machine, run `machine provision --profile service` (or `--profile dev`).

### Provision

```BASH
machine provision
```

During provisioning you will be asked to add your 1Password account
(sign-in address, email, secret key, master password) and to complete the
device-code login for Kubernetes access.

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

## Erlang build fails with "No curses library functions found" or "CFLAGS must contain a -O flag" (macOS)

When building Erlang via `asdf install erlang <version>` inside the nix-shell/devbox environment on macOS, the nix-shell can hide Homebrew's ncurses from the build system. Set these env vars before installing:

```bash
brew install ncurses
export KERL_CONFIGURE_OPTIONS="--with-ssl=$(brew --prefix openssl@3) --without-javac"
export LDFLAGS="-L$(brew --prefix ncurses)/lib -L/opt/homebrew/opt/unixodbc/lib"
export CFLAGS="-O2 -g -I$(brew --prefix ncurses)/include"
asdf install erlang 28.2
```

**Important:** `CFLAGS` must include `-O2 -g` — Erlang's configure requires an optimization flag.

## CLAUDE.md: How to add personal preferences

The team's CLAUDE.md contains what we all can agree on. But everyone works differently. So feel free to enter your personal preferences to ~/.claude/personal-CLAUDE.md.
It will be included automatically and take precedence over the team defaults if conflicting.
