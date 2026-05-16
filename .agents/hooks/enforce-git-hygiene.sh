#!/bin/bash
# PreToolUse hook for the Bash tool.
# Blocks dangerous git patterns that cause CI failures.
set -euo pipefail

input=$(cat)
command=$(echo "$input" | jq -r '.tool_input.command // empty')

if [ -z "$command" ]; then
  exit 0
fi

# Block 'git add .' / 'git add -A' / 'git add --all'
# These stage package-lock.json and other generated files.
if echo "$command" | grep -qE 'git\s+add\s+(\.|--all|-A)(\s|$|;)'; then
  echo "BLOCKED: do not use 'git add .' / 'git add -A' / 'git add --all'. Stage files explicitly by name." >&2
  exit 2
fi

exit 0
