---
title: DTR 2.3 release notes
description: Learn about the new features, bug fixes, and breaking changes for Docker Trusted Registry
keywords: docker trusted registry, whats new, release notes
toc_max: 2
---

Here you can learn about new features, bug fixes, breaking changes, and
known issues for each DTR version.

You can then use [the upgrade instructions](admin/upgrade.md),
to upgrade your installation to the latest release.

## Version 2.3.6

(13 February 2018)

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

## 2.3.5

(20 November 2017)

**Bug fixes**

* Fixed a bug that caused certain vulnerabilities to not be found during scanning.
* Fixed a bug with downloading storage yaml file on Firefox.
* Increased the speed of lock expiration in case of failed joins.
* Changed storage backend to "local" when using the bootstrapper to switch to NFS.
* Fixed the notification when toggling active status of webhooks.
* Fixed a bug where garbage collection ran in a suboptimal mode if scheduled as
a cron from the UI.
* Fixed a potential issue with the way we untar files in uploads of the
vulnerability database.
* Fixed a bug with not backing up repository team permissions correctly.

**General improvements**

* Improved resilience of garbage collection.
* Improved logging of garbage collection.
* Improved memory usage during backup.
* Improved error handling when uploading invalid vulnerability databases.
* Fail faster in case of nfs volume issues.
* Improved resilience of DTR join operations.
* Hide secrets on storage config pages.

## DTR 2.3.4

(12 October 2017)

### Bugs fixed

* High severity
  * Fixed a bug in distribution that caused pull timeouts under load if using
  NFS or local storage. [#2299](https://github.com/docker/distribution/pull/2299)
  * Fixed GCS configuration UI.
* Low severity
  * Fixed missing show password button.
  * Removed incorrectly enforced length limit on repo names from UI.
  * Fixed small UI behavior and appearance inconsistencies.

## DTR 2.3.3

(13 September 2017)

### Bugs fixed

* High severity:
  * Fixed issue with RethinkDB not starting correctly after restarting a DTR
  replica.
  * Fixed issue that prevented UCP 2.1.x from being able to pull images.
* Low severity:
  * Improved error handling in the vulnerability scanner.
  * Fixed issue that caused webhooks not to fire when an image was automatically
  promoted.

## DTR 2.3.2

(25 August 2017)

### Bugs fixed

* High severity:
  * Add the ability to upgrade from 2.3.0 or 2.3.1

## DTR 2.3.1

(24 August 2017)

### Bugs fixed

* High severity:
  * Fixed a bug which caused upgrades to fail when upgrading from 2.2 to 2.3
  when DTR was previously upgraded from 2.1 to 2.2.
  * Make it possible to install DTR when the Docker daemon has SELinux enabled.
* Low severity:
  * In-product documentation for content cache now shows a simplified
  configuration format. The older configuration format is still supported for
  backwards compatibility.
  * When creating teams, the team name is validated in a way that's consistent
  with UCP.
  * The promotion policy creating form was sometimes disabled. This has been
  fixed.
  * Fixed bug that made users show as inactive, when they were actually active.

### Known issues

* You can't upgrade from 2.3.0 to 2.3.1. Upgrade to 2.3.2 directly.

## DTR 2.3.0

(16 August 2017)

### New features

* Repositories can now be marked “immutable”. Tags for images in immutable repos
cannot be changed or updated.
* You can now define promotion policies, to automatically copy images from one
repository to another.
* DTR now has UX for easily creating webhooks for events in repositories such
as image pushes, scans, deletions, and promotions.
* Added support for scanning Windows images for vulnerabilities.
* Support was added for handling manifest-lists for multi-architecture images.
This lets you manage images for different operating systems (eg. Linux and
  Windows) and CPU architectures (eg. x86_64 and s390x) under a single tag.
* You can now use the web UI in Chinese.

## General improvements

#### UI/UX

* The users page is now paginated to decrease long load times.
* Fixed the login page to be more resilient to improperly configured domain names.
Now you can always use the `/login` URL to bypass SSO misconfiguration issues.

#### docker/dtr

* Removed requiring the `--dtr-external-url` flag during DTR installation.
* Improved error handling.
* Added extended help to the installer with the `--help-extended` flag.
* We check kernel versions on install to help avoid known kernel bugs.

#### Storage
* Preserve the NFS hostname instead of resolving it during DTR installation.
* Allow using S3 compatible storage with customized certificates or without TLS
verification.

#### Misc

* Simplified content cache configuration.
* Removed old, insecure 3des from the cipher list.
* Allow refresh tokens with basic auth when using the DTR API.
* Added the ability to rescan all previously scanned images at once.
* Fixed issue with configs not being picked up by containers sometimes which
previously required a restart of all containers.

### Known issues

* When running DTR 2.3.0 + UCP 2.1.x, UCP users cannot pull images from DTR without logging in first.
* When using SSO with UCP 2.1.0, you have to log into DTR and UCP separately.
* Some users are displayed as inactive when they are actually active. This only
happens when using pagination, and can be fixed by refreshing the browser.
