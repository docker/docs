---
title: Navigation
description: Docker extension API
keywords: Docker, extensions, sdk, API
---

`ddClient.desktopUI.navigate` enables navigation to specific screens of Docker Desktop such as the containers tab, the images tab, or a specific container's logs.

For example, navigate to a given container logs:

```typescript
await ddClient.desktopUI.navigate.viewContainerLogs(id);
```

#### Parameters

| Name | Type     | Description                                                                                                                                                                                            |
| :--- | :------- | :----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `id` | `string` | The full container id, e.g. `46b57e400d801762e9e115734bf902a2450d89669d85881058a46136520aca28`. You can use the `--no-trunc` flag as part of the `docker ps` command to display the full container id. |

#### Returns

`Promise`<`void`\>

A promise that fails if the container doesn't exist.

For more details about all navigation methods, see the [Navigation API reference](reference/interfaces/NavigationIntents.md).

> Deprecated navigation methods
>
> These methdos are deprecated and will be removed in a future version. Use the methods specified above.

```typescript
window.ddClient.navigateToContainers();
// id - the full container id, e.g. `46b57e400d801762e9e115734bf902a2450d89669d85881058a46136520aca28`
window.ddClient.navigateToContainer(id);
window.ddClient.navigateToContainerLogs(id);
window.ddClient.navigateToContainerInspect(id);
window.ddClient.navigateToContainerStats(id);

window.ddClient.navigateToImages();
window.ddClient.navigateToImage(id, tag);

window.ddClient.navigateToVolumes();
window.ddClient.navigateToVolume(volume);

window.ddClient.navigateToDevEnvironments();
```
