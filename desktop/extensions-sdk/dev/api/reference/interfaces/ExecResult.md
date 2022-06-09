---
title: Docker extension API reference
description: Docker extension API reference
keywords: Docker, extensions, sdk, API, reference
---

# Interface: ExecResult

## Hierarchy

- [`RawExecResult`](RawExecResult.md)

  ↳ **`ExecResult`**

## Table of contents

### Methods

- [lines](ExecResult.md#lines)
- [parseJsonLines](ExecResult.md#parsejsonlines)
- [parseJsonObject](ExecResult.md#parsejsonobject)

### Properties

- [cmd](ExecResult.md#cmd)
- [killed](ExecResult.md#killed)
- [signal](ExecResult.md#signal)
- [code](ExecResult.md#code)
- [stdout](ExecResult.md#stdout)
- [stderr](ExecResult.md#stderr)

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
