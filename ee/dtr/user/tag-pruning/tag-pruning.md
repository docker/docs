---
title: Tag Pruning
description: Skip the management headache of deciding which tags to delete or preserve by configuring a tag pruning policy or enforcing a tag limit per repository in the Docker Trusted Registry
keywords: registry, tag pruning, tag limit, repo management
---

> BETA DISCLAIMER
>
> This is beta content. It is not yet complete and should be considered a work in progress. This content is subject to change without notice.

## Overview

You can configure the Docker Trusted Registry (DTR) to automatically clean up repository 
images which are abandoned and no longer necessary. You enable tag pruning on each 
repository that you manage by:

* [specifying a tag pruning policy] (#specify-a-tag-pruning-policy)
* [setting a tag limit] (#set-a-tag-limit)

## Specify a tag pruning policy

You can now add tag pruning policies on each repository that you manage. To get started, navigate to `https://<dtr-url>` and log in with your credentials.
 
Select **Repositories** on the left navigation pane, and then click on the name of the repository
that you want to update. Note that you will have to click on the repository name following
the `/` after the specific namespace for your repository.

![tag pruning](../../images/tag-pruning-0.png){: .with-border}

Select the **Pruning** tab, and click **New pruning policy** to specify your tag pruning criteria:

![](../../images/tag-pruning-1.png){: .with-border}


DTR allows you to set your pruning triggers based on the following image attributes:

| Name            | Description                                        |
|:----------------|:---------------------------------------------------|
| Tag name        | Whether the tag name equals, starts with, ends with, contains, is one of, or is not one of your specified string values |
| Component name  | Whether the image has a given component and the component name equals, starts with, ends with, contains, is one of, or is not one of your specified string values |
| Vulnerabilities | Whether the image has vulnerabilities &ndash; critical, major, minor, or all &ndash; and your selected vulnerability filter is greater than or equals, greater than, equals, not equals, less than or equals, or less than your specified number |
| License         | Whether the image uses an intellectual property license and is one of or not one of your specified words 
| Last updated at | Whether the last image update was before your specified number of hours, days, weeks, or months |  

Specify one or more image attributes to your pruning criteria, and choose **Prune future tags** or **Prune all tags**. Upon selection, you will see a confirmation message and will be redirected to your newly updated **Pruning** tab. 

![](../../images/tag-pruning-2.png){: .with-border}


If you have specified multiple pruning policies on the repository, the **Pruning** tab will display a list of your prune triggers and details on when a tag pruning operation was performed based on the trigger, a toggle for deactivating or reactivating the trigger, and a **View** link for modifying or deleting your selected trigger.

![](../../images/tag-pruning-3.png){: .with-border}

All tag pruning policies on your account are evaluated every 15 minutes. Any qualifying tags are then deleted from the metadata store. If a tag pruning policy is modified or created, then the tag pruning policy for the *affected* repository will be evaluated.

This does not cover modification or deletion of pruning policies.

> Tag Pruning
>
> Tag pruning operation only deletes a tag and does not carry out any actual blob deletion. For actual blob deletions, see [Garbage Collection](garbage-collection/index.md).


## Set a tag limit

In addition to pruning policies, you can also set repository tag limits to restrict the number of tags on a given repository. Repository tag limits are processed in a first in first out (FIFO) manner. For example, if you set a tag limit of 2, adding a third tag would push out the first.

![](../../images/tag-pruning-4.png){: .with-border}

To set a tag limit, select the repository that you want to update and click the **Settings** tab. Specify a number in the **Pruning** section and click **Save**. The **Pruning** tab will now display your tag limit above the prune triggers list along with a link to modify this setting.


![](../../images/tag-pruning-5.png){: .with-border}




## Where to go next

- [Garbage collection](garbage-collection/index.md)
