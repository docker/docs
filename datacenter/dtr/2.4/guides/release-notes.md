---
title: DTR 2.4 release notes
description: Learn about the new features, bug fixes, and breaking changes for Docker Trusted Registry
keywords: docker trusted registry, whats new, release notes
toc_max: 2
---

Here you can learn about new features, bug fixes, breaking changes, and
known issues for each DTR version.

You can then use [the upgrade instructions](admin/upgrade.md),
to upgrade your installation to the latest release.

## DTR 2.4.0

(2 November 2017)

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
