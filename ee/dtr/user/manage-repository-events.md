---
title: Manage Repository Events
description: View a list of image events happening within a repository and enable auto-deletion of these events for maintenance.
keywords: registry, events, log, activity stream
---

> BETA DISCLAIMER
>
> This is beta content. It is not yet complete and should be considered a work in progress. This content is subject to change without notice.

## Overview 

Actions at their core are events which happen to a particular image within a particular repository. To provide a quick summary of these events, DTR 2.6 now includes an **Activity** tab on each repository displaying a paginated list of the most recent events. The type of events listed will vary according to your [repository permission level](../admin/manage-users/permission-levels/). Additionally, DTR administrators can enable auto-deletion of repository events as part of maintenance and cleanup.
  
In the following section, we will show you how to:

* View the list of events in a repository, including <a href="#event-types">event types</a> associated with your permission level
* Enable auto-deletion of repository events based on your specified conditions

## View List of Events

To view the list of events within a repository, navigate to `https://<dtr-url>`and log in with your UCP credentials. Select **Repositories** on the left navigation pane, and then click on the name of the repository that you want to view. Note that you will have to click on the repository name following the `/` after the specific namespace for your repository.

![](../images/tag-pruning-0.png){: .with-border}

Select the **Activity** tab. You should see a list of events based on your repository permission level. Pull events are excluded by default and are only visible to repository and DTR administrators. Uncheck "Exclude pull" to view pull events.  

![](../images/manage-repo-events-0.png){: .with-border}


### Streamed Events

Note that the event types may reflect a different friendly name in the web interface and includes the relevant [CRUD operation](https://en.wikipedia.org/wiki/Create,_read,_update_and_delete).

| Event Type          | Description                                        | Permission Level        |
|:----------------|:---------------------------------------------------| :----------------|
| Push        |  | Authenticated Users |
|        |  | Authenticated Users |
| Scan        |  | Authenticated Users |
| Promotion        |  | Repository Admin |
| Delete        |  | Authenticated Users |
| Pull        |  | Repository Admin |
| Mirror        |  | Repository Admin |
| Create repo        |  | Authenticated Users |

## Enable auto-deletion of repository events

In addition to pruning policies, you can also set tag limits on repositories that you manage to restrict the number of tags on a given repository. Repository tag limits are processed in a first in first out (FIFO) manner. For example, if you set a tag limit of 2, adding a third tag would push out the first.

![](../images/tag-pruning-4.png){: .with-border}

To set a tag limit, select the repository that you want to update and click the **Settings** tab. Specify a number in the **Pruning** section and click **Save**. The **Pruning** tab will now display your tag limit above the prune triggers list along with a link to modify this setting.


![](../images/tag-pruning-5.png){: .with-border}

## Where to go next

- [Delete images](delete-images.md)
