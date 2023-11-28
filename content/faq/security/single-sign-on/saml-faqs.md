---
description: Single Sign-on FAQs
keywords: Docker, Docker Hub, SSO FAQs, single sign-on
title: SAML
aliases:
- /single-sign-on/saml-faqs/
---

### Does SAML authentication require additional attributes?

You must provide an email address as an attribute to authenticate through SAML. The ‘Name’ attribute is optional.

### Does the application recognize the NameID/Unique Identifier in the SAMLResponse subject?

The preferred format is your email address, which should also be your Name ID.

### When you enforce SAML SSO, at what stage is the login required for tracking through SAML? At runtime or install time?

At runtime for Docker Desktop if it’s configured to require authentication to the organization.

### Do you have any information on how to use the Docker Desktop application in accordance with the SSO users we provide? How can we verify that we're handling the licensing correctly?

Verify that your users have downloaded the latest version of Docker Desktop. An enhancement in user management observability and capabilities will become available in the future.