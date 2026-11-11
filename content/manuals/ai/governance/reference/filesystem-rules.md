---
title: Filesystem mount rule syntax
linkTitle: Filesystem mount rules
weight: 20
description: Syntax for sandbox filesystem mount rules, including path patterns across macOS, Linux, Windows, and WSL.
keywords: docker sandboxes, filesystem rules, path patterns, workspace mount, rule syntax
---

Filesystem mount rules are written in
[organization filesystem access policies](../organization/policies/filesystem.md).
Developers can't add them with local policy. Filesystem mount rules use the
actions `read` and `write`. Resources are host paths that sandboxes can mount
as workspaces. The rules control which paths a sandbox can mount, not access to
individual files inside an allowed mount.

A workspace mounted with write access must be allowed by both a `read` and a
`write` rule. A read-only workspace needs only `read`. When default deny blocks
a mount, the denial reason names whether read or write access was missing.

`~` expands to the user's home directory on every platform, including Windows,
where it resolves to `%USERPROFILE%`. A single `~/**` rule therefore matches
each user's home tree on macOS, Linux, and Windows. The policy engine expands
only `~`: it does not expand environment variables, so a pattern such as
`%USERPROFILE%\**` or `$HOME/**` matches nothing.

For a path outside the home directory, write it in the format the user's
operating system uses. A rule matches only the format it's written in, so a
location that several platforms share needs a rule for each:

| Operating system | Example path                               |
| ---------------- | ------------------------------------------ |
| macOS, Linux     | `/data/project/**`                         |
| Windows          | `C:\data\project\**`                       |
| WSL              | `\\wsl.localhost\<distro>\data\project\**` |

On Windows, `*:` matches any drive letter, so `*:\data\**` matches the path on
any drive.

Case sensitivity depends on the path format:

- POSIX paths on macOS and Linux are case-sensitive: `/data/Project/**` doesn't
  match `/data/project/file`. Use the directory names' capitalization on disk,
  even on a case-insensitive macOS volume.
- Windows drive paths and UNC paths, including WSL UNC paths, are
  case-insensitive: `C:\Data\Project\**` matches `c:\data\project\file`.

Wildcards behave the same way in every path format:

| Pattern            | Example    | Matches                                                    |
| ------------------ | ---------- | ---------------------------------------------------------- |
| Exact path         | `/data`    | `/data` only                                               |
| Segment wildcard   | `/data/*`  | `/data/project`, one path segment only, not subdirectories |
| Recursive wildcard | `/data/**` | `/data/project`, `/data/project/src`, any depth            |

Use `**` to match a directory tree recursively. A single `*` matches within one
path segment and won't cross a path separator. For example, `~/**` matches all
paths under the home directory, while `~/*` matches only its direct children.
