---
title: Memory
description: Agentic memory for your agent
---

The memory toolset gives your agent the ability to store, retrieve, and delete
persistent memories about the user. This allows the agent to remember important
information across conversations.

## Usage

```yaml
toolsets:
  - type: memory
    path: "./memories.db"
```

## Tools

| Name             | Description                            |
| ---------------- | -------------------------------------- |
| `add_memory`     | Add a new memory to the database       |
| `get_memories`   | Retrieve all stored memories           |
| `delete_memory`  | Delete a specific memory by ID         |

### add_memory

Adds a new memory to the database with an auto-generated ID and timestamp.

Args:

- `memory`: The memory content to store (required)

**Example:**

```json
{
  "memory": "User prefers dark mode for all interfaces"
}
```

### get_memories

Retrieves all stored memories from the database. Returns an array of memory
objects, each containing an ID, creation timestamp, and the memory content.

No arguments required.

**Response format:**

```json
[
  {
    "id": "1234567890",
    "createdAt": "2024-01-15T10:30:00Z",
    "memory": "User prefers dark mode for all interfaces"
  }
]
```

### delete_memory

Deletes a specific memory from the database by its ID.

Args:

- `id`: The ID of the memory to delete (required)

**Example:**

```json
{
  "id": "1234567890"
}
```

## Best Practices

- Use `get_memories` at the start of conversations to retrieve relevant context
- Store important user preferences, facts, and context
- Delete outdated or incorrect memories when necessary
- Create specific, actionable memories rather than vague observations
