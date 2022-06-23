---
skip_read_time: true
description: Docker extension API reference
keywords: Docker, extensions, sdk, API, reference
---

# Interface: ExtensionHost

## Properties

### cli

• `Readonly` **cli**: [`ExtensionCli`](ExtensionCli.md)

Executes a command in the host.

For example, execute the shipped binary `kubectl -h` command in the **host**:

```typescript
await ddClient.extension.host.cli.exec("kubectl", ["-h"]);
```

---

Streams the output of the command executed in the backend container or in the host.

Provided the `kubectl` binary is shipped as part of your extension, you can spawn the `kubectl -h` command in the **host**:

```typescript
await ddClient.extension.host.cli.exec("kubectl", ["-h"], {
  stream: {
    onOutput(data): void {
      // As we can receive both `stdout` and `stderr`, we wrap them in a JSON object
      JSON.stringify(
        {
          stdout: data.stdout,
          stderr: data.stderr,
        },
        null,
        "  "
      );
    },
    onError(error: any): void {
      console.error(error);
    },
    onClose(exitCode: number): void {
      console.log("onClose with exit code " + exitCode);
    },
  },
});
```
