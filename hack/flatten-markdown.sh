#!/bin/sh
set -e

# Flatten markdown files from public/path/to/page/index.md to public/path/to/page.md
# This makes markdown output links work correctly

PUBLIC_DIR="${1:-public}"

[ -d "$PUBLIC_DIR" ] || { echo "Error: Directory $PUBLIC_DIR does not exist"; exit 1; }

find "$PUBLIC_DIR" -type f -name 'index.md' | while read -r file; do
  # Skip the root index.md
  [ "$file" = "$PUBLIC_DIR/index.md" ] && continue

  # Get the directory containing index.md
  dir="${file%/*}"

  # Move index.md to parent directory with directory name
  mv "$file" "${dir%/*}/${dir##*/}.md"
done

echo "Flattened markdown files in $PUBLIC_DIR"
