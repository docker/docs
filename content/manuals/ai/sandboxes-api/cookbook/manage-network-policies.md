---
title: "Manage network policies"
linkTitle: "Manage network policies"
description: "Create, read, update, and delete personal network policies through the SDK."
keywords: "cloud sandboxes, sandboxes api, manage network policies"
weight: 406
params:
  sidebar:
    group: "Security and policy"
---

Create reusable network policies in your Docker account, then attach their IDs when creating sandboxes. The SDK's governance collection uses the same Docker authentication as the sandbox client.

Start with an [authenticated client](connect-to-cloud-with-a-bearer-token.md). These methods manage your personal policies, not organization-wide administration.

## Manage a policy's lifetime {#1-manage-a-policy-s-lifetime}

The example creates a policy, lists the available policies, reads the new one, and replaces its definition. It attempts to delete the temporary policy in cleanup, including when a later step fails. Check cleanup errors; a cancelled request can leave the policy behind.

Supply a policy definition and a replacement using your SDK's policy input type. Include the destinations your workload needs. For a deny-by-default policy, allow the agent's model provider and any required package registries explicitly.

Set `status: 'POLICY_STATUS_ACTIVE'` on both the policy and its replacement. Although the SDK type makes this field optional, the service rejects omitted or other values with `unimplemented`.

An update replaces the definition rather than appending rules. Read the existing policy before deciding what to retain. Personal-policy mutations are not automatically replayed, so inspect the current state after an uncertain result before trying again.

This is a disposable policy-management example: the returned policy ID has been deleted by the time the function finishes. To keep a policy for real workloads, remove that demonstration cleanup and store its ID with your application configuration. Follow [Control what a sandbox can reach](control-what-a-sandbox-can-reach.md) to attach it and inspect enforced rules.

{{< tabs >}}
{{< tab name="TypeScript" >}}

```typescript
const created = await client.governance.policies.create(policy);
try {
  const listed = await client.governance.policies.list();
  const current = await client.governance.policies.get(created.id);
  const updated = await client.governance.policies.update(
    created.id,
    replacement,
  );
  return { listed, current, updated };
} finally {
  await client.governance.policies.delete(created.id);
}
```

<details>
<summary>Complete TypeScript example: policies/manage.ts</summary>

```typescript
import type { GovernancePolicyInput, Sandboxes } from '@docker/sandboxes';

export async function managePolicy(
  client: Sandboxes,
  policy: GovernancePolicyInput,
  replacement: GovernancePolicyInput,
) {
  const created = await client.governance.policies.create(policy);
  try {
    const listed = await client.governance.policies.list();
    const current = await client.governance.policies.get(created.id);
    const updated = await client.governance.policies.update(
      created.id,
      replacement,
    );
    return { listed, current, updated };
  } finally {
    await client.governance.policies.delete(created.id);
  }
}
```

</details>

{{< /tab >}}
{{< /tabs >}}
