---
title: Docker extension API reference
description: Docker extension API reference
keywords: Docker, extensions, sdk, API, reference
---

# Interface: Host

## Table of contents

### Methods

- [openExternal](Host.md#openexternal)

### Properties

- [platform](Host.md#platform)
- [arch](Host.md#arch)
- [hostname](Host.md#hostname)

## Methods

### openExternal

▸ **openExternal**(`url`): `void`

Opens an external URL with the system default browser.

```typescript
ddClient.host.openExternal("https://docker.com");
```

#### Parameters

| Name | Type | Description |
| :------ | :------ | :------ |
| `url` | `string` | The URL the browser will open (must have the protocol `http` or `https`). |

#### Returns

`void`

## Properties

### platform

• **platform**: `string`

Returns a string identifying the operating system platform. See https://nodejs.org/api/os.html#osplatform

___

### arch

• **arch**: `string`

Returns the operating system CPU architecture. See https://nodejs.org/api/os.html#osarch

___

### hostname

• **hostname**: `string`

Returns the host name of the operating system. See https://nodejs.org/api/os.html#oshostname
