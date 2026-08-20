---
title: Share agent skills
linkTitle: Agent skills
weight: 40
description: Import skills from supported host agents into a persistent store shared with Docker Sandboxes.
keywords: docker sandboxes, sbx, agent skills, shared skills, claude code, codex, copilot, cursor, droid
---

Shared agent skills make skills from supported agents on your host available
inside your sandboxes. Importing copies the skills into a persistent store that
survives sandbox deletion and is shared by default with new sandboxes that run
a supported agent.

> [!NOTE]
> Shared agent skills are experimental.

Preview the skills that `sbx` finds without copying them:

```console
$ sbx skills import --dry-run
```

The command scans the following directories in order and copies each skill
subdirectory into the shared store. When the sandbox starts, `sbx` mounts the
store at the path the agent reads inside the sandbox.

| Agent       | Host source         | Sandbox mount target          |
| ----------- | ------------------- | ----------------------------- |
| Claude Code | `~/.claude/skills`  | `/home/agent/.claude/skills`  |
| Codex       | `~/.agents/skills`  | `/home/agent/.agents/skills`  |
| Copilot     | `~/.copilot/skills` | `/home/agent/.copilot/skills` |
| Cursor      | `~/.cursor/skills`  | `/home/agent/.cursor/skills`  |
| Droid       | `~/.factory/skills` | `/home/agent/.factory/skills` |

All imported skills go into the same store, regardless of their source. If
more than one source contains a skill with the same directory name, the skill
from the first source in the table wins and `sbx` warns about the others.

Import the skills:

```console
$ sbx skills import
```

The final output reports the shared store path. The default locations are:

| Platform | Shared store path                                                           |
| -------- | --------------------------------------------------------------------------- |
| macOS    | `~/Library/Application Support/com.docker.sandboxes/sandboxes/agent-skills` |
| Linux    | `~/.local/state/sandboxes/sandboxes/agent-skills`                           |
| Windows  | `%LOCALAPPDATA%\DockerSandboxes\sandboxes\state\agent-skills`               |

On Linux, `sbx` uses `$XDG_STATE_HOME/sandboxes/sandboxes/agent-skills` when
`XDG_STATE_HOME` is set.

When a skill already exists in the store, `sbx` prompts before replacing it.
Use `--force` to replace existing skills without prompts. Importing replaces
the complete skill directory rather than merging files. Run the import command
again when you want to copy updates from the host. Running `sbx reset` clears
the shared store.

Sandboxes created with `sbx` version 0.37.0 or later for a supported agent are
configured to mount the store read-write by default. These sandboxes mount the
current contents of the store each time they start, so you can import skills
before or after creating them. To create a sandbox without the shared store,
use `--no-share-skills`:

```console
$ sbx run --no-share-skills claude
```

Upgrading `sbx` does not enable shared skills for sandboxes created with an
earlier version. Remove and recreate those sandboxes after upgrading. The
`--no-share-skills` option also only applies when the sandbox is created. To
turn off shared skills for an existing sandbox, remove it and recreate it with
the option.

> [!WARNING]
> The shared skills store is mounted read-write. A sandbox can modify any skill
> in the store, and another sandbox can later load the modified instructions or
> run the modified scripts. The store is dedicated sandbox state, so this does
> not by itself execute the modified skill on your host. It does put every
> sandbox that shares the store in the same trust boundary. Use
> `--no-share-skills` to keep a sandbox outside that boundary.

Some agents scan for skills when a session starts. If imported skills don't
appear in an existing session, start another agent session.
