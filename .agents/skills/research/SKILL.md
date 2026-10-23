---
name: research
description: >
  Research a documentation topic — locate affected files, understand the
  problem, identify what to change. Use when investigating an issue, a
  question, or a topic before writing a fix. Triggers on: "research issue
  1234", "investigate what needs changing for #500", "what files are
  affected by #200", "where is X documented", "is our docs page about Y
  accurate", "look into how we document Z".
---

# Research

Investigate the topic and determine whether documentation needs to change.
When an update is warranted, identify exact files, named targets within
those files, and the verified content needed.

## 1. Gather context

If the input is a GitHub issue number, fetch it:

```bash
gh issue view <number> --repo docker/docs \
  --json number,title,body,labels,comments
```

Otherwise, work from what was provided — a description, a URL, a question,
or prior conversation context. Identify the topic, affected feature, or
page to investigate.

## 2. Locate affected files

Search `content/` using the URL or topic from the issue. Remember the
`/manuals` prefix mapping when converting URLs to file paths.

For each candidate file, read the relevant section to confirm it contains
the reported problem.

## 3. Check vendored ownership

Before planning any edit, verify the file is editable locally:

- `_vendor/` — read-only, vendored via Hugo modules
- `data/cli/` — read-only, generated from upstream YAML
- `content/reference/cli/` — read-only, generated from `data/cli/`
- Everything else in `content/` — editable

If the fix requires upstream changes, identify the upstream repo and note
it as out of scope. See the vendored content table in CLAUDE.md.

## 4. Find related content

Look for pages that may need updating alongside the primary fix:

- Pages that link to the affected content
- Include files (`content/includes/`) referenced by the page
- Related pages in the same section describing the same feature

## 5. Verify facts

If the issue makes a factual claim about how a feature behaves, verify it.
Follow external links, read upstream source, check release notes. Do not
plan a fix based on an unverified claim.

If the fix requires a replacement URL and that URL cannot be verified (e.g.
network restrictions), report it as a blocker rather than guessing.

## 6. Assess documentation impact

A code change or release-note entry is evidence to investigate, not a reason
by itself to add documentation. Identify the reader's task and what the
change alters about it. Recommend an update when existing guidance becomes
incorrect, leaves a gap in that task, or omits a capability useful to readers.
An implementation fix may need no update if the docs already describe the
intended behavior.

For each recommended update, connect the reader's task, the documentation
gap, the destination page and section, and the proposed treatment. Prefer
correcting existing prose or examples when they serve that task. Check related
pages for contradictions, without assuming each needs the same explanation.

If there is no clear home, reconsider the value and scope before proposing a
section or page. A separate home can be justified by a distinct reader task,
not by the need to place a release-note bullet. "No documentation change" is
a valid finding; explain it briefly when assessing a set of changes.

## 7. Check the live site (if needed)

For URL or rendering issues, fetch the live page:

```text
https://docs.docker.com/<path>/
```

## 8. Report findings

Summarize what you found — files to change, the specific problem in each,
what the fix should be, and any constraints. This context feeds directly
into the write step.

Be specific: name the file, the section or element within it, and the
verified content needed. "Fix the broken link in networking.md" is not
specific enough. "In `compose/networking.md`, the 'Custom networks' section,
remove the note about `driver_opts` being ignored — this was fixed in
Compose 2.24" is.

## Notes

- Research quality bounds write quality. Vague research produces broad
  changes; precise research produces minimal ones.
- Do not create standalone research files — findings stay in conversation
  context for the write step.
