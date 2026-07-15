---
title: "Terminal UI (TUI)"
description: "docker-agent's default interface is a rich, interactive terminal UI with file attachments, themes, session management, and more."
keywords: docker agent, ai agents, features, terminal ui (tui)
linkTitle: "Terminal UI"
weight: 10
canonical: https://docs.docker.com/ai/docker-agent/features/tui/
---

_docker-agent's default interface is a rich, interactive terminal UI with file attachments, themes, session management, and more._

![docker-agent TUI in action showing an interactive agent session](../../demo.gif)

## Launching the TUI

```bash
# Launch with a config
$ docker agent run agent.yaml

# Start with an initial message
$ docker agent run agent.yaml "Help me refactor this code"

# Auto-approve all tool calls
$ docker agent run agent.yaml --yolo

# Enable debug logging
$ docker agent run agent.yaml --debug

# Override the application name shown in the status bar and window title
$ docker agent run agent.yaml --app-name "My Project"

# Preselect a color theme
$ docker agent run agent.yaml --theme dracula

# Hide the sidebar (cannot be re-enabled via Ctrl+B)
$ docker agent run agent.yaml --sidebar=false

# Disable specific slash commands
$ docker agent run agent.yaml --disable-commands="/cost,/eval,/model"

# Open in read-only mode to review a past session without sending new messages
$ docker agent run agent.yaml --session -1 --session-read-only

# Use the lean TUI for this run
$ docker agent run agent.yaml --lean
```

### Lean TUI

The lean TUI uses a simplified terminal interface with minimal chrome. To make it the default for interactive runs, set `lean` in your user config:

```yaml
# ~/.config/cagent/config.yaml
settings:
  lean: true
```

Omit `lean` or set it to `false` to keep the full TUI as the default. You can still use `--lean` for a single run, or `--lean=false` to use the full TUI when `settings.lean` is enabled.

The lean TUI supports **steering**: messages submitted while the agent is running are queued and delivered to the active session. Pending steering messages appear with muted styling at the end of the live stream so you can see what will be sent next.

## Slash Commands

Type `/` during a session to see available commands, or press <kbd>Ctrl</kbd>+<kbd>K</kbd> for the command palette:

| Command            | Description                                                                          |
| ------------------ | ------------------------------------------------------------------------------------ |
| `/new`             | Start a new conversation                                                             |
| `/clear`           | Clear the current conversation (keep session, drop messages)                         |
| `/compact`         | Summarize and compact the conversation history                                       |
| `/fork`            | Fork the current session into a new branch                                           |
| `/copy`            | Copy the entire conversation to clipboard                                            |
| `/copy-last`       | Copy only the last assistant message to clipboard                                    |
| `/undo`            | Restore file changes from the latest snapshot (only when snapshots are enabled)      |
| `/snapshots`       | List captured snapshots (only when snapshots are enabled)                            |
| `/export`          | Export the session as HTML                                                           |
| `/sessions`        | Browse and load past sessions                                                        |
| `/model`           | Change the model for the current agent                                               |
| `/effort`          | Set the current model's reasoning-effort level (`/effort <none\|minimal\|low\|medium\|high\|xhigh\|max>`, or `/effort` alone to pick from the supported levels; reasoning models only) |
| `/custom`          | Customize the TUI layout: sidebar position, section spacing, and visible sidebar sections |
| `/theme`           | Change the color theme                                                               |
| `/yolo`            | Toggle automatic tool call approval                                                  |
| `/title`           | Set or regenerate session title                                                      |
| `/attach`          | Attach a file to your message                                                        |
| `/shell`           | Open a shell                                                                         |
| `/star`            | Star/unstar the current session                                                      |
| `/context`         | Show a context-window breakdown: estimated tokens per category (system prompt, tool definitions, prompt files, messages, tool results, compaction summary), plus a per-file inventory of attached files and prompt files. Select an attached file with the arrow keys and press <kbd>d</kbd> to drop it |
| `/drop`            | Remove an attached file from the session context (`/drop <path>`, or `/drop` alone to review and drop from the `/context` dialog)  |
| `/cost`            | Show cost breakdown for this session                                                 |
| `/eval`            | Create an evaluation report                                                          |
| `/pause`           | Pause/resume the runtime loop. While the agent is mid-request, the resize handle shows "Pausing…" until the in-flight request completes; once the loop is blocked the indicator changes to "⏸ Paused". Run `/pause` again to resume. |
| `/tools`           | Show every toolset (with lifecycle state) and the tools they expose                  |
| `/skills`          | List skills available to the current agent                                           |
| `/toolset-restart` | Force a supervisor-driven reconnect of the named toolset (`/toolset-restart <name>`) |
| `/permissions`     | Inspect and edit tool permission rules                                               |
| `/split-diff`      | Toggle split-diff view for file edits                                                |
| `/speak`           | Voice input via system speech-to-text (macOS only)                                   |
| `/exit`            | Exit the application (aliases: `/quit`, `/q`)                                        |

