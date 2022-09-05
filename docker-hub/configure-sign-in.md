---
description: Configure registry.json to enforce users to sign into Docker Desktop
keywords: authentication, registry.json, configure,
title: Configure registry.json to enforce sign-in
---

By default, members of your organization can use Docker Desktop on their machines without signing in to any Docker account. To ensure that a user signs in to a Docker account that is a member of your organization and that the
organization’s settings apply to the user’s session, you can use a `registry.json` file.

The `registry.json` file is a configuration file that allows administrators to specify the Docker organization the user must belong to and ensure that the organization’s settings apply to the user’s session. The Docker Desktop installer can create this file on the users’ machines as part of the installation process.

After a `registry.json` file is configured on a user’s machine, Docker Desktop prompts the user to sign in. If a user doesn’t sign in, or tries to sign in using a different organization, other than the organization listed in the `registry.json` file, they will be denied access to Docker Desktop.

Deploying a `registry.json` file and forcing users to authenticate is not required, but offers the following benefits:

 - Allows administrators to configure features such as [Image Access Management](image-access-management.md) which allows team members to:
    - Only have access to Trusted Content on Docker Hub
    - Pull only from the specified categories of images
- Authenticated users get a higher pull rate limit compared to anonymous users. For example, if you are authenticated, you get 200 pulls per 6 hour period, compared to 100 pulls per 6 hour period per IP address for anonymous users. For more information, see [Download rate limit](download-rate-limit.md).
- Blocks users from accessing Docker Desktop until they are added to a specific organization.

{% include configure-registry-json.md %}

## Verify the changes

After you’ve created the `registry.json` file and deployed it onto the users’ machines, you can verify whether the changes have taken effect by asking users to start Docker Desktop.

If the configuration is successful, Docker Desktop prompts the user to authenticate using the organization credentials on start. If the user fails to authenticate, they will see an error message, and they will be denied access to Docker Desktop.
