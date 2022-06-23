---
skip_read_time: true
description: Docker extension API reference
keywords: Docker, extensions, sdk, API, reference
---

# Interface: ExecResult

## Hierarchy

- [`RawExecResult`](RawExecResult.md)

  ↳ **`ExecResult`**

## Methods

### lines

▸ **lines**(): `string`[]

Split output lines.

#### Returns

`string`[]

The list of lines.

___

### parseJsonLines

▸ **parseJsonLines**(): `any`[]

Parse each output line as a JSON object.

#### Returns

`any`[]

The list of lines where each line is a JSON object.

___

### parseJsonObject

▸ **parseJsonObject**(): `any`

Parse a well-formed JSON output.

#### Returns

`any`

The JSON object.

## Properties

### cmd

• `Optional` `Readonly` **cmd**: `string`

#### Inherited from

[RawExecResult](RawExecResult.md).[cmd](RawExecResult.md#cmd)

___

### killed

• `Optional` `Readonly` **killed**: `boolean`

#### Inherited from

[RawExecResult](RawExecResult.md).[killed](RawExecResult.md#killed)

___

### signal

• `Optional` `Readonly` **signal**: `string`

#### Inherited from

[RawExecResult](RawExecResult.md).[signal](RawExecResult.md#signal)

___

### code

• `Optional` `Readonly` **code**: `number`

#### Inherited from

[RawExecResult](RawExecResult.md).[code](RawExecResult.md#code)

___

### stdout

• `Readonly` **stdout**: `string`

#### Inherited from

[RawExecResult](RawExecResult.md).[stdout](RawExecResult.md#stdout)

___

### stderr

• `Readonly` **stderr**: `string`

#### Inherited from

[RawExecResult](RawExecResult.md).[stderr](RawExecResult.md#stderr)
