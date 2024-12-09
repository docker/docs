---
title: Communication and information gathering
description: Gather your company's requirements from key stakeholders and communicate to your developers.
weight: 10
---

## Step one: Communicate with your developers and IT teams

### Docker user communication

You may already have Docker Desktop users within your company, and some steps in this process may affect how they interact with the platform. It's highly recommended to communicate early with users, informing them that as part of the subscription onboarding, they will be upgraded to a supported version of Docker Desktop. 

Additionally, communicate that settings will be reviewed to optimize productivity, and users will be required to sign in to the company’s Docker organization using their business email to fully utilize the subscription benefits.

### MDM team communication

Device management solutions, such as Intune and Jamf, are commonly used for software distribution across enterprises, typically managed by a dedicated MDM team. It is recommended that you engage with this team early in the process to understand their requirements and the lead time for deploying changes.

Several key setup steps in this guide require the use of JSON files, registry keys, or .plist files that need to be distributed to developer machines. It’s a best practice to use MDM tools for deploying these configuration files and ensuring their integrity is preserved.

## Step two: Identify Docker organizations

Some companies may have more than one [Docker organization](/manuals/admin/organization/_index.md) created. These organizations may have been created for specific purposes, or may not be needed anymore. If you suspect your company has more than one Docker organization, it's recommended you survey your teams to see if they have their own organizations. You can also contact your Docker Customer Success representative to get a list of organizations with users whose emails match your domain name.

## Step three: Gather requirements

Through [Settings Management](/manuals/security/for-admins/hardened-desktop/settings-management/_index.md), Docker provides numerous configuration parameters that can be preset. The Docker organization owner, development lead, and infosec representative should review these settings to establish the company’s baseline configuration, including security features and [enforcing sign-in](/manuals/security/for-admins/enforce-sign-in/_index.md) for Docker Desktop users. Additionally, they should decide whether to take advantage of other Docker products, such as [Docker Scout](/manuals/scout/_index.md), which is included in the subscription.

To view the parameters that can be preset, see [Configure Settings Management](/manuals/security/for-admins/hardened-desktop/settings-management/configure-json-file.md#step-two-configure-the-settings-you-want-to-lock-in).

## Optional step four: Meet with the Docker Implementation team

The Docker Implementation team can help you step through setting up your organization, configuring SSO, enforcing sign-in, and configuring Docker. You can reach out to set up a meeting by emailing successteam@docker.com.
