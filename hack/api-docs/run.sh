#!/usr/bin/env bash
set -euo pipefail
ROOT=$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)
cd "$ROOT"
BIN="$ROOT/tmp/api-prototype/bin"
mkdir -p "$BIN" "$ROOT/tmp/api-prototype/reports"
export GOWORK=off
bootstrap() {
  npm --prefix hack/api-docs ci --ignore-scripts --no-audit --no-fund
  (cd hack/api-docs && go build -o "$BIN/api-docs" .)
  if [[ ! -x "$BIN/vacuum" ]] || ! "$BIN/vacuum" version 2>&1 | grep -qE '0\.30\.3'; then
    GOBIN="$BIN" go install github.com/daveshanley/vacuum@v0.30.3
  fi
  if [[ ! -x "$BIN/oasdiff" ]] || ! "$BIN/oasdiff" --version 2>&1 | grep -qE '1\.31\.0'; then
    GOBIN="$BIN" go install github.com/oasdiff/oasdiff@v1.31.0
  fi
}
policy() {
  for api in governance hub dvp registry engine-1.56 engine-1.55; do
    if ! "$BIN/vacuum" lint --no-update-check --remote=false --ruleset hack/api-docs/vacuum.yaml --fail-severity error --min-score 0 --no-banner --no-style --details "prototypes/api-docs/converted/$api.yaml" > "tmp/api-prototype/reports/$api-vacuum.txt" 2>&1; then
      cat "tmp/api-prototype/reports/$api-vacuum.txt"; return 1
    fi
  done
}
index_site() {
  if [[ -n "${PAGEFIND_BIN:-}" ]]; then
    "$PAGEFIND_BIN" --site tmp/api-prototype/site --output-path tmp/api-prototype/site/pagefind
  elif [[ "$(uname -s)" == Linux && "$(getconf PAGESIZE)" -gt 4096 ]]; then
    if [[ ! -x tmp/api-prototype/pagefind-portable/pagefind ]]; then
      mkdir -p tmp/api-prototype/empty
      docker buildx build --file hack/api-docs/pagefind.Dockerfile --output type=local,dest=tmp/api-prototype/pagefind-portable tmp/api-prototype/empty
    fi
    tmp/api-prototype/pagefind-portable/pagefind --site tmp/api-prototype/site --output-path tmp/api-prototype/site/pagefind
  else
    npx --yes pagefind@1.5.2 --site tmp/api-prototype/site --output-path tmp/api-prototype/site/pagefind
  fi
}
prepare() {
  bootstrap
  node hack/api-docs/migrate.mjs verify
  policy
  "$BIN/api-docs" generate "$ROOT" --preview
  "$BIN/oasdiff" diff prototypes/api-docs/converted/engine-1.55.yaml prototypes/api-docs/converted/engine-1.56.yaml --format json > tmp/api-prototype/reports/engine-diff.json
  "$BIN/oasdiff" changelog prototypes/api-docs/converted/engine-1.55.yaml prototypes/api-docs/converted/engine-1.56.yaml --format markdown > tmp/api-prototype/reports/engine-changelog.md
  for api in hub dvp registry governance engine-1.55 engine-1.56; do
    input="prototypes/api-docs/original/$api.yaml"
    [[ "$api" != engine-* ]] || input="prototypes/api-docs/migrations/$api.intermediate.json"
    # Original Registry lacks info.version and original Hub has an unresolved ref.
    # Preserve tool failures as reports; the replay ledger remains mandatory.
    if ! "$BIN/oasdiff" diff "$input" "prototypes/api-docs/converted/$api.yaml" --format json > "tmp/api-prototype/reports/$api-migration-diff.json" 2> "tmp/api-prototype/reports/$api-migration-diff.stderr"; then
      printf '%s\n' 'Comparison unavailable; inspect stderr and the complete replay ledger.' > "tmp/api-prototype/reports/$api-migration-diff.status"
    fi
  done
  node hack/api-docs/report.mjs
}
case "${1:-build}" in
  bootstrap) bootstrap ;;
  verify) node hack/api-docs/migrate.mjs verify ;;
  convert) node hack/api-docs/migrate.mjs convert ;;
  check) bootstrap; node hack/api-docs/migrate.mjs verify; policy; "$BIN/api-docs" check "$ROOT" ;;
  generate) prepare ;;
  report) node hack/api-docs/report.mjs ;;
  test) (cd hack/api-docs && go test ./...); node hack/api-docs/migrate.mjs verify ;;
  build|serve)
    prepare
    [[ -d node_modules/alpinejs ]] || npm ci --no-audit --no-fund
    started=$SECONDS
    rm -rf "$ROOT/tmp/api-prototype/site"
    hugo --config hugo.yaml,prototypes/api-docs/hugo.yaml --destination tmp/api-prototype/site --baseURL "${API_PROTOTYPE_URL:-http://localhost:1314}" --environment api-prototype
    node hack/api-docs/flatten-preview.mjs
    index_site
    node hack/api-docs/verify-output.mjs
    printf '{"buildSeconds":%s}\n' "$((SECONDS-started))" > tmp/api-prototype/timing.json
    if [[ "${1:-build}" == serve ]]; then
      exec python3 -m http.server "${API_PROTOTYPE_PORT:-1314}" --bind 127.0.0.1 --directory tmp/api-prototype/site
    fi
    ;;
  *) printf '%s\n' 'Usage: run.sh bootstrap|verify|convert|check|generate|report|test|build|serve' >&2; exit 2 ;;
esac
