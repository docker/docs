---
title: Supported agents
linkTitle: Agents
weight: 40
description: AI coding agents supported by Docker Sandboxes.
keywords: docker sandboxes, ai agents, claude code, codex, cursor, devin, gemini
---

Docker Sandboxes runs the following agents in local sandboxes:

- [Claude Code](claude-code/)
- [Codex](codex/)
- [Copilot](copilot/)
- [Cursor](cursor/)
- [Devin](devin/)
- [Docker Agent](docker-agent/)
- [Droid](droid/)
- [Gemini](gemini/)
- [Kiro](kiro/)
- [OpenCode](opencode/)
- [Shell](shell/) — agent-less sandbox for manual setup or testing

Want to pre-install tools or customize an agent's environment?
See [Customize](../customize/).

## Forward arguments to an agent

Pass agent arguments after `--` in `sbx run`. If the first argument is a
flag, the launcher keeps the default arguments. A prompt or subcommand replaces
the image's `CMD` defaults while preserving its `ENTRYPOINT` arguments.

The startup commands on the agent pages show the combined command, not how the
image divides it between `ENTRYPOINT` and `CMD`. To check which defaults a
prompt or subcommand keeps, inspect the agent image you use. For image-based
kits, use the [kit image reference](../customize/use-kits.md#run-a-kit), rather
than the sandbox template listed under Base image on the agent page:

```console
$ docker pull <agent-image>
$ docker image inspect <agent-image> --format 'ENTRYPOINT={{json .Config.Entrypoint}} CMD={{json .Config.Cmd}}'
```

For example, an image with `ENTRYPOINT ["agent", "--yolo"]` keeps `--yolo`
when you pass a prompt. An image with `ENTRYPOINT ["agent"]` and
`CMD ["--yolo"]` drops it. Use the exact image tag or digest you run, since
this split can differ between images.