Slash commands (both built-in and named) execute immediately when entered. Regular chat messages are queued and processed in order. This means you can invoke a slash command to interrupt or direct the agent even while it is mid-response.

### Agents Panel

The sidebar's **Agents** section lists every agent in the team. The current agent is shown as a focus **card** (rendered in place at its position in the list) with its name, a wrapped description, its full `provider/model`, and a thinking line. Every other agent is shown as a compact **two-line row** — line 1 is the shortcut/spinner, the agent name (in its accent color), and a right-aligned thinking **gauge**; line 2 is the indented full `provider/model` — so a large team stays scannable while still showing each model. Once an agent has run (in the main session, as a delegated sub-agent, or as a background agent task), line 2 also carries its latest **context usage** as a right-aligned percentage of its context window, so per-agent context accounting is visible at a glance across the whole team. Agents are separated by a blank line so the two-line rows stay visually distinct. The effort **gauge** is the only visual language for thinking; the focus card and the Agent Inspector spell out the exact level alongside it. Left-click any agent to switch to it.

#### Agent inspector

Open a read-only **Agent Inspector** to inspect any agent's full configuration combined with its live state. The instruction/system prompt is deliberately omitted; everything else the agent declares is shown:

- **Right-click any agent** (card or row) to open the inspector without switching to it.
- **<kbd>Ctrl</kbd>+left-click any agent** does the same — a fallback for terminals that don't forward right-clicks.
- **Left-click** always switches to the agent.

The title is rendered in the agent's accent color. Sections appear in this order, and any empty section is omitted:

- **Description** — the agent's wrapped description.
- **Live state** — a `● current agent` line when the inspected agent is the one currently running.
- **Model / Fallback / Thinking** — the `provider/model`, any fallback models, and the gauge + value thinking line (omitted for models with no selectable thinking, e.g. harness-backed agents).
- **Context** — the agent's latest known context usage, e.g. `Context: 12.8K of 128.0K tokens (10%)` (a bare token count when the context limit is unknown; omitted until the agent has run). Sub-agent and background-agent runs are accounted for.
- **Sub-agents (N) / Handoffs (N) / Skills (N)** — compact, inline, comma-separated lists wrapped to the dialog width.
- **Limits** — the configured per-agent limits that are set, e.g. `Limits: max-iter 50 · history 40 · max-tool-calls 5`.
- **Options** — the enabled option flags, e.g. `Options: add-date · add-environment-info · redact-secrets`.
- **Toolsets (N)** — one line per toolset with a status marker, its name, kind, and tool count, followed by the indented tool names.
- **Commands (N)** — the slash commands the agent defines, each with its description.

Each toolset carries a single-width status marker reflecting its **live** lifecycle: `●` started (serving), `○` stopped (not yet started), or `⚠` error. The tools listed under a toolset are the **live** tool names when it has started; for a toolset that has not started, the inspector instead shows its declared `tools:` allow-list prefixed with `declared:` (and shows nothing when the toolset declares no allow-list and therefore serves every tool). This lets you see both what an agent is configured with and what is actually running, even before the agent has been used.

The dialog scrolls when the content is long; press <kbd>Esc</kbd> to close it. Remote runtimes (which hold no local team config) degrade gracefully — the config-derived sections are simply omitted.

Model identifiers on line 2 are truncated **from the left** (e.g. `…claude-sonnet-4-6`) only when they overflow, so the informative tail (variant/version) is preserved. As the sidebar narrows the model keeps its own line, and near the minimum width line 1's gauge collapses to a single cell to keep the name readable.

The thinking state of each model is shown with a gauge + value on the card and a gauge or badge on the row (no `✻` glyph):

