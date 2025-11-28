---
title: Shell
description: Execute shell commands in the user's environment
---

This toolset allows your agent to execute shell commands in the user's default
shell environment. Commands run with full access to environment variables and
can be executed in any working directory.

## Usage

```yaml
toolsets:
  - type: shell
```

## Tools

| Name    | Description                                      |
| ------- | ------------------------------------------------ |
| `shell` | Execute shell commands in the user's environment |

### shell

Executes shell commands in the user's default shell. On Windows, PowerShell
(pwsh/powershell) is used when available; otherwise, cmd.exe is used. On
Unix-like systems, the `$SHELL` environment variable is used, or `/bin/sh` as a
fallback.

Args:

- `cmd`: The shell command to execute (required)
- `cwd`: Working directory to execute the command in (required, use "." for
  current directory)
- `timeout`: Command execution timeout in seconds, default is 30 (optional)

**Features:**

- Supports pipes, redirections, and complex shell operations
- Each command runs in a fresh shell session (no state persists)
- Automatic timeout protection to prevent hanging commands
- Full access to environment variables
- Support for heredocs and multi-line scripts

**Examples:**

Basic command:
```json
{
  "cmd": "ls -la",
  "cwd": "."
}
```

Long-running command with custom timeout:
```json
{
  "cmd": "npm run build",
  "cwd": ".",
  "timeout": 120
}
```

Using pipes:
```json
{
  "cmd": "cat package.json | jq '.dependencies'",
  "cwd": "frontend"
}
```

Complex script with heredoc:
```json
{
  "cmd": "cat << 'EOF' | ${SHELL}\necho 'Hello'\necho 'World'\nEOF",
  "cwd": "."
}
```

## Best Practices

- Use the `cwd` parameter for directory-specific commands
- Quote arguments containing spaces or special characters
- Use the `timeout` parameter for long-running operations (builds, tests, etc.)
- Prefer heredocs over writing temporary script files
- Leverage this tool for batch file operations

## Error Handling

- Commands with non-zero exit codes return error information along with any
  output
- Commands that exceed their timeout are automatically terminated
- Output includes both stdout and stderr combined
