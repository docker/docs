---
description: enforcing sso
keywords: sso, enforce
title: Enforce SSO Login
---

Without SSO enforcement, users can continue to sign in using Docker username and password. If users login with your Domain email, they will authenticate through your identity provider instead. 

You must test your SSO connection first if you’d like to enforce SSO log-in. All users must authenticate with an email address instead of their Docker ID if SSO is enforced


1. In the **Single Sign-On Connections** table, select the Action icon and **Enforce Single Sign-on**.

    > **Note**
    >
    > When you enforce SSO, all members of your organization with a matching domain must authenticate through your IdP. 
2. Continue with the on-screen instructions and verify that you’ve completed the tasks. 
3. Select **Turn on enforcement** to complete. 
