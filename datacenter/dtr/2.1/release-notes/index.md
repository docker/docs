---
title: Docker Trusted Registry release notes
description: Docker Trusted Registry release notes
keywords:
- docker, documentation, about, technology, understanding, enterprise, hub, registry, release notes, Docker Trusted Registry
---


Here you can learn about new features, bug fixes, breaking changes and
known issues for each DTR version.

You can then use [the upgrade instructions](../install/upgrade/upgrade-major.md),
to upgrade your installation to the latest release.

## DTR 2.1 Beta 3

(24 Oct 2016)

**Bug Fixes**

* Fixed Swift configuration error in DTR web UI when using advanced settings
* Fixed bug where organization owner would not see the delete button for repository
* Changed http response when deleting a non-existent tag from 204 to 404
* Changed http response when deleting a non-existent manifest from 500 to 404
* Users now show in the organization member list
* Team name is now being displayed in the UI
* Organization administrators can now delete repositories
* Fixed problem that prevented organization administrators to change a user role
* Fixed problem when removing users from an organization in the UI
* Fixed errors in Internet Explorer 11
* Fixed problem that caused the garbage collection job to treat valid manifests as corrupt

## DTR 2.1 Beta 1

(7 Oct 2016)

**docker/dtr image**

* Added more flags to the docker/dtr image to configure logging and tuning of
etcd when troubleshooting performance issues
* Added support to specify a custom volume or NFS mount to store the Docker images
* Several improvements to make installation command more stable

**Components**

* DTR now contains its own Notary server you can use to store secure image metadata
* Notary is highly-available if DTR is configured for high availability

**UI**

* Improved UI for configuring garbage collection jobs
* Removed user management pages. User management workflows can be done in UCP
* UI now shows author, push time, and layer sizes for tags

Additional tag data is available in the API including Tag Author, Time pushed, Layer sizes, and Dockerfile commands. The GUI includes the following tag data: Tag Author, Time pushed, and Layer sizes.

**General**

* Added support for S3 compatible storage
* Added support for Google Cloud storage
* Several improvements to garbage collection
* Improved DTR health checking API
* API now returns author, push time, layer size and other information for tags
