---
description: Docker Trusted Registry release notes
keywords: docker, documentation, about, technology, understanding, enterprise, hub, registry, release notes, Docker Trusted Registry
title: Docker Trusted Registry release notes
---

Here you can learn about new features, bug fixes, breaking changes and
known issues for each DTR version.

You can then use [the upgrade instructions](../install/upgrade/index.md),
to upgrade your installation to the latest release.

## Version 2.0.4

(13 Oct 2016)

**General improvements**

* Increased limits on pagination for all lists in the UI
* Improved health check endpoints to report DTR status more accurately

**Bug fixes**

* You can now edit the permissions for a team member
* Fixed issue that prevented DTR from being installed in Docker Engine 1.12
* Several improvements to the migrate command
* Improved the reconfigure command to allow tuning the key-value store heartbeat
interval, election timeout, and snapshot count
* Users can now pull from public repositories in an organization namespace
without having to authenticate. Requires UCP 1.1.4

## Version 2.0.3

August 11, 2016

**Bug fixes**

* You can now add descriptions to the repositories using markdown

## Version 2.0.2

July 6, 2016

**docker/dtr image**

* Added the `docker/dtr images` command that lists all images needed for installing DTR.
* Added the `--extra-envs` flag to the install and join commands, that allows to specify
swarm constraints to the DTR containers.

**Misc**

* General improvements to the garbage collection settings screen.
* Deleting a repository is now faster.
* Upgraded the RethinkDB used internally to version 2.3.4, for improved security.
* Upgraded the Nginx used internally to version 1.8.1, for improved security.

**Bug fixes**

* Fixed problem that caused the last run time for the garbage collection job to
be wrong.
* When creating users and organizations, the 'Save & create another' button now
works as expected.
* In the Users screen, the organizations each user belongs to is now listed.

## Version 2.0.1

**New features**

* You can set a maximum duration for the garbage collection job to run. After
this timeout the garbage collection job stops. The job is then resumed on the
next scheduled interval.

* docker/dtr image
  * Added the `docker/dtr upgrade` command, to upgrade an installation.
  * Renamed the `--dtr-load-balancer` flag to `--dtr-external-url`.

**Deprecated features**

*  The user authentication settings page is deprecated and will be removed in a
future release. Configure user authentication in UCP.

**Bug fixes**

* UI
  * The dropdown filter in the Repositories screen works as expected.
  * The global search now returns all results.
  * When removing users from an organization, the page now refreshes
  automatically.
  * When configuring Swift as storage backend, it's now possible to define the
  Chunk size option.
  * When integrating with LDAP, the 'User search config' section is no longer
  being displayed twice.
  * When uploading a new license, the page now refreshes automatically.
  * On the repository screen, the docker pull command is now displayed
  correctly for all users.

* Users and teams
  * When adding users to a team, it's now possible to deselect users.
  * When listing users, the organizations those users belong to are displayed
  correctly.
  * After adding a user to a team, the 'Add user' modal now closes
  automatically.
  * When adding a user, errors are now displayed correctly.

* Browsers
  * When logging in with Microsoft Edge on Windows 10, users are no longer
  redirected indefinitely.
  * You can now update DTR license when using Firefox version 45.
  * You can now add managed users to teams when using Internet Explorer and
  Firefox.
  * You can now list repository tags when using Internet Explorer or Firefox.

* Misc
  * The restore command now gives an error when running restore without input
  redirection.
  * With more than 150k users, the search API is now responsive.
  * When UCP and DTR are running on the same host, restarting the host now
  restarts DTR containers. For this to work you need to use CS Engine
  1.11.1-cs2 and UCP 1.1.1, or newer versions.
  * It's no longer possible to delete the default 'docker-datacenter'
  organization.

**Known issues**

