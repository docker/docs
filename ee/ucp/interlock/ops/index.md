---
title: Update Interlock
description: Learn about Interlock, an application routing and load balancing system
  for Docker Swarm.
keywords: ucp, interlock, load balancing
ui_tabs:
- version: ucp-3.0
  orhigher: false
---

{% if include.version=="ucp-3.0" %}

The following describes how to update Interlock.  There are two parts
to the upgrade.  First, the Interlock configuration must be updated
to specify the new extension and/or proxy image versions.  Then the Interlock
service will be updated.

First we will create the new configuration:

```bash
$> docker config create service.interlock.conf.v2 <path-to-new-config>
```

Then you can update the Interlock service to remove the old and use the new:

```bash
$> docker service update --config-rm service.interlock.conf interlock
$> docker service update --config-add source=service.interlock.conf.v2,target=/config.toml interlock
```

Next update the Interlock service to use the new image:

```bash
$> docker service update \
    --image interlockpreview/interlock@sha256:d173014908eb09e9a70d8e5ed845469a61f7cbf4032c28fad0ed9af3fc04ef51 \
    interlock
```

This will update the Interlock core service to use the `sha256:d173014908eb09e9a70d8e5ed845469a61f7cbf4032c28fad0ed9af3fc04ef51`
version of Interlock.  Interlock will start and check the config object which has the new extension version and will
perform a rolling deploy to update all extensions.

{% endif %}
