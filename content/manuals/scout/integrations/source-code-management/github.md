---
title: Integrate Docker Scout with GitHub
linkTitle: GitHub
description: Integrate Docker Scout using the GitHub app to get remediation advice directly in your repositories
keywords: scout, github, integration, image analysis, supply chain, remediation, source code
params:
  sidebar:
    badge:
      color: gray
      text: Retired
---

> [!IMPORTANT]
>
> The Docker Scout GitHub integration was retired on July 1, 2026. For base-image updates, use
> GitHub Dependabot with `package-ecosystem: "docker"`. For image-to-source
> linkage, build with `--provenance=mode=max`.

> [!NOTE]
>
> This retirement applies only to the Docker Scout GitHub App integration.
> The [`docker/scout-action`](https://github.com/docker/scout-action) GitHub
> Action for CI pipelines is not affected and continues to work.

## Migrate from the GitHub integration

The integration provided two capabilities, each with an alternative.

### Base-image digest repinning

Use GitHub Dependabot with `package-ecosystem: "docker"`. Dependabot opens
PRs to update base image tags and digests on a schedule. When you pin as
`FROM image:tag@sha256:...`, both the tag and digest are updated. The Scout
integration updated only the digest.

Minimal `.github/dependabot.yml`:

```yaml
version: 2
updates:
  - package-ecosystem: "docker"
    directory: "/"
    schedule:
      interval: "weekly"
```

See [Configuring Dependabot version updates](https://docs.github.com/en/code-security/dependabot/dependabot-version-updates/configuring-dependabot-version-updates).

### Image-to-source linkage

Build with `--provenance=mode=max`. Docker Scout reads the resulting provenance
attestation to link an image back to its source repository without requiring
the GitHub app.

```console
$ docker build --provenance=mode=max -t myimage:tag .
```

See [SLSA provenance attestations](/manuals/build/metadata/attestations/slsa-provenance.md).
