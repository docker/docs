---
title: Disable persistent cookies
description: Learn how to disable persistent cookies for Docker Trusted Registry.	
keywords: dtr, browser cookies, sso
---

>{% include enterprise_label_shortform.md %}

If you want your Docker Trusted Registry (DTR) to use session-based authentication cookies that expire when you close your browser, toggle "Disable persistent cookies".
  
![](/ee/dtr/images/disable-persistent-cookies-1.png){: .with-border}

## Verify your DTR cookies setting

You may need to disable Single Sign-On (SSO). From the DTR web UI in a Chrome browser, right-click on any page and click **Inspect**. With the Developer Tools open, select **Application > Storage > Cookies > `https://<dtr-external-url>`**. Verify that the cookies has "Session" as the setting for **Expires / Max-Age**.

## Where to go next

- [Use your own TLS certificates](use-your-own-tls-certificates)
- [Enable authentication using client certificates](/ee/enable-client-certificate-authentication/)
