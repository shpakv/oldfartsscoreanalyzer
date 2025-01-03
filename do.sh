#!/usr/bin/env bash
# shellcheck disable=SC3040,SC3043
# https://github.com/koalaman/shellcheck/wiki/Checks
###############################################################################
# Do - The Simplest Build Tool on Earth.
# Documentation and examples see https://github.com/8gears/do
###############################################################################

###############################################################################
# public functions (commands)
###############################################################################

build() {
  local out=''
    # separate directories to search by spaces
    find cmd -name main.go | grep -v buddy | while read -r file; do
      out="$(
        echo "$file" | \
        sed -e 's/^cmd/bin/g' \
            -e 's/\/cmd\//\/bin\//g' \
            -e 's/main.go/bootstrap/g'
      )"
        go build  -o "$out" "$file"
    done
}

###############################################################################
# private functions
###############################################################################

# See `help set`
# As of this writing, Dash does not support the pipefail option, so using Bash,
# but still maintaining POSIX compatibility.
set -o pipefail
set -o errexit
set -o nounset
set -o errtrace
set -o posix

_SELF="$(basename "$0")"
_QUIET=false
_AWS="aws --no-cli-pager"

_fail() {
  local msg="$1"
  if [ -n "$msg" ]; then
    _log "Failed: $msg"
  fi
  exit 1
}

_error() {
  local msg="$1"
  _log "Error: $msg"
  return 1
}

_warn() {
  local msg="$1"
  _log "Warning: $msg"
}

_info() {
  local msg="$1"
  _log "Info: $msg"
}

_log() {
  printf "%s\n" "$1" 1>&2
}

_testCmd() {
  command -v "$1" >/dev/null
}

_listCommands() {
  grep -E '^\w+().*{' "$_SELF" | \
    sed -e '/^_.*/d' -e 's/\(.*\)().*/\1/g' | sort
}

_usage() {
  cat <<EOT

Usage:

  $_SELF build - Произвести сборку всех приложений, хранящихся в директории cmd.

EOT
  exit 1
}

# required to run in the same directory, simplifies things
test -f ./do.sh ||
  _fail "execute $0 from the same directory, like ./do.sh"

if [ "$#" = "0" ]; then
  _usage
fi

if [ "$1" = "-q" ]; then
  _QUIET=true
  set +o xtrace
  shift
fi

if [ "$1" != "_listCommands" ]; then
  # Do not print commands when doing tab completion.
  if [ $_QUIET = false ]; then
    set -o xtrace
  fi
fi

"$@" # <- execute the task
