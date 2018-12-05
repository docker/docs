---
description: Explains the differece between Classic and Enhanced Automated Builds
keywords: automated, build, images
title: Classic vs Enhanced Automated Builds
---

With the launch of the new Docker Hub, we are introducing Enhanced Automated Build.

Automated builds created using the old Docker Hub are now Classic Automated Builds, and automated builds created 
using the old Docker Cloud are now Enhanced Automated Build. 

All new automated builds created going forward will be Enhanced Automated Build. If you are creating an Enhanced
Automated Build for the first time, see [docs](index.md#configure-automated-build-settings).

In the coming months, we will gradually convert Classic Automated Builds into Enhanced Automated Builds. This should
be a seamless process for most users.


## Managing Classic Automated Builds

You can manage both Classic and Enhanced Automated Builds from the **Builds** tab

Repository with Classic Automated Build:

![A Classic Automated Build dashboard](images/classic-vs-enhanced-classic-only.png)

You can configure build settings pretty much the same way as the old Docker Hub.

If you have previously created an automated build in both the old Docker Hub and Docker Cloud, you can switch between 
Classic and Enhanced Automated Builds.

Enhanced Automated Build is displayed by default. You can switch to Classic Automated Build by clicking on this link at the top

![Switching to Classic Automated Build](images/classic-vs-enhanced-switch-to-classic.png)

Likewise, you can switch back to Enhanced Automated Build by clicking on this link at the top

![Switching to Enhanced Automated Build](images/classic-vs-enhanced-switch-to-enhanced.png)


## Frequently asked questions

1. I've previously linked my Github/Bitbucket account in the old Docker Hub. Do I need to re-link it?

Unless you've linked your Github/Bitbucket account with the old Docker Cloud previously, you will need to [re-link](link-source.md) it
with the new Docker Hub. 

2. What happens to my old Docker Hub automated builds?

They are now Classic Automated Builds. There are no functional difference with the old automated builds and everything 
(build triggers, existing build rules etc) should continue to work seamlessly

3. Is it possible to convert an existing Classic Automated Build into Enhanced Automated Build?

This is currently unsupported. However, we are working to transition all builds into Enhanced Automated Build in
the coming months.
