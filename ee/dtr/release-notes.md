---
title: DTR 2.5 release notes
description: Learn about the new features, bug fixes, and breaking changes for Docker Trusted Registry
keywords: docker trusted registry, whats new, release notes
toc_max: 2
redirect_from:
  - /datacenter/dtr/2.5/guides/release-notes/
---

Here you can learn about new features, bug fixes, breaking changes, and
known issues for each DTR version.

You can then use [the upgrade instructions](admin/upgrade.md),
to upgrade your installation to the latest release.

## 2.5.0 (2018-4-17)

### New features

* You can now configure DTR to automatically create a new repository when
users push to a repository in their personal namespace that doesn't exist yet.
This makes the behavior of DTR consistent with Docker Hub. By default this
setting is disabled, so that DTR continues behaving the same way after an upgrade.
[Learn about creating repositories on push](admin/configure/allow-creation-on-push.md).
* You can create push mirroring policies to automatically push an image to
another DTR deployment or Docker Hub, when the image complies with a policy
of your choice.
[Learn about push mirroring](user/promotion-policies/push-mirror.md).
* You can configure a repository in a DTR deployment to mirror a repository
in a different DTR deployment by constantly monitoring it and pulling new
images when they are available.
[Learn about pull mirroring](user/promotion-policies/pull-mirror.md).
* Added the `emergency-repair` command to the DTR CLI tool. This allows you to
recover your DTR cluster from a loss of quorum and is an alternative to
restoring from a backup.
[Learn about the emergency-repair command](admin/disaster-recovery/repair-a-cluster.md).
* Users can now create access tokens that can be used to authenticate in the
DTR API without providing their credentials.
[Learn about access tokens](user/access-tokens.md).
* You can now configure DTR to run garbage collection jobs without putting DTR
in read-only mode. This is still experimental.
[Learn about garbage collection](admin/configure/garbage-collection.md).
* Administrators can hide vulnerabilities in given image layers if they
know that the vulnerability has been fixed.
[Learn how to override vulnerability reports](user/manage-images/override-a-vulnerability.md)
* You can now connect one DTR deployment to multiple UCPs, allowing you to
use Docker Content Trust in a seamless way between multiple UCPs.
* Added new endpoints to the DTR API to query the results of the Vulnerability
scanner:
  * `/api/v0/imagescan/scansummary/repositories/{namespace}/{reponame}/{tag}` returns
  the scanning summary for a given tag.
  * `/api/v0/imagescan/scansummary/cve/{cve}` gets the scan summary by CVE.
  * `/api/v0/imagescan/scansummary/layer/{layerid}` gets the scan summary by layer SHA.
  * `/api/v0/imagescan/scansummary/license/{license}` gets the scan summary by
  license type.
  * `/api/v0/imagescan/scansummary/component/{component}` get the scan summary by
  component.
* The API endpoint `/api/v0/repositories/{namespace}/{reponame}/manifests/{reference}`
has been deprecated. Use `/api/v0/repositories/{namespace}/{reponame}/tags/{tag}`
instead.

### Bug fixes

* UI
  * Several improvements to the UI to make it more stable
* User accounts
  * When a user changes their password they are automatically logged out.
* Vulnerability scanner
  * Fixed problem causing errors when trying to view scanning information when
an image has not been scanned yet.
* docker/dtr tool
  * When using `docker/dtr reconfigure --log-host`, you now need to also
specify `--log-protocol`.
  * You can now tune the RethinkDB cache size for improved performance. Use the
  `--replica-rethinkdb-cache-mb` option available on install, join, or reconfigure.
* Misc
  * Removed support for manifest schema v1. This doesn't  affect users.

### Known issues

* Web UI
  * The web UI shows "This repository has no tags" in repositories where tags
  have long names. As a workaround, reduce the length of the name for the
  repository and tag.
  * When deleting a repository with signed images, the DTR web UI no longer
  shows instructions on how to delete trust data.
  * There's no UI support to update mirroring policies when rotating the TLS
  certificates used by DTR. Use the API instead.
  * The UI for promotion policies is currently broken if you have a large number
  of repositories.
  * Clicking "Save & Apply" on a promotions policies doesn't work.
* Web hooks
  * There is no web hook event for when an image is pulled.
* Online garbage collection
  * The events API won't report events when tags and manifests are deleted.
  * The events API won't report blobs deleted by the garbage collection job.
* Docker EE Advanced features
  * Scanning any new push after metadatastore migration will not yet work.
  * Pushes to repos with promotion policies (repo as source) are broken when an
  image has a layer over 100MB.
  * On upgrade the scanningstore container may restart with this error message:

  ```
  The data directory was initialized by PostgreSQL version 9.6,
  which is not compatible with this version 10.3.
  ```
  To remedy this remove the scanningstore container and the dtr-postgres volume, and perform a reconfigure
  ```
  docker stop dtr-scanningstore-$replica
  docker rm dtr-scanningstore-$replica
  docker volume rm dtr-postgres-$replica
  docker run --rm -it docker/dtr:2.5.0 reconfigure ...
  ```
  Then resync the vulnerability database from the web UI.


## Earlier versions

- [DTR 2.4 release notes](/datacenter/dtr/2.4/guides/release-notes.md)
- [DTR 2.3 release notes](/datacenter/dtr/2.3/guides/release-notes.md)
- [DTR 2.2 release notes](/datacenter/dtr/2.2/guides/release-notes/index.md)
- [DTR 2.1 release notes](/datacenter/dtr/2.1/guides/release-notes.md)
- [DTR 2.0 release notes](/datacenter/dtr/2.0/release-notes/index.md)
