{{ $product_link := "[Docker Hub](https://hub.docker.com)" }}
{{ $sso_navigation := `Navigate to the SSO settings page for your organization or company.
   - Organization: Select **Organizations**, your organization, **Settings**, and then **Security**.
   - Company: Select **Organizations**, your company, and then **Settings**.` }}

{{ if eq (.Get "product") "admin" }}
  {{ $product_link = "the [Admin Console](https://admin.docker.com)" }}
  {{ $sso_navigation = "Select your organization or company in the left navigation drop-down menu, and then select **SSO & SCIM.**" }}
{{ end }}

> **Important**
>
> If your IdP setup requires an Entity ID and the ACS URL, you must select the
> **SAML** tab in the **Authentication Method** section. For example, if your
> Entra ID (formerly Azure AD) Open ID Connect (OIDC) setup uses SAML configuration within Azure
> AD, you must select **SAML**. If you are [configuring Open ID Connect with Entra ID (formerly Azure AD)](https://docs.microsoft.com/en-us/powerapps/maker/portals/configure/configure-openid-settings) select
> **Azure AD (OIDC)** as the authentication method. Also, IdP initiated connections
> aren't supported at this time.
{ .important}

After your domain is verified, create an SSO connection.

1. Sign in to {{ $product_link }}.
2. {{ $sso_navigation }}
3. In the SSO connections table select **Create Connection**, and create a name for the connection.

   > **Note**
   >
   > You have to verify at least one domain before creating the connections.

4. Select an authentication method, **SAML** or **Azure AD (OIDC)**.
5. Copy the following fields to add to your IdP:

   - SAML: **Entity ID**, **ACS URL**
   - Azure AD (OIDC): **Redirect URL**

   ![SAML](/docker-hub/images/saml-create-connection.png)

   ![Azure AD](/docker-hub/images/azure-create-connection.png)
