---
description: General Frequently asked questions
keywords: desktop, mac, windows, faqs
redirect_from:
- /mackit/faqs/
- /docker-for-mac/faqs/
- /docker-for-windows/faqs/
- /desktop/faqs/
title: Frequently asked questions
---

### What are the system requirements for Docker Desktop?

For information about Docker Desktop system requirements, see:
- [Mac system requirements](../install/mac-install.md#system-requirements) 
- [Windows system requirements](../install/windows-install.md#system-requirements)
- [Linux system requirements](../install/linux-install.md#system-requirements)

### Where does Docker Desktop get installed on my machine?

By default, Docker Desktop is installed at the following location:

- On Mac: `/Applications/Docker.app`
- On Windows: `C:\Program Files\Docker\Docker`
- On Linux: `/opt/docker-desktop`

### Where can I find the checksums for the download files?

You can find the checksums on the [release notes](../release-notes.md) page.

### Do I need to pay to use Docker Desktop?

Docker Desktop remains free for small businesses (fewer than 250 employees AND less than $10 million in annual revenue), personal use, education, and non-commercial open-source projects. It requires a paid subscription for professional use in larger enterprises.
The effective date of these terms is August 31, 2021. When downloading and installing Docker Desktop, you are asked to agree to the [Docker Subscription Service Agreement](https://www.docker.com/legal/docker-subscription-service-agreement){: target="_blank" rel="noopener" class="_"}.

Read the [Blog](https://www.docker.com/blog/updating-product-subscriptions/){: target="_blank" rel="noopener" class="_" id="dkr_docs_subscription_btl"} and [FAQs](https://www.docker.com/pricing/faq){: target="_blank" rel="noopener" class="_" id="dkr_docs_subscription_btl"} to learn how companies using Docker Desktop may be affected. For information about Docker Desktop licensing, see [Docker Desktop License Agreement](../../subscription/index.md#docker-desktop-license-agreement).

### Can I use Docker Desktop offline?

Yes, you can use Docker Desktop offline. However, you
cannot access features that require an active internet
connection. Additionally, any functionality that requires you to sign won't work while using Docker Desktop offline or in air-gapped environments.
This includes:

- The in-app [Quick Start Guide](../get-started.md#quick-start-guide)
- Pulling or pushing an image to Docker Hub
- [Image Access Management](../../docker-hub/image-access-management.md)
- [Vulnerability scanning](../../docker-hub/vulnerability-scanning.md)
- Viewing remote images in the [Docker Dashboard](../dashboard.md)
- Settting up [Dev Environments](../dev-environments/index.md)
- Docker build when using [Buildkit](../../develop/develop-images/build_enhancements.md). You can work around this by disabling
  BuildKit. Run `DOCKER_BUILDKIT=0 docker build .` to disable BuildKit.
- Deploying an app to the cloud through Compose
  [ACI](../../cloud/aci-integration.md) and [ECS](../../cloud/ecs-integration.md)
  integrations
- [Kubernetes](../kubernetes.md) (Images are download when you enable Kubernetes for the first time)
- Check for updates
- [In-app diagnostics](../troubleshoot/overview.md#diagnose-from-the-app) (including the [Self-diagnose tool](../troubleshoot/overview.md#diagnose-from-the-app))
- Tip of the week
- Sending usage statistics

### What is an experimental feature?

{% include experimental.md %}

### Where can I find information about diagnosing and troubleshooting Docker Desktop issues?

You can find information about diagnosing and troubleshooting common issues in the [Troubleshooting topic](../troubleshoot/overview.md).

If you do not find a solution in troubleshooting, browse the Github repositories or create a new issue:

- [docker/for-mac](https://github.com/docker/for-mac/issues){: target="_blank" rel="noopener" class="_"} - - [docker/for-win](https://github.com/docker/for-win/issues){: target="_blank" rel="noopener" class="_"}
- [docker/for-linux](https://github.com/docker/for-linux/issues){: target="_blank" rel="noopener" class="_"}

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

For more information and examples, see [how to connect from a container to a service on the host](../networking.md#i-want-to-connect-from-a-container-to-a-service-on-the-host).

### Can I pass through a USB device to a container?

Unfortunately, it is not possible to pass through a USB device (or a
serial port) to a container as it requires support at the hypervisor level.

### Can I run Docker Desktop in nested virtualization scenarios?

Docker Desktop can run inside a Windows 10 VM running on apps like Parallels or
VMware Fusion on a Mac provided that the VM is properly configured. However,
problems and intermittent failures may still occur due to the way these apps
virtualize the hardware. For these reasons, **Docker Desktop is not supported in
nested virtualization scenarios**. It might work in some cases and not in others.

### Docker Desktop's UI appears green, distorted, or has visual artifacts. How do I fix this?

Docker Desktop uses hardware-accelerated graphics by default, which may cause problems for some GPUs. In such cases,
Docker Desktop will launch successfully, but some screens may appear green, distorted,
or have some visual artifacts.

To work around this issue, disable hardware acceleration by creating a `"disableHardwareAcceleration": true` entry in Docker Desktop's `settings.json` file. You can find this file at:

- **Mac**: `~/Library/Group Containers/group.com.docker/settings.json`
- **Windows**: `C:\Users\[USERNAME]\AppData\Roaming\Docker\settings.json`

After updating the `settings.json` file, close and restart Docker Desktop to apply the changes.

### Can I run Docker Desktop on Virtualized hardware?

No, currently this is unsupported and against terms of use. 

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
