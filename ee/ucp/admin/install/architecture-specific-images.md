---
title: Architecture-specific images
description: Learn how to deploy Docker Universal Control Plane using images that are specific to particular hardware architectures.
keywords: UCP, Docker EE, image, Windows
---

Docker Universal Control Plane deploys images for a number of different
hardware architectures. Some architectures require
pulling images with specific tags or names indicating the target
architecture.

## OS-specific component names

Some UCP component names depend on the node's operating system. Use the
following table to ensure that you're pulling the right images for each node.

| UCP component base name | Windows name   | 
|-------------------------|----------------|
| ucp-agent               | ucp-agent-win  | 
| ucp-dsinfo              | ucp-dsinfo-win |

## Where to go next

- [Join nodes to your cluster](../configure/join-nodes.md)
