#!/usr/bin/env bash
set -euo pipefail
ROOT=$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)
cd "$ROOT"
mkdir -p "$ROOT/tmp/api-reference/reports"
export GOWORK=off
api_docs() {
  (cd hack/api-docs && go run . "$@" "$ROOT")
}
policy() {
  local sources
  sources=$(api_docs sources)
  while IFS=$'\t' read -r api source; do
    vacuum lint --no-update-check --remote=false --ruleset hack/api-docs/validation/rules.yaml --fail-severity error --min-score 0 --no-banner --no-style --details "$source" > "tmp/api-reference/reports/$api-vacuum.txt" 2>&1 || {
      cat "tmp/api-reference/reports/$api-vacuum.txt"; return 1;
    }
  done <<< "$sources"
}
case "${1:-build}" in
  check|generate) policy; api_docs "$1" ;;
  test) (cd hack/api-docs && go test ./...) ;;
  build|serve)
    hugo --destination tmp/api-reference/site --baseURL "${DOCS_URL:-http://localhost:1314}" --cleanDestinationDir
    node hack/flatten-and-resolve.js tmp/api-reference/site
    node hack/api-docs/verify-output.mjs tmp/api-reference/site
    if [[ "${1:-build}" == serve ]]; then
      exec python3 -m http.server "${DOCS_PORT:-1314}" --bind 127.0.0.1 --directory tmp/api-reference/site
    fi
    ;;
  *) printf '%s\n' 'Usage: run.sh check|generate|test|build|serve' >&2; exit 2 ;;
esac
