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

    if git diff --quiet "./data/engine-cli/*.yaml"; then
        printf "\e[32m✅ Already up to date\e[0m\n"
        return 100
    fi

    echo -e "ℹ️ Changes detected:"
    git diff --stat "./data/engine-cli/*.yaml" || true

    NICE_GIT_REF=$(cd "$worktree_dir" && git describe --always --dirty) || return $?

    git add "./data/engine-cli/*.yaml"

    git commit -m "cli: sync docs with docker/cli $NICE_GIT_REF"
    
    printf "\e[32m✅ Committed changes\e[0m\n"
    return 0
}

main "$@"
