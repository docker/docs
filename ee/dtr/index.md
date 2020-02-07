---
title: Docker Trusted Registry overview
description: Learn how to install, configure, and use Docker Trusted Registry.
keywords: registry, repository, images
---

>{% include enterprise_label_shortform.md %}

Docker Trusted Registry (DTR) is Mirantis's enterprise-grade image storage solution. Installed behind the firewall, either on-premises or on a virtual private cloud, DTR provides a secure environment from which users can store and manage Docker images.

Specifically, DTR offers a wide range of benefits that include:

* [Image and Job Management](image-and-job-management)
* [Availability](availability)
* [Efficiency](efficiency)
* [Built-in Access Control](built-in-access-control)
* [Security Scanning](security-scanning)
* [Image Signing](image-signing)

## Image and Job Management
 DTR offers a web-based user interface with which users can browse Docker images and [review repository events](/ee/dtr/user/audit-repository-events/). With the UI, in fact, users can even see which Dockerfile lines were used to produce an image and, if security scanning is enabled, a list of all of the software installed in that image. In addition, the web UI can be used to [review and audit jobs](/ee/dtr/admin/manage-jobs/audit-jobs-via-ui/).

DTR can serve as a Continuous Integration and Continuous Delivery component, in the building, shipping, and running of applications.

## Availability

DTR is highly available through the use of multiple replicas of all containers and metadata. As such, DTR will continue to operate in the event of machine failure, thus allowing for repair.

## Efficiency

DTR is able to reduce the amount of bandwidth used when pulling Docker images by [caching images closer to users](admin/configure/deploy-caches/index.md). in addition, DTR can [clean up unreferenced manifests and layers](admin/configure/garbage-collection.md).

## Built-in Access Control

DTR uses the [same authentication mechanism](https://docs.docker.com/ee/ucp/#built-in-security-and-access-control) as the Universal Control Plane, with which users can be managed manually or synchronized from LDAP or Active Directory. DTR employs [Role Based Access Control](admin/manage-users/index.md) (RBAC), which allow the implementation of fine-grained access control policies for Docker images.

## Security Scanning

A security scanner is built into DTR, which can be used to discover the versions of the software that is in use in your images. This tool scans each layer and aggregates the results, offering a complete picture of what is being shipped as a part of your stack. Most importantly, as the security scanner is kept up-to-date by tapping into a [periodically updated](admin/configure/set-up-vulnerability-scans.md) vulnerability database, it is able to provide [unprecedented insight into your exposure to known security threats](user/manage-images/scan-images-for-vulnerabilities.md).

## Image Signing

DTR ships with [Notary](/notary/getting_started.md) built in, which allows [Docker Content Trust](/engine/security/trust/content_trust.md) to be be put to use in image signing and verification. For more information about managing Notary data in DTR, refer to the [DTR-specific notary documentation](user/manage-images/sign-images/index.md).

## Where to go next

* [DTR architecture](architecture.md)
* [Install DTR](admin/install/index.md)
