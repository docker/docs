{% if include.product == "admin" %}
  {% assign product_link="[Docker Admin](https://admin.docker.com)" %}
  {% assign domain_navigation="Select your organization in the left navigation drop-down menu, and then select **Domain management**." %}
  {% assign sso_link = "[SSO](/admin/organization/security-settings/sso/)" %}
  {% assign scim_link = "[SCIM](/admin/organization/security-settings/scim/)" %}
  {% assign invite_link = "[Invite members](/admin/organization/members/)" %}
{% else %}
  {% assign product_link="[Docker Hub](https://hub.docker.com)" %}
  {% assign domain_navigation="Select **Organizations**, your organization, **Settings**, and then **Security**." %}
  {% assign sso_link = "[SSO](/single-sign-on/)" %}
  {% assign scim_link = "[SCIM](/docker-hub/scim/)" %}
  {% assign invite_link = "[Invite members](/docker-hub/members/)" %}
{% endif %}

Domain audit identifies uncaptured users in an organization. Uncaptured users are Docker users who have authenticated to Docker using an email address associated with one of your verified domains, but they're not a member of your organization in Docker. You can audit domains on organizations that are part of the Docker Business subscription. To upgrade your existing account to a Docker Business subscription, see [Upgrade your subscription](/subscription/upgrade/).

Uncaptured users who access Docker Desktop in your environment may pose a security risk because your organization's security settings, like Image Access Management and Registry Access Management, aren't applied to a user's session. In addition, you won't have visibility into the activity of uncaptured users. You can add uncaptured users to your organization to gain visibility into their activity and apply your organization's security settings.

Domain audit can't identify the following Docker users in your environment:
   * Users who access Docker Desktop without authenticating
   * Users who authenticate using an account that doesn't have an email address associated with one of your verified domains

Although domain audit can't identify all Docker users in your environment, you can enforce sign-in to prevent unidentifiable users from accessing Docker Desktop in your environment. For more details about enforcing sign-in, see [Configure registry.json to enforce sign-in](/docker-hub/configure-sign-in/).

### Audit your domains for uncaptured users

Before you audit your domains, the following prerequisites are required:
   * Your organization must be part of a Docker Business subscription. To upgrade your existing account to a Docker Business subscription, see [Upgrade your subscription](/subscription/upgrade/).
   * You must add and verify your domains.

To audit your domains:

1. Sign in to {{ product_link }}{: target="_blank" rel="noopener" class="_"}.
2. {{ domain_navigation }}
3. In **Domain Audit**, select **Export Users** to export a CSV file of uncaptured users with the following columns:
  - Name: The name of the user.
  - Username: The Docker ID of the user.
  - Email: The email address of the user.

You can invite all the uncaptured users to your organization using the exported CSV file. For more details, see {{ invite_link }}. Optionally, enforce single sign-on or enable SCIM to add users to your organization automatically. For more details, see {{ sso_link }} or {{ scim_link }}.

> **Note**
>
> Domain audit may identify accounts of users who are no longer a part of your organization. If you don't want to add a user to your organization and you don't want the user to appear in future domain audits, you must deactivate the account or update the associated email address.
>
> Only someone with access to the Docker account can deactivate the account or update the associated email address. For more details, see [Deactivating an account](/docker-hub/deactivate-account/).
>
> If you don't have access to the account, you can contact [Docker support](/support/) to discover if more options are available.