---
title: Mastering user and access management
summary: Simplify user access while ensuring security and efficiency in Docker.
description: A guide for managing roles, provisioning users, and optimizing Docker access with tools like SSO and activity logs.
keywords: admin, user management, roles, permissions, sso, provisioning, access control
aliases:
  - /guides/admin-user-management/audit-and-monitor/
  - /guides/admin-user-management/onboard/
  - /guides/admin-user-management/setup/
params:
  tags: [admin]
  featured: false
  time: 20 minutes
  image:
---


Managing roles and permissions is key to securing your Docker environment while enabling easy collaboration and operational efficiency. This guide walks IT administrators through the essentials of user and access management, offering strategies for assigning roles, provisioning users, and using tools like activity logs and Insights to monitor and optimize Docker usage.

## Who's this for?

- IT teams tasked with configuring and maintaining secure user access
- Security professionals focused on enforcing secure access practices
- Project managers overseeing team collaboration and resource management

## What you'll learn

- How to assess and manage Docker user access and align accounts with organizational needs
- When to use team configurations for scalable access control
- How to automate and streamline user provisioning with SSO, SCIM, and JIT
- How to get the most out of Docker's monitoring tools

## Tools integration

This guide covers integration with:

- Okta
- Entra ID SAML 2.0
- Azure Connect (OIDC)

## Setting up roles and permissions in Docker

With the right configurations, you can ensure your developers have easy access to necessary resources while preventing unauthorized access. This page guides you through identifying Docker users so you can allocate subscription seats efficiently within your Docker organization, and assigning roles to align with your organization's structure.

### Identify your Docker users and accounts

Before setting up roles and permissions, it's important to have a clear understanding of who in your organization requires Docker access. Focus on gathering a comprehensive view of active users, their roles within projects, and how they interact with Docker resources. This process can be supported by tools like device management software or manual assessments. Encourage all users to update their Docker accounts to use organizational email addresses, ensuring seamless integration with your subscription.

For steps on how you can do this, see [step 1 of onboarding your organization](/manuals/admin/organization/setup/onboard.md).

### Assign roles strategically

When you invite members to join your Docker organization, you assign them a role.

Docker's predefined roles offer flexibility for various organizational needs. Assigning roles effectively ensures a balance of accessibility and security.

- Member: Non-administrative role. Members can view other members that are in the same organization.
- Editor: Partial administrative access to the organization. Editors can create, edit, and delete repositories. They can also edit an existing team's access permissions.
- Owner: Full organization administrative access. Owners can manage organization repositories, teams, members, settings, and billing.

For more information, see [Roles and permissions](/manuals/enterprise/security/roles-and-permissions.md).

#### Enhance with teams

Teams in Docker provide a structured way to manage member access and they provide an additional level of permissions. They simplify permission management and enable consistent application of policies.

- Organize users into teams aligned with projects, departments, or functional roles. This approach helps streamline resource allocation and ensures clarity in access control.
- Assign permissions at the team level rather than individually. For instance, a development team might have "Read & Write" access to certain repositories, while a QA team has "Read-only" access.
- As teams grow or responsibilities shift, you can easily update permissions or add new members, maintaining consistency without reconfiguring individual settings.

For more information, see [Create and manage a team](/manuals/admin/organization/manage/manage-a-team.md).

#### Example scenarios

- Development teams: Assign the member role to developers, granting access to the repositories needed for coding and testing.
- Team leads: Assign the editor role to team leads for resource management and repository control within their teams.
- Organizational oversight: Restrict the organization owner or company owner roles to a select few trusted individuals responsible for billing and security settings.

#### Best practices

- Apply the principle of least privilege. Assign users only the minimum permissions necessary for their roles.
- Conduct regular reviews of role assignments to ensure they align with evolving team structures and organizational responsibilities.

## Onboarding and managing roles and permissions in Docker

This page guides you through onboarding owners and members, and using tools like SSO and SCIM to future-proof onboarding going forward.

### Invite owners

When you create a Docker organization, you automatically become its sole owner. While optional, adding additional owners can significantly ease the process of onboarding and managing your organization by distributing administrative responsibilities. It also ensures continuity and prevents blockers if the primary owner is unavailable.

For detailed information on owners, see [Roles and permissions](/manuals/enterprise/security/roles-and-permissions.md).

### Invite members and assign roles

Members are granted controlled access to resources and enjoy enhanced organizational benefits. When you invite members to join your Docker organization, you immediately assign them a role.

