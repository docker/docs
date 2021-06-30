---
title: Managing access tokens
description: Learn how to create and manage your personal Docker Hub access tokens to securely push and pull images programmatically.
keywords: docker hub, hub, security, PAT, personal access token
---

Docker Hub lets you create personal access tokens as alternatives to your password. You can use tokens to access Hub images from the Docker CLI.

Using personal access tokens provides some advantages over a password:
* You can investigate when an access token was used last, and disable or delete it if you find any suspicious activity.
* When logged in with an access token, you can't perform any admin activity on the account, including changing the password.

Access tokens are also useful in building integrations, since you can issue
multiple tokens &ndash; one for each integration &ndash; and revoke them at
any time.

> Note: If you have [two-factor authentication (2FA)](/docker-hub/2fa) enabled on your account, you must create at least one personal access token. Otherwise, you will be unable to log in to your account from the Docker CLI.
{: .important }

## Create an access token

You can create as many tokens as you need.

1. Log in to [hub.docker.com](https://hub.docker.com).

2. Click on your username in the top right corner and select **Account
Settings**.

3. Select **Security > New Access Token**.

4. Add a description for your token. Use something that indicates where
the token is going to be used, or set a purpose for the token.

      ![](images/PAT-Menu.png)

5. You can view the following permissions from the drop-down:
 - Read, Write, Delete: lets you to manage team members. This is an admin privilege.
 - Read & Write: let you to push images to any repository managed by your account.
 - Read Only: lets you to view, search, and pull images from any public repositories and any private repositories that you have access to.
 - Public Repo Read-only: lets you view, search, and pull images from any public repositories.

 >Note:
 All users have acccess to Read, Write, and Delete. Only paid users have access to Read & Write, Read-only, and Public Repo Read-only.

6. Copy the token that appears on the screen. Make sure you do this now:
once you close this prompt, Docker will never show the token again.

      ![](images/hub-copy-token.png)

      Treat access tokens like your password and keep them secret. Store your tokens securely (for example, in a credential manager).


## Modify existing tokens

You can rename, deactivate, or delete a token as needed.

1. Access your tokens under **Account Settings > Security**.

2. Select a token and click **Delete** or **Edit**, or use the menu on
the far right of a token row to bring up the edit screen. You can also
select multiple tokens to delete them all once.

      ![](images/Generate-PAT.png)


## Use an access token

You can use an access token anywhere that requires your Docker Hub
password.

When logging in from your Docker CLI client (`docker login --username <username>`),
omit the password in the login command. When you're prompted for
a password, enter your token instead.

If you have 2FA enabled, you must use a personal access token when logging in
from the Docker CLI. If you don't have it enabled, this is an optional (but
more secure) method of authentication.