---
title: Plan and roll out Docker
linkTitle: Plan your rollout
weight: 40
description: Coordinate stakeholders, define Docker Desktop settings, test with a pilot group, and roll out Docker across your organization.
keywords: admin, onboarding, deployment, organization setup, Docker Desktop, rollout
aliases:
  - /guides/admin-set-up/
  - /guides/admin-set-up/comms-and-info-gathering/
  - /guides/admin-set-up/deploy/
  - /guides/admin-set-up/finalize-plans-and-setup/
  - /guides/admin-set-up/testing/
---

Use this workflow to plan a Docker rollout across your organization, test its
settings with a pilot group, and coordinate deployment with your IT teams.
For invitations and sign-in configuration, follow
[Onboard a Docker organization](/manuals/accounts/organization/setup/onboard.md).

## Who needs to be involved

Identify the people responsible for each part of the rollout:

- Organization owners to manage membership and organization settings
- DNS and identity provider teams to verify domains and configure single sign-on (SSO)
- IT and mobile device management (MDM) teams to distribute Docker Desktop and configuration files
- Development leads and security representatives to agree on settings and security requirements
- A small group of developers to test the configuration before wider deployment

## Communication and information gathering

### Communicate with your developers and IT teams

Notify existing Docker Desktop users before changing versions, settings, or
sign-in requirements. Explain which changes they should expect, when the pilot
and wider rollout will happen, and where to report problems. Ask users to use
their business email for their Docker accounts.

Engage your MDM team early to understand deployment requirements and lead
times. Agree on how to distribute configuration files, registry keys, or
`.plist` files for your chosen deployment method.

### Identify Docker organizations

Survey teams to find existing Docker organizations and decide which ones
should be part of the rollout. If you need to administer multiple
organizations together, consider a
[Docker company](/manuals/accounts/company/_index.md). A company groups Docker
organizations; it is separate from the business deploying Docker Desktop.

### Gather requirements

Work with organization owners, development leads, and security representatives
to define a baseline that supports development workflows and company policies.
Agree on:

- Docker Desktop versions and update policies
- Settings developers can change and settings administrators must control
- Sign-in requirements, identity provider integration, and user provisioning
- Approved image sources and registries
- Additional Docker products your teams need and the subscriptions they require

Review [Settings Management](/manuals/desktop/enterprise/hardened-desktop/settings-management/_index.md)
and [Hardened Docker Desktop](/manuals/desktop/enterprise/hardened-desktop/_index.md)
to choose controls available with your subscription.

## Finalize plans and begin setup

### Send finalized settings files to the MDM team

Record the agreed baseline and configure it through
[Docker Home](/manuals/desktop/enterprise/hardened-desktop/settings-management/configure-admin-console.md)
or an [`admin-settings.json` file](/manuals/desktop/enterprise/hardened-desktop/settings-management/configure-json-file.md),
as appropriate for your subscription and deployment method. Coordinate with
your MDM team to prepare the configuration and your chosen method for
[enforcing sign-in](/manuals/desktop/enterprise/enforce-sign-in/_index.md).
Deploy to the pilot group before applying the settings more widely.

### Begin setup

Follow [organization onboarding](/manuals/accounts/organization/setup/onboard.md)
to identify users, invite owners and members, and configure sign-in and
provisioning. SSO and SCIM require Docker Business. Plan authentication and
provisioning together: SSO authenticates users, while provisioning controls
how users are added, updated, and removed.

Choose a supported Docker Desktop version and agree with IT on how users will
receive it through MDM, a software portal, or direct installation. See
[Docker Desktop deployment](/manuals/desktop/enterprise/enterprise-deployment.md).

If your rollout includes other Docker products, confirm access and complete
their setup before the pilot. For example, follow the setup documentation for
[Docker Build Cloud](/manuals/build-cloud/setup/_index.md),
[Docker Scout](/manuals/scout/quickstart.md), or
[Docker Hardened Images](/manuals/dhi/get-started.md).

## Testing

### SSO and SCIM testing

If you configure SSO or SCIM, test sign-in, provisioning, role assignments,
and team membership with pilot accounts. Include a user joining the
organization and a user leaving it. Follow the
[provisioning documentation](/manuals/security/provisioning/_index.md)
for the behavior of your chosen method. Check CLI authentication as well as
Docker Desktop and browser sign-in.

### Test Registry Access Management and Image Access Management

If your subscription and rollout include
[Registry Access Management](/manuals/desktop/enterprise/hardened-desktop/registry-access-management.md)
or [Image Access Management](/manuals/desktop/enterprise/hardened-desktop/image-access-management.md),
verify that pilot users can pull approved images and that configured
restrictions block disallowed sources. Tell pilot users about these
restrictions before testing them.

### Deploy settings and enforce sign in to test group

Apply the baseline settings and sign-in requirements to the pilot group.
Have developers run their normal projects on Docker Desktop and Docker Hub.
Confirm that required settings apply and that developers can still build,
pull, and run the containers they need. Use the
[pilot rollout tutorial](/get-started/tutorials/admin-set-up.md) for a focused
settings policy test.

For each additional product in your rollout, have a pilot developer complete
a representative task: build with a cloud builder, review Scout results, or
pull a Hardened Image. Record problems and adjust the baseline before wider
deployment.

## Deploy your Docker setup

Communicate the rollout schedule and expected changes to users. Confirm that
IT and MDM teams are ready to distribute the tested configuration and support
users who encounter problems.

### Enforce SSO

If your plan includes SSO enforcement, verify that your identity provider
configuration covers all intended users before enforcing it. Follow
[Configure SSO](/manuals/security/authentication/single-sign-on/connect.md)
for setup and enforcement. Enforcing SSO and enforcing Docker Desktop sign-in
are separate controls.

### Deploy configuration settings and enforce sign-in to users

Have the MDM team distribute the tested Docker configuration to the remaining
users. Provide instructions for sign-in and a support contact. Review feedback
after each deployment group before expanding the rollout.

### Next steps

- Review [organization Insights](/manuals/accounts/organization/insights.md),
  if available for your subscription, to track Docker Desktop adoption and versions.
- Use the [user access workflow](/manuals/accounts/organization/manage/user-access.md)
  to review membership and permissions after rollout.
- Monitor [Docker Scout findings](/manuals/scout/explore/analysis.md) for
  repositories where you've configured Scout.
