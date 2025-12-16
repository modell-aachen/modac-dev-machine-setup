#!/usr/bin/env bash

# Logging configuration from environment (set by main provision script)
LOG_FILE="${PROVISION_LOG_FILE:-}"
LOG_LEVEL="${PROVISION_LOG_LEVEL:-INFO}"
LOG_VERBOSE="${PROVISION_VERBOSE:-false}"

# Color codes
LOG_COLOR_RESET='\033[0m'
LOG_COLOR_DEBUG='\033[0;36m'
LOG_COLOR_INFO='\033[0;32m'
LOG_COLOR_WARN='\033[1;33m'
LOG_COLOR_ERROR='\033[1;91m'
LOG_COLOR_MODULE='\033[1;35m'

# Log level priorities
declare -A LOG_LEVELS=( [DEBUG]=0 [INFO]=1 [WARN]=2 [ERROR]=3 )

function log_timestamp() {
    date '+%Y-%m-%d %H:%M:%S'
}

function log_should_print() {
    local level=$1
    local current_priority=${LOG_LEVELS[$LOG_LEVEL]:-1}
    local msg_priority=${LOG_LEVELS[$level]:-1}

    if [ "$level" = "DEBUG" ] && [ "$LOG_VERBOSE" != "true" ]; then
        return 1
    fi

    [ $msg_priority -ge $current_priority ]
}

function log_write() {
    local level=$1
    local message=$2
    local color=$3

    local timestamp=$(log_timestamp)
    local log_entry="[$timestamp] [$level] $message"

    # Write to log file if configured (without color codes)
    if [ -n "$LOG_FILE" ]; then
        echo "$log_entry" >> "$LOG_FILE"
    fi

    # Print to stdout with color
    echo -e "${color}[$timestamp] [$level]${LOG_COLOR_RESET} $message"
}

function log_debug() {
    log_should_print "DEBUG" && log_write "DEBUG" "$1" "$LOG_COLOR_DEBUG"
}

function log_info() {
    log_should_print "INFO" && log_write "INFO" "$1" "$LOG_COLOR_INFO"
}

function log_warn() {
    log_should_print "WARN" && log_write "WARN" "$1" "$LOG_COLOR_WARN"
}

function log_error() {
    log_should_print "ERROR" && log_write "ERROR" "$1" "$LOG_COLOR_ERROR" >&2
}

function log_module_start() {
    local module=$1
    log_write "MODULE" "========== Starting: $module ==========" "$LOG_COLOR_MODULE"
}

function log_module_end() {
    local module=$1
    local duration=$2
    log_write "MODULE" "========== Completed: $module (${duration}s) ==========" "$LOG_COLOR_MODULE"
}

# Export logging functions so executed scripts can use them
export -f log_timestamp
export -f log_should_print
export -f log_write
export -f log_debug
export -f log_info
export -f log_warn
export -f log_error
export -f log_module_start
export -f log_module_end

function install_completion() {
    local shell=$1
    local cmd=$2
    local version=$3
    local shell_path="$HOME/.${shell}rc"
    local completions_path="$HOME/.${shell}_completions"
    local cmd_completion_path="$completions_path/${cmd}_$version.sh"


    if [[ -f "$shell_path" && ! -f "$cmd_completion_path" ]]; then
        mkdir -p "$completions_path"
        rm -f "$completions_path/$cmd"*".sh"

        echo "Installing $cmd completion for $shell under $cmd_completion_path"
        "$cmd" completion "$shell" > "$cmd_completion_path"
    fi
}
