---
title: "Run your own container image"
linkTitle: "Run your own container image"
description: "Start a sandbox from a container image in a registry you name, instead of from a managed image."
keywords: "cloud sandboxes, sandboxes api, run your own container image"
weight: 501
params:
  sidebar:
    group: "Connect and configure"
---

Run your own tools from a container image when a [bundled kit](add-tools-with-kits.md) does not fit the task. You supply the image reference and machine size; the SDK creates an isolated sandbox around it.

You need an [authenticated client](connect-to-cloud-with-a-bearer-token.md) and an OCI image the service can pull. Use a versioned reference or digest when repeatability matters. The image must include `/bin/sh`.

## Create from a registry image {#1-create-from-a-registry-image}

Pass the image reference, CPU count, memory size, display name, and idempotency key. The SDK's image-reference option avoids assembling nested request fields. Do not also supply a managed image or named agent.

The example waits for the sandbox to run. A wait timeout stops your wait; it does not delete a sandbox that was already accepted. Save any sandbox handle retained by the error so you can inspect or clean it up.

The image's startup command runs inside the sandbox. Read `WORKSPACE_DIR` to find the workspace rather than assuming a path.

Next, [run a command](run-your-first-command.md) or [copy in your project files](copy-a-file-into-a-cloud-sandbox.md). [Delete the sandbox](delete-a-cloud-sandbox.md) when the work is complete.

{{< tabs >}}
{{< tab name="TypeScript" >}}

```typescript
const sandbox = await client.create(
  { displayName: name, imageRef, resources },
  { idempotencyKey: requestId },
);
return sandbox.waitUntilRunning();
```

<details>
<summary>Complete TypeScript example: rawimage/create.ts</summary>

```typescript
import type { ClientCreateOptions, Sandboxes } from '@docker/sandboxes';

export async function createFromImageRef(
  client: Sandboxes,
  name: string,
  imageRef: string,
  resources: ClientCreateOptions['resources'],
  requestId: string,
) {
  const sandbox = await client.create(
    { displayName: name, imageRef, resources },
    { idempotencyKey: requestId },
  );
  return sandbox.waitUntilRunning();
}
```

</details>

{{< /tab >}}
{{< /tabs >}}
