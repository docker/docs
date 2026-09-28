---
title: Kit version and interface compatibility
linkTitle: Compatibility
description: Check kit descriptor versions, references, and artifact handling before reusing a kit.
keywords: docker sandboxes, sbx, cloud sandboxes
weight: 5
---

Kit format and launch interface both affect how you can reuse a kit. A kit's
OCI image reference identifies a published artifact, but not every interface
accepts that reference directly.

## Documented launch paths

| Interface | Documented input | What to check |
| --- | --- | --- |
| Local CLI | Built-in agent names, published kits, local paths, and Git references | V3 requires `sbx` v0.45 or later; workload and mixin versions must match |
| Cloud CLI | Built-in agents and cloud kit workflows | Verify descriptor support and host-dependent capabilities before adapting a custom kit |
| Console | Curated kits and public custom kit references | Follow the Console's custom-kit requirements; do not infer backend descriptor support from the authoring format |
| SDK | Names from the bundled catalog, or serialized v2 custom kit artifacts | Launch-by-name does not accept an arbitrary OCI reference; custom artifacts require preparation |
| REST API | Serialized kit artifacts; the SDK guide documents v2 | Prepare artifact bytes; a source reference does not trigger a registry pull |

See [CLI kit usage](/manuals/ai/sandboxes/cli/kits.md), [Console kits](/manuals/ai/sandboxes/console/kits.md), and
[SDK custom kits](/manuals/ai/sandboxes/api/concepts.md) for each launch path.

## V3 and earlier formats

V3 workloads and mixins must be used together. The CLI's built-in agent
shortcuts, such as `claude` and `codex`, select v2 kits and cannot be combined
with v3 mixins. Select a v3 workload explicitly to use v3 mixins.

The [v2 guide](/manuals/ai/sandboxes/author-kits/kits-v2.md) covers maintenance of earlier kits. Do not treat
local v3 support as proof that the same artifact can be passed unchanged to
every cloud interface. Cloud conversion behavior and accepted artifact formats
need to be checked for the interface you use.

## Artifact preparation

The SDK's custom-artifact workflow operates on the descriptor and filesystem
artifact bytes. Resolving a published kit and preparing those bytes is
additional application work. A bundled catalog name, a container image used
as the sandbox base, and a custom kit artifact are different inputs.

See [API resources and kit inputs](/manuals/ai/sandboxes/api/concepts.md) and the
[API reference](/reference/api/sandboxes/latest/) before implementing that path.
