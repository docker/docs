---
title: Garbage collection
description: Save disk space by configuring the garbage collection settings in
  Docker Trusted Registry
keywords: registry, online garbage collection, gc, space, disk space
---

> BETA DISCLAIMER
>
> This is beta content. It is not yet complete and should be considered a work in progress. This content is subject to change without notice.

You can configure the Docker Trusted Registry (DTR) to automatically delete unused image
layers, thus saving you disk space. This process is also known as garbage collection.

## How DTR deletes unused layers

First you configure DTR to run a garbage collection job on a fixed schedule. At
the scheduled time, DTR:

2. Identifies and marks unused image layers.
3. Deletes the marked image layers.

Starting in DTR 2.5, we introduced an experimental feature which lets you run garbage collection jobs
without putting DTR in read-only mode. As of v2.6.0, online garbage collection is no longer in 
experimental mode. This means that the registry no longer has to be in read-only mode (or offline) 
during garbage collection. 


## Schedule garbage collection

In your browser, navigate to `https://<dtr-url>` and log in with your credentials. Select **System** on the left navigation pane, and then click
the **Garbage collection** tab to schedule garbage collection.

![](../../images/garbage-collection-0.png){: .with-border}

Select for how long the garbage collection job should run:
* Until done: Run the job until all unused image layers are deleted.
* For x minutes: Only run the garbage collection job for a maximum of x minutes
at a time.
* Never: Never delete unused image layers.

If you select *Until done* or *For x minutes*, you can specify a recurring schedule in UTC (Coordinated Universal Time) with the following options:
* Custom cron schedule - (Hour, Day of Month, Month, Weekday)
* Daily at midnight UTC
* Every Saturday at 1am UTC
* Every Sunday at 1am UTC
* Do not repeat

![](../../images/garbage-collection-1.png){: .with-border}

Once everything is configured you can choose to **Save & Start** to
run the garbage collection job immediately, or just **Save** to run the job on the next
scheduled interval.

## Review the garbage collection job log

In v2.5, you were notified with a banner under main navigation that no one can push images while a garbage collection job is running. Notice how this is no longer the case
with v2.6.0. If you clicked **Save & Start** previously, verify that the garbage collection routine started by navigating to *Jobs Logs*.

![](../../images/garbage-collection-2.png){: .with-border}

## Under the hood

Each image stored in DTR is made up of multiple files:

* A list of image layers that represent the image filesystem.
* A configuration file that contains the architecture of the image and other
metadata.
* A manifest file containing the list of all layers and configuration file for
an image.

All these files are stored in a content-addressable way in which the name of
the file is the result of hashing the file's content. This means that if two
image tags have exactly the same content, DTR only stores the image content
once, even if the tag name is different.

As an example, if `wordpress:4.8` and `wordpress:latest` have the same content,
they will only be stored once. If you delete one of these tags, the other won't
be deleted.

This means that when users delete an image tag, DTR can't delete the underlying
files of that image tag since it's possible that there are other tags that
also use the same files.

To delete unused image layers, DTR:
1. Becomes read-only to make sure that no one is able to push an image, thus
changing the underlying files in the filesystem.
2. Check all the manifest files and keep a record of the files that are
referenced.
3. If a file is never referenced, that means that no image tag uses it, so it
can be safely deleted.

## Where to go next

- [Deploy DTR caches](deploy-caches/index.md)
