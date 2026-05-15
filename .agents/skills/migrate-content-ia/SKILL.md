---
name: migrate-content-ia
description: >
  Handle Hugo docs information-architecture moves: discover old vs new URLs,
  add front matter aliases (Phase 1), update in-repo links (Phase 2), interactive
  List 2 resolution and fragment validation (Phase 3; no guessing). Supports
  PR-scoped mapping plus whole-content sweeps for inbound links to that mapping,
  or a full-site follow-up. Triggers on: "IA migration", "redirects for moved
  pages", "fix links after content move", "PR-scoped link/anchor pass",
  "aliases for old URLs". After branch work, chain the review-changes skill
  (main...HEAD) before a PR. Agents must run the in-file required procedure
  and definition of done, not the phases alone in isolation.
---

# Migrate content IA (redirects + links + anchors)

Use this skill when pages **move or rename** under `content/` and you must
preserve old public URLs and/or fix cross-references. Work in **phases**;
choose **PR-scoped** vs **full-site** mode per run.

**Read first:** **CLAUDE.md** / **AGENTS.md** (URL rules, vendored areas, external
links, special cases) and **hugo.yaml** (`permalinks`, `refLinksErrorLevel`,
`disablePathToLower`). For **prose and link text**, follow **STYLE.md**; for
**components, front matter, and link examples**, follow **COMPONENTS.md**.

**Related skills:** **research** helps map moves and find inbound links; **write**
commits minimal edits. Run this skill’s phases after the move is identified (or
in parallel with research for large IA work).

## Agent: required procedure (do not skip)

**Common mistake (wrong):** use **`git diff main...HEAD` (or the PR’s file
list) as the full set of places to fix links** for a migration. That set shows
**what *moved***; it is **not** the list of every page that **points *to*** a
moved page. Inbound stragglers are often in files the PR **never** touched. You
must still **sweep the repo** for every string in the **old path and published-URL set**
for this run, not only for “files in the diff.”

**Definition of done (when the migration is *finished*):** **Both** of the
following (unless the user or **AGENTS.md** **explicitly defers** a **List 2**
item in **Phase 3**; document the deferral):

1. **`docker buildx bake validate`** passes for the branch, with no new
   build/link errors from this work.
