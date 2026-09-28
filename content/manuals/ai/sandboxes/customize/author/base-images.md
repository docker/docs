---
title: Base images for sandbox workloads
linkTitle: Base images
weight: 10
description: Choose a Docker-provided agent image or your own Linux base image for a sandbox workload kit.
keywords: sandboxes, sbx, kits, base images, templates, dockerfile, custom agents
aliases:
  - /ai/sandboxes/customize/templates/
---

{{< summary-bar feature_name="Docker Sandboxes sbx" >}}

A base image gives your workload its operating system and starting set of
tools. Docker provides images with agents already installed. You can also start
from another Linux image and prepare it yourself. Choose the base image in your
workload's Dockerfile. Use the kit descriptor to configure network access,
credentials, and other sandbox behavior. A kit set uses the base image from its
workload.

To save and reuse an environment you've configured interactively, see
[Save a sandbox as a template](/manuals/ai/sandboxes/usage.md#saving-a-sandbox-as-a-template).

## Docker-provided images

Docker's sandbox templates are published as
`docker/sandbox-templates:<variant>`. They are based on Ubuntu and run as a
non-root `agent` user with sudo access. Most variants include Git, Docker
CLI, and common development tools like Node.js, Python, Go, and Java.

| Variant               | Agent                                                                |
| --------------------- | -------------------------------------------------------------------- |
| `claude-code`         | [Claude Code](https://claude.ai/download)                            |
| `claude-code-minimal` | Claude Code with a minimal toolset (no Node.js, Python, Go, or Java) |
| `codex`               | [OpenAI Codex](https://github.com/openai/codex)                      |
| `copilot`             | [GitHub Copilot](https://github.com/github/copilot-cli)              |
| `cursor-agent`        | [Cursor](https://cursor.com/cli)                                     |
| `devin`               | [Devin CLI](https://docs.devin.ai/work-with-devin/devin-cli)         |
| `docker-agent`        | [Docker Agent](https://github.com/docker/docker-agent)               |
| `droid`               | [Droid](https://www.factory.ai)                                      |
| `gemini`              | [Gemini CLI](https://github.com/google-gemini/gemini-cli)            |
| `kiro`                | [Kiro](https://kiro.dev)                                             |
| `opencode`            | [OpenCode](https://opencode.ai)                                      |
| `shell`               | No agent pre-installed. Use for manual agent setup.                  |

## Use an image in a v3 workload

Start your workload's Dockerfile with `FROM`, then add the tools and
configuration you need. Install system packages as `root`, and switch back
to `agent` before installing tools in the agent's home directory. Running
those installers as `root` puts files under `/root/`, where the agent can't
use them.

### What the base image provides

`FROM` inherits the image's files and settings. If the image is also a
published kit, its capabilities don't carry over. Declare network access,
credentials, storage, and hooks in your own descriptor. To keep an existing
kit's capabilities and add tools, [compose a kit set](/manuals/ai/sandboxes/customize/author/kit-sets.md).

Your descriptor determines the kit format version. A v3 descriptor creates
a v3 kit, including when you use a Docker template image as its base. See
[Version compatibility](/manuals/ai/sandboxes/customize/_index.md#version-compatibility).

### Package an existing agent image

This example packages Docker's OpenCode image as a v3 workload. The image
already has an agent installed, so you only need to choose its launch command
and describe what it needs to run. You can run the resulting kit directly or
include it in a set with additional tools.

Create a directory with these two files:

```text
opencode-workload/
├── opencode-workload.yaml
└── opencode-workload.dockerfile
```

The Dockerfile selects the base image and tells Docker Sandboxes to run
OpenCode as the `agent` user:

```dockerfile {title="opencode-workload/opencode-workload.dockerfile"}
FROM docker/sandbox-templates:opencode
USER agent
ENTRYPOINT ["opencode"]
CMD []
```

The template supplies OpenCode, Python, uv, and the `agent` user. Create a
YAML descriptor to identify this as a workload and declare what OpenCode
needs: network access, an Anthropic API key, and instructions about its
environment:

```yaml {title="opencode-workload/opencode-workload.yaml"}
# syntax=docker/sandbox-kit:3
schemaVersion: "3"
kind: workload

capabilities:
  - type: com.docker.sandbox/sbx@1
  - type: com.docker.sandbox/network-policy@1
    config:
      runtime:
        allow:
          - api.anthropic.com
          - opencode.ai
          - models.dev
          - registry.npmjs.org
          - pypi.org
          - files.pythonhosted.org
  - type: com.docker.sandbox/credential@1
    config:
      service: anthropic
      phase: runtime
      apiKey:
        name: ANTHROPIC_API_KEY
        proxyManaged: true
        inject:
          - domain: api.anthropic.com
            header: x-api-key
            format: "%s"
  - type: com.docker.sandbox/agent-context@1
    config:
      filename: AGENTS.md
      content: |
        OpenCode runs as the agent user. Python and uv are available.
        Use the project's environment and dependency configuration.
```

The credential entry names the service and describes how to authenticate API
requests. Store the actual API key on your host.

`sbx` can build this directory when you create a sandbox. The result is a
container image with the kit's files and descriptor, which you can also
publish to a registry for others to use.

## Use your own Linux image

If you need a different operating system or set of packages, start from a
Linux image of your choice. You'll need to add the tools and user account
that Docker Sandboxes expects, then install your agent.

The following requirements help you prepare that image. For a step-by-step
example, follow [Build an agent workload](/manuals/ai/sandboxes/customize/author/build-an-agent.md).

### Base image requirements

Prepare your image with the following:

- Tools: Include `curl`, `git`, and trusted CA certificates for accessing
  source repositories and making HTTPS requests.
- Shells: Provide executable `/bin/sh` and `/bin/bash` files for setup commands
  and agent launch.
- User account: Create a non-root `agent` account with UID 1000 and home
  directory `/home/agent`. Add the account to `/etc/passwd` and set `USER
  agent` in the image.
- Launch command: Set `ENTRYPOINT` or `CMD` to the agent or shell you want the
  sandbox to launch.

Docker Sandboxes uses the account's entry in `/etc/passwd` to determine
the user ID, group ID, and home directory when running commands and writing
files.

For a sandbox with a mounted workspace, `WORKDIR` doesn't choose the mount
location. For a sandbox without a mounted workspace, `sbx` uses the image's
absolute `WORKDIR` as the working directory. If the image doesn't specify a
usable absolute path, `sbx` falls back to `/home/agent/workspace`.

### Declare the agent launch contract

`sbx` keeps the sandbox running and starts the agent as a separate process,
using the image's launch command. The `com.docker.sandbox/sbx@1` capability
declares that your workload is prepared to run this way:

```yaml
capabilities:
  - type: com.docker.sandbox/sbx@1
```

This capability takes no `config` and doesn't request extra permissions. It
lets the upstream conformance tests check that your image meets the shell and
user account requirements before you publish.

The host launches the agent through non-interactive Bash. To load persistent
environment settings at launch, set `BASH_ENV` to an absolute path and include
that file in the image. Bash reads this file instead of login profiles or
`.bashrc`.

For the complete contract, see the upstream
[`sbx@1` definition](https://github.com/docker/sandbox-kit-spec/blob/main/docs/spec/capabilities/com.docker.sandbox/sbx@1.md).
