---
title: GitHub app
description:
keywords:
---

{% include atomist/disclaimer.md %}

When installed for a GitHub organization, the Atomist GitHub app links
repository activity to container images.

Install the GitHub app in the organization that contains the source code
repositories for your Docker images.

## Connect to GitHub

1. Go to <https://dso.docker.com/> and sign in using your Docker ID.
2. Open the **Repositories** tab.
3. Click **Connect to GitHub** and follow the authorization flow. You will be
   installing the
   [Atomist GitHub App](https://github.com/apps/atomist "Atomist GitHub App").

   ![install](images/gh-install.png){: width="700px" }

4. Install the app.

   > If your GitHub account is a member of one or more organizations, you will
   > be asked to choose which account to install the app into.

Once you are redirected back to Atomist, you'll be asked to select a set of
repositories that are used to create container images. Atomist will start
watching for images built from the repositories you select here.

![activate-repos](images/activate-repos.png){: width="700px" }

If you are just looking to evaluate Atomist, we recommend that you start by
selecting a few repositories during evaluation. Then, once you are comfortable
using Atomist, you may want to switch on the integration for all repositories.
Selecting **All repositories** also includes any repository created in the
future.

> Note: if Atomist detects `FROM` commands in Dockerfiles in the selected
> repositories, it will begin raising automated pull requests for version
> pinning the image versions in the `FROM` command.

## Manage repository access

If you wish to add or remove repository access for Atomist, go to the
[repositories page](https://dso.docker.com/r/auth/repositories){: target="blank"
rel="noopener" class=""}.

- Select **All repositories** if you want Atomist to be enabled for all
  connected organizations and repositories.
- Select **Only select repositories** if you want to provision access to only a
  subset of repositories.

  Check or un-check organizations or repositories in the list. You can filter
  the list by repository name using the text filter.

## Disconnecting GitHub

You might want to disconnect GitHub when:

- You want to change the GitHub organization or account that is connected to
  your Atomist workspace. To do so, disconnect the old GitHub organization or
  account first, then follow the instructions above for Connecting to GitHub to
  connect the new GitHub organization or account.
- You want to remove Atomist access to a GitHub organization or account when you
  no longer need Atomist to have access.

To disconnect a GitHub account:

1.  Visit **Repositories** and click the **Disconnect** link. This removes the
    connection to your GitHub organization or account.
2.  Go to the
    [GitHub Applications settings page](https://github.com/settings/installations){:
    target="blank" rel="noopener" class=""}, then
    - Find **atomist** on the Installed GitHub Apps tab, select **Configure**,
      then **Uninstall**. This removes the installation of the Atomist GitHub
      App from your GitHub organization or account.
    - Find **atomist** on the Authorized GitHub Apps tab, select **Revoke**.
      This removes the authorization of the Atomist GitHub App from your GitHub
      organization or account.
