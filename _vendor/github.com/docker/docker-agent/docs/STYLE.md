# Documentation style guide

A short reference for writers and reviewers. Goal: keep voice, naming
and examples consistent across every page on this site.

## Product naming

| Context | Use | Don't use |
|---|---|---|
| Prose, headings, marketing | **Docker Agent** (two words, both capitalised — the proper name of the product) | docker-agent, Docker-Agent, docker agent (in prose) |
| The CLI command | `docker agent` (lower-case, two words, in monospace) | `docker-agent`, `Docker Agent run` |
| The repository / module path | `docker/docker-agent` | docker/Docker-Agent |
| Internal identifiers / package names | as defined in code (e.g. `cagent`) — never invent new spellings in prose | mixing internal identifiers into user-facing copy |

A simple rule of thumb:

- **Talking about the product?** → "Docker Agent"
- **Showing a command the user types?** → `docker agent run agent.yaml`

## Voice

- Address the reader as **you**, not "we" or "the user".
- Prefer present tense and active voice ("the agent reads files",
  not "files will be read by the agent").
- Keep sentences short. Two short sentences usually beat one compound
  one.
- Avoid "simply", "just", "easily" — they're rarely accurate and
  often condescending.

## Code samples

- All shell prompts use a dollar sign followed by a space (`$`) and the
  command on the same line. Output, when shown, has no prompt.
- YAML/HCL examples should be runnable as-is when reasonable, or end
  in `# ...` to make truncation explicit.
- The canonical example agent uses `model: anthropic/claude-sonnet-4-5`.
  Use a different model only when the example is *about* that model.
- File names in prose are in `monospace` (`agent.yaml`, not "agent.yaml").

## Callouts

Callouts are written as portable GitHub-style alerts so the same
Markdown renders correctly on docs.docker.com (Hugo), GitHub, and this
site (a small script upgrades them to the styled panels):

```markdown
> [!TIP]
> **When to use it**
>
> Body text.
```

- `[!NOTE]` — neutral context
- `[!TIP]` — positive, "consider this"
- `[!IMPORTANT]` — must-read to succeed
- `[!WARNING]` — caution, breaking, security

The bold line directly after the marker is an optional title; omit it
when the default label (Note, Tip, …) is enough. Don't prefix the
title with an emoji — the icon badge already provides one.

## Links

Internal links are plain relative Markdown paths to the target file,
including the `index.md` filename:

```markdown
See the [Quick Start](../../getting-started/quickstart/index.md).
```

Both Jekyll (`jekyll-relative-links`) and Hugo (docs.docker.com link
render hook) resolve these to the right URL. Never use Liquid
(`relative_url`) or absolute `/path/` links in `docs/**` content —
they break when the page is mounted on docs.docker.com.

## Availability badges

When a page documents a feature that is merged on `main` but not yet
in a tagged release, mark it so readers of the stable docs know what
to expect:

```markdown
> [!NOTE]
> **Coming in v1.99**
>
> This feature is available on `main` and ships in v1.99.
```

Remove the badge in the release PR that tags the version (the
CHANGELOG update is a good reminder to sweep for `Coming in` markers).

## Glossary one-liners

When a page first introduces a term, link to its concept page or use
one of these standard one-liners:

- **Agent** — an LLM with instructions, tools, and (optionally)
  sub-agents, defined in YAML or HCL.
- **Toolset** — a group of related tools the agent can call (e.g.
  `filesystem`, `shell`, `mcp`).
- **MCP** — Model Context Protocol, an open standard for tool servers.
- **A2A** — Agent-to-Agent protocol, used to talk to other agents
  over HTTP.
- **TUI** — Terminal User Interface, the default interactive front end
  Docker Agent ships with.
- **OCI** — Open Container Initiative; the same registry format used
  for Docker images. Docker Agent reuses it for sharing agents.