| Model state            | Card line                      | Row badge              |
| ---------------------- | ------------------------------ | ---------------------- |
| Effort level           | `thinking ▰▰▰▰▱▱ high`         | `▰▰▰▰▱▱` (effort gauge)  |
| Adaptive budget        | `thinking auto adaptive`       | `auto`                 |
| Token budget           | `thinking ◉ 8.2K tokens`        | `◉ 8.2K`               |
| Disabled (capable)     | `thinking ▱▱▱▱▱▱ off` (dimmed)  | `▱▱▱▱▱▱` (empty gauge)  |
| Not reasoning-capable  | _(omitted)_                    | _(omitted)_            |

The **effort gauge** is a fixed-width six-cell indicator (`▰` filled, `▱` empty) so the badge column stays aligned. It maps the six selectable levels one-to-one onto filled-cell counts — `minimal` → `▰▱▱▱▱▱`, `low` → `▰▰▱▱▱▱`, `medium` → `▰▰▰▱▱▱`, `high` → `▰▰▰▰▱▱`, `xhigh` → `▰▰▰▰▰▱`, `max` → `▰▰▰▰▰▰` — so the cell count alone is lossless, with a low→high color ramp as a secondary cue. A capable-but-disabled model shows a dim empty gauge (`▱▱▱▱▱▱` `off`), adaptive budgets show `auto`, and token budgets keep `◉ <count>`. The same gauge + value renders on the focus card, the Agent Inspector, and the row.

Harness-backed agents (e.g. `claude-code`) show the harness type as their model and no thinking gauge. Press **Shift+Tab** to cycle the current model's thinking-effort level; a `✻ Thinking: <level>` toast confirms the change (useful when the sidebar is hidden).

### Context-Usage Gauge

The context percentage shown in the sidebar token-usage section, and the fill bar in the lean TUI status line, both color-escalate as the active session approaches the auto-compaction threshold:

| State   | Color  | Trigger |
| ------- | ------ | ------- |
| Normal  | (default) | Usage below 75% of the compaction threshold |
| Warning | Orange | Usage at or above 75% of the compaction threshold |
| Critical | Red   | Usage at or above 95% of the compaction threshold |

While a compaction is running the percentage is replaced by a **"compacting…"** indicator; token counts remain visible in the lean TUI status line.

The thresholds are proportional to the agent's configured `compaction_threshold` (default `0.9`), so a custom value keeps a predictable visual runway. See [Compaction Threshold](../../configuration/models/index.md#delegating-session-compaction) for configuration details.

### Thinking and Tool Details

Reasoning/thinking blocks are collapsed by default and carry a `Thinking` header badge. When collapsed, the TUI shows a short preview and compact tool summaries. Expand a block to see the full thinking content and the real tool renderers, including detailed tool output such as file edit diffs.

To start new sessions with thinking/tool blocks expanded by default, set `expand_thinking` in your user config:

```yaml
# ~/.config/cagent/config.yaml
settings:
  expand_thinking: true
```

Set it to `false` or omit it to keep the default collapsed behavior.

### Snapshots, `/undo`, and `/snapshots`

Enable shadow-git snapshots globally in `~/.config/cagent/config.yaml`:

```yaml
settings:
  snapshot: true
```

When enabled, docker-agent records filesystem snapshots at turn boundaries. The TUI exposes two slash commands that operate on those snapshots:

- **`/undo`** restores files from the most recent snapshot (one step back).
- **`/snapshots`** opens a dialog showing how many snapshots have been captured and the number of files in each one. Use <kbd>↑</kbd>/<kbd>↓</kbd> (or <kbd>j</kbd>/<kbd>k</kbd>) to highlight an entry, then press <kbd>r</kbd> to reset the workspace to that point. Pick `<original>` to revert every snapshot and bring the workspace back to its pre-agent state. <kbd>Esc</kbd> closes the dialog without changing anything.

Neither command removes messages from the session transcript — they only touch files on disk. Both commands (and the matching command-palette entries) are hidden when snapshots are turned off. Omit `snapshot` or set it to `false` to leave automatic snapshots off; agents can still configure snapshot hooks manually.

See [Snapshots](../snapshots/index.md) for how the shadow-git machinery works and how to wire it per-agent.

## File Attachments

Attach file contents to your messages using the `@` trigger:

1. Type `@` to open the file completion menu
2. Start typing to filter files (respects `.gitignore`)
3. Select a file to insert the reference

```bash
# In the chat input:
Explain what the code in @pkg/agent/agent.go does
```

The agent receives the full file contents in a structured `<attachments>` block, while the UI shows just the reference.

