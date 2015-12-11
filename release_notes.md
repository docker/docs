
+++
title = "Overview"
description = "Docker Universal Control Plane"
[menu.main]
identifier="mn_ucp"
+++


# Release Notes

The latest release is 0.5.  Consult with your Docker sales engineer for the release notes of earlier versions.

## Version 0.5

The following notes apply to this release:

## 3.16 kernel or higher

Your hosts must have a 3.16.0 kernel or higher to use UCP. If you don't have the proper kernel installed, the UCP bootstrapper returns this error.

```
INFO[0000] Verifying your system is compatible with UCP
FATA[0000] Your kernel version 3.13.0 is too old.  UCP requires at least version 3.16.0 for all features to work.  To proceed with an old kernel use the '--old-kernel' flag
```

If you don't want to use Docker's new networking features such as mult-host networking, use the '--old-kernel' flag to proceed anyway.

## New networking

This release includes support for the new networking features added in Docker Engine 1.0. These features include multi-host networking which allows you to configure custom container networks that span across several Docker hosts.

You must enable the networking features manually on your UCP cluster after bootstrapping each node.  To learn the process or to run through it, see [Set up container networking with UCP](networking.md) in the documentation.

## High Availability

This release includes support for setting up High Availability of your UCP Controller. In this mode additional replicas can be created when joining a new node to the swarm cluster.

## External logging

This release supports external logging facilities.

## New upgrade command

This version of UCP includes an `upgrade` command. You cannot upgrade from version 0.4 to version 0.5. You can use the `upgrade` command to upgrade from version 0.5 and future, newer versions of UCP.
