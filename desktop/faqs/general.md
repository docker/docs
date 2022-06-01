---
description: Frequently asked questions
keywords: desktop, mac, windows, faqs
redirect_from:
- /mackit/faqs/
- /docker-for-mac/faqs/
- /docker-for-windows/faqs/
- /desktop/faqs/
title: Frequently asked questions
---

## General

### What are the system requirements for Docker Desktop?

For information about Docker Desktop system requirements, see [Docker Desktop for Mac system requirements](../mac/install.md#system-requirements) and [Docker Desktop for Windows system requirements](../windows/install.md#system-requirements).

### Where does Docker Desktop get installed on my machine?

By default, Docker Desktop is installed at the following location:

- On Mac: `/Applications/Docker.app`
- On Windows: `C:\Program Files\Docker\Docker`
- On Linux: `/opt/docker-desktop`

### Do I need to pay to use Docker Desktop?

Docker Desktop remains free for small businesses (fewer than 250 employees AND less than $10 million in annual revenue), personal use, education, and non-commercial open-source projects. It requires a paid subscription for professional use in larger enterprises.
The effective date of these terms is August 31, 2021. When downloading and installing Docker Desktop, you are asked to agree to the [Docker Subscription Service Agreement](https://www.docker.com/legal/docker-subscription-service-agreement){: target="_blank" rel="noopener" class="_"}.

Read the [Blog](https://www.docker.com/blog/updating-product-subscriptions/){: target="_blank" rel="noopener" class="_" id="dkr_docs_subscription_btl"} and [FAQs](https://www.docker.com/pricing/faq){: target="_blank" rel="noopener" class="_" id="dkr_docs_subscription_btl"} to learn how companies using Docker Desktop may be affected. For information about Docker Desktop licensing, see [Docker Desktop License Agreement](../../subscription/index.md#docker-desktop-license-agreement).

### Can I use Docker Desktop offline?

Yes, you can use Docker Desktop offline. However, you
cannot access features that require an active internet
connection. Additionally, any functionality that requires you to sign won't work while using Docker Desktop offline or in air-gapped environments.
This includes:

- The in-app [Quick Start Guide](../mac/install.md#quick-start-guide)
- Pulling or pushing an image to Docker Hub
- [Image Access Management](../../docker-hub/image-access-management.md)
- [Vulnerability scanning](../../docker-hub/vulnerability-scanning.md)
- Viewing remote images in the [Docker Dashboard](../dashboard.md)
- Settting up [Dev Environments](../dev-environments.md)
- Docker build when using [Buildkit](../../develop/develop-images/build_enhancements.md). You can work around this by disabling
  BuildKit. Run `DOCKER_BUILDKIT=0 docker build .` to disable BuildKit.
- Deploying an app to the cloud through Compose
  [ACI](../../cloud/aci-integration.md) and [ECS](../../cloud/ecs-integration.md)
  integrations
- [Kubernetes](../kubernetes.md) (Images are download when you enable Kubernetes for the first time)
- [Check for updates](../mac/install.md#updates) (manual and automatic)
- [In-app diagnostics](../mac/troubleshoot.md#diagnose-and-feedback) (including the [Self-diagnose tool](../mac/troubleshoot.md#self-diagnose-tool))
- Tip of the week
- Sending usage statistics

### What is an experimental feature?

{% include experimental.md %}

### Where can I find information about diagnosing and troubleshooting Docker Desktop issues?

You can find information about diagnosing and troubleshooting common issues in the Troubleshooting topic. See:
- [Mac Logs and Troubleshooting](../mac/troubleshoot.md)
- [Windows Logs and Troubleshooting](../windows/troubleshoot.md)
- [Linux logs and Troubleshooting](../linux/troubleshoot.md)

If you do not find a solution in Troubleshooting, browse issues on
[docker/for-mac](https://github.com/docker/for-mac/issues){: target="_blank" rel="noopener" class="_"} or [docker/for-win](https://github.com/docker/for-win/issues){: target="_blank" rel="noopener" class="_"} GitHub repository, or create a new one.

### How do I connect to the remote Docker Engine API?

To connect to the remote Engine API, you might need to provide the location of the Engine API for Docker clients and development tools.

Mac and Windows WSL 2 users can connect to the Docker Engine through a Unix socket: `unix:///var/run/docker.sock`.

If you are working with applications like [Apache Maven](https://maven.apache.org/){: target="_blank" rel="noopener" class="_"}
that expect settings for `DOCKER_HOST` and `DOCKER_CERT_PATH` environment
variables, specify these to connect to Docker instances through Unix sockets.

For example:

```console
$ export DOCKER_HOST=unix:///var/run/docker.sock
```

Docker Desktop Windows users can connect to the Docker Engine through a **named pipe**: `npipe:////./pipe/docker_engine`, or **TCP socket** at this URL:
`tcp://localhost:2375`.

For details, see [Docker Engine API](../../engine/api/index.md).

### How do I connect from a container to a service on the host?

Mac, Linux, and Windows have a changing IP address (or none if you have no network access). On both Mac and Windows, we recommend that you connect to the special DNS name `host.docker.internal`, which resolves to the internal IP address used by the host. This is for development purposes and does not work in a production environment outside of Docker Desktop.

For more information and examples, see how to connect from a container to a service on the host
[on Mac](../mac/networking.md#i-want-to-connect-from-a-container-to-a-service-on-the-host) and [on Windows](../windows/networking.md#i-want-to-connect-from-a-container-to-a-service-on-the-host) or [on Linux](../linux/networking.md#i-want-to-connect-from-a-container-to-a-service-on-the-host).

### How do I connect to a container from Mac or Windows?

We recommend that you publish a port, or connect from another container. Port forwarding works for `localhost`; `--publish`, `-p`, or `-P` all work.

For more information and examples, see:
- [I want to connect to a container from Mac](../mac/networking.md#i-want-to-connect-to-a-container-from-the-mac) 
- [I want to connect to a container from Windows](../windows/networking.md#i-want-to-connect-to-a-container-from-windows)
- [I want to connect to a container from Linux](../linux/networking.md#i-want-to-connect-from-a-container-to-a-service-on-the-host)]

### Can I pass through a USB device to a container?

Unfortunately, it is not possible to pass through a USB device (or a
serial port) to a container as it requires support at the hypervisor level.

### Can I run Docker Desktop in nested virtualization scenarios?

Docker Desktop can run inside a Windows 10 VM running on apps like Parallels or
VMware Fusion on a Mac provided that the VM is properly configured. However,
problems and intermittent failures may still occur due to the way these apps
virtualize the hardware. For these reasons, **Docker Desktop is not supported in
nested virtualization scenarios**. It might work in some cases and not in others.

For more information, see [Running Docker Desktop in nested virtualization scenarios](../windows/troubleshoot.md#running-docker-desktop-in-nested-virtualization-scenarios).

### Docker Desktop's UI appears green, distorted, or has visual artifacts. How do I fix this?

Docker Desktop uses hardware-accelerated graphics by default, which may cause problems for some GPUs. In such cases,
Docker Desktop will launch successfully, but some screens may appear green, distorted,
or have some visual artifacts.

To work around this issue, disable hardware acceleration by creating a `"disableHardwareAcceleration": true` entry in Docker Desktop's `settings.json` file. You can find this file at:

- **Mac**: `~/Library/Group Containers/group.com.docker/settings.json`
- **Windows**: `C:\Users\[USERNAME]\AppData\Roaming\Docker\settings.json`

After updating the `settings.json` file, close and restart Docker Desktop to apply the changes.

## Releases

### How do new users install Docker Desktop?

Each Docker Desktop release is also delivered as a full installer for new users. The same applies if you have skipped a version, although this doesn't normally happen as updates are applied automatically.

### How frequent will new releases be?

New releases are available roughly monthly, unless there are critical fixes that need to be released sooner.

### How do I ensure that all users on my team are using the same version?

Previously you had to manage this yourself. Now, it happens automatically as a side effect of all users being on the latest version.

### My colleague has got a new version but I haven’t got it yet.

Sometimes we may roll out a new version gradually over a few days. Therefore, if you wait, it will turn up soon. Alternatively, you can select **Check for Updates** from the Docker menu to jump the queue and get the latest version immediately.

### Where can I find information about Stable and Edge releases?

Starting with Docker Desktop 3.0.0, Stable and Edge releases are combined into a single, cumulative release stream for all users. 

## Support

### Does Docker Desktop offer support?

Yes, Docker Desktop offers support for users with a paid Docker subscription.

For information about Docker subscriptions and to upgrade your existing account, see [Docker pricing](https://www.docker.com/pricing){: target="_blank" rel="noopener" class="_"}.

### What kind of feedback are you looking for?

Everything is fair game. We'd like your impressions on the download-install
process, startup, functionality available, the GUI, usefulness of the app,
command line integration, and so on. Tell us about the issues you are experiencing, what you like, or request a new feature through our public [Docker Roadmap](https://github.com/docker/roadmap){: target="_blank" rel="noopener" class="_"}.

### How is personal data handled in Docker Desktop?

When uploading diagnostics to help Docker with investigating issues, the uploaded diagnostics bundle may contain personal data such as usernames and IP addresses. The diagnostics bundles are only accessible to Docker, Inc.
employees who are directly involved in diagnosing Docker Desktop issues.

By default, Docker, Inc. will delete uploaded diagnostics bundles after 30 days. You may also request the removal of a diagnostics bundle by either specifying the diagnostics ID or via your GitHub ID (if the diagnostics ID is mentioned in a GitHub issue). Docker, Inc. will only use the data in the diagnostics bundle to investigate specific user issues but may derive high-level (non personal) metrics such as the rate of issues from it.

For more information, see [Docker Data Processing Agreement](https://www.docker.com/legal/data-processing-agreement){: target="_blank" rel="noopener" class="_"}.




