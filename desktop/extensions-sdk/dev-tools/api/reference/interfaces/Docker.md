---
title: Docker 
description: Docker extension API reference
keywords: Docker, extensions, sdk, API, reference
---

## Properties

### cli

• `Readonly` **cli**: [`DockerCommand`](DockerCommand.md)

You can also directly execute the docker binary.

```typescript
const output = await ddClient.docker.cli.exec("info", [
  "--format",
  {% raw %}'"{{ json . }}"',{% endraw %}
]);
```

Output:

```
{
  "stderr": "...",
  "stdout": "..."
}
```

In this example the docker command output is a json output.
For convenience, the command result object also has methods to easily parse it. See [ExecResult](ExecResult.md) instead.

---

Streams the output as a result of the execution of a docker command.
It is useful when the output of the command is too long, or you need to get the output as a stream.

```typescript
await ddClient.docker.cli.exec("logs", ["-f", "..."], {
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

## Methods

### listContainers

▸ **listContainers**(`options?`): `Promise`<`unknown`\>

Get the list of running containers (same as `docker ps`).

By default, this will not list stopped containers.
You can use the option `{"all": true}` to list all the running and stopped containers.

```typescript
const containers = await ddClient.docker.listContainers();
```

#### Parameters

| Name       | Type  | Description                                                                                                                                                                                                                                                                                  |
| :--------- | :---- | :------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `options?` | `any` | (Optional). A JSON like `{ "all": true, "limit": 10, "size": true, "filters": JSON.stringify({ status: ["exited"] }), }` For more information about the different properties see [the Docker API endpoint documentation](https://docs.docker.com/engine/api/v1.37/#operation/ContainerList). |

#### Returns

`Promise`<`unknown`\>

---

### listImages

▸ **listImages**(`options?`): `Promise`<`unknown`\>

Get the list of local container images

```typescript
const images = await ddClient.docker.listImages();
```

#### Parameters

| Name       | Type  | Description                                                                                                                                                                                                                                                         |
| :--------- | :---- | :------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| `options?` | `any` | (Optional). A JSON like `{ "all": true, "filters": JSON.stringify({ dangling: ["true"] }), "digests": true }` For more information about the different properties see [the Docker API endpoint documentation](https://docs.docker.com/engine/api/v1.37/#tag/Image). |

#### Returns

`Promise`<`unknown`\>
