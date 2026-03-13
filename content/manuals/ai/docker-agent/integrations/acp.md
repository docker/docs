---
linkTitle: ACP
title: ACP integration
description: Configure your editor or IDE to use agents as coding assistants
keywords: [docker agent, acp, editor, ide, neovim, zed, integration]
weight: 40
---

Run agents directly in your editor using the Agent Client Protocol (ACP).
Your agent gets access to your editor's filesystem context and can read and
modify files as you work. The editor handles file operations while Docker Agent
provides the AI capabilities.

This guide shows you how to configure Neovim, or Zed to run agents with Docker Agent. If
you're looking to expose agents as tools to MCP clients like Claude
Desktop or Claude Code, see [MCP integration](./mcp.md) instead.

## How it works

When you run Docker Agent with ACP, it becomes part of your editor's environment. You
select code, highlight a function, or reference a file - the agent sees what you
see. No copying file paths or switching to a terminal.

Ask "explain this function" and the agent reads the file you're viewing. Ask it
to "add error handling" and it edits the code right in your editor. The agent
works with your editor's view of the project, not some external file system it
has to navigate.

The difference from running Docker Agent in a terminal: file operations go through
your editor instead of the agent directly accessing your filesystem. When the
agent needs to read or write a file, it requests it from your editor. This keeps
the agent's view of your code synchronized with yours - same working directory,
same files, same state.

## Prerequisites

Before configuring your editor, you need:

- **Docker Agent installed** - See the [installation guide](../_index.md#installation)
- **Agent configuration** - A YAML file defining your agent. See the
  [tutorial](../tutorial.md) or [example
  configurations](https://github.com/docker/docker-agent/tree/main/examples)
- **Editor with ACP support** - Neovim, Intellij, Zed, etc.

Your agents will use model provider API keys from your shell environment
(`ANTHROPIC_API_KEY`, `OPENAI_API_KEY`, etc.). Make sure these are set before
launching your editor.

## Editor configuration

### Zed

Zed has built-in ACP support.

1. Add `docker agent` to your agent servers in `settings.json`:

   ```json
   {
     "agent_servers": {
       "my-agent-team": {
         "command": "docker",
         "args": ["agent", "serve", "acp", "agent.yml"]
       }
     }
   }
   ```

   Replace:
   - `my-agent-team` with the name you want to use for the agent
   - `agent.yml` with the path to your agent configuration file.

   If you have multiple agent files that you like to run separately, you can
   create multiple entries under `agent_servers` for each agent.

2. Start a new external agent thread. Select your agent in the drop-down list.

   ![New external thread with Docker Agent in Zed](../images/cagent-acp-zed.avif)

### Neovim

Use the [CodeCompanion](https://github.com/olimorris/codecompanion.nvim) plugin,
which has native support for Docker Agent through a built-in adapter:

1. [Install CodeCompanion](https://codecompanion.olimorris.dev/installation)
   through your plugin manager.
2. Extend the `dockeragent` adapter in your CodeCompanion config:

   ```lua
   require("codecompanion").setup({
     adapters = {
       acp = {
         dockeragent = function()
           return require("codecompanion.adapters").extend("dockeragent", {
             commands = {
               default = {
                 "docker",
                 "agent", 
                 "serve",
                 "acp",
                 "agent.yml",
               },
             },
           })
         end,
       },
     },
   })
   ```

   Replace `agent.yml` with the path to your agent configuration file. If you
   have multiple agent files that you like to run separately, you can create
   multiple commands for each agent.

3. Restart Neovim and launch CodeCompanion:

   ```plaintext
   :CodeCompanion
   ```

4. Switch to the Docker Agent adapter (keymap `ga` in the CodeCompanion buffer, by
   default).

See the [CodeCompanion ACP
documentation](https://codecompanion.olimorris.dev/usage/acp-protocol) for more
information about ACP support in CodeCompanion. Note that terminal operations
are not supported, so [toolsets](../reference/toolsets.md) like `shell` or
`script_shell` are not usable through CodeCompanion.

## Agent references

You can specify your agent configuration as a local file path or OCI registry
reference:

```console
# Local file path
$ docker agent serve acp ./agent.yml

# OCI registry reference
$ docker agent serve acp agentcatalog/pirate
$ docker agent serve acp dockereng/myagent:v1.0.0
```

Use the same syntax in your editor configuration:

```json
{
  "agent_servers": {
    "myagent": {
      "command": "docker",
      "args": ["agent", "serve", "acp", "agentcatalog/pirate"]
    }
  }
}
```

Registry references enable team sharing, version management, and clean
configuration without local file paths. See [Sharing
agents](../sharing-agents.md) for details on using OCI registries.

## Testing your setup

Verify your configuration works:

1. Start the Docker Agent ACP server using your editor's configured method
2. Send a test prompt through your editor's interface
3. Check that the agent responds
4. Verify filesystem operations work by asking the agent to read a file

If the agent starts but can't access files or perform other actions, check:

- Working directory in your editor is set correctly to your project root
- Agent configuration file path is absolute or relative to working directory
- Your editor or plugin properly implements ACP protocol features

## What's next

- Review the [configuration reference](../reference/config.md) for advanced
  agent setup
- Explore the [toolsets reference](../reference/toolsets.md) to learn what tools
  are available
- Add [RAG for codebase search](../rag.md) to your agent
- Check the [CLI reference](../reference/cli.md) for all `docker agent serve acp` options
- Browse [example
  configurations](https://github.com/docker/docker-agent/tree/main/examples) for
  inspiration
