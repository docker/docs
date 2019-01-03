---
title: Uninstall Docker Trusted Registry
description: Learn how to uninstall your Docker Trusted Registry installation.
keywords: dtr, install, uninstall
---

Uninstalling DTR can be done by simply removing all data associated with each
replica. To do that, you just run the remove command once per replica and destroy the last one :

```none
docker run -it --rm \
  docker/dtr:{{ page.dtr_version }} remove \
  --ucp-insecure-tls
```

```none
docker run -it --rm \
  docker/dtr:{{ page.dtr_version }} destroy \
  --ucp-insecure-tls
```

You will be prompted for the UCP URL, UCP credentials, and which replica to
remove/destroy.

To see what options are available in the remove/destroy commands, check the
[remove command reference documentation](../../../reference/cli/remove.md) and [destroy command reference documentation](../../../reference/cli/destroy.md).

## Where to go next

* [Scale your deployment](../configure/set-up-high-availability.md)
* [Install DTR](index.md)
