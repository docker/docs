---
title: Deploy your Docker setup
description: Deploy your Docker setup across your company.
weight: 40
---

> [!WARNING]
>
> Communicate with your users before proceeding, and confirm that your IT and
> MDM teams are prepared to handle any unexpected issues, as these steps will
> affect all existing users signing into your Docker organization.

## Enforce SSO

Enforcing SSO means that anyone who has a Docker profile with an email address
that matches your verified domain must sign in using your SSO connection. Make
sure the Identity provider groups associated with your SSO connection cover all
the developer groups that you want to have access to the Docker subscription.

For instructions on how to enforce SSO, see [Enforce SSO](/manuals/enterprise/security/single-sign-on/connect.md).

## Deploy configuration settings and enforce sign-in to users

Have the MDM team deploy the configuration files for Docker to all users.

## Next steps

Congratulations, you've successfully completed the admin implementation process
for Docker.

To continue optimizing your Docker environment:

- Review your [organization's usage data](/manuals/admin/organization/insights.md) to track adoption
- Monitor [Docker Scout findings](/manuals/scout/explore/analysis.md) for security insights
- Explore [additional security features](/manuals/enterprise/security/_index.md) to enhance your configuration
