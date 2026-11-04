---
title: Gemini (retired)
linkTitle: Gemini
weight: 70
description: |
  Gemini is no longer a built-in agent in Docker Sandboxes. Run it from the
  published Gemini kit, or move to Antigravity.
keywords: docker sandboxes, gemini, google, antigravity, ai agent, sbx, kit
---

{{% include "sandboxes-local-scope.md" %}}

Google deprecated Gemini CLI in favor of Antigravity, so Docker Sandboxes no
longer ships Gemini as a built-in agent. Starting with `sbx` v0.52.0,
`sbx run gemini` and `sbx create gemini` don't create a sandbox. They print a
notice that points to the two options below.

## Run Gemini from the published kit

Gemini CLI is still available as a published kit. Pass the kit reference
instead of the agent name:

```console
$ sbx run docker.io/sbx/gemini-kit:latest ~/my-project
```

The kit sets up Gemini CLI the same way the built-in agent did, including the
`google` credential and the MCP gateway configuration in
`~/.gemini/settings.json`. Store your Google API key with
`sbx secret set google` before you start the sandbox.

To use the kit in an environment file, keep `agent: gemini` and add the kit
reference to `kits:`:

```yaml
agent: gemini
kits:
  - docker.io/sbx/gemini-kit:latest
```

## Move to Antigravity

Antigravity replaces Gemini CLI. Run it from its published kit:

```console
$ sbx run docker.io/sbx/antigravity:latest ~/my-project
```

## Existing Gemini sandboxes

Sandboxes you created with the built-in Gemini agent keep working:

- `sbx run gemini` in the sandbox's workspace re-attaches to it and starts
  `gemini --yolo`, with a retirement warning.
- `sbx run --name <sandbox-name>` re-attaches to a Gemini sandbox by name.
- The sandbox keeps the credentials it already has.

You can't move a Gemini sandbox between local and cloud with `sbx move`, and
`sbx run --new gemini` doesn't create a new one. To start fresh, use the
Gemini kit or Antigravity.
