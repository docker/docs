---
title: "Register and manage an image"
linkTitle: "Register and manage an image"
description: "Register an image with Cloud Sandboxes, get the target to push its content to, check when it is ready, list your images, and delete one you no longer need."
keywords: "cloud sandboxes, sandboxes api, register and manage an image"
weight: 301
params:
  sidebar:
    group: "Packaging and state"
---

Reuse an image across sandboxes by registering it with Docker Cloud Sandboxes. Registration creates an image record and a temporary upload destination. You must push the image content separately with an OCI-compatible registry client.

Start with an [authenticated client](connect-to-cloud-with-a-bearer-token.md) and a built image. If you already have an image in a registry, [run it directly](run-your-own-container-image.md) instead.

## Register an upload destination {#1-register-an-upload-destination}

Choose a display name, startup command, CPU and memory preferences, and an idempotency key. The returned image includes its resource name and a push target.

Use the target's registry reference and temporary credential to push your image. Keep that credential private and finish before it expires. Save the target from the initial response: later image reads do not issue a replacement upload credential.

This example registers the destination. It does not build or push image content.

{{< tabs >}}
{{< tab name="TypeScript" >}}

```typescript
return client.images.create(
  { displayName: name, fromImage: { resources }, startCmd },
  { idempotencyKey: requestId },
);
```

<details>
<summary>Complete TypeScript example: images/create.ts</summary>

```typescript
import type {
  ClientImagesCreateOptions,
  Sandboxes,
} from '@docker/sandboxes';

export async function createImage(
  client: Sandboxes,
  name: string,
  startCmd: string[],
  resources: Extract<
    ClientImagesCreateOptions,
    { fromImage: unknown }
  >['fromImage']['resources'],
  requestId: string,
) {
  return client.images.create(
    { displayName: name, fromImage: { resources }, startCmd },
    { idempotencyKey: requestId },
  );
}
```

</details>

{{< /tab >}}
{{< /tabs >}}

## Check preparation status {#2-check-preparation-status}

After the push finishes, read the image by its resource name. A `COMPLETED` status means it is ready to use. `WAITING_FOR_PUSH` means content has not arrived, `PREPARING` means preparation is in progress, and `FAILED` includes failure details.

The example performs one read. If preparation is still in progress, repeat the read with a delay and deadline. Repeating creation would register another image rather than advance this one.

{{< tabs >}}
{{< tab name="TypeScript" >}}

```typescript
return client.images.get(name);
```

<details>
<summary>Complete TypeScript example: images/ready.ts</summary>

```typescript
import type { Sandboxes } from '@docker/sandboxes';

export async function getImage(client: Sandboxes, name: string) {
  return client.images.get(name);
}
```

</details>

{{< /tab >}}
{{< /tabs >}}

## Find registered images {#3-find-registered-images}

List images with an optional filter, such as `status=completed`. The example follows all pages. List entries are summaries; get the image by name when you need its complete record.

{{< tabs >}}
{{< tab name="TypeScript" >}}

```typescript
return client.images.all({ filter }).collect();
```

<details>
<summary>Complete TypeScript example: images/list.ts</summary>

```typescript
import type { Sandboxes } from '@docker/sandboxes';

export async function listImages(client: Sandboxes, filter: string) {
  return client.images.all({ filter }).collect();
}
```

</details>

{{< /tab >}}
{{< /tabs >}}

## Remove an image {#4-remove-an-image}

Read the image and delete through its handle. The SDK supplies its name and version to protect against deleting a newer record you have not read.

Keep images that future sandbox creations still depend on. Deleting an already absent image is safe for cleanup.

{{< tabs >}}
{{< tab name="TypeScript" >}}

```typescript
await image.delete();
```

<details>
<summary>Complete TypeScript example: images/delete.ts</summary>

```typescript
import type { Image } from '@docker/sandboxes';

export async function deleteImage(image: Image) {
  await image.delete();
}
```

</details>

{{< /tab >}}
{{< /tabs >}}
