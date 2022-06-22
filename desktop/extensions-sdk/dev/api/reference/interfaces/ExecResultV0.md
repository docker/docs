---
title: Docker extension API reference
description: Docker extension API reference
keywords: Docker, extensions, sdk, API, reference
---

# Interface: ExecResultV0

## Table of contents

### Properties

- [cmd](ExecResultV0.md#cmd)
- [killed](ExecResultV0.md#killed)
- [signal](ExecResultV0.md#signal)
- [code](ExecResultV0.md#code)
- [stdout](ExecResultV0.md#stdout)
- [stderr](ExecResultV0.md#stderr)

### Methods

- [lines](ExecResultV0.md#lines)
- [parseJsonLines](ExecResultV0.md#parsejsonlines)
- [parseJsonObject](ExecResultV0.md#parsejsonobject)

## Properties

### cmd

• `Optional` `Readonly` **cmd**: `string`

___

### killed

• `Optional` `Readonly` **killed**: `boolean`

___

### signal

• `Optional` `Readonly` **signal**: `string`

___

### code

• `Optional` `Readonly` **code**: `number`

___

### stdout

• `Readonly` **stdout**: `string`

___

### stderr

• `Readonly` **stderr**: `string`

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
