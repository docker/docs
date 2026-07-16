#!/bin/bash
# Run markdownlint and vale on specific files.
# Usage: scripts/lint.sh <file> [file...]
#
# Scoped output — no repo-wide noise. For full repo validation, use:
#   docker buildx bake validate
set -uo pipefail

if [ $# -eq 0 ]; then
  echo "Usage: $0 <file> [file...]" >&2
  exit 1
fi

exit_code=0

echo "=== markdownlint ==="
if ! npx markdownlint-cli "$@" 2>&1; then
  exit_code=1
fi

echo ""
echo "=== vale ==="
if ! vale "$@" 2>&1; then
  exit_code=1
fi

exit $exit_code
