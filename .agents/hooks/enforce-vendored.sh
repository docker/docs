#!/bin/bash
# PreToolUse hook for Edit and Write tools.
# Blocks edits to vendored content that must be fixed upstream.
set -euo pipefail

input=$(cat)
file_path=$(echo "$input" | jq -r '.tool_input.file_path // empty')

if [ -z "$file_path" ]; then
  exit 0
fi

# Block edits to vendored Hugo modules.
if echo "$file_path" | grep -qE '/_vendor/'; then
  echo "BLOCKED: _vendor/ is vendored from upstream Hugo modules. Fix in the source repo instead." >&2
  exit 2
fi

# Block edits to vendored CLI reference data.
if echo "$file_path" | grep -qE '/data/cli/'; then
  echo "BLOCKED: data/cli/ is generated from upstream repos (docker/cli, docker/buildx, etc.). Fix in the source repo instead." >&2
  exit 2
fi

exit 0
