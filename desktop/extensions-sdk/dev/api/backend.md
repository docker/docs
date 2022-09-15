---
title: Extension Backend
description: Docker extension API
keywords: Docker, extensions, sdk, API
---

The `ddClient.extension.vm` object can be used to communicate with the backend defined in the [vm section](../../extensions/METADATA.md#vm-section) of the extension metadata.

## get

▸ **get**(`url`): `Promise`<`unknown`\>

Performs an HTTP GET request to a backend service.

```typescript
ddClient.extension.vm.service
 .get("/some/service")
 .then((value: any) => console.log(value)
```

See [Service API Reference](reference/interfaces/HttpService.md) for other methods such as POST, UPDATE, and DELETE.

> Deprecated extension backend communication
>
> The methods below that use `window.ddClient.backend` are deprecated and will be removed in a future version. Use the methods specified above.

The `window.ddClient.backend` object can be used to communicate with the backend
defined in the [vm section](../../extensions/METADATA.md#vm-section) of the
extension metadata. The client is already connected to the backend.

Example usages:

```typescript
window.ddClient.backend
  .get("/some/service")
  .then((value: any) => console.log(value));

window.ddClient.backend
  .post("/some/service", { ... })
  .then((value: any) => console.log(value));

window.ddClient.backend
  .put("/some/service", { ... })
  .then((value: any) => console.log(value));

window.ddClient.backend
  .patch("/some/service", { ... })
  .then((value: any) => console.log(value));

window.ddClient.backend
  .delete("/some/service")
  .then((value: any) => console.log(value));

window.ddClient.backend
  .head("/some/service")
  .then((value: any) => console.log(value));

window.ddClient.backend
  .request({ url: "/url", method: "GET", headers: { 'header-key': 'header-value' }, data: { ... }})
  .then((value: any) => console.log(value));
```

## Run a command in the extension backend container

For example, execute the command `ls -l` inside the backend container:

```typescript
await ddClient.extension.vm.cli.exec("ls", ["-l"]);
```

Stream the output of the command executed in the backend container. For example, spawn the command `ls -l` inside the backend container:

```typescript
await ddClient.extension.vm.cli.exec("ls", ["-l"], {
  stream: {
    onOutput(data) {
      if (data.stdout) {
        console.error(data.stdout);
      } else {
        console.log(data.stderr);
      }
    },
    onError(error) {
      console.error(error);
    },
    onClose(exitCode) {
      console.log("onClose with exit code " + exitCode);
    },
  },
});
```

For more details, refer to the [Extension VM API Reference](reference/interfaces/ExtensionVM.md)

> Deprecated extension backend command execution
>
> This method is deprecated and will be removed in a future version. Use the specified method above.

If your extension ships with additional binaries that should be run inside the
backend container, you can use the `execInVMExtension` function:

```typescript
const output = await window.ddClient.backend.execInVMExtension(
  `cliShippedInTheVm xxx`
);
console.log(output);
```

## Invoke an extension binary on the host

You can run binaries defined in the [host section](../../extensions/METADATA.md#host-section)
of the extension metadata.

For example, execute the shipped binary `kubectl -h` command in the host:

```typescript
await ddClient.extension.host.cli.exec("kubectl", ["-h"]);
```

As long as the `kubectl` binary is shipped as part of your extension, you can spawn the `kubectl -h` command in the host and get the output stream:

```typescript
await ddClient.extension.host.cli.exec("kubectl", ["-h"], {
  stream: {
    onOutput(data: { stdout: string } | { stderr: string }): void {
      if (data.stdout) {
        console.error(data.stdout);
      } else {
        console.log(data.stderr);
      }
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

You can stream the output of the command executed in the backend container or in the host.

For more details, refer to the [Extension Host API Reference](reference/interfaces/ExtensionHost.md)

> Deprecated invocation of extension binary
>
> This method is deprecated and will be removed in a future version. Use the method specified above.

To execute a command in the host:

```typescript
window.ddClient.execHostCmd(`cliShippedOnHost xxx`).then((cmdResult: any) => {
  console.log(cmdResult);
});
```

To stream the output of the command executed in the backend container or in the host:

```typescript
window.ddClient.spawnHostCmd(
  `cliShippedOnHost`,
  [`arg1`, `arg2`],
  (data: any, err: any) => {
    console.log(data.stdout, data.stderr);
    // Once the command exits we get the status code
    if (data.code) {
      console.log(data.code);
    }
  }
);
```

> You cannot use this to chain commands in a single `exec()` invocation (like `cmd1 $(cmd2)` or using pipe between commands).
>
> You need to invoke `exec()` for each command and parse results to pass parameters to the next command if needed.
