---
title: Common challenges and questions
description: Explore common challenges and questions related to Docker Scout.
---

<!-- vale Docker.HeadingLength = NO -->

### How is Docker Scout different from other security tools?

Docker Scout takes a broader approach to container security compared to
third-party security tools. Third-party security tools, if they offer
remediation guidance at all, miss the mark on their limited scope of
application security posture within the software supply chain, and often
limited guidance when it comes to suggested fixes. Such tools have either
limitations on runtime monitoring or no runtime protection at all. When they do
offer runtime monitoring, it’s limited in its adherence to key policies.
Third-party security tools offer a limited scope of policy evaluation for
Docker-specific builds. By focusing on the entire software supply chain,
providing actionable guidance, and offering comprehensive runtime protection
with strong policy enforcement, Docker Scout goes beyond just identifying
vulnerabilities in your containers. It helps you build secure applications from
the ground up.

### Can I use Docker Scout with external registries other than Docker Hub?

You can use Scout with registries other than Docker Hub. Integrating Docker Scout
with third-party container registries enables Docker Scout to run image
analysis on those repositories so that you can get insights into the
composition of those images even if they aren't hosted on Docker Hub.

The following container registry integrations are available:

- Artifactory
- Amazon Elastic Container Registry
- Azure Container Registry

Learn more about configuring Scout with your registries in [Integrating Docker Scout with third-party registries](/scout/integrations/#container-registries).

### Does Docker Scout CLI come by default with Docker Desktop?

Yes, the Docker Scout CLI plugin comes pre-installed with Docker Desktop.

### Is it possible to run `docker scout` commands on a Linux system without Docker Desktop?

If you run Docker Engine without Docker Desktop, Docker Scout doesn't come
pre-installed, but you can [install it as a standalone binary](/scout/install/).

### How is Docker Scout using an SBOM?

An SBOM, or software bill of materials, is a list of ingredients that make up
software components. [Docker Scout uses SBOMs](/scout/concepts/sbom/) to
determine the components that are used in a Docker image. When you analyze an
image, Docker Scout will either use the SBOM that is attached to the image (as
an attestation), or generate an SBOM on the fly by analyzing the contents of
the image.

The SBOM is cross-referenced with the advisory database to determine if any of
the components in the image have known vulnerabilities.

<div id="scout-lp-survey-anchor"></div>
