_qe2e-completion() {
    local cur="${COMP_WORDS[COMP_CWORD]}"

    case ${COMP_CWORD} in
        1)
            COMPREPLY=($(compgen -W "run clean --help" -- $cur))
            ;;
        *)
            local subcommand="${COMP_WORDS[COMP_CWORD-1]}"
            case $subcommand in
                -b | --branch )
                    cd "$REPOS_DIRECTORY/QwikiContrib" >/dev/null
                    local branches=$(git branch --format="%(refname)" -a | sed -e "s#^refs/\(heads\|remotes/origin\)/##g")
                    cd - >/dev/null
                    COMPREPLY=($(compgen -W "$branches" -- $cur))
                    ;;
                -w | --which )
                    COMPREPLY=($(compgen -W "regular multisite all" -- $cur))
                    ;;
                *)
                    COMPREPLY=($(compgen -W "--branch --which --legacy" -- $cur))
                    ;;
            esac
            ;;
    esac
}

complete -F _qe2e-completion qe2e
