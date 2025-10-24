---
title: Finalize plans and begin setup
description: Collaborate with your MDM team to distribute configurations and set up SSO and Docker product trials.
weight: 20
---

## Send finalized settings files to the MDM team

After reaching an agreement with the relevant teams about your baseline and
security configurations as outlined in the previous section, configure Settings Management using either the [Docker Admin Console](/manuals/enterprise/security/hardened-desktop/settings-management/configure-admin-console.md) or an
[`admin-settings.json` file](/manuals/enterprise/security/hardened-desktop/settings-management/configure-json-file.md).

Once the file is ready, collaborate with your MDM team to deploy your chosen
settings, along with your chosen method for [enforcing sign-in](/manuals/enterprise/security/enforce-sign-in/_index.md).

> [!IMPORTANT]
>
> Test this first with a small number of Docker Desktop developers to verify the functionality works as expected before deploying more widely.

## Manage your organizations

If you have more than one organization, consider either [consolidating them
into one organization](/manuals/admin/organization/orgs.md) or creating a
[Docker company](/manuals/admin/company/_index.md) to manage multiple
organizations.

## Begin setup

### Set up single sign-on and domain verification

Single sign-on (SSO) lets developers authenticate using their identity
providers (IdPs) to access Docker. SSO is available for a whole company and all associated organizations, or an individual organization that has a Docker
Business subscription. For more information, see the
[documentation](/manuals/enterprise/security/single-sign-on/_index.md).

You can also enable [SCIM](/manuals/enterprise/security/provisioning/scim.md)
for further automation of provisioning and deprovisioning of users.

### Set up Docker product entitlements included in the subscription

[Docker Build Cloud](/manuals/build-cloud/_index.md) significantly reduces
build times, both locally and in CI, by providing a dedicated remote builder
and shared cache. Powered by the cloud, developer time and local resources are
freed up so your team can focus on more important things, like innovation.
To get started, [set up a cloud builder](https://app.docker.com/build/).

[Docker Scout](manuals/scout/_index.md) is a solution for proactively enhancing
your software supply chain security. By analyzing your images, Docker Scout
compiles an inventory of components, also known as a Software Bill of Materials
(SBOM). The SBOM is matched against a continuously updated vulnerability
database to pinpoint security weaknesses. To get started, see
[Quickstart](/manuals/scout/quickstart.md).

[Testcontainers Cloud](https://testcontainers.com/cloud/docs/) allows
developers to run containers in the cloud, removing the need to run heavy
containers on your local machine.

[Docker Hardened Images](/manuals/dhi/_index.md) are minimal, secure, and production-ready container base and application images maintained by Docker.
Designed to reduce vulnerabilities and simplify compliance, DHIs integrate
easily into your existing Docker-based workflows with little to no retooling
required.

### Ensure you're running a supported version of Docker Desktop

> [!WARNING]
>
> This step could affect the experience for users on older versions of Docker
Desktop.

Existing users may be running outdated or unsupported versions of
Docker Desktop. All users should update to a supported version. Docker Desktop
versions released within the past 6 months from the latest release are supported.

Use an MDM solution to manage the version of Docker Desktop for users. Users
may also get Docker Desktop directly from Docker or through a company software
portal.
