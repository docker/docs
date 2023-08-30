---
title: Docker Scout
keywords: scout, supply chain, vulnerabilities, packages, cves, scan, analysis, analyze
description: Docker Scout analyzes your images to help you understand their dependencies
  and potential vulnerabilities
aliases:
 - /atomist/
 - /atomist/try-atomist/
 - /atomist/configure/settings/
 - /atomist/configure/advisories/
 - /atomist/integrate/github/
 - /atomist/integrate/deploys/
 - /engine/scan/
---

{{< include "scout-early-access.md" >}}

Container images are often built from layers of other container images and
software packages. These layers and packages can contain vulnerabilities that
make your containers and the applications they run vulnerable to attack.

Docker Scout can proactively help you find and fix these vulnerabilities,
helping you create a more secure software supply chain. It does this by analyzing your images and creating a full inventory of the
packages and layers called a [Software bill of materials (SBOM)](https://ntia.gov/sites/default/files/publications/sbom_at_a_glance_apr2021_0.pdf).
It then correlates this inventory with a continuously updated vulnerability
database to identify vulnerabilities in your images.

You can use Docker Scout in Docker Desktop, Docker Hub, the Docker CLI, and in
the [Docker Scout Dashboard](./dashboard.md). Docker Scout also supports
integrations with third-party systems, refer to [Integrating Docker
Scout](./integrations/index.md) for more information.

{{< include "scout-plans.md" >}}

## Quickstart

_The following video shows an end-to-end workflow of using Docker Scout to remediate a reported vulnerability_.

<div style="position: relative; padding-bottom: 64.86486486486486%; height: 0;"><iframe src="https://www.loom.com/embed/e066986569924555a2546139f5f61349?sid=6e29be62-78ba-4aa7-a1f6-15f96c37d916" frameborder="0" webkitallowfullscreen mozallowfullscreen allowfullscreen style="position: absolute; top: 0; left: 0; width: 100%; height: 100%;"></iframe></div>

> **Quickstart with Docker Scout**
>
> For a self-guided quickstart that shows you how to use Docker Scout to identify and remediate vulnerabilities in your images, read the [quickstart](./quickstart.md).
{ .tip }

## Enabling Docker Scout

_The following video shows how to enable Docker Scout on your repositories_.

<div style="position: relative; padding-bottom: 64.86486486486486%; height: 0;"><iframe src="https://www.loom.com/embed/a6fb14ede0a94d0d984edf6cf16604e0?sid=ba34f694-32a6-4b74-b3f8-9cc6b80ef66f" frameborder="0" webkitallowfullscreen mozallowfullscreen allowfullscreen style="position: absolute; top: 0; left: 0; width: 100%; height: 100%;"></iframe></div>
