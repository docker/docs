---
title: DTR 2.4 release notes
description: Learn about the new features, bug fixes, and breaking changes for Docker Trusted Registry
keywords: docker trusted registry, whats new, release notes
toc_max: 2
---

> **DTR 2.5 now available**. [You can check the release notes here](/ee/dtr/release-notes.md).

Here you can learn about new features, bug fixes, breaking changes, and
known issues for each patch release of DTR 2.4.
You can then use [the upgrade instructions](admin/upgrade.md),
to upgrade your installation to the latest release.

## Version 2.4.3 (2018-03-19)

**Security**

* Dependencies updated to consume upstream CVE patches.

## Version 2.4.2 (13 February 2018)

**Security notice**

The log driver is now disabled for containers started by backup and HA cluster
join operations. This is a critical security fix for customers that rely on
Docker Trusted Registry 2.2, 2.3 and 2.4 with a log driver to capture logs from
all containers across the platform.

Caution is advised when applying this update, make sure you redeploy DTR, and in
the process you will create new credentials because the previous ones were
potentially disclosed due to the vulnerability.

Use the `--log-driver=none` option for `docker run` when running a DTR backup, HA
cluster join or dumpcerts.

## 2.4.1 (20 November 2017)

**Bug fixes**

* Fixed a bug that cause certain vulnerabilities to not be found during scanning.
* Increased speed of lock expiration in case of failed joins.
* Fixed notification when toggling active status of webhooks.
* Speed up detection of dead jobrunners.
* Fixed a bug where garbage collection ran in a suboptimal mode if scheduled as
a cron from the UI.
* Fixed a potential issue with the way we untar files in uploads of the
vulnerability database.
* Fixed scanning issue with some windows images.
* Fixed a bug with not backing up repository team permissions correctly.

**General improvements**

* Improved resilience of garbage collection.
* Improved logging of garbage collection.
* Improved memory usage during backup.
* Improved error handling when uploading invalid vulnerability databases.
* Improve resilience of DTR join operations.
* Hide secrets on storage config pages.

**Deprecations**

* The `api/v0/imagescan/layer/{layerid}` endpoint is deprecated, and will be
removed in DTR 2.5. You can use the
`/api/v0/imagescan/repositories/{namespace}/{reponame}/{tag}` endpoint instead.


## DTR 2.4.0 (2 November 2017)

**New features**

* Upgraded to Swagger 2.0 and Swagger UI 3.0.
* DTR can now be deployed on IBM Z (s390x architecture).
* Updated the `docker/dtr-rethink` images to include `rethinkcli` for easier troubleshooting.
* Notary now allows you to see audit logs using the `/v2/_trust/changefeed`, and
`/v2/<repository>/_trust/changefeed` endpoints.

**Bug fixes**

* When setting up periodic garbage collection, it used to run in a different,
less thorough mode than when run manually. Now garbage collection always runs in
the correct mode.
* Fix error when garbage collecting manifest lists.
* Fixed issue when reconfiguring DTR from a non-local storage to NFS, causing the
change to not be persisted.
* Backported Docker Distribution race fix. [2299](https://github.com/docker/distribution/pull/2299)
* Reduced unnecessary logs in Jobrunner.
* Other general reliability improvements.

**Known issues**

* Backup uses too much memory and can cause out of memory issues for large databases.
* The `--nfs-storage-url` option uses the system's default NFS version instead
of testing the server to find which version works.
