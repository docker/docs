---
description: Docker Trusted Registry release notes
keywords: docker trusted registry, whats new, release notes
title: DTR 2.2 release notes
---

Here you can learn about new features, bug fixes, breaking changes and
known issues for each DTR version.

You can then use [the upgrade instructions](index.md),
to upgrade your installation to the latest release.

## DTR 2.2.0

(9 Feb 2017)

**New features**

* DTR can now scan the binaries contained in the image layers, and report
security vulnerabilities
* You can now configure multiple caches, so that users can pull images faster
* You can now configure webhooks to run automated tasks that are triggered by
events like image push, repository creation, and others

**General improvements**

* UI/UX
  * Improved error messages to be more meaningful and help troubleshoot the problem
  * Several UI/UX improvements to the DTR configuration page and user settings page
  * Several improvements to the search bar used in the UI

* docker/dtr image
  * The `docker/dtr install` command now shows all the nodes that are part of a
  UCP cluster for you choose on which node to deploy DTR
  * The install command was improved to avoid deploying DTR to a node where it
  cannot run due to port collisions
  * The `docker/dtr install --ucp-node` flag is now mandatory
  * The install command no longer allows deploying replicas with duplica ids
  * The upgrade command now validates if all tags were migrated to the latest
  version before trying to migrate blob links

**Bugs fixed**

* When creating a repository, the length of the repository now is consistent
between the UI and API
* The UI now validate and doesn't allow creating repository names with uppercase
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
* When you don't have permissions to see the repository details, the UI now
shows that you don't have permissions instead of saying it has no manifests
* Jobs are retried if the worker running them stops unexpectedly

**Deprecation**

The `/load_balancer_status` is deprecated and is going to be removed in future
versions. Use the `/health` endpoint instead.
