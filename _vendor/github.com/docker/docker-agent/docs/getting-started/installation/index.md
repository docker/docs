---
title: "Installation"
description: "Get docker-agent running on your system in minutes."
keywords: docker agent, ai agents, getting started, installation
weight: 20
canonical: https://docs.docker.com/ai/docker-agent/getting-started/installation/
---

_Get docker-agent running on your system in minutes._

## Prerequisites

- An API key for at least one AI provider (OpenAI, Anthropic, Google, etc.)
- **Optional:** [Docker Desktop](https://www.docker.com/products/docker-desktop/) — for running containerized MCP tools and Docker Model Runner

## Docker Desktop (Pre-installed)

Starting with [Docker Desktop 4.63](https://docs.docker.com/desktop/release-notes/#4630), **docker-agent is already available**. No separate installation needed — just open a terminal and run:

```bash
$ docker agent version
```

> [!TIP]
> Docker Desktop bundles docker-agent and keeps it up to date. This is the easiest way to get started, especially if you want to use Docker MCP tools and Docker Model Runner.

## Homebrew (macOS / Linux)

Install docker-agent using [Homebrew](https://brew.sh/):

```bash
# Install
$ brew install docker-agent

# Verify
$ docker-agent version
```

You can also install docker-agent as a docker CLI plugin, by copying `docker-agent` binary in `~/.docker/cli-plugins`. You can then run `docker agent version`.

## Download Binary Releases

Download [prebuilt binary releases](https://github.com/docker/docker-agent/releases) for Windows, macOS, and Linux from the GitHub Releases page.

### macOS / Linux

```bash
# Download the latest release
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m); case "$ARCH" in x86_64) ARCH=amd64;; aarch64) ARCH=arm64;; esac
curl -L "https://github.com/docker/docker-agent/releases/latest/download/docker-agent-${OS}-${ARCH}" -o docker-agent
chmod +x docker-agent
sudo mv docker-agent /usr/local/bin/
docker-agent version

# or alternatively, instead of moving to /usr/local/bin:
mkdir -p ~/.docker/cli-plugins
sudo mv docker-agent ~/.docker/cli-plugins
docker agent version
```

### Windows

Download `docker-agent-windows-amd64.exe` from the [releases page](https://github.com/docker/docker-agent/releases), rename it to `docker-agent.exe` and add it to your PATH. Alternatively you can move it to `~/.docker/cli-plugins`

## Optional Self-Updates

When docker-agent is installed from a standalone GitHub release binary, you can opt in to automatic self-updates by setting `DOCKER_AGENT_AUTO_UPDATE` to a truthy value (`1`, `true`, `yes`, or `on`):

```bash
# Enable for one command
DOCKER_AGENT_AUTO_UPDATE=1 docker agent run

# Or enable for the current shell session
export DOCKER_AGENT_AUTO_UPDATE=1
docker agent run
```

With self-updates enabled, docker-agent checks the latest GitHub release before normal commands run. If a newer release exists and your session is interactive, docker-agent asks whether you want to install it or keep running your current version. When the answer is yes (or the session is non-interactive, such as CI or piped input, in which case the update proceeds automatically), it downloads the asset for your OS and architecture, verifies the release-provided SHA-256 digest/checksum, replaces the current binary, and restarts the command with the same arguments.

Self-updates are fail-safe: if checking, downloading, verifying, installing, or restarting fails, docker-agent keeps running the current binary. Version/help/completion commands and Docker CLI plugin metadata handshakes do not trigger self-updates.

> [!NOTE]
> **Package-manager installs**
>
> Docker Desktop and Homebrew already manage docker-agent updates. Prefer those update mechanisms when you installed docker-agent that way. Self-updates are mainly intended for standalone release binaries.

## Build from Source

For the latest features, or to contribute, build from source:

### Prerequisites

- [Go 1.26](https://go.dev/dl/) or higher
- [Task](https://taskfile.dev/installation/) (build tool)
- [golangci-lint](https://golangci-lint.run/docs/welcome/install/local/) (for linting)

```bash
# Clone the repository
git clone https://github.com/docker/docker-agent.git
cd docker-agent

# Build the binary
task build

# The binary is at ./bin/docker-agent
./bin/docker-agent --help
```

> [!TIP]
> **Building on Windows**
>
> On Windows, use `task build-local` instead of `task build`. This builds the binary inside a Docker container using Docker Buildx, which avoids issues with Windows-specific toolchain setup and CGo cross-compilation. The output goes to the `./dist` directory.

## Set Up API Keys

docker-agent needs API keys for the model providers you want to use. Set them as environment variables:

```bash
# Pick one (or more) depending on your provider
export OPENAI_API_KEY="sk-..."           # OpenAI
export ANTHROPIC_API_KEY="sk-ant-..."    # Anthropic
export GOOGLE_API_KEY="AI..."            # Google Gemini (or GEMINI_API_KEY)
export MISTRAL_API_KEY="..."             # Mistral
export OPENROUTER_API_KEY="..."          # OpenRouter
```

See [Configuration Overview](../../configuration/overview/index.md#environment-variables) for the full list of supported providers and environment variables.

> [!NOTE]
> You only need the key(s) for the provider(s) you configure in your agent YAML. If you use Docker Model Runner (DMR), no API key is needed — models run locally.

## Verify Installation

```bash
# Check the version
$ docker agent version

# Run the default agent
$ docker agent run

# Or try a built-in example
$ docker agent run agentcatalog/pirate
```

## What's Next?

- [**Quick Start**](../quickstart/index.md) — create and run your first agent in under 5 minutes.
- [**Troubleshooting**](../../community/troubleshooting/index.md) — something not working? Debug mode, common issues, and solutions.
