---
description: Configure registry.json to enforce users to sign into Docker Desktop
keywords: authentication, registry.json, configure,
title: Configure registry.json to enforce sign in
---

The `registry.json` file is a configuration file that allows administrators to
specify the Docker organization the user must belong to, and thereby ensure
that the organization's settings are applied to the user's session. Docker
Desktop installation requires admin access. In large enterprises where admin
access is restricted, administrators can create a `registry.json` file and
deploy it to the users’ machines using a device management software as part of
the Docker Desktop installation process.

After you deploy a `registry.json` file to a user’s machine, it prompts the user to sign into Docker Desktop. If a user doesn’t sign in, or tries to sign in using a different organization, other than the organization listed in the `registry.json` file, they will be denied access to Docker Desktop.
Deploying a `registry.json` file and forcing users to authenticate offers the following benefits:

1. Allows administrators to configure features such as [Image Access Management](image-access-management.md) which allows team members to:
    - Only have access to Trusted Content on Docker Hub
    - Pull only from the specified categories of images
2. Authenticated users get a higher pull rate limit compared to anonymous users. For example, if you are authenticated, you get 200 pulls per 6 hour period, compared to 100 pulls per 6 hour period per IP address for anonymous users. For more information, see [Download rate limit](download-rate-limit.md).
3. Blocks users from accessing Docker Desktop until they are added to a specific organization.

{% include configure-registry-json.md %}

## Verify the changes

After you’ve created the `registry.json` file and deployed it onto the users’ machines, you can verify whether the changes have taken effect by asking users to start Docker Desktop.

If the configuration is successful, Docker Desktop prompts the user to authenticate using the organization credentials on start. If the user fails to authenticate, they will see an error message, and they will be denied access to Docker Desktop.
