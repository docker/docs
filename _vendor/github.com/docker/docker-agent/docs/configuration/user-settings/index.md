---
title: "User Settings"
description: "Full reference for the global settings block in ~/.config/cagent/config.yaml: theme, layout, YOLO, sound, snapshots, permissions, hooks, keybindings, and how they interact with agent config and CLI flags."
keywords: docker agent, ai agents, configuration, yaml, user settings, user config
linkTitle: "User Settings"
weight: 110
canonical: https://docs.docker.com/ai/docker-agent/configuration/user-settings/
---

_Full reference for the global settings block in your user config file._

## Where Settings Live

Docker Agent reads a single user-level config file, independent of any agent YAML:

```text
~/.config/cagent/config.yaml
```

The `settings:` block inside it holds preferences that apply to every agent you run — appearance, behavior, notifications, and a few global safety defaults. Everything under `settings:` is optional; an unset field falls back to the documented default.

```yaml
# ~/.config/cagent/config.yaml
settings:
  theme: dracula
  lean: false
  sound: true
```

You rarely need to hand-edit this file. Most fields are managed from the TUI's `/settings` dialog (**Appearance**, **Behavior**, **Notifications** tabs) — press <kbd>Enter</kbd> there to apply and persist a change. A few fields (`permissions`, `hooks`, `keybindings`) have no dialog UI and are only set by editing the file directly.

> [!NOTE]
> This page documents `settings:`. The user config file also has top-level sections outside `settings:` — `aliases:`, `providers:`, `board:`, `credential_helper:`, `sandbox_allowlist:` — which are not covered here.

## Settings Reference

