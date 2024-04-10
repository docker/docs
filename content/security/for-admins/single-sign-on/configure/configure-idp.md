---
description: Learn how to set up SSO in your IdP and take the next steps for enabling SSO.
keywords: configure, sso, docker hub, hub, docker admin, admin, security
title: Configure your IdP
---

The steps to set up your SSO configuration are:

1. [Add and verify the domain or domains](/security/for-admins/single-sign-on/configure#step-one-add-and-verify-your-domain) that your members use to sign in to Docker.
2. [Create your SSO connection](/security/for-admins/single-sign-on/configure#step-two-create-an-sso-connection-in-docker) in Docker.
3. [Configure your IdP](#step-three-configure-your-idp-to-work-with-docker) to work with Docker.
4. [Complete your SSO connection](/security/for-admins/single-sign-on/connect/) in Docker.

This page walks through step 3 for common IdPs.

## Prerequisites

Make sure you have completed the following before you begin:

- Your domain is verified
- You have created your SSO connection in Docker
- You have copied the necessary fields from Docker to paste in your IdP:
    - SAML: **Entity ID**, **ACS URL**
    - Azure AD (OIDC): **Redirect URL**

## SSO attributes

When a user signs in using SSO, Docker obtains the following attributes from the IdP:

- **Email address** - unique identifier of the user
- **Full name** - name of the user
- **Groups (optional)** - list of groups to which the user belongs
- **Docker Org (optional)** - the organization to which the user belongs
- **Docker Team (optional)** - the team within an organization that a user has been added to
- **Docker Role (optional)** - the role for the user that grants their permissions in an organization

If you use SAML for your SSO connection, Docker obtains these attributes from the SAML assertion message. Your IdP may use different naming for SAML attributes than those in the previous list. The following table lists the possible SAML attributes that can be present in order for your SSO connection to work.

> **Important**
>
>SSO uses Just-in-Time (JIT) provisioning by default. If you [enable SCIM](../../scim.md), JIT values still overwrite the attribute values set by SCIM provisioning whenever users log in. To avoid conflicts, make sure your JIT values match your SCIM values. For example, to make sure that the full name of a user displays in your organization, you would set a `name` attribute in your SAML attributes and ensure the value includes their first name and last name. The exact method for setting these values (for example, constructing it with `user.firstName + " " + user.lastName`) varies depending on your IdP.
{.important}

> **Beta feature**
>
> Optional Just-in-Time (JIT) provisioning is available in Private Beta when you use the Admin Console. If you're participating in this program, you can avoid conflicts between SCIM and JIT by disabling JIT provisioning in your SSO connection. See [SSO authentication with JIT provisioning disabled](/security/for-admins/group-mapping/#sso-authentication-with-jit-provisioning-disabled).
{ .experimental }

You can also configure attributes to override default values, such as default team or organization. See [role mapping](../../scim.md#set-up-role-mapping).

| SSO attribute | SAML assertion message attributes |
| ---------------- | ------------------------- |
| Email address    | `"http://schemas.xmlsoap.org/ws/2005/05/identity/claims/nameidentifier"`, `"http://schemas.xmlsoap.org/ws/2005/05/identity/claims/upn"`, `"http://schemas.xmlsoap.org/ws/2005/05/identity/claims/emailaddress"`, `email`                           |
| Full name        | `"http://schemas.xmlsoap.org/ws/2005/05/identity/claims/name"`, `name`, `"http://schemas.xmlsoap.org/ws/2005/05/identity/claims/givenname"`, `"http://schemas.xmlsoap.org/ws/2005/05/identity/claims/surname"`  |
| Groups (optional) | `"http://schemas.xmlsoap.org/claims/Group"`, `"http://schemas.microsoft.com/ws/2008/06/identity/claims/groups"`, `Groups`, `groups` |
| Docker Org (optional)        | `dockerOrg`   |
| Docker Team (optional)     | `dockerTeam`  |
| Docker Role (optional)      | `dockerRole`  |

> **Important**
>
> If none of the email address attributes listed in the previous table are found, SSO returns an error. Also, if the `Full name` attribute isn't set, then the name will be displayed as the value of the `Email address`.
{ .important}

## Step three: Configure your IdP to work with Docker

The user interface for your IdP may differ slightly from the following steps. You can refer to the documentation for your IdP to verify.

{{< tabs >}}
{{< tab name="Okta" >}}

1. Go to the Okta admin portal.
2. Go to **Applications > Applications > Create App Integration**.
3. Select **SAML 2.0**, then select **Next**.
4. Enter App Name "Docker Hub" and optionally upload a logo for the app, then select **Next**.
5. To configure SAML, enter the following into Okta:
    - ACS URL: Single Sign On URL
    - Entity ID: Audience URI (SP Entity ID)
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
9. After you create the app, go to your app and select **View SAML setup instructions**.
10. Here you can find the **SAML Sign-in URL** and the **x509 Certificate**. Open the certificate file in a text editor and paste the contents of the file in the **x509 Certificate** field in Docker Hub or Admin Console. Then, paste the value of the **SAML Sign-in URL** and paste it into the corresponding field in Docker Hub or Admin Console.

{{< /tab >}}
{{< tab name="Entra ID SAML 2.0" >}}

> **Tip**
>
> When you create the application for your SSO connection in Entra ID (formerly Azure AD) we recommend that you don't assign the app to all the users in the directory.
> Instead, you can create a security group and assign the app to the group. This way, you can control who in your organization has access to Docker.
> To change the default setting for assignment, go to the main properties for your app and find the **Assignment required** setting. Set it to **Yes**.
{ .tip }

See [More resources](#more-resources) for a video overview on how to set up SSO with SAML in Entra ID (formerly Azure AD).

1. Go to Azure AD admin portal.
2. Go to **Default Directory > Add > Enterprise Application > Create your own application**.
3. Enter “Docker” for application name and select **non-gallery** option.
4. After the application is created, go to **Single Sign-On** and select **SAML**.
5. Select **Edit** on the **Basic SAML configuration** section.
6. Add the following settings from Docker Hub:
    - Entity ID: Identifier
    - ACS URL: Reply URL
7. Save configuration.
8. From section **SAML Signing Certificate** download **Certificate (Base64)**.
9. Open the certificate file in a text editor and paste the contents of the file in the **x509 Certificate** field in Docker Hub or Admin Console.
10. From the section **Set up Docker**, copy **Login URL** and paste it into the **SAML Sign-in URL** field in Docker Hub or Admin Console.

{{< /tab >}}
{{< tab name="Azure Connect (OIDC)" >}}

### Create app registration

1. Go to Azure AD admin portal.
2. Select **App Registration > New Registration**.
3. Enter “Docker Hub SSO” or similar for the application name.
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
4. Next, you need to add additional permissions. Select **Add a permission**.
5. Select **Delegated permissions**.
6. Search for `Directory.Read.All`, and select this option.
7. Then, search for `User.Read`, and select this option.
8. Select **Add permissions**.

You can verify admin consent was granted for each permission correctly by checking the **Status** column.

### Assign users to the SSO app

1. Navigate to your Azure AD dashboard, then select **Enterprise Applications > APP NAME**.
2. Select **1. Assign users and groups**.
3. Add users that will be allowed to use the app.

In the Docker console, paste the following values obtained in the previous steps to continue configuration:

- **Client ID**
- **Client Secret**
- **Azure AD Domain**

{{< /tab >}}
{{< /tabs >}}

## More resources

The following videos demonstrate how to configure your IdP with your Docker SSO connection.

- [Video: SSO connection with Okta](https://youtu.be/c56YECO4YP4?feature=shared&t=633)
- [Video: SSO connection with Azure Connect (OIDC)](https://youtu.be/bGquA8qR9jU?feature=shared&t=630)
- [Video: SSO connection with Entra ID (Azure) SAML](https://youtu.be/bGquA8qR9jU?feature=shared&t=1246)

## What's next?

[Complete your connection](../connect/_index.md) in the Docker console, then test your connection.
