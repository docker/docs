---
title: Docker Trusted Registry overview
description: Learn how to install, configure, and use Docker Trusted Registry.
keywords: registry, repository, images
---

>{% include enterprise_label_shortform.md %}

Docker Trusted Registry (DTR) is Mirantis's enterprise-grade image storage solution. Installed behind the firewall, either on-premises or on a virtual private cloud, DTR provides a secure environment from which users can store and manage Docker images.

Specifically, DTR offers a wide range of benefits that include:

* [Image and job management](image-and-job-management)
* [Availability](availability)
* [Efficiency](efficiency)
* [Built-in access control](built-in-access-control)
* [Security scanning](security-scanning)
* [Image signing](image-signing)

## Image and job management

 DTR has a web-based user interface that you ca use to browse images and [audit repository events](/ee/dtr/user/audit-repository-events/). With the UI, you can see which Dockerfile lines produced an image and, if security scanning is enabled, a list of all of the software installed in that image. You can also [audit jobs with the web interface](/ee/dtr/admin/manage-jobs/audit-jobs-via-ui/).

DTR can serve as a Continuous Integration and Continuous Delivery (CI/CD) component, in the building, shipping, and running of applications.

## Availability

DTR is highly available through the use of multiple replicas of all containers and metadata. As such, DTR will continue to operate in the event of machine failure, thus allowing for repair.

## Efficiency

DTR is able to reduce the bandwidth used when pulling Docker images by [caching images closer to users](admin/configure/deploy-caches/index.md). In addition, DTR can [clean up unreferenced manifests and layers](admin/configure/garbage-collection.md).

## Built-in access control

As with Universal Control Plane (UCP), DTR uses [Role Based Access Control (RBAC)](admin/manage-users/index.md), which allows you to manage image access, either manually, with LDAP, or with Active Directory.

## Security scanning

A security scanner is built into DTR, which can be used to discover the versions of the software that is in use in your images. This tool scans each layer and aggregates the results, offering a complete picture of what is being shipped as a part of your stack. Most importantly, as the security scanner is kept up-to-date by tapping into a [periodically updated](admin/configure/set-up-vulnerability-scans.md) vulnerability database, it is able to provide [unprecedented insight into your exposure to known security threats](user/manage-images/scan-images-for-vulnerabilities.md).

## Image signing

DTR ships with [Notary](/notary/getting_started.md), which allows you to sign and verify images using [Docker Content Trust](/engine/security/trust/content_trust.md). For more information on managing Notary data in DTR, refer to the [Using Notary to sign an image](user/manage-images/sign-images/index.md).

## Where to go next

* [DTR architecture](architecture.md)
* [Install DTR](admin/install/index.md)
