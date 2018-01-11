---
title: Uninstall Docker Trusted Registry
description: Learn how to uninstall your Docker Trusted Registry installation.
keywords: dtr, install, uninstall
ui_tabs:
- version: dtr-2.5
  orhigher: false
- version: dtr-2.4
  orlower: true
next_steps:
- path: ../configure/set-up-high-availability/
  title: Scale your deployment
- path: ../install/
  title: Install DTR
---
{% if include.version=="dtr-2.5" %}

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

{% elsif include.version=="dtr-2.4" %}

Learn about [uninstalling Docker Trusted Registry](/datacenter/dtr/2.4/guides/admin/install/uninstall.md).

{% endif %}
