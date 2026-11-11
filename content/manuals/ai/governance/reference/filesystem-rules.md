---
title: Filesystem mount rule syntax
linkTitle: Filesystem mount rules
weight: 20
description: Syntax for sandbox filesystem mount rules, including path patterns across macOS, Linux, Windows, and WSL.
keywords: docker sandboxes, filesystem rules, path patterns, workspace mount, rule syntax
---

Filesystem mount rules make up
[organization filesystem access policies](../organization/policies/filesystem.md).
Each rule lists host paths that sandboxes can mount as workspaces, and the
access it allows or denies on them: `read`, `write`, or both. To add rules to a
policy, see
[Create and manage organization policies](../organization/policies/create-manage.md).

When a sandbox mounts a workspace, Docker Sandboxes checks the path against
these rules for the access the mount needs. A read-only workspace needs only
`read`. A writable workspace needs both `read` and `write`, so the path must be
allowed for both actions.

## Path patterns

### Home directory

`~` expands to the user's home directory on every platform, including Windows,
where it resolves to `%USERPROFILE%`. A single `~/**` rule therefore matches
each user's home tree on macOS, Linux, and Windows. The policy engine expands
only `~`: it does not expand environment variables, so a pattern such as
`%USERPROFILE%\**` or `$HOME/**` matches nothing.

### Operating system formats

For a path outside the home directory, write it in the format the user's
operating system uses. A macOS or Linux path doesn't match a Windows or WSL
path, so a location that several platforms share needs a rule for each:

| Operating system | Example path                               |
| ---------------- | ------------------------------------------ |
| macOS, Linux     | `/data/project/**`                         |
| Windows          | `C:\data\project\**`                       |
| WSL              | `\\wsl.localhost\<distro>\data\project\**` |

On Windows, `*:` matches any drive letter, so `*:\data\**` matches the path on
any drive.

### Wildcards

Wildcards behave the same way in every path format:

| Pattern            | Example    | Matches                                                    |
| ------------------ | ---------- | ---------------------------------------------------------- |
| Exact path         | `/data`    | `/data` only                                               |
| Segment wildcard   | `/data/*`  | `/data/project`, one path segment only, not subdirectories |
| Recursive wildcard | `/data/**` | `/data/project`, `/data/project/src`, any depth            |

For example, `~/**` matches all paths under the home directory, while `~/*`
matches only its direct children.

### Case and separators

Matching ignores some differences in how a path is written:

| Difference                                         | Ignored | Example                                       |
| -------------------------------------------------- | ------- | --------------------------------------------- |
| Slash direction                                    | Yes     | `C:/data/**` and `C:\data\**` are the same    |
| Case in Windows drive and UNC paths, including WSL | Yes     | `C:\Data\**` matches `c:\data\project`        |
| Case in macOS and Linux paths                      | No      | `/Data/**` doesn't match `/data/project`      |
| WSL prefix                                         | Yes     | `\\wsl$\` and `\\wsl.localhost\` are the same |

macOS paths are case sensitive even on a case-insensitive volume, so write rules
with the directory names' capitalization on disk.

## How mount paths are resolved

Docker Sandboxes checks the rules against the resolved workspace path, not the
path a developer passes to `sbx run` or `sbx create`. It expands `~`, makes
the path absolute, removes `.` and `..` segments, and resolves the path as
follows:

| Developer passes             | Path checked      |
| ---------------------------- | ----------------- |
| A relative path, such as `.` | The absolute path |
| A symbolic link              | The link's target |

Because symbolic links are followed, a link under an allowed path is denied
when it points to a directory that no rule allows.
