---
title: Single sign-on overview
linkTitle: Single sign-on
description: >-
  How Docker single sign-on works for Business organizations and companies,
  and the steps to configure a domain, IdP connection, and optional
  enforcement.
keywords: single sign-on, SSO, sign-in, admin, Docker Hub, Docker Home,
  security, identity provider, IdP, SSO configuration, SSO setup, Docker
  Business, AI Governance, SAML, domain verification, enforce SSO, JIT
  provisioning, personal access token
aliases:
  - /single-sign-on/
  - /admin/company/settings/sso/
  - /admin/organization/security-settings/sso-management/
  - /security/for-admins/single-sign-on/
  - /enterprise/security/single-sign-on/
weight: 10
---

{{< summary-bar feature_name="SSO" >}}

Single sign-on (SSO) lets users access Docker by authenticating through your
identity provider (IdP). You can configure SSO for an entire company,
including all associated organizations, or for a single organization. SSO
requires a Docker Business subscription or an
[AI Governance](/manuals/accounts/organization/manage/manage-licenses.md)
license.

## How SSO works

When a user signs in to Docker Desktop or Docker Hub, Docker sends them to
your IdP's sign-in page. The IdP accepts or denies the sign-in.

```mermaid {title="Sign-in when the IdP accepts the user" caption="The user has no Docker session, and the IdP accepts them."}
sequenceDiagram
  actor User
  participant App as Docker Desktop<br/>or<br/>Docker Hub
  participant Login as Docker Login
  participant IdP as Your IdP

  User->>App: Open the app
  App->>Login: Start authentication
  Login->>IdP: Send the user to your IdP
  IdP->>User: Prompt for credentials
  User->>IdP: Enter credentials
  IdP->>Login: Confirm the user
  Login->>App: Return the user to the app
```

## Set up SSO

Setting up SSO is a multi-step procedure that moves between Docker and your
IdP.

```mermaid {title="SSO setup flow" caption="The hexagon marks the IdP step. After the connection works, complete the post-setup tasks."}
flowchart LR
  %% Link the columns to each other. A link into a step inside a
  %% column makes that column run left to right.
  subgraph configure["Configure SSO"]
    direction TB
    A["1. Add a domain"] --> B["2. Verify your domain"]
    B --> C["3. Create an SSO connection in Docker"]
    C --> D{{"4. Create an SSO connection in your IdP"}}
    D --> E["5. Connect Docker to your IdP"]
    E --> F["6. Test the connection"]
  end

  subgraph post["Post-setup"]
    direction TB
    G["Provision users"]
    H["Optional: Enforce SSO / Desktop sign-in"]
    I["Manage SSO configuration"]
    G --> H --> I
  end

  configure -->|"After the connection works"| post

  classDef step fill:#ffffff,stroke:#6b7280,color:#111827
  class A,B,C,D,E,F,G,H,I step
  style configure fill:#dbeafe,stroke:#2563eb,stroke-width:2px
  style post fill:#fef9c3,stroke:#ca8a04,stroke-width:2px
```

After you finish the setup, users sign in to Docker services with their
company email address. When
[Just-in-Time provisioning](/manuals/security/provisioning/just-in-time.md)
is turned on, which is the default for an SSO connection, Docker adds each
user to your organization and team on first sign-in.

> [!IMPORTANT]
>
> When SSO is enforced, the Docker CLI rejects password sign-in. Use a
> personal access token (PAT) to sign in to the CLI.

## Next steps

- Start [configuring SSO](connect.md).
- Read the [FAQs](/manuals/faqs/security.md).
- [Troubleshoot SSO issues](troubleshoot-sso.md).
