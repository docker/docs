---
description: Learn how to set up SSO in your IdP and take the next steps for enabling SSO.
keywords: configure, sso, docker hub, hub, docker admin, admin, security
title: Configure your IdP
---

This page describes the general steps to continue setting up your single-sign on (SSO) connection.

## Prerequisites

Make sure you have completed the following before you begin:

- Your domain is verified
- You have created your SSO connection in Docker
- You have copied the necessary fields from Docker to paste in your IdP:
    - SAML: **Entity ID**, **ACS URL**
    - Azure AD (OIDC): **Redirect URL**

## Configure your IdP to work with Docker

{{< tabs >}}
{{< tab name="Okta" >}}

1. Go to the Okta admin portal.
2. Go to **Applications > Applications > Create App Integration**.
3. Select **SAML 2.0**, then select **Next**.
4. Enter App Name "Docker Hub" and optionally upload a logo for the app, then select **Next**.
5. To configure SAML, enter the following into Okta:
    - ACS URL - Single Sign On URL
    - Entity ID - Audience URI (SP Entity ID)
    - Name ID format: `EmailAddress`
    - Application username: `Email`
    - Update application on: `Create or Update`
    - Attribute Statements: `add`. You can define your attribute statement like the following: 

     | Attribute name | Name format | Value                                    |
     | :------------- | :---------- | :--------------------------------------- |
     | name           | Unspecified | username.firstName + " " + user.lastName |

6. Select **Next**.
7. Select **I'm an Okta customer adding an internal app**.
8. Select **Finish**.

{{< /tab >}}
{{< tab name="Entra ID SAML 2.0" >}}

1. Go to Azure AD admin portal.
2. Go to **Default Directory > Add > Enterprise Application > Create your own application**.
3. Enter “Docker” for application name and select **non-gallery** option.
4. After the application is created, go to Single Sign-On and select **SAML**.
5. Select **Edit** on the **Basic SAML configuration** section.
6. Add the following settings from Docker Hub:
    - Entity ID → Identifier
    - ACS URL → Reply URL
7. Save configuration
8. From section **SAML Signing Certificate** download **Certificate (Base64)**
9. Open the certificate file in a text editor and paste the contents of the file in the **x509 Certificate** field in Docker Hub or Admin Console.
10. From the section **Set up Docker**, copy **Login URL** and paste it into the **SAML Sign-in URL** field in Docker Hub or Admin Console.

{{< /tab >}}
{{< tab name="Azure Connect (OIDC)" >}}

### Create app registration

1. Go to Azure AD admin portal.
2. Select **App Registration > New Registration**.
3. Name the App to “Docker Hub SSO” or pick any name you wish for the app.
4. Under **Supported account types**, specify who can use this application or access the app.
5. In the **Redirect URI** section, select **Web** from the dropdown menu and paste the **Redirect URI** value from the Docker console into this field.
6. Select **Register** to register the app.
7. Take note of the **Client ID** from the app's overview page. You need this information to continue configuring SSO on Docker Hub.

### Create client secrets for your Docker app

1. Go to the Docker Hub SSO app that you created in the previous steps, then select **Certificates & secrets**.
2. Select **+ New client secret**.
3. Specify the description of the secret and set how long the keys can be used on Azure.
4. Select **Add** to continue.
5. Copy the secret **Value** field and keep it somewhere safe so you can use it to configure Docker SSO later on.

### Configure API permission for Docker SSO and grant admin consent

1. Go to the Docker Hub SSO app that you created in the previous steps.
2. Navigate to the **API permission** category in your app settings.
3. Select **Grant admin consent for YOUR TENANT NAME > Yes**.

### Assign users to the SSO app

1. Navigate to your Azure AD dashboard, then select **Enterprise Applications > APP NAME**.
2. Select **1. Assign users and groups**.
3. Add users that will be allowed to use the app.

In the Docker Console, paste the following values obtained in the previous steps to continue configuration:

- **Client ID**
- **Client Secret**
- **Azure AD Domain**

{{< /tab >}}
{{< tab name="OneLogin" >}}

1. Go to the OneLogin admin portal.
2. Go to **Applications > Applications > Add App**
3. Use the search input and search for “**SCIM Provisioner with SAML (SCIM v2 Core)**” and select the only result item.
4. Enter a display name, for example “Docker Hub”.
5. Optional: Upload Docker icons for the app and add an appropriate description for the app.
6. Select **Save**.
7. After saving the new app, more tabs on the left will appear. Go to the **Configuration** tab.
8. Open a separate tab in your browser and go to **Docker Hub > Settings > Security > SSO > SAML**.
9. Copy from the Docker console and paste the following into OneLogin:
   - Entity ID - SAML Audience URL
   - ACS URL - SAML Consumer URL
   - SCIM Base URL - SCIM Base URL
   - Custom Headers:
     - `Content-Type: application/scim+json`
     - `User-Agent: OneLogin SCIM`
   - SCIM Bearer Token - SCIM Bearer Token
   - SCIM JSON Template:

     ```json
        {
            "schemas": [
            "urn:ietf:params:scim:schemas:core:2.0:User"
            ],
            "userName": "{$parameters.scimusername}",
            "name": {
                "givenName": "{$user.firstname}",
                "familyName": "{$user.lastname}"
            },
            "emails": [
                {
                    "value": "{$user.email}",
                    "primary": true
                }
            ]
        }
     ```
10. Select **API Connection > Enable**, to change **API Status** to **Enable**.
11. Select **Save** button again, then go to the **Parameters** tab.
12. Select **scimusername** and set the value to `Email`.
13. Select **Save**, then go to the **SSO** tab.
14. Copy **SAML 2.0 Endpoint (HTTP)** url to paste in the Docker console in the **SAML Sign-on URL** field.
15. Go to **X.509 Certificate** and select **View Details** to copy the PEM certificate (—-BEGIN CERTIFICATE —- ….) to paste in the Docker console in the **x509 Certificate** field.

{{< /tab >}}
{{< /tabs >}}

## What's next?

[Complete your connection](../connect/_index.md) in the Docker console, then test your connection.
