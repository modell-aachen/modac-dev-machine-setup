qypress () {
    usage() {
        cat << USAGE
Usage:
    qypress [options] <host> <cypress command> [-- options for cypress]

Options:
    -h, --help          shows this message
    -m, --multisite     runs multisite fixtures
    -r, --remote        runs E2E Tests against cypress dashboard reports to locally checked out branch.
                        only works for RMS testsystems (Mailpit).

Example
    qypress -r "https://e2ehotfix.testing.modac.eu" run     runs all E2E on rms testsystem e2ehotfix and reports on the branch you have checked out locally
USAGE
    }

    OPTS=`getopt -o h --long help -o m --long multisite -o r --long remote -- "$@"`
    if [ $? != 0 ] ; then echo "Failed parsing options." >&2 ; return ; fi

    eval set -- "$OPTS"

    local integrationFolder="cypress/integration"
    local organizationalUnit="false"
    local remote="false"

    while true; do
        case "$1" in
            -h | --help )
                usage
                return
                shift ;;
            -m | --multisite )
                organizationalUnit="true"
                integrationFolder="cypress/integration_multisite"
                shift ;;
            -r | --remote )
                remote="true"
                shift ;;
            -- )
                shift
                break ;;
            * )
                break ;;
        esac
    done

    shift $(expr $OPTIND - 1 )
    local host=$1
    shift

    pushd "$REPOS_DIRECTORY/QwikiContrib/e2e-tests"

    local remainingArgs=$@

    local mailPitUrl="$host:8025"
    if [[ "$host" == *qluster.localhost ]]; then
        mailPitUrl="https://mailpit.qluster.localhost"
        oidcMockAuthorizeUrl="https://oidc.qluster.localhost/connect/authorize"
        oidcMockTokenUrl="http://qwiki-local-oidc.default.svc.cluster.local/connect/token"
    fi

    yarn
    if [[ "$remote" == "true" ]];
    then
        if [ -z "$CYPRESS_RECORD_KEY" ]; then
            echo "Missing CYPRESS_RECORD_KEY environment variable"
            return
        fi
        if [ -z "$CYPRESS_PROJECT_ID" ]; then
            echo "Missing CYPRESS_PROJECT_ID environment variable"
            return
        fi

        mailPitUrl=$(sed "s/^https:\/\//http:\/\/mailpit./g" <<< $host)

        CYPRESS_baseUrl=$host/ \
        CYPRESS_mailPitUrl=$mailPitUrl \
        CYPRESS_oidcMockAuthorizeUrl=$oidcMockAuthorizeUrl \
        CYPRESS_oidcMockTokenUrl=$oidcMockTokenUrl \
        CYPRESS_integrationFolder=$integrationFolder \
        CYPRESS_organizationalUnit=$organizationalUnit \
        CYPRESS_PROJECT_ID=$CYPRESS_PROJECT_ID \
            yarn cypress $@ --record --key $CYPRESS_RECORD_KEY

    else

        CYPRESS_baseUrl=$host/ \
        CYPRESS_mailPitUrl=$mailPitUrl \
        CYPRESS_oidcMockAuthorizeUrl=$oidcMockAuthorizeUrl \
        CYPRESS_oidcMockTokenUrl=$oidcMockTokenUrl \
        CYPRESS_integrationFolder=$integrationFolder \
        CYPRESS_organizationalUnit=$organizationalUnit \
            yarn cypress $@
    fi

    popd
}

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