* UI
  * It's not possible to delete tags already deleted by the garbage collection
  job.
  * On the Users page, the pagination only works until the 20th page.
  * The navigation bar doesn't show what page is active.
  * After migrating from DTR 1.4.3 some tags are displayed on the UI but cannot
  be pulled because they are tags that are not associated with any image.
  Contact Docker support for help in removing these tags.
  * In the Repository screens, the filter dropdown doesn't reset after
  selecting an organization.

* Misc
  * When using Firefox 38.8.0, some buttons are not responsive and tags are not
  displayed. As a workaround use a more recent version of Firefox, or a
  different browser.
  * The docker search command doesn't work with DTR 2.0 and CS Engine 1.10+.
  This issue should be fixed on CS Engine 1.12.
  * When the DTR proxy container stops, it may seem that the DTR UI is
  responding, but it shows an "empty" notification when saving settings.
  * When configuring DTR to use a proxy, connections to the Docker Notary server
  are routed through the proxy. If the Docker Notary is not reachable through
  the proxy, you need to add the Notary server url to the NO_PROXY settings.

## Version 2.0.0

4 May 2016

**Features**

* Core
  * Support for high-availability and horizontal scalability.
  * UCP and DTR are now using a unified authentication service.
  * LDAP configurations should be performed on UCP. In future releases that
  configuration will be removed from DTR.
  * Users and teams created in the DTR ‘Datacenter’ organization are displayed
  in UCP.
  * Credentials of LDAP and Active Directory users are no longer stored on disk
  Instead a token is used for pushing and pulling images.
  * Added anonymous event tracking to allow us to improve DTR. The data is
  completely anonymized and you can turn this off in the DTR settings page.
* DTR installer
  * The DTR installer is now similar to UCP.
  * The DTR installer now runs on a container and has commands to install
  configure, and backup DTR.
* UI
  * Several changes to the UI to make it consistent with UCP.

**Bug fixes**

* Improved search functionality.
* Search autocomplete now displays all users.
* Image tags now take less time to be displayed.

**Known issues**

* UI
  * The dropdown filter in the Repositories screen doesn't open, but you can
  click on the input box.
  * The global search is returning a limited number of results.
  * When configuring Swift as storage backend, it's not possible to define
  the Chunk size option.
  * When logging in with Microsoft Edge on Windows 10, users are redirected
  indefinitely.
  * It's not possible to delete tags already delete by the garbage collector.
  * When integrating with LDAP, the 'User search config' section is sometimes
  displaying twice.
  * On the Users page, the pagination only works until the 20th page.
  * When navigating to detail pages, the navigation bar doesn't show what page
  is active.

* Users and teams
  * When removing a user from an organization, the page doesn't refresh
  automatically.
  * When adding users to a team, it's not possible to deselect users.
  * When uploading a new license, the page doesn't refresh automatically.
  * After adding a user to a team, the 'Add user' modal doesn't close
  automatically.

* API and CLI
  * With more than 150k users, the search API takes too long to return a
  response.
  * The docker search command does not work with DTR 2.0.
  * The docker/dtr restore command doesn't return a meaningful error when
  invoking the command without redirecting input from a file.

* Misc
  * When using Docker Engine 1.11 and DTR and UCP are running on the same host,
  restarting the host might cause the DTR containers not start. This is caused
  by a [bug in Docker Engine](https://github.com/moby/moby/issues/22486).
  You can restart the DTR containers from the UCP UI.
  * When the DTR proxy container stops, it may seem that the DTR UI is
  responding but it shows an "empty" notification when saving settings.
  * Adding a proxy to the DTR configuration forces the connections to the
  notary server to be routed through the proxy. If the notary server is not
  reachable through the proxy, add the notary URL to the NO_PROXY setting.
  * UCP has a single default organization “docker-datacenter”. If you delete
  this organization on DTR, you won’t be able to manage users and teams from
  UCP. To recover from this, contact Docker support.

## Prior versions

You can find the release notes for older versions of DTR on the
[release notes archive](prior-release-notes.md).