2. A **sweep of the old path and published-URL set for this run** (see
   [Sweep commands](#sweep-commands) below) finds **no** remaining
   migration-relevant **inbound** reference—**including**:
   - links to an old **source** path (plain `.md` and equivalent `ref` forms),
   - links that use the old path **and** a `#fragment`,
   - and, where your mapping includes them, old **published-style** `link:` /
   `url:` / full-site URL strings,  
   **except** intentional entries to keep: for example `aliases` on the **new**
   canonical page, or **redirects.yml** *sources* you must not edit per policy.
   (A hit on a **source** that is only an `alias` line on the new page is
   **expected**—do not “fix” that away; distinguish alias rows from straggler
   links in body or nav config.)

**Chaining (policy):** when this branch’s content work is ready for handoff,
**run the [review-changes](../review-changes/SKILL.md) skill** on
**`main...HEAD`** (or **`merge-base`…`HEAD`** for a different target branch) so
the **whole branch** is re-read for cross-page issues before opening a PR. Do
not treat phases 0–3 alone as the final check.

**Run in order (mandatory for agents):**

1. **Scope the moves (mapping input):** set the Git range like **review-changes**
   (for a PR to `main`: `git diff --name-only main...HEAD`; for another target:
   `BASE=$(git merge-base <target-branch> HEAD)` then
   `git diff --name-only $BASE...HEAD`, as in **Phase 0.5**). Include
   renames; build the **old → new** table (source and published) per **Phase
   0**.
2. **Sweep and list:** for every **old** path/URL in that table, run
   [Sweep commands](#sweep-commands) on the **allowed** trees. Record
   every hit as **List 1** (no `#`) or **List 2** (old path with `#...`) per
   **Phase 0.5**.
3. **Phased edits:** **Phase 1** (`aliases`), then **Phase 2** (List 1), then
   **Phase 3** (List 2) with **no guessing**—as in the sections below.
4. **Re-sweep** the same old-path set, then run **`docker buildx bake
   validate`**. The **Definition of done** above is met or you have **explicit
   defers** for the remainder.
5. **review-changes:** run **[review-changes](../review-changes/SKILL.md)**
   on the branch vs **`main`…`HEAD`** (or the correct base) before a PR.

### Sweep commands

Use a **repository** search (e.g. `rg` / your IDE) so **nothing** in the
allowed scope is only eyeballed.

**Trees to include** (at minimum): all of `content/`, plus **`data/`** and
**`layouts/`** when a migration can appear in config, `link:`-like fields,
shortcodes, or hardcoded path strings. Follow **Vendored / generated** rules in
**AGENTS.md**; do not edit disallowed files.

**What to search for (repeat per row in the old side of the mapping):**

- **Hugo / source form:** path segments that identify the *old* file, e.g.
  `manuals/.../old-segment/...` or `../old-segment/.../page.md` as your tree
  uses; include variants that still appear in the repo.
- **Published / site form:** e.g. `/admin/.../old-slug/` in front matter, nav
  `url:`, or `https://docs.docker.com/...` in allowed files—**match the
  file’s** established pattern, per **Conventions** below.
- **Anchors:** search for the **old path string**; matches that also include
  `#...` belong on **List 2** for **Phase 3** unless the whole link is
  a pure path-only case.

[scripts/scope-pr-files.sh](scripts/scope-pr-files.sh) (if present) prints
**`PR_SCOPE_FILES` only**—it does **not** replace this sweep. Use it to build
the **old → new** table, **not** to list where inbound links were fixed.

## Progressive disclosure (optional)

The procedure below stays in this file. If a run produces a very large
**old → new** URL table, store that table in **`reference.md`** in this skill
directory and link it from the task summary, so the agent reads the long
mapping only when needed.

## Modes

- **PR-scoped (typical for a single PR)**  
  - **What the PR “owns” (focus):** use `git diff` / `base...HEAD` to know which
    pages and renames the branch actually moves (`PR_SCOPE_FILES`). The **old →
    new** mapping and **List 1 / List 2** for this migration are defined from
    **that** work, not from unrelated areas.
  - **Where to look for stale references (sweep):** search broadly—typically all
    of `content/` (and config, shortcodes, layouts, per Conventions)—for **inbound**
    links and fields whose **target** is an **old** path or URL in **this** PR’s
    mapping. Inbound stragglers are often in files the PR never touched; finding
    them is **in scope** for this migration.  
  - **What to edit:** update **any** file in the allowed trees that contains a
    **migration-relevant** reference (target ∈ this PR’s old path set) according
    to the phases below. **Do not** treat `PR_SCOPE_FILES` as a hard limit on
    *which files you may save* for **inbound** link repairs (unless
    project policy for a given PR says otherwise; then follow policy and
    **defer** out-of-PR file fixes).
  - **Out of scope (defer / ignore in this run):** link and anchor problems that
    are **not** about this PR’s old→new map—e.g. a different area’s own slug
    issues, rot unrelated to the remapped path set. *Example:* a PR that only
    remaps `content/strawberry/...` should not “fix the whole site”; it **should**
    still fix a link under `mango/…` that **points at** an old `strawberry/…` path
    in the mapping, and **should not** chase **mango/**-only issues that do
    not involve those old targets.

- **Full-site (complete migration after the PR)**  
  - Update stragglers **across the repo** (or all inbound links to moved
    sections), including config-driven `link:` fields if policy allows.  
  - Still make **minimal** edits; no drive-by rewrites to **unrelated** targets
    outside the run’s **declared** mapping and lists.

### No guessing

- The agent must **not** guess **replacement paths, published URLs, or fragment
  IDs** (including for consolidated pages, renamed headings, or
  “semantic” remaps of `#anchor` → new `#…`). If the user has not given an
  explicit new target, **ask**, **defer**, or **stop** per **AGENTS.md**; never
  infer, autocomplete, or substitute a plausible fragment from the target page’s
  heading list. That rule applies in **every** phase, including after validation
  in Phase 3.

---

## Conventions (links, anchors, redirects)

### Front matter `aliases` (redirects)

- Per **COMPONENTS.md**, `aliases` are **URLs that redirect to this page**.
- Add or **merge** on the **new canonical** page; do not drop unrelated
  entries. Match local examples: **published-style paths** (leading `/`), and
  **trailing `/`** when that matches existing pages in the same area.
- **No** speculative redirects for URLs that were never published.
- **Collision check** before adding: no other page or redirect may already
  own the same old path.
- If the site also uses **`data/redirects.yml`**, only add entries when
  project policy requires it; avoid duplicating the same old URL in
  `aliases` **and** `redirects.yml` unless maintainers do.

### Internal links in Markdown (STYLE.md + COMPONENTS.md)

- Use **relative paths to source files** (e.g. `../section/page.md`) with
  **`.md`**, following **COMPONENTS.md** examples, unless the file already
  uses an established pattern (e.g. some `link:` or nav fields use **published**
  paths without `manuals` or `.md` — **match the surrounding file**).
- Keep **CLAUDE.md** / **AGENTS.md** rules: internal ref targets under
  `content/manuals/...` often use the full **`/manuals/...`** path; published
  URLs omit the `manuals` segment—do not confuse the two when fixing links.
- **Link text (STYLE.md):** descriptive, ~**5 words**; no “click here” or
  “learn more”; **no** end punctuation **inside** the link text; **no** bold/italic
  on link text unless normal in the sentence.
- **Headings (STYLE):** **sentence case**; do not rename headings in passing
  unless the migration requires it (heading changes break fragments).

### Shortcodes and layouts (links not only in Markdown)

- **Phase 2–3 scope includes** any **shortcode or layout partial** (under
  **Modes**, search broadly for inbound links to the migration; **edits** follow
  the same file-level rules as for Markdown) that emits links: e.g. `ref` /
  `relref`, `link` fields in shortcode args, or hardcoded
  `docs.docker.com` / path strings. Grep for old paths, slugs, and fragments
  under `layouts/shortcodes/` (and `layouts/_default/` if partials build nav).
- Match each file’s existing pattern; do not rewrite working shortcode style
  just to “clean up.”

### Fragments / anchors (Phase 3)

- List 1 / List 2: fragment-bearing **cross-references to old paths** are tracked
  on **List 2** in Phase 0.5; do not bulk-rewrite them in the **List 1** pass
  (Phase 2). See Phase 0.5 and Phase 2.
- **Valid `#fragment` values:** after the user supplies a new fragment, it should
  match the **target** page’s **generated** heading ID (Hugo slugification; see
  **CLAUDE.md** / **AGENTS.md**). The agent still **validates** (see Phase 3) and
  must **not** “pick” a different id from the page to replace a bad answer—**No
  guessing**.
- Same-page: `[Text](#section-id)`.
- Cross-page: when user-provided, `#fragment` must still be checked against the
  **target** file. Validate fragments in shortcodes the same way as in body
  Markdown.

### External URLs (**AGENTS.md**)

- Do not commit **guessed** replacement URLs. If a URL cannot be verified,
  treat as blocked or drop the fragment per AGENTS guidance. See also **No
  guessing** above; internal and external link targets are treated the same for
  inference: **none** without user input or a verified source.

### Special cases (**AGENTS.md**)

- **Engine API version** pages: respect coordinated **`/latest/` `aliases`**
  rules—never leave two version files both owning `/latest/`.
- **Vendored / generated** trees: read-only; see CLAUDE.md. Do not “fix” links
  there if policy forbids.

---

## Phase 0 — Discovery (read-only; may use whole repo)

1. Read **hugo.yaml** (permalinks, `refLinksErrorLevel`, `disablePathToLower`).
2. From the branch (diff, renames), build a **mapping table**:
   - old source path → new source path  
   - old published URL → new published URL (from permalink rules)
3. **Case:** with `disablePathToLower: true`, filesystem path **case** appears in
   URLs—**directory and link casing must match** (e.g. `setup` vs `Setup`).
4. When planning **inbound link** fixes, treat old-path references as two
   categories: **no fragment** vs **with `#fragment`**. That split feeds
   **List 1** and **List 2** in Phase 0.5 and drives Phase 2 ordering (see
   there).

---

## Phase 0.5 — PR-scoped evaluation (required before edits in PR mode)

1. **Set `PR_SCOPE_FILES` (Git scope for PR mode)**  
   - When the PR **targets `main`**, use the same triple-dot form as
     **review-changes**:  
     `git diff --name-only main...HEAD`  
   - For a **different target branch** or a custom base, use the merge base:  
     `BASE=$(git merge-base <target-branch> HEAD)`  
     then:  
     `git diff --name-only "$BASE"...HEAD`  
   - Those paths define **what moved** in the branch; they are the primary input
     to the **old → new** path/URL table. They are **not** a hard cap on *where
     to search* for **inbound** links (see **Modes**): sweeps for links **to** old
     paths usually cover all of `content/` (and other trees per Conventions).  
   - If project policy **limits edits** to the diff for a given PR, follow that
     and **defer** link fixes in files outside the diff; note the exception in
     the task if the user relaxes that policy.

2. Build checklists (see **Modes** for sweep vs area-of-work):
   - path/URL mapping this run must honor (old source path → new; old published
     → new, from the **PR’s** moves in PR-scoped mode, or the **declared** full
     migration in full-site mode)
   - **List 1 — old path, no fragment:** every **inbound** reference, found on
     the **sweep** surface, to a moved **old** path that does **not** include a
     `#...` fragment (e.g. `…/banana.md` in the repo’s link style for that
     file).
   - **List 2 — old path with fragment:** every **inbound** reference, found on
     the same sweep, to a moved **old** path that **includes** a `#...` fragment
     (e.g. `…/banana.md#anchor` or the published-style equivalent in context). The
     **same** old path string may appear on **both** List 1 and List 2 for
     different links; duplication across the two lists is OK.
   - **Matching rules:** when recording List 1 / List 2, use **one** consistent
     path representation for comparison (e.g. relative `../path/banana.md` vs
     root-anchored) **per the conventions in this doc** and the **surrounding
     file’s** established pattern. Agents compare and skip List 2 links in the
     List 1 pass using the **same** representation rules.
3. **Out of scope** for the lists: only include references whose **old** target
   is in this run’s **mapping**. Do not build List 1/2 for unrelated **mango/**
   (or other) problems unless those links also target an **old** path that this
   migration renames. Defer those issues separately (see **Modes**).

---

## Phase 1 — `aliases` (old published URLs)

1. On each **new** canonical page, add or merge **`aliases`** for every **real**
   former public URL.
2. Do not strip existing unrelated aliases.
3. **PR-scoped:** add aliases only where the canonical file is in scope or the
   project requires it; otherwise list missing alias targets for follow-up.

---

## Phase 2 — In-repo link reference updates

1. **List 1 first (path only):** update references that belong to **List 1**
   (old path, **no** fragment). Replace old source paths or old published URLs
   with the **new** targets; preserve each file’s link pattern (relative vs
   root-anchored `.md` paths). **Do not** apply the same bulk path replacement to
   links that appear in **List 2** (old path **with** `#...`) during this
   sub-step—**leave** every **List 2** link **unchanged** for now.
2. **After List 1 is complete:** **re-scan** the **same** **sweep** surface as
   in Phase 0.5 (e.g. all of `content/` plus config) or **print** a clear list of
   all **remaining** **List 2** entries. Those links should still point at the
   **old** path and **old** fragment until Phase 3.  
3. **Full-site (extra sweep):** after steps 1–2, still use **AGENTS “Page
   deletion checklist”**-style thoroughness for **config / front matter**
   `link:` and similar so nav and grids are not left on old slugs. Apply the
   **List 1 / List 2** rules there too: path-only old references first; defer
   fragment-bearing rewrites in line with **List 2** until Phase 3.
4. **PR-scoped (which files to change):** apply List 1 and later Phase 3 updates
   to **every** file the **sweep** finds with a **migration-relevant** reference
   (inbound to an **old** path in the mapping), including files **not** in
   `PR_SCOPE_FILES`, per **Modes**. **Log** and **defer** (do not “fix”)
   unrelated stragglers. If policy forbids out-of-PR file edits, defer per step 1
   of Phase 0.5.  
5. Include **shortcodes and layout partials** (see Conventions and **Modes** for
   sweep vs focus).

---

## Phase 3 — List 2: interactive path and fragment resolution

**Prerequisites:** Phase 2 has updated **List 1**; **List 2** still lists **old
path + `#...`** (unchanged) for this migration. See **Modes** for which files
may be edited; **No guessing** applies.

1. **Print List 2** to the user: every remaining **old path** + `#anchor` (in the
   agreed representation), so nothing is hidden before the loop.
2. **For each distinct** `old-path#oldAnchor` (or process in the order the user
   prefers, one at a time):  
   - Ask: **What is the new path (and fragment, if any) for this content?** The
     user may give a new source path, published URL, and/or `#newAnchor` per
     project conventions.  
   - **Validate** the user’s answer: open the **target** page (or resolve the
     target) and check that `#newAnchor` (if any) **exists** as a real heading
     / generated id on that page, per **CLAUDE.md** / **AGENTS.md** (same rules
     as the rest of the site). **Do not** replace the user’s fragment with a
     “better” one from the file.  
   - If validation **fails** (unknown target file, or `#newAnchor` not found on
     the page): **warn** clearly (what failed: path vs missing fragment), then
     **ask again** for a corrected path and/or fragment. **Repeat** until
     validation passes or the user **defers** / **drops** the fragment (per
     **AGENTS.md**). **Never** guess a new fragment to fix the problem.  
   - When validation **passes:** update **all** in-repo references that match
     that **same** `old-path#oldAnchor` to the user-approved `new-path#newAnchor`
     (respect each file’s link style; include shortcodes/layouts on the same
     **sweep** surface as Phase 2).  
3. **Repeat** from step 1: **re-print** or **re-scan** for **List 2** until it is
   **empty** or the user defers the remainder.  
4. **PR-scoped / full-site:** the **loop** is the same. **Edits** follow **Modes**:
   migration-relevant **inbound** links may live in any file on the sweep; do
   not expand into **unrelated** link debt from other areas. Defer as in **Modes**
   and Phase 0.5.

---

## Optional: scripts helper

This skill includes a small **scope helper** so agents do not re-derive Git
recipes. See [scripts/scope-pr-files.sh](scripts/scope-pr-files.sh) — it prints
paths in PR scope for a given target branch (default `main`).

---

## Verification

```bash
docker buildx bake validate
```

Use the **Definition of done** in **Agent: required procedure (do not skip)**
as the final bar: **validate** must pass, and the **sweep** must be clean for
**plain** and **`#fragment`** old-path references, **or** the remainder must be
**explicitly deferred** in **Phase 3** per **AGENTS.md** / the user. Mid-run,
**Phase 2** may still leave **List 2** links unchanged **until** Phase 3; that
intermediate state is **not** the finished migration.
