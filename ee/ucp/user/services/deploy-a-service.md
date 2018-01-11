---
title: Deploy a service
description: Learn how to deploy services to a cluster managed by Universal Control Plane.
keywords: ucp, deploy, service
ui_tabs:
- version: ucp-3.0
  orlower: true
cli_tabs:
- version: docker-cli-linux
---
{% if include.ui %}

{% if include.version=="ucp-3.0" %}

You can deploy and monitor your services from the UCP web UI. In this example
we'll deploy an [NGINX](https://www.nginx.com/) web server and make it
accessible on port `8000`.

In your browser, navigate to the UCP web UI and click **Services**. On the
**Create a Service** page, click **Create Service** to configure the
NGINX service.

Fill in the following fields:

| Field         | Value |
|:--------------|:------|
| Service name  | nginx |
| Image name    | nginx:latest |

![](../../images/deploy-a-service-1.png){: .with-border}

In the left pane, click **Network**. In the **Ports** section,
click **Publish Port** and fill in the following fields:

| Field         | Value |
|:--------------|:------|
| Internal port | 80    |
| Protocol      | tcp   |
| Publish mode  | Ingress |
| Public port   | 8000  | 

![](../../images/deploy-a-service-2.png){: .with-border}

Click **Confirm** to map the ports for the NGINX service. 

Once you've specified the service image and ports, click **Create** to
deploy the service into the UCP cluster.

Once the service is up and running, you'll be able to see the default NGINX
page, by going to `http://<node-ip>:8000`.

![](../../images/deploy-a-service-4.png){: .with-border}

{% endif %}
{% endif %}

{% if include.cli %}

{% if include.version=="docker-cli-linux" %}

You can also deploy the same service from the CLI. Once you've set up your
[UCP client bundle](../access-ucp/cli-based-access.md), run:

```bash
docker service create --name nginx \
  --publish 8000:80 \
  --label com.docker.ucp.access.owner=<your-username> \
  nginx
```

{% endif %}
{% endif %}