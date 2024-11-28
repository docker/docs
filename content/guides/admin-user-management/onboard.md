---
title: Onboard your users
description:
weight: 10
---

Efficiently onboarding users and managing access is critical to maintaining security and productivity with Docker.

This page provides a top-level look at th tools and techniques to simplify user onboarding and ensure your organization remains secure and scalable.

Step 1: Identify your Docker users and their Docker accounts
Identifying your users will ensure that you allocate your subscription seats efficiently and that all your Docker users receive the benefits of your subscription.

Identify the Docker users in your organization.
If your organization uses device management software, like MDM or JAMF, you may use the device management software to help identify Docker users. See your device management software's documentation for details. You can identify Docker users by checking if Docker Desktop is installed at the following location on each user's machine:
Mac: /Applications/Docker.app
Windows: C:\Program Files\Docker\Docker
Linux: /opt/docker-desktop
If your organization doesn't use device management software or your users haven't installed Docker Desktop yet, you may survey your users.
Instruct all your Docker users in your organization to update their existing Docker account's email address to an address that's in your organization's domain, or to create a new account using an email address in your organization's domain.
To update an account's email address, instruct your users to sign in to Docker Hub, and update the email address to their email address in your organization's domain.
To create a new account, instruct your users to go sign up using their email address in your organization's domain.
Ask your Docker sales representative or contact sales to get a list of Docker accounts that use an email address in your organization's domain.



## Step one: Secure user authentication with Single Sign-On (SSO)

Single Sign-On (SSO) is essential for centralized and secure user authentication. By integrating Docker with your identity provider, you can enforce compliance and streamline user access.

### Benefits of SSO

Centralized Authentication: Users log in with their organizational credentials.
Improved Security: Eliminates password-related vulnerabilities.
Simplified Onboarding: SSO works seamlessly with SCIM and group mapping for automated user provisioning.

### Steps to Configure SSO

Enable SSO in the Admin Console under Organization Settings.
Connect your identity provider using SAML or OIDC.
Test the SSO setup with a small group of users before rolling it out organization-wide.

[Detailed SSO setup documentation](/manuals/security/for-admins/)

## Step two: Onboarding Users
Adding users to your Docker Business organization involves a mix of manual invitations and automated provisioning methods. Choose the approach that best suits your organization’s needs.

### Manual User Onboarding

Access the Admin Console: Use the Admin Console in Docker Hub to invite and manage users.
→ Admin Console Overview
Send Invitations: Select Invite Users and enter their email addresses. Assign a default role during the invitation process.
→ How to Invite Users
Monitor Invitations: Regularly check pending invitations to ensure all users complete the sign-up process.

### Automated User Provisioning with SCIM

For larger organizations, System for Cross-domain Identity Management (SCIM) automates user provisioning and de-provisioning, reducing manual workload and minimizing errors.

How SCIM Works: SCIM connects Docker to your identity provider (e.g., Okta, Azure AD) to automatically add, update, or remove users based on changes in your directory.
Steps to Enable SCIM:
Configure SCIM in your identity provider.
Generate a SCIM token in the Docker Admin Console.
Link the token to your identity provider.

## Step three: Streamlining User Access with Group Mapping
Group mapping simplifies access management by linking identity provider groups to Docker roles and teams. This feature ensures users are automatically assigned the correct permissions based on their directory group membership.

### How Group Mapping Works
Identity provider groups (e.g., “Developers” or “Admins”) are mapped to specific Docker roles and teams.
When users are added to these groups in your directory, their Docker permissions are automatically updated.
### Benefits of Group Mapping
Reduces manual assignment errors.
Ensures consistent access control policies.
Simplifies scaling permissions as teams grow or change.

## Step four: Future proofing onboarding with Just-in-Time (JIT) Provisioning
Just-in-time provisioning ensures that users are added to your Docker organization the first time they log in, based on their identity provider credentials. This feature eliminates the need for pre-inviting users while still enforcing role-based access control.

### How JIT Works
Users authenticate via SSO.
During their first login, they are automatically added to your Docker organization and assigned roles based on group mapping.
### Benefits of JIT Provisioning
Streamlines onboarding for large or distributed teams.
Reduces admin overhead by removing the need for manual user invites.
Works seamlessly with SCIM and SSO for a fully automated provisioning process.


## Best Practices for User Management
Combine SCIM and Group Mapping: Use SCIM for user synchronization and group mapping to automate role assignments.
Leverage JIT Provisioning: Enable JIT for dynamic onboarding without manual invites.
Monitor Activity: Regularly review audit logs to track access and changes.
→ Using Audit Logs
Regularly Review Permissions: Periodically check and adjust group mappings and roles to align with organizational changes.


## More resources

https://docs.docker.com/admin/organization/onboard/