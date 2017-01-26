---
description: Learn how to uninstall your Docker Trusted Registry installation.
keywords: docker, dtr, install, uninstall
title: Uninstall Docker Trusted Registry
---

Uninstalling DTR is a two-step process. You first scale your DTR deployment down
to a single replica. Then you uninstall the last DTR replica, which permanently
removes DTR and deletes all its data.

Start by [scaling down your DTR deployment](scale-your-deployment.md) to a
single replica.

When your DTR deployment is down to a single replica, you can use the
`docker/dtr destroy` command to permanently remove DTR and all its data:

1. Use ssh to log into any node that is part of UCP.
2. Uninstall DTR:

```none
docker run -it --rm \
  {{ page.docker_image }} destroy \
  --ucp-insecure-tls
```

To see what options are available in the destroy command, check the
[destroy command reference documentation](../../reference/cli/destroy.md).

## Where to go next

* [Scale your deployment](scale-your-deployment.md)
* [Install DTR](index.md)
