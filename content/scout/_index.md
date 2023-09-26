---
title: Docker Scout
keywords: scout, supply chain, vulnerabilities, packages, cves, scan, analysis, analyze
description:
  Docker Scout analyzes your images to help you understand their dependencies
  and potential vulnerabilities
aliases:
  - /atomist/
  - /atomist/try-atomist/
  - /atomist/configure/settings/
  - /atomist/configure/advisories/
  - /atomist/integrate/github/
  - /atomist/integrate/deploys/
  - /engine/scan/
grid:
  - title: Quickstart
    link: /scout/quickstart/
    description: Learn what Docker Scout can do, and how to get started.
    icon: explore
  - title: Image analysis
    link: /scout/image-analysis/
    description: Reveal and dig into the composition of your images.
    icon: radar
  - title: Advisory database
    link: /scout/advisory-db-sources/
    description: Learn about the information sources that Docker Scout uses.
    icon: database
  - title: Integrations
    description: |
      Connect Docker Scout with your CI, registries, and other third-party services.
    link: /scout/integrations/
    icon: multiple_stop
  - title: Dashboard
    link: /scout/dashboard/
    description: |
      The web interface for Docker Scout.
    icon: dashboard
  - title: Policy {{< badge color=violet text="Early Access" >}}
    link: /scout/policy/
    description: |
      Ensure that your artifacts align with supply chain best practices.
    icon: policy
  - title: Upgrade
    link: /billing/scout-billing/
    description: |
      The free plan includes up to 3 repositories. Upgrade for more.
    icon: upgrade
---

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

{{< grid >}}
