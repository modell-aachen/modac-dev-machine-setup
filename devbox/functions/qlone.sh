qlone() {


    usage() {
        printf -v text "%s" \
            "qlone <modell-aachen repository> [OPTION...]\n" \
            "    -f  --force        clones repository even when old exists\n" \
            "        --no-cd        don't change into cloned directory\n" \
            "    -h, --help         shows this help message\n"
        printf "$text"
    }


    local is_cloning_forced=0
    local is_switching_directory=1

    OPTS=`getopt -o fh --long no-cd,force,help -- "$@"`
    if [ $? != 0 ] ; then echo "Failed parsing options." >&2 ; return 1 ; fi

    eval set -- "$OPTS"

    while true; do
        case "$1" in
            -f | --force )
                is_cloning_forced=1
                shift ;;
            --no-cd )
                is_switching_directory=0
                shift ;;
            -h | --help )
                usage
                return
                shift ;;
            -- )
                shift
                break ;;
            * )
                break ;;
        esac
    done

    shift $(expr $OPTIND - 1 )
    local repo=${1:-$(basename $(pwd))}

    if [ -z "$REPOS_DIRECTORY" ] || [ ! -d "$REPOS_DIRECTORY" ] ; then
        echo "repos directory '$REPOS_DIRECTORY' does not exist!"
        return
    fi

    local repo_path="$REPOS_DIRECTORY/$repo"
    if [ -z "$repo" ] ; then
        repo_path=`pwd`
    fi

    local branch=""

    if [ -d "$repo_path" ] && [ "$is_cloning_forced" = 1 ] ; then
        cd "$repo_path"
        branch=`git rev-parse --abbrev-ref HEAD`
        cd $HOME
        rm -rf "$repo_path"
    fi

    if [ ! -d "$repo_path" ] ; then
        git clone git@github.com:modell-aachen/$repo.git "$repo_path"
    fi

    if [ "$is_switching_directory" = 1 ] ; then
        cd "$repo_path"
    fi

    if [ "$branch" ] ; then
        git checkout "$branch"
    fi
}


_qlone-completion()
{
    COMPREPLY=()
    local cur="${COMP_WORDS[COMP_CWORD]}"

    local repos="
        QwikiContrib
        RMS
        deploy
        devops-machine-provisioner
        dotfiles
        dotfiles-pandoc
        latex-modac
        qwiki-cli
        qwiki-gitops
        qwikinow-deployment
        terraform
    "

    COMPREPLY=( $(compgen -W "$repos" -- ${cur}) )
}

complete -F _qlone-completion qlone
