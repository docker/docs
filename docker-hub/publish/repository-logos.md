---
title: Repository logos
description: Override the organization logo on a per-repository basis
keywords: dvp, verified, publisher, repository, logo, icons
---

Docker Verified Publishers (DVP) and Docker Sponsored Open Source (DSOS)
organizations can upload custom images for individual repositories on Docker Hub.
This lets you override the default organization-level logo on a per-repository basis.

Only a user with administrative access (owner or team member with Admin permission)
over the repository can change the repository logo.

## Image requirements

- The supported filetypes for the logo image are JPEG and PNG.
- The minimum allowed image size in pixels is 120×120.
- The maximum allowed image size in pixels is 1000×1000.
- The maximum allowed image file size is 5MB.

## Set the repository logo

1. Sign in to Docker Hub.
2. Go to the page of the repository that you want to change the logo for.
3. Select the upload logo button, represented by a camera icon
   (![Camera icon](./images/upload_logo_sm.png){: .inline height="22px" })
   overlaying the current repository logo.
4. In the dialog that opens, select the PNG image that you want to upload to
   set it as the logo for the repository.

## Remove the logo

Remove a logo using the clear button (![Clear button](../images/clear_logo_sm.png){: .inline height="22px" }).

Removing the logo makes the repository fallback to using the organization logo, if set, and the default logo if not.

![Default logo which is a 3D grey cube](./images/default_logo_sm.png){: .inline height="22px" }
