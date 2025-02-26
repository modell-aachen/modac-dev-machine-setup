machine() {

    usage() {
        cat << USAGE
Usage:
  machine [command]

Available Commands:
  edit-config   edit custom inventory file
  provision     provision this machine

Flags:
  -h, --help   shows this help message
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
    local machineDir=$HOME/modac-dev-machine-setup
    local inventoryFile=$HOME/.inventory_local.yml

    case "$subcommand" in
        provision )
            cd $machineDir
            ./dev-provision -i $inventoryFile $@
            cd - >/dev/null
            ;;
        edit-config )
            editor $inventoryFile
            ;;
        "" )
            cd $machineDir
            ;;
    esac

}

_machine-completion()
{
    local cur="${COMP_WORDS[COMP_CWORD]}"
    local prev=${COMP_WORDS[COMP_CWORD-1]}
    case ${COMP_CWORD} in
        1)
            COMPREPLY=($(compgen -W "--help provision edit-config" -- $cur))
            ;;
        *)
            case $prev in
                provision )
                    local local_tooling="--help dotfiles packages tooling scripts upgrade-provisioner"
                    COMPREPLY=($(compgen -W "$local_tooling" -- $cur))
                    ;;
            esac
            ;;
    esac
}

complete -F _machine-completion machine

