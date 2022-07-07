---
description: Docker extension API reference
keywords: Docker, extensions, sdk, API, reference
skip_read_time: true
---

# Interface: ExecStreamOptions

## Methods

### onOutput

▸ `Optional` **onOutput**(`data`): `void`

Invoked when receiving output from command execution.
By default, the output is split into chunks at arbitrary boundaries.
If you prefer the output to be split into complete lines, set `splitOutputLines`
to true. The callback is then invoked once for each line.

#### Parameters

| Name | Type | Description |
| :------ | :------ | :------ |
| `data` | { `stdout`: `string` ; `stderr?`: `undefined`  } \| { `stdout?`: `undefined` ; `stderr`: `string`  } | output content. Can include either stdout string, or stderr string, one at a time. |

#### Returns

`void`

___

### onError

▸ `Optional` **onError**(`error`): `void`

Invoked to report error if the executed command errors.

#### Parameters

| Name | Type | Description |
| :------ | :------ | :------ |
| `error` | `any` | the error happening in the executed command |

#### Returns

`void`

___

### onClose

▸ `Optional` **onClose**(`exitCode`): `void`

Invoked when process exits.

#### Parameters

| Name | Type | Description |
| :------ | :------ | :------ |
| `exitCode` | `number` | the process exit code |

#### Returns

`void`

## Properties

### splitOutputLines

• `Optional` `Readonly` **splitOutputLines**: `boolean`

Specifies the behaviour invoking `onOutput(data)`. Raw output by default, splitting output at any position. If set to true, `onOutput` will be invoked once for each line.
