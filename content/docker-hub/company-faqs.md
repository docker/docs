---
description: Company FAQs
keywords: Docker, Docker Hub, SSO FAQs, single sign-on
title: FAQs on companies
---

### Are existing subscriptions affected when a company is created and organizations are added to the company?

You can manage subscriptions and related billing details at the organization level.

### Some of my organizations don’t have a Docker Business subscription. Can I still use a parent company?

Yes, but you can only add organizations with a Docker Business subscription to a company.

### What happens if one of my organizations downgrades from Docker Business, but I still need access as a company owner?

To access and manage child organizations, the organization must have a Docker Business subscription. If the organization isn’t included in this subscription, the owner of the organization must manage the organization outside of the company.

### Does my organization need to prepare for downtime during the migration process?

No, you can continue with business as usual.

### How many company owners can I add?

You can add a maximum of 10 company owners to a single company account.

### What permissions does the company owner have in the associated/nested organizations?

Company owners can navigate to the **Organizations** page to view all their nested organizations in a single location. They can also view or edit organization members and change SSO and SCIM settings. Changes to company settings impact all users in each organization under the company. See [Roles and permissions](../security/for-admins/roles-and-permissions.md).

### What features are supported at the company level?

You can manage domain verification, Single Sign-on, and System for Cross-domain Identity Management (SCIM) at the company level. The following features aren't supported at the company level, but you can manage them at the organization level:

- Image Access Management
- Registry Access Management
- User management
- Billing

To view and manage users across all the organizations under your company, you can [manage users at the company level](../admin/company/users.md) when you use Docker Admin.

Domain audit is not supported for companies or organizations within a company.

### What's required to create a company name?

A company name must be unique to that of its child organization. If a child organization requires the same name as a company, we suggest modifying slightly. For example, **Docker Inc** (parent company), **Docker** (child organization).

### How does a company owner add an organization to the company?

You can add organizations to a company in [Docker Admin](../admin/company/organizations.md/#add-organizations-to-a-company.md) or [Docker Hub](./new-company.md/#add-organizations-to-a-company.md).

### How does a company owner manage SSO/SCIM settings from my new parent company?

See your [SCIM](scim.md) and [SSO](../security/for-admins/single-sign-on/configure/index.md) settings.

### How does a company owner enable group mapping in my IdP?

See [SCIM](scim.md) and [Group mapping](../security/for-admins/group-mapping.md) for more information.

### What's the definition of a company vs an organization?

A company is a collection of organizations that are managed together. An organization is a collection of repositories and teams that are managed together.