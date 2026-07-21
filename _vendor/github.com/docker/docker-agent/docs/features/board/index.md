---
title: "Kanban Board"
description: "Orchestrate multiple agents from a Kanban TUI: each card runs an agent in a tmux session on an isolated git worktree."
keywords: docker agent, ai agents, features, board, kanban, orchestration
linkTitle: "Kanban Board"
weight: 15
canonical: https://docs.docker.com/ai/docker-agent/features/board/
---

_Board is a Kanban TUI for orchestrating agents. Each card launches an agent
in a tmux session on an isolated git worktree, and moving a card forward
through the pipeline sends the destination column's prompt to its agent._

## Launching the board

```bash
$ docker agent board
```

Requirements: `tmux` and `git` must be installed.

## How it works

- **Cards run agents.** Creating a card (`n`) launches `docker agent run` in
  a dedicated tmux session, working in a fresh git worktree branched from the
  project's upstream default branch. The card's title, running/idle status,
  and failures are mirrored live from the agent's control plane.
- **Startup phases.** While an agent is coming up its card moves through three
  intermediate statuses before reaching **running**: `starting` (tmux session
  created, process booting, no worktree yet) → `loading` (worktree present;
  agent loading config, models, and tools) → `attaching` (control-plane socket
  bound; board waiting for the first snapshot).
- **Columns are a pipeline.** The default pipeline is
  Dev → Review → Push → Done, and it's fully customizable: manage columns
  from the board (`c`) or in the config file. Moving a card forward (`]`)
  sends the destination column's prompt to the card's agent; moving it back
  (`[`) sends nothing.
- **Attach anytime.** Press `enter` (or double-click a card) to attach your
  terminal to the agent's session and interact with it directly; `ctrl+q`
  detaches and returns to the board.
- **Everything is recoverable.** Quitting the board leaves agents running in
  tmux; restarting it reattaches to them. If an agent process dies, the board
  relaunches it and resumes the same conversation and worktree. An agent that
  keeps crashing at startup turns its card red instead of relaunching
  forever: attach to it (`enter`) to read the error output, then move the
  card forward to relaunch it, or delete it.

## Key bindings

| Key           | Action                                              |
| ------------- | --------------------------------------------------- |
| `n`           | Create a card (project + prompt)                    |
| `enter`       | Attach to the card's agent (`ctrl+q` detaches)      |
| `d`           | View the card's worktree diff                       |
| `o`           | Open the card's worktree in `$DOCKER_AGENT_BOARD_EDITOR` (`code`) |
| `s`           | Open an interactive shell in the card's worktree    |
| `[` / `]`     | Move the card back / forward                        |
| `1`-`9`       | Move the card to column N                           |
| `x`           | Delete the card, its session, worktree, and branch  |
| `p`           | Manage projects (add, edit, reorder, remove)        |
| `c`           | Manage columns (add, edit, reorder, remove)         |
| `e`           | Edit the selected column's prompt                   |
| `←↓↑→` `hjkl` | Navigate                                            |
| mouse         | Click selects, double-click attaches, drag moves, wheel scrolls |
| `?`           | Help                                                |
| `q`           | Quit (agents keep running)                          |

## Configuration

Everything is configured in the global config file
(`~/.config/cagent/config.yaml`) or through the TUI itself (`p` for projects,
`c` for columns, `e` for column prompts):

```yaml
board:
  projects:
    - name: my-app
      path: /Users/me/src/my-app
      agent: coder # any agent ref; defaults to the built-in agent
  columns:
    - id: dev
      name: Dev
      emoji: 🔨
    - id: review
      name: Review
      emoji: 🔍
      prompt: Review the local changes and fix any issues you find.
    - id: done
      name: Done
      emoji: ✅
```

Omitting `columns` keeps the default pipeline. Column `id`s identify a
column across renames (cards remember the column they are in by id); when
omitted, the id is derived from the column's name. When a card enters a
column with a `prompt`, that prompt is delivered to the card's agent as its
next message.
