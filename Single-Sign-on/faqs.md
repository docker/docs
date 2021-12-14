---
description: Single Sign-on FAQs
keywords: Docker, Docker Hub, SSO FAQs, single sign-on
title: Single Sign-on FAQs
toc_max: 2
---

## General

### Q: How does Docker SSO work?

Docker Single Sign-on (SSO) allows users to authenticate using their identity providers (IdPs) to access Docker. Docker currently supports Azure AD and identity providers that support SAML 2.0. When you enable SSO, users are redirected to your provider’s authentication page to authenticate using their email and password.

### Q: What SSO flows are supported by Docker?

Docker currently supports Service Provider Initiated (SP-initiated) SSO flow. This means, users must sign in to Docker Hub or Docker Desktop to initiate the SSO authentication process.

### Q: Can I enable SSO in all organizations?

You can enable SSO on organizations that are part of the Docker Business subscription.

### Q: We currently have a Docker Team subscription. How do we enable SSO?

Docker SSO is available with a Docker Business subscription. To enable SSO, you must first upgrade your subscription to a Docker Business subscription. To learn how to upgrade your existing account, see [Upgrade your subscription](/subscription/upgrade/){:target="blank" rel="noopener" class=""}.

### Q: Where can I find detailed instructions on how to configure Docker SSO?