Attached files are also recorded on the session so sub-agents spawned by task transfer can read them. To review what is attached, open `/context`: the dialog lists every attached file (and resolved prompt file) with a per-file token estimate. Use <kbd>↑</kbd>/<kbd>↓</kbd> to select an attached file and press <kbd>d</kbd> (or <kbd>x</kbd>/<kbd>Del</kbd>) to drop it, or run `/drop <path>` directly. Dropping stops sharing the file with sub-agents and skills; content already inlined in earlier messages stays in the conversation until compaction, and the file can always be re-attached with `@` or `/attach`.

## Runtime Model Switching

Change the AI model during a session with `/model` or <kbd>Ctrl</kbd>+<kbd>M</kbd>:

1. Press <kbd>Ctrl</kbd>+<kbd>M</kbd> or type `/model`
2. Select from config models or type a custom `provider/model`
3. The model switch is saved with the session and restored on reload

When a models gateway is configured (`--models-gateway`) and it exposes an OpenAI-style `/v1/models` endpoint, the picker lists the models actually served by the gateway (merged with the models defined in the agent config). When the gateway doesn't expose `/v1/models`, the picker falls back to the regular catalog.

> [!TIP]
> Use model switching to try a more capable model for complex tasks, or a cheaper one for simple queries — without modifying your YAML config.

## Editable Messages

Edit any previous user message to branch the conversation. Click on a past message to modify it — the agent will re-process from that point, while the original session history is preserved. This is great for exploring alternative approaches without losing your work.

## Error Recovery

When an agent turn fails (fatal model error, hook block, loop detection, tool-setup failure), the TUI displays the error in the message stream and persists it to the session store. Errors survive a reload and are shown exactly where they occurred, making them visible in shared or remote sessions.

Each error message includes a clickable **↻ retry** button. Clicking it resumes the conversation from the point of failure — without retyping your last message. This lets you recover from transient failures (rate limits, network blips, model API errors) in one click.

## Session Management

docker-agent automatically saves your sessions. Use `/sessions` to browse past conversations:

- **Browse** past sessions with search and filtering
- **Workspace grouping**: sessions started in the current directory are listed first under "This workspace", everything else under "Other locations" with its originating directory shown next to each entry; press <kbd>Ctrl</kbd>+<kbd>G</kbd> in the browser to cycle between all, current-directory only, and other-directory views. Restoring a session reopens it in its original directory, so the label always matches where a restore will land
- **Star** important sessions with `/star`
- **Branch** conversations by editing any previous user message — preserving the original session history
- **Resume** sessions with `docker agent run config.yaml --session <id>`
- **Relative refs**: `--session -1` for the last session, `-2` for the one before

### Session Title Editing

Customize session titles to make them more meaningful and easier to find. By default, docker-agent auto-generates titles based on your first message, but you can override or regenerate them at any time.

**Using the `/title` command:**

```bash
/title                     # Regenerate title using AI (based on recent messages)
/title My Custom Title     # Set a specific title
```

**Using the sidebar:**

1. Click the pencil icon (✎) next to the session title in the sidebar
2. Type your new title
3. Press <kbd>Enter</kbd> to save, or <kbd>Escape</kbd> to cancel

> [!NOTE]
> Manually set titles are preserved and won’t be overwritten by auto-generation. Title changes are persisted immediately to the session.

## Keyboard Shortcuts

| Shortcut   | Action                                          |
| ---------- | ----------------------------------------------- |
| Ctrl+K     | Open command palette                            |
| Ctrl+M     | Switch model                                    |
| Ctrl+R     | Reverse history search (search previous inputs) |
| Ctrl+G     | Cancel reverse history search                   |
| Ctrl+S     | Cycle to next agent in the team                 |
| Shift+Tab  | Cycle the current model's thinking-effort level (shows a `✻ Thinking: <level>` toast) |
| Ctrl+1 – 9 | Switch directly to agent _N_ in the team list   |
| Ctrl+T     | Open a new tab (additional agent session)       |
| Ctrl+W     | Close the current tab                           |
| Ctrl+N     | Next tab                                        |
| Ctrl+P     | Previous tab                                    |
| Ctrl+B     | Toggle the sidebar (full-UI mode only; disabled when --sidebar=false) |
| Ctrl+Y     | Toggle YOLO mode (auto-approve tool calls)      |
| Ctrl+O     | Toggle hide tool results                        |
| Ctrl+Z     | Suspend TUI to background (resume with `fg`)    |
| Ctrl+X     | Clear queued messages                           |
| Escape     | Cancel current operation                        |
| Enter      | Send message (or newline with Shift+Enter)      |
| Up/Down    | Navigate message history                        |

