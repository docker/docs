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

## Version 2.0 Beta2

(24 Oct 2016)

** Features **

* Core
	* Added HTTP Routing Mesh feature to enable hostname routing for services on 
	top of Swarm's existing network mesh
	* It's now possible to install a new UCP cluster using a backup via the
	`install --from-backup` operation. Note that this will preserve the UCP cluster 
	configs but not any application services or networks from the old cluster
* UI/UX
	* Added "DDC Highlights" section to the main dashboard which provides sample 
	workflows to help users get started with DDC
	* Added a new color for node state, gray, which is used to notify the user
	to wait on using the node until it finishes an operation it is undergoing
	* The Create Network screen now has additional options for encryption, MTU,
	and multiple IPAM configs
	* It's now possible to load image .tar files directly from the UCP UI
	* It's now possible to edit a permission label directly on the Team screen
	* It's now possible to edit multiple node parameters prior to saving changes
	* You can now see a container's size in its detailed description view
	
** Bug Fixes **

* Core
	* Pre-loading UCP images onto a node is no longer requires prior to joining
	it as a manager or promoting it from worker to manager
	* Uninstall now properly removes all UCP containers
	* Auth now correctly normalizes capitalized letters for LDAP authentication
	* It is now possible to demote or remove manager nodes which are "down"
	* Support dumps should now download correctly despite stalling or timeouts
	from the cluster
	* Networks RBAC now works correctly for `service create` and `service update`
	* Admins can now successfully create/update services with bind mounts
	* Only admins are now allowed to stop UCP system containers
	* DAB deploy screen now correctly displays networks
* docker/ucp image
	* UCP now informs you if docker/ucp command has been deprecated or removed
	* License injection during install now works correctly
	* The `--root-ca-only` flag has been removed from backup/restore commands
* UI/UX
	* The Users tab page limits now supports up to 1 million LDAP/AD users
	* The Nodes page now displays more useful status updates for commands for
	joining, removing, and promoting/demoting nodes
	* Removed unnecessary scrollbars from Deploy Services wizard
	* Errors in the Services tab are now more descriptive
	* It's now more clear where to find and manage DDC system images
	* UCP now provides more clear warnings for issues related to upgrades, node
	management, backups, DAB deploys, and uploading certs 

## Version 2.0 Beta1

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
