---
title: Finalize plans and begin setup
description: Collaborate with your MDM team to distribute configurations and set up SSO and Docker product trials.
weight: 20
---

## Step one: Send finalized settings files to the MDM team 

After reaching an agreement with the relevant teams about your baseline and security configurations as outlined in module one, configure Settings Management using either the [Docker Admin Console](/manuals/security/for-admins/hardened-desktop/settings-management/configure-admin-console.md) or an [`admin-settings.json` file](/manuals/security/for-admins/hardened-desktop/settings-management/configure-json-file.md).

Once the file is ready, collaborate with your MDM team to deploy your chosen settings, along with your chosen method for [enforcing sign-in](/manuals/security/for-admins/enforce-sign-in/_index.md).

> [!IMPORTANT]
>
> It’s highly recommended that you test this first with a small number of Docker Desktop developers to verify the functionality works as expected before deploying more widely.

## Step two: Manage your organizations

If you have more than one organization, it’s recommended that you either consolidate them into one organization or create a [Docker company](/manuals/admin/company/_index.md) to manage multiple organizations. Work with the Docker Customer Success and Implementation teams to make this happen.

## Step three: Begin setup

### Set up single sign-on SSO domain verification

Single sign-on (SSO) lets developers authenticate using their identity providers (IdPs) to access Docker. SSO is available for a whole company, and all associated organizations, or an individual organization that has a Docker Business subscription. For more information, see the [documentation](/manuals/security/for-admins/single-sign-on/_index.md).

You can also enable [SCIM](/manuals/security/for-admins/provisioning/scim.md) for further automation of provisioning and deprovisioning of users.

### Set up Docker product entitlements included in the subscription

[Docker Build Cloud](/manuals/build-cloud/_index.md) significantly reduces build times, both locally and in CI, by providing a dedicated remote builder and shared cache. Powered by the cloud, developer time and local resources are freed up so your team can focus on more important things, like innovation. To get started, [set up a cloud builder](https://app.docker.com/build/). 

[Docker Scout](manuals/scout/_index.md) is a solution for proactively enhancing your software supply chain security. By analyzing your images, Docker Scout compiles an inventory of components, also known as a Software Bill of Materials (SBOM). The SBOM is matched against a continuously updated vulnerability database to pinpoint security weaknesses. To get started, see [Quickstart](/manuals/scout/quickstart.md).

### Ensure you're running a supported version of Docker Desktop

> [!WARNING]
>
> This step could affect the experience for users on older versions of Docker Desktop.  

Existing users may be running outdated or unsupported versions of Docker Desktop. It is highly recommended that all users update to a supported version. Docker Desktop versions released within the past 6 months from the latest release are supported.

It's recommended that you use a MDM solution to manage the version of Docker Desktop for users. Users may also get Docker Desktop directly from Docker or through a company software portal.  
