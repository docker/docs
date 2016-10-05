<!--[metadata]>
+++
alias=["/ucp/release_notes/"]
title ="Release Notes"
description="Release notes for Docker Universal Control Plane. Learn more about the changes introduced in the latest versions."
keywords = ["Docker, UCP", "Release notes", "Versions"]
[menu.main]
identifier="ucp-release-notes-current"
parent="ucp_menu_release_notes"
weight=0
+++
<![end-metadata]-->

# UCP Release Notes

Here you can learn about new features, bug fixes, breaking changes and
known issues for the latest UCP version.
You can then use [the upgrade instructions](../installation/upgrade-major.md), to
upgrade your installation to the latest release.

## Version 2.0 Beta

(7 Oct 2016)

**Docker swarm**

* UCP now leverages the swarm capabilities provided by Docker Engine 1.12
* Improved performance and scalability since engine-discovery mode and etcd are
no longer used for swarm inventory and overlay networks
* Smooth transition from container-based workflows to service-based workflows
* You run `docker logs`, `docker exec` and other container operations on
service tasks
* Adding nodes to UCP is easier, just run `docker swarm join`
* Now you can deploy and update services from the UI
* Experimental support for deploying Distributed Application Bundles

**docker/ucp image**

* Renamed the `uninstall` command to `uninstall-cluster`
* Checks if clocks are synchronized when joining nodes
* Uninstaller now removes all UCP volumes
* Improved help output and error messages

**UI/UX**

* Added tooltips and other messages to provide better in-context help
* All edit pages now use a side-modal for creating in-context
* Improved input validation across all the UCP UI
* Added a wizard for guiding users when deploying services
* Improved the dashboard page
* Improved UI for joining nodes to UCP
* Improved error page for 404 errors
* Improved page for customizing certificates

**General**

* Improved performance and scalability
* Several improvements to the authentication and authorization service
