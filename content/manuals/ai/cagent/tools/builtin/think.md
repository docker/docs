---
title: Think
description: Thought recording and reasoning tool
---

This toolset provides your agent with a scratchpad for reasoning and planning.
It allows the agent to record thoughts without performing actions or changing
data, making it useful for complex reasoning tasks.

## Usage

```yaml
toolsets:
  - type: think
```

## Tools

| Name    | Description                                    |
| ------- | ---------------------------------------------- |
| `think` | Record thoughts for reasoning and planning     |

### think

Use this tool to think about something. It will not obtain new information or
change the database, but will append the thought to a log. This is useful when
complex reasoning or cache memory is needed.

Args:

- `thought`: The thought to think about (required)

**Use Cases:**

- List specific rules that apply to the current request
- Check if all required information has been collected
- Verify that planned actions comply with policies
- Iterate over tool results for correctness
- Break down complex problems into steps
- Record intermediate reasoning steps

**Example:**

```json
{
  "thought": "The user wants to create a new feature. I need to: 1) Check existing code structure, 2) Identify dependencies, 3) Create new files, 4) Update configuration"
}
```

## Best Practices

- Use the think tool generously before taking actions
- Record your reasoning process step-by-step
- Use it to verify tool results and plan next steps
- Helpful for maintaining context during multi-step tasks
- Use it to check compliance with rules and policies before proceeding

## Output

The tool returns all accumulated thoughts in the session, allowing you to review
your reasoning process:

```
Thoughts:
First thought here
Second thought here
Third thought here
```
