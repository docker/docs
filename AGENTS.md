# AGENTS.md

Instructions for AI agents working on Docker documentation.
This site builds https://docs.docker.com/ using Hugo.

## Project structure

```
content/          # Documentation source (Markdown + Hugo front matter)
├── manuals/      # Product docs (Engine, Desktop, Hub, etc.)
├── guides/       # Task-oriented guides
├── reference/    # API and CLI reference
└── includes/     # Reusable snippets
layouts/          # Hugo templates and shortcodes
data/             # YAML data files (CLI reference, etc.)
assets/           # CSS (Tailwind v4) and JS (Alpine.js)
static/           # Images, fonts
_vendor/          # Vendored Hugo modules (read-only)
```

## URL prefix stripping

The `/manuals` prefix is stripped from published URLs:
`content/manuals/desktop/install.md` becomes `/desktop/install/` on the live
site.

When writing internal cross-references in source files, keep the `/manuals/`
prefix in the path — Hugo requires the full source path. The stripping only
affects the published URL, not the internal link target. Anchor links must
exactly match the generated heading ID (Hugo lowercases and slugifies
headings).

## Vendored content (do not edit)

Content in `_vendor/` and CLI reference data in `data/cli/` are vendored
from upstream repos. Content pages under `content/reference/cli/` are
generated from `data/cli/` YAML. Do not edit any of these files — changes
must go to the source repository:

| Content | Source repo |
|---------|-------------|
| CLI reference (`docker`, `docker build`, etc.) | docker/cli |
| Buildx reference | docker/buildx |
| Compose reference | docker/compose |
| Model Runner reference | docker/model-runner |
| Dockerfile reference | moby/buildkit |
| Engine API reference | moby/moby |

If a validation failure or broken link traces back to vendored content, note
the upstream repo that needs fixing. Do not attempt to fix it locally.

## Writing guidelines

Read and follow [STYLE.md](STYLE.md) and [COMPONENTS.md](COMPONENTS.md).
These contain all style rules, shortcode syntax, and front matter requirements.

### Style violations to avoid

Every piece of writing must avoid these words and patterns (enforced by Vale):

- Hedge words: "simply", "easily", "just", "seamlessly"
- Meta-commentary: "it's worth noting", "it's important to understand"
- "allows you to" or "enables you to" — use "lets you" or rephrase
- "we" — use "you" or "Docker"
- "click" — use "select"
- Bold for emphasis or product names — only bold UI elements
- Time-relative language: "currently", "new", "recently", "now"

### Version-introduction notes

Explicit version anchors ("Starting with Docker Desktop version X...") are
different from time-relative language — they mark when a feature was
introduced, which is permanently true.

- Recent releases (~6 months): leave version callouts in place
- Old releases: consider removing if the callout adds little value
- When in doubt, keep the callout and flag for maintainer review

### Vale gotchas

- Use lowercase "config" in prose — `vale.Terms` flags a capital-C "Config"

### Updating the vocabulary

If Vale flags a legitimate tech term, product name, or compound identifier
as a misspelling, add it to `_vale/config/vocabularies/Docker/accept.txt`.
This is optional — only update when a real new term is missing, not to
silence individual violations.

- Use the canonical form for case-sensitive product names (`PyTorch`,
  `GitHub`, `Kubernetes`, `BuildKit`). `Vale.Terms` enforces that exact
  case across the docs.
