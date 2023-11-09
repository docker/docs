---
description: Company FAQs
keywords: Docker, Docker Hub, SSO FAQs, single sign-on
title: FAQs on companies
---

### Are existing subscriptions affected when a company is created and organizations are added to the company?

Subscriptions and related billing details continue to be managed at the organization level at this time.

### Some of my organizations don’t have a Docker Business subscription. Can I still use a parent company?

Yes, but only organizations with a Docker Business subscription can be placed under a company.

### What happens if one of my organizations downgrades from Docker Business, but I still need access as a company owner?

To access and manage child organizations, the organization must have a Docker Business subscription. If the organization isn’t included in this subscription, the owner of the organization must manage the organization outside of the company.

### Does my organization need to prepare for downtime during the migration process?

No, you can continue with business as usual.

### How many company owners can I add?

A maximum of 10 company owners can be added to a single company account.

### What permission does the company owner have in the associated/nested organizations?

Company owners can navigate to the **Organizations** page to view all their nested organizations in a single location. They can also view or edit organization members and change SSO and SCIM settings. Changes to company settings impact all users in each organization under the company.

### What features are supported at the company level?

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

See your [SCIM](scim.md) and [SSO](../security/for-admins/single-sign-on/configure/index.md) settings.

### How does a company owner enable group mapping in my IdP?

See [SCIM](scim.md) for more information.

### What's the definition of a company vs an organization?

A company is a collection of organizations that are managed together. An organization is a collection of repositories and teams that are managed together.