#!/bin/bash
# PostToolUse hook for Edit/Write in the write skill.
# Auto-formats Markdown files with rumdl after each edit.
set -euo pipefail

input=$(cat)
file_path=$(echo "$input" | jq -r '.tool_input.file_path // empty')

[ -z "$file_path" ] && exit 0
[[ "$file_path" != *.md ]] && exit 0

npx --no-install rumdl fmt "$file_path" 2>/dev/null
