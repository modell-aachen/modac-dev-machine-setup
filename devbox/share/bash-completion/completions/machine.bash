_machine-completion()
{
    local cur="${COMP_WORDS[COMP_CWORD]}"
    local prev=${COMP_WORDS[COMP_CWORD-1]}
    case ${COMP_CWORD} in
        1)
            COMPREPLY=($(compgen -W "--help install provision edit-config backup" -- $cur))
            ;;
        *)
            case $prev in
                provision )
                    local local_tooling="--help --filter --skip-install list-modules"
                    COMPREPLY=($(compgen -W "$local_tooling" -- $cur))
                    ;;
                backup )
                    local sub_commands="--help create restore"
                    COMPREPLY=($(compgen -W "$sub_commands" -- $cur))
                    ;;
            esac
            ;;
    esac
}

complete -F _machine-completion machine

