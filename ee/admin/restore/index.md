---
title: Restore Docker Enterprise
description: Learn how to restore Docker Enterprise platform from a backup.
keywords: enterprise, restore, recovery
---

You should only restore Docker Enterprise Edition from a backup as a last resort. If you're running Docker
Enterprise Edition in high-availability mode, you can remove unhealthy nodes from the
swarm and join new ones to bring the swarm to an healthy state.

Restore components individually and in the following order:

1. [Restore Docker Swarm](restore-swarm.md).
2. [Restore Universal Control Plane (UCP)](restore-ucp.md).
3. [Restore Docker Trusted Registry (DTR)](restore-dtr.md).

## Where to go next

- [Restore Docker Swarm](restore-swarm.md)
