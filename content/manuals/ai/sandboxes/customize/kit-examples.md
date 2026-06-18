---
title: Kit examples
linkTitle: Examples
description: Copy-and-adapt spec.yaml snippets for common mixin and sandbox kit patterns — static files, install commands, background services, initFiles, Claude Code skills, and agent forks.
keywords: sandboxes, sbx, kits, mixins, examples, patterns, skills
weight: 25
---

{{< summary-bar feature_name="Docker Sandboxes sbx" >}}

> [!NOTE]
> Kits are experimental. The kit file format, CLI commands, and experience
> for creating, loading, and managing kits are subject to change as the
> feature evolves. Share feedback and bug reports in the
> [docker/sbx-releases](https://github.com/docker/sbx-releases) repository.

Each section below shows one `spec.yaml` snippet that demonstrates a
single kit pattern. These aren't complete, distributable kits — they're
small, focused examples you can lift into your own kit. For the full
spec reference, see [Kit spec reference](kit-reference.md).

## Drop a shared config file

Use static files under `files/workspace/` when the content is the same
across every sandbox and doesn't need any runtime values substituted
in. Typical use cases: linter rules, editor settings, a shared
`.editorconfig`, team dotfiles.

```text
ruff-lint/
├── spec.yaml
└── files/
    └── workspace/
        └── ruff.toml
```

```yaml {title="ruff-lint/spec.yaml"}
schemaVersion: "2"
kind: mixin
name: ruff-lint
displayName: Ruff
description: Python linting with shared team config

commands:
  install:
    - command: "uv tool install ruff@latest"
      user: "1000"
```

```toml {title="ruff-lint/files/workspace/ruff.toml"}
line-length = 80

[lint]
select = ["E", "F", "I"]
```

## Install a tool at sandbox creation

`commands.install` runs once per sandbox, at creation time. It's where
anything that needs to land in the image goes — package managers
(`apt-get`, `pip`, `npm`), binary downloads, or vendor install scripts.

```yaml
commands:
  install:
    - command: "apt-get update && apt-get install -y jq"
    - command: "curl -fsSL https://example.com/install.sh | sh"
```

Install commands run as root by default. Set `user: "1000"` when the
step should run as the agent user — for example, `npm install -g`
against a user-scoped prefix, or anything that writes to
`/home/agent/`.

## Install an internal CA certificate

If your organization uses a proxy that inspects HTTPS traffic, install
the proxy's internal root CA in the sandbox trust store. This helps
agents and SDKs trust certificates signed by the proxy.

```text
internal-ca/
├── spec.yaml
└── files/
    └── home/
        └── internal-ca.crt
```

Use a PEM-encoded certificate with a `.crt` extension. Files under
`files/home/` land in `/home/agent/` in the sandbox, so
`files/home/internal-ca.crt` becomes `/home/agent/internal-ca.crt` —
which is the path the install command reads from. If traffic can be
signed by more than one internal proxy, include each proxy's root CA in
the kit and install each certificate before running
`update-ca-certificates`.

```yaml {title="internal-ca/spec.yaml"}
schemaVersion: "2"
kind: mixin
name: internal-ca

commands:
  install:
    - command: "install -m 0644 /home/agent/internal-ca.crt /usr/local/share/ca-certificates/internal-ca.crt && update-ca-certificates"
      user: "0"
      description: Install internal CA certificate
```

`update-ca-certificates` adds the certificate to the system trust
store, so tools and SDKs that read the system bundle trust the proxy's
certificates without further configuration.

## Run a background service

`commands.startup` runs on every sandbox start. To keep a long-running
service such as a dev server or daemon alive, set `background: true`. The
sandbox runs the command in the background and replays startup commands on
each start, so the service comes back after a stop/start cycle:

```yaml
commands:
  startup:
    - command: ["my-service", "--port", "8080"]
      user: "1000"
      background: true
```

A background service doesn't write to your terminal. To capture its output
for debugging, wrap the command in a shell and redirect to a log file. Let
`background: true` run the command in the background rather than adding a
trailing `&` yourself:

```yaml
commands:
  startup:
    - command:
        - sh
        - -c
        - my-service --port 8080 > /tmp/my-service.log 2>&1
      user: "1000"
      background: true
```

An empty log file tells you the wrapper ran; a populated one tells you why
the service failed.

## Bake runtime values into a file with initFiles

When a config file needs a value that isn't known until sandbox start
— most often the absolute workspace path — use `commands.initFiles`.
The `${WORKDIR}` placeholder expands to the primary workspace path
when the file is written.

```yaml
commands:
  initFiles:
    - path: /home/agent/.local/bin/start-code-server.sh
      content: |
        exec code-server --bind-addr 0.0.0.0:8080 --auth none "${WORKDIR}"
      mode: "0755"
  startup:
    - command:
        - sh
        - -c
        - nohup /home/agent/.local/bin/start-code-server.sh > /tmp/code-server.log 2>&1 &
      user: "1000"
```

`mode: "0755"` makes the generated file executable so the startup
command can invoke it directly.

Use `initFiles` instead of a static file whenever the content depends
on a runtime value. Use a static file otherwise.

> [!TIP]
> This snippet is lifted from the
> [code-server kit](https://github.com/docker/sbx-kits-contrib/tree/main/code-server)
> in the contrib repository, which is also a runnable sample that demonstrates
> the full pattern.

## Ship a Claude Code skill

Claude Code reads project-scoped skills from
`.claude/skills/<name>/SKILL.md` in the workspace. Drop one into
`files/workspace/` and it's available in the sandbox.

```text
docker-review/
├── spec.yaml
└── files/
    └── workspace/
        └── .claude/
            └── skills/
                └── docker-review/
                    └── SKILL.md
```

```yaml {title="docker-review/spec.yaml"}
schemaVersion: "2"
kind: mixin
name: docker-review
displayName: Dockerfile review skill
description: Ships a Claude Code skill that reviews Dockerfiles
```

```markdown {title="docker-review/files/workspace/.claude/skills/docker-review/SKILL.md"}
---
name: docker-review
description: Review a Dockerfile for best practices. Use when the user asks to review, audit, or improve a Dockerfile.
---

When reviewing a Dockerfile, check:

1. Base image — pinned tag or digest, appropriate for the workload
2. Layer order — dependencies copied before application source
3. Image size — multi-stage builds, `.dockerignore`, package-manager cache flags
4. Security — non-root `USER`, no secrets in `ARG`/`ENV`
5. Reproducibility — pinned package versions, frontend directive where relevant
```

Kits have to target the workspace rather than `~/.claude/` because
sandboxes don't pick up user-level agent configuration from the host.
See the
[FAQ](../faq.md#why-doesnt-the-sandbox-use-my-user-level-agent-configuration)
for details.

## Override agent settings

Sandboxes seed settings files for some built-in agents during setup.
For example, the sandbox writes `/home/agent/.claude/settings.json`
for the `claude` agent. This happens after the kit's static files and
`initFiles`, so a kit can't override those paths with either mechanism.
Use `commands.startup` instead, which runs after the sandbox seeds its
files:

```yaml
commands:
  startup:
    - command:
        - sh
        - -c
        - |
          mkdir -p /home/agent/.claude
          cat > /home/agent/.claude/settings.json <<'JSON'
          {"permissions": {"allow": ["Bash(echo:*)"]}}
          JSON
      user: "1000"
      description: Write user-scope claude settings
```

Startup commands replay on every sandbox start, so the script must be
idempotent. The heredoc pattern overwrites cleanly each time.

## Fork an existing agent

Sandbox kits (`kind: sandbox`) define a full agent from scratch. The most
common variant is a fork of a built-in agent — same image and
credentials, but a different entrypoint. This example reproduces the
built-in `claude` agent but drops `--dangerously-skip-permissions` so
every tool call prompts for approval:

```yaml {title="claude-safe/spec.yaml"}
schemaVersion: "2"
kind: sandbox
name: claude-safe
displayName: Claude Code (with approval prompts)
description: Claude Code without --dangerously-skip-permissions

sandbox:
  image: "docker/sandbox-templates:claude-code-docker"
  aiFilename: CLAUDE.md
  entrypoint:
    run: [claude]

caps:
  network:
    allow:
      - "claude.com:443"

credentials:
  - service: anthropic
    apiKey:
      name: ANTHROPIC_API_KEY
      inject:
        - domain: api.anthropic.com
          header: x-api-key
          format: "%s"
        - domain: console.anthropic.com
          header: x-api-key
          format: "%s"
```

Launch with the kit's `name:` as the agent argument to `sbx run`:

```console
$ sbx run claude-safe --kit ./claude-safe
```

For a step-by-step walkthrough of building a new sandbox kit from
scratch, see [Build an agent](build-an-agent.md).

## More examples

These patterns are all drawn from working kits in the
[sbx-kits-contrib](https://github.com/docker/sbx-kits-contrib)
repository, which contains each example as a complete, loadable kit.
Use it to study the full shape of a kit, or load one directly:

```console
$ sbx run claude --kit "git+https://github.com/docker/sbx-kits-contrib.git#dir=<kit>"
```
