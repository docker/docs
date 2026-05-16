---
name: write
description: >
  Write a documentation fix on a branch. Makes the minimal change, formats,
  self-reviews, and commits. Use after research has identified what to change.
  "write the fix", "make the changes", "implement the fix for #1234".
hooks:
  PostToolUse:
    - matcher: "Edit|Write"
      hooks:
        - type: command
          command: "bash ${CLAUDE_SKILL_DIR}/scripts/post-edit.sh"
---

# Write

Make the minimal change that resolves the issue. Research has already
identified what to change — this skill handles the edit, formatting,
self-review, and commit.

## 1. Create a branch

```bash
git checkout -b fix/issue-<number>-<short-desc> main
```

Use a short kebab-case description derived from the issue title (3-5 words).

## 2. Read then edit

Always read each file before modifying it. Make the minimal change that
fixes the issue. Do not improve surrounding content, add comments, or
address adjacent problems.

Follow the writing guidelines in CLAUDE.md, STYLE.md, and COMPONENTS.md.

## 3. Front matter check

Every content page requires `title`, `description`, and `keywords` in its
front matter. If any are missing from a file you touch, add them.

## 4. Validate

Prettier runs automatically after each edit via the PostToolUse hook.
Run lint manually after all edits are complete:

```bash
scripts/lint.sh <changed-files>
```

The lint script runs markdownlint and vale on only the files you pass it,
so the output is scoped to your changes. Fix any errors it reports.

## 5. Self-review

Re-read each changed file: right file, right lines, change is complete,
front matter is present. Run `git diff` and verify only intended changes
are present.

## 6. Commit

Stage only the changed files:

```bash
git add <files>
git diff --cached --name-only  # verify — no package-lock.json or other noise
git commit -m "$(cat <<'EOF'
docs: <short description under 72 chars> (fixes #NNNN)

<What was wrong: one sentence citing the specific problem.>
<What was changed: one sentence describing the exact edit.>

Co-Authored-By: Claude <noreply@anthropic.com>
EOF
)"
```

The commit body is mandatory. A reviewer reading only the commit should
understand the problem and the fix without opening the issue.

## Notes

- Never edit `_vendor/` or `data/cli/` — these are vendored
- If a file doesn't exist, check for renames:
  `git log --all --full-history -- "**/filename.md"`
- If the fix requires a URL that cannot be verified, stop and report a
  blocker rather than guessing
