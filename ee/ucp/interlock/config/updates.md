---
<<<<<<< HEAD
<<<<<<< HEAD
title: Update Interlock services
description: Learn how to update the UCP layer 7 routing solution services
keywords: routing, proxy, interlock
---

There are two parts to the update process:

1. Update the Interlock configuration to specify the new extension and/or proxy image versions.
2. Update the Interlock service to use the new configuration and image.

## Update the Interlock configuration
=======
=======
>>>>>>> df4abbfc665cd5b9e518a8f6d91bd686f1bf8ce5
title: Updating Interlock services
description: Learn how to update the UCP layer 7 routing solution services
keywords: routing, proxy
---

# Updating Interlock services
There are two parts to the update process:

1. Updating the Interlock configuration to specify the new extension and/or proxy image versions.
2. Updating the Interlock service to use the new configuration and image.

## Updating the Interlock configuration
<<<<<<< HEAD
>>>>>>> Raw content addition
=======
>>>>>>> df4abbfc665cd5b9e518a8f6d91bd686f1bf8ce5
Create the new configuration:

```bash
$> docker config create service.interlock.conf.v2 <path-to-new-config>
```

<<<<<<< HEAD
<<<<<<< HEAD
## Update the Interlock service
Remove the old configuration and specify the new configuration:

```bash
$> docker service update --config-rm service.interlock.conf ucp-interlock
$> docker service update --config-add source=service.interlock.conf.v2,target=/config.toml ucp-interlock
```

Next, update the Interlock service to use the new image. To pull the latest version of UCP, run the following:

```bash
$> docker pull docker/ucp:latest
```

### Example output

```bash
latest: Pulling from docker/ucp
cd784148e348: Already exists 
3871e7d70c20: Already exists 
cad04e4a4815: Pull complete 
Digest: sha256:63ca6d3a6c7e94aca60e604b98fccd1295bffd1f69f3d6210031b72fc2467444
Status: Downloaded newer image for docker/ucp:latest
docker.io/docker/ucp:latest
```

Next, list all the latest UCP images. To learn more about `docker/ucp images` and available options, 
see [the reference page](/reference/ucp/3.1/cli/images/).

```bash
$> docker run --rm docker/ucp images --list
```

### Example output

```bash
docker/ucp-agent:{{ page.ucp_version }}
docker/ucp-auth-store:{{ page.ucp_version }}
docker/ucp-auth:{{ page.ucp_version }}
docker/ucp-azure-ip-allocator:{{ page.ucp_version }}
docker/ucp-calico-cni:{{ page.ucp_version }}
docker/ucp-calico-kube-controllers:{{ page.ucp_version }}
docker/ucp-calico-node:{{ page.ucp_version }}
docker/ucp-cfssl:{{ page.ucp_version }}
docker/ucp-compose:{{ page.ucp_version }}
docker/ucp-controller:{{ page.ucp_version }}
docker/ucp-dsinfo:{{ page.ucp_version }}
docker/ucp-etcd:{{ page.ucp_version }}
docker/ucp-hyperkube:{{ page.ucp_version }}
docker/ucp-interlock-extension:{{ page.ucp_version }}
docker/ucp-interlock-proxy:{{ page.ucp_version }}
docker/ucp-interlock:{{ page.ucp_version }}
docker/ucp-kube-compose-api:{{ page.ucp_version }}
docker/ucp-kube-compose:{{ page.ucp_version }}
docker/ucp-kube-dns-dnsmasq-nanny:{{ page.ucp_version }}
docker/ucp-kube-dns-sidecar:{{ page.ucp_version }}
docker/ucp-kube-dns:{{ page.ucp_version }}
docker/ucp-metrics:{{ page.ucp_version }}
docker/ucp-pause:{{ page.ucp_version }}
docker/ucp-swarm:{{ page.ucp_version }}
docker/ucp:{{ page.ucp_version }}
```

Interlock starts and checks the config object, which has the new extension version, and 
=======
=======
>>>>>>> df4abbfc665cd5b9e518a8f6d91bd686f1bf8ce5
## Updating the Interlock service
Remove the old configuration and specify the new configuration:

```bash
$> docker service update --config-rm service.interlock.conf interlock
$> docker service update --config-add source=service.interlock.conf.v2,target=/config.toml interlock
```

Next, update the Interlock service to use the new image. The following example updates the Interlock core service to use the `sha256:d173014908eb09e9a70d8e5ed845469a61f7cbf4032c28fad0ed9af3fc04ef51`
version of Interlock. Interlock starts and checks the config object, which has the new extension version, and 
<<<<<<< HEAD
>>>>>>> Raw content addition
=======
>>>>>>> df4abbfc665cd5b9e518a8f6d91bd686f1bf8ce5
performs a rolling deploy to update all extensions.

```bash
$> docker service update \
<<<<<<< HEAD
<<<<<<< HEAD
    --image {{ page.ucp_org }}/ucp-interlock:{{ page.ucp_version }} \
    ucp-interlock
=======
    --image interlockpreview/interlock@sha256:d173014908eb09e9a70d8e5ed845469a61f7cbf4032c28fad0ed9af3fc04ef51 \
    interlock
>>>>>>> Raw content addition
=======
    --image interlockpreview/interlock@sha256:d173014908eb09e9a70d8e5ed845469a61f7cbf4032c28fad0ed9af3fc04ef51 \
    interlock
>>>>>>> df4abbfc665cd5b9e518a8f6d91bd686f1bf8ce5
```
