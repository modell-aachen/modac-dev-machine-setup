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

## Inventory creation

### Clone repo
```
git clone https://github.com/modell-aachen/modac-dev-machine-setup.git
```

### Create inventroy blueprint
**For development users**
```
cp $HOME/modac-dev-machine-setup/provisioning/inventory_custom_example.yml $HOME/.inventory_local.yml
```
**For service users**
```
cp $HOME/modac-dev-machine-setup/provisioning/inventory_service_example.yml $HOME/.inventory_local.yml
```

## Initial local configuration and adjustments
```BASH
vim $HOME/.inventory_local.yml
```
1) remove unused packages and snaps, insert your keys, e.g.
2) create a local ssh key and set as SSH key to your GitHub Account (https://github.com/settings/keys)
3) check successfull authentication against github.com
    ```bash
    ssh -T git@github.com
    ```
3) create a directory for your repositories and set REPOS_DIRECTORY to the created path
7) set NEXUS_BOT_TOKEN to the value of 1Password entry at https://start.1password.com/open/i?a=CXJNQFCHNNGSLNOEP6SLPHLZQ4&v=3mhbwhicfwkifyq7bc2nrnhywa&i=32jtgjyel43vabz2wbukslo5bq&h=modac.1password.eu

### Additional, development only, configuration and adjustments
5) set GITHUB_AUTH_TOKEN to the value of 1Password entry at https://start.1password.com/open/i?a=CXJNQFCHNNGSLNOEP6SLPHLZQ4&h=modac.1password.eu&i=dwpktyrfuj6cyjfy6y74q3ifiy&v=6u4nznoclnkg7467ne4ntutcgq
4) set FONTAWESOME_NPM_AUTH_TOKEN to the value of 1Password entry at https://start.1password.com/open/i?a=CXJNQFCHNNGSLNOEP6SLPHLZQ4&h=modac.1password.eu&i=xmhedekcuokrqrch62bsuvr5lu&v=6u4nznoclnkg7467ne4ntutcgq
1) set HARBOR_USERNAME to your modac mail address
2) set HARBOR_PASSWORD to your CLI token (https://harbor.modac.cloud -> Login -> user profile [top right corner] -> User Profile -> CLI secret)


## Initial setup
```BASH
cd $HOME/modac-dev-machine-setup/
./dev-provision -i ~/.inventory_local.yml
```

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
