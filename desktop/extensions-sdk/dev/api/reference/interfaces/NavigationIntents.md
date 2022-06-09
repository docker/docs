---
title: Docker extension API reference
description: Docker extension API reference
keywords: Docker, extensions, sdk, API, reference
---

# Interface: NavigationIntents

## Table of contents

### Container Methods

- [viewContainers](NavigationIntents.md#viewcontainers)
- [viewContainer](NavigationIntents.md#viewcontainer)
- [viewContainerLogs](NavigationIntents.md#viewcontainerlogs)
- [viewContainerInspect](NavigationIntents.md#viewcontainerinspect)
- [viewContainerStats](NavigationIntents.md#viewcontainerstats)

### Images Methods

- [viewImages](NavigationIntents.md#viewimages)
- [viewImage](NavigationIntents.md#viewimage)

### Other Methods

- [viewDevEnvironments](NavigationIntents.md#viewdevenvironments)

### Volume Methods

- [viewVolumes](NavigationIntents.md#viewvolumes)
- [viewVolume](NavigationIntents.md#viewvolume)

## Container Methods

### viewContainers

▸ **viewContainers**(): `Promise`<`void`\>

Navigate to the containers window in Docker Desktop.

```typescript
ddClient.desktopUI.navigate.viewContainers()
```

#### Returns

`Promise`<`void`\>

___

### viewContainer

▸ **viewContainer**(`id`): `Promise`<`void`\>

Navigate to the container window in Docker Desktop.

```typescript
await ddClient.desktopUI.navigate.viewContainer(id)
```

#### Parameters

| Name | Type | Description |
| :------ | :------ | :------ |
| `id` | `string` | The full container id, e.g. `46b57e400d801762e9e115734bf902a2450d89669d85881058a46136520aca28`. You can use the `--no-trunc` flag as part of the `docker ps` command to display the full container id. |

#### Returns

`Promise`<`void`\>

A promise that fails if the container doesn't exist.

___

### viewContainerLogs

▸ **viewContainerLogs**(`id`): `Promise`<`void`\>

Navigate to the container logs window in Docker Desktop.

```typescript
await ddClient.desktopUI.navigate.viewContainerLogs(id)
```

#### Parameters

| Name | Type | Description |
| :------ | :------ | :------ |
| `id` | `string` | The full container id, e.g. `46b57e400d801762e9e115734bf902a2450d89669d85881058a46136520aca28`. You can use the `--no-trunc` flag as part of the `docker ps` command to display the full container id. |

#### Returns

`Promise`<`void`\>

A promise that fails if the container doesn't exist.

___

### viewContainerInspect

▸ **viewContainerInspect**(`id`): `Promise`<`void`\>

Navigate to the container inspect window in Docker Desktop.

```typescript
await ddClient.desktopUI.navigate.viewContainerInspect(id)
```

#### Parameters

| Name | Type | Description |
| :------ | :------ | :------ |
| `id` | `string` | The full container id, e.g. `46b57e400d801762e9e115734bf902a2450d89669d85881058a46136520aca28`. You can use the `--no-trunc` flag as part of the `docker ps` command to display the full container id. |

#### Returns

`Promise`<`void`\>

A promise that fails if the container doesn't exist.

___

### viewContainerStats

▸ **viewContainerStats**(`id`): `Promise`<`void`\>

Navigate to the container stats to see the CPU, memory, disk read/write and network I/O usage.

```typescript
await ddClient.desktopUI.navigate.viewContainerStats(id)
```

#### Parameters

| Name | Type | Description |
| :------ | :------ | :------ |
| `id` | `string` | The full container id, e.g. `46b57e400d801762e9e115734bf902a2450d89669d85881058a46136520aca28`. You can use the `--no-trunc` flag as part of the `docker ps` command to display the full container id. |

#### Returns

`Promise`<`void`\>

A promise that fails if the container doesn't exist.

___

## Images Methods

### viewImages

▸ **viewImages**(): `Promise`<`void`\>

Navigate to the images window in Docker Desktop.

```typescript
await ddClient.desktopUI.navigate.viewImages()
```

#### Returns

`Promise`<`void`\>

___

### viewImage

▸ **viewImage**(`id`, `tag`): `Promise`<`void`\>

Navigate to a specific image referenced by `id` and `tag` in Docker Desktop.
In this navigation route you can find the image layers, commands, created time and size.

```typescript
await ddClient.desktopUI.navigate.viewImage(id, tag)
```

#### Parameters

| Name | Type | Description |
| :------ | :------ | :------ |
| `id` | `string` | The full image id (including sha), e.g. `sha256:34ab3ae068572f4e85c448b4035e6be5e19cc41f69606535cd4d768a63432673`. |
| `tag` | `string` | The tag of the image, e.g. `latest`, `0.0.1`, etc. |

#### Returns

`Promise`<`void`\>

A promise that fails if the image doesn't exist.

___

## Other Methods

### viewDevEnvironments

▸ **viewDevEnvironments**(): `Promise`<`void`\>

Navigate to the Dev Environments window in Docker Desktop.

```typescript
ddClient.desktopUI.navigate.viewDevEnvironments()
```

#### Returns

`Promise`<`void`\>

___

## Volume Methods

### viewVolumes

▸ **viewVolumes**(): `Promise`<`void`\>

Navigate to the volumes window in Docker Desktop.

```typescript
ddClient.desktopUI.navigate.viewVolumes()
```

#### Returns

`Promise`<`void`\>

___

### viewVolume

▸ **viewVolume**(`volume`): `Promise`<`void`\>

Navigate to a specific volume in Docker Desktop.

```typescript
await ddClient.desktopUI.navigate.viewVolume(volume)
```

#### Parameters

| Name | Type | Description |
| :------ | :------ | :------ |
| `volume` | `string` | The name of the volume, e.g. `my-volume`. |

#### Returns

`Promise`<`void`\>
