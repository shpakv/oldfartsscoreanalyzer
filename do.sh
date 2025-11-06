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
  local sorry_bro="${SORRY_BRO:-}"
  local telegram_token="${TELEGRAM_BOT_TOKEN:-}"
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
  echo ""

  # Find all main.go files (excluding buddy)
  find cmd -name main.go | grep -v buddy | while read -r file; do
    # Extract package name from directory (e.g., cmd/teambuilder/main.go -> teambuilder)
    local pkg_name
    pkg_name="$(basename "$(dirname "$file")")"

    # Extract base directory (e.g., cmd/logs/stats/main.go -> logs/stats)
    local pkg_path
    pkg_path="$(dirname "$file" | sed 's|^cmd/||')"

    _info "Building ${pkg_name}..."

    # Create output directory
    local out_dir="bin/${pkg_path}"
    mkdir -p "$out_dir"

    # Build for current platform
    local current_out="${out_dir}/${pkg_name}"
    if [ -n "$ldflags" ]; then
      eval go build "$ldflags" -o "$current_out" "$file"
    else
      go build -o "$current_out" "$file"
    fi
    _info "  ✓ Built: ${current_out} (current platform)"

    # Build for Linux (for AWS Lambda bootstrap)
    local linux_out="${out_dir}/bootstrap"
    if [ -n "$ldflags" ]; then
      # Combine -s -w with custom ldflags
      local linux_ldflags=$(echo "$ldflags" | sed 's|-ldflags="|-ldflags="-s -w|')
      eval GOOS=linux GOARCH=amd64 go build "$linux_ldflags" -o "$linux_out" "$file"
    else
      GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o "$linux_out" "$file"
    fi
    _info "  ✓ Built: ${linux_out} (linux/amd64 for Lambda)"

    # Build for Windows
    local windows_out="${out_dir}/${pkg_name}.exe"
    if [ -n "$ldflags" ]; then
      eval GOOS=windows GOARCH=amd64 go build "$ldflags" -o "$windows_out" "$file"
    else
      GOOS=windows GOARCH=amd64 go build -o "$windows_out" "$file"
    fi
    _info "  ✓ Built: ${windows_out} (windows/amd64)"

    echo ""
  done

  _info "✓ All builds completed successfully!"
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
