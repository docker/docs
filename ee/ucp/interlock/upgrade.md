---
title: Layer 7 routing upgrade
description: Learn how to route layer 7 traffic to your swarm services
keywords: routing, proxy, hrm
---

The [HTTP routing mesh](/datacenter/ucp/2.2/guides/admin/configure/use-domain-names-to-access-services.md)
functionality was redesigned in UCP 3.0 for greater security and flexibility.
The functionality was also renamed to "layer 7 routing", to make it easier for
new users to get started.

[Learn about the new layer 7 routing functionality](index.md).

To route traffic to your service you apply specific labels to your swarm
services, describing the hostname for the service and other configurations.
Things work in the same way as they did with the HTTP routing mesh, with the
only difference being that you use different labels.

You don't have to manually update your services. During the upgrade process to
3.0, UCP updates the services to start using new labels.

This article describes the upgrade process for the routing component, so that
you can troubleshoot UCP and your services, in case something goes wrong with
the upgrade.

# UCP upgrade process

If you are using the HTTP routing mesh, and start an upgrade to UCP 3.0:

1. UCP starts a reconciliation process to ensure all internal components are
deployed. As part of this, services using HRM labels are inspected.
2. UCP creates the `com.docker.ucp.interlock.conf-<id>` based on HRM configurations.
3. The HRM service is removed.
4. The `ucp-interlock` service is deployed with the configuration created.
5. The `ucp-interlock` service deploys the `ucp-interlock-extension` and
`ucp-interlock-proxy-services`.

The only way to rollback from an upgrade is by restoring from a backup taken
before the upgrade. If something goes wrong during the upgrade process, you
need to troubleshoot the interlock services and your services, since the HRM
service won't be running after the upgrade.

[Learn more about the interlock services and architecture](architecture.md).

## Check that routing works

After upgrading to UCP 3.0, you should check if all swarm services are still
routable.

For services using HTTP:

```bash
curl -vs http://<ucp-url>:<hrm-http-port>/ -H "Host: <service-hostname>"
```

For services using HTTPS:

```bash
curl -vs http://<ucp-url>:<hrm-https-port>
```

After the upgrade, check that you can still use the same hostnames to access
the swarm services.

## The ucp-interlock services are not running

After the upgrade to UCP 3.0, the following services should be running:

* `ucp-interlock`: monitors swarm workloads configured to use layer 7 routing.
* `ucp-interlock-extension`: Helper service that generates the configuration for
the `ucp-interlock-proxy` service.
* `ucp-interlock-proxy`: A service that provides load balancing and proxying for
swarm workloads.

To check if these services are running, use a client bundle with administrator
permissions and run:

```bash
docker ps --filter "name=ucp-interlock"
```

* If the `ucp-interlock` service doesn't exist or is not running, something went
wrong with the reconciliation step.
* If this still doesn't work, it's possible that UCP is having problems creating
the `com.docker.ucp.interlock.conf-1`, due to name conflicts. Make sure you
don't have any configuration with the same name by running:
   ```
   docker config ls --filter "name=com.docker.ucp.interlock"
   ```
* If either the `ucp-interlock-extension` or `ucp-interlock-proxy` services are
not running, it's possible that there are port conflicts.
As a workaround re-enable the layer 7 routing configuration from the
[UCP settings page](deploy/index.md). Make sure the ports you choose are not
being used by other services.

## Workarounds and clean-up

If you have any of the problems above, disable and enable the layer 7 routing
setting on the [UCP settings page](deploy/index.md). This redeploys the
services with their default configuration.

When doing that make sure you specify the same ports you were using for HRM,
and that no other services are listening on those ports.

You should also check if the `ucp-hrm` service is running. If it is, you should
stop it since it can conflict with the `ucp-interlock-proxy` service.

## Optionally remove labels

As part of the upgrade process UCP adds the
[labels specific to the new layer 7 routing solution](usage/labels-reference.md).

You can update your services to remove the old HRM labels, since they won't be
used anymore.

## Optionally segregate control traffic

Interlock is designed so that all the control traffic is kept separate from
the application traffic.

If before upgrading you had all your applications attached to the `ucp-hrm`
network, after upgrading you can update your services to start using a
dedicated network for routing that's not shared with other services.
[Learn how to use a dedicated network](usage/index.md).

If before upgrading you had a dedicate network to route traffic to each service,
Interlock will continue using those dedicated networks. However the
`ucp-interlock` will be attached to each of those networks. You can update
the `ucp-interlock` service so that it is only connected to the `ucp-hrm` network.
