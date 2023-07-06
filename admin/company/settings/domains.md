---
description: Domain management in Docker Admin
keywords: domains, SCIM, SSO, Docker Admin
title: Domain management
---

{% include admin-early-access.md %}

Use domain management to manage your domains for Single Sign-On and SCIM.

## Add and verify a domain

1. Sign in to [Docker Admin](https://admin.docker.com){: target="_blank" rel="noopener" class="_"}.
2. In the left navigation, select your company in the drop-down menu.
3. Select **Domain management**.
4. Select **Add a domain** and continue with the on-screen instructions to add the TXT Record Value to your domain name system (DNS).

    >**Note**
    >
    > Format your domains without protocol or www information, for example, `yourcompany.example`. This should include all email domains and subdomains users will use to access Docker, for example `yourcompany.example` and `us.yourcompany.example`. Public domains such as `gmail.com`, `outlook.com`, etc. aren’t permitted. Also, the email domain should be set as the primary email.

5. Once you have waited 72 hours for the TXT Record verification, you can then select **Verify** next to the domain you've added, and follow the on-screen instructions.