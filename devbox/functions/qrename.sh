qrename() {

    usage() {
        cat << USAGE
Usage:
  qrename [source] [destination]

Examples:
  qrename Qwiki/Tenant Qwiki/Administration/Model/Tenant
  qrename Qwiki/Tenant.pm Qwiki/Administration/Model/Tenant.pm

Flags:
  -h, --help   shows this help message
      --dry    shows found source packages
USAGE
    }

    local isRenamed=1

    OPTS=`getopt -o h --long help,dry -- "$@"`
    if [ $? != 0 ] ; then echo "Failed parsing options." >&2 ; return 1 ; fi

    eval set -- "$OPTS"

    while true; do
        case "$1" in
            --dry )
                isRenamed=0
                shift ;;
            -h | --help )
                usage
                return ;;
            -- )
                shift
                break ;;
            * )
                break ;;
        esac
    done

    shift $(expr $OPTIND - 1 )
    local src=$1
    local dest=$2

    module() {
        local path=$1
        echo "$path" | sed "s#/#::#g" | sed 's#\.pm$##'
    }

    allFiles() {
        find . \( -name "*.pl" -o -name "*.pm" -o -name "*.t" -o -wholename "./tools/*" -o -wholename "./bin/*" \) -exec $@ {} \;
    }

    local qwikiContrib=$REPOS_DIRECTORY/QwikiContrib
    cd "$qwikiContrib/core"

    local srcModule=$(module "$src")
    local destModule=$(module "$dest")

    if [ $isRenamed = "0" ]; then
        if [ -d "$qwikiContrib/core/lib/$src" ]; then
            allFiles grep --color "\([^:][^:]\)$srcModule::[A-Z]"
        else
            allFiles grep --color "\([^:][^:]\)$srcModule\([^[:alnum:]:]\|::[a-z]\|$\)"
        fi
    else
        echo "renames from '$srcModule' to '$destModule'"

        if [ -d "$qwikiContrib/core/lib/$src" ]; then
            allFiles sed -i "s/\([^:][^:]\)$srcModule\(::[A-Z]\)/\1$destModule\2/g"
            mkdir -p "lib/$dest"
            mv lib/$src/* lib/$dest/
        else
            allFiles sed -i "s/\([^:][^:]\)$srcModule\([^[:alnum:]:]\|::[a-z]\|$\)/\1$destModule\2/g"

            mkdir -p "$(dirname "lib/$dest")"
            mv "lib/$src" "lib/$dest"
        fi
    fi

    cd - >/dev/null
}

_qrename-completion()
{
    cd "$REPOS_DIRECTORY/QwikiContrib/core/lib"

    local cur="${COMP_WORDS[COMP_CWORD]}"

    local depth=$( tr -dc '/' <<< "/$cur" | wc -c )

    local completions=$(find . -maxdepth $depth -wholename "./$cur*" | sed "s#^\./##")

    COMPREPLY=($(compgen -W "-h --help --dry $completions" -- $cur))
}

complete -F _qrename-completion qrename
