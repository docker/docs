---
title: "File Tool"
description: "Read, write, and edit individual files."
keywords: docker agent, ai agents, tools, toolsets, file tool
linkTitle: "File"
weight: 9
canonical: https://docs.docker.com/ai/docker-agent/tools/file/
---

_Read, write, and edit individual files without exposing directory or search operations._

## Overview

The file toolset is a focused subset of the [filesystem toolset](../filesystem/index.md). It reuses the filesystem implementations and exposes only these tools:

| Tool | Description |
| --- | --- |
| `read_file` | Read a whole file or a selected line range. Supports text and images. |
| `write_file` | Create a file or completely overwrite an existing file. |
| `edit_file` | Apply one or more exact text replacements to an existing file. |

Use this toolset when an agent needs to work with known files but should not list directories, build directory trees, search file contents, or create and remove directories.

## Configuration

```yaml
toolsets:
  - type: file
```

The file toolset supports the filesystem toolset's path and editing options:

| Property | Type | Default | Description |
| --- | --- | --- | --- |
| `post_edit` | array | `[]` | Commands to run after `write_file` or `edit_file` for matching paths. |
| `allow_list` | array | `[]` | Directories the tools may access. Empty means unrestricted. |
| `deny_list` | array | `[]` | Directories the tools must not access. Takes precedence over `allow_list`. |

Paths resolve relative to the agent's working directory. Home-directory expansion, absolute paths, access controls, symlink protection, post-edit hooks, and `.agentsignore` handling behave exactly as documented for the [filesystem toolset](../filesystem/index.md).

```yaml
toolsets:
  - type: file
    allow_list:
      - "."
    deny_list:
      - ".env"
    post_edit:
      - path: "*.go"
        cmd: "gofmt -w ${file}"
```
