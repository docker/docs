---
title: "Page through and filter lists"
linkTitle: "Page through and filter lists"
description: "Walk every page of a list with the SDK's iterators, and narrow a list on the server with a filter and an ordering."
keywords: "cloud sandboxes, sandboxes api, page through and filter lists"
weight: 604
params:
  sidebar:
    group: "Requests and responses"
---

List resources without losing results at page boundaries. SDK collection iterators request each page for you; collecting materializes the complete result in memory.

Use an [authenticated client](connect-to-cloud-with-a-bearer-token.md). For a large account, process iterator items as they arrive instead of collecting them all.

## Walk every sandbox {#1-walk-every-sandbox}

Set a page size and iterate the sandbox collection. The example collects the results for convenience. A page size controls each request, not the total number of returned resources.

When using one-page methods directly, pass the returned next-page token unchanged and keep the filter, order, and page size stable. Stop when no next token is returned, not when a page contains fewer items than requested.

{{< tabs >}}
{{< tab name="TypeScript" >}}

```typescript
return client.all({ pageSize }).collect();
```

<details>
<summary>Complete TypeScript example: pagination/walk.ts</summary>

```typescript
import type { Sandboxes } from '@docker/sandboxes';

export async function walkSandboxes(client: Sandboxes, pageSize: number) {
  return client.all({ pageSize }).collect();
}
```

</details>

{{< /tab >}}
{{< /tabs >}}

## Filter and order the result {#2-filter-and-order-the-result}

Pass a filter and ordering supported by the collection. The example uses images; choose a filter such as `status=completed` to select ready images.

Filters are strings interpreted by the service. Use the field names and operators documented for that list method rather than a language object's property names. An invalid filter should be fixed, not silently removed and retried as an unfiltered list.

{{< tabs >}}
{{< tab name="TypeScript" >}}

```typescript
return client.images.all({ filter, orderBy }).collect();
```

<details>
<summary>Complete TypeScript example: pagination/filter.ts</summary>

```typescript
import type { Sandboxes } from '@docker/sandboxes';

export async function filterImages(
  client: Sandboxes,
  filter: string,
  orderBy: string,
) {
  return client.images.all({ filter, orderBy }).collect();
}
```

</details>

{{< /tab >}}
{{< /tabs >}}
