+++
title = "regen-certs.md"
description = "Regenerate certificates for Docker Universal Control Plane."
keywords= ["install, ucp"]
[menu.main]
parent = "ucp_ref"
identifier = "ucp_ref_regen_certs"
+++

# docker/ucp regen-certs

The `regen-certs` command is no longer used.

## Usage

```bash
docker run --rm -it \
           --name ucp \
           -v /var/run/docker.sock:/var/run/docker.sock \
           docker/ucp \
           regen-certs [command options]
```

## Description

The `regen-certs` command is no longer used. Certificates are automatically renewed
before they expire and when SANs for a node are updated. To update SANs for a node,
edit the node in the Node Management Web UI.
