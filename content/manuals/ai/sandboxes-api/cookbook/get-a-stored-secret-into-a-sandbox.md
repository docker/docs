---
title: "Get a stored secret into a sandbox"
linkTitle: "Get a stored secret into a sandbox"
description: "Store a credential once, then pass its resource name when creating a sandbox so the sandbox starts with that secret attached."
keywords: "cloud sandboxes, sandboxes api, get a stored secret into a sandbox"
weight: 405
params:
  sidebar:
    group: "Security and policy"
---

Give an agent access to its provider without placing the real credential in its command arguments or source files. Store the credential once, then attach its secret name when creating each sandbox.

Start with an [authenticated client](connect-to-cloud-with-a-bearer-token.md). Obtain the provider credential from that provider and read it from your application's secret manager. For Claude Code, use an Anthropic API key and service type `anthropic`; the Docker login token is not an Anthropic key.

## Store the provider credential {#1-store-the-provider-credential}

Supply a display name, service type, token value, and an idempotency key. Save the returned secret name. The response contains metadata, never the stored token.

Choose the service type that matches the workload's provider. A token for one provider does not authenticate another agent.

{{< tabs >}}
{{< tab name="TypeScript" >}}

```typescript
return client.secrets.create(
  { displayName: name, serviceType, token: { value: token } },
  { idempotencyKey: requestId },
);
```

<details>
<summary>Complete TypeScript example: credinject/store.ts</summary>

```typescript
import type { Sandboxes } from '@docker/sandboxes';

export async function storeSecret(
  client: Sandboxes,
  name: string,
  serviceType: string,
  token: string,
  requestId: string,
) {
  return client.secrets.create(
    { displayName: name, serviceType, token: { value: token } },
    { idempotencyKey: requestId },
  );
}
```

</details>

{{< /tab >}}
{{< /tabs >}}

## Attach the secret at creation {#2-attach-the-secret-at-creation}

Put the secret name in the sandbox's storage options. This example uses a managed image; the same storage options can be passed when [launching a named kit](add-tools-with-kits.md).

Attach the secret when creating the sandbox, before starting the agent. Passing a secret's display name instead of its resource name will not select it. Keep the real token on the client side of the storage step rather than duplicating it in environment variables.

The SDK returns a handle after the example waits for the sandbox to run. Use that handle to run the agent. To replace or delete the credential later, follow [Manage cloud secrets](manage-cloud-secrets.md).

{{< tabs >}}
{{< tab name="TypeScript" >}}

```typescript
const sandbox = await client.create(
  { displayName: name, image, storage: { secrets: secretNames } },
  { timeoutMs: 300_000, idempotencyKey: requestId },
);
return sandbox.waitUntilRunning();
```

<details>
<summary>Complete TypeScript example: credinject/attach.ts</summary>

```typescript
import type { Sandboxes } from '@docker/sandboxes';

export async function createWithSecrets(
  client: Sandboxes,
  name: string,
  image: string,
  secretNames: string[],
  requestId: string,
) {
  const sandbox = await client.create(
    { displayName: name, image, storage: { secrets: secretNames } },
    { timeoutMs: 300_000, idempotencyKey: requestId },
  );
  return sandbox.waitUntilRunning();
}
```

</details>

{{< /tab >}}
{{< /tabs >}}
