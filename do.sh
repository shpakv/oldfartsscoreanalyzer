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
  local sorry_bro="${SORRY_BRO:-}"
  local telegram_token="${TELEGRAM_TOKEN:-}"
  local ldflags=""

  # Build ldflags if variables are set
  if [ -n "$sorry_bro" ] || [ -n "$telegram_token" ]; then
    ldflags="-ldflags="
    local flags=""

    if [ -n "$sorry_bro" ]; then
      flags="${flags} -X 'main.SorryBro=${sorry_bro}'"
    fi

    if [ -n "$telegram_token" ]; then
      flags="${flags} -X 'oldfartscounter/internal/telegram.Token=${telegram_token}'"
    fi

    ldflags="${ldflags}\"${flags}\""
  fi

  _info "Building with:"
  _info "  SorryBro: ${sorry_bro:-<not set>}"
  _info "  Token: ${telegram_token:+<set>}${telegram_token:-<not set>}"

  # separate directories to search by spaces
  find cmd -name main.go | grep -v buddy | while read -r file; do
    out="$(
      echo "$file" | \
      sed -e 's/^cmd/bin/g' \
          -e 's/\/cmd\//\/bin\//g' \
          -e 's/main.go/oldfartbuilder/g'
    )"

    _info "Building $(basename "$(dirname "$file")")..."

    if [ -n "$ldflags" ]; then
      eval go build "$ldflags" -o "$out" "$file"
    else
      go build -o "$out" "$file"
    fi

    _info "✓ Built: $out"
  done
}

terInit() {
  terraform -chdir=./terraform init
}

plan() {
  terInit
  terraform -chdir=./terraform plan
}

deploy() {
  terInit
  terraform -chdir=./terraform apply -auto-approve
}

destroy() {
  terInit
  terraform -chdir=./terraform destroy -auto-approve
}

preCommit() {
  _info "Running pre-commit checks..."

  _info "Step 1/3: Running goimports..."
  if ! _testCmd goimports; then
    _warn "goimports not found, installing..."
    go install golang.org/x/tools/cmd/goimports@latest
  fi
  goimports -w .

  _info "Step 2/3: Running golangci-lint..."
  if ! _testCmd golangci-lint; then
    _error "golangci-lint not found. Please install: https://golangci-lint.run/usage/install/"
  fi
  golangci-lint run

  _info "Step 3/3: Running tests..."
  go test ./...

  _info "✓ All pre-commit checks passed!"
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

Build with custom variables:

  SORRY_BRO="PlayerName" TELEGRAM_TOKEN="bot123:ABC" $_SELF build

Examples:

  # Build without variables (uses default/empty values)
  $_SELF build

  # Build with SorryBro player
  SORRY_BRO="Looka" $_SELF build

  # Build with Telegram token
  TELEGRAM_TOKEN="123456789:ABCdefGHIjklMNOpqrsTUVwxyz" $_SELF build

  # Build with both
  SORRY_BRO="Mr. Titspervert" TELEGRAM_TOKEN="bot_token_here" $_SELF build

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
