+++
title = "Release Notes"
description = "Docker Universal Control Plane"
[menu.ucp]
weight="-99"
+++


# UCP Release Notes

The latest release is 0.6.  Consult with your Docker sales engineer for the
release notes of earlier versions.

## Version 0.6

The following notes apply to this release:

## Licensing functionality enabled

This release enables UPC licensing functionality. UCP starts in "unlicensed" mode. This mode does not limit any functionality. However, when you start UCP, an informational banner appears on the application noting it is unlicensed.

For the purposes of this beta, you can use an existing DTR license and remove the banner by doing the following:

1. Install UCP.

2. Log into the controller.

3. Go to the **Settings** page.

    ![License](../images/license.png)

4. Upload a valid DTR license.

    After a successful upload, the banner disappears.

## UI

- Sidebar is now permanently visible and is responsive for smaller screens
- Fixed issue with hidden item count on Applications View
- Dashboard chart enhancements
- Disable stats and exec on stopped container
- UI Link to UCP documentation
- Tags displayed on image removal dialog
- Enhanced breadcrumb tracking
- Licensing configuration

## Images

- CFSSL updated to 1.1.0

## Misc

- Banner is now shown reporting high availability status
