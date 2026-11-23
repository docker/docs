---
title: Customize an environment with API v2
description: Compose a sandbox from a v3 workload kit and compatible mixin kits.
keywords: docker sandboxes, API, SDK, local sandboxes, cloud sandboxes
weight: 10
linkTitle: Customize an environment
---

Compose a repeatable environment from a workload kit and compatible mixins.
This workflow applies to local and cloud sandboxes, subject to the capabilities
required by each kit.

A workload supplies the base environment and launch command. Mixins add tools,
configuration, or runtime behavior. Choose published v3 kit references and
preview the composition to understand its requirements before creating a
sandbox. The API resolves the kits, so your application does not need to package
their files into the request.

Kit format support alone does not guarantee that a kit works in both local and
cloud sandboxes. Check its required capabilities, credential requirements, and
supported execution environments.

For kit concepts and authoring, see
[Kits](/manuals/ai/sandboxes/customize/_index.md) and
[Author kits](/manuals/ai/sandboxes/customize/author/_index.md).

<!-- Add a verified preview/create example with one workload and one mixin.
Cover references, arguments, requirements, readiness, and cleanup. Keep CLI kit
source support distinct from the API's accepted registry references. -->
