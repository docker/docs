---
description: Docker Hub Automated Builds with GitHub
keywords: Docker, docker, registry, accounts, plans, Dockerfile, Docker Hub, docs, documentation, trusted, builds, trusted builds,  automated builds, GitHub
title: Configure automated builds from GitHub
---

If you've previously linked Docker Hub to your GitHub account,
skip to [Creating an Automated Build](github.md#creating-an-automated-build).

## Linking Docker Hub to a GitHub account

> Automated Build Permissions

> Automated Builds currently require *read* and *write* access since
> [Docker Hub](https://hub.docker.com) needs to set up a GitHub webhook.
> We have no choice here &ndash; this is how GitHub manages permissions.
> We do guarantee nothing else is touched in your account.

To set up an Automated Build of a repository on GitHub, you need to
link [Docker Hub](https://hub.docker.com/account/authorized-services/) to your
GitHub account. This allows the registry to see your GitHub repositories.

To add, remove or view your linked account, log in to your Docker Hub account. Select **Settings > Linked Accounts & Services**.

Linking to Github grants Docker Hub access to all of your repositories. Follow the
onscreen instructions to authorize and link your GitHub account to Docker Hub. 
Once it is linked, you can choose a source repository from which to create the Automatic Build.

You can review and revoke Docker Hub's access by visiting the
[GitHub User's Applications settings](https://github.com/settings/applications).

> **Note**: If you delete the connection to the GitHub account that is used for one of your
> automated build repositories, previously built images are still available.
> If you re-link to that GitHub account later, the automated build can be started
> using the **Start Build** button on the Hub, or if the webhook on the GitHub repository
> still exists, it is triggered by any subsequent commits.

## Changing the GitHub user link

If you want to remove, or change the level of linking between your GitHub
account and Docker Hub, you need to make the change in two places.

First, remove the **Linked Account** from your Docker Hub **Settings**. Then go to
your GitHub account's Personal settings, and in the **Applications** section,
***Revoke access***.

You can now re-link your account at any time.

## GitHub organizations

GitHub organizations and private repositories forked from organizations are
made available for autobuilds using the "Docker Hub Registry" application, which
needs to be added to the organization - and then applied to all users.

To verify or request access, go to your GitHub **Settings** page. Select the
**Applications** section from the left side bar, then click the **View** button for
"Docker Hub Registry".

![Check User access to GitHub](images/gh-check-user-org-dh-app-access.png)

The organization's administrators may need to go to the Organization's **Third
party access** screen in **Settings** to grant or deny access to Docker Hub
Registry application. This change applies to all organization members.

![Check Docker Hub application access to Organization](images/gh-check-admin-org-dh-app-access.png)

More detailed access controls to specific users and GitHub repositories can be
managed using the GitHub **People and Teams** interfaces.

## Creating an Automated Build

You can [create an Automated Build](
https://hub.docker.com/add/automated-build/github/) from any of your
public or private GitHub repositories that have a `Dockerfile`.

Once you've selected the source repository, you can then configure:

- The Hub user/org namespace the repository is built to - either your Docker ID name, or the name of any Hub organizations your account is in
- The Docker repository name the image is built to
- The description of the repository
- If the visibility of the Docker repository is "Public" or "Private",
  you can change the accessibility options after the repository has been created.
  If you add a Private repository to a Hub user namespace, then you can only add other users
  as collaborators, and those users can view and pull all images in that
  repository. To configure more granular access permissions, such as using teams of
  users or allow different users access to different image tags, then you need
  to add the Private repository to a Hub organization for which your user has Administrator
  privileges.
- Enable or disable rebuilding the Docker image when a commit is pushed to the
  GitHub repository

You can also select one or more:
- The git branch/tag
- A repository sub-directory to use as the context
- The Docker image tag name

You can modify the description for the repository by clicking the "Description" section
of the repository view.
The "Full Description" is overwritten by the **README.md** file when the
next build is triggered.

## GitHub private submodules

If your GitHub repository contains links to private submodules, your build fails.

Normally, Docker Hub sets up a deploy key in your GitHub repository.
Unfortunately, GitHub only allows a repository deploy key to access a single
repository.

To work around this, you can create a dedicated user account in GitHub and
attach the automated build's deploy key that account. This dedicated build
account can be limited to read-only access to just the repositories required to
build.

<table class="table table-bordered">
  <thead>
    <tr>
      <th>Step</th>
      <th>Screenshot</th>
      <th>Description</th>
    </tr>
  </thead>
  <tbody>
    <tr>
      <td>1.</td>
      <td><img src="/docker-hub/images/gh_org_members.png"></td>
      <td>First, create the new account in GitHub. It should be given read-only
      access to the main repository and all submodules that are needed.</td>
    </tr>
    <tr>
      <td>2.</td>
      <td><img src="/docker-hub/images/gh_team_members.png"></td>
      <td>This can be accomplished by adding the account to a read-only team in
      the organization(s) where the main GitHub repository and all submodule
      repositories are kept.</td>
    </tr>
    <tr>
      <td>3.</td>
      <td><img src="/docker-hub/images/gh_repo_deploy_key.png"></td>
      <td>Next, remove the deploy key from the main GitHub repository. This can be done in the GitHub repository's "Deploy keys" Settings section.</td>
    </tr>
    <tr>
      <td>4.</td>
      <td><img src="/docker-hub/images/deploy_key.png"></td>
      <td>Your automated build's deploy key is in the "Build Details" menu
      under "Deploy keys".</td>
    </tr>
    <tr>
      <td>5.</td>
      <td><img src="/docker-hub/images/gh_add_ssh_user_key.png"></td>
      <td>In your dedicated GitHub User account, add the deploy key from your
      Docker Hub Automated Build.</td>
    </tr>
  </tbody>
</table>

## GitHub webhook

A GitHub webhook allows GitHub to notify Docker Hub when something has
been committed to a given Git repository.

When you create an Automated Build, a webhook should get automatically added to your GitHub
repository.

To add, confirm, or modify the webhook, log in to GitHub, then navigate to
the repository. Within the repository,  select **Settings > Webhooks**.
You must have admin privileges on the repository to view or modify
this setting. Click **Add webhook**, and use the following settings:


| Field | Value |
| ------|------ |
| Payload URL | https://registry.hub.docker.com/hooks/github |
| Content type | application/json |
| Which events would you like to trigger this webhook? | Just the push event |
| Active | checked |

The image below shows the **Webhooks/Add webhook** form with the above settings reflected:

![github-webhook-add](images/github-webhook-add.png)

If configured correctly, you'll see this in the **Webhooks** view
![github-webhook](images/github-webhook.png)


