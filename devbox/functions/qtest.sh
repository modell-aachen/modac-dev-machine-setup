qtest() {
    local workingDir=$(pwd)
    qlone QwikiContrib
    ./core/backend-test "$@"
    cd "$workingDir"
}

_qtest-completion()
{
    local cur="${COMP_WORDS[COMP_CWORD]}"
    local prev="${COMP_WORDS[COMP_CWORD-1]}"

    case "$prev" in
        --filter )
            pushd "$REPOS_DIRECTORY/QwikiContrib" > /dev/null
                COMPREPLY=( $( compgen -d -S '/' -- ${cur} ) )
                compopt -o nospace
            popd > /dev/null
            ;;
        * )
            local options="-c --no-clear --image -i --integration -p --no-pull -v --verbose -w --watch --help -h --filter --version"
            COMPREPLY=( $(compgen -W "$options" -- ${cur}) )
            ;;
    esac
}

complete -F _qtest-completion qtest
