---
description: Docker Desktop support
keywords: Support, Docker Desktop, Linux, Mac, Windows
title: Get support
---

Find information on how to get support, and the scope of Docker Desktop support.

{% include upgrade-cta.html
  body="Docker Desktop offers support for developers subscribed to a Pro, Team, or a Business tier. Upgrade now to benefit from Docker Support."
  header-text="This feature requires a paid Docker subscription"
  target-url="https://www.docker.com/pricing?utm_source=docker&utm_medium=webreferral&utm_campaign=docs_driven_upgrade_desktop_support"
%}

### How do I get Docker Desktop support?

If you have a paid Docker subscription, raise a ticket through [Docker Desktop support](https://hub.docker.com/support/desktop/){:target="_blank" rel="noopener" class="_"}.

Docker Community users can get support through our Github repos for-win, for-mac, for-linux, where we respond on a best-effort basis.

### What support can I get?

If you have a paid Docker subscription, you can request for support on the following types of issues:

* Desktop upgrade issues
* Desktop installation issues
  * Installation crashes
  * Failure to launch Docker Desktop on first run
* Usage issues
  * Crash closing software
  * Docker Desktop not behaving as expected
* Configuration issues
* Basic product ‘how to’ questions

For Windows users, you can also request support on:
* Enabling virtualization in BIOS
* Enabling Windows features

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
* Docker Support excludes training, customization, and integration

### What versions are supported?

We currently only offer support for the latest version of Docker Desktop. If you are running an older version, you may be asked to upgrade before we investigate your support request.

### How many machines can I get support for Docker Desktop on?

As a Pro user you can get support for Docker Desktop on a single machine.
As a Team, you can get support for Docker Desktop for the number of machines equal to the number of seats as part of your plan.

### What OS’s are supported?

Docker Desktop is available for Mac, Linux, and Windows. The supported version information can be found on the following pages:

* [Mac system requirements](install/mac-install.md#system-requirements)
* [Windows system requirements](install/windows-install.md#system-requirements)
* [Linux system requirements](install/linux-install.md#system-requirements)

### How is personal diagnostic data handled in Docker Desktop?

When uploading diagnostics to help Docker with investigating issues, the uploaded diagnostics bundle may contain personal data such as usernames and IP addresses. The diagnostics bundles are only accessible to Docker, Inc.
employees who are directly involved in diagnosing Docker Desktop issues.

By default, Docker, Inc. will delete uploaded diagnostics bundles after 30 days. You may also request the removal of a diagnostics bundle by either specifying the diagnostics ID or via your GitHub ID (if the diagnostics ID is mentioned in a GitHub issue). Docker, Inc. will only use the data in the diagnostics bundle to investigate specific user issues but may derive high-level (non personal) metrics such as the rate of issues from it.

For more information, see [Docker Data Processing Agreement](https://www.docker.com/legal/data-processing-agreement){: target="_blank" rel="noopener" class="_"}.