Press <kbd>Ctrl</kbd>+<kbd>H</kbd> to view the complete list of all available keyboard shortcuts.

### Custom Keybindings

You can remap the shortcuts above by adding a `keybindings` list to the `settings` block of your `~/.config/cagent/config.yaml`. Each entry maps an action to one or more key combinations in [Bubbles key format](https://github.com/charmbracelet/bubbles) (for example `ctrl+q`, `alt+enter`, `f2`). Unlisted actions keep their defaults.

This is the recommended way to replace the `Ctrl+J` newline fallback, which conflicts with common editor/terminal shortcuts (for example inside VS Code).

```yaml
settings:
  keybindings:
    # Insert a newline with Alt+Enter instead of Ctrl+J. Shift+Enter still
    # works automatically on terminals that report it.
    - action: "editor_newline"
      keys: ["alt+enter"]
    # Allow several keys for one action.
    - action: "commands"
      keys: ["f2", "ctrl+k"]
    - action: "quit"
      keys: ["ctrl+q"]
```

**Valid actions:**

| Action                     | Default      | Description                            |
| -------------------------- | ------------ | -------------------------------------- |
| `editor_send`              | `enter`      | Send the current message               |
| `editor_newline`           | `ctrl+j`     | Insert a newline in the input          |
| `quit`                     | `ctrl+c`     | Quit (opens the exit confirmation)     |
| `switch_focus`             | `tab`        | Switch focus between panels            |
| `commands`                 | `ctrl+k`     | Open the command palette               |
| `help`                     | `ctrl+h`     | Show the help dialog                   |
| `toggle_yolo`              | `ctrl+y`     | Toggle YOLO mode                       |
| `toggle_hide_tool_results` | `ctrl+o`     | Toggle hiding tool results             |
| `cycle_agent`              | `ctrl+s`     | Cycle to the next agent                |
| `model_picker`             | `ctrl+m`     | Open the model picker                  |
| `clear_queue`              | `ctrl+x`     | Clear queued messages                  |
| `suspend`                  | `ctrl+z`     | Suspend the TUI                        |
| `toggle_sidebar`           | `ctrl+b`     | Toggle the sidebar                     |
| `edit_external`            | `ctrl+g`     | Edit input in an external editor       |
| `history_search`           | `ctrl+r`     | Incremental history search             |

`Shift+Enter` for newline is detected from your terminal's capabilities and is always available where supported, independent of `editor_newline`.

Invalid entries are ignored with a warning (visible with `--debug`) so a bad config never breaks the TUI: unknown actions, empty or malformed keys, and keys that would collide with another action are dropped while every other binding keeps working.

## History Search

Press <kbd>Ctrl</kbd>+<kbd>R</kbd> to enter incremental history search mode. Start typing to filter through your previous inputs. Press <kbd>Enter</kbd> to select a match, or <kbd>Escape</kbd> to cancel.

## Layout Customization

Run `/custom` to open the layout dialog. It shows a live schematic preview of the resulting layout and applies every change immediately to the UI behind the dialog:

- **Sidebar position**: `Right` (default), `Left`, `Top`, or `Bottom`. Left/right keep the full vertical sidebar next to the chat; top/bottom render it as a compact horizontal band above or below the chat (session title, working directory, usage, plus a one-line summary of the current agent, tools, and todos).
- **Section spacing**: `Compact`, `Normal` (default), or `Relaxed`, the number of blank lines between the sidebar sections (1, 2, or 3).
- **Sidebar sections**: toggle the visibility of the **Token usage**, **Agents**, **Tools**, and **Todos** sections. The session block (title and working directory) is always shown.

Press <kbd>Enter</kbd> to apply and persist, or <kbd>Escape</kbd> to cancel and restore the previous layout. The settings are saved globally in `~/.config/cagent/config.yaml`:

```yaml
# ~/.config/cagent/config.yaml
settings:
  layout:
    sidebar_position: left # right (default), left, top, bottom
    section_spacing: compact # normal (default), compact, relaxed
    hide_usage: true
    hide_agents: false
    hide_tools: false
    hide_todos: false
```

## Theming

Customize the TUI appearance with built-in or custom themes:

```bash
# Switch themes interactively
/theme
```

### Built-in Themes

`default`, `default-light`, `catppuccin-latte`, `catppuccin-mocha`, `dracula`, `gruvbox-dark`, `gruvbox-light`, `nord`, `one-dark`, `solarized-dark`, `tokyo-night`

### Auto Theme (match the terminal)

The special theme `auto` follows the terminal's light/dark background instead of naming a fixed theme. Select **Auto (match terminal)** in the `/theme` picker, pass `--theme auto`, or set it in your user config:

```yaml
settings:
  theme: auto
  theme_dark: default # optional, theme used on dark backgrounds (default: default)
  theme_light: default-light # optional, theme used on light backgrounds (default: default-light)
```

At startup the terminal background is queried (OSC 11) to pick the dark or light theme of the pair; non-interactive runs (pipes, CI) fall back to the dark theme. In terminals that report appearance changes (DEC mode 2031 — Ghostty, kitty, contour, …), flipping the OS or terminal appearance while docker-agent is running switches the theme live. Terminals without that mode re-sync when the window regains focus.

### Custom Themes

Create theme files in `~/.cagent/themes/` as YAML. Theme files are **partial overrides** — you only need to specify the colors you want to change. Any omitted keys fall back to the built-in default theme values.

```yaml
# ~/.cagent/themes/my-theme.yaml
name: "My Custom Theme"

colors:
  # Backgrounds
  background: "#1a1a2e"
  background_alt: "#16213e"

  # Text colors
  text_bright: "#ffffff"
  text_primary: "#e8e8e8"
  text_secondary: "#b0b0b0"
  text_muted: "#707070"

  # Accent colors
  accent: "#4fc3f7"
  brand: "#1d96f3"

  # Status colors
  success: "#4caf50"
  error: "#f44336"
  warning: "#ff9800"
  info: "#00bcd4"

# Optional: Customize syntax highlighting colors
chroma:
  comment: "#6a9955"
  keyword: "#569cd6"
  literal_string: "#ce9178"

# Optional: Customize markdown rendering colors
markdown:
  heading: "#4fc3f7"
  link: "#569cd6"
  code: "#ce9178"
```

### Applying Themes

**In user config** (`~/.config/cagent/config.yaml`):

```yaml
settings:
  theme: my-theme # References ~/.cagent/themes/my-theme.yaml
```

**At launch:** Pass `--theme <name>` to `docker agent run` to preselect a theme for that session. This overrides `settings.theme` in your config but is not saved. Invalid theme names print an error at startup listing the available options. Has no effect in `--exec` mode. `--theme auto` enables the [auto theme](#auto-theme-match-the-terminal) for the session.

**At runtime:** Use the `/theme` command to open the theme picker and select from available themes. Your selection is saved globally in `~/.config/cagent/config.yaml` under `settings.theme` and persists across sessions.

> [!TIP]
> **Hot Reload**
>
> Custom themes auto-reload when you save changes to the file — no restart needed. This makes it easy to tweak colors in real-time.

> [!WARNING]
> **Partial overrides**
>
> All user themes are applied on top of the `default` theme. If you want to customize a built-in theme (e.g., `dracula`), copy its full YAML from the [built-in themes on GitHub](https://github.com/docker/docker-agent/tree/main/pkg/tui/styles/themes) into `~/.cagent/themes/` and edit the copy. Otherwise, omitted values will use `default` colors, not the original theme's colors.

## Tool Permissions

When an agent calls a tool, docker-agent shows a confirmation dialog by default. You can:

- **Approve once** — Allow this specific call
- **Always allow** — Permanently approve this tool/command for the session
- **Deny** — Reject the tool call

**Granular permissions:** The permission system supports pattern-based matching. When you “Always allow” a specific tool command, only that exact pattern is auto-approved — other commands from the same tool still require confirmation. This lets you auto-approve safe, read-only operations while maintaining control over destructive ones.

> [!TIP]
> **YOLO mode**
>
> Use `--yolo` or the `/yolo` command to auto-approve all tool calls. You can also toggle this mid-session. For aliases, set `--yolo` when creating the alias: `docker agent alias add fast agentcatalog/coder --yolo`.

## Notifications

The TUI displays transient notification banners for agent warnings, errors, and other runtime events. Notifications auto-dismiss after a short delay unless the mouse is hovering over them — hovering pauses the timer so you have time to read the message.

| Interaction | Behaviour |
| ----------- | --------- |
| Hover       | Pauses auto-dismiss; the notification stays visible until the mouse moves away |
| Click       | Copies the notification text to the clipboard |
| × (close)   | Dismisses immediately; the glyph turns red when hovered |

Hint text in the top-left corner of the notification border shows the available actions at a glance.
