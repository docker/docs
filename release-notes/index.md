<!--[metadata]>
+++
aliases = ["/docker-trusted-registry/release-notes/release-notes/"]
title = "Trusted Registry release notes"
description = "Docker Trusted Registry release notes "
keywords = ["docker, documentation, about, technology, understanding, enterprise, hub, registry, release notes, Docker Trusted Registry"]
[menu.main]
parent="dtr_menu_release_notes"
identifier="dtr_release_notes"
weight=0
+++
<![end-metadata]-->

# Docker Trusted Registry release notes

Here you can learn about new features, bug fixes, breaking changes and
known issues for each DTR version.

You can then use [the upgrade instructions](../install/upgrade/upgrade-major.md),
to upgrade your installation to the latest release.

## DTR 2.1

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
