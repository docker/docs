---
title: Docker extension API reference
description: Docker extension API reference
keywords: Docker, extensions, sdk, API, reference
---

# Interface: ExtensionVM

## Table of contents

### Properties

- [cli](ExtensionVM.md#cli)
- [service](ExtensionVM.md#service)

## Properties

### cli

• `Readonly` **cli**: [`ExtensionCli`](ExtensionCli.md)

Executes a command in the backend container.

Example: Execute the command `ls -l` inside the **backend container**:

```typescript
await ddClient.extension.vm.cli.exec("ls", ["-l"]);
```

Streams the output of the command executed in the backend container.

Example: Spawn the command `ls -l` inside the **backend container**:

```typescript
await ddClient.extension.vm.cli.exec("ls", ["-l"], {
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

**`param`** Command to execute.

**`param`** Arguments of the command to execute.

**`param`** The callback function where to listen from the command output data and errors.

---

### service

• `Optional` `Readonly` **service**: [`HttpService`](HttpService.md)
