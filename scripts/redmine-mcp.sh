#!/bin/sh
# Launcher for the rdflow MCP server, for MCP host configs to use as the
# command (sh redmine-mcp.sh mcp serve). GUI-launched hosts often get a
# minimal PATH, so look in the usual install locations when `rdflow` is not
# on PATH.

if command -v rdflow >/dev/null 2>&1; then
  exec rdflow "$@"
fi

# _REDMINE_MCP_TEST_SYSTEM_DIRS replaces the fixed system directories. It is
# for the tests only and not a supported setting.
for dir in ${_REDMINE_MCP_TEST_SYSTEM_DIRS:-/opt/homebrew/bin /usr/local/bin} "${HOME:-}/.local/bin" "${GOBIN:-}" "${HOME:-}/go/bin"; do
  case "$dir" in
    /*) ;;
    *) continue ;;
  esac
  if [ -f "$dir/rdflow" ] && [ -x "$dir/rdflow" ]; then
    exec "$dir/rdflow" "$@"
  fi
done

# Builtins only: PATH may not contain the directory with cat.
printf '%s\n' \
  "rdflow: the 'rdflow' binary was not found on PATH or in the usual install locations." \
  "" \
  "Install it with:" \
  "  curl -fsSL https://raw.githubusercontent.com/TruthUnderlier/RedmineFlow/main/install.sh | REDMINE_CLI_SKIP_ATTESTATION=1 bash" \
  "" \
  "Docs: https://github.com/TruthUnderlier/RedmineFlow" >&2
exit 127