For step by step instructions on how to configure Docker SSO, see [Single Sign-on docs](/single-sign-on/#index.md).

### Q: Is it possible to use more than one IdP with Docker SSO?

No. You can only configure Docker SSO to work with a single IdP. A domain can only be associated with a single IdP.Docker currently supports Azure AD and identity providers that support SAML 2.0

### Q: Is it possible to change my identity provider after configuring SSO?

Yes. You must delete your existing IdP configuration in Docker Hub and follow the instructions in the [Single Sign-on docs](/single-sign-on/#index.md) to configure SSO using your IdP.

###  Q: Is Docker SSO available for all paid subscriptions?

Docker SSO is only available with the Docker Business subscription. [Upgrade](/subscription/upgrade/){:target="blank" rel="noopener" class=""} your existing subscription to start using Docker SSO.

### Q: Does Docker SSO support multi-factor authentication (MFA)?

When SSO is being used by an organization, MFA is determined at the idP level and not the Docker system.

### Q: How does service accounts work with SSO?

Service accounts work like any other user when SSO is turned on. If the service account is using an email for a domain with SSO turned on, it needs a [PAT](/docker-hub/access-tokens/) for CLI and API usage.


## Configuration

### Q: What information do I need from my Identity providers to configure SSO?

To enable SSO in Docker, you need the following from your IdP:

* **SAML 2.0**: Entity ID, ACS URL, Single Logout URL and Certificate Download URL
* **Azure AD**: Client ID, Client Secret, AD Domain

### Q: Is DNS verification required to enable SSO?

Yes. You must verify a domain before using it with an SSO connection.

### Q: Does Docker SSO support authenticating through the command line?

Yes. When SSO is enabled, you can access the Docker CLI through Personal Access Tokens (PATs).  Each user must create a PAT to access the CLI. To learn how to create a PAT, see [Managing access tokens](/docker-hub/access-tokens/). Before we transition to PATs, CLI can continue logging in using their personal credentials until early next year to mitigate the risk of interrupting CI/CD pipelines.

###  Q:  How does SSO affect our automation systems and CI/CD pipelines?

Before enforcing SSO, you must create PATs for automation systems and CI/CD pipelines and use the tokens instead of a password.

When SSO is enforced, password-based authentication no longer works on your automation systems and CI/CD pipelines.

### Q: Do I need a specific version of Docker Desktop for SSO?

Yes, all users in your organization must upgrade to Docker Desktop version 4.4.0 or higher. Users on older versions of Docker Desktop will not be able to sign in after enforcing SSO.

### Q: Does SAML authentication require additional attributes?

You must provide an email address as an attribute to authenticate via SAML. The ‘Name’ attribute is currently optional.

### Q: When SAML SSO is enforced, at what stage is the login required to be tracked through SAML? At runtime or install time?

Runtime for Docker Desktop if it’s configured to require authentication to the organization.

## Managing users

### Q: How do I manage users when using SSO?

Users are managed through organizations in Docker Hub. When you configure SSO in Docker, you need to make sure an account exists for each user in your IdP account. When a user signs into Docker for the first time using their domain email address, they will be automatically added to the organization after a successful authentication.

### Q: Do I need to manually add users to my organization?

No, you don’t need to manually add users to your organization in Docker Hub. You just need to make sure an account for your users exists in your IdP and then invite them to your organization using the **Invite Member** option in Docker Hub.

When a user signs into Docker for the first time using their domain email address, they will be automatically added to the organization after a successful authentication.

### Q: Can users in my organization use different email addresses to authenticate via SSO?

During the SSO setup, you’ll have to specify the company email domains that are allowed to authenticate. All users in your organization must authenticate using the email domain specified during SSO setup. Some of your users may want to maintain a different account for their personal projects.

Users with a public domain email address will be added as guests.

### Q: Can Docker Org Owners/Admins approve users to an organization and use a seat, rather than having them automatically added when SSO Is enabled?

Admins and organization owners can currently approve users by configuring their permissions through their IdP. That is, if the user account is configured in the IdP, the user will be automatically added to the organization in Docker Hub as long as there’s an available seat.

### Q: How will users be made aware that they are being made a part of a Docker Org?

When SSO is enabled, users will be prompted to authenticate through SSO the next time they try to sign in to Docker Hub or Docker Desktop. The system will see the end-user has a domain email associated with the docker ID they are trying to authenticate with, and prompts them to sign in with SSO email and credentials instead.

If users attempt to log in through the CLI, they must authenticate using a personal access token (PAT).

### Q: Is it possible to force users of Docker Desktop to authenticate, and/or authenticate using their company’s domain?

Yes. Admins can force users to authenticate with Docker Desktop by provisioning a registry.json configuration file. The registry.json file will force users to authenticate as a user that is configured in the **allowedOrgs** list in the **registry.json file**.

Once SSO enforcement is set up on their DB org on Hub, when the user is forced to auth with Docker Desktop, the SSO enforcement will also force users to authenticate through SSO with their IdP (instead of authenticating using their username and password).

Users may still be able to authenticate as a "guest" account to the organization using a non-domain email address. However, they can only authenticate as guests if that non-domain email was invited to the organization by the organization owner

### Q: Is it possible to convert existing users from non-SSO to SSO accounts?

Yes, you can convert existing users to an SSO account. To convert users from a non-SSO account:

* Ensure your users have a company domain email address and they have an account in your IdP
* Verify that all users have Docker Desktop version 4.4.0 or higher installed on their machines
* Each user has created a PAT to replace their passwords to allow them to log in through Docker CLI
* Confirm that all CI/CD pipelines automation systems have replaced their passwords with PATs.

For detailed prerequisites and for instruction on how to enable SSO, see [Single Sign-on](/single-sign-on/#index.md).

### Q: What impact can users expect once we start onboarding them to SSO accounts?

When SSO is enabled and enforced, your users just have to sign in using the email address and password.

### Q: Is Docker SSO fully synced with Active Directory (AD)?

Docker doesn’t currently support a full sync with AD. That is, if a user leaves the organization, administrators must sign in to Docker Hub and manually remove the user from the organization.

Additionally, you can use our APIs to complete this process.


## Troubleshooting

### Q: What happens if my IdP goes down when SSO is enabled?

It is not possible to access Docker Hub when your IdP is down. However, you can access Docker Hub images from the CLI using your Personal Access Token. Or, if you had an existing account before the SSO enforcement, you can use your username and password to access Docker Hub images during the grace period for your organization.

### Q: What happens when I turn off SSO for my organization?

When you turn off SSO, authentication through your Identity Provider will no longer be required to access Docker. Users may continue to log in through Single Sign-On as well as Docker ID and password.

### Q: What happens if my existing certificate expires?

If your existing certificate has expired, you need to contact your identity provider to generate a new x509 certificate. The new certificate must be added to the SSO configuration settings page on Docker Hub.






