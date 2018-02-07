---
description: Release notes for Docker EE for IBM Cloud. Learn more about the changes introduced in the latest versions.
keywords: ibm cloud, ibm, iaas, release notes
title: Docker EE for IBM Cloud (beta) release notes
---

Here you can learn about new features, bug fixes, breaking changes, and known issues for the latest Docker Enterprise Edition (EE) for IBM Cloud (beta) version.

## Version 1.0.2 (closed beta)

(26 January 2017)

Start using version 1.0.2 of Docker EE for IBM Cloud today!

1. Update the CLI plug-in:

   ```bash
   $ bx plugin update docker-for-ibm-cloud -r Bluemix
   ```

2. [Deploy a new cluster](administer-swarmd.md#create-swarms).

The second release of the closed beta includes the following enhancements:

* [IBM Cloud Security Groups](https://console.bluemix.net/docs/infrastructure/security-groups/sg_overview.html#about-security-groups) securely control access to cluster traffic.
* Users get their UCP password in the output of the [cluster create process](administer-swarmd.md#create-swarms).
* By default, IBM Cloud Swift API Object Storage is used for the DTR container instead of local storage volume to improve high availability.
* The [previous known issue](#service-load-balancer-has-the-configuration-of-an-older-service) about the service load balancer having the configuration of an older service is resolved.

## Version 1.0.1 (closed beta)

(20 December 2017)

Docker Enterprise Edition for IBM Cloud is a Container-as-a-Service platform that helps you modernize and extend your applications. It provides an unmanaged, native Docker environment running within IBM Cloud, giving you the ability to enhance your apps with services from the IBM Cloud catalog such as Watson, Internet of Things, Data, Logging, Monitoring, and many more.

The beta is based on the latest Docker EE version 17.06. You receive a 90-day Docker EE license for 20 Linux x86-64 nodes that you use when creating a cluster with the IBM Cloud `bx d4ic` CLI plug-in.

[Sign up for the closed beta](mailto:sealbou@us.ibm.com). Then use the [quick start](quickstart.md) to spin up your first Docker EE for IBM Cloud swarm.

### Known issues

#### Service load balancer has the configuration of an older service

**What it is**: You update (`docker service update`) or a create a new service after removing an old one (`docker service rm` then immediately `docker service create`) with a change to [the certificate or health check path](load-balancer.md#labels-for-ssl-termination-and-health-check-paths) values. The service load balancer still has the configuration of the older service.

**Why it's happening**: The InfraKit ingress controller does not check if the certificate or health check path has changed between the service load balancer configuration and the Docker swarm listener configuration. This results in the old configuration not getting updated to the new values.

**What to do about it**: Avoid updating a service with the `docker service update` command.

First remove the old service with `docker service rm`, and then wait for a time so that the service load balancer configuration updates. Then, create the new service with `docker service create`.

If you must use `docker service update`, remove the configuration for the port of the service load balancer in the [IBM Cloud infrastructure web UI](https://control.softlayer.com/). The InfraKit ingress controller then recreates the configuration with the new values.
