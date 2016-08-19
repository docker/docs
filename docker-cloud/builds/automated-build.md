---
aliases:
- /docker-cloud/feature-reference/automated-build/
description: Automated builds
keywords:
- automated, build, images
menu:
  main:
    parent: builds
    weight: -50
title: Automated builds
---

# Automated builds

> **Note**: Docker Cloud's Build functionality is in BETA.

Docker Cloud can automatically build images from source code in an external
repository and automatically push the built image to your Docker
repositories.

When you set up automated builds (also called autobuilds), you create a list of
branches and tags of the images that you want to build. When you push code to a
source code branch (for example in Github) for one of those listed image tags,
the push triggers a new build. The build artifact is then pushed back to the
image repository in Docker Cloud or in an external registry.

If you have automated tests configured, these run after building but before
pushing to the registry, which you can use to create a continuous integration
workflow. Automated tests do not push images to the registry on their own.
[Learn more about automated image testing here.](automated-testing.md)

You can also push pre-built images to these repositories, even if you have
configured automatic builds.

#### Autobuild for Teams

When you create an automated build repository in your own account namespace, you can start, cancel, and retry builds, and edit and delete your own repositories.

These same actions are also available for team repositories from Docker Hub if
you are a member of the Organization's "Owners" team. If you are a member of a
team with `write` permissions you can start, cancel and retry builds in your
team's repositories, but you cannot edit the team repository settings or delete
the team repositories. If your user account has `read` permission, or if you're
a member of a team with `read` permission, you can view the build configuration
including any testing settings.

## Configure automated build settings

You can configure your repositories in Docker Cloud so that they automatically
build an image each time you push new code to your source provider. If you have
[automated tests](automated-testing.md) configured, the new image is only pushed
when the tests succeed.

Before you set up automated builds you need to [link to your source code provider](link-source.md), and [have a repository](repos.md) to build.

1. From the **Repositories** section, click into a repository to view its details.

2. Click the **Builds** tab.

3. The first time you configure automated builds for a repository, you'll see
buttons that allow you to link to an external source code repository. Select the
repository service where the image's source is stored.

    (If you haven't yet linked a source provider, follow the instructions
    [here](link-source.md) to link your account.)

    If you are editing an existing the build settings for an existing automated
    build, click **Configure automated builds**.

4. If necessary, select the **source repository** to build the repository from.

5. Select the **source repository** to build the Docker images from.

    You might need to specify an organization or user from the source code
    provider to find the code repository you want to build.

6. Choose where to run your build processes.

    You can either run the process on your own infrastructure and optionally
    [set up specific nodes to build on](#set-up-builder-nodes), or use the
    hosted build service offered on Docker Cloud's infrastructure. If you use
    Docker's infrastructure, select a builder size to run the build process on.
    This hosted build service is free while it is in Beta.

    ![](images/edit-repository-builds.png)

7. Optionally, enable [autotests](automated-testing.md#enable-automated-tests-on-a-repository).

8. In the **Tag mappings** section, enter one or more tags to build.

    For each tag:

    * Select the **Source type** to build: either a **tag** or a
    **branch**. This tells the build system what to look for in the source code
    repository.

    * Enter the name of the **Source** branch or tag you want to build.

        You can enter a name, or use a regex to match which source branch or tag
        names to build. To learn more, see
        [regexes](#regexes-and-automated-builds).

    * Specify the **Dockerfile location** as a path relative to the root of the source code repository. (If the Dockerfile is at the repository root, leave this path set to `/`.)

    * Enter the tag to apply to Docker images built from this source.
        If you configured a regex to select the source, you can reference the
        capture groups and use its result as part of the tag. To learn more, see
        [regexes](#regexes-and-automated-builds).

9. For each branch or tag, enable or disable the **Autobuild** toggle.

    Only branches or tags with autobuild enabled are built, tested,
    *and* pushed. Branches with autobuild disabled will be built for testing purposes if enabled, but not pushed.

10. Click **Save** to save the settings, or click **Save and build** to save and
run an initial test.

    A webhook is automatically added to your source code repository to notify
    Docker Cloud on every push. Only pushes to branches that are listed as the
    source for one or more tags will trigger a build.

## Check your active builds

1. To view active builds, go to the repository view and click **Timeline**.

    The Timeline displays the pending, in progress, successful and failed builds
    for any tag of the repository.

2. Click to expand a timeline entry to check the build logs.

You can click the **Cancel** button for pending builds and builds in progress.
If a build fails, the cancel button is replaced by a **Retry** button.

![](images/cancel-build.png)

## Disable an automated build

Automated builds are enabled per branch or tag, and can be disabled and
re-enabled easily. You might do this when you want to only build manually for
awhile, for example when you are doing major refactoring in your code. Disabling
autobuilds does not disable [autotests](automated-testing.md).

To disable an automated build:

1. From the **Repositories** page, click into a repository, and click the **Builds** tab.

2. Click **Configure automated builds** to edit the repository's build settings.

3. In the **Tag mappings** section, locate the branch or tag you no longer want
to automatically build.

4. Click the **autobuild** toggle next to the branch configuration line.

    The toggle turns gray when disabled.

5. Click **Save** to save your changes.

## Regexes and automated builds

You can specify a regular expression (regex) so that only matching branches or
tags are built. You can also use the results of the regex to create the Docker
tag that is applied to the built image.

You can use the variable `{sourceref}` to use the branch or tag name that
matched the regex. (The variable includes the whole source name, not just the
portion that matched the regex.) You can also use up to nine regular expression
capture groups (expressions enclosed in parentheses) to select a source to
build, and reference these in the Docker Tag field using `{/1}` through `{/9}`.

**Regex example: build from version number branch and tag with version number**

You might want to automatically build any branches that end with a number
formatted like a version number, and tag their resulting Docker images using a
name that incorporates that branch name.

To do this, specify a `branch` build with the regex `/[0-9.]+$/` in the
**Source** field, and use the formula `version-{sourceref}` in the **Docker
tag** field.

<!-- Not a priority
#### Regex example: build from version number branch and tag with version number

You could also use capture groups to build and label images that come from various sources. For example, you might have

`/(alice|bob)-v([0-9.]+)/` -->

## What's Next?

### Customize your build process

Additional advanced options are available for customizing your automated builds,
including utility environment variables, hooks, and build phase overrides. To
learn more see [Advanced options for Autobuild and Autotest](advanced.md).

### Set up builder nodes

If you are building on your own infrastructure, you can run the build process on
specific nodes by adding the `builder` label to them. If no builder nodes are
specified, the build containers are deployed using an "emptiest node" strategy.

You can also limit the number of concurrent builds (including `autotest` builds)
on a specific node by using a `builder=n` tag, where the `n` is the number of
builds to allow. For example a node tagged with `builder=5` only allows up to
five concurrent builds or autotest-builds at the same time.

### Autoredeploy services on successful build

You can configure your services to automatically redeploy once the build
succeeds. [Learn more about autoredeploy](../apps/auto-redeploy.md)

### Add automated tests

To test your code before the image is pushed, you can use
Docker Cloud's [Autotest](automated-testing.md) feature which
integrates seamlessly with autobuild and autoredeploy.

> **Note**: While the Autotest feature builds an image for testing purposes, it
does not push the resulting image to Docker Cloud or the external registry.
