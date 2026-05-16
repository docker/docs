---
title: FAQ
weight: 70
description: Frequently asked questions about Docker Sandboxes.
keywords: docker sandboxes, sbx, faq, sign in, telemetry
---

{{< summary-bar feature_name="Docker Sandboxes sbx" >}}

## Why do I need to sign in?

Docker Sandboxes is built around the idea that you and your agents are a team.
Signing in gives each sandbox a verified identity, which lets Docker:

- **Tie sandboxes to a real person.** Governance matters when agents can build
  containers, install packages, and push code. Your Docker identity is the
  anchor.
- **Enable team features.** Team-scale features like
  [organization governance](security/governance.md), shared environments, and
  audit logs need a concept of "who," and adding that later would be worse for
  everyone.
- **Authenticate against Docker infrastructure.** Sandboxes pull images, run
  daemons, and talk to Docker services. A Docker account makes that seamless.

Your Docker account email is only used for authentication, not marketing.

## Can I enforce sandbox policies across my organization?

Yes. Admins can centrally manage network and filesystem policies from the
Docker Admin Console. Rules defined there apply to every sandbox in the
organization and take precedence over local rules set with `sbx policy`.
Admins can optionally delegate specific rule types back to local control so
developers can add additional allow rules.

See [Organization governance](security/governance.md). This feature requires
a separate paid subscription —
[contact Docker Sales](https://www.docker.com/products/ai-governance/#contact-sales)
to get started.

## Does the CLI collect telemetry?

The `sbx` CLI collects basic usage data about CLI invocations:

- Which command you ran
- Whether it succeeded or failed
- How long it took
- If you're signed in, your Docker username is included

Docker Sandboxes doesn't monitor sessions, read your prompts, or access your
code. Your code stays in the sandbox and on your host.

To opt out of all analytics, set the `SBX_NO_TELEMETRY` environment variable:

```console
$ export SBX_NO_TELEMETRY=1
```

## How do I set custom environment variables inside a sandbox?

The [`sbx secret`](/reference/cli/sbx/secret/) command only supports a fixed set
of [services](security/credentials.md#built-in-services) (Anthropic, OpenAI,
GitHub, and others). If your agent needs an environment variable that isn't
tied to a supported service, such as `BRAVE_API_KEY` or a custom internal
token, write it to `/etc/sandbox-persistent.sh` inside the sandbox. This
file is sourced on every shell login, so the variable persists across agent
sessions for the sandbox's lifetime.

Use `sbx exec` to append the export:

```console
$ sbx exec -d <sandbox-name> bash -c "echo 'export BRAVE_API_KEY=your_key' >> /etc/sandbox-persistent.sh"
```

The `bash -c` wrapper is required so the `>>` redirect runs inside the
sandbox instead of on your host.

> [!NOTE]
> Unlike `sbx secret`, which injects credentials through a host-side proxy
> without exposing them to the agent, this approach stores the value inside
> the sandbox. The agent process can read it directly. Only use this for
> credentials where proxy-based injection isn't available.

Variables in `/etc/sandbox-persistent.sh` are sourced automatically when
bash runs inside the sandbox, including interactive sessions and agents
started with `sbx run`. If you run a command directly with
`sbx exec <name> <command>`, the command runs without a shell, so the
persistent environment file is not sourced. Wrap the command in `bash -c`
to load the environment:

```console
$ sbx exec <sandbox-name> bash -c "your-command"
```

To verify the variable is set, open a shell in the sandbox:

```console
$ sbx exec -it <sandbox-name> bash
$ echo $BRAVE_API_KEY
```

## Why do agents run without approval prompts?

The sandbox itself is the safety boundary. Because agents run inside an
isolated microVM with [network policies](security/policy.md),
[credential isolation](security/credentials.md), and no access to your host
system outside the workspace, the usual reasons for approval prompts (preventing
destructive commands, network access, file modifications) are handled by the
sandbox isolation layers instead.

If you prefer to re-enable approval prompts, change the permission mode
inside the session. Most agents let you switch permission modes after
startup. In Claude Code, use the `/permissions` command to change the mode
interactively.

To make approval prompts the default for every session, define a custom
agent kit that overrides the agent's entrypoint to drop the
permission-skipping flag. For example, a kit that launches Claude Code
without `--dangerously-skip-permissions`:

```yaml {title="claude-safe/spec.yaml"}
schemaVersion: "1"
kind: agent
name: claude-safe
agent:
  image: "docker/sandbox-templates:claude-code-docker"
  entrypoint:
    run: [claude]
```

Run it with `sbx run claude-safe --kit ./claude-safe/`. See
[Agent kits](customize/kits.md#agent-kits) for the full pattern.

## How do I know if my agent is running in a sandbox?

Ask the agent. The agent can see whether or not it's running inside a sandbox.
In Claude Code, use the `/btw` slash command to ask without interrupting an
in-progress task:

```text
/btw are you running in a sandbox?
```

## Why doesn't the sandbox use my user-level agent configuration?

Sandboxes don't pick up user-level agent configuration from your host. This
includes directories like `~/.claude` for Claude Code or `~/.codex` for Codex,
where hooks, skills, and other settings are stored. Only project-level
configuration in the working directory is available inside the sandbox.

To make configuration available in a sandbox, copy or move what you need into
your project directory before starting a session:

```console
$ cp -r ~/.claude/skills .claude/skills
```

Don't use symlinks — a sandboxed agent can't follow symlinks to paths outside
the sandbox.

Collocating skills and other agent configuration with the project itself is a
good practice regardless of sandboxes. It's versioned alongside the code and
evolves with the project as it changes.
