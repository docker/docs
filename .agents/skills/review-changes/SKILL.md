---
name: review-changes
description: >
  Review uncommitted or recently committed documentation changes for
  correctness, coherence, and style compliance. Use before creating a PR
  to catch issues. "review my changes", "review the diff", "check the fix
  before submitting", "does this look right".
context: fork
model: opus
---

# Review Changes

Evaluate whether the changes correctly and completely solve the stated
problem, without introducing new issues. Start with no assumptions — the
change may contain mistakes. Your job is to catch what the writer missed,
not to rubber-stamp the diff.

## 1. Identify what changed

Determine the scope of changes to review:

```bash
# Uncommitted changes
git diff --name-only

# Last commit
git diff --name-only HEAD~1

# Entire branch vs main
git diff --name-only main...HEAD
```

Pick the right comparison for what's being reviewed. If reviewing a branch,
use `main...HEAD` to see all changes since the branch diverged.

## 2. Read each changed file in full

Do not just read the diff. For every changed file, read the entire file to
understand the full context the change lives in. A diff can look correct in
isolation but contradict something earlier on the same page.

Then read the diff for the detailed changes:

```bash
# Adjust the comparison to match step 1
git diff --unified=10              # uncommitted
git diff --unified=10 HEAD~1       # last commit
git diff --unified=10 main...HEAD  # branch
```

## 3. Follow cross-references

For each changed file, check what links to it and what it links to:

- Search for other pages that reference the changed content (grep for the
  filename, heading anchors, or key phrases)
- Read linked pages to verify the change doesn't create contradictions
  across pages
- Check that anchor links in cross-references still match heading IDs

A change that's correct on its own page can break the story told by a
related page.

## 4. Verify factual accuracy

Don't assume the change is factually correct just because it reads well.

- If the change describes how a feature behaves, verify against upstream
  docs or source code
- If the change includes a URL, check that it resolves
- If the change references a CLI flag, option, or API field, confirm it
  exists

## 5. Evaluate as a reader

Consider someone landing on this page from a search result, with no prior
context:

- Does the page make sense on its own?
- Is the changed section clear without having read the issue or diff?
- Would a reader be confused by anything the change introduces or leaves
  out?

## 6. Review code and template changes

For non-Markdown changes (JS, HTML, CSS, Hugo templates):

- Trace through the common execution path
- Trace through at least one edge case (no stored preference, Alpine fails
  to load, first visit vs returning visitor)
- Ask whether the change could produce unexpected browser or runtime
  behavior that no automated tool would catch

## 7. Decision

**Approve** if the change is correct, coherent, complete, and factually
accurate.

**Request changes** if:
- The change does not correctly solve the stated problem
- There is a factual error or contradiction (on-page or cross-page)
- A cross-reference is broken or misleading
- A reader would be confused

When requesting changes, be specific: quote the exact text that is wrong,
explain why, and suggest the correct fix.
