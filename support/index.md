---
description: Get support
keywords: Support, Docker Desktop, Docker Hub, Hub, Linux, Mac, Windows
title: Get support
redirect_from:
- /desktop/support/
---

Find information on how to get support, and the scope of Docker support.

{% include upgrade-cta.html
  body="Docker offers support for developers subscribed to a Pro, Team, or a Business tier. Upgrade now to benefit from Docker Support."
  header-text="This feature requires a paid Docker subscription"
  target-url="https://www.docker.com/pricing?utm_source=docker&utm_medium=webreferral&utm_campaign=docs_driven_upgrade_desktop_support"
%}

## How do I get support?

If you have a paid Docker subscription, raise a ticket through [Docker support](https://hub.docker.com/support/contact/){:target="_blank" rel="noopener" class="_"}.

All Docker users can seek support through the following resources, where we or the community respond on a best effort basis.
   - [Docker Desktop for Windows GitHub repo](https://github.com/docker/for-win){:target="_blank" rel="noopener" class="_"} 
   - [Docker Desktop for Mac GitHub repo](https://github.com/docker/for-mac){:target="_blank" rel="noopener" class="_"}
   - [Docker Desktop for Linux GitHub repo](https://github.com/docker/for-linux){:target="_blank" rel="noopener" class="_"}
   - [Docker Community Forums](https://forums.docker.com/){:target="_blank" rel="noopener" class="_"}
   - [Docker Community Slack](https://dockercommunity.slack.com/){:target="_blank" rel="noopener" class="_"}


## What support can I get?

If you have a paid Docker subscription, you can get support for the following types of issues:

   * Account management related issues
   * Automated builds
   * Basic product 'how to' questions
   * Billing or subscription issues
   * Configuration issues
   * Desktop installation issues
      * Installation crashes
      * Failure to launch Docker Desktop on first run
   * Desktop update issues
   * Login issues in both the command line interface and Docker Hub user interface
   * Push or pull issues, including rate limiting
   * Usage issues
      * Crash closing software
      * Docker Desktop not behaving as expected

   For Windows users, you can also request support on:
   * Enabling virtualization in BIOS
   * Enabling Windows features
   * Running inside [certain VM or VDI environments](../desktop/vm-vdi.md) (Docker Business customers only)


## What is not supported?

Docker excludes support for the following types of issues:
   * Altered or modified Docker software
   * Any version of the Docker software other than the latest version
   * Defects in the Docker software due to hardware malfunction, abuse, or improper use
   * Docker Support excludes training, customization, and integration
   * Features labeled as experimental
   * Reimbursing and expenses spent for third-party services not provided by Docker
   * Routine product maintenance (data backup, cleaning disk space and configuring log rotation)
   * Running containers of a different architecture using emulation
   * Running on unsupported operating systems, including beta or preview versions of operating systems
   * Scale deployment/multi-machine installation of Desktop
   * Support for Kubernetes
   * Support for the Docker Engine, Docker CLI, or other bundled Linux components
   * Supporting Desktop as a production runtime
   * System/Server administration activities
   * Third-party applications not provided by Docker
   * Use on or in conjunction with hardware or software other than that specified in the applicable documentation

## What Docker Desktop versions are supported?

For Docker Business customers, we offer support for versions up to six months older than the latest version, although any fixes will be on the latest version.

For Pro and Team customers, we only offer support for the latest version of Docker Desktop. If you are running an older version, you may be asked to update before we investigate your support request.

## How many machines can I get support for Docker Desktop on?

As a Pro user you can get support for Docker Desktop on a single machine.
As a Team, you can get support for Docker Desktop for the number of machines equal to the number of seats as part of your plan.

## What OS’s are supported?

Docker Desktop is available for Mac, Linux, and Windows. The supported version information can be found on the following pages:

* [Mac system requirements](../desktop/install/mac-install.md/#system-requirements)
* [Windows system requirements](../desktop/install/windows-install.md/#system-requirements)
* [Linux system requirements](../desktop/install/linux-install.md/#system-requirements)

## How is personal diagnostic data handled in Docker Desktop?

When uploading diagnostics to help Docker with investigating issues, the uploaded diagnostics bundle may contain personal data such as usernames and IP addresses. The diagnostics bundles are only accessible to Docker, Inc.
employees who are directly involved in diagnosing Docker Desktop issues.

By default, Docker, Inc. will delete uploaded diagnostics bundles after 30 days. You may also request the removal of a diagnostics bundle by either specifying the diagnostics ID or via your GitHub ID (if the diagnostics ID is mentioned in a GitHub issue). Docker, Inc. will only use the data in the diagnostics bundle to investigate specific user issues but may derive high-level (non personal) metrics such as the rate of issues from it.

For more information, see [Docker Data Processing Agreement](https://www.docker.com/legal/data-processing-agreement){: target="_blank" rel="noopener" class="_"}.

## What can I do before seeking support?

Before seeking support, you can perform basic troubleshooting. See [Diagnose and Troubleshooting](../desktop/troubleshoot/overview.md) for more information.

You can also see if an answer already exists in the following FAQs:
- [Docker Business or Team onboarding](../docker-hub/onboarding-faqs.md)
- [Docker Desktop](../desktop/faqs/general.md)
- [Docker Desktop for Linux](../desktop/faqs/linuxfaqs.md)
- [Docker Desktop for Mac](../desktop/faqs/macfaqs.md)
- [Docker Desktop for Windows](../desktop/faqs/windowsfaqs.md)
- [Single Sign-on](../single-sign-on/faqs.md)
