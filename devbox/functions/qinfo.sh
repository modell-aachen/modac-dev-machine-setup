qinfo () {
    usage() {
        cat << USAGE
Usage:
    qinfo [options] <tenantId>

Options:
    -h, --help          shows this message

Example
    qinfo e8307f3e-612c-44a5-81e2-c6c64d5cf591    returns tenant information for tenant with the specific tenant id.
USAGE
    }

    OPTS=`getopt -o h --long help -- "$@"`
    if [ $? != 0 ] ; then echo "Failed parsing options." >&2 ; return ; fi

    eval set -- "$OPTS"

    while true; do
        case "$1" in
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
    local tenantId=$1
    shift

    QWIKI_POD=$(kubectl get pods -l app.kubernetes.io/component=frontend --no-headers -o custom-columns=":metadata.name" --context gcpeuw3-gke-qwikiproduction --namespace qwiki | head -n 1)
    kubectl exec $QWIKI_POD --context gcpeuw3-gke-qwikiproduction --namespace qwiki -- /var/www/qwikis/core/tools/tenant info --tenant-id $tenantId
}
