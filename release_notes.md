+++
title = "Release Notes"
description = "Docker Universal Control Plane"
[menu.ucp]
weight="-99"
+++


# UCP Release Notes

The latest release is 0.7.  Consult with your Docker sales engineer for the
release notes of earlier versions.

## Version 0.7

The following notes apply to this release:

### Public images on Docker Hub.

The UCP images now live within the `docker` organization on Docker Hub, and
are publicly accessible.

```bash
docker run --rm -it \
     --name ucp \
     -v /var/run/docker.sock:/var/run/docker.sock \
     docker/ucp \
     install --help
```

### Upgrade

Prior versions 0.5-0.6 of the beta can now be upgraded using the UCP tool.
Run the following command to review usage information.

```bash
docker run --rm -it \
    --name ucp \
    -v /var/run/docker.sock:/var/run/docker.sock \
    docker/ucp \
    upgrade --help
```

### UI

- New Volumes UI
- Paged views for nodes, containers, images, networks, volumes, and accounts to better handle large
  scale deployments.
- Support for browser based client cert login.
- Cluster Controllers listed on Nodes screen
- Ability to disable anonymous usage reporting


### Misc

- Fix restarting daemon or rebooting breaks cfssl and client bundle generation
- Controller port can be changed at install time (default 443)
- Fix TLS cipher support for direct admin access to engine proxy
- Fix event streaming and logs with follow
