---
description: Docker Trusted Registry release notes
keywords: docker trusted registry, whats new, release notes
redirect_from:
- /docker-trusted-registry/release-notes/release-notes/
- /docker-trusted-registry/release-notes/
title: Docker Trusted Registry release notes
---

Here you can learn about new features, bug fixes, breaking changes and
known issues for each DTR version.

You can then use [the upgrade instructions](install/upgrade.md),
to upgrade your installation to the latest release.

## DTR 2.1.7

(17 July 2017)

**Bugs fixed**

* Fixed registry DoS vulnerability.  CVE-2017-11468.  Severity: high
* Fixed small memory leak when handling batch jobs.  Severity: small

**General improvements**
* Added registry pprof endpoint if pprof is enabled.

## DTR 2.1.6

(13 April 2017)

**Bugs fixed**

* High impact
  * Fixed memory leaks causing DTR to use all RAM available

## DTR 2.1.5

(17 March 2017)

**Bug fixes**

* Fix memory leaks in the DTR API and Registry containers. Severity: high
* Fixed an issue causing RethinkDB to not start due to DNS errors when
the RethinkDB containers were not restarted at the same time. Severity: high
* The UI seemed to load indefinitely when viewing V1 manifests that didn't
include a 'size' property. To fix this, manifests without this field are no
longer displayed to the user. Severity: medium
* If you're using NFS as a storage backend, DTR had a chance of hanging when
reconfiguring or upgrading. This has now been fixed. Severity: medium
* When `--dtr-external-url` was set with uppercase letters, users would not be
able to log into DTR. This is now fixed. Severity: low
* The Users page was always showing a form to create new users. Now this form
only shows after clicking the 'New user' button. Severity: low


## DTR 2.1.4

(17 Jan 2017)

**Bug fixes**

* Fixed garbage collection UI slowdown when changing settings
* Fixed storage settings UI missing "Save" button when changing storage backends
* Fixed bug which was showing image tags as "outdated" for Notary signed images
* Removed `--log-tls-*` options which were not working correctly

## DTR 2.1.3

(20 Dec 2016)

**Bug fixes**

* docker/dtr image
  * Restore command now correctly prints error messages
  * Improved join command to retry after failure
* DTR web UI
  * UI now renders correctly when hiding the left navigation bar
  * You can now create organizations that use hyphens in their name
  * DTR now displays a UI banner when migrating tag data
  * Tag and manifest tags now render faster

## DTR 2.1.2

(8 Dec 2016)

**Features**

* The web UI now alerts when no backups have been made in a week


**Bug fixes**

* Restore operation now prints logs
* Google Cloud Storage driver now throttles data if there's heavy load, instead
of generating errors
* Upgraded Alpine images used by the DTR services to fix a [security
vulnerability with Expat2](https://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2016-4472)
* Fix for tag migration when pushing non-standard manifests
* Fix for tag migration failing during upgrade due to database timeouts


## DTR 2.1.1

(28 Nov 2016)

**Features**

* Updated backend storage configuration to use AWS v4 headers
* Added support for Scality, an Amazon S3 compatible object storage

**Other Improvements**

* Health check now reports failures after 3 consecutive failures
* Restore command now restores Notary server data
* Fix subsequent joins after a failed join


## DTR 2.1.0

(10 Nov 2016)

**Features**

* Out of the box integration between UCP and DTR. You no longer need to
configure UCP to trust DTR and vice versa. Requires UCP 2.0 or higher
* DTR now contains its own Notary server you can use to store secure image
metadata
* Notary is highly-available if DTR is configured for high availability
* Added support of Google Cloud Storage driver using YML configurations
* Added support for Amazon S3 compatible storages like Cleversafe object store
by IBM

**Installer**

Made several improvements to the DTR installer, and added more configuration
flag, for more customization at install time.

* Several improvements to make installation more stable
* Added the `--log-tls-ca-cert`, `--log-tls-cert`, `--log-tls-key`,
`--log-tls-skip-verify` for specifying the TLS certificates to be used
with the DTR logging driver
* Added the `--enable-pprof` to enable pprof profiling of the server
* Added the `--etcd-heartbeat-interval`, `--etcd-election-timeout`, and
`--etcd-snapshot-count` options to configure the key-value store used by DTR
* Added the  `--nfs-storage-url`, and `--dtr-storage-volume` options to allow
configuring where Docker images are stored

**Web UI**

* Web UI now displays information about tag metadata and logs
* Improved garbage collection settings

**General improvements**

* Better integration with NFS storage driver to store Docker images
* Better integration with Filesystem storage driver to store Docker images
* Improved garbage collection performance and efficiency
* Improved health checking API for more granularity

**Known issues**

* When upgrading to this version, tag metadata is migrated to DTR's internal
database. Depending on how many images are stored in DTR this can take some
time to complete.
