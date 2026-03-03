# docker model launch

<!---MARKER_GEN_START-->
Launch an app configured to use Docker Model Runner.

Without arguments, lists all supported apps.

Supported apps: anythingllm, claude, codex, openclaw, opencode, openwebui

Examples:
  docker model launch
  docker model launch opencode
  docker model launch claude -- --help
  docker model launch openwebui --port 3000
  docker model launch claude --config

### Options

| Name        | Type     | Default | Description                                     |
|:------------|:---------|:--------|:------------------------------------------------|
| `--config`  | `bool`   |         | Print configuration without launching           |
| `--detach`  | `bool`   |         | Run containerized app in background             |
| `--dry-run` | `bool`   |         | Print what would be executed without running it |
| `--image`   | `string` |         | Override container image for containerized apps |
| `--model`   | `string` |         | Model to use (for opencode)                     |
| `--port`    | `int`    | `0`     | Host port to expose (web UIs)                   |


<!---MARKER_GEN_END-->

