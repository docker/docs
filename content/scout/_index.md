---
title: Docker Scout
keywords: scout, supply chain, vulnerabilities, packages, cves, scan, analysis, analyze
description:
  Get an overview on Docker Scout to proactively enhance your software supply chain security
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

Container images consist of layers and software packages, which are susceptible to vulnerabilities.
These vulnerabilities can compromise the security of containers and applications.

Docker Scout is a solution for proactively enhancing your software supply chain security.
By analyzing your images, Docker Scout compiles an inventory of components, also known as a Software Bill of Materials (SBOM).
The SBOM is matched against a continuously updated vulnerability database to pinpoint security weaknesses.

Docker Scout is a standalone service and platform that you can interact with
using Docker Desktop, Docker Hub, the Docker CLI, and the Docker Scout Dashboard.
Docker Scout also facilitates integrations with third-party systems, such as container registries and CI platforms.

{{< grid >}}
