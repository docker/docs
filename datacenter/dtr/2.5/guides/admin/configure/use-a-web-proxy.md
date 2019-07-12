---
title: Use a web proxy
description: Learn how to configure Docker Content Trust to use a web proxy to
  reach external services.
keywords: dtr, configure, http, proxy
---

Docker Trusted Registry makes outgoing connections to check for new versions,
automatically renew its license, and update its vulnerability database.
If DTR can't access the internet, then you need to manually apply updates.

One option to keep your environment secure while still allowing DTR access to
the internet is to use a web proxy. If you have an HTTP or HTTPS proxy, you
can configure DTR to use it. To avoid downtime you should do this configuration
outside business peak hours.

As an administrator, log into a node where DTR is deployed, and run:

```
docker run -it --rm \
  {{ page.dtr_org }}/{{ page.dtr_repo }}:{{ page.dtr_version }} reconfigure \
  --http-proxy http://<domain>:<port> \
  --https-proxy https://<doman>:<port> \
  --ucp-insecure-tls
```

To confirm how DTR is configured, check the **Settings** page on the web UI.

![DTR settings](../../images/use-a-web-proxy-1.png){: .with-border}

## Where to go next

* [Configure garbage collection](garbage-collection.md)
