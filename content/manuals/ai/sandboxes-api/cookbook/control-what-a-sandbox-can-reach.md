---
title: "Control what a sandbox can reach"
linkTitle: "Control what a sandbox can reach"
description: "Attach your account's network policies to a new sandbox, read the policy the service enforces, and review which destinations it allowed or blocked."
keywords: "cloud sandboxes, sandboxes api, control what a sandbox can reach"
weight: 401
params:
  sidebar:
    group: "Security and policy"
---

Control which destinations an agent or command can reach. A sandbox's effective policy combines account rules with the policies attached to it.

You need an [authenticated client](connect-to-cloud-with-a-bearer-token.md), a managed image, and policy IDs from your account. Create policies in the Console or with [Manage network policies](manage-network-policies.md).

## Attach policies at creation {#1-attach-policies-at-creation}

Pass the policy IDs when creating the sandbox. The example waits until the sandbox is running and returns its handle.

Attachments refer to existing policies; they do not contain policy definitions. Keep the IDs distinct from display names. Attached policies restrict access alongside account policy, so they cannot grant access that your account denies.

Include every destination the workload needs, including its model provider and package registries. A kit's network requirements do not override account restrictions.

{{< tabs >}}
{{< tab name="TypeScript" >}}

```typescript
const sandbox = await client.create(
  { displayName: name, image, network: { policyIds } },
  { idempotencyKey: requestId },
);
return sandbox.waitUntilRunning();
```

<details>
<summary>Complete TypeScript example: netpolicy/reference.ts</summary>

```typescript
import type { Sandboxes } from '@docker/sandboxes';

export async function createWithPolicyIds(
  client: Sandboxes,
  name: string,
  image: string,
  policyIds: string[],
  requestId: string,
) {
  const sandbox = await client.create(
    { displayName: name, image, network: { policyIds } },
    { idempotencyKey: requestId },
  );
  return sandbox.waitUntilRunning();
}
```

</details>

{{< /tab >}}
{{< /tabs >}}

## Inspect the rules in force {#2-inspect-the-rules-in-force}

Read the sandbox's effective policy before troubleshooting an application timeout. The result includes combined allow and deny rules, their origins, and an exact policy representation.

The effective view is the useful starting point for questions such as “can this sandbox reach api.anthropic.com?” Check the default mode as well as matching rules. A deny-by-default policy needs an explicit allowance for the destination.

{{< tabs >}}
{{< tab name="TypeScript" >}}

```typescript
return sandbox.networkPolicies.get();
```

<details>
<summary>Complete TypeScript example: netpolicy/effective.ts</summary>

```typescript
import type { Sandbox } from '@docker/sandboxes';

export async function getNetworkPolicies(sandbox: Sandbox) {
  return sandbox.networkPolicies.get();
}
```

</details>

{{< /tab >}}
{{< /tabs >}}

## Check allowed and blocked traffic {#3-check-allowed-and-blocked-traffic}

List policy log entries for the sandbox, using `domain=network` to select network decisions. The example follows pagination and returns all matching entries.

Use each entry's destination, decision, reason, and count to identify blocked dependencies. An empty result means no matching entries were returned; it does not prove that a destination is reachable. The log records recent activity, not a permanent audit history.

{{< tabs >}}
{{< tab name="TypeScript" >}}

```typescript
return sandbox.networkPolicies.logs.all({ filter }).collect();
```

<details>
<summary>Complete TypeScript example: netpolicy/logs.ts</summary>

```typescript
import type { Sandboxes } from '@docker/sandboxes';

export async function listPolicyLogEntries(
  client: Sandboxes,
  sandboxName: string,
  filter: string,
) {
  const sandbox = await client.get(sandboxName);
  return sandbox.networkPolicies.logs.all({ filter }).collect();
}
```

</details>

{{< /tab >}}
{{< /tabs >}}
