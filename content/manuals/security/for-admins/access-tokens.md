---
title: Organization access tokens
description: Learn how to create and manage organization access tokens
  to securely push and pull images programmatically.
keywords: docker hub, security, OAT, organization access token
linkTitle: Organization access tokens
params:
  sidebar:
    badge:
      color: blue
      text: Beta
---

{{< summary-bar feature_name="OATs" >}}

> [!WARNING]
>
> Organization access tokens (OATs) are not intended to be used with Docker Desktop, and are incompatible.
>
> OATs are also currently incompatible with the following services:
>
> - Docker Build Cloud
> - Docker Scout
> - Docker REST APIs
>
> If you use Docker Desktop or one of these services, you must use personal access tokens instead.

An organization access token (OAT) is like a [personal access token
(PAT)](/security/for-developers/access-tokens/), but an OAT is associated with
an organization and not a single user account. Use an OAT instead of a PAT to
let business-critical tasks access Docker Hub repositories without connecting
the token to single user. You must have a [Docker Team or Business
subscription](/subscription/core-subscription/details/) to use OATs.

OATs provide the following advantages:

- You can investigate when the OAT was last used and then disable or delete it
  if you find any suspicious activity.
- You can limit what each OAT has access to, which limits the impact if an OAT
  is compromised.
- All organization owners can manage OATs. If one owner leaves the organization,
  the remaining owners can still manage the OATs.
- OATs have their own Docker Hub usage limits that don't count towards your
  personal account's limits.

If you have existing [service accounts](/docker-hub/service-accounts/), Docker recommends that you replace the service accounts with OATs. OATs offer the following advantages over service accounts:

- Access permissions are easier to manage with OATs. You can assign access
  permissions to OATs, while service accounts require using teams for access
  permissions.
- OATs are easier to manage. OATs are centrally managed in the Admin Console.
  For service accounts, you may need to sign in to that service account to
  manage it. If using single sign-on enforcement and the service account is not
  in your IdP, you may not be able to sign in to the service account to manage
  it.
- OATs are not associated with a single user. If a user with access to the
  service account leaves your organization, you may lose access to the service
  account. OATs can be managed by any organization owner.

## Create an organization access token

> [!IMPORTANT]
>
> Treat access tokens like a password and keep them secret. Store your tokens securely in a credential manager for example.

Organization owners can create up to 10 organization access tokens (OATs) for
organizations with a Team subscription and up to 100 OATs for organizations with
a Business subscription. Expired tokens count towards the total amount of
tokens.

To create an OAT:

1. Sign in to the [Admin Console](https://app.docker.com/admin).

2. Select the organization you want to create an access token for.

3. Under **Security and access**, select **Access tokens**.

4. Select **Generate access token**.

5. Add a label and optional description for your token. Use something that indicates the use case or purpose of the token.

6. Select the expiration date for the token.

7. Select the repository access for the token.

   The access permissions are scopes that set restrictions in your repositories.
   For example, for Read & Write permissions, an automation pipeline can build
   an image and then push it to a repository. However, it can't delete the
   repository. You can select one of the following options:

   - **Public repositories (read only)**
   - **All repositories**: You can select read access, or read and write access.
   - **Select repositories**: You can select up to 50 repositories, and then
     select read access, or read and write access for each repository.

8. Select **Generate token** and then copy the token that appears on the screen
   and save it. You won't be able to retrieve the token once you exit the
   screen.

## Use an organization access token

You can use an organization access token when you sign in using Docker CLI.

Sign in from your Docker CLI client with the following command, replacing
`YOUR_ORG` with your organization name:

```console
$ docker login --username <YOUR_ORG>
```

When prompted for a password, enter your organization access token instead of a
password.

## Modify existing tokens

You can rename, update the description, update the repository access,
deactivate, or delete a token as needed.

1. Sign in to the [Admin Console](https://app.docker.com/admin).

2. Select the organization you want to modify an access token for.

3. Under **Security and access**, select **Access tokens**.

4. Select the actions menu on the far right of a token row, then select
   **Deactivate**, **Edit**, or **Delete** to modify the token. For **Inactive**
   tokens, you can only select **Delete**.

5. If editing a token, select **Save** after specifying your modifications.
