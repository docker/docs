---
title: Manage webhooks
description: Learn how to create, configure, and test webhooks in Docker Trusted Registry.
keywords: registry, webhooks
---

DTR has webhooks so that you can run custom logic when an event happens. This
lets you build complex CI and CD pipelines with your Docker images.

## Create a webhook

To create a webhook, navigate to the **repository details** page, choose
the **Webhooks** tab, and click **New Webhook**.

![](../images/manage-webhooks-1.png){: .with-border}

Select the event that will trigger the webhook, and set the URL to send
information about the event. Once everything is set up, click **Test** for
DTR to send a JSON payload to the URL you set up, so that you can validate
that the integration is working. You get an event that looks like this:

```
{
  "contents": {
    "architecture": "amd64",
    "author": "",
    "digest": "sha256:b5bb9d8014a0f9b1d61e21e796d78dccdf1352f23cd32812f4850b878ae4944c",
    "imageName": "example.com/foo/bar:latest",
    "namespace": "foo",
    "os": "linux",
    "pushedAt": "2015-01-02T15:04:05Z",
    "repository": "bar",
    "tag": "latest"
  },
  "createdAt": "2017-06-20T01:29:53.046620425Z",
  "location": "/repositories/foo/bar/tags/latest",
  "type": "TAG_PUSH"
}
```

Once you save, your webhook is active and starts sending notifications when
the event is triggered.

![](../images/manage-webhooks-2.png){: .with-border}

## Where to go next

* [Create deployment policies](create-deployment-policies.md)
