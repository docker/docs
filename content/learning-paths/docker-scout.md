---
title: Securing your software supply chain with Docker Scout
summary: |
  Enhance container security by automating vulnerability detection and
  remediation, ensuring compliance, and protecting your development workflow.
description: |
  Learn how to use Docker Scout to enhance container security by automating
  vulnerability detection and remediation, ensuring compliance, and protecting
  your development workflow.
params:
  image: images/learning-paths/scout.png
  skill: Beginner
  time: 10 minutes
  prereq: None
---

{{< columns >}}

When container images are insecure, significant risks can arise. Around 60% of
organizations have reported experiencing at least one security breach or
vulnerability incident within a year, resulting in operational
disruption.[^CSA] These incidents often result in considerable downtime, with
44% of affected companies experiencing over an hour of downtime per event. The
financial impact is substantial, with the average data breach cost reaching
$4.45 million.[^IBM] This highlights the critical importance of maintaining
robust container security measures.

Docker Scout enhances container security by providing automated vulnerability
detection and remediation, addressing insecure container images, and ensuring
compliance with security standards.

[^CSA]: https://cloudsecurityalliance.org/blog/2023/09/21/2023-global-cloud-threat-report-cloud-attacks-are-lightning-fast
[^IBM]: https://www.ibm.com/reports/data-breach

<!-- break -->

## What you'll learn

- Define secure software supply chain (SSSC)
- Review SBOMs and how to use them
- Detect and monitor vulnerabilities

## Tools integration

Works well with Docker Desktop, GitHub Actions, Jenkins, Kubernetes, and
other CI solutions.

{{< /columns >}}

## Who’s this for?

- DevOps engineers who need to integrate automated security checks into CI/CD
  pipelines to enhance the security and efficiency of their workflows.
- Developers who want to use Docker Scout to identify and remediate
  vulnerabilities early in the development process, ensuring the production of
  secure container images.
- Security professionals who must enforce security compliance, conduct
  vulnerability assessments, and ensure the overall security of containerized
  applications.

## Modules

{{< accordion large=true title=`Why Docker Scout?` icon=`play_circle` >}}

Organizations face significant challenges from data breaches,
including financial losses, operational disruptions, and long-term damage to
brand reputation and customer trust. Docker Scout addresses critical problems
such as identifying insecure container images, preventing security breaches,
and reducing the risk of operational downtime due to vulnerabilities.

Docker Scout provides several benefits:

- Secure and trusted content
- A system of record for your Software Development Lifecycle (SDLC)
- Continuous security posture improvement

Docker Scout offers automated vulnerability detection and remediation, helping
organizations identify and fix security issues in container images early in the
development process. It also integrates with popular development tools like
Docker Desktop and GitHub Actions, providing seamless security management and
compliance checks within existing workflows.

**Duration**: 5 minutes

{{< youtube-embed "-omsQ7Uqyc4" >}}

{{< /accordion >}}

{{< accordion large=true title=`Docker Scout Demo` icon=`play_circle` >}}

Docker Scout has powerful features for enhancing containerized application
security and ensuring a robust software supply chain.

- Define vulnerability remediation
- Discuss why remediation is essential to maintain the security and integrity
  of containerized applications
- Discuss common vulnerabilities
- Implement remediation techniques: updating base images, applying patches,
  removing unnecessary packages
- Verify and validate remediation efforts using Docker Scout

**Duration**: 5 minutes

{{< youtube-embed "TkLwJ0p46W8" >}}

{{< /accordion >}}

{{< accordion large=true title=`Common challenges and questions` icon=`quiz` >}}

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

{{< /accordion >}}

{{< accordion large=true title=`Resources` icon=`link` >}}

- [Docker Scout overview](/scout/)
- [Docker Scout quickstart](/scout/quickstart/)
- [Install Docker Scout](/scout/install/)
- [Software Bill of Materials](/scout/concepts/sbom/)

<!-- vale Docker.HeadingLength = YES -->

{{< /accordion >}}
