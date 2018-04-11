---
title: Uninstall Docker Trusted Registry
description: Learn how to uninstall your Docker Trusted Registry installation.
keywords: dtr, install, uninstall
---

Uninstalling DTR can be done by simply removing all data associated with each
replica. To do that, you just run the destroy command once per replica:

```bash
docker run -it --rm \
  docker/dtr:{{ page.dtr_version }} destroy \
  --ucp-insecure-tls
```

You will be prompted for the UCP URL, UCP credentials, and which replica to
destroy.

To see what options are available in the destroy command, check the
[destroy command reference documentation](/reference/dtr/2.5/cli/destroy.md).

## Where to go next

- [Scale your deployment](../configure/set-up-high-availability.md)
- [Install DTR](index.md)
