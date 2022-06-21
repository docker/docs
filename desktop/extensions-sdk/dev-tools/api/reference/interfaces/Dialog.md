---
title: Dialog
description: Docker extension API reference
keywords: Docker, extensions, sdk, API, reference
---

Allows opening native dialog boxes.

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
