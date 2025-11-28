---
title: Filesystem
description: Access and manage the filesystem
---

This toolset gives your agent access to your filesystem. By default the
filesystem tool has access to the current working directory, this can be
modified in your agent file, an agent can also decide, thanks to
`add_allowed_directory`, to ask for permission for a directory access.

## Usage

```yaml title="agent.yaml"
toolsets:
  - type: filesystem
```

## Tools

| Name                        | Description                                         |
| --------------------------- | --------------------------------------------------- |
| `add_allowed_directory`     | Adds a directory to the list of allowed directories |
| `create_directory`          | Creates a directory                                 |
| `directory_tree`            | Returns a recursive view of the directory tree      |
| `edit_file`                 | Modifies a file                                     |
| `get_file_info`             | Get stat info about a file/directory                |
| `list_allowed_directories`  | Returns the list of currently allowed directories   |
| `list_directory`            | Lists the contents of a directory                   |
| `list_directory_with_sizes` | Desc                                                |
| `move_file`                 | Moves a file                                        |
| `read_file`                 | Reads a file                                        |
| `read_multiple_files`       | Reads multiple files                                |
| `search_files`              | Search for files by filename/pattern                |
| `search_files_content`      | Grep-like search                                    |
| `write_file`                | Writes a file                                       |

### add_allowed_directory

By default the filesystem tool only has access to the current working directory.
This tool allows the agent to request access to other directories.

Args:

- `path`: The directory path to add to the list of allowed directories.

### create_directory

Creates a directory at the specified path.

Args:

- `path`: The directory path to create.

### directory_tree

Returns a recursive view of the directory tree starting from the specified path.

Args:

- `path`: The directory path to view.
- `max_depth`: The maximum depth to recurse into the directory tree.

### edit_file

Modifies a file at the specified path using a series of edit operations.

Args:

- `path`: The file path to edit.
- `edits`: Array of edit operations.

### get_file_info

Gets stat info about a file or directory.

Args:

- `path`: The file or directory path to get info about.

### list_allowed_directories

Returns the list of currently allowed directories.

### list_directory

Lists the contents of a directory.

Args:
- `path`: The directory path to list.

### list_directory_with_sizes

Lists the contents of a directory along with the sizes of each item.

Args:
- `path`: The directory path to list.

### move_file

Moves a file from one location to another.

Args:
- `source`: The source file path.
- `destination`: The destination file path.

### read_file

Reads a file at the specified path.

Args:
- `path`: The file path to read.

### read_multiple_files

Reads multiple files at the specified paths.

Args:
- `paths`: An array of file paths to read.
- `json` (optional): If true, returns the contents as a JSON object.

### search_files

Search for files by filename or pattern.

Args:
- `pattern`: The filename or pattern to search for.
- `path`: The starting directory path.
- `excludePatterns`: Patterns to exclude from search.

### search_files_content

Grep-like search for file contents.

Args:
- `pattern`: The content pattern to search for.
- `path`: The starting directory path.
- `query`: The text or regex pattern to search for.
- `isRegex`: If true, treat query as regex; otherwise literal text.
- `excludePatterns`: Patterns to exclude from search.

### write_file

Writes a file at the specified path.

Args:
- `path`: The file path to write.
- `content`: The content to write to the file.
