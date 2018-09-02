---
title: Enable SAML authentication
description: Learn how configure user authentication with SAML 2.0
keywords: SAML, ucp, authentication
---

> Beta disclaimer
>
> This is beta content. It is not yet complete and should be considered a work in progress. This content is subject to change without notice.

Security Assertion Markup Language (SAML) is an open standard for exchanging authentication data between an identity provider and a service provider. SAML-based single sign-on (SSO) gives you access to UCP through an identity provider of your choice. For more information about SAML, see the [SAML XML website] (http://saml.xml.org/)

A list of the identity providers that we support is found in the [Docker Compatibility matrix] (http://success.docker.com/article/compatibility-matrix).

## Prerequisites

Before you can enable SAML authentication, you must first be set up with your identity provider of choice. This process varies from provider to provider, so consult your provider's documentation for details. There are specific bits of information you need from the identity provider to enable UCP to authenticate with that identity.
