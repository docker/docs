---
title: Supported agents
linkTitle: Agents
weight: 40
description: AI coding agents supported by Docker Sandboxes.
keywords: docker sandboxes, ai agents, claude code, codex, copilot, cursor, droid, gemini, kiro, sandbox kits
aliases:
  - /ai/sandboxes/agents/copilot/
  - /ai/sandboxes/agents/droid/
  - /ai/sandboxes/agents/kiro/
---

Docker Sandboxes runs the following agents out of the box:

- [Claude Code](claude-code/)
- [Codex](codex/)
- [Cursor](cursor/)
- [Docker Agent](docker-agent/)
- [Gemini](gemini/)
- [OpenCode](opencode/)
- [Shell](shell/) — agent-less sandbox for manual setup or testing

Copilot, Droid, and Kiro are available as community kits:

```console
$ sbx run docker.io/sbx/copilot-kit:latest
$ sbx run docker.io/sbx/droid-kit:latest
$ sbx run docker.io/sbx/kiro-kit:latest
```

Want to pre-install tools or customize an agent's environment?
See [Customize](../customize/).
