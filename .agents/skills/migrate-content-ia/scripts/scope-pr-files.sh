#!/bin/bash
# List files changed on the current branch since merge-base with the target
# branch — suitable for PR_SCOPE_FILES in PR-scoped migrate-content-ia runs.
#
# Usage:
#   ./scope-pr-files.sh [target-branch]
# Default target-branch: main
#
# Example (Bash / Git Bash, from repo root):
#   bash .agents/skills/migrate-content-ia/scripts/scope-pr-files.sh
# Example (other base):
#   bash .../scope-pr-files.sh upstream/main
set -euo pipefail

target="${1:-main}"

if ! base=$(git merge-base "$target" HEAD 2>/dev/null); then
  echo "Error: could not merge-base with '$target'. Fetch remotes or pass a valid branch." >&2
  exit 1
fi

git diff --name-only "$base"...HEAD
