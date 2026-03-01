# docker model skills

<!---MARKER_GEN_START-->
Install Docker Model Runner skills for AI coding assistants.

Skills are configuration files that help AI coding assistants understand
how to use Docker Model Runner effectively for local model inference.

Supported targets:
  --codex     Install to ~/.codex/skills (OpenAI Codex CLI)
  --claude    Install to ~/.claude/skills (Claude Code)
  --opencode  Install to ~/.config/opencode/skills (OpenCode)
  --dest      Install to a custom directory

Example:
  docker model skills --claude
  docker model skills --codex --claude
  docker model skills --dest /path/to/skills

### Options

| Name            | Type     | Default | Description                                             |
|:----------------|:---------|:--------|:--------------------------------------------------------|
| `--claude`      | `bool`   |         | Install skills for Claude Code (~/.claude/skills)       |
| `--codex`       | `bool`   |         | Install skills for OpenAI Codex CLI (~/.codex/skills)   |
| `--dest`        | `string` |         | Install skills to a custom directory                    |
| `-f`, `--force` | `bool`   |         | Overwrite existing skills without prompting             |
| `--opencode`    | `bool`   |         | Install skills for OpenCode (~/.config/opencode/skills) |


<!---MARKER_GEN_END-->

