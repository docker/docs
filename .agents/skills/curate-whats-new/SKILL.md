---
name: curate-whats-new
description: Curate noteworthy Docker launches from documentation pull requests merged during a requested period. Use when generating or reviewing data/whats-new.json, preparing the Docker Docs What's new timeline, or deciding which documented releases merit Docker-wide highlights.
---

# Curate What's New

Treat merged documentation as evidence that a capability shipped, but do not
treat a documentation change as news by itself.

Include an item only when the merged documentation directly shows all of the
following:

- A released user-facing feature, material enhancement, or broader availability
  milestone that was not available before the period
- A substantial capability or workflow, not new syntax or a small control
  within an existing workflow
- Enough Docker-wide editorial significance to merit proactively telling users
  about it outside product release notes
- A useful published page and a factual title and description

Apply a high bar. The result is a curated launch archive, not a complete
changelog. A specialized feature can qualify when its user impact is
substantial. A quiet period can produce few or no items.

## Exclusions

Exclude documentation maintenance; fixes; rewrites; guidance for old behavior;
routine release or generated-content syncs; limitations, prerequisites, and
workarounds; narrow flags, settings, command variants, protocols, and
compatibility changes; incremental UI, safety, permissions, or observability
improvements; and lower-level Engine, Build, networking, or storage changes.
These qualify only when they are part of an independently newsworthy
product-level launch.

Judge the user outcome, not PR size, product popularity, labels, changed lines,
a dedicated page, or the existence of a new API or command.

## Select highlights

Include every qualifying launch; do not impose a quota. Mark the five most
important as `featured: true`, or all items when fewer than five qualify. Rank
by the magnitude and distinctness of the user outcome and the value of helping
its audience discover it. Breadth can matter, but a major capability for a
specialized audience can outrank a smaller change for a broad audience. Recency
and product variety are not ranking goals.

Create one item per launch and combine PRs that document the same launch.
Preserve existing copy while it remains accurate and qualifies. Change featured
status only when the relative importance of the candidate set changes.

## Procedure

1. Read `data/whats-new.json`.
2. List every PR merged in the requested period:

   ```console
   $ gh pr list --repo docker/docs --state merged --search 'merged:START..END' --limit 200
   ```

3. Inspect the diff and resulting pages for every plausible candidate.
4. Decide what qualifies using only evidence in the merged documentation.
5. Replace `period_start`, `period_end`, and `items` in
   `data/whats-new.json`. Sort items by `published` date, newest first.
6. Write `.pr-body.md` with the publication period, selected highlights and
   source PRs, plus concise reasons for plausible exclusions.

Each item must contain `product`, `title`, `description`, `url`, `published`,
`source_prs`, and `featured`. Use the canonical product name, a published
internal URL, the merge date in `YYYY-MM-DD` format, and source PR numbers.

Write factual, restrained copy. Avoid superlatives, promotional language, and
claims about ease or importance. Do not modify tracked files other than
`data/whats-new.json`.
