---
name: fix-issues
description: >
  Fix one or more GitHub issues by creating branches, writing fixes, and opening PRs.
  Use this skill whenever the user provides GitHub issue numbers and wants them fixed,
  or says things like "fix issue 1234", "address these issues", "create PRs for issues
  1234 and 5678". Triggers on any request involving GitHub issue numbers paired with
  fixing, addressing, resolving, or creating PRs. Also triggers for "fix #1234" shorthand.
---

# Fix Issues

Given one or more GitHub issue numbers from the docker/docs repository, fix each
issue end-to-end: analyze, branch, fix, commit, push, and open a PR.

## Workflow

### 1. Fetch all issues in parallel

Use `gh issue view <number> --repo docker/docs --json title,body,labels` for each
issue number. Launch all fetches in parallel since they're independent.

### 2. Fix each issue sequentially

Process each issue one at a time in the main context. Do NOT use background
subagents for this — they can't get interactive Bash permission approval, which
blocks all git operations. Sequential processing in the main context is faster
than agents that stall on permissions.

For each issue:

#### a. Analyze

Read the issue body to understand:
- Which file(s) need changes
- What the problem is
- What the fix should be

#### b. Create a branch

```bash
git checkout -b fix/issue-<number>-<short-description> main
```

Use a short kebab-case description derived from the issue title (3-5 words max).

#### c. Read and fix

- Read each affected file before editing
- Make the minimal change that addresses the issue
- Don't over-engineer or add unrequested improvements
- Follow the repository's STYLE.md and COMPONENTS.md guidelines

#### d. Format

Run prettier on every changed file:

```bash
npx prettier --write <file>
```

#### e. Self-review

Re-read the changed file to verify:
- The fix addresses the issue correctly
- No unintended changes were introduced
- Formatting looks correct

#### f. Commit

Stage only the changed files (not `git add -A`).

#### g. Push and create PR

```bash
git push -u origin fix/issue-<number>-<short-description>
```

Then create the PR:

```bash
gh pr create --repo docker/docs \
  --title "<Short summary matching commit>" \
  --body "$(cat <<'EOF'
## Summary

- <1-3 bullet points describing what changed and why>

Closes #<issue-number>

🤖 Generated with [Claude Code](https://claude.com/claude-code)
EOF
)"
```

#### h. Label and assign reviewers

```bash
gh pr edit <pr-number> --repo docker/docs \
  --add-label "status/review" \
  --add-reviewer docker/docs-team
```

### 3. Switch back to main

After all issues are processed:

```bash
git checkout main
```

### 4. Report results

Present a summary table:

| Issue | PR | Change |
|-------|-----|--------|
| #N | #M | Brief description |

## Important notes

- Always work from `main` as the base for each branch
- Each issue gets its own branch and PR — don't combine issues
- If an issue references a file that doesn't exist, check for renames or
  reorganization before giving up (files move around in this repo)
- Validation commands (`docker buildx bake lint vale`) are available but slow;
  only run them if the user asks or the changes are complex enough to warrant it
