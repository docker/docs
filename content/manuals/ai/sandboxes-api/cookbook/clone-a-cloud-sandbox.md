---
title: "Clone a cloud sandbox"
linkTitle: "Clone a cloud sandbox"
description: "Read how a running sandbox was configured, then create a second sandbox with the same image and environment."
keywords: "cloud sandboxes, sandboxes api, clone a cloud sandbox"
weight: 702
params:
  sidebar:
    group: "Long-running work"
---

Create another sandbox using selected settings from an existing one. This example copies the image source and environment, not the sandbox's files or running processes. For a saved filesystem or memory state, use a [snapshot](snapshot-and-fork-a-sandbox.md).

You need an [authenticated client](connect-to-cloud-with-a-bearer-token.md), the source sandbox's resource name, and a display name for the new sandbox.

## Read the source configuration {#1-read-the-source-configuration}

Get the source sandbox and inspect its image and environment. Review environment values before copying them; they may contain credentials or settings specific to the original job.

Do not treat the complete response as a create request. It includes read-only state and connection details.

{{< tabs >}}
{{< tab name="TypeScript" >}}

```typescript
const sandbox = await client.get(name);
return sandbox.core;
```

<details>
<summary>Complete TypeScript example: clone/inspect.ts</summary>

```typescript
import type { Sandboxes } from '@docker/sandboxes';

export async function inspectSandbox(client: Sandboxes, name: string) {
  const sandbox = await client.get(name);
  return sandbox.core;
}
```

</details>

{{< /tab >}}
{{< /tabs >}}

## Create a fresh sandbox {#2-create-a-fresh-sandbox}

The example copies the image source and environment. For a registry image it also copies resource settings. It waits for the new sandbox to run.

The call returns the new sandbox's handle, ready for running commands or transferring files.

This is a limited configuration copy, not a complete clone. It does not copy kits, attached policies, secrets, volumes, files, or process state. A copied image alone may not reproduce a kit's setup. For repeated agent environments, launch the same [named kit](add-tools-with-kits.md) with the options your application saved.

Attach required network policies and credentials before running work in the new sandbox. The source is unchanged, and both sandboxes need their own cleanup.

{{< tabs >}}
{{< tab name="TypeScript" >}}

```typescript
const copy = await client.create(
  {
    displayName: name,
    ...(core.image
      ? { image: core.image, agent: core.agent }
      : {
          imageRef: core.imageRef!,
          resources: core.resources,
        }),
    environment: core.environment,
  },
  { timeoutMs: 300_000, idempotencyKey: requestId },
);
return copy.waitUntilRunning();
```

<details>
<summary>Complete TypeScript example: clone/rebuild.ts</summary>

```typescript
import type { Sandbox, Sandboxes } from '@docker/sandboxes';

export async function rebuildSandbox(
  client: Sandboxes,
  core: Sandbox['core'],
  name: string,
  requestId: string,
) {
  if (!core.image && !core.imageRef)
    throw new TypeError('Source has no image');
  const copy = await client.create(
    {
      displayName: name,
      ...(core.image
        ? { image: core.image, agent: core.agent }
        : {
            imageRef: core.imageRef!,
            resources: core.resources,
          }),
      environment: core.environment,
    },
    { timeoutMs: 300_000, idempotencyKey: requestId },
  );
  return copy.waitUntilRunning();
}
```

</details>

{{< /tab >}}
{{< /tabs >}}
