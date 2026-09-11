---
title: Kit examples
linkTitle: Examples
description: Build schema v3 kits that add tools, shared files, runtime configuration, lifecycle hooks, and agent instructions to a sandbox workload.
keywords: sandboxes, sbx, kits, mixins, examples, capabilities, build, lifecycle
weight: 40
---

{{< summary-bar feature_name="Docker Sandboxes sbx" >}}

> [!NOTE]
> Kits are experimental. The kit file format, CLI commands, and experience
> for creating, loading, and managing kits are subject to change. Share
> feedback in the [docker/sbx-releases](https://github.com/docker/sbx-releases)
> repository.

These schema v3 examples show how to add tools, configuration, and instructions
to a workload. Each example includes the files needed to run it locally.
For schema v2 patterns, see [Schema v2 kit examples](kits-v2/kit-examples.md).
For field definitions, see the [Kit spec reference](kit-reference.md).

## Create a workload for the examples

Create a `shell-v3` directory and save this descriptor inside it:

```yaml {title="shell-v3/shell-v3.yaml"}
# syntax=docker/runtime-kit:3
schemaVersion: "3"
kind: workload
displayName: Example shell

build: |
  FROM docker/sandbox-templates:shell
  USER agent
  ENTRYPOINT ["bash"]
  CMD ["-l"]

capabilities:
  - type: com.docker.runtime/agent-context@1
    config:
      filename: AGENTS.md
      content: This shell is an environment for testing composed kits.
```

The inline `build` is a Dockerfile that produces a complete workload. Run it
from the parent directory, using your current directory as the workspace:

```console
$ sbx run --name kit-shell ./shell-v3 .
```

Use this workload with the mixins that follow. Each example uses a different
sandbox name because adding kits requires creating a sandbox. Schema v3
mixins need a schema v3 workload. The built-in agent names use earlier kit
formats and can't be combined with these mixins.

## Copy shared configuration

Ship static files in the image, then use a lifecycle install hook to copy
them to a destination that only exists when the sandbox is created. This is
useful for workspace defaults, because the workspace is mounted at runtime
and its path varies between projects.

Create this directory:

```text
team-config/
├── team-config.yaml
├── team-config.dockerfile
└── editorconfig
```

```ini {title="team-config/editorconfig"}
root = true

[*]
charset = utf-8
indent_style = space
indent_size = 2
insert_final_newline = true
```

The Dockerfile places the file outside the workspace, where a runtime mount
won't hide it:

```dockerfile {title="team-config/team-config.dockerfile"}
FROM scratch
COPY editorconfig /usr/share/team-config/editorconfig
```

The descriptor copies the default into the workspace only when the project
has no `.editorconfig`. It also keeps a copy in the agent's home directory:

```yaml {title="team-config/team-config.yaml"}
# syntax=docker/runtime-kit:3
schemaVersion: "3"
kind: mixin
displayName: Team configuration

capabilities:
  - type: com.docker.runtime/lifecycle@1
    config:
      install:
        - command: |
            set -eu
            if [ ! -e "$WORKSPACE_DIR/.editorconfig" ]; then
              cp /usr/share/team-config/editorconfig "$WORKSPACE_DIR/.editorconfig"
            fi
            mkdir -p /home/agent/.config/team
            cp /usr/share/team-config/editorconfig /home/agent/.config/team/editorconfig
          user: agent
          env: [WORKSPACE_DIR]
          description: Copy team defaults into the workspace and home
```

Run it with the example workload:

```console
$ sbx run --name kit-team-config ./shell-v3 --kit ./team-config .
```

`WORKSPACE_DIR` is a runtime value. Declaring it in the hook's `env` list
makes it available to the command. Running as `agent` keeps the copied files
writable by the workload user. With a directly mounted workspace, creating
`.editorconfig` also creates that file in the host project.

Schema v3 has no automatic `files/workspace/` or `files/home/` placement.
The Dockerfile defines where static content lives in the image. The hook
handles the runtime destination.

## Install an internal CA certificate

If your organization uses a proxy that inspects HTTPS traffic, add its root
CA certificate to the sandbox's trust store. Create an `internal-ca`
directory and save the PEM-encoded certificate as `internal-ca.crt` beside
these two files:

```dockerfile {title="internal-ca/internal-ca.dockerfile"}
FROM scratch
COPY internal-ca.crt /usr/local/share/ca-certificates/team-internal-ca.crt
```

```yaml {title="internal-ca/internal-ca.yaml"}
# syntax=docker/runtime-kit:3
schemaVersion: "3"
kind: mixin
displayName: Internal CA certificate

capabilities:
  - type: com.docker.runtime/lifecycle@1
    config:
      install:
        - command: update-ca-certificates
          user: "0"
          description: Add the internal CA to the sandbox trust store
```

The build includes the certificate at a fixed path with a `.crt` extension.
The install hook updates the workload's system trust store after the overlay
is applied. Tools that use that trust store can then verify certificates
signed by the internal CA.

```console
$ sbx run --name kit-ca ./shell-v3 --kit ./internal-ca .
```

Use a distinct certificate filename for each CA. If the proxy uses several
root certificates, copy each one into `/usr/local/share/ca-certificates/`
before running `update-ca-certificates`.

## Build a tool overlay

Build tools into a mixin so each sandbox can use the same image layers.
This example compiles [`gojq`](https://github.com/itchyny/gojq), a JSON query
tool with `jq` syntax, and ships its binary without the Go compiler.

Create a `gojq` directory with these two files:

```yaml {title="gojq/gojq.yaml"}
# syntax=docker/runtime-kit:3
schemaVersion: "3"
kind: mixin
displayName: JSON query tool

args:
  version:
    default: "0.12.17"
    pattern: '^[0-9]+\.[0-9]+\.[0-9]+$'
    buildArg: GOJQ_VERSION

provides: ["gojq@${{ kit.args.version }}"]
```

```dockerfile {title="gojq/gojq.dockerfile"}
FROM golang:1.25 AS build
ARG GOJQ_VERSION
RUN CGO_ENABLED=0 go install github.com/itchyny/gojq/cmd/gojq@v${GOJQ_VERSION}

FROM scratch
COPY --from=build /go/bin/gojq /usr/local/bin/gojq
```

The final `FROM scratch` stage contains the tool overlay. The binary is
compiled without C dependencies, so it doesn't require shared libraries
from the workload. For other tools, include their runtime libraries or
declare a dependency on a compatible environment.

There is no lifecycle install hook: the compiler and module downloads run
when the kit is built. The resulting tool is available as soon as the
composed filesystem is ready.

```console
$ sbx run --name kit-gojq ./shell-v3 --kit ./gojq .
```

Inside the sandbox, query some JSON:

```console
$ printf '%s\n' '{"kit": "gojq", "ready": true}' | gojq '.ready'
true
```

Another mixin can declare `requires: ["gojq >= 0.12.17"]`. You must include
both mixins in the launch command: a requirement checks the supplied set
and doesn't download a provider. See [Kit composition](kits.md).

## Write runtime configuration

Use lifecycle `files` when a file contains values chosen at sandbox creation.
This example writes a team settings file using a validated kit argument.
It needs no Dockerfile because it contributes only declarations:

```yaml {title="workspace-config/workspace-config.yaml"}
# syntax=docker/runtime-kit:3
schemaVersion: "3"
kind: mixin
displayName: Workspace settings

args:
  project:
    default: demo
    pattern: '^[a-z][a-z0-9-]*$'

capabilities:
  - type: com.docker.runtime/lifecycle@1
    config:
      files:
        - path: /home/agent/.config/team/project.json
          content: |
            {"project": "${{ kit.args.project }}"}
          mode: "0644"
```

Create the `workspace-config` directory, save the descriptor, and supply the
project name when creating the sandbox:

```console
$ sbx run --name kit-project ./shell-v3 --kit ./workspace-config \
    --kit-arg workspace-config.project=payments .
```

The runtime expands `${{ kit.args.project }}` and writes the file before the
workload starts. The argument prefix is the local kit directory's name.
The file is written during sandbox creation. Add `overwrite: false` if it
should only seed a default and preserve an existing file.

Use absolute paths under `/home/agent/` for files the agent owns. A variable
such as `$HOME` in `path` isn't expanded. To write to a runtime workspace
path, use a hook with `env: [WORKSPACE_DIR]`, as in the shared configuration
example.

## Run a hook on every start

Use lifecycle `startup` for work that must repeat after the sandbox stops
and starts. This example records the time and workspace at each start:

```yaml {title="start-log/start-log.yaml"}
# syntax=docker/runtime-kit:3
schemaVersion: "3"
kind: mixin
displayName: Sandbox start log

capabilities:
  - type: com.docker.runtime/lifecycle@1
    config:
      startup:
        - command: |
            printf '%s %s\n' "$(date -u +%FT%TZ)" "$WORKSPACE_DIR" >> /home/agent/sandbox-starts.log
          user: agent
          env: [WORKSPACE_DIR]
          description: Record each sandbox start
```

```console
$ sbx run --name kit-start-log ./shell-v3 --kit ./start-log .
```

Startup hooks run separately from the workload launch. Use install hooks
or lifecycle files for configuration the agent must read when it starts.
Design repeated setup so running it again leaves the environment usable.

For a long-running service, set `background: true` on its startup hook and
redirect output to a log file. The service binary must already be present in
the workload or a composed tool overlay.

## Contribute agent instructions

Use the agent-context capability to tell an agent how to use tools and
configuration supplied by a kit. A mixin contributes instructions to the
workload's context profile without choosing its filename:

```yaml {title="team-review/team-review.yaml"}
# syntax=docker/runtime-kit:3
schemaVersion: "3"
kind: mixin
displayName: Team review instructions

capabilities:
  - type: com.docker.runtime/agent-context@1
    config:
      content: |
        When reviewing a Dockerfile, check the base image version, layer
        ordering, cache use, and whether secrets appear in ARG or ENV.
        Explain the effect of each suggested change and run available
        project checks before reporting completion.
```

The example shell workload chooses `AGENTS.md`. When you compose this mixin,
`sbx` adds a kit entry to that profile and puts the instructions in a separate
file for the agent to read on demand. An agent workload can choose another
profile, such as `CLAUDE.md`, and the same mixin contributes to that profile.

```console
$ sbx run --name kit-team-review ./shell-v3 --kit ./team-review .
```

For longer instructions, set `contentFile: ./review.md` instead of `content`
and keep the Markdown file beside the descriptor. The build frontend stages
the file in the image. Only workload kits can set `filename`.

To build an agent workload step by step, see [Build an agent](build-an-agent.md).