| Setting | Type | Default | Description |
| --- | --- | --- | --- |
| `hide_tool_results` | boolean | `false` | Hide tool call results in the TUI by default. Mirrors the `--hide-tool-results` flag and the <kbd>Ctrl</kbd>+<kbd>O</kbd> toggle. |
| `expand_thinking` | boolean | `false` | Start new sessions with thinking/tool blocks expanded instead of collapsed. |
| `split_diff_view` | boolean | `true` | Render file-edit diffs side-by-side instead of unified. |
| `render_images` | boolean | `true` | Render images in the TUI using the Kitty graphics protocol. Applies to both tool-result images and Markdown images in agent responses. Automatically disabled when the terminal does not support Kitty. |
| `theme` | string | `default` | Theme name, loaded from a built-in theme or `~/.cagent/themes/<name>.yaml`. The special value `auto` follows the terminal's light/dark background. See [Theming](../../features/tui/index.md#theming). |
| `theme_dark` | string | `default` | Theme applied when `theme: auto` and the terminal background is dark. |
| `theme_light` | string | `default-light` | Theme applied when `theme: auto` and the terminal background is light. |
| `YOLO` | boolean | `false` | Auto-approve all tool calls globally, across every agent you run. Mirrors the `--yolo` flag and the `/yolo` command. |
| `lean` | boolean | `false` | Make the [lean TUI](../../features/tui/index.md#lean-tui) (simplified, minimal-chrome interface) the default for interactive runs instead of the full TUI. |
| `tab_title_max_length` | int | `20` | Maximum display length for tab titles; longer titles are truncated with an ellipsis. |
| `restore_tabs` | boolean | `false` | Restore previously open tabs when launching the TUI. |
| `sound` | boolean | `false` | Play a notification sound on task success or failure. |
| `sound_threshold` | int | `10` | Minimum duration in seconds a task must run before a success sound plays (failures always play). |
| `snapshot` | boolean | `false` | Enable automatic shadow-git snapshots at turn boundaries globally. See [Snapshots](../../features/snapshots/index.md). |
| `cache_stable_prompts` | boolean | `false` | Keep changing trusted context (date, environment info, dynamic prompt files) out of the frozen system prefix and append chronological updates instead, improving prompt-cache hit rates on long sessions. |
| `warn_on_cache_miss` | boolean | `false` | Warn when a model call after the first one in a session reports no cached input tokens (a prompt-cache miss). Managed from the **Notifications** tab of `/settings`. |
| `busy_send_mode` | string | `steer` | What happens to a message sent while the agent is working: `steer` injects it into the ongoing stream; `queue` holds it until the current turn ends. |
| `permissions` | object | _unset_ | Global tool-permission rules (`allow` / `ask` / `deny`), merged with agent-level and session-level permissions. See [Permissions](../permissions/index.md#global-permissions). |
| `hooks` | object | _unset_ | Global lifecycle hooks applied to every agent, additive with agent-config and CLI hooks. See [Global (user-level) hooks](../hooks/index.md#global-user-level-hooks). |
| `keybindings` | array | _unset_ | Remap TUI keyboard shortcuts. See [Custom Keybindings](../../features/tui/index.md#custom-keybindings) for the full list of actions and syntax. |
| `layout` | object | _unset_ | Sidebar position and section visibility. See [Layout Settings](#layout-settings) below. |

## Layout Settings

`layout` customizes the TUI's sidebar. The zero value (an omitted `layout:` block, or any field left out) is the default: sidebar on the right, every section visible, normal spacing.

| Field | Type | Default | Description |
| --- | --- | --- | --- |
| `sidebar_position` | string | `right` | `right`, `left`, `top`, or `bottom`. Left/right keep a full vertical sidebar; top/bottom render a compact horizontal band. |
| `section_spacing` | string | `normal` | `compact`, `normal`, or `relaxed` — the number of blank lines between sidebar sections. |
| `hide_session_path` | boolean | `false` | Hide the working-directory (session path) line, including its git branch. |
| `hide_usage` | boolean | `false` | Hide the token-usage section. |
| `hide_agents` | boolean | `false` | Hide the Agents section. |
| `hide_tools` | boolean | `false` | Hide the Tools section. |
| `hide_todos` | boolean | `false` | Hide the Todos section. |

```yaml
settings:
  layout:
    sidebar_position: left
    section_spacing: compact
    hide_usage: true
```

## Complete Example

```yaml
# ~/.config/cagent/config.yaml
settings:
  theme: auto
  theme_dark: dracula
  theme_light: default-light
  lean: false
  expand_thinking: false
  split_diff_view: true
  render_images: true
  hide_tool_results: false
  sound: true
  sound_threshold: 10
  snapshot: true
  cache_stable_prompts: true
  warn_on_cache_miss: true
  busy_send_mode: queue
  restore_tabs: true
  tab_title_max_length: 24
  layout:
    sidebar_position: right
    section_spacing: normal
  permissions:
    deny:
      - "shell:cmd=sudo*"
    allow:
      - "read_*"
  hooks:
    session_start:
      - type: command
        command: "~/.config/cagent/hooks/session-start.sh"
  keybindings:
    - action: "commands"
      keys: ["f2", "ctrl+k"]
```

## Precedence Rules

User settings are the **lowest-priority** source: they establish defaults, and anything more specific wins.

- **CLI flags over user settings — except plain boolean flags going from `true` to `false`.** Where a `docker agent run` flag mirrors a setting, passing the flag for a specific run takes precedence over the setting for that run only, and the flag never modifies the saved user config file. This holds cleanly for `--lean` / `lean` and `--theme` / `theme`, which track whether the flag was explicitly passed on the command line. `--yolo` / `YOLO` and `--hide-tool-results` / `hide_tool_results` don't: they're plain booleans with no "was this explicitly set" tracking, so passing `--yolo=false` or `--hide-tool-results=false` cannot turn a saved `YOLO: true` / `hide_tool_results: true` setting off for that run — the saved `true` wins and is reapplied on top of the flag. Passing the flag to turn either *on* (`--yolo`, `--hide-tool-results`) works as expected regardless of the saved setting.
- **Aliases sit between CLI flags and user settings.** An [alias](../../features/cli/index.md#docker-agent-alias) (`docker agent alias add ...`) can bundle its own `yolo`, `model`, `hide_tool_results`, and `sandbox` defaults; those apply when the corresponding flag was not explicitly passed, the same way user settings do, but are resolved after user settings so an alias's own choices take priority over your global defaults.
- **Permissions are merged, not overridden.** Global `settings.permissions` and an agent's own `permissions:` are combined into a single set of `deny` → `allow` → `ask` patterns before evaluation — a global deny always blocks, regardless of what the agent config allows. See [Merging Behavior](../permissions/index.md#merging-behavior).
- **Hooks are additive, not overridden.** For a given lifecycle event, hooks from the agent config, `settings.hooks`, `hooks.d/` drop-ins, and `--hook-*` CLI flags **all** run, in that order. Global hooks cannot be suppressed by an individual agent.
- **Everything else is a plain default.** Fields with no CLI or agent-config equivalent (`sound`, `sound_threshold`, `restore_tabs`, `tab_title_max_length`, `split_diff_view`, `render_images`, `cache_stable_prompts`, `warn_on_cache_miss`, `busy_send_mode`, `keybindings`, `layout`) only ever come from `settings:` (or the `/settings` dialog that writes it) — there is nothing to override them per run.
