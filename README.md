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
./dev-provision -h
./dev-provision -i ~/.inventory_local.yml packages
./dev-provision -i ~/.inventory_local.yml tooling
```

## Updates
Update your `$HOME/.inventory_local.yml`:
```BASH
machine edit-config
```

Apply updates:
```BASH
machine provision -i ~/.inventory_local.yml packages
machine provision -i ~/.inventory_local.yml tooling
```

## Local Q.wikis
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
