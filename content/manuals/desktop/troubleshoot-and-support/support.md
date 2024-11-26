---
description: See what support is available for Docker Desktop
keywords: Support, Docker Desktop, Linux, Mac, Windows
title: Get support
weight: 20
aliases:
 - /desktop/support/
 - /support/
---

> [!NOTE]
> 
> Docker Desktop offers support for developers with a [Pro, Team, or Business subscription](https://www.docker.com/pricing?utm_source=docker&utm_medium=webreferral&utm_campaign=docs_driven_upgrade_desktop_support).

### How do I get Docker Desktop support?

If you have a paid Docker subscription, you can [contact the Support team](https://hub.docker.com/support/contact/).

All Docker users can seek support through the following resources, where Docker or the community respond on a best effort basis.
   - [Docker Desktop for Windows GitHub repo](https://github.com/docker/for-win) 
   - [Docker Desktop for Mac GitHub repo](https://github.com/docker/for-mac)
   - [Docker Desktop for Linux GitHub repo](https://github.com/docker/desktop-linux)
   - [Docker Community Forums](https://forums.docker.com/)
   - [Docker Community Slack](http://dockr.ly/comm-slack)

### What support can I get?

* Account management related issues
   * Automated builds
   * Basic product 'how to' questions
   * Billing or subscription issues
   * Configuration issues
   * Desktop installation issues
      * Installation crashes
      * Failure to launch Docker Desktop on first run
   * Desktop update issues
   * Sign-in issues in both the command line interface and Docker Hub user interface
   * Push or pull issues, including rate limiting
   * Usage issues
      * Crash closing software
      * Docker Desktop not behaving as expected

   For Windows users, you can also request support on:
   * Turning on virtualization in BIOS
   * Turning on Windows features
   * Running inside [certain VM or VDI environments](/manuals/desktop/setup/vm-vdi.md) (Docker Business customers only)

### What is not supported?

Docker Desktop excludes support for the following types of issues:

* Use on or in conjunction with hardware or software other than that specified in the applicable documentation
* Running on unsupported operating systems, including beta/preview versions of operating systems
* Running containers of a different architecture using emulation
* Support for the Docker engine, Docker CLI, or other bundled Linux components
* Support for Kubernetes
* Features labeled as experimental
* System/Server administration activities
* Supporting Desktop as a production runtime
* Scale deployment/multi-machine installation of Desktop
* Routine product maintenance (data backup, cleaning disk space and configuring log rotation)
* Third-party applications not provided by Docker
* Altered or modified Docker software
* Defects in the Docker software due to hardware malfunction, abuse, or improper use
* Any version of the Docker software other than the latest version
* Reimbursing and expenses spent for third-party services not provided by Docker
* Docker support excludes training, customization, and integration
* Running multiple instances of Docker Desktop on a single machine

> [!NOTE]
>
> Support for [running Docker Desktop in a VM or VDI environment](/manuals/desktop/setup/vm-vdi.md) is only available to Docker Business customers.

### What versions are supported?

For Docker Business customers, Docker offers support for versions up to six months older than the latest version, although any fixes will be on the latest version.

For Pro and Team customers, Docker only offers support for the latest version of Docker Desktop. If you are running an older version, Docker may ask you to update before investigating your support request.

### How many machines can I get support for Docker Desktop on?

As a Pro user you can get support for Docker Desktop on a single machine.
As a Team, you can get support for Docker Desktop for the number of machines equal to the number of seats as part of your plan.

### What OS’s are supported?

Docker Desktop is available for Mac, Linux, and Windows. The supported version information can be found on the following pages:

* [Mac system requirements](/manuals/desktop/setup/install/mac-install.md#system-requirements)
* [Windows system requirements](/manuals/desktop/setup/install/windows-install.md#system-requirements)
* [Linux system requirements](/manuals/desktop/setup/install/linux/_index.md#system-requirements)

### How is personal diagnostic data handled in Docker Desktop?

When uploading diagnostics to help Docker with investigating issues, the uploaded diagnostics bundle may contain personal data such as usernames and IP addresses. The diagnostics bundles are only accessible to Docker, Inc.
employees who are directly involved in diagnosing Docker Desktop issues.

By default, Docker, Inc. will delete uploaded diagnostics bundles after 30 days. You may also request the removal of a diagnostics bundle by either specifying the diagnostics ID or via your GitHub ID (if the diagnostics ID is mentioned in a GitHub issue). Docker, Inc. will only use the data in the diagnostics bundle to investigate specific user issues but may derive high-level (non personal) metrics such as the rate of issues from it.

For more information, see [Docker Data Processing Agreement](https://www.docker.com/legal/data-processing-agreement).
