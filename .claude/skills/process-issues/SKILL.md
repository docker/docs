---
name: process-issues
description: >
  Process a batch of GitHub issues end-to-end: fetch unlabeled issues, triage
  them, apply labels, fix the actionable ones, and babysit the resulting PRs.
  Use this skill when the user wants to work through the issue backlog
  autonomously — e.g. "process issues", "work through the backlog", "run the
  issue pipeline". Accepts an optional batch size: "process issues 5".
  Pairs well with /loop for continuous background processing.
---

# Process Issues

Fetches a batch of unprocessed GitHub issues, triages each one, labels it,
fixes the actionable ones, and babysits the resulting PRs.

## Arguments

- Batch size (optional, default 10): number of issues to process per run.
  E.g. "process issues 5" or "process issues --batch 5".

## Labels

These labels track agent processing state. Create any that don't exist yet
with `gh label create --repo docker/docs <name> --color <hex>`.

| Label | Meaning |
|-------|---------|
| `agent/triaged` | Agent has analyzed this issue; verdict in a comment |
| `agent/fix` | Agent has opened a PR for this issue |
| `agent/skip` | Triaged; not actionable (STALE, UPSTREAM, INDETERMINATE, CLOSEABLE_FIXED) |

## Workflow

### 1. Resolve fork username and fetch a batch

Get the authenticated user's GitHub login — don't hardcode it:

```bash
FORK_USER=$(gh api user --jq '.login')
```

Fetch up to N open issues that have none of the `agent/*` labels:

```bash
gh issue list --repo docker/docs \
  --state open \
  --limit <N> \
  --json number,title,labels \
  --jq '[.[] | select(
    ([.labels[].name] | map(startswith("agent/")) | any) | not
  )]'
```

If there are no unprocessed issues, report that the backlog is clear and stop.

### 2. Triage each issue

Follow the full **triage-issues** skill workflow for all fetched issues,
running fetches in parallel. Produce a verdict for each:
OPEN, CLOSEABLE_FIXED, UPSTREAM, INDETERMINATE, or STALE.

### 3. Label each issue immediately after verdict

Apply `agent/triaged` to every issue regardless of verdict — it means "we
looked at it." Then apply a second label based on the outcome:

```bash
# All issues — mark as triaged
gh issue edit <number> --repo docker/docs --add-label "agent/triaged"

# Non-actionable (STALE, UPSTREAM, INDETERMINATE)
gh issue edit <number> --repo docker/docs --add-label "agent/skip"

# Already resolved (CLOSEABLE_FIXED) — close with explanation
gh issue close <number> --repo docker/docs \
  --comment "Closing — <one sentence explaining what resolved it>."

# Actionable (OPEN, no existing PR) — no extra label yet; agent/fix applied after PR is created
```

Leave a comment on every issue summarising the verdict and reasoning in one
sentence. Do this immediately — don't batch it for the end.

### 4. Fix actionable issues

For each issue with verdict OPEN and no existing open PR, follow the
**fix-issues** skill workflow.

Skip issues where:
- An open PR already exists
- Verdict is anything other than OPEN
- The fix requires changes to `_vendor/` or `data/cli/` (upstream owned)

After the PR is created, apply `agent/fix` to the issue:

```bash
gh issue edit <number> --repo docker/docs --add-label "agent/fix"
```

### 5. Babysit PRs

After opening PRs, schedule a recurring check with `/loop` so babysitting
continues asynchronously after the batch summary is reported:

```
/loop 5m babysit PRs <#N, #M, …> in docker/docs — check for failing checks,
new review comments, and requested changes; investigate and fix anything that
needs attention; stop looping once all PRs are merged or closed
```

At each check, for every open PR:

- **Failing checks**: investigate the failure, fix the cause, force-push an
  updated commit to the branch via the GitHub API
- **Review comments**: read them, address the feedback, push an update, reply
  to the comment
- **All clear**: note it and move on

Don't just report status — act on anything that needs attention.

### 6. Report results

```
## Batch summary

Processed: <N> issues
PRs opened: <n>
Skipped: <n> (STALE: n, UPSTREAM: n, INDETERMINATE: n)
Closed: <n> (already resolved)

### PRs opened
| Issue | PR | Checks | Review |
|-------|-----|--------|--------|
| #N    | #M  | ✅     | pending |

### Skipped
| Issue | Verdict | Reason |
|-------|---------|--------|
| #N    | STALE   | ...    |
```

## Notes

- **Fork username**: always resolve dynamically with `gh api user --jq '.login'`
- **One issue, one PR**: never combine multiple issues in a single branch
- **Validation**: skip `docker buildx bake lint vale` unless the change is
  complex — it's slow and the basic checks run automatically on the PR
- **Resumability**: labels are applied immediately at triage time, so if the
  session ends mid-run the next run skips already-processed issues automatically
