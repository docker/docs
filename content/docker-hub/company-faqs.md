---
description: Company FAQs
keywords: Docker, Docker Hub, SSO FAQs, single sign-on
title: Frequently asked questions
---

### As a Docker Business customer, what steps should I follow to create a company?

Contact your designated CSM team member or Docker Support.

### Are existing subscriptions affected when a company is created and organizations are added to the company?

Subscriptions and related billing details will continue to be managed at the organization level at this time.

### Some of my organizations don’t have a Docker Business subscription. Can I still use a parent company?

Yes, but only organizations with a Docker Business subscription are placed under a company.

### What happens if one of my organizations downgrades from Docker Business, but I still need access as a company owner?

To access and manage child organizations, the organization must have a Docker Business subscription. If the organization isn't a part of this subscription, the owner of the organization must manage the org from the company. 

### Does my organization need to prepare for downtime during the migration process?

No, you can continue with business as usual.

### How many company owners can I add?

A maximum of 10 company owners can be added to a single company account.

### What permission does the company owner have in the associated/nested organizations?

Company owners can navigate to the **Organization** page, view/edit organization members, change SSO/SCIM settings that may impact all users in each organization under the company. However, a company owner can't change any organization repositories. 

### What features are supported at the company level? Will this change over time?

Domain verification, Single Sign-on, and System for Cross-domain Identity Management (SCIM) are supported at the company level. The following aren't supported:

- Image Access Management
- Registry Access Management
- User management
- Billing

### What's required to create a company name?

A company name must be unique to that of its child organization. If a child organization requires the same name as a company, we suggest modifying slightly. For example, **Docker Inc** (parent company), **Docker** (child organization).

### How does a company owner add an organization to the company?

Contact your designated CSM team member or Docker Support with a list of the Docker Business organizations you want to add to the new company.

### How does a company owner manage SSO/SCIM settings from my new parent company?

See Manage your [SCIM](../docker-hub/company-scim.md) and [SSO](../docker-hub/sso-connection.md) settings.

### How does a company owner enable group mapping in my IdP?

See [group mapping](../docker-hub/group-mapping.md) for your IdP.

### What's the definition of a company vs an organization?

A company is a collection of organizations that's managed together. An organization is a collection of repositories and teams that's managed together. 

### What are the different permissions for an organization owner?

Organization owners can create, view, push, and pull repositories from their organization. As a company owner, you don’t have these privileges.

### If an organization isn't part of a company, would SSO or SCIM settings change?

No, the SSO or SCIM settings won't change for that organization. 
