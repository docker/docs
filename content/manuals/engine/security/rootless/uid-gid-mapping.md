---
description: How container UIDs and GIDs are mapped to the host in rootless mode
keywords: security, namespaces, rootless, uid, gid, subuid, subgid
title: UID/GID mapping
weight: 15
---

Rootless mode and [`userns-remap` mode](../userns-remap.md) map container UIDs
and GIDs to the host differently.

- In `userns-remap` mode, container UID `0` is mapped to the first subordinate
  UID listed in `/etc/subuid` for the remap user, and container UID `n` is
  mapped to `subuid + n`.
- In rootless mode, container UID `0` is mapped to the host UID of the user
  running rootless Docker (the result of `id -u`); container UID `n` (for
  `n >= 1`) is mapped to `subuid + (n - 1)`.

GIDs follow the same rules using `/etc/subgid`.

This difference matters when setting file permissions on bind-mounted
directories: in rootless mode, files owned by your host user appear as owned
by `root` inside the container.
