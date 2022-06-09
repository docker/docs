---
title: Docker extension API reference
description: Docker extension API reference
keywords: Docker, extensions, sdk, API, reference
---

# Interface: Dialog

Allows opening native dialog boxes.

## Table of contents

### Methods

- [showOpenDialog](Dialog.md#showopendialog)

## Methods

### showOpenDialog

▸ **showOpenDialog**(`dialogProperties`): `Promise`<[`OpenDialogResult`](OpenDialogResult.md)\>

Display a native open dialog, allowing to select a file or a folder.

```typescript
ddClient.desktopUI.dialog.showOpenDialog({properties: ['openFile']});
```

#### Parameters

| Name | Type | Description |
| :------ | :------ | :------ |
| `dialogProperties` | `any` | Properties to specify the open dialog behaviour, see https://www.electronjs.org/docs/latest/api/dialog#dialogshowopendialogbrowserwindow-options. |

#### Returns

`Promise`<[`OpenDialogResult`](OpenDialogResult.md)\>
