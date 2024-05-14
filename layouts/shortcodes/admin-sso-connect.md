{{ $product_link := "[Docker Hub](https://hub.docker.com)" }}
{{ $sso_navigation := `Navigate to the SSO settings page for your organization or company.
   - Organization: Select **Organizations**, your organization, **Settings**, and then **Security**.
   - Company: Select **Organizations**, your company, and then **Settings**.` }}

{{ if eq (.Get "product") "admin" }}
  {{ $product_link = "the [Admin Console](https://admin.docker.com)" }}
  {{ $sso_navigation = "Select your organization or company in the left navigation drop-down menu, and then select **SSO & SCIM.**" }}
{{ end }}

1. In {{ $product_link }}, select the verified domains you want to apply the connection to.
2. To provision your users, select the organization(s) and/or team(s).
3. Review your summary and select **Create Connection**.

## Test your SSO configuration

After you’ve completed the SSO configuration process in Docker, you can test the configuration when you sign in to {{ $product_link }} using an incognito browser. Sign in to {{ $product_link }} using your domain email address. You are then redirected to your IdP's login page to authenticate.

1. Authenticate through email instead of using your Docker ID, and test the login process.
2. To authenticate through CLI, your users must have a PAT before you enforce SSO for CLI users.

>**Important**
>
> SSO has Just-in-Time (JIT) provisioning enabled by default, unless you have [disabled it](/security/for-admins/provisioning/just-in-time/#sso-authentication-with-jit-provisioning-disabled). This means your users are auto-provisioned to your organization on Docker Hub.
>
> You can change this on a per-app basis. To prevent auto-provisioning users, you can create a security group in your IdP and configure the SSO app to authenticate and authorize only those users that are in the security group. Follow the instructions provided by your IdP:
>
> - [Okta](https://help.okta.com/en-us/Content/Topics/Security/policies/configure-app-signon-policies.htm)
> - [Entra ID (formerly Azure AD)](https://learn.microsoft.com/en-us/azure/active-directory/develop/howto-restrict-your-app-to-a-set-of-users)
{ .important}

The SSO connection is now created. You can continue to set up SCIM without enforcing SSO log-in. For more information about setting up SCIM, see [Set up SCIM](/security/for-admins/provisioning/scim/).

## Optional: Enforce SSO

1. Sign in to {{ $product_link }}.
2. {{ $sso_navigation }}
3. In the SSO connections table, select the **Action** icon and then **Enable enforcement**.

   When SSO is enforced, your users are unable to modify their email address and password, convert a user account to an organization, or set up 2FA through Docker Hub. You must enable 2FA through your IdP.

4. Continue with the on-screen instructions and verify that you’ve completed the tasks.
5. Select **Turn on enforcement** to complete.

Your users must now sign in to Docker with SSO.

> **Important**
>
> If SSO isn't enforced, users can choose to sign in with either their Docker ID or SSO.
{ .important}
