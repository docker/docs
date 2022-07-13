---
description: Docker extension API reference
keywords: Docker, extensions, sdk, API, reference
skip_read_time: true
---

# Interface: BackendV0

## Container Methods

### execInContainer

▸ **execInContainer**(`container`, `cmd`): `Promise`<[`ExecResultV0`](ExecResultV0.md)\>

Executes a command inside a container.

```typescript
const output = await window.ddClient.backend.execInContainer(container, cmd);

console.log(output);
```

**`deprecated`** :warning: It will be removed in a future version.

#### Parameters

| Name | Type | Description |
| :------ | :------ | :------ |
| `container` | `string` | - |
| `cmd` | `string` | The command to be executed. |

#### Returns

`Promise`<[`ExecResultV0`](ExecResultV0.md)\>

___

## HTTP Methods

### get

▸ **get**(`url`): `Promise`<`unknown`\>

Performs an HTTP GET request to a backend service.

```typescript
window.ddClient.backend
 .get("/some/service")
 .then((value: any) => console.log(value));
```

**`deprecated`** :warning: It will be removed in a future version. Use [HttpService.get](HttpService.md#get) instead.

#### Parameters

| Name | Type | Description |
| :------ | :------ | :------ |
| `url` | `string` | The URL of the backend service. |

#### Returns

`Promise`<`unknown`\>

___

### post

▸ **post**(`url`, `data`): `Promise`<`unknown`\>

Performs an HTTP POST request to a backend service.

```typescript
window.ddClient.backend
 .post("/some/service", { ... })
 .then((value: any) => console.log(value));
```

**`deprecated`** :warning: It will be removed in a future version. Use [HttpService.post](HttpService.md#post) instead.

#### Parameters

| Name | Type | Description |
| :------ | :------ | :------ |
| `url` | `string` | The URL of the backend service. |
| `data` | `any` | The body of the request. |

#### Returns

`Promise`<`unknown`\>

___

### put

▸ **put**(`url`, `data`): `Promise`<`unknown`\>

Performs an HTTP PUT request to a backend service.

```typescript
window.ddClient.backend
 .put("/some/service", { ... })
 .then((value: any) => console.log(value));
```

**`deprecated`** :warning: It will be removed in a future version. Use [HttpService.put](HttpService.md#put) instead.

#### Parameters

| Name | Type | Description |
| :------ | :------ | :------ |
| `url` | `string` | The URL of the backend service. |
| `data` | `any` | The body of the request. |

#### Returns

`Promise`<`unknown`\>

___

### patch

▸ **patch**(`url`, `data`): `Promise`<`unknown`\>

Performs an HTTP PATCH request to a backend service.

```typescript
window.ddClient.backend
 .patch("/some/service", { ... })
 .then((value: any) => console.log(value));
```

**`deprecated`** :warning: It will be removed in a future version. Use [HttpService.patch](HttpService.md#patch) instead.

#### Parameters

| Name | Type | Description |
| :------ | :------ | :------ |
| `url` | `string` | The URL of the backend service. |
| `data` | `any` | The body of the request. |

#### Returns

`Promise`<`unknown`\>

___

### delete

▸ **delete**(`url`): `Promise`<`unknown`\>

Performs an HTTP DELETE request to a backend service.

```typescript
window.ddClient.backend
 .delete("/some/service")
 .then((value: any) => console.log(value));
```

**`deprecated`** :warning: It will be removed in a future version. Use [HttpService.delete](HttpService.md#delete) instead.

#### Parameters

| Name | Type | Description |
| :------ | :------ | :------ |
| `url` | `string` | The URL of the backend service. |

#### Returns

`Promise`<`unknown`\>

___

### head

▸ **head**(`url`): `Promise`<`unknown`\>

Performs an HTTP HEAD request to a backend service.

```typescript
window.ddClient.backend
 .head("/some/service")
 .then((value: any) => console.log(value));
```

**`deprecated`** :warning: It will be removed in a future version. Use [HttpService.head](HttpService.md#head) instead.

#### Parameters

| Name | Type | Description |
| :------ | :------ | :------ |
| `url` | `string` | The URL of the backend service. |

#### Returns

`Promise`<`unknown`\>

___

### request

▸ **request**(`config`): `Promise`<`unknown`\>

Performs an HTTP request to a backend service.

```typescript
window.ddClient.backend
 .request({ url: "/url", method: "GET", headers: { 'header-key': 'header-value' }, data: { ... }})
 .then((value: any) => console.log(value));
```

**`deprecated`** :warning: It will be removed in a future version. Use [HttpService.request](HttpService.md#request) instead.

#### Parameters

| Name | Type | Description |
| :------ | :------ | :------ |
| `config` | [`RequestConfigV0`](RequestConfigV0.md) | The URL of the backend service. |

#### Returns

`Promise`<`unknown`\>

___

## VM Methods

### execInVMExtension

▸ **execInVMExtension**(`cmd`): `Promise`<[`ExecResultV0`](ExecResultV0.md)\>

Executes a command inside the backend container.
If your extensions ships with additional binaries that should be run inside the backend container you can use the `execInVMExtension` function.

```typescript
const output = await window.ddClient.backend.execInVMExtension(
  `cliShippedInTheVm xxx`
);

console.log(output);
```

**`deprecated`** :warning: It will be removed in a future version. Use [ExtensionCli.exec](ExtensionCli.md#exec) instead.

#### Parameters

| Name | Type | Description |
| :------ | :------ | :------ |
| `cmd` | `string` | The command to be executed. |

#### Returns

`Promise`<[`ExecResultV0`](ExecResultV0.md)\>

___

### spawnInVMExtension

▸ **spawnInVMExtension**(`cmd`, `args`, `callback`): `void`

Returns a stream from the command executed in the backend container.

```typescript
window.ddClient.spawnInVMExtension(
  `cmd`,
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

**`deprecated`** :warning: It will be removed in a future version. Use {@link ExtensionCli.spawn} instead.

#### Parameters

| Name | Type | Description |
| :------ | :------ | :------ |
| `cmd` | `string` | The command to be executed. |
| `args` | `string`[] | The arguments of the command to execute. |
| `callback` | (`data`: `any`, `error`: `any`) => `void` | The callback function where to listen from the command output data and errors. |

#### Returns

`void`
