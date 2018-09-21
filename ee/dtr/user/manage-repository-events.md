---
title: Manage Repository Events
description: View a list of image events happening within a repository.
keywords: registry, events, log, activity stream
---

> BETA DISCLAIMER
>
> This is beta content. It is not yet complete and should be considered a work in progress. This content is subject to change without notice.

## Overview 

Starting in DTR 2.6, each repository page includes an **Activity** tab which displays a sortable and paginated list of the most recent events within the repository. This offers better visibility along with the ability to audit events. Event types listed will vary according to your [repository permission level](../admin/manage-users/permission-levels/). Additionally, DTR admins can [enable auto-deletion of repository events](../admin/configure/auto-delete-repo-events/) as part of maintenance and cleanup.
  
In the following section, we will show you how to:

* Review actions or events that you have access to
* View the list of events in a repository, including <a href="#event-types">event types</a> associated with your permission level

## Review Repository Permissions

To review the type of events you have access to, do the following:
1. Navigate to `https://<dtr-url>`and log in with your UCP credentials. 

2. Select **Repositories** on the left navigation pane, and then click on the name of the repository that you want to view. Note that you will have to click on the repository name following the `/` after the specific namespace for your repository.

![](../images/tag-pruning-0.png){: .img-fluid .with-border}

3. You should see the **Info** tab by default. Notice **Your Permission** under **Docker Pull Command**. 

4. Hover over the question mark next to [your permission level](../admin/manage-users/permission-levels/) to view the list of repository events you have access to.
	   
![](../images/manage-repo-events-2.png){: .img-fluid .with-border}

> KNOWN LIMITATIONS
>
> Your repository permissions list may include events that are not displayed in the **Activity** tab. It is also not an exhaustive list of event types on your activity stream.

## View List of Events

To view the list of events within a repository, select the **Activity** tab. You should see a paginated list of the latest events based on your repository permission level. By default, **Activity** shows the latest `10` events and excludes pull events, which are only visible to repository and DTR admins. 
   * If you're a repository or a DTR admin, uncheck "Exclude pull" to view pull events. This should give you a better understanding of who is consuming your images.
   * To update your event view, select a different time filter from the drop-down list.  

![](../images/manage-repo-events-0.png){: .img-fluid .with-border}

### Activity Stream
 
Let's break down the data included in an event. We will use the highlighted "Create Promotion Policy" event as our example.

| Event Detail          | Description                                        | Example |
|:----------------|:-------------------------------------------------|:--------|
| Label        |  Friendly name of the event. | `Create Promotion Policy`
| Repository  | This will always be the repository in review following the `<user-or-org>/<repository_name>` convention outlined in [Create a Repository](../user/manage-images/#create-a-repository). | `test-org/test-repo-1` |
| Tag        | Tag affected by the event, when applicable. | `test-org/test-repo-1:latest` where `latest` is the affected tag| 
| SHA | The [digest value](/registry/spec/api/#content-digests) for `CREATE` operations such as creating a new image tag or a promotion policy. | `sha256:bbf09ba3` |
| Type | Event type. Possible values are: `CREATE`, `GET`, `UPDATE`, `DELETE`, `SEND`, `FAIL` and `SCAN` | `CREATE` |
| Initiated by | The actor responsible for the event. For user-initiated events, this will reflect the user ID and link to that user's profile. For image events triggered by a policy &ndash; pruning, pull / push mirroring, or promotion &ndash; this will reflect the relevant policy ID except for manual promotions where it reflects `PROMOTION MANUAL_P`, and link to the relevant policy page. Other event actors may not include a link.  | `PROMOTION CA5E7822` |
| Date and Time | When the event happened in your configured time zone. | `9/13/2018 9:59 PM` |  

### Event Audits

Given the level of detail on each event, it should be easy for DTR and security admins to determine what events have taken place inside of DTR.  For example, when an image which shouldn’t have been deleted ends up getting deleted, the security admin can determine when and who initiated the deletion.

### Event Permissions

For more details on different permission levels within DTR, see [Authentication and authorization in DTR](../admin/manage-users/) to understand the minimum level required to view the different repository events.  

| Repository Event          | Description                                        | Minimum Permission Level        |
|:----------------|:---------------------------------------------------| :----------------|
| Push        |  Refers to "Create Manifest" and "Update Tag" events. Learn more about [pushing images](../user/manage-images/pull-and-push-images/#push-the-image). | Authenticated Users |
| Scan        | Requires [security scanning to be set up](../admin/configure/set-up-vulnerability-scans/) by a DTR admin. Once enabled, this will display as a `SCAN` event type.  | Authenticated Users |
| Promotion        |  Refers to a "Create Promotion Policy" event which links to the **Promotions** tab of the repository where you can edit the existing promotions. See [Promotion Policies](../user/promotion-policies/) for different ways to promote an image. | Repository Admin |
| Delete        |  Refers to "Delete Tag" events. Learn more about [deleting images](../user/manage-images/pull-and-push-images/#delete-images). | Authenticated Users |
| Pull        | Refers to "Get Tag" events. Learn more about [pulling images](../user/manage-images/pull-and-push-images/#pull-the-image). | Repository Admin |
| Mirror        |Refers to "Pull mirroring" and "Push mirroring" events. See [Mirror images to another registry](../user/promotion-policies/#mirror-images-to-another-registry) and [Mirror images from another registry](../user/promotion-policies/#mirror-images-from-another-registry) for more details. | Repository Admin |
| Create repo        | Refers to "Create Repository" events. See [Create a repository](../user/manage-images/) for more details. | Authenticated Users |

## Where to go next

- [Enable auto-deletion of repository events](../admin/auto-delete-repo-events.md)
