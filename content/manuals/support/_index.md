---
title: Get support for Docker products
linkTitle: Support
description: >-
  Find Docker support options, including paid subscription response times,
  the Support portal, Docker Desktop support scope, and community resources.
keywords: docker support, support portal, subscription support, premium
  support, support response times, severity levels, community forums,
  docker desktop support
weight: 50
params:
  sidebar:
    group: Accounts and admin
---

Docker offers multiple support channels depending on your subscription
level and needs.

## Community support

All Docker users can seek support through community resources, where
Docker or the community responds on a best-effort basis:

- [Docker Community Forums](https://forums.docker.com/)
- [Docker Community Slack](http://dockr.ly/comm-slack)

## Paid subscription support

All Docker Pro, Team, and Business subscribers receive email support for
Docker products.

> [!TIP]
>
> Before reaching out for support, review the troubleshooting documentation
> for your product.

Paid subscribers create and track cases in the
[Docker Support portal](/manuals/support/portal.md).
[Who gets access](/manuals/support/portal.md#who-gets-access) covers
enrollment, including automatic access for Premium Support. See
[Severity levels](/manuals/support/portal.md#severity-levels) and
[Case statuses](/manuals/support/portal.md#case-statuses) on that page.

### Support response times

- Docker Pro: within 5 business days
- Docker Team: within 2 business days, 24×5 availability
- Docker Business: within 1 business day, 24×5 availability

Docker Verified Publisher (DVP) organizations follow the response times
for their Docker Pro, Team, or Business subscription. Without one of those
subscriptions, the following response times apply:

- DVP Starter: within 3 business days
- DVP Growth: within 2 business days, 24×5 availability

> [!TIP]
>
> Premium Support with faster response times and 24×7 availability is
> available as an add-on for
> [Docker Business and Docker Hardened Images subscribers](https://www.docker.com/pricing?ref=Docs&refAction=DocsSupport).

## Docker Desktop support

Docker Desktop support is available with a paid subscription.

### Scope of support

{{< tabs >}}
{{< tab name="Covered" >}}

Docker Desktop support includes:

- Configuration and installation issues
- Desktop updates
- Account management and billing
- Sign-in issues
- Push or pull issues, including rate limiting
- Automated builds
- Basic product how-to questions
- Application crashes or unexpected behavior
- On Windows: turning on virtualization in BIOS
- On Windows: turning on Windows features
- On Windows: running inside certain
  [VM or VDI environments](/manuals/desktop/setup/vm-vdi.md)
  (Docker Business only)

{{< /tab >}}
{{< tab name="Not covered" >}}

Docker Desktop support excludes:

- Unsupported operating systems, including beta and preview versions
- Versions older than the latest release (except Docker Business)
- Training, customization, and integration
- Altered or modified Docker software
- Docker Engine, Docker CLI, or other bundled Linux components
- Third-party applications not provided by Docker
- Features labeled as experimental
- Hardware malfunction, abuse, or improper use
- System and server administration activities
- Routine product maintenance (data backup, disk space, log rotation)
- Desktop as a production runtime
- Running containers of a different architecture using emulation
- Kubernetes
- Scale deployment and multi-machine installation
- Running multiple instances on a single machine

{{< /tab >}}
{{< /tabs >}}

### Supported versions

Docker Desktop support is scoped by plan:

| Plan | Supported versions | Machines per subscription |
| :--- | :--- | :--- |
| Docker Pro | Latest only | One |
| Docker Team | Latest only | One per subscription seat |
| Docker Business | Latest plus versions up to six months older (fixes on latest only) | Unlimited |

### Diagnostic data and privacy

- Diagnostics bundles may contain personal data such as usernames and IP
  addresses.
- Access is limited to Docker, Inc. employees directly involved in
  diagnosing the issue.
- Docker, Inc. deletes uploaded bundles after 30 days by default.
- To request earlier removal, provide the diagnostics ID or your GitHub
  ID.
- Docker, Inc. uses the data only to investigate specific user issues but
  may derive high-level (non-personal) metrics.

For more information, see the
[Docker Data Processing Agreement](https://www.docker.com/legal/data-processing-agreement).

## Next steps

- Create and track a case in the
  [Docker Support portal](/manuals/support/portal.md)
- Review system requirements:
  [Mac](/manuals/desktop/setup/install/mac-install.md#system-requirements),
  [Windows](/manuals/desktop/setup/install/windows-install.md#system-requirements),
  [Linux](/manuals/desktop/setup/install/linux/_index.md#system-requirements)
- Explore the [Docker Desktop issue tracker](https://github.com/docker/desktop-feedback)
