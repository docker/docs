---
title: "Manage cloud secrets"
linkTitle: "Manage cloud secrets"
description: "Store a token as a stored secret, list your secrets' metadata, replace the token, and delete the secret."
keywords: "cloud sandboxes, sandboxes api, manage cloud secrets"
weight: 402
params:
  sidebar:
    group: "Security and policy"
---

Store workload credentials separately from application code. Secret reads return metadata only, so keep the original credential in your secret manager if you need it elsewhere.

Use an [authenticated client](connect-to-cloud-with-a-bearer-token.md). These examples store a token for a named service. For instructions on giving a kit access to that token, see [Get a stored secret into a sandbox](get-a-stored-secret-into-a-sandbox.md).

## Create a secret {#1-create-a-secret}

Pass a display name, service type, token value, and idempotency key. Save the returned resource name for attachments and future updates. Do not log the request or token.

{{< tabs >}}
{{< tab name="TypeScript" >}}

```typescript
return client.secrets.create(
  { displayName: name, serviceType, token: { value: token } },
  { idempotencyKey: requestId },
);
```

<details>
<summary>Complete TypeScript example: secrets/create.ts</summary>

```typescript
import type { Sandboxes } from '@docker/sandboxes';

export async function createSecret(
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

## List stored credentials {#2-list-stored-credentials}

List secret metadata to find a credential by name or label. The example follows each page. Listing does not reveal secret values.

Use the resource name, not the display name, when attaching a secret to a sandbox.

{{< tabs >}}
{{< tab name="TypeScript" >}}

```typescript
return client.secrets.all().collect();
```

<details>
<summary>Complete TypeScript example: secrets/list.ts</summary>

```typescript
import type { Sandboxes } from '@docker/sandboxes';

export async function listSecrets(client: Sandboxes) {
  return client.secrets.all().collect();
}
```

</details>

{{< /tab >}}
{{< /tabs >}}

## Rotate the credential {#3-rotate-the-credential}

Read the existing secret, then update it through a handle with the replacement value and service type. The handle carries the version it read so an intervening change is detected.

The update replaces the record's writable content. Send the service type again instead of assuming omitted values are preserved. Keep the returned metadata for subsequent operations.

If another writer changed the secret, read it again and decide whether to apply your replacement. Do not silently overwrite another rotation.

{{< tabs >}}
{{< tab name="TypeScript" >}}

```typescript
return secret.update(
  { serviceType, token: { value: token } },
  { idempotencyKey: requestId },
);
```

<details>
<summary>Complete TypeScript example: secrets/update.ts</summary>

```typescript
import type { Secret } from '@docker/sandboxes';

export async function updateSecret(
  secret: Secret,
  serviceType: string,
  token: string,
  requestId: string,
) {
  return secret.update(
    { serviceType, token: { value: token } },
    { idempotencyKey: requestId },
  );
}
```

</details>

{{< /tab >}}
{{< /tabs >}}

## Delete a secret {#4-delete-a-secret}

Delete through the secret handle when no workload needs the credential. Deleting an already absent secret is safe for repeated cleanup.

Removing a stored secret is not a substitute for revoking a leaked credential with its provider. Revoke compromised credentials and replace them before starting new workloads.

{{< tabs >}}
{{< tab name="TypeScript" >}}

```typescript
await secret.delete();
```

<details>
<summary>Complete TypeScript example: secrets/delete.ts</summary>

```typescript
import type { Secret } from '@docker/sandboxes';

export async function deleteSecret(secret: Secret) {
  await secret.delete();
}
```

</details>

{{< /tab >}}
{{< /tabs >}}
