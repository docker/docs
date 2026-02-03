---
title: Get support for Docker products
linkTitle: Support
description: Learn about support options for Docker products including paid subscriptions and community resources
keywords: support, help, docker desktop, subscriptions, community, troubleshooting
weight: 5
params:
  sidebar:
    group: Platform
---

Docker offers multiple support channels depending on your subscription level and needs.

## Paid subscription support

All Docker Pro, Team, and Business subscribers receive email support for Docker products.

### Support response times

- Docker Pro: 3 business day response
- Docker Team: 2 business day response, 24×5 availability
- Docker Business: 1 business day response, 24×5 availability

> [!NOTE]
>
> Premium Support with faster response times and 24×7 availability is available as an add-on for Docker Business subscribers.

For detailed support features and response times, see [Docker Pricing](https://www.docker.com/pricing/).

### Support severity levels

| Level    | Description                                                                                                                                                                |
| :------- | :------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| Critical | Widespread or company-wide service outage affecting many customers or all users within a single organization. Business operations are halted with no workaround available. |
| High     | Team or department-level impact preventing significant users from accessing core functionality. Severe business impact with no workaround exists.                          |
| Medium   | Individual user or small group impact causing partial loss of functionality. Business operations continue, often with workarounds available but reduced productivity.      |

### Request support

> [!TIP]
>
> Before reaching out for support, review the troubleshooting documentation for your product.

If you have a paid Docker subscription, [contact the Support team](https://hub.docker.com/support/contact/).

## Community support

All Docker users can seek support through community resources, where Docker or the community respond on a best effort basis:

- [Docker Community Forums](https://forums.docker.com/)
- [Docker Community Slack](http://dockr.ly/comm-slack)

## Docker Desktop support

Docker Desktop support is available with a paid subscription.

### Scope of support

{{< tabs >}}
{{< tab name="Covered">}}

Docker Desktop support includes:

- Account management and billing
- Configuration and installation issues
- Desktop updates
- Sign-in issues
- Push or pull issues, including rate limiting
- Application crashes or unexpected behavior
- Automated builds
- Basic product 'how to' questions

**Windows-specific:**

- Turning on virtualization in BIOS
- Turning on Windows features
- Running inside [certain VM or VDI environments](/manuals/desktop/setup/vm-vdi.md) (Docker Business only)

{{< /tab >}}
{{< tab name="Not covered">}}

Docker Desktop support excludes:

- Unsupported operating systems, including beta/preview versions
- Running containers of a different architecture using emulation
- Docker Engine, Docker CLI, or other bundled Linux components
- Kubernetes
- Features labeled as experimental
- System/Server administration activities
- Desktop as a production runtime
- Scale deployment/multi-machine installation
- Routine product maintenance (data backup, disk space, log rotation)
- Third-party applications not provided by Docker
- Altered or modified Docker software
- Hardware malfunction, abuse, or improper use
- Versions older than the latest release (except Docker Business)
- Training, customization, and integration
- Running multiple instances on a single machine

> [!NOTE]
>
> Support for [running Docker Desktop in a VM or VDI environment](/manuals/desktop/setup/vm-vdi.md) is only available to Docker Business customers.

{{< /tab >}}
{{< /tabs >}}

### Supported versions

- Docker Business: Versions up to six months older than the latest version (fixes applied to latest version only)
- Docker Pro and Team: Latest version only

### Number of machines

- Docker Pro: One machine
- Docker Team: Number of machines equal to subscription seats
- Docker Business: Unlimited machines

### Supported operating systems

- [Mac system requirements](/manuals/desktop/setup/install/mac-install.md#system-requirements)
- [Windows system requirements](/manuals/desktop/setup/install/windows-install.md#system-requirements)
- [Linux system requirements](/manuals/desktop/setup/install/linux/_index.md#system-requirements)

### Community resources

- [Docker Desktop issue tracker](https://github.com/docker/desktop-feedback)

### Diagnostic data and privacy

When uploading diagnostics, the bundle may contain personal data such as usernames and IP addresses. Diagnostics bundles are only accessible to Docker, Inc. employees directly involved in diagnosing issues.

By default, Docker, Inc. deletes uploaded diagnostics bundles after 30 days. You may request removal of a diagnostics bundle by specifying the diagnostics ID or your GitHub ID. Docker, Inc. only uses the data to investigate specific user issues but may derive high-level (non-personal) metrics.

For more information, see [Docker Data Processing Agreement](https://www.docker.com/legal/data-processing-agreement).
