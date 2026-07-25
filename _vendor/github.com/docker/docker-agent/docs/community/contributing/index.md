---
title: "Contributing"
description: "docker-agent is open source. Here's how to set up your development environment and contribute."
keywords: docker agent, ai agents, community, contributing
weight: 10
canonical: https://docs.docker.com/ai/docker-agent/community/contributing/
---

_docker-agent is open source. Here's how to set up your development environment and contribute._

## Development Setup

### Prerequisites

- [Go 1.26](https://go.dev/dl/) or higher
- API key(s) for your chosen AI provider
- [Task](https://taskfile.dev/installation/)
- [golangci-lint](https://golangci-lint.run/docs/welcome/install/local/)

> [!NOTE]
> **Platform Support**
>
> macOS and Linux are fully supported for development. On Windows, use `task build-local` to build via Docker.

### Build from Source

```bash
# Clone and build
git clone https://github.com/docker/docker-agent.git
cd docker-agent
task build

# Set API keys
export OPENAI_API_KEY=your_key_here
export ANTHROPIC_API_KEY=your_key_here

# Run an example
./bin/docker-agent run examples/code.yaml
```

### Development Commands

| Command            | Description                                     |
| ------------------ | ----------------------------------------------- |
| `task build`       | Build the binary to `./bin/docker-agent`        |
| `task test`        | Run all tests (clears API keys for determinism) |
| `task lint`        | Run golangci-lint                               |
| `task format`      | Format code                                     |
| `task dev`         | Run lint, test, and build in sequence           |
| `task build-local` | Build for local platform via Docker             |
| `task cross`       | Cross-platform builds (all architectures)       |

## Dogfooding

Use docker-agent to work on docker-agent! The project includes a specialized developer agent:

```bash
cd docker-agent
docker agent run ./golang_developer.yaml
```

This agent is an expert Go developer that understands the docker-agent codebase. Ask it questions, request fixes, or have it implement features.

## Core Concepts

- **Root Agent** — Main entry point that coordinates the system
- **Sub-Agents** — Specialized agents for specific domains
- **Tools** — External capabilities via MCP
- **Models** — AI provider configurations

## Code Style

The project uses `golangci-lint` with strict rules. As long as `task lint` passes, the code is stylistically acceptable.

Key conventions:

- Use `fmt.Errorf("context: %w", err)` for error wrapping
- Always pass `context.Context` as the first parameter
- Use `slog` for structured logging
- Use functional options pattern for constructors
- In tests: use `t.Context()`, `t.TempDir()`, `t.Setenv()`, and `t.Parallel()`

## Opening Issues

File issues on the [GitHub issue tracker](https://github.com/docker/docker-agent/issues). Please:

> [!NOTE]
> **See also**
>
> [Troubleshooting](../troubleshooting/index.md) — Common issues and debug mode. [Telemetry](../telemetry/index.md) — What data is collected and how to opt out.

- Use the included issue template
- Search for existing issues before creating new ones
- Only use issues for bugs and feature requests (not support)

## Submitting Pull Requests

1. **Fork** the repository and create a branch for your changes
2. **Write** your code following the style and testing guidelines above
3. **Test** your changes: run `task lint` and `task test`
4. **Sign** your commits with `git commit -s` (DCO required)
5. **Open a pull request** against the `main` branch

> [!TIP]
> Use the dogfooding agent (`docker agent run ./golang_developer.yaml`) to help write and review your changes before submitting.

## Sign Your Work

All contributions require a Developer Certificate of Origin (DCO) sign-off:

```bash
# Sign commits automatically
git config user.name "Your Name"
git config user.email "your.email@example.com"
git commit -s -m "Your commit message"
```

## Community

Find us on [Slack](https://dockercommunity.slack.com/archives/C09DASHHRU4) for questions and discussions.

## Code of Conduct

We want to keep the docker-agent community welcoming, inclusive, and collaborative. Key guidelines:

- **Be nice** — Be courteous, respectful, and polite. No abuse of any kind will be tolerated.
- **Encourage diversity** — Make everyone feel welcome regardless of background.
- **Keep it legal** — Share only content you own and don't break the law.
- **Stay on topic** — Post to the correct channel and avoid off-topic discussions.

The governance for this repository is handled by Docker Inc.
