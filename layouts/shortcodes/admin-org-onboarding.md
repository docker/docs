{{ $sso_link := "[Configure SSO](/security/for-admins/single-sign-on/)" }}
{{ $scim_link := "[Configure SCIM](/security/for-admins/scim/)" }}
{{ $members_link := "[Invite members](/docker-hub/members/)" }}
{{ $audit_link := "[Audit your domains](/docker-hub/domain-audit/)" }}

{{ if eq (.Get "product") "admin" }}
  {{ $sso_link = "[Configure SSO](/security/for-admins/single-sign-on/)" }}
  {{ $scim_link = "[Configure SCIM](/security/for-admins/scim/)" }}
  {{ $members_link = "[Invite members](/admin/organization/members/)" }}
  {{ $audit_link = "[Audit your domains](/admin/organization/security-settings/domains/)" }}
{{ end }}

Onboarding your organization allows you to gain visibility into the activity of your users and enforce security settings. In addition, members of your organization receive increased pull limits and other organization wide benefits. For more details, see [Docker subscriptions and features](/subscription/details/).

## Prerequisites

Before you start to onboard your organization, ensure that you:
- Have a Docker Team or Business subcription. See [Pricing & Subscriptions](https://www.docker.com/pricing/) for details.

  > **Note**
  >
  > When purchasing a subscription through [Pricing & Subscriptions](https://www.docker.com/pricing/), the on-screen instructions guide you through creating an organization. If you have purchased a subscription through Docker Sales and you have not yet created an organization, see [Create an organization](/docker-hub/orgs#create-an-organization).

- Familiarize yourself with Docker concepts and terminology in the [glossary](/glossary/) and [FAQs](/docker-hub/general-faqs/).

## Step 1: Identify your Docker users and their Docker accounts

Identifying your users will ensure that you allocate your subscription seats efficiently and that all your Docker users receive the benefits of your subscription.

1. Identify the Docker users in your organization.
   - If your organization uses device management software, like MDM or JAMF, you may use the device management software to help identify Docker users. See your device management software's documentation for details. You can identify Docker users by checking if Docker Desktop is installed at the following location on each user's machine:
      - Mac: `/Applications/Docker.app`
      - Windows: `C:\Program Files\Docker\Docker`
      - Linux: `/opt/docker-desktop`
   - If your organization doesn't use device management software or your users haven't installed Docker Desktop yet, you may survey your users.
2. Instruct all your Docker users in your organization to update their existing Docker account's email address to an address that's in your organization's domain, or to create a new account using an email address in your organization's domain.
   - To update an account's email address, instruct your users to sign in to [Docker Hub](https://hub.docker.com), and update the email address to their email address in your organization's domain.
   - To create a new account, instruct your users to go [sign up](https://hub.docker.com/signup) using their email address in your organization's domain.
3. Ask your Docker sales representative or [contact sales](https://www.docker.com/pricing/contact-sales/) to get a list of Docker accounts that use an email address in your organization's domain.

## Step 2: Invite owners

When you create an organization, you are the only owner. You may optionally add additional owners. Owners can help you onboard and manage your organization.

To add an owner, invite a user and assign them the owner role. For more details, see {{ $members_link }}.

## Step 3: Invite members

When you add users to your organization, you gain visibility into their activity and you can enforce security settings. In addition, members of your organization receive increased pull limits and other organization wide benefits.

To add a member, invite a user and assign them the member role. For more details, see {{ $members_link }}.

## Step 4: Manage members with SSO and SCIM

Configuring Single Sign-On (SSO) or System for Cross-domain Identity Management (SCIM) is optional and only available to Docker Business subscribers. To upgrade a Docker Team subscription to a Docker Business subscription, see [Upgrade your subscription](/subscription/upgrade/).

You can manage your members in your identity provider and automatically provision them to your Docker organization with SSO or SCIM. See the following for more details.
   - {{ $sso_link }} to automatically provision and add members when they sign in to Docker through your identity provider.
   - {{ $scim_link }} to automatically provision, add, and de-provision members to Docker through your identity provider.


## Step 5: Enforce sign-in for Docker Desktop

By default, members of your organization can use Docker Desktop on their machines without signing in to any Docker account. You must enforce sign-in to ensure that users receive the benefits of your Docker subscription and that security settings are enforced. For details, see [Configure registry.json to enforce sign-in](/docker-hub/configure-sign-in/).

## What's next

- [Create](/docker-hub/repos/create/) and [manage](/docker-hub/repos/) repositories.
- Create [teams](/docker-hub/manage-a-team/) for fine-grained repository access.
- Configure [Hardened Docker Desktop](/desktop/hardened-desktop/) to improve your organization’s security posture for containerized development.
- {{ $audit_link }} to ensure that all Docker users in your domain are part of your organization.

Your Docker subscription provides many more additional features. To learn more, see [Docker subscriptions and features](/subscription/details/).
