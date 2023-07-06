---
title: Docker Scout
keywords: scout, supply chain, vulnerabilities, packages, cves, scan, analysis, analyze
description: >
  Docker Scout analyzes your images to help you understand their dependencies and potential vulnerabilities
redirect_from:
  - /atomist/
  - /atomist/try-atomist/
  - /atomist/get-started/
  - /atomist/configure/settings/
  - /atomist/configure/advisories/
  - /atomist/integrate/github/
  - /atomist/integrate/deploys/
  - /engine/scan/
---

{% include scout-early-access.md %}

Container images are often built from layers of other container images and
software packages. These layers and packages can contain vulnerabilities that
make your containers and the applications they run vulnerable to attack.

Docker Scout can proactively help you find and fix these vulnerabilities,
helping you create a more secure software supply chain. It does this by analyzing your images and creating a full inventory of the
packages and layers called a [Software bill of materials (SBOM)](https://ntia.gov/sites/default/files/publications/sbom_at_a_glance_apr2021_0.pdf).
It then correlates this inventory with a continuously updated vulnerability
database to identify vulnerabilities in your images.

You can use Docker Scout in Docker Desktop, Docker Hub, the Docker CLI, and in
the Docker Scout Dashboard. If you host your images in JFrog Artifactory, you
can also use Docker Scout to analyze your images there.

_The following video shows an end-to-end workflow of using Docker Scout to remediate a reported vulnerability_.

<div style="position: relative; padding-bottom: 64.86486486486486%; height: 0;"><iframe src="https://www.loom.com/embed/e066986569924555a2546139f5f61349?sid=6e29be62-78ba-4aa7-a1f6-15f96c37d916" frameborder="0" webkitallowfullscreen mozallowfullscreen allowfullscreen style="position: absolute; top: 0; left: 0; width: 100%; height: 100%;"></iframe></div>


{% include scout-plans.md %}

## Docker Desktop

Docker Scout analyzes all images stored locally in Docker Desktop, providing you
with up-to-date vulnerability information as you build your images.

For more information, read the [Advanced image analysis guide](./advanced-image-analysis.md).

## Docker Hub

If you enable [Advanced image analysis](./advanced-image-analysis.md) for a
repository in Docker Hub, Docker Scout analyzes your images every time you push
them to Docker Hub. Docker Scout shows analysis results on every tag view for
that repository.

The analysis updates continuously, meaning that the vulnerability report for an
image is always up to date as Docker Scout becomes aware of new CVEs. No need to
re-scan an image.

For more information, read the [Advanced image analysis guide](./advanced-image-analysis.md).

## Docker Scout CLI plugin {#docker-scout-cli}

The `docker scout` CLI plugin provides a terminal interface for using Docker
Scout with local and remote images.

Using the CLI, you can analyze images and view the analysis report in text
format. You can print the results directly to stdout, or export them to a file
using a structured format, such as Static Analysis Results Interchange Format
(SARIF).

For more information about how to use the `docker scout` CLI, see the
[reference documentation](../engine/reference/commandline/scout_cves.md).

The plugin is available in Docker Desktop starting with version 4.17 and
available as a standalone binary.

To install the plugin, run the following command:

```console
$ curl -fsSL https://raw.githubusercontent.com/docker/scout-cli/main/install.sh -o install-scout.sh
$ sh install-scout.sh
```

> **Note**
>
> Always examine scripts downloaded from the internet before running them locally.
> Before installing, make yourself familiar with potential risks and limitations
> of the convenience script.

If you want to install the plugin manually, you can find full instructions in
the [plugin's repository](https://github.com/docker/scout-cli).

The plugin is also available as [a container image](https://hub.docker.com/r/docker/scout-cli)
and as [a GitHub action](https://github.com/docker/scout-action).

## Docker Scout Dashboard

The [Docker Scout Dashboard](https://scout.docker.com){: target="\_blank" rel="noopener" }
helps you share the analysis and security status of images in
an organization with your team. You can also [use the dashboard to enable analysis of multiple repositories at once](./dashboard.md#repository-settings).

For more information, read the [Docker Scout Dashboard guide](./dashboard.md).

## JFrog Artifactory integration

Users of JFrog Artifactory, or JFrog Container Registry, can integrate Docker
Scout to enable automatic analysis of images locally and remotely. For more
information, see [Artifactory integration](./artifactory.md).

_The following video shows how to enable Docker Scout on your repositories_.

<div style="position: relative; padding-bottom: 64.86486486486486%; height: 0;"><iframe src="https://www.loom.com/embed/a6fb14ede0a94d0d984edf6cf16604e0?sid=ba34f694-32a6-4b74-b3f8-9cc6b80ef66f" frameborder="0" webkitallowfullscreen mozallowfullscreen allowfullscreen style="position: absolute; top: 0; left: 0; width: 100%; height: 100%;"></iframe></div>