#### Benefits of inviting members

- Enhanced visibility: Gain insights into user activity, making it easier to monitor access and enforce security policies.
- Streamlined collaboration: Help members collaborate effectively by granting access to shared resources and repositories.
- Improved resource management: Organize and track users within your organization, ensuring optimal allocation of resources.
- Access to enhanced features: Members benefit from organization-wide perks, such as increased pull limits and access to premium Docker features.
- Security control: Apply and enforce security settings at an organizational level, reducing risks associated with unmanaged accounts.

For detailed information, see [Manage organization members](/manuals/admin/organization/manage/members.md).

### Future-proof user management

A robust, future-proof approach to user management combines automated provisioning, centralized authentication, and dynamic access control. Implementing these practices ensures a scalable, secure, and efficient environment.

#### Secure user authentication with single sign-on (SSO)

Integrating Docker with your identity provider streamlines user access and enhances security.

SSO:

- Simplifies sign in, as users sign in with their organizational credentials.
- Reduces password-related vulnerabilities.
- Simplifies onboarding as it works seamlessly with SCIM and group mapping for automated provisioning.

For more information, see the [SSO documentation](/manuals/enterprise/security/single-sign-on/_index.md).

#### Automate onboarding with SCIM and JIT provisioning

Streamline user provisioning and role management with [SCIM](/manuals/enterprise/security/provisioning/scim/_index.md) and [Just-in-Time (JIT) provisioning](/manuals/enterprise/security/provisioning/just-in-time.md).

With SCIM you can:

- Sync users and roles automatically with your identity provider.
- Automate adding, updating, or removing users based on directory changes.

With JIT provisioning you can:

- Automatically add users upon first sign in based on [group mapping](#simplify-access-with-group-mapping).
- Reduce overhead by eliminating pre-invite steps.

#### Simplify access with group mapping

Group mapping automates permissions management by linking identity provider groups to Docker roles and teams.

It also:

- Reduces manual errors in role assignments.
- Ensures consistent access control policies.
- Help you scale permissions as teams grow or change.

For more information on how it works, see [Group mapping](/manuals/enterprise/security/provisioning/scim/group-mapping.md).

## Monitoring and insights

Activity logs and Insights are useful tools for user and access management in Docker. They provide visibility into user actions, team workflows, and organizational trends, helping enhance security, ensure compliance, and boost productivity.

### Activity logs

Activity logs track events at the organization and repository levels, offering a clear view of activities like repository changes, team updates, and billing adjustments.

Activity logs are available for Docker Team or Docker Business plans, with data retained for three months.

#### Key features

- Change tracking: View what changed, who made the change, and when.
- Comprehensive reporting: Monitor critical events such as repository creation, deletion, privacy changes, and role assignments.

#### Example scenarios

- Audit trail for security: A repository’s privacy settings were updated unexpectedly. The activity logs reveal which user made the change and when, helping administrators address potential security risks.
- Team collaboration review: Logs show which team members pushed updates to a critical repository, ensuring accountability during a development sprint.
- Billing adjustments: Track who added or removed subscription seats to maintain budgetary control and compliance.

For more information, see [Activity logs](/manuals/admin/activity-logs.md).

### Insights

Insights provide data-driven views of Docker usage to improve team productivity and resource allocation.

#### Key benefits

- Standardized environments: Ensure consistent configurations and enforce best practices across teams.
- Improved visibility: Monitor metrics like Docker Desktop usage, builds, and container activity to understand team workflows and engagement.
- Optimized resources: Track license usage and feature adoption to maximize the value of your Docker subscription.

#### Example scenarios

- Usage trends: Identify underutilized licenses or resources, allowing reallocation to more active teams.
- Build efficiency: Track average build times and success rates to pinpoint bottlenecks in development processes.
- Container utilization: Analyze container activity across departments to ensure proper resource distribution and cost efficiency.

For more information, see [Insights](/manuals/admin/insights.md).

### Next steps

Now that you've mastered user and access management in Docker, you can:

- Review your [activity logs](/manuals/admin/activity-logs.md) regularly to maintain security awareness
- Check your [Insights dashboard](/manuals/admin/insights.md) to identify opportunities for optimization
- Explore [advanced security features](/manuals/enterprise/security/_index.md) to further enhance your Docker environment
- Share best practices with your team to ensure consistent adoption of security policies
