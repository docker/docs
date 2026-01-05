#!/bin/bash
set -euo pipefail

main() {
    local branch_name="${1:-upstream/master}"
    local cli_source="${2:-$HOME/src/cli}/.git"
    local worktree_dir="./internal-update-cli-docs"

    (
        set -e
        GIT_DIR="$cli_source"
        GIT_DIR="$GIT_DIR" git fetch upstream
        GIT_DIR="$GIT_DIR" git worktree add "$worktree_dir" "$branch_name"
    ) || return $?
    trap "GIT_DIR=\"$cli_source\" git worktree remove \"$worktree_dir\" --force 2>/dev/null || true" EXIT

    (set -e; cd "$worktree_dir"; make -f docker.Makefile yamldocs || { printf "::error::Failed to generate YAML docs!\n"; exit 1; }) || return $?
    cp "$worktree_dir"/docs/yaml/*.yaml ./data/engine-cli/

    NICE_GIT_REF=$(cd "$worktree_dir" && git describe --always --dirty) || return $?
    echo "git_ref=$NICE_GIT_REF" >> "${GITHUB_OUTPUT:-/dev/stdout}"

    printf "\e[32m✅ Synced CLI docs from %s\e[0m\n" "$NICE_GIT_REF"
    return 0
}

main "$@"
