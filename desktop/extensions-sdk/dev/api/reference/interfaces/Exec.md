---
description: Docker extension API reference
keywords: Docker, extensions, sdk, API, reference
skip_read_time: true
---

# Interface: Exec

## Callable

### Exec

▸ **Exec**(`cmd`, `args`, `options?`): `Promise`<[`ExecResult`](ExecResult.md)\>

Executes a command.

#### Parameters

| Name | Type | Description |
| :------ | :------ | :------ |
| `cmd` | `string` | The command to execute. |
| `args` | `string`[] | The arguments of the command to execute. |
| `options?` | [`ExecOptions`](ExecOptions.md) | The list of options. |

#### Returns

`Promise`<[`ExecResult`](ExecResult.md)\>

A promise that will resolve once the command finishes.

### Exec

▸ **Exec**(`cmd`, `args`, `options`): [`ExecProcess`](ExecProcess.md)

Streams the result of a command if `stream` is specified in the `options` parameter.

Specify the `stream` if the output of your command is too long or if you need to stream things indefinitely (for example container logs).

#### Parameters

| Name | Type | Description |
| :------ | :------ | :------ |
| `cmd` | `string` | The command to execute. |
| `args` | `string`[] | The arguments of the command to execute. |
| `options` | [`SpawnOptions`](SpawnOptions.md) | The list of options. |

#### Returns

[`ExecProcess`](ExecProcess.md)

The spawned process.
