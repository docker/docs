---
title: Install Docker Skills
linkTitle: Install
description: Install Docker's official skills in compatible coding agents, Docker Sandboxes, or Docker Agent, and verify that your agent can use them.
keywords: [docker skills, install skills, agent skills, docker sandboxes, docker agent]
weight: 10
---

Choose the installation method managed by your agent or organization. Docker's
official skill collection evolves, so available skills and product coverage may
change. Native plugins, the Gemini CLI extension, and the cross-client skills
CLI are separate options; the skills CLI is not a prerequisite. For the current
collection, see the [repository catalog](https://github.com/docker/skills#readme).

## Claude Code {#claude-code}

In Claude Code, add Docker's marketplace, then install its plugin:

```text
/plugin marketplace add docker/skills
/plugin install docker-skills@docker
```

Use Claude Code's plugin manager to check the installation scope and receive
updates. If the plugin isn't found, confirm that the marketplace was added and
that your organization permits it. See [Claude Code's plugin
docs](https://code.claude.com/docs/en/discover-plugins).

<!-- vale off -->
## GitHub Copilot CLI {#github-copilot-cli}
<!-- vale on -->

Add Docker's marketplace in the GitHub Copilot CLI, then install its plugin:

```console
$ copilot plugin marketplace add docker/skills
$ copilot plugin install docker-skills@docker
```

Use Copilot CLI's plugin manager for scope and updates. If the plugin isn't
found, confirm that the marketplace was added and that your organization
permits it. See [Copilot CLI's plugin
documentation](https://docs.github.com/en/copilot/how-tos/copilot-cli/customize-copilot/plugins-finding-installing).

## Cursor {#cursor}

In Cursor's plugin or marketplace interface, add `docker/skills` if your
organization enables repository-backed marketplaces. Manage scope and updates
in that interface. If it isn't available or a skill isn't discovered, check
that the plugin is enabled or use the [skills CLI](#skills-cli) instead. See
[Cursor's plugin docs](https://cursor.com/docs/plugins).

## Codex {#codex}

Where the Codex client offers repository-backed plugins, select the Docker
marketplace and the `docker-skills` plugin in its marketplace interface. Manage
scope and updates in the client. If the interface isn't available or the
plugin isn't found, use the [skills CLI](#skills-cli) instead.

## Gemini CLI {#gemini-cli}

Install the repository as a Gemini CLI extension:

```console
$ gemini extensions install https://github.com/docker/skills
```

Manage its scope and updates in Gemini CLI; the repository URL follows its
upstream revision, not an immutable release. If the extension or its skills
don't appear, check that it is enabled and restart the client. See the
[Gemini CLI extension docs](https://geminicli.com/docs/extensions/).

## Google Antigravity {#google-antigravity}

Antigravity reads standard `SKILL.md` directories from your workspace's
`.agents/skills/`. From the project directory, use the [skills CLI](#skills-cli)
and select Antigravity and the skills you need:

```console
$ npx skills add docker/skills --agent antigravity
```

Alternatively, [clone a reviewed Docker Skills release](#git-clone-or-manual-copy)
and copy each selected directory from its `skills/` folder into
`<workspace-root>/.agents/skills/<skill-folder>/`, keeping its `SKILL.md` and
supporting files together. Start a new conversation and ask for a matching
Docker task to verify discovery. See Google's [Antigravity agent skills
docs](https://antigravity.google/docs/skills) for other discovery locations.
The Gemini CLI extension is a separate installation method, not an Antigravity
plugin.

## Any agent with the skills CLI {#skills-cli}

From your project directory, run the interactive installer:

```console
$ npx skills add docker/skills
```

Select the agent and skills when prompted. The default installation is
project-scoped; use `--global` for user scope. For CLI-managed installations,
run `npx skills list` to check the selection and `npx skills update` to update
it (add `--global` for the user scope). If the agent can't find a skill,
check the selected agent and scope; if it doesn't follow the installed link,
reinstall with `--copy`. See the [skills CLI docs](https://skills.sh/docs).

## Docker Sandboxes {#docker-sandboxes}

> [!NOTE]
> The `sbx skills` installer is experimental and may change or be removed.
> Check `sbx skills --help` in your installed release before automating it.

Install from the repository into the shared sandbox skill store:

```console
$ sbx skills add docker/skills
```

Run `sbx skills ls` to check the store and `sbx skills update` to refresh
repository-installed skills. Start or restart a sandbox to check that the skill
is linked. If `sbx skills` is unavailable, use the [skills CLI](#skills-cli)
or a [manual copy](#git-clone-or-manual-copy). Learn more about
[shared agent skills](/ai/sandboxes/workflows/agent-skills/).

## Docker Agent {#docker-agent}

Docker Agent consumes skills installed in its documented
[discovery paths](/ai/docker-agent/features/skills/); it doesn't install or
update them. Install with the [skills CLI](#skills-cli) or copy the selected
skills to a [documented discovery path](#git-clone-or-manual-copy), then select
those skills in `agent.yaml`:

```yaml
agents:
  root:
    model: dmr/ai/qwen3
    instruction: Help with Docker development tasks.
    skills: true
    toolsets:
      - type: filesystem
```

`skills: true` includes all discovered skills; use a list of skill names to
limit the selection. Update with the installer that owns the files and start a
new `docker agent run ./agent.yaml` after changes. If a skill isn't available,
check the discovery path and the `agents.<name>.skills` filter.

## Git clone or manual copy {#git-clone-or-manual-copy}

When a managed installer doesn't fit your client, clone a reviewed
[Docker Skills release](https://github.com/docker/skills/releases), replacing
`vX.Y.Z` with the tag you selected:

```console
$ git clone --branch vX.Y.Z --depth 1 https://github.com/docker/skills.git
```

Copy each selected directory under `skills/` into a discovery path documented
by your agent. Keep its `SKILL.md` and all supporting files together. The
destination determines project or user scope; update by replacing the complete
skill directory from another reviewed tag. If the skill isn't found, check the
client's discovery path and restart it. If supporting files are missing, recopy
the entire directory.

## Verify the installation

Check your client's installed-plugin or extension view, `npx skills list`, or
`sbx skills ls`, as appropriate. For a manual copy, check that the destination
contains a readable `SKILL.md` and its supporting files. Start a new agent
session and ask for a matching Docker task, such as “Review my Dockerfile for
cache efficiency and non-root execution.” Confirm that the agent uses the
expected guidance before relying on it.

## Update and pin

Update native plugins and extensions with their client, CLI-managed skills with
`npx skills update` in the matching scope, and shared sandbox skills with
`sbx skills update`. Docker Agent relies on the installer that owns its
discovered files. For a fixed snapshot, choose a tag from
[Docker Skills releases](https://github.com/docker/skills/releases) and use a
tagged source rather than a rolling repository URL. Move to another reviewed
tag deliberately; manually copied directories don't update themselves.

## Troubleshooting

- If a plugin or extension doesn't appear, check that it is enabled in the
  intended scope and that your organization permits the repository.
- If the agent doesn't find an installed skill, check its discovery path and
  selected agent and scope, then start a new session. For Docker Agent, also
  check the `skills:` filter; for Sandboxes, confirm the shared store is mounted.
- If a link isn't followed or support files are missing, reinstall with a
  physical copy of the complete skill directory.
- If updates don't appear, use the installer that owns the installation rather
  than a different client or CLI; manually copied or pinned sources require a
  deliberate replacement.
