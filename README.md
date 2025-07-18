# Required packages

**MacOS**: Please start with [MacOS required packages](README_mascos.md).

The machine provisioner requires some basic packages.
If your using ubuntu >= 20.04, you can follow along with the "Usage" section. On other operating system, i.e. MacOS, please satisfy the following requirements manually:

* python >= 3.10
* pip >= 20.3
* pip3 >= 20.3
* git
* vim

**Important** The command `python3` needs to point to the python-executable for Python >= 3.10, `python` might not work.

## MacOS

These packages can be installed using xcode-select, eg. run `xcode-select --install` to launch the installer.

# Usage

## Preparations
```BASH
sudo apt update
sudo apt install -y software-properties-common git vim
cd ~
```

### Configure 1Password
As preparation you need to
* log into 1Password app
* integrate with 1Password cli.
  Follow the instructions here: https://developer.1password.com/docs/cli/get-started/#step-2-turn-on-the-1password-desktop-app-integration

After following the instructions you should be able to run
```BASH
op vault list
```
and see a list of your vaults.

1) check successfull authentication against github.com
    ```bash
    ssh -T git@github.com
    ```

if not, one possibility is :
```bash
sudo apt install gh
gh auth login
```
do auth with ssh when prompted

### Add Harbor secrets to 1Password
1) in 1Password: New Item
1) Add Login
1) Change `Login` title to `Harbor`
1) set username to your modac email address
1) set password (https://harbor.modac.cloud -> Login -> user profile [top right corner] -> User Profile -> CLI secret)


## Installation

### Clone repo
```
cd ~
git clone https://github.com/modell-aachen/modac-dev-machine-setup.git
cd modac-dev-machine-setup
```

### Install provisioner and devbox packages
```BASH
./devbox/provision
```

### By `Please logout and login again to use docker without sudo` restart your laptop and then:
check again that you can login to 1password:
```BASH
op vault list
```
Then resume:
```BASH
 cd modac-dev-machine-setup/
./devbox/provision
```

## Provision dev machine

```BASH
source $HOME/.bashrc
machine provision
```

After that open up a new terminal to have an updated PATH with all the tools available.

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
