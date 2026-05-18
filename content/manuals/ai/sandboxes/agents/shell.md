---
title: Shell
weight: 90
description: Run an agent-less sandbox with a Bash login shell for manual setup, testing custom agent implementations, or inspecting a running environment.
keywords: sandboxes, sbx, shell, agent, manual setup, testing
---

{{< summary-bar feature_name="Docker Sandboxes sbx" >}}

`sbx run shell` drops you into a Bash login shell inside a sandbox with no
pre-installed agent binary. It's useful for installing and configuring
agents manually, testing custom implementations, or inspecting a running
environment.

```console
$ sbx run shell ~/my-project
```

The workspace path defaults to the current directory. To run a one-off
command instead of an interactive shell, pass it after `--`:

```console
$ sbx run shell -- -c "echo 'Hello from sandbox'"
```

Set your API keys as environment variables so the sandbox proxy can inject
them into API requests automatically. Credentials are never stored inside
the VM:

```console
$ export ANTHROPIC_API_KEY=sk-ant-xxxxx
$ export OPENAI_API_KEY=sk-xxxxx
```

Once inside the shell, you can install agents using their standard methods,
for example `npm install -g @continuedev/cli`. For complex setups, build a
[custom template](../customize/templates.md) instead of installing
interactively each time.

## Base image

The shell sandbox uses the `shell` base image — the common base environment
without a pre-installed agent.
