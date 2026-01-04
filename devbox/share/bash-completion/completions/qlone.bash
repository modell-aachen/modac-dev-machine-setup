_qlone-completion()
{
    COMPREPLY=()
    local cur="${COMP_WORDS[COMP_CWORD]}"

    local repos="
        QwikiContrib
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
