<!--[metadata]>
+++
aliases = ["/docker-trusted-registry/prior-release-notes/"]
title = "Prior Trusted Registry release notes"
description = "Archived release notes for Docker Trusted Registry"
keywords = ["docker, documentation, about, technology, understanding, enterprise, hub, registry, Docker Trusted Registry, release"]
[menu.main]
parent="dtr_menu_release_notes"
weight=30
+++
<![end-metadata]-->

# Docker Trusted Registry release notes archive

This document contains the previous versions of the Docker Trusted Registry
release notes.

## Docker Trusted Registry 1.3.3
(18 September 2015) (amended: 2 November 2015)

This release corrects the following issues in Docker Trusted Registry 1.3.2

* Fixed an issue related to LDAP integration for users of Oracle Virtual Directory.

* Corrected an issue where Docker Trusted Registry would not accept a given certificate if the configured domain was only in the Subject Alternative Names
(SANs) field and not in the Common Name (CN) field of the certificate.

* Docker, Inc. discovered an issue in which the tokens used in authorization caused a break in certain deployments that utilized a load balancer in front of
multiple Trusted Registry instances to achieve high availability. We regret any
inconvenience this may have caused you and is working on a future fix.

## Docker Trusted Registry 1.3.2
(16 September 2015)

This release addresses the following change in Docker Trusted Registry 1.3.2 and is only available to customers who purchased DTR through Amazon Web Services (AWS) Marketplace.

* Docker Trusted Registry (DTR) now supports Amazon Web
Services (AWS) Integrated Billing. Previously, AWS users were required to
separately purchase a DTR license from Docker. AWS users can try DTR
out-of-the-box.

## Docker Trusted Registry 1.3.1
(31 August 2015)

This release corrects the following issues in Docker Trusted Registry 1.3.0

* The dashboard page was calculating incorrect stats.
* LDAP group sync failed to handle paginated results for extremely large groups.
* The repo delete endpoint returned incorrect error codes under certain conditions.

## Docker Trusted Registry 1.3.0
(26 August 2015)

This release addresses a few bugs and issues in Docker Trusted Registry 1.2.0 and introduces some new features and functionality, including:

* A completely new user-interface for the Admin application brings Docker Trusted Registry in line with other Docker products and provides greater ease-of-use.

* A new Accounts & Repos API provides new fine-grained role-based access control down to the per-repo level. See the [API's documentation](http://docs.docker.com/apidocs/v1.3.3/) for more information.

* Improvements to the handling of configuration changes so that fewer restarts are required.

* Multiple security improvements and bug fixes.

## Docker Trusted Registry 1.2.0
(23 July 2015)

This release adds CentOS support and addresses a few bugs and issues in Docker Trusted Registry 1.1.0:

* Fixes an issue where for certain configurations of Docker Trusted Registry, proxy configuration settings and variables were not being passed to all Docker Trusted Registry containers and thus were not being respected.
* Documentation links in the UI now point to correct docs.
* Generated support info bundles have been scrubbed to remove highly sensitive data.
* Certifies support for CentOS 7.1.

## Docker Trusted Registry 1.1.0
(23 June 2015)

This release of Docker Trusted Registry (formerly DHE) adds major integration
with the AWS and Azure marketplaces, giving customers a smoother installation
path. Docker Trusted Registry 1.1 also adds finer-grained permissions and
improvements and additions to the UI and logging. Bugs in LDAP/AD integration
have also been remediated, improving the stability and usability of Docker
Trusted Registry. See below for specifics.

### New Features

* New, more granular, [roles for users](../user-management/permission-levels.md). Docker Trusted Registry users can now be assigned different levels of access
(admin, r/w, r/o) to the repositories. **Important:** Existing Docker Trusted
Registry users should make sure to see the note [below](#dhe-1-0-upgrade-warning) regarding migrating users before upgrading.
* A new storage status indicator for storage space. The dashboard now shows used and available storage space for supported storage drivers.
* A new diagnostics tool gathers and bundles Docker Trusted Registry logs, system information, container
information, and other configuration settings for use by Docker support or as a
backup.
* Performance and reliability improvements to the S3 storage backend.
* Docker Trusted Registry images are now available on the Amazon AWS and Microsoft Azure marketplaces.

### Fixes

The following notable issues have been remediated:

* Fixed an issue that caused Docker Trusted Registry logins to fail if some LDAP servers were unreachable.
* Fixed a resource leak in Docker Trusted Registry storage.

### DHE 1.0 Upgrade Warning

Customers who are currently using DHE 1.0 **must** follow the [upgrading instructions](https://forums.docker.com/t/upgrading-docker-hub-enterprise-to-docker-trusted-registry/1925) in our support Knowledge Base. These instructions will show you how to modify existing authentication data and storage volume
settings to move to Docker Trusted Registry. Note that automatic upgrading has
been disabled for DHE users because of these issues.

## Docker Trusted Registry 1.0.1
(11 May 2015)

- Addresses compatibility issue with 1.6.1 CS Docker Engine

## Docker Trusted Registry 1.0.0
(23 Apr 2015)

- First release
