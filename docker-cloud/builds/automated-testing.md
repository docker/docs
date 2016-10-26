---
aliases:
- /docker-cloud/feature-reference/automated-testing/
description: Automated tests
keywords:
- Automated, testing, repository
menu:
  main:
    parent: builds
    weight: -50
title: Automated repository tests
---

# Automated repository tests

Docker Cloud can automatically test changes pushed to your source code
repositories using containers. You can enable `Autotest` on [any Docker Cloud repository](repos.md) to run tests at each push to the source code repository,
similar to a continuous integration testing service.

Enabling `Autotest` builds an image for testing purposes, but does **not**
automatically push the built image to the Docker repository. If you want to push
built images to your Docker Cloud repository, enable [Automated Builds](automated-build.md).

## Set up automated test files

To set up your automated tests, create a `docker-compose.test.yml` file in the
root of your source code repository which defines a `sut` service that lists the
tests to be run.

For example:

```yml
sut:
  build: .
  command: run_tests.sh
```

The example above builds the repository, and runs the `run_tests.sh` file inside
a container using the built image.

You can define any number of linked services in this file. The only requirement
is that `sut` is defined. Its return code determines if tests passed or not:
tests **pass** if the `sut` service returns `0`, and **fail** otherwise.

You can define more than one `docker-compose.test.yml` file if needed. Any file
that ends in `.test.yml` is used for testing, and the tests run sequentially.

> **Note**: If you enable Automated builds, they will also run any tests defined
in the `test.yml` files.

## Enable automated tests on a repository

To enable testing on a source code repository, you must first create an
associated build-repository in Docker Cloud.  Your `Autotest` settings are
configured on the same page as [automated builds](automated-build.md), however
you do not need to enable Autobuilds to use `Autotest`. Autobuild is enabled per
branch or tag, and you do not need to enable it at all.

Only branches that are configured to use **Autobuild** will push images to the
Docker repository, regardless of the Autotest settings.

1. Log in to Docker Cloud and select **Repositories** in the left navigation.

3. Select the repository you want to enable `Autotest` on.

4. From the repository view, click the **Builds** tab.

4. Click **Configure automated builds**.

5. Configure the automated build settings as explained in [Automated Builds](automated-build.md).

    At minimum you must configure:

    * The source code repository
    * the build location
    * at least one build rule

8. Choose your **Autotest** option.

    The following options are available:

    * `Off`: no additional tests. Test commits only to branches that are using Autobuild to build and push images.
    * `Source Repository`: test commits to all branches of the source code repository, regardless of their Autobuild setting.
    * `Source Repository & External Pull Requests`: tests commits to all branches of the source code repository, including any pull requests opened against it.

    > **Note**: For security purposes, autotest on _external pull requests_ is
    disabled on public repositories. If you select this option on a public
    repository, tests will still run on pushes to the source code repository,
    but not on pull requests.

9. Click **Save** to save the settings, or click **Save and build** to save and
run an initial test.

## Check your test results

From the repository's details page, click **Timeline**.

From this tab you can see any pending, in-progress, successful, and failed
builds and test runs for the repository.

You can click any timeline entry to view the logs for each test run.
