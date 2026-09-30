---
title: "Attach persistent storage"
linkTitle: "Attach persistent storage"
description: "Create a volume, mount it into a new sandbox so its data outlives the sandbox, and delete the volume when you no longer need it."
keywords: "cloud sandboxes, sandboxes api, attach persistent storage"
weight: 505
params:
  sidebar:
    group: "Connect and configure"
---

Keep project data after a sandbox is deleted by mounting a persistent volume. A volume is an independent resource: deleting a sandbox does not delete the volume.

Start with an [authenticated client](connect-to-cloud-with-a-bearer-token.md), a managed image from [Register and manage an image](register-and-manage-an-image.md), and an absolute mount path such as `/workspace/data`.

## Create a volume {#1-create-a-volume}

Choose a display name and an idempotency key for the create request. Save the returned volume name, such as `volumes/…`. That name identifies the volume in later requests; its display name is only a label.

{{< tabs >}}
{{< tab name="TypeScript" >}}

```typescript
return client.volumes.create(
  { displayName: name },
  { idempotencyKey: requestId },
);
```

<details>
<summary>Complete TypeScript example: volumes/create.ts</summary>

```typescript
import type { Sandboxes } from '@docker/sandboxes';

export async function createVolume(
  client: Sandboxes,
  name: string,
  requestId: string,
) {
  return client.volumes.create(
    { displayName: name },
    { idempotencyKey: requestId },
  );
}
```

</details>

{{< /tab >}}
{{< /tabs >}}

## Mount it in a new sandbox {#2-mount-it-in-a-new-sandbox}

Supply the volume name and destination path when creating the sandbox. The example accepts a map of volume names to mount paths and waits for the sandbox to run. For Cloud Sandboxes, supply one entry and use exclusive attachment. Volume access must be enabled for your account; memory-snapshot images do not support volume attachment.

Attachments belong to the sandbox's creation settings. To mount the volume elsewhere, create another sandbox with the attachment. With exclusive attachment, release the first sandbox before attaching the volume to another one. Avoid overlapping mount paths.

Write a file under the mount path, delete the sandbox, then attach the volume to a new sandbox to read it again. Files outside the mount remain on the sandbox's own disk.

{{< tabs >}}
{{< tab name="TypeScript" >}}

```typescript
const sandbox = await client.create(
  {
    displayName: name,
    image,
    storage: {
      volumes: [...mountPaths].map(([volume, target]) => ({
        volume,
        target,
      })),
    },
  },
  { timeoutMs: 300_000, idempotencyKey: requestId },
);
return sandbox.waitUntilRunning();
```

<details>
<summary>Complete TypeScript example: volumes/attach.ts</summary>

```typescript
import type { Sandboxes } from '@docker/sandboxes';

export async function createWithVolumes(
  client: Sandboxes,
  name: string,
  image: string,
  mountPaths: Map<string, string>,
  requestId: string,
) {
  const sandbox = await client.create(
    {
      displayName: name,
      image,
      storage: {
        volumes: [...mountPaths].map(([volume, target]) => ({
          volume,
          target,
        })),
      },
    },
    { timeoutMs: 300_000, idempotencyKey: requestId },
  );
  return sandbox.waitUntilRunning();
}
```

</details>

{{< /tab >}}
{{< /tabs >}}

## Delete the volume {#3-delete-the-volume}

Delete the sandbox using the volume first, then delete the volume through its handle. If the service reports that it is still in use, wait for sandbox deletion to finish before retrying.

Deleting the volume permanently removes its contents. Copy out anything you need to keep.

{{< tabs >}}
{{< tab name="TypeScript" >}}

```typescript
await volume.delete();
```

<details>
<summary>Complete TypeScript example: volumes/delete.ts</summary>

```typescript
import type { Volume } from '@docker/sandboxes';

export async function deleteVolume(volume: Volume) {
  await volume.delete();
}
```

</details>

{{< /tab >}}
{{< /tabs >}}
