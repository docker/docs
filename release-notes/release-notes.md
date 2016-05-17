<!--[metadata]>
+++
aliases = ["/docker-trusted-registry/release-notes/"]
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

## Version 2.0

(4 May 2016)

**Features**

* Core
  * Support for high-availability and horizontal scalability,
  * UCP and DTR are now using a unified authentication service,
  * LDAP configurations should be performed on UCP. In future releases that
  configuration will be removed from DTR,
  * Users and teams created in the DTR ‘Datacenter’ organization are displayed
  in UCP
  * Credentials of LDAP and Active Directory users are no longer stored on disk.
  Instead a token is used for pushing and pulling images,
  * Added anonymous event tracking to allow us to improve DTR. The data is
  completely anonymized and you can turn this off in the DTR settings page.
* DTR installer
  * The DTR installer is now similar to UCP,
  * The DTR installer now runs on a container and has commands to install,
  configure, and backup DTR,
* UI
  * Several changes to the UI to make it consistent with UCP.

**Bug fixes**

* Improved search functionality,
* Search autocomplete now displays all users,
* Image tags now take less time to be displayed.

**Known issues**

* UI
  * The dropdown filter in the Repositories screen doesn't open, but you can
  click on the input box,
  * The global search is returning a limited number of results,
  * When configuring Swift as storage backend, it's not possible to define the Chunk size option,
  * When logging in with Microsoft Edge on Windows 10, users are redirected indefinitely,
  * It's not possible to delete tags already delete by the garbage collector,
  * When integrating with LDAP, the 'User search config' section is sometimes displaying twice,
  * On the Users page, the pagination only works until the 20th page,
  * When navigating to detail pages, the navigation bar doesn't show what page is active,

* Users and teams
  * When removing a user from an organization, the page doesn't refresh automatically,
  * When adding users to a team, it's not possible to deselect users,
  * When uploading a new license, the page doesn't refresh automatically,
  * After adding a user to a team, the 'Add user' modal doesn't close automatically,

* API and CLI
  * With more than 150k users, the search API takes too long to return a response,
  * The docker search command does not work with DTR 2.0,
  * The docker/dtr restore command doesn't return a meaningful error when invoking the
  command without redirecting input from a file.

* Misc
  * When using Docker Engine 1.11 and DTR and UCP are running on the same host, restarting
  the host might cause the DTR containers not start. This is caused by a [bug in
  Docker Engine](https://github.com/docker/docker/issues/22486).
  You can restart the DTR containers from the UCP UI.
  * When the DTR proxy container stops, it may seem that the DTR UI is responding,
  but it shows an "empty" notification when saving settings,
  * Adding a proxy to the DTR configuration forces the connections to the
  notary server to be routed through the proxy. If the notary server is not
  reachable through the proxy, add the notary URL to the NO_PROXY setting,
  * UCP has a single default organization “docker-datacenter”. If you delete
  this organization on DTR, you won’t be able to manage users and teams from
  UCP. To recover from this, contact Docker support.

## Prior versions

You can find the release notes for older versions of DTR on the
[relese notes archive](prior-release-notes.md).
