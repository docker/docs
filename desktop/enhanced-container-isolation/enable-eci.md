---
description: Instructions on how to set up enhanced container isolation
title: Enable Enhanced Container Isolation
keywords: set up, enhanced container isolation, rootless, security
---

How to configure it if you are an admin

what you will see as a developer

## How to enable/ get ECI
(e.g. currently developers in Docker Business customers, requires authentication, etc)

requires an Apply and restart
- Admins can lock in the use of the ‘Enhanced container isolation’ mode within their org via the ‘Admin Controls’ feature <link to Admin Controls docs>

To enable Hardened Docker Desktop, Docker Business administrators simply have to toggle on the ‘Hardened Desktop’ option within the Settings panel of their Organization’s space on Docker Hub. Your developers must then authenticate to your organization in Docker Desktop for the settings to be applied. You can follow this simple guide for ensuring developers authenticate to your organization before using Docker Desktop.

Anything that you have the opportunity to configure as an admin, will be locked. Including:

Registry Access Management
Docker Engine runtime will be locked as Sysbox
Proxy settings (TBD)
Other Docker Engine configs (TBD)
Other Docker Desktop configs (TBD)