---
title: Communication and information gathering
description: Gather your company's requirements from key stakeholders and communicate to your developers.
weight: 10
---

## Communicate with your developers and IT teams

Before rolling out Docker Desktop across your organization, coordinate with key stakeholders to ensure a smooth transition.

### Notify Docker Desktop users

You may already have Docker Desktop users within your company. Some steps in
this onboarding process may affect how they interact with the platform.

Communicate early with users to inform them that:

- They'll be upgraded to a supported version of Docker Desktop as part of the subscription onboarding
- Settings will be reviewed and optimized for productivity
- They'll need to sign in to the company's Docker organization using their
business email to access subscription benefits

### Engage with your MDM team

Device management solutions, such as Intune and Jamf, are commonly used for
software distribution across enterprises. These tools are typically managed by a dedicated MDM team.

Engage with this team early in the process to:

- Understand their requirements and lead time for deploying changes
- Coordinate the distribution of configuration files

Several setup steps in this guide require JSON files, registry keys, or .plist
files to be distributed to developer machines. Use MDM tools to deploy these configuration files and ensure their integrity.

## Identify Docker organizations

Some companies may have more than one
[Docker organization](/manuals/admin/organization/_index.md) created. These
organizations may have been created for specific purposes, or may not be
needed anymore.

If you suspect your company has multiple Docker organizations:

- Survey your teams to see if they have their own organizations
- Contact your Docker Support to get a list of organizations with users whose
    emails match your domain name

## Gather requirements

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

## Optional: Meet with the Docker Implementation team

The Docker Implementation team can help you set up your organization,
configure SSO, enforce sign-in, and configure Docker Desktop.

To schedule a meeting, email successteam@docker.com.
