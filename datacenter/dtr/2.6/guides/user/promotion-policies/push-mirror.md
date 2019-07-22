---
title: Mirror images to another registry
description: Learn how to create a promotion policy that promotes images to an external registry, creating a DTR mirror.
keywords: registry, promotion, mirror
---

Docker Trusted Registry allows you to create mirroring policies for a repository.
When an image gets pushed to a repository and meets the mirroring criteria,
DTR automatically pushes it to a repository in a remote Docker Trusted or Hub registry.

This not only allows you to mirror images but also allows you to create
image promotion pipelines that span multiple DTR deployments and datacenters.

In this example we will create an image mirroring policy such that:

1. Developers iterate and push their builds to `dtr-example.com/dev/website` &endash; the
repository in the DTR deployment dedicated to development.
2. When the team creates a stable build, they make sure their image is tagged
with `-stable`.
3. When a stable build is pushed to `dtr-example.com/dev/website`, it will
automatically be pushed to `qa-example.com/qa/website`, mirroring the image and
promoting it to the next stage of development.

With this mirroring policy, the development team does not need access to the
QA cluster, and the QA team does not need access to the development
cluster.

You need to have permissions to push to the destination repository in order to set up the mirroring policy.

## Configure your repository

Once you have [created the repository](../manage-images/index.md), navigate to
the repository page on the web interface, and select the
**Mirrors** tab.

![create integration](../../images/push-mirror-2.png){: .with-border}

Click **New mirror**, and define where the image will be pushed if
it meets the mirroring criteria. Make sure the account you use for the integration
has permissions to write to the remote repository. Under **Mirror direction**, choose **Push to remote registry**. 

In this example, the image gets pushed to the `qa/website` repository of a
DTR deployment available at `qa-example.com` using a service account
that was created just for mirroring images between repositories. Note that you may use a password or access token to log in to your remote registry.

If the destination DTR deployment is using self-signed TLS certificates or
certificates issued by your own certificate authority, click
**Show advanced settings** to provide the CA certificate used by the
DTR where the image will be pushed.

You can get that CA certificate by accessing `https://<destination-dtr>/ca`.

Once you're done, click **Connect** to test the integration.

![test connection](../../images/push-mirror-3.png){: .with-border}

DTR allows you to set your mirroring policy based on the following image attributes:

| Name            | Description                                        | Example           |
|:----------------|:---------------------------------------------------| :----------------|
| Tag name        | Whether the tag name equals, starts with, ends with, contains, is one of, or is not one of your specified string values | Copy image to remote repository if Tag name ends in `stable`|
| Component name  | Whether the image has a given component and the component name equals, starts with, ends with, contains, is one of, or is not one of your specified string values | Copy image to remote repository if Component name starts with `b` |
| Vulnerabilities | Whether the image has vulnerabilities &ndash; critical, major, minor, or all &ndash; and your selected vulnerability filter is greater than or equals, greater than, equals, not equals, less than or equals, or less than your specified number | Copy image to remote repository if Critical vulnerabilities = `3` |
| License         | Whether the image uses an intellectual property license and is one of or not one of your specified words | Copy image to remote repository if License name = `docker` | 

Finally you can choose to keep the image tag, or transform the tag into
something more meaningful in the remote registry by using a tag template.

![choose policy](../../images/push-mirror-4.png){: .with-border}

In this example, if an image in the `dev/website` repository is tagged with
a word that ends in "stable", DTR will automatically push that image to
the DTR deployment available at `qa-example.com`. The image is pushed to the
`qa/website` repository and is tagged with the timestamp of when the image
was promoted.

Everything is set up! Once the development team pushes an image that complies
with the policy, it automatically gets promoted to `qa/website` in the remote trusted registry at `qa-example.com`.

## Metadata persistence

When an image is pushed to another registry using a mirroring policy, scanning
and signing data is not persisted in the destination repository.

If you have scanning enabled for the destination repository, DTR is going to scan
the image pushed. If you want the image to be signed, you need to do it manually.

## Where to go next

* [Mirror images from another registry](pull-mirror.md)
