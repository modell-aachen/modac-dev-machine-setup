machine() {

    usage() {
        cat << USAGE
Usage:
  machine [command]

Available Commands:
  edit-config    edit custom inventory file
  provision      provision this machine
  pull-changes   pull changes for dev machine provisioner
  backup-devbox  backup devbox configuration devbox.json

Flags:
  -h, --help     shows this help message
USAGE
    }

    while [[ "$1" == * ]] ; do
        case "$1" in
            -h | --help )
                usage
                return
                ;;
            -- )
                shift
                break
                ;;
            -* )
                echo "failed parsing option '$1'" >&2
                return
                ;;
            * )
                break
                ;;
        esac
        shift
    done

    local subcommand=$1
    if [ $# -gt 0 ]; then
        shift
    fi

    case "$subcommand" in
        provision )
            pushd "$PROVISIONER_DIRECTORY" > /dev/null
            ./devbox/provision $@
            popd > /dev/null
            ;;
        backup-devbox )
            pushd "$PROVISIONER_DIRECTORY" > /dev/null
            ./devbox/backup-config $@
            popd > /dev/null
            ;;
        restore-devbox )
            pushd "$PROVISIONER_DIRECTORY" > /dev/null
            ./devbox/restore-config $@
            popd > /dev/null
            ;;
        pull-changes )
            pushd "$PROVISIONER_DIRECTORY" > /dev/null
            echo "Pulling changes for dev machine provisioner ..."
            git pull origin main
            popd > /dev/null
            ;;
        edit-config )
            editor "$(devbox global path)/devbox.json"
            ;;
        "" )
            cd "$PROVISIONER_DIRECTORY"
            ;;
    esac

}

_machine-completion()
{
    local cur="${COMP_WORDS[COMP_CWORD]}"
    local prev=${COMP_WORDS[COMP_CWORD-1]}
    case ${COMP_CWORD} in
        1)
            COMPREPLY=($(compgen -W "--help provision edit-config pull-changes" -- $cur))
            ;;
        *)
            case $prev in
                provision )
                    local local_tooling="--help"
                    COMPREPLY=($(compgen -W "$local_tooling" -- $cur))
                    ;;
            esac
            ;;
    esac
}

complete -F _machine-completion machine

