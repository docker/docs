---
title: "Expose a port from a cloud sandbox"
linkTitle: "Expose a port from a cloud sandbox"
description: "Publish a TCP port that a program inside a cloud sandbox listens on, read back the URL that reaches it, list the sandbox's published ports, and withdraw the publication."
keywords: "cloud sandboxes, sandboxes api, expose a port from a cloud sandbox"
weight: 204
params:
  sidebar:
    group: "Working in a sandbox"
---

Reach a web application running inside a sandbox through an HTTPS URL. Publishing creates a route to the application; it does not start the application.

You need an [authenticated client](connect-to-cloud-with-a-bearer-token.md), the sandbox name, and an HTTP application listening on the port you will publish. Bind the application to an address reachable inside the sandbox, such as `0.0.0.0`.

## Publish the application's port {#1-publish-the-application-s-port}

Supply the sandbox name and the application's port number. The example uses TCP. Use the URL returned by the service rather than constructing a hostname.

This is HTTP application access, not a general-purpose TCP tunnel. Docker account authentication does not protect the application URL. Configure authentication in your application before exposing sensitive data, and never send your Docker account token to that URL.

{{< tabs >}}
{{< tab name="TypeScript" >}}

```typescript
return sandbox.ports.create(
  { number, protocol: 'tcp' },
  { idempotencyKey: requestId },
);
```

<details>
<summary>Complete TypeScript example: ports/publish.ts</summary>

```typescript
import type { Sandboxes } from '@docker/sandboxes';

export async function createPort(
  client: Sandboxes,
  sandboxName: string,
  number: number,
  requestId: string,
) {
  const sandbox = await client.get(sandboxName);
  return sandbox.ports.create(
    { number, protocol: 'tcp' },
    { idempotencyKey: requestId },
  );
}
```

</details>

{{< /tab >}}
{{< /tabs >}}

## List published ports {#2-list-published-ports}

Read the sandbox's port collection to find existing publications and their URLs. Check the application's readiness separately: a published route does not prove the application has started.

{{< tabs >}}
{{< tab name="TypeScript" >}}

```typescript
return (await sandbox.ports.list()).published ?? [];
```

<details>
<summary>Complete TypeScript example: ports/list.ts</summary>

```typescript
import type { Sandbox } from '@docker/sandboxes';

export async function listPorts(sandbox: Sandbox) {
  return (await sandbox.ports.list()).published ?? [];
}
```

</details>

{{< /tab >}}
{{< /tabs >}}

## Withdraw a publication {#3-withdraw-a-publication}

Delete the published port when access is no longer needed. This removes the route but leaves the application process running.

Stopping the sandbox does not remove its port publications. Deleting the sandbox does.

{{< tabs >}}
{{< tab name="TypeScript" >}}

```typescript
await port.delete();
```

<details>
<summary>Complete TypeScript example: ports/unpublish.ts</summary>

```typescript
import type { Port } from '@docker/sandboxes';

export async function deletePort(port: Port) {
  await port.delete();
}
```

</details>

{{< /tab >}}
{{< /tabs >}}
