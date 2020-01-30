---
title: Docker Trusted Registry overview
description: Learn how to install, configure, and use Docker Trusted Registry.
keywords: registry, repository, images
---

>{% include enterprise_label_shortform.md %}

Docker Trusted Registry (DTR) is Mirantis's enterprise-grade image storage solution. Installed behind the firewall, either on-premises or on a virtual private cloud, DTR provides a secure environment from which users can store and manage Docker images.

Specifically, DTR offers a wide range of benefits, including:

* [Image and Job Management](image-and-job-management)
* [Availability](availability)
* [Efficiency](efficiency)
* [Built-in Access Control](built-in-access-control)
* [Security Scanning](security-scanning)
* [Image Signing](image-signing)

## Image and Job Management

DTR offers a web user interface that allows authorized users to browse Docker images and [review repository events](/ee/dtr/user/audit-repository-events/). With this UI you can see which Dockerfile lines were used to produce an image, and if security scanning is enabled you can also view a list of all of the software installed in that image. In addition, the web UI can be used to [review and audit jobs](/ee/dtr/admin/manage-jobs/audit-jobs-via-ui/).

DTR can serve as a Continuous Integration and Continuous Delivery component, in the building, shipping, and running of applications.

## Availability

DTR is highly available through the use of multiple replicas of all containers
and metadata such that if a machine fails, DTR continues to operate and can be repaired.

## Efficiency

DTR has the ability to [cache images closer to users](admin/configure/deploy-caches/index.md)
to reduce the amount of bandwidth used when pulling Docker images.

DTR has the ability to [clean up unreferenced manifests and layers](admin/configure/garbage-collection.md).

## Built-in access control

DTR uses the same authentication mechanism as Docker Universal Control Plane.
Users can be managed manually or synchronized from LDAP or Active Directory. DTR
uses [Role Based Access Control](admin/manage-users/index.md) (RBAC) to allow you to implement fine-grained
access control policies for your Docker images.

## Security scanning

DTR has a built-in security scanner that can be used to discover what versions
of software are used in your images. It scans each layer and aggregates the
results to give you a complete picture of what you are shipping as a part of
your stack. Most importantly, it correlates this information with a
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
