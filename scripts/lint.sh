#!/bin/bash
# Run rumdl and Vale on specific Markdown files.
# Usage: scripts/lint.sh <markdown-file> [markdown-file...]
#
# Scoped output — no repo-wide noise. For full repo validation, use:
#   docker buildx bake validate
set -uo pipefail

if [ $# -eq 0 ]; then
  echo "Usage: $0 <markdown-file> [markdown-file...]" >&2
  exit 1
fi

for file in "$@"; do
  if [[ "$file" != *.md ]]; then
    echo "Error: scripts/lint.sh only accepts Markdown files: $file" >&2
    exit 2
  fi
done

repo_root=$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")/.." && pwd)

exit_code=0

echo "=== rumdl ==="
if ! npx --no-install rumdl check "$@" 2>&1; then
  exit_code=1
fi

echo ""
echo "=== vale ==="
if ! vale --no-global --config="$repo_root/.vale.ini" sync 2>&1; then
  exit_code=1
elif ! vale --no-global --config="$repo_root/.vale.ini" "$@" 2>&1; then
  exit_code=1
fi

echo ""
echo "Review Vale warnings and suggestions on changed lines; they don't affect the exit status."

exit $exit_code
