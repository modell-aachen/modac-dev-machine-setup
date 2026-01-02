_qypress_completion() {
    
    local cur="${COMP_WORDS[COMP_CWORD]}"

    case ${COMP_CWORD} in
        1)
            COMPREPLY=(
                $( compgen -W "--help" -- $cur )
                $( compgen -W '"https://dev.qluster.localhost"' -- $cur)
            )
            ;;
        *)
            COMPREPLY=($( compgen -W "--help --multisite --remote help version run open install verify cache info" -- $cur ))
            ;;
    esac
}

complete -F _qypress_completion qypress
