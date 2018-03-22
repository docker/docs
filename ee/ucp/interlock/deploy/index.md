---
title: Enable ayer 7 routing
description: Learn how to enable the layer 7 routing solution for UCP, that allows
  you to route traffic to swarm services.
keywords: routing, proxy
ui_tabs:
- version: ucp-3.0
  orhigher: false
- version: ucp-2.2
---

{% if include.version=="ucp-3.0" %}

To enable support for layer 7 routing, also known as HTTP routing mesh,
log in to the UCP web UI as an administrator, navigate to the **Admin Settings**
page, and click the **Routing Mesh** option. Check the **Enable routing mesh** option.

![http routing mesh](../../images/interlock-install-3.png){: .with-border}

By default, the routing mesh service listens on port 80 for HTTP and port
8443 for HTTPS. Change the ports if you already have services that are using
them.

Once you save, the layer 7 routing service can be used by your swarm services.

{% elsif include.version=="ucp-2.2" %}

* [Configure UCP 2.2 HTTP routing mesh](/datacenter/ucp/2.2/guides/admin/configure/use-domain-names-to-access-services.md)

{% endif %}
