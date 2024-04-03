{{ $product_link := "[Docker Hub](https://hub.docker.com)" }}
{{ $sso_navigation := `Navigate to the SSO settings page for your organization or company.
   - Organization: Select **Organizations**, your organization, **Settings**, and then **Security**.
   - Company: Select **Organizations**, your company, and then **Settings**.` }}
{{ $member_navigation := "Select **Organizations**, your organization, and then **Members**." }}
{{ $invite_button := "**Invite members**" }}
{{ $remove_button := "**Remove member**" }}

{{ if eq (.Get "product") "admin" }}
  {{ $product_link = "the [Admin Console](https://admin.docker.com)" }}
  {{ $invite_button = "**Invite**" }}
  {{ $sso_navigation = "Select your organization or company in the left navigation drop-down menu, and then select **SSO & SCIM**." }}
  {{ $member_navigation := `Navigate to the user management page for your organization or company. 
    - Organization: Select your organization in the left navigation drop-down menu, and then select **Members**.
    - Company: Select your company in the left navigation drop-down menu, and then select **Users**.` }}
  {{ $remove_button = "**Remove member**, if you're an organization, or **Remove user**, is you're a company" }}
{{ end }}

### Remove a domain from an SSO connection

1. Sign in to {{ $product_link }}.
2. {{ $sso_navigation }}
3. In the SSO connections table, select the **Action** icon and then **Edit connection**.
4. Select **Next** to navigate to the section where the connected domains are listed.
5. In the **Domain** drop-down, select the **x** icon next to the domain that you want to remove.
6. Select **Next** to confirm or change the connected organization(s).
7. Select **Next** to confirm or change the default organization and team provisioning selections.
8. Review the **Connection Summary** and select **Save**.

> **Note**
>
> If you want to re-add the domain, a new TXT record value is assigned. You must then complete the verification steps with the new TXT record value.

## Manage SSO connections

### Edit a connection

1. Sign in to {{ $product_link }}.
2. {{ $sso_navigation }}
3. In the SSO connections table, select the **Action** icon.
4. Select **Edit connection** to edit your connection.
5. Follow the on-screen instructions to edit the connection.

### Delete a connection

1. Sign in to {{ $product_link }}.
2. {{ $sso_navigation }}
3. In the SSO connections table, select the **Action** icon.
4. Select **Delete connection**.
5. Follow the on-screen instructions to delete a connection.

### Deleting SSO

When you disable SSO, you can delete the connection to remove the configuration settings and the added domains. Once you delete this connection, it can't be undone. Users must authenticate with their Docker ID and password or create a password reset if they don't have one.

## Manage users

> **Important**
>
> SSO has Just-In-Time (JIT) Provisioning enabled by default. This means your users are auto-provisioned to your organization.
>
> You can change this on a per-app basis. To prevent auto-provisioning users, you can create a security group in your IdP and configure the SSO app to authenticate and authorize only those users that are in the security group. Follow the instructions provided by your IdP:
>
> - [Okta](https://help.okta.com/en-us/Content/Topics/Security/policies/configure-app-signon-policies.htm)
> - [Entra ID (formerly Azure AD)](https://learn.microsoft.com/en-us/azure/active-directory/develop/howto-restrict-your-app-to-a-set-of-users)
{ .important}

> **Beta feature**
>
> Optional Just-in-Time (JIT) provisioning is available in private beta. If you're participating in this program, you have the option to turn off this default provisioning and disable JIT. This configuration is recommended if you're using SCIM to auto-provision users.
{ .experimental }

### Add guest users when SSO is enabled

To add a guest if they aren’t verified through your IdP:

1. Sign in to {{ $product_link }}.
2. {{ $member_navigation }}
3. Select {{ $invite_button }}.
4. Follow the on-screen instructions to invite the user.

### Remove users from the SSO company

To remove a user:

1. Sign in to {{ $product_link }}.
2. {{ $member_navigation }}
3. Select the action icon next to a user’s name, and then select {{ $remove_button }}.
4. Follow the on-screen instructions to remove the user.
