---
title: "Read files out of a sandbox"
linkTitle: "Read files out of a sandbox"
description: "List a directory in a running cloud sandbox, download files from it into memory, then move or remove paths inside it."
keywords: "cloud sandboxes, sandboxes api, read files out of a sandbox"
weight: 203
params:
  sidebar:
    group: "Working in a sandbox"
---

Read build results, inspect project files, or copy data out before deleting a sandbox. Start with a running sandbox handle and an absolute path inside it.

The examples return content to your application. To keep it on your machine, write the received bytes to a local file.

## List a directory {#1-list-a-directory}

Walk the directory through the file collection's iterator. It follows pagination and returns the entries. Listing gives metadata, not file content.

{{< tabs >}}
{{< tab name="TypeScript" >}}

```typescript
return sandbox.files.all(path).collect();
```

<details>
<summary>Complete TypeScript example: download/list.ts</summary>

```typescript
import type { Sandbox } from '@docker/sandboxes';

export async function listDirectory(sandbox: Sandbox, path: string) {
  return sandbox.files.all(path).collect();
}
```

</details>

{{< /tab >}}
{{< /tabs >}}

## Read a small text file {#2-read-a-small-text-file}

Use the bounded read helper for text you want in memory. The example limits the read to one MiB. Choose a bound that fits your application's memory budget, or download a larger file as a stream.

{{< tabs >}}
{{< tab name="TypeScript" >}}

```typescript
return sandbox.files.read(path, {
  encoding: 'utf8',
  maxBytes: 1024 * 1024,
});
```

<details>
<summary>Complete TypeScript example: download/read.ts</summary>

```typescript
import type { Sandbox } from '@docker/sandboxes';

export async function readFile(sandbox: Sandbox, path: string) {
  return sandbox.files.read(path, {
    encoding: 'utf8',
    maxBytes: 1024 * 1024,
  });
}
```

</details>

{{< /tab >}}
{{< /tabs >}}

## Inspect a path {#3-inspect-a-path}

Read metadata before deciding whether to download, move, or remove a path. A successful metadata read does not reserve the file: another process may change it afterward.

{{< tabs >}}
{{< tab name="TypeScript" >}}

```typescript
return sandbox.files.stat(path);
```

<details>
<summary>Complete TypeScript example: download/stat.ts</summary>

```typescript
import type { Sandbox } from '@docker/sandboxes';

export async function statFile(sandbox: Sandbox, path: string) {
  return sandbox.files.stat(path);
}
```

</details>

{{< /tab >}}
{{< /tabs >}}

## Download a file as a stream {#4-download-a-file-as-a-stream}

Pass the sandbox path and a callback or writer that consumes bytes. The example closes the transfer when it finishes or fails.

Treat the download as complete only when it ends successfully. If it fails halfway through, discard or separately identify the partial local file before retrying.

{{< tabs >}}
{{< tab name="TypeScript" >}}

```typescript
const download = await sandbox.files.download(path);
try {
  for await (const chunk of download) write(chunk);
} finally {
  await download.close();
}
```

<details>
<summary>Complete TypeScript example: download/download.ts</summary>

```typescript
import type { Sandbox } from '@docker/sandboxes';

export async function downloadFiles(
  sandbox: Sandbox,
  path: string,
  write: (bytes: Uint8Array) => void,
) {
  const download = await sandbox.files.download(path);
  try {
    for await (const chunk of download) write(chunk);
  } finally {
    await download.close();
  }
}
```

</details>

{{< /tab >}}
{{< /tabs >}}

## Move a path {#5-move-a-path}

Pass the current path and destination. This moves data inside the sandbox; it does not download anything to your computer.

{{< tabs >}}
{{< tab name="TypeScript" >}}

```typescript
await sandbox.files.move(from, to);
```

<details>
<summary>Complete TypeScript example: download/move.ts</summary>

```typescript
import type { Sandbox } from '@docker/sandboxes';

export async function movePath(
  sandbox: Sandbox,
  from: string,
  to: string,
) {
  await sandbox.files.move(from, to);
}
```

</details>

{{< /tab >}}
{{< /tabs >}}

## Remove a path {#6-remove-a-path}

Remove a file, or enable recursive removal for a directory tree. Check the result for a failed path rather than assuming that every requested removal succeeded.

Recursive removal is destructive. Keep user-supplied paths constrained to the directory your application owns.

{{< tabs >}}
{{< tab name="TypeScript" >}}

```typescript
const result = await sandbox.files.remove(path, { recursive });
if (result.failedPath)
  throw new Error(`Remove ${path} stopped at ${result.failedPath}`);
```

<details>
<summary>Complete TypeScript example: download/remove.ts</summary>

```typescript
import type { Sandbox } from '@docker/sandboxes';

export async function removePath(
  sandbox: Sandbox,
  path: string,
  recursive: boolean,
) {
  const result = await sandbox.files.remove(path, { recursive });
  if (result.failedPath)
    throw new Error(`Remove ${path} stopped at ${result.failedPath}`);
}
```

</details>

{{< /tab >}}
{{< /tabs >}}
