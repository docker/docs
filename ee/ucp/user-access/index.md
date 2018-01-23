---
title: Web-based access
description: Learn how to access Docker Universal Control Plane from the web browser.
keywords: ucp, web, administration
ui_tabs:
- version: ucp-3.0
  orlower: true
next_steps:
- path: ../../authorization/
  title: Authorization
- path: cli-based-access/
  title: Access UCP from the CLI
redirect_from:
  - /ee/ucp/user/access-ucp/
---
{% if include.version=="ucp-3.0" %}

Docker Universal Control Plane allows you to manage your cluster in a visual
way, from your browser.

![](../../../../../images/ucp.png){: .with-border}


Docker UCP secures your cluster by using
[role-based access control](../../access-control/index.md).
From the browser, administrators can:

* Manage cluster configurations,
* Manage the permissions of users, teams, and organizations,
* See all images, networks, volumes, and containers.
* Grant permissions to users for scheduling tasks on specific nodes
  (with the Docker EE Advanced license).  

![](../../images/web-based-access-2.png){: .with-border}

Non-admin users can only see and change the images, networks, volumes, and
containers, and only when they're granted access by an administrator.

{% endif %}
