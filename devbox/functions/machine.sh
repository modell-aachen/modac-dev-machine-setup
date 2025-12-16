machine() {

    usage() {
        cat << USAGE
Usage:
  machine [command]

Available Commands:
  edit-config    edit custom inventory file
  provision      provision this machine
                   Options:
                     --filter <modules>    Comma-separated list of modules to run (implies --skip-install)
                     --skip-install        Skip running install, devbox shellenv, and op signin
                     --help                Show provision help
                   Subcommands:
                     list-modules       List all available modules
  pull-changes   pull changes for dev machine provisioner
  backup         create / restore backup files to/from 1Password configured in devbox.json

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
        backup )
            pushd "$PROVISIONER_DIRECTORY" > /dev/null
            ./devbox/backup $@
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
            COMPREPLY=($(compgen -W "--help provision edit-config pull-changes backup" -- $cur))
            ;;
        *)
            case $prev in
                provision )
                    local local_tooling="--help --filter list-modules"
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

