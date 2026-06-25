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

### Provide needed secrets to 1Password for Harbor
1) in 1Password: New Item
1) Add Login
1) Change `Login` title to `Harbor`
1) set username to your modac email address
1) set password (https://harbor.modac.cloud -> Login -> user profile [top right corner] -> User Profile -> CLI secret)
1) enable cli integration (1Password app > ... > Settings > Developer > Command-Line Interface)

check that you can login to 1password:
```BASH
op vault list
```

### Provide a GitHub token to 1Password for Claude Code plugins

Claude Code refreshes private plugin marketplaces (e.g. `modell-aachen/claude-skills`)
at startup via a background `git pull`. That pull cannot use your interactive
`gh` login, so it needs `GITHUB_TOKEN` in the environment — otherwise the marketplace
clone goes stale and new plugins/skills silently never appear.

1) create a **classic** Personal Access Token (https://github.com/settings/tokens)
   on your own GitHub account with scopes `repo` and `read:org`
1) in 1Password (vault `Entwicklung`): New Item → Password
1) title it `GitHub Plugin Marketplace Token`
1) put the PAT in the `credential` field

> Note: because `GITHUB_TOKEN` is exported into your shell, `gh` CLI commands use it
> in place of your keyring login — so it must be **your own** PAT, not a shared/bot
> token. Git push/clone are unaffected (they use SSH).

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
