#!/bin/bash
# Update the Github release notes to match the release notes from the docs.
set -euo pipefail

version="${1:-}"

if [ -z "$version" ]; then
    echo "Usage: $0 <version>"
    exit 1
fi

moby="moby/moby"
version="${version#v}" # remove the v prefix
major_version="${version%%.*}" # get the major version (e.g. 29.0.3 -> 29)
moby_tag="docker-v${version}"

docs_notes=$(mktemp   -t new)
github_notes=$(mktemp -t old)

# Get the release notes from the docs.
grep -A 10000 "## ${version}" "content/manuals/engine/release-notes/${major_version}.md" | \
    grep -m 2 -A 0 -B 10000 '^## ' | \
    sed '$d' | \
    sed '/{{< release-date /{N;d;}' \
    > "$docs_notes" || { echo "Release notes for ${version} not found" && exit 1; }

release_json=$(gh api repos/moby/moby/releases -q ".[] | select(.name==\"v${version}\" or .tag_name==\"$moby_tag\")")

if [ -z "$release_json" ]; then
    echo "Release not found"
    printf '\033[0;34mTo create the draft Github release run the following command:\033[0m\n\n'
    echo gh -R moby/moby release create --draft \
        --notes-file '"'$docs_notes'"' \
        --title '"'v${version}'"' \
        ${moby_tag}
    exit 1
fi

release_id=$(echo "$release_json" | jq -r '.id')
echo "$release_json" | jq -r '.body' | sed 's/\r$//' > "$github_notes"

docs_notes_diff=$(mktemp -t diff)
# Copy docs_notes content and ensure it has exactly 1 blank line at the end
sed -e :a -e '/^\n*$/{$d;N;ba' -e '}' "$docs_notes" | sed '$a\' > "$docs_notes_diff"

# Compare the release notes.
if diff -u --color=auto "$github_notes" "$docs_notes_diff"; then
    printf '\033[0;32mThe release notes are already up to date.\033[0m\n'
    exit 0
fi

echo '========================================'
printf '\033[0;34mTo update the release notes run the following command:\033[0m\n\n'
echo gh api "repos/moby/moby/releases/$release_id" -X PATCH -f '"body=$(cat "'$docs_notes'")"'
