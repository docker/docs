---
title: Communication and information gathering
description: Gather your company's requirements from key stakeholders and communicate to your developers.
weight: 10
---


Docker user communication
You may already have Docker Desktop users in your company.  Some of the steps in this process may cause changes in how they use Docker Desktop.  It’s recommended that you send out a communication up front to the users letting them know that as part of the subscription onboarding process you will be upgrading existing Docker Desktop users to a supported version of the product, reviewing settings to help user productivity, and requiring users to sign in to the company’s Docker org with their business email so they are using the subscription.

MDM team communication
Device management solutions like Intune and Jamf are a standard way to distribute software across enterprises.  There is typically a MDM team that manages this tool.  We recommend talking with that team early in the process to understand their requirements and lead time on distributing changes.  The Docker configurations can include both JSON files and/or registry key/plist entries that will be distributed to developer machines.  It is recommended to use MDM tooling to both distribute configuration files, and ensure their contents don’t change.

Identify Organizations
Some companies may have more than one Docker organization created.  These organizations may have been created for specific purposes, or may not be needed anymore.  If you suspect your company has more than one organization, it's recommended you survey your teams to see if they have their own organizations.  You can also contact your Docker CS representative to get a list of organizations with users whose emails match your domain name.

Baseline configuration discussions
Docker offers a significant number of configuration parameters that can be preset, including enforcing sign in for Docker Desktop users.  The Docker organization owner and the development lead should review the settings to determine which of those settings to configure to create the company’s baseline configuration.  There are also settings for the free trials of other Docker products included in the subscription.  The list of configurations that can be preset is located here.

Security configuration discussions
Docker offers a number of security related features that have configuration parameters that can be preset.  The infosec representative, Docker organization owner, and the development lead should review those features to determine which they want to enable as part of the company’s baseline configuration.  The list of security related features is located here.

Meet with the Docker implementation team
The Docker Implementation Team can help you step through setting up your organization, configuring SSO, enforcing sign in, and configuring Docker.  You can reach out to set up a meeting by emailing successteam@docker.com

SSO domain verification
The SSO process has multiple steps involving different teams, so it's recommended that the process is started right away.  The first step is domain verification.  This step ensures that the person setting up SSO actually controls the domain they are requesting.  The detailed steps to verify a domain are located here.  Your DNS team will need to be involved in this step.