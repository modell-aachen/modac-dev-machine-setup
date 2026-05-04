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

alias kctx='kubie ctx'
alias kns='kubie ns'

# kubie-only: current-context in ~/.kube/config entfernen, damit kubectl
# ausserhalb einer kubie-Subshell keinen impliziten Context hat.
# Nur wenn KUBECONFIG leer ist (sonst wuerden wir die temp-Config von kubie zerschiessen).
if [ -z "$KUBECONFIG" ] && [ -f "$HOME/.kube/config" ] && command -v kubectl >/dev/null; then
    if kubectl --kubeconfig="$HOME/.kube/config" config current-context >/dev/null 2>&1; then
        kubectl --kubeconfig="$HOME/.kube/config" config unset current-context >/dev/null
    fi
fi
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
