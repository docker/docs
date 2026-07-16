---
name: triage-issue
description: >
  Analyze a single GitHub issue for docker/docs — check whether the problem
  still exists, determine a verdict, and report findings. Use when asked to
  triage, assess, or review an issue, even if the user doesn't say "triage"
  explicitly: "triage issue 1234", "is issue 500 still valid", "should we
  close #200", "look at this issue", "what's going on with #200".
argument-hint: "<issue-number>"
context: fork
---

# Triage Issue

Given GitHub issue **$ARGUMENTS** from docker/docs, figure out whether
it's still a real problem and say what should happen next.

## 1. Fetch the issue

```bash
gh issue view $ARGUMENTS --repo docker/docs \
  --json number,title,body,state,labels,createdAt,updatedAt,closedAt,assignees,author,comments
```

## 2. Understand the problem

Read the issue body and all comments. Identify:

- What is the reported problem?
- What content, URL, or file does it reference?
- Has anyone already proposed a fix or workaround in the comments?

Check for linked PRs in the issue timeline, not only in the issue body or
comments:

```bash
gh api repos/docker/docs/issues/$ARGUMENTS/timeline --paginate \
  --jq '.[] | select(.event=="cross-referenced" or .event=="connected" or .event=="referenced") | {event, created_at, source: .source.issue.html_url, title: .source.issue.title, state: .source.issue.state}'
```

If an open PR already addresses the issue, don't open another PR. Review the
existing PR instead, and report that the issue already has an associated PR. A
merged PR is strong evidence the issue is fixed. A closed-without-merge PR means
the issue is likely still open.

## 3. Follow URLs

Find all `docs.docker.com` URLs in the issue body and comments. For each:

- Fetch the URL to check if it still exists (404 = content removed or moved)
- Check whether the content still contains the problem described
- Note when the page was last updated relative to when the issue was filed

For non-docs URLs (GitHub links, external references), fetch them too if
they are central to understanding the issue.

## 4. Check the repository

If the issue references specific files, content sections, or code:

- Find and read the current version of that content
- Check whether the problem has been fixed, content moved, or file removed
- Remember the `/manuals` prefix mapping when looking up files

## 5. Check for upstream ownership

If the issue is about content in `_vendor/` or `data/cli/`, it cannot be
fixed here. Identify which upstream repo owns it (see the vendored content
table in CLAUDE.md).

## 6. Decide and act

After investigating, pick one of these verdicts and take the corresponding
action on the issue:

- **Close it** — the problem is already fixed, the content no longer exists,
  or the issue is too outdated to be useful. Close the issue with a comment
  explaining why:

  ```bash
  gh issue close $ARGUMENTS --repo docker/docs \
    --comment "Closing: <one-sentence reason>"
  ```

- **Fix it** — the problem is real and fixable in this repo. Name the
  file(s) and what needs to change. Label the issue `status/confirmed` and
  remove `status/triage` if present:

  ```bash
  gh api repos/docker/docs/issues/$ARGUMENTS/labels \
    --method POST --field 'labels[]=status/confirmed'
  gh api repos/docker/docs/issues/$ARGUMENTS/labels/status%2Ftriage \
    --method DELETE || true
  ```

- **Escalate upstream** — the problem is real but lives in vendored content.
  Name the upstream repo. Label the issue `status/upstream` and remove
  `status/triage` if present:

  ```bash
  gh api repos/docker/docs/issues/$ARGUMENTS/labels \
    --method POST --field 'labels[]=status/upstream'
  gh api repos/docker/docs/issues/$ARGUMENTS/labels/status%2Ftriage \
    --method DELETE || true
  ```

- **Leave it open** — you can't determine the current state, or the issue
  needs human judgment. Label the issue `status/needs-analysis`:

  ```bash
  gh api repos/docker/docs/issues/$ARGUMENTS/labels \
    --method POST --field 'labels[]=status/needs-analysis'
  gh api repos/docker/docs/issues/$ARGUMENTS/labels/status%2Ftriage \
    --method DELETE || true
  ```

Don't overthink the classification. An old issue isn't stale if the problem
still exists. An upstream issue is still valid — it's just not fixable here.

Also apply the most relevant `area/` label based on the content affected.
Available area labels: `area/accounts`, `area/admin`, `area/ai`,
`area/api`, `area/billing`, `area/build`, `area/build-cloud`, `area/cli`,
`area/compose`, `area/compose-spec`, `area/config`, `area/contrib`,
`area/copilot`, `area/desktop`, `area/dhi`, `area/engine`,
`area/enterprise`, `area/extensions`, `area/get-started`, `area/guides`,
`area/hub`, `area/install`, `area/networking`, `area/offload`,
`area/release-notes`, `area/samples`, `area/scout`, `area/security`,
`area/storage`, `area/subscription`, `area/swarm`, `area/ux`. Pick one
(or at most two if the issue clearly spans areas). Skip if none fit.

```bash
gh api repos/docker/docs/issues/$ARGUMENTS/labels \
  --method POST --field 'labels[]=area/<name>'
```

## 7. Report

Write a short summary: what the issue reports, what you found, and what
should happen next. Reference the specific files, URLs, or PRs that support
your conclusion. Skip metadata fields — the issue itself has the dates and
labels. Mention the action you took (closed, labeled, etc.).

## Notes

- Always check timeline cross-references before deciding to fix an issue
- Do not narrate your process — produce the final report
- End every issue comment with an accurate agent-disclosure footer that names
  the active coding agent, for example `Generated by Codex` or `Generated by
Claude Code`.
.
