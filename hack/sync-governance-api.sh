#!/bin/bash
# Re-vendor the AI Governance OpenAPI spec from docker/governor-services.
#
# The spec at content/reference/api/ai-governance/api.yaml is a verbatim copy
# of the upstream openapi.yaml — no transformation is applied. This script
# fetches the latest version and overwrites the local copy, then prints a diff
# summary so you can sanity-check before committing.
#
# Requires: gh (authenticated with access to the private docker/governor-services
# repo). No secrets are stored in this repo — auth comes from your own gh login.
#
# Usage: hack/sync-governance-api.sh [ref]
#   ref   optional git ref (branch, tag, or SHA). Defaults to the repo default branch.
set -euo pipefail

REPO="docker/governor-services"
SRC_PATH="governor-service-api/openapi.yaml"
DEST="content/reference/api/ai-governance/api.yaml"
REF="${1:-}"

# Resolve this repo's root so the script works from any working directory.
ROOT="$(git -C "$(dirname "$0")" rev-parse --show-toplevel)"
DEST_ABS="$ROOT/$DEST"

if ! command -v gh >/dev/null 2>&1; then
  echo "error: gh CLI not found. Install it and run 'gh auth login'." >&2
  exit 1
fi

api_path="repos/$REPO/contents/$SRC_PATH"
[ -n "$REF" ] && api_path="$api_path?ref=$REF"

echo "Fetching $SRC_PATH from $REPO${REF:+@$REF}..."
tmp="$(mktemp)"
trap 'rm -f "$tmp"' EXIT

if ! gh api "$api_path" --jq '.content' | base64 -d > "$tmp"; then
  echo "error: failed to fetch spec. Check 'gh auth status' and repo access." >&2
  exit 1
fi

if [ ! -s "$tmp" ]; then
  echo "error: fetched spec is empty — aborting without overwriting $DEST." >&2
  exit 1
fi

if diff -q "$DEST_ABS" "$tmp" >/dev/null 2>&1; then
  echo "Already up to date — $DEST is identical to upstream."
  exit 0
fi

cp "$tmp" "$DEST_ABS"
echo "Updated $DEST. Review the change:"
echo
git -C "$ROOT" --no-pager diff --stat -- "$DEST"
