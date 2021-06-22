# Usage

## Preparations
```BASH
sudo apt install git
cd ~
git clone https://github.com/modell-aachen/modac-dev-machine-setup.git
cp $HOME/modac-dev-machine-setup/provisioning/inventory_custom_example.yml $HOME/.inventory_local.yml
```
## Local adjustments
```BASH
vim $HOME/.inventory_local.yml
```
- remove unused packages and snaps, insert your keys, e.g.
  * FONTAWESOME_NPM_AUTH_TOKEN: https://password.int.modac.eu/WebClient/Main?itemId=717a1db1-9bd2-457f-9dc3-6b4f25670524
- create a local ssh key

## Initial setup
```BASH
cd $HOME/modac-dev-machine-setup/
./dev-provision -h
machine provision -i ~/.inventory_local.yml packages
machine provision -i ~/.inventory_local.yml tooling
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
