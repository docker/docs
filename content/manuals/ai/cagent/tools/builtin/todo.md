---
title: Todo
description: Task tracking for your agent
---

This toolset provides your agent with task tracking capabilities. It allows the
agent to create, update, and list todo items to maintain progress through
complex multi-step tasks.

## Usage

```yaml
toolsets:
  - type: todo
```

## Tools

| Name            | Description                                    |
| --------------- | ---------------------------------------------- |
| `create_todo`   | Create a new todo item with a description      |
| `create_todos`  | Create multiple todo items at once             |
| `update_todo`   | Update the status of a todo item               |
| `list_todos`    | List all current todos with their status       |

### create_todo

Creates a single new todo item with a unique ID and sets its initial status to
"pending".

Args:

- `description`: Description of the todo item (required)

**Example:**

```json
{
  "description": "Implement user authentication module"
}
```

### create_todos

Creates multiple todo items at once. Useful for breaking down a complex task
into multiple steps at the beginning.

Args:

- `descriptions`: Array of todo item descriptions (required)

**Example:**

```json
{
  "descriptions": [
    "Read existing code structure",
    "Design new feature architecture",
    "Implement core functionality",
    "Write tests",
    "Update documentation"
  ]
}
```

### update_todo

Updates the status of an existing todo item. Valid statuses are: `pending`,
`in-progress`, and `completed`.

Args:

- `id`: ID of the todo item (required)
- `status`: New status - `pending`, `in-progress`, or `completed` (required)

**Example:**

```json
{
  "id": "todo_1",
  "status": "completed"
}
```

### list_todos

Lists all current todos with their ID, description, and status. No arguments
required.

**Response format:**

```
Current todos:
- [todo_1] Implement user authentication module (Status: completed)
- [todo_2] Design new feature architecture (Status: in-progress)
- [todo_3] Write tests (Status: pending)
```

## Best Practices

- **Always create todos before starting complex tasks** - Break down work into
  manageable steps
- **Use list_todos frequently** - Check remaining work before responding to
  users
- **Update status regularly** - Mark todos as completed to track progress
- **Never skip steps** - Ensure all todos are addressed before considering a
  task complete
- **Use create_todos for batch creation** - More efficient than multiple
  create_todo calls

## Workflow

1. **Before starting:** Use `create_todos` to break down the task into steps
2. **While working:** Use `list_todos` to check what remains
3. **After each step:** Use `update_todo` to mark completed items
4. **Before finishing:** Use `list_todos` to verify all steps are done
