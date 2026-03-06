# AGENTS.md

Instructions for AI agents working on the Docker documentation
repository. This site builds https://docs.docker.com/ using Hugo.

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

The `/manuals` prefix is stripped from published URLs:
`content/manuals/desktop/` → `/desktop/` on the live site.

## Writing guidelines

Read and follow [STYLE.md](STYLE.md) and [COMPONENTS.md](COMPONENTS.md).
These contain all style rules, shortcode syntax, and front matter
requirements.

## Vendored content (do not edit)

Content in `_vendor/` and CLI reference pages generated from
`data/cli/` are vendored from upstream repos. Don't edit these
files — changes must go to the source repository:

- docker/cli, docker/buildx, docker/compose, docker/model-runner → CLI reference YAML in `data/cli/`
- moby/buildkit → Dockerfile reference in `_vendor/`
- moby/moby → Engine API docs in `_vendor/`

If a validation failure traces back to vendored content, note the
upstream repo that needs fixing but don't block on it.

## Commands

```sh
npx prettier --write <file>        # Format before committing
docker buildx bake validate        # Run all validation checks
docker buildx bake lint            # Markdown linting only
docker buildx bake vale            # Style guide checks only
docker buildx bake test            # HTML and link checking
```

## Verification loop

1. Make changes
2. Format with prettier
3. Run `docker buildx bake lint vale`
4. Run a full build with `docker buildx bake`

## Self-improvement

After every correction or mistake, update this file with a rule to
prevent repeating it. End corrections with: "Now update AGENTS.md so
you don't make that mistake again."

## Mistakes to avoid

- Don't use hedge words: "simply", "easily", "just", "seamlessly"
- Don't use meta-commentary: "it's worth noting that...", "it's important to understand that..."
- Don't use "allows you to" or "enables you to" — use "lets you" or rephrase
- Don't use "we" — use "you" or "Docker"
- Don't use "click" — use "select"
- Don't bold product names or for emphasis — only bold UI elements
- Don't use time-relative language: "currently", "new", "recently", "now"
- Don't edit vendored content in `_vendor/` or `data/cli/`
