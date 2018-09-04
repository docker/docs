---
title: Enable SAML authentication
description: Learn how configure user authentication with SAML 2.0
keywords: SAML, ucp, authentication
---

> Beta disclaimer
>
> This is beta content. It is not yet complete and should be considered a work in progress. This content is subject to change without notice.

Ping Identity integration requires these values:

SAML is commonly supported by enterprise authentication systems. SAML-based single sign-on (SSO) gives you access to UCP through a SAML 2.0-compliant identity provider. For more information about SAML, see the [SAML XML website] (http://saml.xml.org/).

SAML-based single sign-on (SSO) gives you access to UCP through a SAML 2.0-compliant identity provider. UCP supports SAML for authentication as a service provider integrated with your identity provider.

UCP supports these identity providers:

- (Okta) [https://www.okta.com/]
- (ADFS) [https://docs.microsoft.com/en-us/windows-server/identity/active-directory-federation-services]
- (Ping Identity) [https://www.pingidentity.com/en/platform/single-sign-on/sso-overview.html]


## Configure identity provider integration

There are values your identity provider needs for successful integration with UCP. These values can vary between identity providers. Consult your identity provider documentation for instructions on providing these values as part of their integration process.

### Okta integration values

Okta integration requires these values:

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

ADFS integration requires these values:

- (need values)
-

### Ping Identity integration values

Ping Identity integration requires these values:

- (need values)

## Configure the SAML integration

To enable SAML authentication:

1 Go to the UCP web UI.
2. Navigate to the **Admin Settings**.
3. Select **Authentication & Authorization**.

![Enabling SAML in UCP](../../images/saml_enabled.png)

4. In the **SAML Enabled** section, select **Yes** to display the required settings.

![Configuring SAML in UCP](../../images/saml_settings.png)

5. In **IdP Metadata URL** enter the URL for the identity provider's metadata.
6. In **UCP Host** enter the URL that includes the IP address of your UCP console.
7. Select **Save** to complete the integration.

## Security considerations

You can download a client bundle to access UCP. To ensure that access from the client bundle is synced with the identity provider, we recommend the following steps. Otherwise, a previously-authorized user could get access to UCP through their existing client bundle.

- Remove the user account from UCP granting client bundle access if access is removed from the identity provider.
- If group membership in the identity provider changes, replicate this change in UCP.
- Continue to use LDAP to sync group membership.
