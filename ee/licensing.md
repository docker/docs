---
title: Docker Enterprise licensing
description: Learn about the Docker Enterprise licensing
keywords: license, enterprise, engine enterprise, ucp, dtr, desktop enterprise
toc_min: 1
toc_max: 2
---

>{% include enterprise_label_shortform.md %}

This page provides information about Docker Enterprise 3.0 licensing. Docker Enterprise 3.0 is a soft bundle of products that deliver the complete desktop-to-cloud workflow.

> Important
>
> Docker Enterprise 3.0 consists of two separately purchased and licensed products:
> - **Universal Control Plane (UCP) and Docker Trusted Registry (DTR) with Docker Engine - Enterprise**: Installed on servers and licensed based on the size of the environment.
>
> - **Docker Desktop Enterprise**: Installed on developer workstations and separately licensed addition to the platform.

# Licensing Docker Enterprise

Complete the following steps to license Docker Enterprise 3.0:

1. Set up your Docker Hub ID, Organizations, and Team.
1. Register your Docker license key.
1. Activate your subscription.
1. Access and download your license.
1. Apply the license.

## Set up Docker Hub IDs, Organization, and Teams

1. Before you begin, identify who needs access to the Docker license files. This may be any number of users across your company, agency, or organization.
2. Create a Docker Hub ID for each person that requires access to the license file. Visit [Docker Hub](https://hub.docker.com/) to create a Docker ID. This requires email verification. Ensure that you provide a corporate email address.
3. Decide on the name of the Hub Organization that will be used to access licenses. The designated team leader must log into Docker Hub and create the Hub Organization. For more information about creating organizations, see [Teams and Organizations](https://docs.docker.com/docker-hub/orgs/).

    > The name of the Hub Organization must be unique across Docker Hub. If you have multiple independent organizations within your company, for example, if you represent 'OrgB' within 'CompanyA', use an Organization name such as 'CompanyAOrgB' instead of 'CompanyA'. Ensure that the company's licenses are registered to the organization, rather than an individual. Otherwise, the team will not be able to access the license(s).

4. On the **Organizations** page, select your organization, and then click on the **Teams** tab. You should see the name of your organization, and the **Choose Team** option below that. You should also see that you have a special team called **owners**. The owners team has full access to all repositories within the organization.
5. Click **Add User** and add the Docker ID of every member that needs access to the license files. After entering each Docker ID, ensure the IDs are displayed in the **Members** section within the **owners** team.

## Register your Docker license key

When your license order has been processed by Docker, your support administrator, designated support contact, billing contact, or the primary sales point of contact receives a welcome email. This email contains a link to an onboarding experience and your license keys.

It is important to note that only one recipient, preferably the support administrator or technical primary point of contact, must click on the links and follow the instructions in the welcome email.

Depending on which type of email your team receives, there may be several intermediate steps before you are directed towards activating your subscription.

## Activate your subscription

1. Log into the [Docker Hub](https://hub.docker.com/procurement) using your Docker ID and enter the activation key provided in your welcome email.
2. Click the **Submit Access Key** button.
3. From the **Subscribe as** drop-down menu, select the ID you would like the subscription applied to.
4. Click the check box to agree to the terms of service and then click **Confirm and Subscribe**.
5. Your subscription is now available. You will then be redirected to the **My Content** page.
6. On the **My Content** page, click the **Setup** button on your Docker Enterprise subscription to access your license keys.

## Access and download your license

When you click the **Setup** button in step 6 above, you will be redirected to the setup instructions page where you can click the License Key link and download your Enterprise license. This page contains links and instructions on obtaining software such as UCP, DTR, Desktop Enterprise, and Engine - Enterprise.

## Apply the license

After you’ve downloaded the license keys, you can apply it to your Docker Enterprise component. The following section contains information on how to apply the license to UCP, DTR, Enterprise - Engine, and Docker Desktop Enterprise components.

### UCP

1. Log into the UCP web UI using the administrator credentials.
2. Navigate to the **Admin Settings** page.
3. On the left pane, click **License** and then **Upload License**. The license is refreshed immediately.

For details, see [Licensing UCP](ucp/admin/configure/license-your-installation.md).

### DTR

1. Navigate to {% raw %} https://<dtr-url>{% endraw %} and log in with your credentials.
2. Select **System** from the left navigation pane.
3. Click **Apply new license** and upload your license key.
    
For details, see [Licensing DTR](dtr/admin/configure/license-your-installation.md).

### Engine - Enterprise

When you license UCP, the same license is applied to the underlying engines in the cluster. Docker recommends that Enterprise customers use UCP to manage their license.

### Desktop Enterprise

> Docker Desktop Enterprise licenses are not included as part of your UCP, DTR, and Engine - Enterprise license. It is a separate license installed on developer workstations. Please contact your Sales team to obtain [Docker Desktop Enterprise](https://docs.docker.com/desktop/enterprise/) licenses.

Install the Docker Desktop Enterprise license file at the following location:

On macOS:

`/Library/Group Containers/group.com.docker/docker_subscription.lic`

On Windows:

`%ProgramData%\DockerDesktop\docker_subscription.lic`

You must create the path if it doesn’t already exist. If the license file is missing, you will be asked to provide it when you try to run Docker Desktop Enterprise. Contact your system administrator to obtain the license file.

# What happens when my Enterprise license expires?

If there is a lapse in your Docker Enterprise entitlements, you will be alerted in the product until a new license is applied. However, you will not lose access to the software.

**Engine**: Docker Engine doesn't depend on the license being installed for ongoing functionality. It only requires licensing to access the package repositories. Note that an expired license may affect the node's ability to upgrade.

**UCP**: UCP components continue to work as expected when the license expires. However, warning banners regarding the license expiry will appear in the UCP web UI.

**DTR**: Image pushes to DTR will be disabled once the license expires. All other functionality will persist.

**Desktop**: Warnings regarding the license expiry appear in the Desktop UI. Note that an expired license may affect the software's ability to upgrade.

Please work with your sales team to ensure that your licenses are renewed before the expiration date on your licenses.

# Where to go next?

Refer to the following articles for more information:

- [Get Ready for your Licenses and Access to Technical Support](https://success.docker.com/article/get-ready-for-licenses-and-support)
- [Commercial Support Service Levels](https://success.docker.com/article/commercial-support-service-levels)
- [Where is my Docker Enterprise license?](https://success.docker.com/article/where-is-my-docker-enterprise-edition-license)
