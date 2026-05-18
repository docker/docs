---
name: create-lab-guide
description: "Clone a dockersamples Labspace repo, extract learning objectives and module structure from labspace.yaml, and produce a Hugo guide page under content/guides/ with correct frontmatter, labspace-launch shortcode, and Docker docs style compliance. Use when asked to create a lab guide, write a Labspace page, add a Docker lab tutorial, migrate a lab to docs, or document a hands-on lab."
---

# Create Lab Guide

Create a guide page for a Docker Labspace: clone the source repo, extract
structure from `labspace.yaml`, write the Hugo markdown page, and validate.

## Inputs

- **REPO_NAME**: GitHub repo in the `dockersamples` org (e.g. `labspace-ai-fundamentals`)

## Step 1: Clone the labspace repo

```bash
TMPDIR=$(mktemp -d)
git clone --depth 1 https://github.com/dockersamples/{REPO_NAME}.git "$TMPDIR/{REPO_NAME}"
```

## Step 2: Extract key information

Read these files from the cloned repo:

| File | Purpose |
|------|---------|
| `README.md` | Lab purpose and overview |
| `labspace/labspace.yaml` | Module structure and content paths |
| `labspace/*.md` | Module content (only files listed in `labspace.yaml`) |
| `.github/workflows/*.yml` | Published Compose file URL for the launch command |
| `compose.override.yaml` | Check for top-level `model` specs (triggers `model-download` param) |

Extract:
1. A short description for the `description` and `summary` frontmatter fields.
2. Learning objectives from the module content.
3. Whether a model download is required (`compose.override.yaml` → top-level `model` key).

## Step 3: Write the guide markdown

Place the file at `content/guides/lab-{GUIDE_ID}.md`.

```markdown
---
title: "Lab: { Short title }"
linkTitle: "Lab: { Short title }"
description: |
  A short description of the lab for SEO and social sharing.
summary: |
  A short summary of the lab for the guides listing page. 2-3 lines.
keywords: AI, Docker, Model Runner, agentic apps, lab, labspace
aliases: # Include only for AI-related labs
  - /labs/docker-for-ai/{REPO_NAME_WITHOUT_LABSPACE_PREFIX}/
params:
  tags: [ai, labs]
  time: 20 minutes
  resource_links:
    - title: A resource link pointing to relevant documentation or code
      url: /ai/model-runner/
    - title: Labspace repository
      url: https://github.com/dockersamples/{REPO_NAME}
---

Short explanation of the lab and what it covers.

## Launch the lab

{{< labspace-launch image="dockersamples/{REPO_NAME}" >}}

## What you'll learn

By the end of this Labspace, you will have completed the following:

- Objective #1
- Objective #2
- Objective #3

## Modules

| # | Module | Description |
|---|--------|-------------|
| 1 | Module #1 | Description of module #1 |
| 2 | Module #2 | Description of module #2 |
| 3 | Module #3 | Description of module #3 |
```

Conditional rules:
- All lab guides **must** include `labs` in `params.tags`.
- AI-related labs: also add `ai` tag and an alias under `/labs/docker-for-ai/`.
- If a model download is required: add `model-download: true` to the `labspace-launch` shortcode.

## Step 4: Apply Docker docs style rules

Follow STYLE.md and COMPONENTS.md. Key rules:

| Avoid | Use instead |
|-------|-------------|
| "we", "let's" | Imperative voice or "you" |
| "simply", "easily", "just" | Remove the hedge word |
| "allows you to" / "enables you to" | "lets you" or rephrase |
| "click" | "select" |
| Bold for emphasis / product names | Bold only for UI elements |
| "currently", "new", "recently" | Remove time-relative language |

Use `console` as the language hint for shell blocks with `$` prompts.
Use contractions ("it's", "you're", "don't").

## Step 5: Validate

1. Confirm frontmatter has `title`, `description`, `keywords`, and `params.tags` including `labs`.
2. Run `npx prettier --write <file>` to format.
3. Run `docker buildx bake lint vale` and fix any errors.
4. Re-read the file and verify: correct shortcode syntax, objectives match source content, modules match `labspace.yaml`, no vendored paths edited.

Do not proceed to commit until validation passes.

