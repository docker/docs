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

Thoroughly investigate the topic at hand and produce a clear plan for
the fix. The goal is to identify exact files, named targets within those
files, and the verified content needed for the fix.

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

## 6. Check the live site (if needed)

For URL or rendering issues, fetch the live page:

```
https://docs.docker.com/<path>/
```

## 7. Report findings

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
