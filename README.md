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

# Usage

## Preparations
```BASH
sudo apt update
sudo apt install -y software-properties-common git vim python3 python3-pip python-is-python3
cd ~
```

## Install and configure 1Password
As preparation you need to install 1password CLI and be able to login.
Follow the instructions here: https://developer.1password.com/docs/cli/get-started

After following the instructions you should be able to run
```BASH
op vault list
```
and see a list of your vaults.

## Inventory creation

### Clone repo
```
git clone https://github.com/modell-aachen/modac-dev-machine-setup.git
cd modac-dev-machine-setup
```



### Install provisioner and devbox packages
```BASH
./devbox/install
```

## Initial local envs
```BASH
vim $HOME/.env
```
1) check successfull authentication against github.com
    ```bash
    ssh -T git@github.com
    ```
1) set NEXUS_BOT_TOKEN to the value of 1Password entry at https://start.1password.com/open/i?a=CXJNQFCHNNGSLNOEP6SLPHLZQ4&v=3mhbwhicfwkifyq7bc2nrnhywa&i=32jtgjyel43vabz2wbukslo5bq&h=modac.1password.eu

### Additional, development only, configuration and adjustments
6) set GITHUB_AUTH_TOKEN to the value of 1Password entry at https://start.1password.com/open/i?a=CXJNQFCHNNGSLNOEP6SLPHLZQ4&h=modac.1password.eu&i=dwpktyrfuj6cyjfy6y74q3ifiy&v=6u4nznoclnkg7467ne4ntutcgq
1) set FONTAWESOME_NEXUS_AUTH_TOKEN to the value of 1Password entry at https://start.1password.com/open/i?a=CXJNQFCHNNGSLNOEP6SLPHLZQ4&v=6u4nznoclnkg7467ne4ntutcgq&i=5xls52q24au3eie2itucvc635u&h=modac.1password.eu
1) set HARBOR_USERNAME to your modac mail address
1) set HARBOR_PASSWORD to your CLI token (https://harbor.modac.cloud -> Login -> user profile [top right corner] -> User Profile -> CLI secret)


## Provision dev machine

```BASH
./devbox/provision
```

After that open up a new terminal to have an updated PATH with all the tools available.

## Updates
Update your `$HOME/.inventory_local.yml`:
```BASH
machine edit-config
```

Apply updates:
```BASH
machine provision
```

# FAQ
## Problem lösen "**GitHub Error Message - Permission denied (publickey)**"
Source: https://stackoverflow.com/questions/12940626/github-error-message-permission-denied-publickey
Solution: Write that before starting the script
```BASH
ssh-agent -s
ssh-add ~/.ssh/id_rsa
```

## DNS resolver
Local Q.wiki (e.g. dev.modac) can't be resolved  ( < Ubuntu 21.04 only):

```BASH
sudo systemctl restart lxd-host-dns.service
```
## ZSH
If you use zsh, source `.env` and `bashrc.sh` in your `.zshrc`
```BASH
[ -f $HOME/.env ] && source $HOME/.env
[ -f $HOME/.modac-bash/bashrc.sh ] && source $HOME/.modac-bash/bashrc.sh
```

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
