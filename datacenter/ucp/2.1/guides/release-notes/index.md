---
title: UCP 2.1 release notes
description: Release notes for Docker Universal Control Plane. Learn more about the
 changes introduced in the latest versions.
keywords: Docker, UCP, release notes
redirect_from:
- /datacenter/ucp/2.1/guides/admin/upgrade/release-notes/
---

Here you can learn about new features, bug fixes, breaking changes, and
known issues for the latest UCP version.
You can then use [the upgrade instructions](../admin/upgrade.md), to
upgrade your installation to the latest release.

## Version 2.1.8 (2018-04-17)

* Fixed an issue that allows users to incorrectly interact with local volumes.

## Version 2.1.7

(13 February 2018)

**Security Notice**

The user must use `--log-driver=none` to disable the log driver for containers
started by backup operations. This is a critical security fix for customers that
rely on Universal Control Plane 2.1 and a log driver to capture logs from all
containers across the platform.

Caution is advised: any sensitive information that has already been disclosed in
the logs will NOT be removed by this update. Sensitive information needs to be
purged manually from the logs.
Use the backup encryption mechanism with the `--passphrase` option when running a
UCP backup.

A full credentials re-generation and update transition procedure is available:
[https://success.docker.com/article/KB000623](https://success.docker.com/article/KB000623)

This is a breaking change on UCP backup operation. It is now mandatory to specify
`--log-driver none` option for `docker run` for all UCP backups.

## Version 2.1.6

(16 January 2018)

**Bug fixes**

* Security
  * Role-based access control is now enforced for volumes managed by 3rd party
volume plugins (for example using the NetApp or other volume plugins). This is a
critical security fix for customers that use 3rd party volume drivers and rely
on Docker Universal Control Plane for tenant isolation of workloads and data.
**Caution** is advised when applying this update because users or automated
workflows may have come to rely on lack of access control enforcement when
manipulating volumes created with 3rd party volume plugins.

## Version 2.1.5

(20 July 2017)

**Security Update**

* Remediated a privilege escalation where an authenticated user could obtain
admin-level privileges

This issue affects UCP versions 2.0.0-2.0.3 and 2.1.0-2.1.4. It was discovered
by our development team during internal testing

**Bug Fixes**

* Core
    * Fixed an issue where clients misusing the events API (e.g. slowly reading
    or failing to read events) leads to unresponsive behavior from the cluster
    * Fixed an issue where app services pulling DTR private images using
    integrated single-sign-on would fail due to token expiration
    * UCP resource metrics now correctly display CPU utilization on newer Linux
    kernels
    * Fixed an issue where UCP incorrectly reported 100% memory usage on a node
    due to the usage of memory constraints on containers
    * Network and volume label filters now work correctly on UCP (for example
    when using `docker volume ls --filter label="foo"="bar")`
    * UCP can now be installed correctly when SELinux enforcement mode is
    enabled (e.g. `--selinux-enabled`)
    * Fixed an issue where rejoining (or demoting and promoting) a manager node
    caused `ucp-kv` to become unhealthy due to a stale KV cache
    * UCP now exposes a Registry field in `docker info` output, so that
    deploying with registry credentials (e.g. `docker stack deploy --with-registry-auth`
    now works correctly
    * UCP now reports percentage progress while pulling images
    * `docker images -f dangling=true` now correctly lists untagged `<none>`
    images instead of listing all images
    * Added a network diagnostic tool to `ucp-dsinfo` image to aid in troubleshooting
    issues related to overlay networks
    * Added additional diagnostic information about `docker stacks` to support dumps
    for troubleshooting purposes
    * UCP now provides a more informative warning banner and clearer logs when
    `ucp-auth-store` is unhealthy
    * Reduced the default cache size for `ucp-auth-store` to free up memory on the UCP manager.
    This cache can be adjusted via the `RethinkDBCacheSize` parameter in the UCP Config API
    * Various performance improvements made to `ucp-auth-store` to reduce overhead when the API
    is being repeatedly accessed in a short period of time
    * Fixed an issue where one `ucp-auth-store` instances would fail to join the HA
    cluster if started in the wrong order
    * Fixed an issue where a UCP manager might get stuck in a restart loop due to
    being unable to correctly access the root CA
    * Fixed an issue where users with view-only permissions received an access denied
    error when attempting to deploy stacks via the Compose UI, despite having been granted
    label access to do so


## Version 2.1.4

(4 May 2017)

**Bug Fixes**

* Core
	* Fixed an issue where updating the UCP server certificates, the web UI would
	report success, but not make any changes
	* UCP no longer shows an `invalid memory address` or `nil pointer dereference
	panic` when inspecting containers created with Docker 1.10 or older
	* It is no longer possible to create a service with the same published ingress
	port as the UCP controller's port, thereby rendering UCP inaccessible
	* Fixed an issue where usernames with special language characters (such as ä)
	were unable to login to the system
	* Fixed an issue where a Compose stack deploy could not update an existing service
	due to access control conflicts with the `com.docker.ucp.access.owner` label

* docker/ucp image
	* UCP support dumps now include `docker stats` output

* UI/UX
	* Fixed an issue where an application deployed using `docker stack deploy`
	in the CLI did not show up in the web UI
	* Fixed an issue where deploying a Compose application via UI with a slow network
	connection might display a websocket error despite successful deployment


## Version 2.1.3

(4 Apr 2017)

**Known issues**

In UCP 2.1.3, if you try to upload externally-signed controller
certificates through the **Admin Settings** page on the UI, you see a
"Success" message, but the certificates isn't updated on any
of the controller nodes.

The workaround is to update the contents of the `ucp-controller-server-certs`
volume manually on each manager node with the new `ca.pem`, `cert.pem`, and
`key.pem` contents. Update all three of these files approximately
simultaneously, to avoid issues with reconciliation.

**Bug fixes**

* Core
	* Fixed known issue where worker nodes would be left in a pending state
	  after upgrading from UCP 1.1.z.
	* Nodes will no longer be reported as unhealthy if the `ucp-reconcile`
	  container is removed.
	* Fixed an issue where nodes in the same subnet may report incorrect
	  hostnames in the UCP node list.

* UI/UX
	* UCP support dumps and client bundles can now be downloaded on IE10/11.
	* The task counter in the services page should now correctly omit tasks that
	  have not been assigned to a node yet.

## Version 2.1.2

(29 Mar 2017)

**Known issues**

There is known issue in UCP 2.1 where upgrading from UCP 1.1.z can cause swarm
to leave worker nodes in a pending state with the message:

```
[Pending] Completing node registration
```

There are two workarounds for rectifying this issue:

1. When upgrading from UCP 1.1.z, first upgrade to UCP 2.0.z, and then to UCP
2.1.z. This will prevent the issue from happening, and is the recommended upgrade path.
2. If you have already upgraded from UCP 1.1.z directly to UCP 2.1.z, you can
fix the issue by restarting the ucp-swarm-manager container on each of your UCP
controller nodes.

This issue will be fixed in UCP 2.1.3.

**Bug fixes**

* Core
	* `ucp-reconcile` service now correctly brings up `ucp-kv` container if it
	has stopped or become unreachable
	* Fixed known issue in which users are unable to log into UCP UI after upgrading
	from UCP 2.1.0 to 2.1.1 because the parameter for maximum concurrent users was
	incorrectly defaulted to `0`
	* Fixed an issue where the UCP manager becomes unresponsive and requires a restart
	if `docker ps` or `docker info` calls to engine take a long time for a response
	* HTTP Routing Mesh now correctly provides httplog for debug logging of services
	* `docker node ls -f` now correctly filters when run against a UCP cluster
	* `docker inspect task` no longer returns errors when run against a UCP cluster
	* UCP now correctly reports progress when loading an image from CLI

* docker/ucp image
	* UCP support dumps now include Docker Engine daemon logs
	* Host address IPs are now automatically added to SANs during install
	* UCP now reports its version number in the CLI after being installed

* UI/UX
	* Deploying Compose-based applications in the GUI now works correctly when
	Docker Content Trust "Run Only Signed Images" is turned on
	* Fixed an issue where UI temporarily showed more tasks for a service than
	actually existed
	* Fixed an issue in which metrics incorrectly displayed `0%` in the UI


## Version 2.1.1

(14 Mar 2017)

**Known issues**

If you are currently running UCP 2.1.0 and previously customized the sessions
lifetime parameter in the Authentication settings UI, upgrading to UCP 2.1.1 may
cause users to not be able to log into UCP and DTR. This is caused by a faulty
default value which sets maximum concurrent user sessions to zero.

You can either wait for UCP 2.1.2 to be released so that the problem is
automatically fixed, or upgrade to 2.1.1, and use the following steps to fix
the problem.

Start by getting the current configuration for user sessions by running:

```bash
curl -u admin "https://$UCP_HOST/enzi/v0/config/sessions"
```

The command will prompt for the `admin` user's password and then return
the current sessions config which should look something like:

```json
{
  "lifetimeHours": 72,
  "renewalThresholdHours": 24,
  "perUserLimit": 0
}
```

If `perUserLimit` is set to `0`, you need to set it to a value between 1 and 100.
The recommended value is 5. You should also customize the command below with
the `lifetimeHours` and `perUserLimit` values returned by the first command.

```bash
curl -u admin "https://$UCP_HOST/enzi/v0/config/sessions" \
  -X PUT \
  -H 'Content-Type: application/json' \
  -d '{"lifetimeHours": 72, "renewalThresholdHours": 24, "perUserLimit": 5}'
```

You can now log into UCP and DTR.

**New features**

* Core
  * Administrators can now configure the frequency with which UCP polls metrics.
  Use `docker service update --env-add METRICS_SCRAPE_INTERVAL=10m ucp-agent`,
  and the frequency can be in s/m/h.
  * Administrators can now configure the frequency with which UCP gathers disk usage data.
  Use `docker service update --env-add METRICS_DISK_USAGE_INTERVAL=12h ucp-agent`,
  and the frequency can be in s/m/h.
  * Support for syncing users and teams from multiple LDAP servers/domains
  (e.g. a separate server to use for `dc=domain2,dc=example,dc=com`)
  * Support for limiting the number of maximum concurrent login sessions any
  user may have

**Bug fixes**

* Core
  * Fixed an issue in which UCP manager would panic and be unable to return
  the right system status after the cluster became unhealthy
  * `ucp-hrm` container now provides debug logs through `stdout`
  * HTTP Routing Mesh now checks to ensure an ingress port is not already
  in use by UCP or DTR before becoming active
  * Fixed an issue in which UCP did not use swarm-mode node IDs, preventing
  usage of node constraints and other features when using cloned VMs as UCP nodes
  * Fixed an issue in which certain Docker API 1.26 commands were not correctly supported
  * Disk usage metrics no longer display 0% when using devicemapper filesystem
  * Disk usage metrics are now collected every 2 hours by default, and can be tunned
  * Fixed an issue causing Content Trust enforcement to ignore an optional `tag` for
  `/images/create`, causing some signed content to not run correctly
  * LDAP sync logs now take up less disk space on manager nodes
  * UCP support dumps are now correctly compressed to take up less disk space,
  and provide information on HTTP Routing Mesh and metrics
* docker/ucp image
    * UCP install now correctly fails and presents an error when trying to
    specify `host-address` to an existing swarm-mode cluster
    * Clarified upgrade message to make it clear that the upgrade command now
    works at once for the entire cluster rather than needing to be run on every
    node
* UI/UX
    * UI now displays a warning if there is significant latency or network issues
    in communications between UCP manager nodes
    * UI no longer incorrectly displays 'No Services' while still loading the
    Services tab
    * UI no longer displays errors when global tasks are removed due to node
    constraints
    * UI now displays a warning when underlying engines in the swarm-mode
    cluster are running different versions
    * UI now displays an error when 'Load Image' command fails
    * 'KV Store Timeout' option now displays correct units (milliseconds)
    * Dashboard now correctly displays errors when metrics are unavailable
    * The DTR deployment page now validates if a DTR replica ID is valid or not

## Version 2.1.0

(9 Feb 2017)

This version of UCP extends the functionality provided by CS Docker Engine
1.13. Before installing or upgrading this version, you need to install CS
Docker Engine 1.13 in the nodes that you plan to manage with UCP.

**New features**

* Core
  * Support for managing secrets (e.g. sensitive information such as passwords
  or private keys) and using them when deploying services. You can store secrets
  securely on the cluster and configure who has access to them, all without having
  to give users access to the sensitive information directly
  * Support for Compose yml 3.1 to deploy stacks of services, networks, volumes,
  and secrets.
  * HTTP Routing Mesh now generally available. It now supports HTTPS passthrough
  where the TLS termination is performed by your services, Service Name Indication
  (SNI) extension of TLS, multiple networks for app isolation, and Sticky Sessions
  * Granular label-based access control for secrets and volumes
  (NOTE: unlike other resources controlled via label-based access control, a
  volume without a label is accessible by all UCP users with Restricted Control
  or higher default permissions)

* UI/UX
  * You can now view and manage application stacks directly from the UI
  * You can now view cluster and node level resource usage metrics
  * When updating a service, the UI now shows more information about the service status
  * Rolling update for services now have `failure-action` which you can use to
  * Several improvements to service lifecycle management
  specify rollback, pausing, or continuing if the update fails for a task
  * LDAP synching has more configuration options for extra flexibility
  * UCP now warns when the cluster has nodes with different Docker Engine versions
  * The HTTP routing mesh settings page now lists all services using the
  routing mesh, with details on parameters and health status
  * Admins can now view team membership in a user's details screen
  * You can now customize session timeouts in the authentication settings page
  * Can now mount `tmpfs` or existing local volumes to a service when deploying
  services from the UI
  * Added more tooltips to guide users on the above features

**Bug fixes**

* Core
    * HTTP routing mesh can now be enabled or reconfigured when UCP is configured
    to only run images signed by specific teams
    * Fixed an error in which `_ping` calls were causing multiple TCP connections
    to open up on the cluster
    * Fixed an issue in which UCP install occasionally failed with the error
    "failed to change temp password"
    * Fixed an issue where multiple rapid updates of HTTP Routing Mesh configuration
    would not register correctly
    * Demoting a manager while in HA configuration no longer causes the `ucp-auth-api`
     container to provide errors

* UI/UX
    * When creating a user, pressing enter on keyboard no longer causes problems
    * Fixed assorted icon and text visibility glitches
    * Installing DTR no longer fails when "Enable scheduling on UCP controllers and
    DTR nodes" is unchecked.
    * Publishing a port to both TCP and UDP in a service via UI now works correctly

**Known issues**


The `docker stats` command is sometimes wrongly reporting high CPU usage.
Use the `top` command to confirm the real CPU usage of your node.
[Learn more](https://github.com/moby/moby/issues/28941).


**Version compatibility**

UCP 2.1 requires minimum versions of the following Docker components:

* Docker Engine 1.13.0
* Docker Remote API 1.25
* Compose 1.9
