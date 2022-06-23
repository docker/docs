---
skip_read_time: true
description: Docker extension API reference
keywords: Docker, extensions, sdk, API, reference
---

# Interface: Toast

Toasts provide a brief notification to the user.
They appear temporarily and shouldn't interrupt the user experience.
They also don't require user input to disappear.

## Methods

### success

▸ **success**(`msg`): `void`

Display a toast message of type success.

```typescript
ddClient.desktopUI.toast.success("message");
```

#### Parameters

| Name | Type | Description |
| :------ | :------ | :------ |
| `msg` | `string` | The message to display in the toast. |

#### Returns

`void`

___

### warning

▸ **warning**(`msg`): `void`

Display a toast message of type warning.

```typescript
ddClient.desktopUI.toast.warning("message");
```

#### Parameters

| Name | Type | Description |
| :------ | :------ | :------ |
| `msg` | `string` | The message to display in the warning. |

#### Returns

`void`

___

### error

▸ **error**(`msg`): `void`

Display a toast message of type error.

```typescript
ddClient.desktopUI.toast.error("message");
```

#### Parameters

| Name | Type | Description |
| :------ | :------ | :------ |
| `msg` | `string` | The message to display in the toast. |

#### Returns

`void`
