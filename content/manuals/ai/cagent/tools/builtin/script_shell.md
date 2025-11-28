---
title: Script Shell
description: Define custom shell command as tools
---

This toolset allows you to define custom shell command tools with typed
parameters in your agent configuration. This enables you to create reusable,
parameterized shell commands as first-class tools for your agent.

## Usage

```yaml
toolsets:
  - type: script_shell
    tools:
      deploy:
        cmd: "./deploy.sh"
        description: "Deploy the application to production"
        workingDir: "./scripts"
        args:
          environment:
            type: string
            description: "Target environment (staging/production)"
          version:
            type: string
            description: "Version to deploy"
        required:
          - environment
          - version
      
      run_tests:
        cmd: "go test -v -race ./..."
        description: "Run Go tests with race detection"
        workingDir: "."
        args:
          package:
            type: string
            description: "Specific package to test (optional)"
```

## Configuration

Each custom tool is defined with the following properties:

- `cmd`: The shell command to execute (required)
- `description`: Human-readable description of what the tool does (optional,
  defaults to showing the command)
- `workingDir`: Working directory to execute the command in (optional)
- `args`: Object defining typed parameters that can be passed to the command
  (optional)
- `required`: Array of required parameter names (optional, defaults to all args
  being required)

## Parameters

Parameters defined in `args` are passed to the command as environment variables.
Each parameter can specify:

- `type`: The parameter type (string, number, boolean, etc.)
- `description`: Description of what the parameter is for

## Examples

### Simple command without parameters

```yaml
toolsets:
  - type: script_shell
    tools:
      build:
        cmd: "npm run build"
        description: "Build the frontend application"
        workingDir: "./frontend"
```

### Command with required parameters

```yaml
toolsets:
  - type: script_shell
    tools:
      create_migration:
        cmd: "migrate create -ext sql -dir ./migrations -seq $name"
        description: "Create a new database migration"
        args:
          name:
            type: string
            description: "Name of the migration"
        required:
          - name
```

### Command with optional parameters

```yaml
toolsets:
  - type: script_shell
    tools:
      run_benchmark:
        cmd: |
          if [ -n "$package" ]; then
            go test -bench=. -benchmem $package
          else
            go test -bench=. -benchmem ./...
          fi
        description: "Run Go benchmarks"
        args:
          package:
            type: string
            description: "Specific package to benchmark (optional)"
        required: []
```

## How It Works

1. Parameters are passed as environment variables to the shell command
2. Commands execute in the specified `workingDir` or current directory
3. The command runs in the user's default shell (`$SHELL` on Unix, or `/bin/sh`)
4. stdout and stderr are combined and returned as the tool result

## Best Practices

- **Use descriptive names** - Tool names should clearly indicate their purpose
- **Document parameters** - Provide clear descriptions for all parameters
- **Set working directories** - Use `workingDir` to ensure commands run in the
  correct context
- **Handle optional parameters** - Use shell conditionals when parameters are
  optional
- **Keep commands focused** - Each tool should do one thing well
- **Use shell scripts for complex logic** - For multi-step operations, call a
  shell script rather than inlining complex commands

## Use Cases

- Deployment automation
- Running tests with specific configurations
- Database migrations
- Code generation
- Build system integration
- CI/CD operations
- Custom project-specific workflows
