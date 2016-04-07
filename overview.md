<!--[metadata]>
+++
aliases = ["/docker-hub-enterprise/"]
title = "Docker Trusted Registry overview"
description = "Learn how to install, configure, and use Docker Trusted Registry."
keywords = ["docker, registry, repository, images"]
[menu.main]
parent="workw_dtr"
weight=0
+++
<![end-metadata]-->

# Docker Trusted Registry overview

Docker Trusted Registry (DTR) is the enterprise-grade image storage solution
from Docker. You install it behind your firewall so that you can securely store
and manage the Docker images you use in your applications.

<!--  TODO: add screenshot -->

## Image management

Docker UCP can be installed on-premises, or on a virtual private cloud.
And with it, you can store your Docker images securely, behind your firewall.

You can use DTR as part of your Continuous Integration (CI), and Continuous
Delivery (CD) processes, to build, run, and ship your applications.


## Built-in security and access control

Docker UCP has its own built-in authentication mechanism, and supports LDAP
and Active Directory. It also supports Role Based Access Control (RBAC).

This allows you to implement fine-grain access control policies, on who has
access to your Docker images.

<!--  TODO: add screenshot -->
