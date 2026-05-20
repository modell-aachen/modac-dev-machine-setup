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

# Kubie spawnt seine Subshell ueber bare 'bash' PATH-Lookup. Wenn irgendwo
# ein bash-minimal (kein progcomp, keine readline) vor dem devbox-Profil im
# PATH liegt, ist die Subshell unbrauchbar (shopt/complete-Fehler, keine
# Pfeiltasten). Wir ziehen das devbox-Profil-bin gezielt fuer den kubie-Call
# nach vorne, damit unsere bashInteractive zuerst gefunden wird.
__modac_kubie() {
    local devbox_bin="$HOME/.local/share/devbox/global/default/.devbox/nix/profile/default/bin"
    PATH="$devbox_bin:$PATH" command kubie "$@"
}
kctx() { __modac_kubie ctx "$@"; }
kns()  { __modac_kubie ns  "$@"; }

# Completion fuer kctx/kns. Quelle ist 'kubectl config get-contexts -o name'
# bzw. 'kubectl get namespaces -o name' — bewusst kubectl-basiert statt
# kubies _kubie aufzurufen, damit es ohne compinit-Reihenfolge-Magie und
# in bash wie zsh identisch funktioniert.
if [ -n "$ZSH_VERSION" ]; then
    _modac_kctx_zsh() {
        local -a items
        items=(${(f)"$(kubectl config get-contexts -o name 2>/dev/null)"})
        _describe -t contexts 'kube context' items
    }
    _modac_kns_zsh() {
        local -a items
        items=(${(f)"$(kubectl get namespaces -o name 2>/dev/null | sed 's|^namespace/||')"})
        _describe -t namespaces 'kube namespace' items
    }
    (( $+functions[compdef] )) || { autoload -Uz compinit && compinit -u 2>/dev/null; }
    compdef _modac_kctx_zsh kctx
    compdef _modac_kns_zsh kns
elif [ -n "$BASH_VERSION" ]; then
    _modac_kctx_bash() {
        local cur="${COMP_WORDS[COMP_CWORD]}"
        COMPREPLY=( $(compgen -W "$(kubectl config get-contexts -o name 2>/dev/null)" -- "$cur") )
    }
    _modac_kns_bash() {
        local cur="${COMP_WORDS[COMP_CWORD]}"
        COMPREPLY=( $(compgen -W "$(kubectl get namespaces -o name 2>/dev/null | sed 's|^namespace/||')" -- "$cur") )
    }
    complete -F _modac_kctx_bash kctx
    complete -F _modac_kns_bash kns
fi

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
