---
title: Manage Repository Events
description: View a list of events happening within a repository and enable auto-deletion of events for repository maintenance.
keywords: registry, events, log, activity stream
---

> BETA DISCLAIMER
>
> This is beta content. It is not yet complete and should be considered a work in progress. This content is subject to change without notice.

## Overview 
Actions at their core are events which happen to a particular image in a repository. They happen when an image was pushed, pulled, promoted, scanned, or mirrored. As of v2.6, repositories now include an Activity tab which shows a paginated list of the most recent actions which have happened within a given repository. The range of activity types displayed will vary according to your repository privileges. If you are a DTR admin, you will have the option to compact the activity stream history as part of repository maintenance and cleanup.
  
In the following section, we will show you how to:

* View the list of activities in a repository
* Compact the activity stream history of a repository

## View List of Events

To view the list of events within a repository, navigate to `https://<dtr-url>`and log in with your UCP credentials. Select **Repositories** on the left navigation pane, and then click on the name of the repository that you want to view. Note that you will have to click on the repository name following
the `/` after the specific namespace for your repository.

![](../images/tag-pruning-0.png){: .with-border}

Select the **Pruning** tab, and click **New pruning policy** to specify your tag pruning criteria:

![](../images/tag-pruning-1.png){: .with-border}


DTR allows you to set your pruning triggers based on the following image attributes:

| Name            | Description                                        | Example           |
|:----------------|:---------------------------------------------------| :----------------|
| Tag name        | Whether the tag name equals, starts with, ends with, contains, is one of, or is not one of your specified string values | Tag name = `test`|
| Component name  | Whether the image has a given component and the component name equals, starts with, ends with, contains, is one of, or is not one of your specified string values | Component name starts with `b` |
| Vulnerabilities | Whether the image has vulnerabilities &ndash; critical, major, minor, or all &ndash; and your selected vulnerability filter is greater than or equals, greater than, equals, not equals, less than or equals, or less than your specified number | Critical vulnerabilities = `3` |
| License         | Whether the image uses an intellectual property license and is one of or not one of your specified words | License name = `docker` | 
| Last updated at | Whether the last image update was before your specified number of hours, days, weeks, or months. For details on valid time units, see [Go's ParseDuration function](https://golang.org/pkg/time/#ParseDuration). |  Last updated at: Hours = `12` |

Specify one or more image attributes to add to your pruning criteria, then choose:
 
	**Prune future tags** to save the policy and apply your selection to future tags. Only qualifying tags pushed after the policy addition will be pruned during garbage collection.

 	**Prune all tags** to save the policy and evaluate existing and future tags on your repository. 

Upon selection, you will see a confirmation message and will be redirected to your newly updated **Pruning** tab. 


![](../images/tag-pruning-2.png){: .with-border}


If you have specified multiple pruning policies on the repository, the **Pruning** tab will display a list of your prune triggers and details on when the last tag pruning was performed based on the trigger, a toggle for deactivating or reactivating the trigger, and a **View** link for modifying or deleting your selected trigger.

![](../images/tag-pruning-3.png){: .with-border}

All tag pruning policies on your account are evaluated every 15 minutes. Any qualifying tags are then deleted from the metadata store. If a tag pruning policy is modified or created, then the tag pruning policy for the *affected* repository will be evaluated.

## Set a tag limit

In addition to pruning policies, you can also set tag limits on repositories that you manage to restrict the number of tags on a given repository. Repository tag limits are processed in a first in first out (FIFO) manner. For example, if you set a tag limit of 2, adding a third tag would push out the first.

![](../images/tag-pruning-4.png){: .with-border}

To set a tag limit, select the repository that you want to update and click the **Settings** tab. Specify a number in the **Pruning** section and click **Save**. The **Pruning** tab will now display your tag limit above the prune triggers list along with a link to modify this setting.


![](../images/tag-pruning-5.png){: .with-border}

## Where to go next

- [Garbage collection](../../admin/configure/garbage-collection.md)
To get started, 

| Name            | Description                                        | Example           |
|:----------------|:---------------------------------------------------| :----------------|
| Tag name        | Whether the tag name equals, starts with, ends with, contains, is one of, or is not one of your specified string values | Tag name = `test`|
| Component name  | Whether the image has a given component and the component name equals, starts with, ends with, contains, is one of, or is not one of your specified string values | Component name starts with `b` |
| Vulnerabilities | Whether the image has vulnerabilities &ndash; critical, major, minor, or all &ndash; and your selected vulnerability filter is greater than or equals, greater than, equals, not equals, less than or equals, or less than your specified number | Critical vulnerabilities = `3` |
| License         | Whether the image uses an intellectual property license and is one of or not one of your specified words | License name = `docker` | 
| Last updated at | Whether the last image update was before your specified number of hours, days, weeks, or months. For details on valid time units, see [Go's ParseDuration function](https://golang.org/pkg/time/#ParseDuration). |  Last updated at: Hours = `12` |

![](../../images/view-repo-events-0.png){: .img-fluid .with-border}

## Where to go next

- [Delete images](delete-images.md)
