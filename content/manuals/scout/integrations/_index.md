---
description: How to setup Docker Scout with other systems.
keywords: supply chain, security, integrations, registries, ci, environments
title: Integrating Docker Scout with other systems
linkTitle: Integrations
weight: 80
---

By default, Docker Scout integrates with your Docker organization and your
Docker Scout-enabled repositories on Docker Hub. You can integrate Docker Scout
with additional third-party systems to get access to even more insights,
including real-time information about you running workloads.

## Integration categories

You'll get different insights depending on where and how you choose to integrate
Docker Scout.

### Container registries

Integrating Docker Scout with third-party container
registries enables Docker Scout to run image analysis on those repositories,
so that you can get insights into the composition of those images even if they
aren't hosted on Docker Hub.

The following container registry integrations are available:

- [Amazon Elastic Container Registry](./registry/ecr.md)
- [Azure Container Registry](./registry/acr.md)
- [JFrog Artifactory](./registry/artifactory.md)

### Continuous Integration

Integrating Docker Scout with Continuous Integration (CI) systems is a great way
to get instant, automatic feedback about your security posture in your inner
loop. Analysis running in CI also gets the benefit of additional context that's
useful for getting even more insights.

The following CI integrations are available:

- [GitHub Actions](./ci/gha.md)
- [GitLab](./ci/gitlab.md)
- [Microsoft Azure DevOps Pipelines](./ci/azure.md)
- [Circle CI](./ci/circle-ci.md)
- [Jenkins](./ci/jenkins.md)

### Environment monitoring

Environment monitoring refers to integrating Docker Scout with your deployments.
This can give you information in real-time about your running container workloads.

Integrating with environments lets you compare production workloads to other
versions, in your image repositories or in your other environments.

The following environment monitoring integrations are available

- [Sysdig](./environment/sysdig.md)

For more information about environment integrations, see
[Environments](./environment/_index.md).

### Code quality

Integrating Docker Scout with code analysis tools enables quality checks
directly on source code, helping you keep track of bugs, security issues, test
coverage, and more. In addition to image analysis and environment monitoring,
code quality gates let you shift left your supply chain management with Docker
Scout.

Once you enable a code quality integration, Docker Scout includes the code
quality assessments as policy evaluation results for the repositories where
you've enabled the integration.

The following code quality integrations are available:

- [SonarQube](sonarqube.md)

### Source code management

Integrate Docker Scout with your version control system to get guided
remediation advice on how to address issues detected by Docker Scout image
analysis, directly in your repositories.

The following source code management integrations are available:

- [GitHub](source-code-management/github.md) {{< badge color=blue text=Beta >}}

### Team collaboration

Integrations in this category let you integrate Docker Scout with collaboration
platforms for broadcasting notifications about your software supply chain in
real-time to team communication platforms.

The following team collaboration integrations are available:

- [Slack](./team-collaboration/slack.md)
