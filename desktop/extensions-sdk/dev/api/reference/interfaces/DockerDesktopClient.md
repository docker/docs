---
title: Docker extension API reference
description: Docker extension API reference
keywords: Docker, extensions, sdk, API, reference
---

# Interface: DockerDesktopClient

When we moved from the v0 to v1 schema, we made sure
window.ddClient satisfied both interfaces. This combined type
describes the resulting API. We should delete it when we stop providing
the v0 API.

## Table of contents

### Properties

- [backend](DockerDesktopClient.md#backend)
- [extension](DockerDesktopClient.md#extension)
- [desktopUI](DockerDesktopClient.md#desktopui)
- [host](DockerDesktopClient.md#host)
- [docker](DockerDesktopClient.md#docker)

### Container Methods

- [listContainers](DockerDesktopClient.md#listcontainers)

### Image Methods

- [listImages](DockerDesktopClient.md#listimages)

### Navigation Methods

- [navigateToContainers](DockerDesktopClient.md#navigatetocontainers)
- [navigateToContainer](DockerDesktopClient.md#navigatetocontainer)
- [navigateToContainerLogs](DockerDesktopClient.md#navigatetocontainerlogs)
- [navigateToContainerInspect](DockerDesktopClient.md#navigatetocontainerinspect)
- [navigateToContainerStats](DockerDesktopClient.md#navigatetocontainerstats)
- [navigateToImages](DockerDesktopClient.md#navigatetoimages)
- [navigateToImage](DockerDesktopClient.md#navigatetoimage)
- [navigateToVolumes](DockerDesktopClient.md#navigatetovolumes)
- [navigateToVolume](DockerDesktopClient.md#navigatetovolume)
- [navigateToDevEnvironments](DockerDesktopClient.md#navigatetodevenvironments)

### Other Methods

- [execHostCmd](DockerDesktopClient.md#exechostcmd)
- [spawnHostCmd](DockerDesktopClient.md#spawnhostcmd)
- [execDockerCmd](DockerDesktopClient.md#execdockercmd)
- [spawnDockerCmd](DockerDesktopClient.md#spawndockercmd)
- [openExternal](DockerDesktopClient.md#openexternal)

### Toast Methods

- [toastSuccess](DockerDesktopClient.md#toastsuccess)
- [toastWarning](DockerDesktopClient.md#toastwarning)
- [toastError](DockerDesktopClient.md#toasterror)

## Properties

### backend

• `Readonly` **backend**: `undefined` \| [`BackendV0`](BackendV0.md)

The `window.ddClient.backend` object can be used to communicate with the backend defined in the vm section of
the extension metadata.
The client is already connected to the backend.

**`deprecated`** :warning: It will be removed in a future version. Use [DockerDesktopClient.extension](DockerDesktopClient.md#extension) instead.

#### Inherited from

DockerDesktopClientV0.backend

---

### extension

• `Readonly` **extension**: [`Extension`](Extension.md)

The `ddClient.extension` object can be used to communicate with the backend defined in the vm section of the
extension metadata.
The client is already connected to the backend.

#### Inherited from

DockerDesktopClientV1.extension

---

### desktopUI

• `Readonly` **desktopUI**: [`DesktopUI`](DesktopUI.md)

#### Inherited from

DockerDesktopClientV1.desktopUI

---

### host

• `Readonly` **host**: [`Host`](Host.md)

#### Inherited from

DockerDesktopClientV1.host

---

### docker

• `Readonly` **docker**: [`Docker`](Docker.md)

#### Inherited from

DockerDesktopClientV1.docker

## Container Methods

### listContainers

▸ **listContainers**(`options`): `Promise`<`unknown`\>

Get the list of running containers (same as `docker ps`).

By default, this will not list stopped containers.
You can use the option `{"all": true}` to list all the running and stopped containers.

```typescript
const containers = await window.ddClient.listContainers();
```

**`deprecated`** :warning: It will be removed in a future version. Use [Docker.listContainers](Docker.md#listcontainers) instead.

#### Parameters

| Name      | Type    | Description                                                                                                                                                                                                                                                                                  |
| :-------- | :------ | :------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `options` | `never` | (Optional). A JSON like `{ "all": true, "limit": 10, "size": true, "filters": JSON.stringify({ status: ["exited"] }), }` For more information about the different properties see [the Docker API endpoint documentation](https://docs.docker.com/engine/api/v1.37/#operation/ContainerList). |

#### Returns

`Promise`<`unknown`\>

#### Inherited from

DockerDesktopClientV0.listContainers

---

## Image Methods

### listImages

▸ **listImages**(`options`): `Promise`<`unknown`\>

Get the list of images

```typescript
const images = await window.ddClient.listImages();
```

**`deprecated`** :warning: It will be removed in a future version. Use [Docker.listImages](Docker.md#listimages) instead.

#### Parameters

| Name      | Type    | Description                                                                                                                                                                                                                                                         |
| :-------- | :------ | :------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| `options` | `never` | (Optional). A JSON like `{ "all": true, "filters": JSON.stringify({ dangling: ["true"] }), "digests": true }` For more information about the different properties see [the Docker API endpoint documentation](https://docs.docker.com/engine/api/v1.37/#tag/Image). |

#### Returns

`Promise`<`unknown`\>

#### Inherited from

DockerDesktopClientV0.listImages

---

## Navigation Methods

### navigateToContainers

▸ **navigateToContainers**(): `void`

Navigate to the containers window in Docker Desktop.

```typescript
window.ddClient.navigateToContainers();
```

**`deprecated`** :warning: It will be removed in a future version. Use [NavigationIntents.viewContainers](NavigationIntents.md#viewcontainers) instead.

#### Returns

`void`

#### Inherited from

DockerDesktopClientV0.navigateToContainers

---

### navigateToContainer

▸ **navigateToContainer**(`id`): `Promise`<`any`\>

Navigate to the container window in Docker Desktop.

```typescript
await window.ddClient.navigateToContainer(id);
```

**`deprecated`** :warning: It will be removed in a future version. Use [NavigationIntents.viewContainer](NavigationIntents.md#viewcontainer) instead.

#### Parameters

| Name | Type     | Description                                                                                                                                                                                            |
| :--- | :------- | :----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `id` | `string` | The full container id, e.g. `46b57e400d801762e9e115734bf902a2450d89669d85881058a46136520aca28`. You can use the `--no-trunc` flag as part of the `docker ps` command to display the full container id. |

#### Returns

`Promise`<`any`\>

A promise that fails if the container doesn't exist.

#### Inherited from

DockerDesktopClientV0.navigateToContainer

---

### navigateToContainerLogs

▸ **navigateToContainerLogs**(`id`): `Promise`<`any`\>

Navigate to the container logs window in Docker Desktop.

```typescript
await window.ddClient.navigateToContainerLogs(id);
```

**`deprecated`** :warning: It will be removed in a future version. Use {@link DockerDesktopClient.viewContainerLogs} instead.

#### Parameters

| Name | Type     | Description                                                                                                                                                                                            |
| :--- | :------- | :----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `id` | `string` | The full container id, e.g. `46b57e400d801762e9e115734bf902a2450d89669d85881058a46136520aca28`. You can use the `--no-trunc` flag as part of the `docker ps` command to display the full container id. |

#### Returns

`Promise`<`any`\>

A promise that fails if the container doesn't exist.

#### Inherited from

DockerDesktopClientV0.navigateToContainerLogs

---

### navigateToContainerInspect

▸ **navigateToContainerInspect**(`id`): `Promise`<`any`\>

Navigate to the container inspect window in Docker Desktop.

```typescript
await window.ddClient.navigateToContainerInspect(id);
```

**`deprecated`** :warning: It will be removed in a future version. Use {@link DockerDesktopClient.viewContainerInspect} instead.

#### Parameters

| Name | Type     | Description                                                                                                                                                                                            |
| :--- | :------- | :----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `id` | `string` | The full container id, e.g. `46b57e400d801762e9e115734bf902a2450d89669d85881058a46136520aca28`. You can use the `--no-trunc` flag as part of the `docker ps` command to display the full container id. |

#### Returns

`Promise`<`any`\>

A promise that fails if the container doesn't exist.

#### Inherited from

DockerDesktopClientV0.navigateToContainerInspect

---

### navigateToContainerStats

▸ **navigateToContainerStats**(`id`): `Promise`<`any`\>

Navigate to the container stats to see the CPU, memory, disk read/write and network I/O usage.

```typescript
await window.ddClient.navigateToContainerStats(id);
```

**`deprecated`** :warning: It will be removed in a future version. Use {@link DockerDesktopClient.viewContainerStats} instead.

#### Parameters

| Name | Type     | Description                                                                                                                                                                                            |
| :--- | :------- | :----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `id` | `string` | The full container id, e.g. `46b57e400d801762e9e115734bf902a2450d89669d85881058a46136520aca28`. You can use the `--no-trunc` flag as part of the `docker ps` command to display the full container id. |

#### Returns

`Promise`<`any`\>

A promise that fails if the container doesn't exist.

#### Inherited from

DockerDesktopClientV0.navigateToContainerStats

---

### navigateToImages

▸ **navigateToImages**(): `void`

Navigate to the images window in Docker Desktop.

```typescript
await window.ddClient.navigateToImages(id);
```

**`deprecated`** :warning: It will be removed in a future version. Use [NavigationIntents.viewImages](NavigationIntents.md#viewimages) instead.

#### Returns

`void`

#### Inherited from

DockerDesktopClientV0.navigateToImages

---

### navigateToImage

▸ **navigateToImage**(`id`, `tag`): `Promise`<`any`\>

Navigate to a specific image referenced by `id` and `tag` in Docker Desktop.
In this navigation route you can find the image layers, commands, created time and size.

```typescript
await window.ddClient.navigateToImage(id, tag);
```

**`deprecated`** :warning: It will be removed in a future version. Use [NavigationIntents.viewImage](NavigationIntents.md#viewimage) instead.

#### Parameters

| Name  | Type     | Description                                                                                                        |
| :---- | :------- | :----------------------------------------------------------------------------------------------------------------- |
| `id`  | `string` | The full image id (including sha), e.g. `sha256:34ab3ae068572f4e85c448b4035e6be5e19cc41f69606535cd4d768a63432673`. |
| `tag` | `string` | The tag of the image, e.g. `latest`, `0.0.1`, etc.                                                                 |

#### Returns

`Promise`<`any`\>

A promise that fails if the container doesn't exist.

#### Inherited from

DockerDesktopClientV0.navigateToImage

---

### navigateToVolumes

▸ **navigateToVolumes**(): `void`

Navigate to the volumes window in Docker Desktop.

```typescript
await window.ddClient.navigateToVolumes();
```

**`deprecated`** :warning: It will be removed in a future version. Use [NavigationIntents.viewVolumes](NavigationIntents.md#viewvolumes) instead.

#### Returns

`void`

#### Inherited from

DockerDesktopClientV0.navigateToVolumes

---

### navigateToVolume

▸ **navigateToVolume**(`volume`): `void`

Navigate to a specific volume in Docker Desktop.

```typescript
window.ddClient.navigateToVolume(volume);
```

**`deprecated`** :warning: It will be removed in a future version. Use [NavigationIntents.viewVolume](NavigationIntents.md#viewvolume) instead.

#### Parameters

| Name     | Type     | Description                               |
| :------- | :------- | :---------------------------------------- |
| `volume` | `string` | The name of the volume, e.g. `my-volume`. |

#### Returns

`void`

#### Inherited from

DockerDesktopClientV0.navigateToVolume

---

### navigateToDevEnvironments

▸ **navigateToDevEnvironments**(): `void`

Navigate to the Dev Environments window in Docker Desktop.

```typescript
window.ddClient.navigateToDevEnvironments();
```

**`deprecated`** :warning: It will be removed in a future version. Use [NavigationIntents.viewDevEnvironments](NavigationIntents.md#viewdevenvironments) instead.

#### Returns

`void`

#### Inherited from

DockerDesktopClientV0.navigateToDevEnvironments

---

## Other Methods

### execHostCmd

▸ **execHostCmd**(`cmd`): `Promise`<[`ExecResultV0`](ExecResultV0.md)\>

You can run binaries defined in the host section in the extension metadata.

```typescript
window.ddClient.execHostCmd(`cliShippedOnHost xxx`).then((cmdResult: any) => {
  console.log(cmdResult);
});
```

**`deprecated`** :warning: It will be removed in a future version. Use [ExtensionCli.exec](ExtensionCli.md#exec) instead.

#### Parameters

| Name  | Type     | Description                 |
| :---- | :------- | :-------------------------- |
| `cmd` | `string` | The command to be executed. |

#### Returns

`Promise`<[`ExecResultV0`](ExecResultV0.md)\>

#### Inherited from

DockerDesktopClientV0.execHostCmd

---

### spawnHostCmd

▸ **spawnHostCmd**(`cmd`, `args`, `callback`): `void`

Invoke an extension binary on your host and getting the output stream.

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

**`deprecated`** :warning: It will be removed in a future version. Use [ExtensionCli.exec](ExtensionCli.md#exec) instead.

#### Parameters

| Name       | Type                                      | Description                                                                    |
| :--------- | :---------------------------------------- | :----------------------------------------------------------------------------- |
| `cmd`      | `string`                                  | The command to be executed.                                                    |
| `args`     | `string`[]                                | The arguments of the command to execute.                                       |
| `callback` | (`data`: `any`, `error`: `any`) => `void` | The callback function where to listen from the command output data and errors. |

#### Returns

`void`

#### Inherited from

DockerDesktopClientV0.spawnHostCmd

---

### execDockerCmd

▸ **execDockerCmd**(`cmd`, ...`args`): `Promise`<[`ExecResultV0`](ExecResultV0.md)\>

You can also directly execute the docker binary.

```typescript
const output = await window.ddClient.execDockerCmd(
  "info",
  "--format",
  {% raw %}'"{{ json . }}"'{% endraw %}
);
```

**`deprecated`** :warning: It will be removed in a future version. Use [DockerCommand.exec](DockerCommand.md#exec) instead.

#### Parameters

| Name      | Type       | Description                              |
| :-------- | :--------- | :--------------------------------------- |
| `cmd`     | `string`   | The command to execute.                  |
| `...args` | `string`[] | The arguments of the command to execute. |

#### Returns

`Promise`<[`ExecResultV0`](ExecResultV0.md)\>

The result will contain both the standard output and the standard error of the executed command:

```
{
  "stderr": "...",
  "stdout": "..."
}
```

In this example the docker command output is a json output.

For convenience, the command result object also has methods to easily parse it:

- `output.lines(): string[]` splits output lines.
- `output.parseJsonObject(): any` parses a well-formed json output.
- `output.parseJsonLines(): any[]` parses each output line as a json object.

If the output of the command is too long, or you need to get the output as a stream you can use the
spawnDockerCmd function:

```typescript
window.ddClient.spawnDockerCmd("logs", ["-f", "..."], (data, error) => {
  console.log(data.stdout);
});
```

#### Inherited from

DockerDesktopClientV0.execDockerCmd

---

### spawnDockerCmd

▸ **spawnDockerCmd**(`cmd`, `args`, `callback`): `void`

**`deprecated`** :warning: It will be removed in a future version. Use [DockerCommand.exec](DockerCommand.md#exec) instead.

#### Parameters

| Name       | Type                                      |
| :--------- | :---------------------------------------- |
| `cmd`      | `string`                                  |
| `args`     | `string`[]                                |
| `callback` | (`data`: `any`, `error`: `any`) => `void` |

#### Returns

`void`

#### Inherited from

DockerDesktopClientV0.spawnDockerCmd

---

### openExternal

▸ **openExternal**(`url`): `void`

Opens an external URL with the system default browser.

```typescript
window.ddClient.openExternal("https://docker.com");
```

**`deprecated`** :warning: It will be removed in a future version. Use [Host.openExternal](Host.md#openexternal) instead.

#### Parameters

| Name  | Type     | Description                                                               |
| :---- | :------- | :------------------------------------------------------------------------ |
| `url` | `string` | The URL the browser will open (must have the protocol `http` or `https`). |

#### Returns

`void`

#### Inherited from

DockerDesktopClientV0.openExternal

---

## Toast Methods

### toastSuccess

▸ **toastSuccess**(`msg`): `void`

Display a toast message of type success.

```typescript
window.ddClient.toastSuccess("message");
```

**`deprecated`** :warning: It will be removed in a future version. Use [Toast.success](Toast.md#success) instead.

#### Parameters

| Name  | Type     | Description                          |
| :---- | :------- | :----------------------------------- |
| `msg` | `string` | The message to display in the toast. |

#### Returns

`void`

#### Inherited from

DockerDesktopClientV0.toastSuccess

---

### toastWarning

▸ **toastWarning**(`msg`): `void`

Display a toast message of type warning.

```typescript
window.ddClient.toastWarning("message");
```

**`deprecated`** :warning: It will be removed in a future version. Use [Toast.warning](Toast.md#warning) instead.

#### Parameters

| Name  | Type     | Description                          |
| :---- | :------- | :----------------------------------- |
| `msg` | `string` | The message to display in the toast. |

#### Returns

`void`

#### Inherited from

DockerDesktopClientV0.toastWarning

---

### toastError

▸ **toastError**(`msg`): `void`

Display a toast message of type error.

```typescript
window.ddClient.toastError("message");
```

**`deprecated`** :warning: It will be removed in a future version. Use [Toast.error](Toast.md#error) instead.

#### Parameters

| Name  | Type     | Description                          |
| :---- | :------- | :----------------------------------- |
| `msg` | `string` | The message to display in the toast. |

#### Returns

`void`

#### Inherited from

DockerDesktopClientV0.toastError
