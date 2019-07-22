---
title: Docker Enterprise release notes
description: Learn about the new features, bug fixes, and breaking changes for Docker Enterprise.
keywords: engine enterprise, ucp, dtr, desktop enterprise, whats new, release notes
---

This page provides information about Docker Enterprise 3.0. For 
detailed information about for each enterprise component, refer to the individual component release notes 
pages listed in the following **Docker Enterprise components install and upgrade** section.

## What’s New?

| Feature | Component | Component version |
|---------|-----------|-------------------|
| [Group Managed Service Accounts (gMSA)](/engine/swarm/services.md#gmsa-for-swarm) | UCP | 3.2.0 |
| [Open Security Controls Assessment Language (OSCAL)](/compliance/oscal/) | UCP | 3.2.0 |
| [Container storage interface (CSI)](/ee/ucp/kubernetes/storage/use-csi/) | UCP | 3.2.0 |
| [Internet Small Computer System Interface (iSCSI)](/ee/ucp/kubernetes/storage/use-iscsi/) | UCP | 3.2.0 |
| [System for Cross-domain Identity Management (SCIM)](/ee/ucp/admin/configure/integrate-scim/) | UCP | 3.2.0 |
| [Pod Security Policies](/ee/ucp/kubernetes/pod-security-policies/) | UCP | 3.2.0 |
| [Docker Registry CLI (Experimental)](/engine/reference/commandline/registry/) | DTR | 2.7.0 |
| [App Distribution](/ee/dtr/user/manage-applications/) | DTR | 2.7.0 |
| [Client certificate-based Authentication](/ee/enable-client-certificate-authentication/) | DTR and UCP|2.7.0 (DTR) and 3.2.0 (UCP)|
| [Application Designer](/ee/desktop/app-designer/) | Docker Desktop Enterprise | 0.1.4  |
| [Docker App (Experimental)](/app/working-with-app/) |CLI | 0.8.0 |
| [Docker Assemble (Experimental)](/assemble/install/) | CLI | 0.36.0 |
| [Docker Buildx (Experimental)](/buildx/working-with-buildx/)| CLI | 0.2.2 |
| [Docker Cluster](/cluster/) | CLI | 1.0.0 |
| [Docker Template CLI (Experimental)](/app-template/working-with-template/) | CLI | 0.1.4 |


## Docker Enterprise components install and upgrade

| Component Release Notes | Version | Install | Upgrade |
|---------|-----------|-------------------|-------------- |
| [Docker Engine - Enterprise](/engine/release-notes/) | 19.03 | [Install Docker Engine - Enterprise](/ee/supported-platforms/) | [Upgrade Docker Engine - Enterprise](/ee/upgrade/) |
| [Docker UCP](/ee/ucp/release-notes/) | 3.2 | [Install UCP](/ee/ucp/admin/install/) | [Upgrade UCP](/ee/ucp/admin/install/upgrade/) |
| [DTR](/ee/dtr/release-notes/) | 2.7 | [Install DTR](/ee/dtr/admin/install/) | [Upgrade DTR](/ee/dtr/admin/upgrade/) |
| [Docker Desktop Enterprise](/ee/desktop/release-notes/) | 2.1.0 |Install Docker Desktop Enterprise [Mac](/ee/desktop/admin/install/mac/), [Windows](/ee/desktop/admin/install/windows/) | Upgrade Docker Desktop Enterprise  [Mac](/ee/desktop/admin/install/mac/), [Windows](/ee/desktop/admin/install/windows/) |

Refer to the [Compatibility Matrix](https://success.docker.com/article/compatibility-matrix) and the [Maintenance Lifecycle](https://success.docker.com/article/maintenance-lifecycle) for compatibility and software maintenance details.

 
 

