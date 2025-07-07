{{ $product_link := "[Docker Hub](https://hub.docker.com)" }}
{{ $sso_navigation := `Navigate to the SSO settings page for your organization. Select **My Hub**, your organization, **Settings**, and then **Security**.` }}

{{ if eq (.Get "product") "admin" }}
  {{ $product_link = "[Docker Home](https://app.docker.com) and select your organization. Note that when an organization is part of a company, you must select the company and configure SSO for that organization at the company level. Each organization can have its own SSO configuration and domain, but it must be configured at the company level." }}
  {{ $sso_navigation = "Select **Admin Console**, then **SSO and SCIM**." }}
{{ end }}

### Remove a domain from an SSO connection

> [!IMPORTANT]
>
> Docker supports multiple IdP configurations, where a single domain is used for multiple SSO identity providers. If you want to remove a domain from multiple SSO connections, you must remove it from each connection individually.

1. Sign in to {{ $product_link }}.
2. {{ $sso_navigation }}
3. In the SSO connections table, select the **Action** icon and then **Edit connection**.
4. Select **Next** to navigate to the section where the connected domains are listed.
5. In the **Domain** drop-down, select the **x** icon next to the domain that you want to remove.
6. Select **Next** to confirm or change the connected organization(s).
7. Select **Next** to confirm or change the default organization and team provisioning selections.
8. Review the **Connection Summary** and select **Update connection**.

> [!NOTE]
>
> If you want to re-add the domain, a new TXT record value is assigned. You must then complete the verification steps with the new TXT record value.
