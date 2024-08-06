---
description: Learn how to troubleshoot common SSO issues.
keywords: sso, troubleshoot, single sign-on
title: Troubleshoot single sign-on
---

While configuring or using single sign-on (SSO), you may encounter issues that
can stem from your identity provider (IdP) or Docker configuration. The
following sections describe how to view the error messages in the Docker Admin
Console as well as some common errors and possible solutions. You can also see
your identity provider's documentation to learn if you can view error logs in
their service.

## View SSO and SCIM error logs

1. Sign in to the [Admin Console](https://app.docker.com/admin/).
2. Select your organization or company in the left navigation drop-down menu,
   and then select **SSO and SCIM**.
3. In the SSO connections table, select the **Actions** icon and **View error
   logs**. The **Connection errors** page appears with a list of errors that
   have occurred in the past 7 days.
4. In the **Connection errors** page, select **View error details** next to an
   error message for more details. A modal appears with a JSON object containing
   more details.

## Common SSO errors and solutions

[View the SSO and SCIM error logs](#view-sso-and-scim-error-logs) and then use
the following sections for solutions to common configuration errors.

### IdP-initiated login is not enabled for connection

An error message, similar to the following, appears in the error logs for this
issue.

```text
IdP-Initiated login is not enabled for connection "$ssoConnection".
```

Docker doesn't support an IdP-initiated SAML flow. This error can occur when a
user attempts to authenticate from the IdP, for example using the Docker SSO App
tile on the dashboard.

Possible solutions:

 * The user must initiate authentication from Docker apps (Hub, Desktop, etc).
   The user needs to enter their email address and they will get redirected to
   the configured SSO IdP for their domain.
 * (Optional) Configure the Docker SSO App as not visible to users on your IdP
   so users don’t attempt to start authentication from there.

### Not enough seats in organization

An error message, similar to the following, appears in the error logs for this
issue.

```text
Not enough seats in organization '$orgName'. Please add more seats. Please contact your company administrator. TraceID: XXXXXXXXXXXXXX
```

This error can occur when attempting to provision a user into the organization
via SSO Just-in-Time provisioning or SCIM, but the organization has no available
seats for the user.

Possible solutions:

 * Add more Docker Business subscription seats to the organization. For details,
   see [Add seats to your
  subscription](/subscription/core-subscription/add-seats/).
 * Remove some users or pending invitations from your organization to make more
   seats available. For more details, see [Manage organization
   members](/admin/organization/members/).

### Domain is not verified for SSO connection

An error message, similar to the following, appears in the error logs for this
issue.

```text
Domain '$emailDomain' is not verified for your SSO connection. Please contact your company administrator. TraceID: XXXXXXXXXXXXXX
```

This error occurs if the IdP authenticated a user through SSO and the UPN
returned to Docker doesn’t match any of the verified domains associated to the
SSO connection configured in Docker.

Possible solutions:

 * Make sure the IdP SSO connection is returning the correct UPN value as part
   of the assertion attributes (attributes mapping).
 * Add and verify all domains and subdomains that are used as UPN by your IdP
   and associate them to your Docker SSO connection. For more details, see [Add
   and verify your
   domain](/security/for-admins/single-sign-on/configure/#step-one-add-and-verify-your-domain).

### Back button pressed

An error message, similar to the following, appears in the error logs for this
issue.

```text
You may have pressed the back button, refreshed during login, opened too many login dialogs, or there is some issue with cookies, since we couldn't find your session. Try logging in again from the application and if the problem persists please contact the administrator.
```

This error typically occurs during the authentication flow when a user presses
the back or the refresh button on the browser. This causes the sign-in flow to
lose track of the initial authentication request, which is required to complete
all authentication flows.

Possible solutions:

 * Avoid pressing the back or refresh buttons during sign in.
 * Close the browser’s tab and start the authentication flow from the beginning
   in the app (Docker Desktop, Hub, etc.)

### User is not assigned to the organization

An error message, similar to the following, appears in the error logs for this
issue.

```text
User '$username' is not assigned to this SSO organization. Please contact your company administrator. TraceID: XXXXXXXXXXXXX
```

This error occurs if SSO Just-In-Time (JIT) provisioning is disabled. JIT
provisioning ensures that a user is added to an organization after they
authenticate via SSO. JIT provisioning can be optionally disabled to prevent
users taking seats in the organization automatically or when SCIM is used as
the only option for user provisioning.

Possible solutions:

 * Review your SSO connection configuration and enable JIT provisioning to add
   all users to the organization after authenticating via SSO. For more details,
   see [Just-in-Time
   provisioning](/security/for-admins/provisioning/just-in-time/).
 * If JIT provisioning should remain disabled, then add the user to the
   organization by manually inviting them. Next time the user authenticates via
   SSO they will get added to the org because they are invited. For more
   details, see [Manage organization members](/admin/organization/members/).
 * If SCIM should provision the user, then ensure that the IdP controlling SCIM
   provisioning is properly configured to synchronize users with Docker as soon
   as they get assigned to the app. For more details, refer to your identity
   provider's documentation.

### Name ID is not an email address

An error message, similar to the following, appears in the error logs for this
issue.

```text
The name ID sent by the identity provider is not an email address. Please contact your company administrator.
```

This error can occur during SAML authentication, when your IdP sends back a Name
ID (UPN) that doesn't comply with the email address format required. The Docker
SSO app requires a name identifier to be the primary email address of the user.

Possible solutions:

 * Ensure that the Name ID attribute format is `EmailAddress`.
