---
name: review-pr
description: >
  Review one or more incoming Docker documentation pull requests as a
  maintainer. Independently validate technical claims, assess editorial fit
  and information architecture, choose a verdict, and draft exact inline or
  PR-wide feedback behind a confirmation gate. Use for requests such as
  "review PR 123", "is this PR correct?", "does this information belong
  here?", "validate this PR", or "help review backlog PRs". Do not use to
  maintain or fix a PR you own; use maintain-pr for that.
---

# Review PR

Review incoming contributions for factual correctness and whether they make
the documentation better as a whole. Treat a technically true addition as
insufficient when it is misplaced, overemphasized, redundant, or unhelpful
to the page's intended reader.

## Preserve the write boundary

Perform the review in two phases:

1. Research the PR, decide a verdict, and present the exact proposed
   comment or review text.
2. Wait for explicit user confirmation, then post only the confirmed text.

Before confirmation, do not post comments, submit a GitHub review, approve or
request changes, resolve threads, push commits, edit labels, or otherwise
mutate GitHub. A request to review or draft feedback is not confirmation to
post it. Ask `Post these comments?` and stop. Treat revisions to a draft as
unconfirmed until the user explicitly asks to post them.

## 1. Gather the full context

For each PR, inspect its metadata, body, commits, changed files, checks,
conversation, reviews, and linked issues. Always fetch inline comments
separately because `gh pr view --json reviews` omits them.

```bash
gh pr view <PR> --repo docker/docs \
  --json number,title,url,state,author,body,baseRefName,headRefName,headRefOid,commits,files,comments,reviews,reviewDecision,statusCheckRollup
gh api repos/docker/docs/pulls/<PR>/comments \
  --jq '[.[] | {id, author: .user.login, body, path, line, side, commit_id}]'
gh pr diff <PR> --repo docker/docs
```

Read linked issues and relevant discussion. An issue is evidence that a
reader was confused, but it does not establish the reporter's diagnosis or
justify a new highlighted note by itself. Green CI establishes only that automated
checks passed, not that the content is correct.

Fetch the PR head when local inspection is useful. Compare it with the
canonical upstream base rather than assuming the local branch is fresh.
Read each changed file in full, not only its diff.

## 2. Research independently

Verify every material claim against authoritative sources such as product
source code, upstream documentation, specifications, release notes, or safe
local reproduction. Do not accept the PR description, issue diagnosis, or
existing review feedback as fact.

Search the documentation for related explanations and canonical pages. Read
`STYLE.md`, `COMPONENTS.md`, and applicable repository instructions. Check
whether a changed file is generated or maintained upstream and identify the correct
upstream repository instead of proposing a local edit.

Distinguish among:

- a wrong fact
- a correct fact expressed inaccurately
- a correct fact placed on the wrong page
- content already explained elsewhere
- a real discovery problem better addressed with a short signpost and link
- a request that needs no documentation change.

If an external claim or replacement URL cannot be verified, report that
limitation instead of guessing.

## 3. Assess editorial fit

Apply these questions to each addition:

- Does it change a reader's decision or next action on this page?
- Is this the canonical page for the concept?
- Is the fact general, or specific to this page, feature, or component?
- Is the information already documented elsewhere?
- Would a concise local signpost to canonical coverage solve the discovery
  problem better than duplicating the explanation?
- Is the visual and textual weight proportional to the information's value?
- Does it preserve the page's scope, flow, and character?

Prefer one coherent explanation in the canonical location. Add local context
only when it helps the reader complete the task at hand. Avoid stray notes,
callouts, and exhaustive edge cases whose prominence exceeds their value.

## 4. Choose a decisive verdict

Lead with one of these outcomes:

- **Approve**: correct, useful, well placed, and ready to merge.
- **Approve with optional polish**: ready to merge; suggestions are genuinely
  non-blocking.
- **Focused rewrite**: the underlying need is valid, but wording, scope,
  placement, or structure should change before merge.
- **Close / no docs change**: incorrect, redundant, out of scope, or not a
  documentation problem.

Explain the verdict with evidence. When wording is the issue, provide exact
replacement text rather than a vague request to improve it.

## 5. Place feedback deliberately

Use an inline comment when the finding is anchored to a narrow changed line
or range and acting on it is local. Examples include an inaccurate sentence,
an ambiguous option description, a broken link, or a precise wording
replacement.

Use a PR-wide comment for scope, information architecture, overall approach,
multiple intertwined edits, or a proposed replacement section. Do not attach
holistic feedback to an arbitrary line.

Use both when appropriate: put the overall direction in the PR-wide comment
and line-specific corrections inline. Do not repeat the same point in both.
Consolidate related feedback so the author receives the fewest comments that
remain clear and actionable.

For every proposed inline comment, resolve and display the current changed
file path and right-side diff line. If the target line is not part of the
current diff or cannot be identified reliably, use a PR-wide comment that
quotes the target text instead. Never guess a line number.

End comments posted on the user's behalf with an accurate agent-disclosure
footer, such as `Generated by Codex`.

## 6. Present drafts and stop

Before any GitHub write, show the review in this form, omitting empty
sections:

```markdown
## Verdict

Focused rewrite

## Findings

- <finding and evidence>

## Proposed inline comments

1. `path/to/file.md:42`
   > Exact comment text

## Proposed PR-wide comment

> Exact comment text

Post these comments?
```

For multiple PRs, give each PR its own verdict and comment set. Make the
confirmation scope unambiguous. Do not interpret approval of one PR's drafts
as approval to post comments on the others.

## 7. Post only confirmed feedback

Immediately before posting, re-fetch the PR head SHA and diff. If either the
head or an inline target changed, stop and show the updated draft or
placement for confirmation.

Post confirmed inline comments as a single comment-only review when
practical. Use the current head SHA and right-side diff lines:

```bash
gh api repos/docker/docs/pulls/<PR>/reviews --method POST --input <payload>
```

The JSON payload contains `commit_id`, `event: "COMMENT"`, and a `comments`
array whose entries contain `path`, `line`, `side: "RIGHT"`, and `body`.
Submitting a review with `APPROVE` or `REQUEST_CHANGES` requires separate,
explicit user authorization; a verdict alone does not grant it.

Post confirmed holistic feedback separately:

```bash
gh pr comment <PR> --repo docker/docs --body-file <file>
```

Use a safely created temporary file or API input so Markdown, backticks, and
shell substitutions are preserved literally. Post exactly the confirmed
text. Verify the resulting review/comments and report their URLs and
placements. If GitHub rejects an inline location, do not silently fall back
to a PR-wide comment; report the failure and prepare a revised placement for
confirmation.

## Definition of done

- Verify technical claims with authoritative evidence.
- Evaluate usefulness, placement, duplication, and proportionality.
- Give a decisive verdict and exact actionable wording.
- Choose inline and PR-wide placement based on the feedback's scope.
- Show every exact draft and target before any GitHub mutation.
- Post only after explicit confirmation and verify what was posted.
