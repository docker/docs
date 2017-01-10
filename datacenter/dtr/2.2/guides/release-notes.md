---
description: Docker Trusted Registry release notes
keywords: docker trusted registry, whats new, release notes
title: DTR 2.2 release notes
---

Here you can learn about new features, bug fixes, breaking changes and
known issues for each DTR version.

You can then use [the upgrade instructions](install/upgrade.md),
to upgrade your installation to the latest release.

## DTR 2.2.0 beta 1

(10 Jan 2017)

**New features**

* DTR can now scan the binaries contained in the image layers, and report
security vulnerabilities
* You can now configure multiple caches, so that users can pull images faster
* You can now configure webhooks to run automated tasks that are triggered by
events like image push, repository creation, and others

**General improvements**

* Improved error messages to be more meaningful and help troubleshoot the problem
* Several UI/UX improvements to the DTR configuration page and user settings page
* Several UI/UX improvements to the user settings page
* Several improvements to the search bar used in the UI

**Bugs fixed**

* When creating a repository, the length of the repository now is consistent
between the UI and API
* The UI now validate and doesn't allow create repository names using uppercase
letters
* You can now create organizations with dashes in the name
* Fixed a bug that didn't allow deleting users immediately after they were
created
* The copy to clipboard button on the repository page now works on Firefox
* The repository page now renders properly the repository permissions
* You can now delete a users's full name from the UI
* Organization administrators can now see the repositories owned by the organization
* The garbage collection settings now show the correct cron values
* You can now specify DTR to use port 443 when installing DTR

**Known issues**

* When viewing the result of a security scan, clicking on a layer sometimes
highlights two different layers
* The `docker search` command is not returning exact matches to the user namespace
* The UI becomes slow when synchronizing more than 400k LDAP users
* The flag `docker/dtr reconfigure --tls-syslog-certs` may break the connection
to syslog since not all replicas have the correct TLS certificates locally
* The `docker/dtr remove` sometimes doesn't remove all the DTR volumes
