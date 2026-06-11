---
title: Set up your company for success with Docker
linkTitle: Admin set up
summary: Get the most out of Docker by streamlining workflows, standardizing development environments, and ensuring smooth deployments across your company.
description: Learn how to onboard your company and take advantage of all of the Docker products and features.
keywords: admin, onboarding, deployment, organization setup, docker business, rollout
aliases:
  - /guides/admin-set-up/comms-and-info-gathering/
  - /guides/admin-set-up/deploy/
  - /guides/admin-set-up/finalize-plans-and-setup/
  - /guides/admin-set-up/testing/
params:
  tags: [admin]
  time: 20 minutes
  image:
      url: "https://www.docker.com/pricing?ref=Docs&refAction=DocsGuidesAdminSetup"
---


Docker's tools provide a scalable, secure platform that empowers your
developers to create, ship, and run applications faster. As an administrator,
you can streamline workflows, standardize development environments, and ensure
smooth deployments across your organization.

By configuring Docker products to suit your company's needs, you can optimize
performance, simplify user management, and maintain control over resources.
This guide helps you set up and configure Docker products to maximize
productivity and success for your team while meeting compliance and security
policies.

## Who’s this for?

- Administrators responsible for managing Docker environments within their
  organization
- IT leaders looking to streamline development and deployment workflows
- Teams aiming to standardize application environments across multiple users
- Organizations seeking to optimize their use of Docker products for greater
  scalability and efficiency
