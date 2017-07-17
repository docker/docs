---
title: DTR release notes
description: Learn about the new features, bug fixes, and breaking changes for Docker Trusted Registry
keywords: docker trusted registry, whats new, release notes
---

Here you can learn about new features, bug fixes, breaking changes and
known issues for each DTR version.

You can then use [the upgrade instructions](../admin/upgrade.md),
to upgrade your installation to the latest release.

## DTR 2.3.0 beta 1

(17 July 2017)

### New Features

* Repositories can now be configured to be “immutable” in that any existing
image can not be overwritten or deleted.
* You can now define promotion policies, to automatically copy images from one
repository to another.
* DTR now has UX for easily creating webhooks for events in repositories such
as image pushes, scans, deletions, and promotions.
* Added support for scanning Windows images for vulnerabilities.
* Support was added for handling manifest-lists for multi-architecture images.
This allows pulling different images with the same tag such as Linux and Windows.
* You can now use the web UI in Chinese.

##  General Improvements

### UI/UX

* The users page is now paginated to decrease long load times.
* Fixed the login page to be more resilient to improperly configured domain names.

### docker/dtr

* Removed requiring the `--dtr-external-url` flag during DTR installation.
* Improved error handling.
* Added extended help to the installer with the `--help-extended` flag.

### Storage
* Preserve the NFS hostname instead of resolving it during DTR installation.
* Allow using S3 compatible storage with customized certificates or without TLS
verification.

### Misc

* Simplified content cache configuration.
* Removed old, insecure 3des from the cipher list.
* Allow refresh tokens with basic auth when using the DTR API.
* Added the ability to automatically rescan any previously scanned image.

## Known issues

* The web UI works best with Chrome and Firefox. Full support for
Internet Explorer will be added on the GA release
