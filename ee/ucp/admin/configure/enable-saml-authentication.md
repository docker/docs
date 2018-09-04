---
title: Enable SAML authentication
description: Learn how configure user authentication with SAML 2.0
keywords: SAML, ucp, authentication
---

> Beta disclaimer
>
> This is beta content. It is not yet complete and should be considered a work in progress. This content is subject to change without notice.

UCP supports Security Assertion Markup Language (SAML) for authentication as a service provider. SAML is an open standard for exchanging authentication data between an identity provider and a service provider, such as UCP. SAML is commonly supported by enterprise authentication systems. SAML-based single sign-on (SSO) gives you access to UCP through a SAML 2.0-compliant identity provider. For more information about SAML, see the [SAML XML website] (http://saml.xml.org/)

UCP supports these identity providers:

- (Okta) [https://www.okta.com/]
- (ADFS) [https://docs.microsoft.com/en-us/windows-server/identity/active-directory-federation-services]
- (Ping Identity) [https://www.pingidentity.com/en/platform/single-sign-on/sso-overview.html]

## Configure identity provider integration

There are values your identity provider needs for successful integration with UCP. These values can vary between identity providers. Consult your identity provider documentation for instructions on providing these values as part of their integration process.


### Okta integration values

The integration values required by Okta are:

- URL for single signon (SSO). This value is the URL for UCP, qualified with `/enzi/v0/saml/acs`. For example, `https://111.111.111.111/enzi/v0/saml/acs`.
- Service provider audience URI. This value is the URL for UCP, qualified with `/enzi/v0/saml/metadata`. For example, `https://111.111.111.111/enzi/v0/saml/metadata`.
- NameID format. Select Unspecified.
- Application username. Email (For example, a custom `${f:substringBefore(user.email, "@")}` specifies the username in the email address.
- Attribute Statements:
    - Name: `fullname`, Value: `user.displayName`.
    - Group Attribute Statement:
Name: `member-of`, Filter: (user defined) for associate group membership. The group name is returned with the assertion.
Name: `is-admin`, Filter: (user defined) for identifying if the user is an admin.



### ADFS integration values

The integration values required by ADFS are:

- (need values)
-
### Ping integration values

The integration values required by Ping Identity are:

- (need values)

## Configure the SAML integration

To enable SAML authentication, go to the UCP web UI, then navigate to the **Admin Settings**. Select **Authentication & Authorization** to enable SAML.

![Enabling SAML in UCP](../../images/saml_enabled.png)

In the **SAML Enabled** section, select **Yes** to display the required settings

![Configuring SAML in UCP](../../images/saml_settings.png)

1. In **IdP Metadata URL** enter the URL for the identity provider's metadata.
2. In **UCP Host** enter the URL that includes the IP address of your UCP console.
3. Select **Save** to complete the integration.

## Security considerations

You can download a client bundle to access UCP. To ensure that access from the client bundle is synced with the identity provider, we recommend the following steps. Otherwise, a previously-authorized user could get access to UCP through their existing client bundle.

- Remove the user account from UCP granting client bundle access if access is removed from the identity provider.
- If group membership in the identity provider changes, replicate this change in UCP.
- Continue to use LDAP to sync group membership.
