---
description: Docker extension API reference
keywords: Docker, extensions, sdk, API, reference
skip_read_time: true
---

# Interface: ExtensionVM

**`since`** 0.2.0

## Properties

### cli

• `Readonly` **cli**: [`ExtensionCli`](ExtensionCli.md)

Executes a command in the backend container.

Example: Execute the command `ls -l` inside the **backend container**:

```typescript
await ddClient.extension.vm.cli.exec(
  "ls",
  ["-l"]
);
```

Streams the output of the command executed in the backend container.

When the extension defines its own `docker-compose.yaml` file
with multiple containers, the command is executed on the first container defined.
Change the order in which containers are defined to execute commands on another
container.

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

___

### service

• `Optional` `Readonly` **service**: [`HttpService`](HttpService.md)
