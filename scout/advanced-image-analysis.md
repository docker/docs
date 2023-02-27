---
description: Image analysis with Docker Scout
keywords: scanning, vulnerabilities, Hub, supply chain, security
title: Advanced image analysis
---

> **Note**
>
> Docker Scout is an [early access](../release-lifecycle.md#early-access-ea)
> product, and requires a Docker Pro, Team, or Business subscription.

Advanced image analysis is a Docker Scout feature for Docker Hub.

When you activate Advanced image analysis for a repository, new tags are
automatically analyzed when you push to that repository. Advanced image analysis
is more than point-in-time scanning - your analysis gets reevaluated
continuously, meaning you don't need to re-scan the image to see an updated
vulnerability report. Information about packages and CVEs displays on the image
tag page on Docker Hub, as well as in the image details view in Docker Desktop.

Advanced image analysis supports the following package ecosystems:

- .NET
- GitHub packages
- Go
- Java
- JavaScript
- PHP
- Python
- RPM
- Ruby
- `alpm` (Arch Linux)
- `apk` (Alpine Linux)
- `deb` (Debian Linux and derivatives)

## Activate Advanced image analysis

Advanced image analysis is an early access feature. It can be activated on a
per-repository basis for organizations with a
[Docker Pro, Team, or Business subscription](../subscription/index.md).

> **Note**
>
> Only repository owners and administrators can activate Advanced image analysis
> on a repository.

To activate Advanced image analysis:

1. Log into your Docker Hub account.
2. Click **Repositories** from the main menu and select a repository from the
   list.
3. Go to the **Settings** tab
4. Under **Image insight settings**, select **Advanced image analysis provided
   by Docker Scout**.
5. Select **Save**.

## Analyze an image

To trigger Advanced image analysis, push the image to a Docker Hub repository
where Advanced image analysis is active:

1. Sign in to your Docker ID, either using the `docker login` command or the
   **Sign in** button in Docker Desktop.
2. Tag the image to analyze. For example, to tag a Redis image, run:

   ```console
   $ docker tag redis <org>/<imagename>:latest
   ```

3. Push the image to Docker Hub to trigger analysis of the image:

   ```console
   $ docker push <org>/<imagename>:latest
   ```

## View the vulnerability report

To view the vulnerability report on Docker Hub:

1. Go to Docker Hub and open the repository page. A vulnerability summary is
   available in the **Tags** section of the page.

   It may take a few minutes for the vulnerability report to appear in your
   repository. If your vulnerability summary doesn't display yet, wait a moment
   and then refresh the page.

2. Click on the tag in the table. This opens the details page for the tag.

3. Select the **Vulnerabilities** tab on the right side of the page.

   This tab displays a deep-dive view of your vulnerability exposure. If your
   image contains packages known to be affected by vulnerabilities, those
   packages show up here.

   For more information about how to interpret the vulnerability report, see
   [Image details view](./image-details-view.md)

Expanding any of the packages in the list shows you more information about the
vulnerabilities that affect a given package. Clicking on the hyperlink of a
vulnerability opens it in Docker's image vulnerability database, which provides
even more information on the vulnerability and what images it affects.

## Deactivate Advanced image analysis

> **Note**
>
> Only repository owners and administrators can deactivate Advanced image
> analysis on a repository.

To deactivate Advanced image analysis:

1. Go to Docker Hub and sign in.
2. Select **Repositories** from the main menu and select a repository from the
   list.
3. Go to the **Settings** tab.
4. Under **Image insight settings**, and select one of the following options:

   - **None** if you want to completely turn off vulnerability detection
   - **Hub vulnerability scanning** if you want to use the basic scanning
     feature.

5. Select **Save**.

## Feedback

Thank you for trying out the Advanced image analysis feature. Give feedback or
report any bugs you may find through the issues tracker on the
[hub-feedback](https://github.com/docker/hub-feedback/issues){: target="_blank"
rel="noopener" class="_"} GitHub repository.
