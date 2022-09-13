---
description: admin controls for desktop
keywords: admin, controls, rootless, enhanced container isolation
title: Desktop Managed (Early Access)
--- 

>Note
>
>Desktop Managed is currently in [Early Access](../../release-lifecycle.md#early-access-ea) and available to Docker Business customers only. 

## What it is

When admins create a setting and have ‘Rootless Containers’ enabled, the setting cannot be modified by developers within their org

##  What are the benefits



we need to think about how best to convey to users via the docs that these settings will only be respected in the VM if the user has ‘Enhanced Container Isolation’ toggled on.


should note our competitive advantage here, e.g. that when the ‘Enhanced container isolation’ setting is configured, these settings cannot be modified by developers (loop in Cesar / Rodny to advise on wording)

Admins for Docker Business organizations will get Admin Controls, a feature allowing them to enforce certain Docker Desktop settings for their organization. Specifically, admins will be able to enforce:

Send usage statistics is also locked

Admins can lock down any values configurable via the Resources > Network tab via the admin-settings.json. For the time being, this is just the Docker subnet value (screenshot enclosed)


Main usage for this is to lock enterprise related settings:
 - proxies (so users don’t have to set up them / to know that they exist / to bypass them)
 - vpnkit CIDR (to avoid network clashes)
 - VM settings
 - block telemetry
 - auto update behavior
 - daemon config (we want to lock some fields with optional value but keep other ones free to use)

 If ‘Software Updates’ are locked by the admin:

‘Preferences’ section is not shown on ‘Software Updates’ panel

User cannot see A new update is ready to download text

User will get the following message You're currently on version X. The latest version is Y. Updates are managed by your admin.

Users cannot see the Download update button

Users are still able to see the description of the new release as well as the associated Release notes button

The admin should be able to configure all proxy values available via the Docker Desktop Preferences > Resources >  Proxies UI, via the admin-settings.json file.

Acceptance criteria

Admins should have the ability to enforce the use of Hyper-V OR WSL2



## What can be set?

 the Admin can lock via the admin-settings.json



these details should include the exact syntax / options that the admin can use in the admin-settings.json to configure each setting

where ‘enhanced container isolation’ is a prerequisite to ensure that some settings are enforced within the Docker Desktop Linux VM, clearly denote this to the user (Cesar and Rodny can advise on this one)

## What do developers see 

Explain what happens on the developer side, e.g. once you configure your settings, your developer will see that they are locked by their org admin in the Docker Desktop UI

## How to set it up

Details on where the admin-settings.json should be placed on Windows and macOS

Explain that for this feature to take effect, developers must authenticate to their Docker Business org. In order to ensure that this happens, admins must use the registry.json file (link to relevant doc)

