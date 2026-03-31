---
title: Custom environments
weight: 80
description: Customize agent sandbox environments or use the shell sandbox for manual setup.
---

{{< summary-bar feature_name="Docker Sandboxes sbx" >}}

Every sandbox is customizable — agents install packages, pull images, and
configure tools as they work, and those changes persist for the sandbox's
lifetime. This page covers two things beyond that: a general-purpose shell
sandbox for manual setup, and custom templates that capture a configured
environment into a reusable image so you don't have to set it up again every
time.

## Shell sandbox

`sbx run shell` drops you into a Bash login shell inside a sandbox
with no pre-installed agent binary. It's useful for installing and
configuring agents manually, testing custom implementations, or inspecting a
running environment.

```console
$ sbx run shell ~/my-project
```

The workspace path defaults to the current directory. To run a one-off command
instead of an interactive shell, pass it after `--`:

```console
$ sbx run shell -- -c "echo 'Hello from sandbox'"
```

Set your API keys as environment variables so the sandbox proxy can inject them
into API requests automatically. Credentials are never stored inside the VM:

```console
$ export ANTHROPIC_API_KEY=sk-ant-xxxxx
$ export OPENAI_API_KEY=sk-xxxxx
```

Once inside the shell, you can install agents using their standard methods,
for example `npm install -g @continuedev/cli`. For complex setups, build a
[custom template](#custom-templates) instead of installing interactively each
time.

The shell sandbox uses the `shell` base image — the common base environment
without a pre-installed agent.

## Custom templates

Custom templates are reusable sandbox images that extend one of the built-in
agent environments with additional tools and configuration baked in. Instead of
asking the agent to install packages every time, build a template once and
reuse it across sandboxes and team members.

Templates make sense when multiple people need the same environment, when setup
involves steps that are tedious to repeat, or when you need pinned versions of
specific tools. For one-off work, the default image is fine — ask the agent to
install what's needed.

> [!NOTE]
> Custom templates customize an existing agent's environment — they don't
> create new agent runtimes. The agent that launches inside the sandbox is
> determined by the base image variant you extend and the agent you specify
> in the `sbx run` command, not by binaries installed in the template.

### Base images

All sandbox templates are published as
`docker/sandbox-templates:<variant>`. They are based on Ubuntu and run as a
non-root `agent` user with sudo access. Most variants include Git, Docker
CLI, and common development tools like Node.js, Python, Go, and Java.

| Variant               | Agent                                                                |
| --------------------- | -------------------------------------------------------------------- |
| `claude-code`         | [Claude Code](https://claude.ai/download)                            |
| `claude-code-minimal` | Claude Code with a minimal toolset (no Node.js, Python, Go, or Java) |
| `codex`               | [OpenAI Codex](https://github.com/openai/codex)                      |
| `copilot`             | [GitHub Copilot](https://github.com/github/copilot-cli)              |
| `docker-agent`        | [Docker Agent](https://github.com/docker/docker-agent)               |
| `gemini`              | [Gemini CLI](https://github.com/google-gemini/gemini-cli)            |
| `kiro`                | [Kiro](https://kiro.dev)                                             |
| `opencode`            | [OpenCode](https://opencode.ai)                                      |
| `shell`               | No agent pre-installed. Use for manual agent setup.                  |

Each variant also has a `-docker` version (for example,
`claude-code-docker`) that includes a full Docker Engine running inside the
sandbox — no local Docker daemon required. The `-docker` variants are the
defaults used by `sbx run` on macOS and Linux. These variants run in
privileged mode inside the microVM (not on your host), with a dedicated block
volume at `/var/lib/docker`, and `dockerd` starts automatically inside the
sandbox.

The block volume defaults to 50 GB and uses a sparse file, so it only
consumes disk space as Docker writes to it. On Windows, the volume is not
sparse and the full 50 GB is allocated at creation time, which increases
startup time. For this reason, the non-docker variants are the default on
Windows.

To override the volume size, set the `DOCKER_SANDBOXES_DOCKER_SIZE`
environment variable to a size string before starting the sandbox:

```console
$ DOCKER_SANDBOXES_DOCKER_SIZE=10g sbx run claude
```

Use the non-docker variant if you don't need to build or run containers
inside the sandbox and want a lighter, non-privileged environment. Specify
it explicitly with `--template`:

```console
$ sbx run claude --template docker.io/docker/sandbox-templates:claude-code
```

### Build a custom template

Building a custom template requires [Docker Desktop](https://docs.docker.com/desktop/).

Write a Dockerfile that extends one of the base images. Pick the variant that
matches the agent you plan to run. For example, extend `claude-code` to
customize a Claude Code environment, or `codex` to customize an OpenAI Codex
environment.

The following example creates a Claude Code template with Rust and
protocol buffer tools pre-installed:

```dockerfile
FROM docker/sandbox-templates:claude-code
USER root
RUN apt-get update && apt-get install -y protobuf-compiler
USER agent
RUN curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh -s -- -y
```

Use `root` for system-level package installations (`apt-get`), and switch back
to `agent` before installing user-level tools. Tools that install into the
home directory, such as `rustup`, `nvm`, or `pyenv`, must run as `agent` —
otherwise they install under `/root/` and aren't available in the sandbox.

Build the image and push it to an OCI registry, such as Docker Hub:

```console
$ docker build -t my-org/my-template:v1 --push .
```

> [!NOTE]
> The Docker daemon used by Docker Sandboxes pulls templates from a registry
> directly; it doesn't share the image store of your local Docker daemon on
> the host.

Unless you use the permissive `allow-all` network policy, you may also need to
allow-list any domains that your custom tools depend on:

```console
$ sbx policy allow network "*.example.com:443,example.com:443"
```

Then run a sandbox with your template. The agent you specify must match the
base image variant your template extends:

```console
$ sbx run --template docker.io/my-org/my-template:v1 claude
```

Because this template extends the `claude-code` base image, you run it with
`claude`. If you extend `codex`, use `codex`; if you extend `shell`, use
`shell` (which drops you into a Bash shell with no agent).

> [!NOTE]
> Unlike Docker commands, `sbx` does not automatically resolve the Docker Hub
> domain (`docker.io`) in image references.

### Template caching

Template images are cached locally. The first use pulls from the registry;
subsequent sandboxes reuse the cache. Cached images persist across sandbox
creation and deletion, and are cleared when you run `sbx reset`.