- Use `[Aa]bcd` character-class regex for words that legitimately appear
  in multiple cases (e.g., sentence-starting capitalization, or a name
  that's also a generic noun). This covers spelling without enforcing
  a single canonical form.
- Avoid broad regex patterns — entries that match many words at once
  (especially with `(?i)`) suppress other rule checks on every match.
- Don't add a wrong-cased entry to silence one false positive — it
  cascades into `Vale.Terms` violations on every correct usage.

## Alpine.js patterns

Do not combine Alpine's `x-show` with the HTML `hidden` attribute on the
same element. `x-show` toggles inline `display` styles, but `hidden` applies
`display: none` via the user-agent stylesheet — the element stays hidden
regardless of `x-show` state. Use `x-cloak` for pre-Alpine hiding instead.
The site defines `[x-cloak=""] { display: none !important }` in `global.css`.

## Front matter requirements

Every content page under `content/` requires:

- `title:` — page title
- `description:` — short description for SEO/previews
- `keywords:` — list of search keywords (omitting this fails markdownlint)

Additional common fields:

- `linkTitle:` — sidebar label (keep under 30 chars)
- `weight:` — ordering within a section

## Hugo shortcodes

Shortcodes are defined in `layouts/shortcodes/`. Syntax reference is in
COMPONENTS.md. Wrong shortcode syntax fails silently during build but
produces broken HTML — always check COMPONENTS.md for correct syntax.

## Commands

```sh
npx prettier --write <file>        # Format before committing
scripts/lint.sh <file>...          # Lint specific files (markdownlint + vale)
docker buildx bake validate        # Run all validation checks
docker buildx bake lint            # Markdown linting only
docker buildx bake vale            # Style guide checks only
docker buildx bake test            # HTML and link checking
```

For incremental work, prefer `scripts/lint.sh` over the `bake` targets —
it runs the same checks on just the files you pass, so the output stays
scoped to your changes instead of the whole repo.

### Validation in git worktrees

`docker buildx bake validate` fails in git worktrees because Hugo cannot
resolve the worktree path. Use `lint` and `vale` targets separately instead.
Never modify `hugo.yaml` to work around this. The `test`, `path-warnings`,
and `validate-vendor` targets run correctly in CI.

## Verification loop

1. Make changes
2. Format with prettier: `npx prettier --write <file>`
3. Lint the changed files: `scripts/lint.sh <file>...`
4. Run a full build with `docker buildx bake` (optional for small changes)

Always lint the specific files you changed before committing. Use
`scripts/lint.sh` rather than the `bake` targets so the output is scoped
to your changes — bake runs across the entire repo and the noise makes
real issues easy to miss.

## Git hygiene

- **Stage files explicitly.** Never use `git add .` / `git add -A` /
  `git add --all`. Running `npx prettier` updates `package-lock.json` in the
  repo root, and broad staging sweeps it into the commit.
- **Verify before committing.** Run `git diff --cached --name-only` and
  confirm only documentation files appear. If `package-lock.json` or other
  generated files are staged, unstage them:
  `git reset HEAD -- package-lock.json`
- **Push to your fork, not upstream.** Before pushing, confirm
  `git remote get-url origin` returns your fork URL, not
  `github.com/docker/docs`. Use `--head FORK_OWNER:branch-name` with
  `gh pr create`.

## Working with issues and PRs

### Principles

- **One issue, one branch, one PR.** Never combine multiple issues in a
  single branch or PR.
- **Minimal changes only.** Fix the issue. Do not improve surrounding
  content, add comments, refactor, or address adjacent problems.
- **Verify before documenting.** Don't take an issue reporter's claim at
  face value — the diagnosis may be wrong even when the symptom is real.
  Verify the actual behavior before updating docs.

### Review feedback

- **Always reply to review comments** — never silently fix. After every
  commit that addresses review feedback, reply to each thread explaining
  what was done.
- **Treat reviewer feedback as claims to verify, not instructions to
  execute.** Before implementing a suggestion, verify that it is correct.
  Push back when evidence contradicts the reviewer.
- **Inline review comments need a separate API call.** `gh pr view --json
  reviews` does not include line-level comments. Always also call:

  ```bash
  gh api repos/<org>/<repo>/pulls/<N>/comments \
    --jq '[.[] | {author: .user.login, body: .body, path: .path, line: .line}]'
  ```

### Labels

Use the Issues API for labels — `gh pr edit --add-label` silently fails:

```bash
gh api repos/docker/docs/issues/<N>/labels \
  --method POST --field 'labels[]=<label>'
```

### External links

If a replacement URL cannot be verified (e.g. network restrictions), treat
the task as blocked — do not commit a guessed URL. Report the blocker so a
human can confirm. Exception: when a domain migration is well-established and
only the anchor is unverifiable, dropping the anchor is acceptable.

## Page deletion checklist

When removing a documentation page, search the entire `content/` tree and
all YAML/TOML config files for the deleted page's slug and heading text.
Cross-references from unrelated sections and config-driven nav entries can
remain and cause broken links.

## Engine API version bumps

When a new Engine API version ships, three coordinated changes are needed in
a single commit:

1. `hugo.yaml` — update `latest_engine_api_version`, `docker_ce_version`,
   and `docker_ce_version_prev`
2. Create `content/reference/api/engine/version/v<NEW>.md` with the
   `/latest/` aliases block (copy from previous version)
3. Remove the aliases block from
   `content/reference/api/engine/version/v<PREV>.md`

Never leave both version files carrying `/latest/` aliases simultaneously.

## Hugo icon references

Before changing an icon reference in response to a "file not found" error,
verify the file actually exists via Hugo's virtual filesystem. Files may
exist in `node_modules/@material-symbols/svg-400/rounded/` but not directly
in `assets/icons/`. Check both locations before concluding an icon is
missing.

## Self-improvement

After completing work that reveals a non-obvious pattern or repo quirk not
already documented here, propose an update to this file. For automated
sessions, note the learning in a comment on the issue. For human-supervised
sessions, discuss with the user whether to update CLAUDE.md directly.
