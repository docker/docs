---
description: Set up automated builds
keywords: automated, build, images, Docker Hub
redirect_from:
- /docker-hub/builds/automated-build/
- /docker-cloud/feature-reference/automated-build/
- /docker-cloud/builds/automated-build/
- /docker-cloud/builds/
- /docker-hub/builds/classic/
title: Set up Automated Builds
---

{% include upgrade-cta.html
  header-text="This feature requires a Docker subscription"
  target-url="https://www.docker.com/pricing?utm_source=docker&utm_medium=webreferral&utm_campaign=docs_driven_upgrade_auto_builds"
%}

This page contains information on:
-[Configuring automate builds](#configure-automated-build-settings)
- 

## Configure automated build settings

You can configure repositories in Docker Hub so that they automatically
build an image each time you push new code to your source provider. If you have
[automated tests](automated-testing.md) configured, the new image is only pushed
when the tests succeed.

You can add builds to existing repositories, or add them when you create a repository.

1. From the **Repositories** section, select a repository to view its details.

2. Select the **Builds** tab.

3. If you are setting up automated builds for the first time, select the code
   repository service (GitHub or Bitbucket) where the image's source code is stored.

   > Note
   >
   > You may be redirected to the settings page to [link](link-source.md) the
   > code repository service. Otherwise, if you are editing the build settings
   > for an existing automated build, click **Configure automated builds**.

4. Select the **source repository** to build the Docker images from.

   > Note
   > You might need to specify an organization or user (the _namespace_) from
   > the source code provider. Once you select a namespace, its source code
   > repositories appear in the **Select repository** dropdown list.

5. Optionally, enable [autotests](automated-testing.md#enable-automated-tests-on-a-repository).

6. Review the default **Build Rules**, and optionally select the
   **plus sign** to add and configure more build rules.

    _Build rules_ control what Docker Hub builds into images from the contents
    of the source code repository, and how the resulting images are tagged
    within the Docker repository.

    A default build rule is set up for you, which you can edit or delete. This
    default set builds from the `Branch` in your source code repository called
    `master`, and creates a Docker image tagged with `latest`.

7. For each branch or tag, enable or disable the **Autobuild** toggle.

    Only branches or tags with autobuild enabled are built, tested, and have
    the resulting image pushed to the repository. Branches with autobuild
    disabled are built for test purposes (if enabled at the repository
    level), but the built Docker image isn't pushed to the repository.

8. For each branch or tag, enable or disable the **Build Caching** toggle.

    [Build caching](../../develop/develop-images/dockerfile_best-practices.md#leverage-build-cache)
    can save time if you are building a large image frequently or have
    many dependencies. Leave the build caching disabled to
    make sure all of your dependencies are resolved at build time, or if
    you have a large layer that's quicker to build locally.

9. Click **Save** to save the settings, or click **Save and build** to save and
   run an initial test.

    > Note
    >
    > A webhook is automatically added to your source code repository to notify
    > Docker Hub on every push. Only pushes to branches that's listed as the
    > source for one or more tags trigger a build.

### Set up build rules

By default when you set up automated builds, a basic build rule is created for you.
This default rule watches for changes to the `master` branch in your source code
repository, and builds the `master` branch into a Docker image tagged with
`latest`.

In the **Build Rules** section, enter one or more sources to build.

For each source:

* Select the **Source type** to build either a tag or a branch. This
  tells the build system what to look for in the source code repository.

* Enter the name of the **Source** branch or tag you want to build.

  The first time you configure automated builds, a default build rule is set up
  for you. This default set builds from the `Branch` in your source code called
  `master`, and creates a Docker image tagged with `latest`.

  You can also use a regex to select which source branches or tags to build.
  To learn more, see
  [regexes](index.md#regexes-and-automated-builds).

* Enter the tag to apply to Docker images built from this source.

  If you configured a regex to select the source, you can reference the
  capture groups and use its result as part of the tag. To learn more, see
  [regexes](index.md#regexes-and-automated-builds).

* Specify the **Dockerfile location** as a path relative to the root of the source code repository. If the Dockerfile is at the repository root, leave this path set to `/`.

> **Note**
>
> When Docker Hub pulls a branch from a source code repository, it performs a
> shallow clone (only the tip of the specified branch). Refer to
> [Advanced options for Autobuild and Autotest](advanced.md#source-repository--branch-clones)
> for more information.

### Environment variables for builds

You can set the values for environment variables used in your build processes
when you configure an automated build. Add your build environment variables by
clicking the plus sign next to the **Build environment variables** section, and
then entering a variable name and the value.

When you set variable values from the Docker Hub UI, you can use them by the
commands you set in `hooks` files. However, they're stored so that only users who have `admin` access to the Docker Hub repository can see their values. This
means you can use them to store access tokens or other information that
should remain secret.

> **Note**
>
> The variables set on the build configuration screen are used during
> the build processes only and shouldn't get confused with the environment
> values used by your service (for example to create service links).



## Cancel or retry a build

While a build is in queue or running, a **Cancel** icon appears next to its build
report link on the General tab and on the Builds tab. You can also click the
**Cancel** on the build report page, or from the Timeline tab's logs
display for the build.

![List of builds showing the cancel icon](images/build-cancelicon.png)

## Advanced automated build options

At the minimum you need a build rule composed of a source branch (or tag) and
destination Docker tag to set up an automated build. You can also change where
the build looks for the Dockerfile, set a path to the files the build use
(the build context), set up multiple static tags or branches to build from, and
use regular expressions (regexes) to dynamically select source code to build and
create dynamic tags.

All of these options are available from the **Build configuration** screen for
each repository. Select **Repositories** from the left navigation, and select the name of the repository you want to edit. Select the **Builds** tab, and click **Configure Automated builds**.

### Tag and branch builds

You can configure your automated builds so that pushes to specific branches or tags triggers a build.

1. In the **Build Rules** section, click the plus sign to add more sources to build.

2.  Select the **Source type** to build: either a tag or a branch.

    > Note
    >
    > This tells the build system what type of source to look for in the code
    > repository.

3. Enter the name of the **Source** branch or tag you want to build.

    > Note
    >
    > You can enter a name, or use a regex to match which source branch or tag
    > names to build. To learn more, see [regexes](index.md#regexes-and-automated-builds).

4. Enter the tag to apply to Docker images built from this source.

   > Note
   >
   > If you configured a regex to select the source, you can reference the
   > capture groups and use its result as part of the tag. To learn more, see
   > [regexes](index.md#regexes-and-automated-builds).

5. Repeat steps 2 through 4 for each new build rule you set up.

### Set the build context and Dockerfile location

Depending on how you arrange the files in your source code repository, the
files required to build your images may not be at the repository root. If that's
the case, you can specify a path where the build looks for the files.

The _build context_ is the path to the files needed for the build, relative to
the root of the repository. Enter the path to these files in the **Build context** field. Enter `/` to set the build context as the root of the source code repository.

> **Note**
>
> If you delete the default path `/` from the **Build context** field and leave
> it blank, the build system uses the path to the Dockerfile as the build
> context. However, to avoid confusion it's recommended that you specify the
> complete path.

You can specify the **Dockerfile location** as a path relative to the build
context. If the Dockerfile is at the root of the build context path, leave the
Dockerfile path set to `/`. (If the build context field is blank, set the path
to the Dockerfile from the root of the source repository.)

### Regexes and automated builds

You can specify a regular expression (regex) so that only matching branches or
tags are built. You can also use the results of the regex to create the Docker
tag that's applied to the built image.

You can use up to nine regular expression capture groups
(expressions enclosed in parentheses) to select a source to build, and reference
these in the **Docker Tag** field using `{\1}` through `{\9}`.

<!-- Capture groups Not a priority
#### Regex example: build from version number branch and tag with version number

You could also use capture groups to build and label images that come from various
sources. For example, you might have

`/(alice|bob)-v([0-9.]+)/` -->

### Build images with BuildKit

Autobuilds use the BuildKit build system by default. If you want to use the legacy
Docker build system, add the [environment variable](index.md#environment-variables-for-builds){: target="_blank" rel="noopener" class="_"}
`DOCKER_BUILDKIT=0`. Refer to the [BuildKit](../../build/buildkit/index.md)
page for more information on BuildKit.

## Autobuild for Teams

When you create an automated build repository in your own account namespace, you
can start, cancel, and retry builds, and edit and delete your own repositories.

These same actions are also available for team repositories from Docker Hub if
you are a member of the Organization's `Owners` team. If you are a member of a
team with `write` permissions you can start, cancel, and retry builds in your
team's repositories, but you cannot edit the team repository settings or delete
the team repositories. If your user account has `read` permission, or if you're
a member of a team with `read` permission, you can view the build configuration
including any testing settings.

| Action/Permission     | read | write | admin | owner |
| --------------------- | ---- | ----- | ----- | ----- |
| view build details    |  x   |   x   |   x   |   x   |
| start, cancel, retry  |      |   x   |   x   |   x   |
| edit build settings   |      |       |   x   |   x   |
| delete build          |      |       |       |   x   |

### Service users for team autobuilds

> **Note**: Only members of the `Owners` team can set up automated builds for teams.

When you set up automated builds for teams, you grant Docker Hub access to
your source code repositories using OAuth tied to a specific user account. This
means that Docker Hub has access to everything that the linked source provider
account can access.

For organizations and teams, it's recommended you create a dedicated service account (or "machine user") to grant access to the source provider. This ensures that no
builds break as individual users' access permissions change, and that an
individual user's personal projects aren't exposed to an entire organization.

This service account should have access to any repositories to be built,
and must have administrative access to the source code repositories so it can
manage deploy keys. If needed, you can limit this account to only a specific
set of repositories required for a specific build.

If you are building repositories with linked private submodules (private
dependencies), you also need to add an override `SSH_PRIVATE` environment
variable to automated builds associated with the account.

1. Create a service user account on your source provider, and generate SSH keys for it.
2. Create a "build" team in your organization.
3. Ensure that the new "build" team has access to each repository and submodule you need to build.

    Go to the repository's **Settings** page. On GitHub, add the new "build" team
    to the list of **Collaborators and Teams**. On Bitbucket, add the "build" team
    to the list of approved users on the **Access management** screen.

4. Add the service user to the "build" team on the source provider.

5. Log in to Docker Hub as a member of the `Owners` team, switch to the organization, and follow the instructions to [link to source code repository](link-source.md) using the service account.

    > **Note**: You may need to log out of your individual account on the source code provider to create the link to the service account.

6. Optionally, use the SSH keys you generated to set up any builds with private submodules, using the service account and [the instructions above](index.md#build-repositories-with-linked-private-submodules).

## What's Next?

- [Customize your build process](advanced.md) with environment variables, hooks, and more
- [Add automated tests](automated-testing.md)
- [Manage your builds](manage-builds.md)
- [Troubleshoot](troubleshoot.md)