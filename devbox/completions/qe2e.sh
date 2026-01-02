qe2e() {
    _usage() {
        cat << USAGE
Usage:
  qe2e

Available Commands:
  run                                     Triggers jenkins pipeline that configures the RMS and runs e2e tests

  clean                                   Deletes feature branch in rms


Required environment variables:
  JENKINS_USER                            Most likely your xxx@modell-aachen.de email address

  JENKINS_TOKEN                           Your user token generated on https://jenkins.modac.cloud/user/xxx@modell-aachen.de/configure

  RMS_AUTH_TOKEN                          Authentication token for the https://rms.modac.eu api (see 1password)


Flags:
  -h, --help                              Show this help message

  All subcommands:

  -b, --branch <branch_name>              Specify branch for which the end to tests should run
                                          (Default is the currently checked out QwikiContrib branch)

  Subcommand 'run':

  -w, --which  <all|regular|multisite>    Specify which end to end tests to run
                                          (Defaults to 'all')

  -l, --legacy                            Trigger e2e tests on legacy jenkins (https://jenkins.int.modac.eu)
                                          This is necessary for all end to end tests <= Q.wiki 5.11
                                          Multisite end to end tests will not be triggered
USAGE
    }

    _get_crumb() {
        curl -s -u ${JENKINS_USER}:${JENKINS_TOKEN} 'https://jenkins.modac.cloud/crumbIssuer/api/json' | jq -r .crumb
    }

    _trigger_run() {
        local branch=$1
        local buildParams="${2:-}"

        if [ -z "$JENKINS_USER" ]; then
            echo "Missing JENKINS_USER environment variable"
            return
        fi

        if [ -z "$JENKINS_TOKEN" ]; then
            echo "Missing JENKINS_TOKEN environment variable"
            return
        fi

        local crumb=$(_get_crumb)

        echo "Triggering end to end tests on jenkins.modac.cloud for branch '${branch}' with params '${buildParams}' ..."
        curl -X POST -s -u ${JENKINS_USER}:${JENKINS_TOKEN} -H "Jenkins-Crumb:${crumb}" "https://jenkins.modac.cloud/job/qwiki_e2e_pre_merge/buildWithParameters?BRANCH=${branch}${buildParams}"
    }

    _trigger_run_legacy() {
        local branch=$1
        local buildParams=$2

        echo "Triggering end to end tests on jenkins.int.modac.eu (legacy) for branch '${branch}' with params '${buildParams}' ..."
        curl "https://jenkins.int.modac.eu/job/Pre-merge%20E2E%20Test%20Run/buildWithParameters?BRANCH=${branch}${buildParams}"
    }

    _trigger_run_multisite() {
        local branch=$1
        local e2eConfigId="1554816780804" # ModellAachenDev master/multisite latest
        local integrationFolder="cypress/integration_multisite"
        local cypressProjectId="bq1n27"
        _trigger_run $branch "${buildParams}&E2E_RMS_CONFIG_ID=${e2eConfigId}&INTEGRATION_FOLDER=${integrationFolder}&PARALLELIZATION_DEPTH=1&CYPRESS_PROJECT_ID=${cypressProjectId}&ORGANIZATIONAL_UNIT=true"
    }

    _curl_rms() {
        local method=$1
        local endpoint=$2
        echo "$(curl -X "${method}" -s "https://rms.modac.eu/api/${endpoint}" -H "rms-auth-token: ${RMS_AUTH_TOKEN}")"
    }

    _clean_rms_branch() {
        local branch=$1

        if [ -z "$RMS_AUTH_TOKEN" ]; then
            echo "Missing RMS_AUTH_TOKEN environment variable"
            return
        fi

        echo "Deleting RMS feature branch ${branch}..."
        branchId=$(_curl_rms 'GET' 'featureBranch/get_all' | jq -r ".[] | select(.name == \"${branch}\") | .id")
        _curl_rms 'DELETE' "featureBranch/${branchId}" >/dev/null
    }

    cd "$REPOS_DIRECTORY/QwikiContrib" >/dev/null
    local branch=$(git symbolic-ref --short HEAD)
    cd - >/dev/null

    local legacy=false
    local which=all
    local subcommand=""
    while [[ "$1" != "" ]] ; do
        case "$1" in
            -h | --help )
                _usage
                return
                ;;
            -b | --branch )
                branch=$2
                shift
                ;;
            -w | --which )
                which=$2
                shift
                ;;
            -l | --legacy )
                legacy=true
                ;;
            run | clean )
                subcommand=$1
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

    case "$subcommand" in
        run )
            if [ "$legacy" = true ] ; then
                _trigger_run_legacy $branch ""
            else
                if [ "$which" = "all" ] || [ "$which" = "regular" ] ; then
                    _trigger_run $branch
                fi

                if [ "$which" = "all" ] || [ "$which" = "multisite" ] ; then
                    _trigger_run_multisite $branch
                fi
            fi
            return
            ;;
        clean )
            _clean_rms_branch $branch
            return
            ;;
    esac
}

_qe2e-completion() {
    local cur="${COMP_WORDS[COMP_CWORD]}"

    case ${COMP_CWORD} in
        1)
            COMPREPLY=($(compgen -W "run clean --help" -- $cur))
            ;;
        *)
            local subcommand="${COMP_WORDS[COMP_CWORD-1]}"
            case $subcommand in
                -b | --branch )
                    cd "$REPOS_DIRECTORY/QwikiContrib" >/dev/null
                    local branches=$(git branch --format="%(refname)" -a | sed -e "s#^refs/\(heads\|remotes/origin\)/##g")
                    cd - >/dev/null
                    COMPREPLY=($(compgen -W "$branches" -- $cur))
                    ;;
                -w | --which )
                    COMPREPLY=($(compgen -W "regular multisite all" -- $cur))
                    ;;
                *)
                    COMPREPLY=($(compgen -W "--branch --which --legacy" -- $cur))
                    ;;
            esac
            ;;
    esac
}

complete -F _qe2e-completion qe2e
