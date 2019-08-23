# echo "Link VSCode"
# DOTFILES_DIR=$(pwd)
# VSCODE_USER=~/.config/Code/User
# rm -f $VSCODE_USER/keybindings.json
# ln -s $DOTFILES_DIR/vscode/user/keybindings.json $VSCODE_USER/keybindings.json
# rm -f $VSCODE_USER/projects.json
# ln -s $DOTFILES_DIR/vscode/user/projects.json $VSCODE_USER/projects.json
# rm -f $VSCODE_USER/settings.json
# ln -s $DOTFILES_DIR/vscode/user/settings.json $VSCODE_USER/settings.json
# rm -rf $VSCODE_USER/snippets.json
# ln -s $DOTFILES_DIR/vscode/user/snippets $VSCODE_USER/snippets.json

# echo "Install/Updating Fonts.."
# ./setup_fonts.sh

# echo "Install and Configure GIT"
# sh ./git/setup.sh
# sh ./git/configure.sh

# echo "Done."