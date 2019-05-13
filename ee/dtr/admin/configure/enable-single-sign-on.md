---
title: Enable single sign-on
description: Learn how to set up single sign-on between UCP and DTR, so that your users only have to authenticate once
keywords: dtr, login, sso
---

Users are shared between UCP and DTR by default, but the applications have separate browser-based interfaces which require authentication.

To only authenticate once, you can configure DTR to have single sign-on (SSO) with UCP.

> **Note**: After configuring single sign-on with DTR, users accessing DTR via 
> `docker login` should create an [access token](/ee/dtr/user/access-tokens/) and use it to authenticate. 

## At install time

When [installing DTR](/reference/dtr/2.7/install/), pass the `--dtr-external-url <url>`
option to enable SSO. This makes it so that when you access DTR's web interface, you are redirected to the UCP login page for authentication. Upon successfully logging in, you are then redirected to your specified DTR external URL during installation.

[Specify the Fully Qualified Domain Name (FQDN)](/use-your-own-tls-certificates/) of your DTR, or a load balancer, to load-balance requests across multiple DTR replicas.

## Post-installation

1. Navigate to `https://<dtr-url>` and log in with your credentials. 
2. Select **System** from the left navigation pane, and scroll down to **Domain & Proxies**. 
3. Update the **Load balancer / Public Address** field with the external URL where users
should be redirected once they are logged in. Click **Save** to apply your changes.
4. Toggle **Single Sign-on** to automatically redirect users to UCP for logging in.

## Where to go next

- [Use your own TLS certificates](use-your-own-tls-certificates)
- [Enable authentication using client certificates](/ee/enable-authentication-via-client-certs/)
