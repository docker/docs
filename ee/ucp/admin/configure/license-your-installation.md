---
title: License your installation
description: Learn how to license your Docker Universal Control Plane installation.
keywords: Universal Control Plane, UCP, install, license
ui_tabs:
- version: ucp-3.0
  orhigher: false
- version: ucp-2.2
  orlower: true
next_steps:
- path: ../install/
  title: Install UCP
- path: ../install/install-offline/
  title: Install UCP offline
---
{% if include.version=="ucp-3.0" %}

After installing Docker Universal Control Plane, you need to license your
installation. Here's how to do it.

## Download your license

Go to [Docker Store](https://www.docker.com/enterprise-edition) and
download your UCP license, or get a free trial license.

![](../../images/license-ucp-1.png){: .with-border}

## License your installation

Once you've downloaded the license file, you can apply it to your UCP
installation. 

In the UCP web UI, log in with administrator credentials and
navigate to the **Admin Settings** page.

In the left pane, click **License** and click **Upload License**. The
license refreshes immediately, and you don't need to click **Save**.

![](../../images/license-ucp-2.png){: .with-border}

{% elsif include.version=="ucp-2.2" %}

Learn about [licensing your installation](/datacenter/ucp/2.2/guides/admin/configure/license-your-installation.md).

{% endif %}
