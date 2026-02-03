package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

const aliasesContent = ` qlone() {
    if [ "$#" -ne 1 ]; then
      echo "Usage: qlone <modell-aachen repository>" >&2
      exit 1
    fi

    local repo="$1"

    if [ -z "$REPOS_DIRECTORY" ] || [ ! -d "$REPOS_DIRECTORY" ] ; then
        echo "repos directory '$REPOS_DIRECTORY' does not exist!"
        return
    fi

    local repo_path="$REPOS_DIRECTORY/$repo"

    if [ ! -d "$repo_path" ] ; then
        gh repo clone "modell-aachen/$repo" "$repo_path"
    fi

    cd "$repo_path"
}

_qlone-completion()
{
    COMPREPLY=()
    local cur="${COMP_WORDS[COMP_CWORD]}"

    local repos="
        QwikiContrib
        modac-dev-machine-setup
        modac-shell-helper
        dotfiles
        dotfiles-pandoc
        latex-modac
        qwiki-cli
        qwiki-gitops
        qwikinow-deployment
    "

    COMPREPLY=( $(compgen -W "$repos" -- ${cur}) )
}

complete -F _qlone-completion qlone

alias k='kubectl'
complete -o default -F __start_kubectl k
`

var aliasesCmd = &cobra.Command{
	Use:   "aliases",
	Short: "Output shell aliases and completions",
	Long:  "Sets aliases and functions in shell like qlone, k",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Print(aliasesContent)
		return nil
	},
}
