---
description: How to setup Docker Scout with other systems.
keywords: supply chain, security, integrations, registries, ci, environments
title: Integrating Docker Scout with other systems
---

{% include scout-early-access.md %}

By default, Docker Scout integrates with your Docker organization and your
Docker Scout-enabled repositories on Docker Hub. You can integrate Docker Scout
with additional third-party systems to get access to even more insights,
including realtime information about you running workloads.

## Integration categories

You'll get different insights depending on where and how you choose to integrate
Docker Scout.

### Container registries

Integrating Docker Scout with third-party container
registries enables Docker Scout to run image analysis on those repositories,
so that you can get insights into the composition of those images even if they
aren't hosted on Docker Hub.

The following container registry integrations are available:

- [Artifactory](./registry/artifactory.md)
- Amazon ECR (coming soon)

### Continuous Integration

Integrating Docker Scout with Continuous Integration (CI) systems is a great
way to get instant, automatic feedback about your security posture in your inner
loop. Analysis running in CI also gets the benefit of additional context that's
useful for getting even more insights.

The following CI integrations are available:

- [GitHub Actions](./ci/gha.md)
- [GitLab](./ci/gitlab.md)
- [Microsoft Azure DevOps Pipelines](./ci/azure.md)
- [Circle CI](./ci/circle-ci.md)
- [Jenkins](./ci/jenkins.md)
