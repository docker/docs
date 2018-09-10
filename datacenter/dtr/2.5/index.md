---
title: Docker Trusted Registry overview
description: Learn how to install, configure, and use Docker Trusted Registry.
keywords: registry, repository, images
---

Docker Trusted Registry (DTR) is the enterprise-grade image storage solution
from Docker. You install it behind your firewall so that you can securely store
and manage the Docker images you use in your applications.

## Image management

DTR can be installed on-premises, or on a virtual private
cloud. And with it, you can store your Docker images securely, behind your
firewall.

You can use DTR as part of your continuous integration, and continuous
delivery processes to build, ship, and run your applications.

DTR has a web based user interface that allows authorized users in your
organization to browse docker images. It provides information about
who pushed what image at what time. It even allows you to see what dockerfile
lines were used to produce the image and, if security scanning is enabled, to
see a list of all of the software installed in your images.

## Availability

DTR is highly available through the use of multiple replicas of all containers
and metadata such that if a machine fails, DTR continues to operate and can be repaired.

## Efficiency

DTR has the ability to [cache images closer to users](admin/configure/deploy-caches/index.md)
to reduce the amount of bandwidth used during docker pulls.

DTR has the ability to [clean up unreferenced manifests and layers](admin/configure/garbage-collection.md).

## Built-in access control

DTR uses the same authentication mechanism as Docker Universal Control Plane.
Users can be managed manually or synched from LDAP or Active Directory. DTR
uses [Role Based Access Control](admin/manage-users/index.md) (RBAC) to allow you to implement fine-grained
access control policies for who has access to your Docker images.

## Security scanning

DTR has a built in security scanner that can be used to discover what versions
of software are used in your images. It scans each layer and aggregates the
results to give you a complete picture of what you are shipping as a part of
your stack. Most importantly, it co-relates this information with a
vulnerability database that is kept up to date through [periodic
updates](admin/configure/set-up-vulnerability-scans.md). This
gives you [unprecedented insight into your exposure to known security
threats](user/manage-images/scan-images-for-vulnerabilities.md).

## Image signing

DTR ships with [Notary](/notary/getting_started.md)
built in so that you can use
[Docker Content Trust](/engine/security/trust/content_trust.md) to sign
and verify images. For more information about managing Notary data in DTR see
the [DTR-specific notary documentation](user/manage-images/sign-images/index.md).

## Where to go next

* [DTR architecture](architecture.md)
* [Install DTR](admin/install/index.md)