- Organizations with a
  [Docker Business subscription](https://www.docker.com/pricing?ref=DocsGuides&refAction=DocsGuidesCTAClicked)

## What you’ll learn

- Why signing into your company's Docker organization provides access to usage
  data and enhanced functionality
- How to standardize Docker Desktop versions and settings to create a consistent
  baseline for all users, while allowing flexibility for advanced developers
- Strategies for implementing Docker's security configurations to meet company
  IT and software development security requirements without hindering developer productivity

## Features covered

This guide covers the following Docker features:

- [Organizations](/manuals/admin/organization/_index.md): The core structure
  for managing your Docker environment, grouping users, teams, and image
  repositories. Your organization was created with your subscription and is
  managed by one or more owners. Users signed into the organization are
  assigned seats based on the purchased subscription.
- [Enforce sign-in](/manuals/enterprise/security/enforce-sign-in/_index.md):
  By default, Docker Desktop doesn't require sign-in. You can configure
  settings to enforce this and ensure your developers sign in to your
  Docker organization.
- [SSO](/manuals/enterprise/security/single-sign-on/_index.md): Without SSO,
  user management in a Docker organization is manual. Setting
  up an SSO connection between your identity provider and Docker ensures
  compliance with your security policy and automates user provisioning. Adding
  SCIM further automates user provisioning and de-provisioning.
- General and security settings: Configuring key settings ensures smooth
  onboarding and usage of Docker products within your environment. You can also
  enable security features based on your company's specific security needs.

## Who needs to be involved

- Docker organization owner: Must be involved in the process and is required
  for several key steps
- DNS team: Needed during the SSO setup to verify the company domain
- MDM team: Responsible for distributing Docker-specific configuration files to
  developer machines
- Identity Provider team: Required for configuring the identity provider and
  establishing the SSO connection during setup
- Development lead: A development lead with knowledge of Docker configurations
  to help establish a baseline for developer settings
- IT team: An IT representative familiar with company desktop policies to
  assist with aligning Docker configuration to those policies
- Infosec: A security team member with knowledge of company development
  security policies to help configure security features
- Docker testers: A small group of developers to test the new settings and
  configurations before full deployment

## Tools integration

This guide covers integration with:

- Okta
- Entra ID SAML 2.0
- Azure Connect (OIDC)
- MDM solutions like Intune

## Communication and information gathering

### Communicate with your developers and IT teams

Before rolling out Docker Desktop across your organization, coordinate with key stakeholders to ensure a smooth transition.

#### Notify Docker Desktop users

You may already have Docker Desktop users within your company. Some steps in
this onboarding process may affect how they interact with the platform.

Communicate early with users to inform them that:

- They'll be upgraded to a supported version of Docker Desktop as part of the subscription onboarding
- Settings will be reviewed and optimized for productivity
- They'll need to sign in to the company's Docker organization using their
  business email to access subscription benefits

#### Engage with your MDM team

Device management solutions, such as Intune and Jamf, are commonly used for
software distribution across enterprises. These tools are typically managed by a dedicated MDM team.

Engage with this team early in the process to:

- Understand their requirements and lead time for deploying changes
- Coordinate the distribution of configuration files

Several setup steps in this guide require JSON files, registry keys, or .plist
files to be distributed to developer machines. Use MDM tools to deploy these configuration files and ensure their integrity.

### Identify Docker organizations

Some companies may have more than one
[Docker organization](/manuals/admin/organization/_index.md) created. These
organizations may have been created for specific purposes, or may not be
needed anymore.

If you suspect your company has multiple Docker organizations:

- Survey your teams to see if they have their own organizations
- Contact your Docker Support to get a list of organizations with users whose
  emails match your domain name

### Gather requirements

[Settings Management](/manuals/enterprise/security/hardened-desktop/settings-management/_index.md) lets you preset numerous configuration parameters for Docker Desktop.

Work with the following stakeholders to establish your company's baseline
configuration:

- Docker organization owner
- Development lead
- Information security representative

Review these areas together:

- Security features and
  [enforcing sign-in](/manuals/enterprise/security/enforce-sign-in/_index.md)
  for Docker Desktop users
- Additional Docker products included in your subscriptions

To view the parameters that can be preset, see [Configure Settings Management](/manuals/enterprise/security/hardened-desktop/settings-management/configure-json-file.md#step-two-configure-the-settings-you-want-to-lock-in).

### Optional: Meet with the Docker Implementation team

The Docker Implementation team can help you set up your organization,
configure SSO, enforce sign-in, and configure Docker Desktop.

To schedule a meeting, email successteam@docker.com.

## Finalize plans and begin setup

### Send finalized settings files to the MDM team

After reaching an agreement with the relevant teams about your baseline and
security configurations as outlined in the previous section, configure Settings Management using either the [Docker Admin Console](/manuals/enterprise/security/hardened-desktop/settings-management/configure-admin-console.md) or an
[`admin-settings.json` file](/manuals/enterprise/security/hardened-desktop/settings-management/configure-json-file.md).

Once the file is ready, collaborate with your MDM team to deploy your chosen
settings, along with your chosen method for [enforcing sign-in](/manuals/enterprise/security/enforce-sign-in/_index.md).

> [!IMPORTANT]
>
> Test this first with a small number of Docker Desktop developers to verify the functionality works as expected before deploying more widely.

### Manage your organizations

If you have more than one organization, consider either [consolidating them
into one organization](/manuals/admin/organization/setup/orgs.md) or creating a
[Docker company](/manuals/admin/company/_index.md) to manage multiple
organizations.

### Begin setup

#### Set up single sign-on and domain verification

Single sign-on (SSO) lets developers authenticate using their identity
providers (IdPs) to access Docker. SSO is available for a whole company and all associated organizations, or an individual organization that has a Docker
Business subscription. For more information, see the
[documentation](/manuals/enterprise/security/single-sign-on/_index.md).

You can also enable [SCIM](/manuals/enterprise/security/provisioning/scim/_index.md)
for further automation of provisioning and deprovisioning of users.

#### Set up Docker product entitlements included in the subscription

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

#### Ensure you're running a supported version of Docker Desktop

> [!WARNING]
>
> This step could affect the experience for users on older versions of Docker
> Desktop.

Existing users may be running outdated or unsupported versions of
Docker Desktop. All users should update to a supported version. Docker Desktop
versions released within the past 6 months from the latest release are supported.

Use an MDM solution to manage the version of Docker Desktop for users. Users
may also get Docker Desktop directly from Docker or through a company software
portal.

## Testing

### SSO and SCIM testing

Test SSO and SCIM by signing in to Docker Desktop or Docker Hub with the email
address linked to a Docker account that is part of the verified domain.
Developers who sign in using their Docker usernames remain unaffected by the
SSO and SCIM setup.

> [!IMPORTANT]
>
> Some users may need CLI based logins to Docker Hub, and for this they will
> need a [personal access token (PAT)](/manuals/security/access-tokens.md).

### Test Registry Access Management and Image Access Management

> [!WARNING]
>
> Communicate with your users before proceeding, as this step will impact all
> existing users signing into your Docker organization.

If you plan to use [Registry Access Management (RAM)](/manuals/enterprise/security/hardened-desktop/registry-access-management.md) and/or [Image Access Management (IAM)](/manuals/enterprise/security/hardened-desktop/image-access-management.md):

1. Ensure your test developer signs in to Docker Desktop using their
   organization credentials
2. Have them attempt to pull an unauthorized image or one from a disallowed
   registry via the Docker CLI
3. Verify they receive an error message indicating that the registry is
   restricted by the organization

### Deploy settings and enforce sign in to test group

Deploy the Docker settings and enforce sign-in for a small group of test users
via MDM. Have this group test their development workflows with containers on
Docker Desktop and Docker Hub to ensure all settings and the sign-in enforcement
function as expected.

### Test Docker Build Cloud capabilities

Have one of your Docker Desktop testers [connect to the cloud builder you created and use it to build](/manuals/build-cloud/usage.md).

### Test Testcontainers Cloud

Have a test developer [connect to Testcontainers Cloud](https://testcontainers.com/cloud/docs/#getting-started) and run a container in
the cloud to verify the setup is working correctly.

### Verify Docker Scout monitoring of repositories

Check the [Docker Scout dashboard](https://scout.docker.com/) to confirm that
data is being properly received for the repositories where Docker Scout has
been enabled.

### Verify access to Docker Hardened Images

Have a test developer attempt to [pull a Docker Hardened Image](/manuals/dhi/get-started.md) to confirm that
the team has proper access and can integrate these images into their workflows.

## Deploy your Docker setup

> [!WARNING]
>
> Communicate with your users before proceeding, and confirm that your IT and
> MDM teams are prepared to handle any unexpected issues, as these steps will
> affect all existing users signing into your Docker organization.

### Enforce SSO

Enforcing SSO means that anyone who has a Docker profile with an email address
that matches your verified domain must sign in using your SSO connection. Make
sure the Identity provider groups associated with your SSO connection cover all
the developer groups that you want to have access to the Docker subscription.

For instructions on how to enforce SSO, see [Enforce SSO](/manuals/enterprise/security/single-sign-on/connect.md).

### Deploy configuration settings and enforce sign-in to users

Have the MDM team deploy the configuration files for Docker to all users.

### Next steps

Congratulations, you've successfully completed the admin implementation process
for Docker.

To continue optimizing your Docker environment:

- Review your [organization's usage data](/manuals/admin/insights.md) to track adoption
- Monitor [Docker Scout findings](/manuals/scout/explore/analysis.md) for security insights
- Explore [additional security features](/manuals/enterprise/security/_index.md) to enhance your configuration
