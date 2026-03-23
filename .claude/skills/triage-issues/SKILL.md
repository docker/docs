---
name: triage-issues
description: >
  Triage one or more GitHub issues for the docker/docs repository. Analyzes each
  issue's content, checks whether the problem still exists in the repo and on the
  live site, and produces a structured verdict. Use this skill whenever the user
  asks to triage, analyze, review, or assess GitHub issues — e.g. "triage issue
  1234", "what's the status of these issues", "which of these can be closed",
  "look at issues 100 200 300 and tell me what to do with them".
---

# Triage Issues

Given one or more GitHub issue numbers from the docker/docs repository, analyze
each issue and produce a structured verdict on its current status.

## Workflow

### 1. Fetch all issues in parallel

For each issue number, fetch everything in a single call:

```bash
gh issue view <number> --repo docker/docs \
  --json number,title,body,state,labels,createdAt,updatedAt,closedAt,assignees,author,comments
```

When triaging multiple issues, fetch all of them in parallel before starting
analysis — don't process one at a time.

### 2. Analyze each issue

For each issue, work through these checks:

#### a. Understand the problem

Read the issue body and all comments. Identify:
- What is the reported problem?
- What content, URL, or file does it reference?
- Are there linked PRs? Check whether they were merged or closed without merge.
- Has anyone already proposed a fix or workaround in the comments?

#### b. Follow URLs

Find all `docs.docker.com` URLs in the issue body and comments. For each one:
- Fetch the URL to check if it still exists (404 = content removed or moved)
- Check whether the content still contains the problem described
- Note when the page was last updated relative to when the issue was filed

For non-docs URLs (GitHub links, external references), fetch them too if they
are central to understanding the issue.

#### c. Check the repository

If the issue references specific files, content sections, or code:
- Use file tools to find and read the current version of that content
- Check whether the problem has been fixed, the content moved, or the file removed
- Remember the `/manuals` prefix mapping when looking up files

#### d. Check for upstream ownership

If the issue is about content in `_vendor/` or `data/cli/`, it cannot be fixed
here. Identify which upstream repo owns it and note that in your analysis.

### 3. Determine verdict

Assign one of these verdicts based on what you found:

| Verdict | When to use |
|---------|-------------|
| **OPEN** | Issue is valid and still unfixed in this repo |
| **CLOSEABLE_FIXED** | Content has been updated, corrected, or removed since the issue was filed |
| **UPSTREAM** | Problem exists but originates in vendored/upstream content |
| **INDETERMINATE** | Not enough information to determine current state |
| **STALE** | Outdated with no recent activity; references content or features that no longer exist; context has changed enough that a new issue would be more appropriate |

Be confident when evidence is clear. Use INDETERMINATE only when you genuinely
cannot determine the current state after checking.

### 4. Report results

#### Single issue

Print a structured report:

```
## Issue #<number>: <title>

**Verdict:** <VERDICT>
**Confidence:** <high|medium|low>
**Filed:** <creation date>
**Last activity:** <last comment or update date>
**Labels:** <labels or "none">

### Summary
<One or two sentences describing the reported problem.>

### Analysis
<What you checked and what you found. Reference specific URLs, files, or
content. Note any linked PRs or related issues.>

### Recommendation
<Concrete next step: close with a comment, fix specific content, escalate
to upstream repo, request more info from reporter, etc.>
```

#### Multiple issues

Start with a summary table, then print the full report for each issue:

```
| Issue | Title | Verdict | Confidence |
|-------|-------|---------|------------|
| #123  | ...   | OPEN    | high       |
| #456  | ...   | STALE   | medium     |

---

## Issue #123: ...
[full report]

---

## Issue #456: ...
[full report]
```

## Notes

- A merged PR linked to an issue is strong evidence the issue is fixed
- A closed-without-merge PR means the issue is likely still open
- Check creation date and last activity to assess staleness — issues with no
  activity in over a year that reference old product versions are candidates
  for STALE
- Do not narrate your process; just produce the final report(s)
