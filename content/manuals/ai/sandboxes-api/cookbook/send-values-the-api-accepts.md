---
title: "Send values the API accepts"
linkTitle: "Send values the API accepts"
description: "Build a create request whose durations, counts, and flags mean what you intend, then read back the values the service settled on."
keywords: "cloud sandboxes, sandboxes api, send values the api accepts"
weight: 605
params:
  sidebar:
    group: "Requests and responses"
---

Build SDK options without confusing an omitted value with an explicit false or zero. This matters for settings whose default is chosen by the service.

Start with an [authenticated client](connect-to-cloud-with-a-bearer-token.md). Use the option types exported by your SDK so your editor can show the supported inputs.

## Construct creation options {#1-construct-creation-options}

Supply the image source and lifecycle settings, using the duration units required by your SDK. The example accepts seconds and converts them at the boundary.

This example uses a managed image, which supplies its resource defaults. For a registry image instead, use the image-reference option and provide a supported CPU and memory pair.

Preserve absence for an optional boolean when you want the service default. Explicit false is a choice, not a missing value. Validate user-supplied values before building the request.

{{< tabs >}}
{{< tab name="TypeScript" >}}

```typescript
return {
  displayName: name,
  image,
  platform,
  lifecycle: { timeoutMs: lifeSeconds * 1_000, autoResume },
};
```

<details>
<summary>Complete TypeScript example: values/build.ts</summary>

```typescript
import type { ClientCreateOptions, Sandboxes } from '@docker/sandboxes';

export function createRequest(
  name: string,
  image: string,
  lifeSeconds: number,
  platform: ClientCreateOptions['platform'],
  autoResume: boolean | undefined,
): ClientCreateOptions {
  return {
    displayName: name,
    image,
    platform,
    lifecycle: { timeoutMs: lifeSeconds * 1_000, autoResume },
  };
}

export async function sendCreate(
  client: Sandboxes,
  request: ClientCreateOptions,
) {
  const sandbox = await client.create(request, { timeoutMs: 300_000 });
  return sandbox.waitUntilRunning();
}
```

</details>

{{< /tab >}}
{{< /tabs >}}

## Read the effective values {#2-read-the-effective-values}

Get the sandbox and inspect the reported platform, resources, and expiration. These describe what the service recorded, which can differ from omitted or defaulted input values.

Keep the language's native numeric and duration representations. Avoid converting large integers through a floating-point type merely to display or serialize them.

{{< tabs >}}
{{< tab name="TypeScript" >}}

```typescript
return {
  expiresAt: sandbox.effectiveFeatures?.timeouts?.expiresAt,
  platform: sandbox.core.platform,
  cpus: sandbox.core.resources?.cpus ?? undefined,
};
```

<details>
<summary>Complete TypeScript example: values/read.ts</summary>

```typescript
import type { Sandbox, Sandboxes } from '@docker/sandboxes';

export function reportedValues(sandbox: Sandbox) {
  return {
    expiresAt: sandbox.effectiveFeatures?.timeouts?.expiresAt,
    platform: sandbox.core.platform,
    cpus: sandbox.core.resources?.cpus ?? undefined,
  };
}

export async function readSandbox(client: Sandboxes, name: string) {
  return reportedValues(await client.get(name));
}
```

</details>

{{< /tab >}}
{{< /tabs >}}
