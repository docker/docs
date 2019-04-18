---
title: Update Interlock services
description: Learn how to update the UCP layer 7 routing solution services
keywords: routing, proxy, interlock
---

There are two parts to the update process:

1. Update the Interlock configuration to specify the new extension and/or proxy image versions.
2. Update the Interlock service to use the new configuration and image.

## Update the Interlock configuration
Create the new configuration:

```bash
$> docker config create service.interlock.conf.v2 <path-to-new-config>
```

## Update the Interlock service
Remove the old configuration and specify the new configuration:

```bash
$> docker service update --config-rm service.interlock.conf interlock
$> docker service update --config-add source=service.interlock.conf.v2,target=/config.toml interlock
```

Next, update the Interlock service to use the new image. The following example updates the Interlock core service to use the `sha256:d173014908eb09e9a70d8e5ed845469a61f7cbf4032c28fad0ed9af3fc04ef51`
version of Interlock. Interlock starts and checks the config object, which has the new extension version, and 
performs a rolling deploy to update all extensions.

```bash
$> docker service update \
    --image interlockpreview/interlock@sha256:d173014908eb09e9a70d8e5ed845469a61f7cbf4032c28fad0ed9af3fc04ef51 \
    interlock
```
