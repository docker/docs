#!/usr/bin/env bash
set -euo pipefail
ROOT=$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)
cd "$ROOT"
BIN="$ROOT/tmp/api-reference/bin"
mkdir -p "$BIN" "$ROOT/tmp/api-reference/reports"
export GOWORK=off
bootstrap() {
  (cd hack/api-docs && go build -o "$BIN/api-docs" .)
  if [[ ! -x "$BIN/vacuum-v0.30.3" ]]; then
    GOBIN="$BIN" go install github.com/daveshanley/vacuum@v0.30.3
    mv "$BIN/vacuum" "$BIN/vacuum-v0.30.3"
  fi
}
policy() {
  local sources
  sources=$("$BIN/api-docs" sources "$ROOT")
  while IFS=$'\t' read -r api source; do
    "$BIN/vacuum-v0.30.3" lint --no-update-check --remote=false --ruleset hack/api-docs/validation/rules.yaml --fail-severity error --min-score 0 --no-banner --no-style --details "$source" > "tmp/api-reference/reports/$api-vacuum.txt" 2>&1 || {
      cat "tmp/api-reference/reports/$api-vacuum.txt"; return 1;
    }
  done <<< "$sources"
}
generate() {
  bootstrap
  policy
  "$BIN/api-docs" generate "$ROOT" --allow-known-issues
}
case "${1:-build}" in
  bootstrap) bootstrap ;;
  check) bootstrap; policy; "$BIN/api-docs" check "$ROOT" ;;
  generate) generate ;;
  test) (cd hack/api-docs && go test ./...) ;;
  build|serve)
    generate
    hugo --destination tmp/api-reference/site --baseURL "${DOCS_URL:-http://localhost:1314}" --cleanDestinationDir
    node hack/api-docs/flatten.mjs tmp/api-reference/site
    node hack/api-docs/verify-output.mjs tmp/api-reference/site
    if [[ "${1:-build}" == serve ]]; then
      exec python3 -m http.server "${DOCS_PORT:-1314}" --bind 127.0.0.1 --directory tmp/api-reference/site
    fi
    ;;
  *) printf '%s\n' 'Usage: run.sh bootstrap|check|generate|test|build|serve' >&2; exit 2 ;;
esac
