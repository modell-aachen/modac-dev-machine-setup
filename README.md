# Required packages
The machine provisioner requires some basic packages.
If your using ubuntu >= 20.04, you can follow along with the "Usage" section. On other operating system, i.e. MacOS, please satisfy the following requirements manually:

* python >= 3.8
* pip >= 20.3
* pip3 >= 20.3
* git
* vim

**Important** The command `python` needs to point to the python-executable for Python >= 3.8, `python3` will not work.


# Usage

## Preparations
```BASH
sudo apt update
sudo apt install -y software-properties-common git vim python3 python3-pip python-is-python3
cd ~
git clone https://github.com/modell-aachen/modac-dev-machine-setup.git
cp $HOME/modac-dev-machine-setup/provisioning/inventory_custom_example.yml $HOME/.inventory_local.yml
```
## Initial local configuration and adjustments
```BASH
vim $HOME/.inventory_local.yml
```
1) remove unused packages and snaps, insert your keys, e.g.
2) create a local ssh key and set as SSH key to your GitHub Account (https://github.com/settings/keys)
3) create a directory for your repositories and set REPOS_DIRECTORY to the created path
4) set FONTAWESOME_NPM_AUTH_TOKEN to https://start.1password.com/open/i?a=CXJNQFCHNNGSLNOEP6SLPHLZQ4&h=modac.1password.eu&i=xmhedekcuokrqrch62bsuvr5lu&v=6u4nznoclnkg7467ne4ntutcgq
5) set GITHUB_AUTH_TOKEN to https://start.1password.com/open/i?a=CXJNQFCHNNGSLNOEP6SLPHLZQ4&h=modac.1password.eu&i=dwpktyrfuj6cyjfy6y74q3ifiy&v=6u4nznoclnkg7467ne4ntutcgq
6) set RMS_AUTH_TOKEN to https://start.1password.com/open/i?a=CXJNQFCHNNGSLNOEP6SLPHLZQ4&h=modac.1password.eu&i=loyd7k5wwwnxkp5ncbkkh7pnmq&v=6u4nznoclnkg7467ne4ntutcgq

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
machine provision -i ~/.inventory_local.yml
```

## Local LXC Q.wiki deployments
E.g. dev, master, etc.
### Setup
```BASH
qontainer init
```

### Build
E.g. dev (accessible in browser as `dev.qwiki`)
```BASH
qontainer create dev
```

### Connect
E.g. dev
```BASH
qontainer login dev
```

## Local kubernetes (container) Q.wiki deployment
E.g. dev, master, etc.
### Setup
```BASH
qluster init

qaffold init
```

### Build
Deploy current checkedout Q.wiki
```BASH
qaffold deploy
```

### Connect
```BASH
qaffold login
```

# FAQ
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
