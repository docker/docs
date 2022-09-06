---
description: admin controls for desktop
keywords: admin, controls, rootless, enhanced container isolation
title: Desktop Managed (Early Access)
--- 

we need to think about how best to convey to users via the docs that these settings will only be respected in the VM if the user has ‘Enhanced Container Isolation’ toggled on.

High-level introduction to feature and it’s benefits (can grab some of this from the PR-FAQ, when ready). Include that this is for Docker Business customers only

should note our competitive advantage here, e.g. that when the ‘Enhanced container isolation’ setting is configured, these settings cannot be modified by developers (loop in Cesar / Rodny to advise on wording)

Details on each setting that the Admin can lock via the admin-settings.json

these details should include the exact syntax / options that the admin can use in the admin-settings.json to configure each setting

where ‘enhanced container isolation’ is a prerequisite to ensure that some settings are enforced within the Docker Desktop Linux VM, clearly denote this to the user (Cesar and Rodny can advise on this one)

Explain what happens on the developer side, e.g. once you configure your settings, your developer will see that they are locked by their org admin in the Docker Desktop UI

Details on where the admin-settings.json should be placed on Windows and macOS

Explain that for this feature to take effect, developers must authenticate to their Docker Business org. In order to ensure that this happens, admins must use the registry.json file (link to relevant doc)