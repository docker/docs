---
title: "Name a sandbox and find it again"
linkTitle: "Name a sandbox and find it again"
description: "Give a sandbox a readable label when you create it, and keep the resource name the service returns so you can address the same sandbox later."
keywords: "cloud sandboxes, sandboxes api, name a sandbox and find it again"
weight: 107
params:
  sidebar:
    group: "Get started"
---

Save a sandbox's resource name so another request or program run can find it. A display name is a human-readable label; the resource name is the identifier you pass to the SDK.

Start with an [authenticated client](connect-to-cloud-with-a-bearer-token.md). The creation example accepts a managed image and agent configuration. For a bundled starting point, use a [named kit](add-tools-with-kits.md).

## Choose a display name {#1-choose-a-display-name}

Create the sandbox with a label meaningful to your application, such as a job or project name. Use an idempotency key for the logical create operation.

Save the returned resource name and UID. Do not construct a resource name from the display name or assume labels are unique.

{{< tabs >}}
{{< tab name="TypeScript" >}}

```typescript
const sandbox = await client.create(
  { agent, displayName: name, image },
  { idempotencyKey: requestId },
);
return sandbox.waitUntilRunning();
```

<details>
<summary>Complete TypeScript example: naming/name.ts</summary>

```typescript
import type { Sandboxes } from '@docker/sandboxes';

export async function createNamedSandbox(
  client: Sandboxes,
  name: string,
  image: string,
  agent: string,
  requestId: string,
) {
  const sandbox = await client.create(
    { agent, displayName: name, image },
    { idempotencyKey: requestId },
  );
  return sandbox.waitUntilRunning();
}
```

</details>

{{< /tab >}}
{{< /tabs >}}

## Read it by resource name {#2-read-it-by-resource-name}

Pass the saved name to the client's get method to obtain a current sandbox handle, including its state and connection information. The example returns that handle, ready for commands and lifecycle operations.

If the sandbox is absent, do not silently replace it under the same application job without considering the original work. List and filter resources when you need to discover them; use get when you already know the name.

{{< tabs >}}
{{< tab name="TypeScript" >}}

```typescript
return client.get(name);
```

<details>
<summary>Complete TypeScript example: naming/resolve.ts</summary>

```typescript
import type { Sandboxes } from '@docker/sandboxes';

export async function readByName(client: Sandboxes, name: string) {
  return client.get(name);
}
```

</details>

{{< /tab >}}
{{< /tabs >}}
