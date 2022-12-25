---
description: Image analysis and reporting of packages and vulnerabilities
keywords: analysis, vulnerabilities, Hub, scanning
title: Image analysis
redirect_from:
  - /docker-hub/vulnerability-scanning/
---

{% include upgrade-cta.html
  body="Docker Hub image analysis is available for users subscribed to Docker Pro, Team, or a Business subscription. Upgrade now to get automatic vulnerability reports for your images."
  header-text="This feature requires a paid Docker subscription"
  target-url="https://www.docker.com/pricing?utm_source=docker&utm_medium=webreferral&utm_campaign=docs_driven_upgrade_scan"
%}

When image analysis is active for a repository, pushing a new tag to that
repository triggers an automatic content analysis. The analysis identifies
software artifacts included in your image, and detects whether any of the
artifacts are associated with a known vulnerability.

This allows developers and teams to review the security state of the container
images and take actions to fix issues identified during the analysis, resulting
in more secure deployments. The analysis result includes the source of the
vulnerability, such as OS packages and libraries, version in which it was
introduced, and a recommended fixed version (if available) to remediate the
vulnerabilities discovered.

## Analyze images on Docker Hub

Image analysis allows repository owners and administrators of a Docker Pro,
Team, or a Business tier to activate and deactivate image analysis. When
analysis is active for a specific repository, anyone with push access can
trigger an analysis by pushing an image to Docker Hub.

Repository owners in a Docker Pro subscription, and team members in a Team, or a
Business subscription, can view the detailed analysis result.

> **Note**
>
> Docker Hub currently supports analyzing images using the `linux/x86_64` (or
> `linux/amd64`) platform, and are less than 10 GB in size.

### Activate image analysis

Repository owners and administrators can activate image analysis on a
repository.

To activate image analysis:

1. Log into your Docker Hub account.
2. Click **Repositories** from the main menu and select a repository from the
   list.
3. By default, image analysis is inactive for all repositories. You can activate
   analysis for each repository individually. Go to the **Settings** tab and
   select **Activate image analysis**.

### Start image analysis

To trigger an image analysis, push the image to a Docker Hub repository where
analysis is active:

1. Ensure you have installed Docker locally.
2. Sign in to your Docker ID, either using the `docker login` command or the
   **Sign in** button in Docker Desktop.
3. Tag the image that you’d like to analyze. For example, to tag a Redis image,
   run:

   ```console
   $ docker tag redis <org>/<imagename>:latest
   ```

4. Push the image to Docker Hub to trigger analysis of the image:

   ```console
   $ docker push <org>/<imagename>:latest
   ```

### View the vulnerability report

To view the vulnerability report on Docker Hub:

1. Go to Docker Hub and open the repository page. A vulnerability summary is
   available in the **Tags** section of the page.

   It may take a few minutes for the vulnerability report to appear in your
   repository. If your vulnerability summary doesn't display yet, wait a moment
   and then refresh the page.

   <!-- TODO: add screenshot -->

2. Click on the tag in the table. This opens the details page for the tag.

3. Select the **Vulnerabilities** tab on the right side of the page.

   This tab displays a deep-dive view of your vulnerability exposure. If your
   image contains any packages known to be affected by vulnerabilities, those
   packages show up here.

   <!-- TODO: add screenshot -->

Expanding any of the packages in the list shows you more information about the
vulnerabilities that affect a given package. Clicking on the hyperlink of a
vulnerability opens it in the Docker image vulnerability database, which
provides even more information on the vulnerability and what images it affects.

## Deactivate image analysis

Repository owners and administrators can deactivate image analysis for a
repository. To deactivate analysis:

1. Go to Docker Hub and sign in.
2. Select **Repositories** from the main menu and select a repository from the
   list.
3. Go to the **Settings** tab and select **Deactivate image analysis**.

## Fixing vulnerabilities

Once a list of vulnerabilities have been identified, there are a couple of
actions you can take to remediate the vulnerabilities. For example, you can:

1. Specify an updated base image in the Dockerfile, check your application-level
   dependencies, rebuild the Docker image and then push the new image to Docker
   Hub.
2. Rebuild the Docker image, run an update command on the OS packages, and push
   a newer version of image to Docker Hub.
3. Edit the Dockerfile to manually remove or update specific libraries that
   contain vulnerabilities, rebuild the image, and push the new image to Docker
   Hub

After you have followed the steps suggested above, browse the new vulnerability
report to view the updated analysis result.

## Feedback

Thank you for trying out the image analysis feature. Give feedback or report any
bugs you may find through the issues tracker on the
[hub-feedback](https://github.com/docker/hub-feedback/issues){: target="_blank"
rel="noopener" class="_"} GitHub repository.
