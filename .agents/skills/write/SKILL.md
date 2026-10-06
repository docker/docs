---
name: write
description: >
  Write or edit reader-facing technical prose for immediate comprehension.
  Use for documentation, PR titles and descriptions, release notes, design
  documents, user-facing explanations, and substantive comments. Also use
  when asked to make writing clearer, more natural, less AI-generated, or
  easier to read. For documentation fixes, handles edits, formatting,
  self-review, and commits after research identifies what to change.
  Do not use for code-only tasks with no prose deliverable.
hooks:
  PostToolUse:
    - matcher: "Edit|Write"
      hooks:
        - type: command
          command: "bash ${CLAUDE_SKILL_DIR}/scripts/post-edit.sh"
---

# Write

Follow the human-readable prose guidance in AGENTS.md and the style rules in
STYLE.md. Apply the prose review below to every prose deliverable, including
titles and short comments.

For a prose-only request, return the requested text after reviewing it. Use
the numbered workflow for a documentation fix after research has identified
what to change. Branch and commit steps apply to that workflow, not to drafting
a title, description, or comment. Publishing prose requires authorization
from the user.

## Prose review

Before returning or committing prose, read it as if explaining the change to
an experienced colleague. Fix sentences that require a second reading even
when their grammar is correct. Lint passing does not replace this review.

- Replace invented labels and stacks of nouns with the action or relationship
  they describe: "post-update configuration validation" becomes "validate
  the configuration after the update"
- Use direct verbs: "perform an evaluation of" becomes "evaluate"
- State concrete behavior instead of vague claims: "improve navigation
  discoverability" becomes "show the current page in the sidebar" when that
  is the actual change
- Keep established technical terms, exact identifiers, and qualifications
  needed for accuracy. Clarify the surrounding sentence instead of replacing
  a precise term with a vague one
- Check titles and headings separately. Use enough words to make their meaning
  clear without relying on the body to explain an invented label

Make these edits silently. Return the requested prose without a report of
this review unless the user asks for one.

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

rumdl runs automatically after each edit via the PostToolUse hook.
Run lint manually after all edits are complete:

```bash
scripts/lint.sh <changed-markdown-files>
```

The lint script runs rumdl and Vale on only the files you pass it,
so the output is scoped to your changes. Fix errors and warnings on lines you
added or changed, and review each suggestion. Vale can exit successfully when
warnings or suggestions remain in its output.

## 5. Self-review

Re-read each changed file: right file, right lines, change is complete,
front matter is present. Run `git diff` and verify only intended changes
are present. Apply the prose review to the changed text and the commit message.

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
