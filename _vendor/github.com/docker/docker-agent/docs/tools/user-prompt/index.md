---
title: "User Prompt Tool"
description: "Ask the user questions and collect interactive input during agent execution."
keywords: docker agent, ai agents, tools, toolsets, user prompt tool
linkTitle: "User Prompt"
weight: 190
canonical: https://docs.docker.com/ai/docker-agent/tools/user-prompt/
---

_Ask the user questions and collect interactive input during agent execution._

## Overview

The user prompt tool allows agents to ask questions and collect input from users during execution. This enables interactive workflows where the agent needs clarification, confirmation, or additional information before proceeding.

> [!NOTE]
> **When to Use**
>
> - When the agent needs clarification before proceeding
> - Collecting credentials or configuration values
> - Presenting choices and getting user decisions
> - Confirming destructive or important actions

## Configuration

```yaml
agents:
  assistant:
    model: openai/gpt-4o
    description: Interactive assistant
    instruction: |
      You are a helpful assistant. When you need information
      from the user, use the user_prompt tool to ask them.
    toolsets:
      - type: user_prompt
      - type: filesystem
      - type: shell
```

## Tool Interface

The `user_prompt` tool takes these parameters:

| Parameter | Type   | Required | Description                                                                                        |
| --------- | ------ | -------- | -------------------------------------------------------------------------------------------------- |
| `message` | string | ✓        | The question or prompt to display.                                                                 |
| `title`   | string | ✗        | Optional title for the dialog window in the TUI. Defaults to `"Question"` when not provided.       |
| `schema`  | object | ✗        | JSON Schema defining the expected response structure (object or primitive).                        |

## Response Format

The tool returns a JSON response:

```json
{
  "action": "accept",
  "content": {
    "field1": "user value",
    "field2": true
  }
}
```

### Action Values

| Action    | Meaning                                    |
| --------- | ------------------------------------------ |
| `accept`  | User provided a response (check `content`) |
| `decline` | User declined to answer                    |
| `cancel`  | User cancelled the prompt                  |

## Schema Examples

### Simple String Input

```json
{
  "type": "string",
  "title": "API Key",
  "description": "Enter your API key"
}
```

### Multiple Choice

```json
{
  "type": "string",
  "enum": ["development", "staging", "production"],
  "title": "Environment",
  "description": "Select the target environment"
}
```

### Boolean Confirmation

```json
{
  "type": "boolean",
  "title": "Confirm",
  "description": "Are you sure you want to proceed?"
}
```

### Object with Multiple Fields

```json
{
  "type": "object",
  "properties": {
    "username": {
      "type": "string",
      "description": "Your username"
    },
    "password": {
      "type": "string",
      "description": "Your password"
    },
    "remember": {
      "type": "boolean",
      "description": "Remember credentials"
    }
  },
  "required": ["username", "password"]
}
```

### Number Input

```json
{
  "type": "integer",
  "title": "Port Number",
  "description": "Enter the port number (1024-65535)",
  "minimum": 1024,
  "maximum": 65535
}
```

## Example Usage

Here's how an agent might use the user prompt tool:

```text
Agent: I need to deploy this application. Let me ask which environment to target.

[Calls user_prompt with message: "Which environment should I deploy to?"
 and schema with enum: ["development", "staging", "production"]]

User selects: "staging"

Agent: Great, I'll deploy to staging. Let me confirm this action.

[Calls user_prompt with message: "Deploy to staging? This will replace the current version."
 and schema with type: "boolean"]

User confirms: true

Agent: Deploying to staging...
```

## UI Presentation

How the prompt appears depends on the interface:

- **TUI**: Displays an interactive dialog with appropriate input controls
- **CLI (exec mode)**: Prints the prompt and reads from stdin
- **API/MCP**: Returns an elicitation request to the client

> [!TIP]
> **Best Practice**
>
> Provide clear, concise messages. Include context about why you're asking and what the information will be used for. Use schemas with descriptions to guide users on expected input format.

## Handling Responses

The agent should handle all possible actions:

- **accept**: Process the `content` and continue
- **decline**: Acknowledge and try an alternative approach or explain what's needed
- **cancel**: Stop the current operation gracefully

> [!WARNING]
> **Context Requirement**
>
> The user prompt tool requires an elicitation handler to be configured. It works in the TUI and CLI modes but may not be available in all contexts (e.g., some MCP client configurations).
