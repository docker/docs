---
title: Manage repository webhooks via the web interface
description: Learn how to create, configure, and test repository webhooks for DTR using the web interface.
keywords: dtr, webhooks, ui, web interface, registry
---

## Prerequisites

- You must have admin privileges to the repository in order to create a webhook.
- See [Webhook types](/ee/dtr/admin/manage-webhooks/index.md#webhook-types) for a list of events you can trigger notifications for using the web interface.

## Create a webhook for your repository

1. In your browser, navigate to `https://<dtr-url>` and log in with your credentials.

2. Select **Repositories** from the left navigation pane, and then click on the name of the repository that you want to view. Note that you will have to click on the repository name following the `/` after the specific namespace for your repository.

3. Select the **Webhooks** tab, and click **New Webhook**.

     ![](/ee/dtr/images/manage-webhooks-1.png){: .with-border}

4. From the drop-down list, select the event that will trigger the webhook.
5. Set the URL which will receive the JSON payload. Click **Test** next to the **Webhook URL** field, so that you can validate that the integration is working. At your specified URL, you should receive a JSON payload for your chosen event type notification.

	```json
	{
	  "type": "TAG_PUSH",
	  "createdAt": "2019-05-15T19:39:40.607337713Z",
	  "contents": {
	    "namespace": "foo",
	    "repository": "bar",
	    "tag": "latest",
	    "digest": "sha256:b5bb9d8014a0f9b1d61e21e796d78dccdf1352f23cd32812f4850b878ae4944c",
	    "imageName": "foo/bar:latest",
	    "os": "linux",
	    "architecture": "amd64",
	    "author": "",
	    "pushedAt": "2015-01-02T15:04:05Z"
	  },
	  "location": "/repositories/foo/bar/tags/latest"
	}
	```

6. Expand "Show advanced settings" to paste the TLS certificate associated with your webhook URL. For testing purposes, you can test over HTTP instead of HTTPS.

7. Click **Create** to save. Once saved, your webhook is active and starts sending POST notifications whenever your chosen event type is triggered.

     ![](/ee/dtr/images/manage-webhooks-2.png){: .with-border}

As a repository admin, you can add or delete a webhook at any point. Additionally, you can create, view, and delete webhooks for your organization or trusted registry [using the API](use-the-api).

## Where to go next

- [Manage webhooks via the API](use-the-api)
