---
title: Platform FAQs
linkTitle: Platform
description: Frequently asked questions about Docker platform security, containers, networking, and VMs.
keywords: Docker security, FAQs, authentication, vulnerability reporting, session management, container security, docker desktop isolation, enhanced container isolation, file sharing, docker desktop networking, virtualization, hyper-v, wsl2, network security, firewall
weight: 40
tags: [FAQ]
toc_max: 3
aliases:
  - /faq/security/general/
  - /security/faqs/general/
  - /faq/security/containers/
  - /security/faqs/containers/
  - /faq/security/networking-and-vms/
  - /security/faqs/networking-and-vms/
---

## General

### How do I report a vulnerability?

If you've discovered a security vulnerability in Docker, report it responsibly to security@docker.com so Docker can quickly address it.

### Does Docker lockout users after failed sign-ins?

Docker Hub locks out users after 10 failed sign-in attempts within 5 minutes. The lockout duration is 5 minutes. This policy applies to Docker Hub, Docker Desktop, and Docker Scout authentication.

### Do you support physical multi-factor authentication (MFA) with YubiKeys?

You can configure physical multi-factor authentication (MFA) through SSO using your identity provider (IdP). Check with your IdP if they support physical MFA devices like YubiKeys.

### How are sessions managed and do they expire?

Docker uses tokens to manage user sessions with different expiration periods:

- Docker Desktop: Signs you out after 90 days, or 30 days of inactivity
- Docker Hub and Docker Home: Sign you out after 24 hours

Docker also supports your IdP's default session timeout through SAML attributes. For more information, see [SSO attributes](/manuals/security/provisioning/_index.md#sso-attributes).

### How does Docker distinguish between employee users and contractor users?

Organizations use verified domains to distinguish user types. Team members with email domains other than verified domains appear as "Guest" users in the organization.

### How long are activity logs available?

Docker activity logs are available for 90 days. You're responsible for exporting logs or setting up drivers to send logs to your internal systems for longer retention.

### Can I export a list of users with their roles and privileges?

Yes, use the [Export Members](/manuals/accounts/organization/manage/members.md#export-members-csv-file) feature to export a CSV file containing your organization's users with role and team information.

### How do I remove users who aren't part of my IdP when using SSO without SCIM?

If SCIM isn't turned on, you must manually remove users from the organization. SCIM can automate user removal, but only for users added after SCIM is turned on. Users added before SCIM was turned on must be removed manually.

For more information, see [Manage organization members](/manuals/accounts/organization/manage/members.md).

### What metadata does Scout collect from container images?

For information about metadata stored by Docker Scout, see [Data handling](/manuals/scout/deep-dive/data-handling.md).

### How are Marketplace extensions vetted for security?

Security vetting for extensions isn't implemented. Extensions aren't covered as part of Docker's Third-Party Risk Management Program.

### Can I prevent users from pushing images to Docker Hub private repositories?

No direct setting exists to disable private repositories. However, [Registry Access Management](/manuals/enterprise/security/hardened-desktop/registry-access-management.md) lets administrators control which registries developers can access through Docker Desktop via Docker Home.

## Docker Desktop

### How does Docker Desktop handle authentication information?

Docker Desktop uses the host operating system's secure key management to store authentication tokens:

- macOS: [Keychain](https://support.apple.com/guide/security/keychain-data-protection-secb0694df1a/web)
- Windows: [Security and Identity API via Wincred](https://learn.microsoft.com/en-us/windows/win32/api/wincred/)
- Linux: [Pass](https://www.passwordstore.org/).

### Containers

#### How are containers isolated from the host in Docker Desktop?

Docker Desktop runs all containers inside a customized Linux virtual machine (except for native Windows containers). This adds strong isolation between containers and the host machine, even when containers run as root.

Important considerations include:

- Containers have access to host files configured for file sharing via Docker Desktop settings
- Containers run as root with limited capabilities inside the Docker Desktop VM by default
- Privileged containers (`--privileged`, `--pid=host`, `--cap-add`) run with elevated privileges inside the VM, giving them access to VM internals and Docker Engine

With Enhanced Container Isolation turned on, each container runs in a dedicated Linux user namespace inside the Docker Desktop VM. Even privileged containers only have privileges within their container boundary, not the VM. ECI uses advanced techniques to prevent containers from breaching the Docker Desktop VM and Docker Engine.

#### Which portions of the host filesystem can containers access?

Containers can only access host files that are:

1. Shared using Docker Desktop settings
1. Explicitly bind-mounted into the container (e.g., `docker run -v /path/to/host/file:/mnt`)

#### Can containers running as root access admin-owned files on the host?

No. Host file sharing uses a user-space file server (running in `com.docker.backend` as the Docker Desktop user), so containers can only access files that the Docker Desktop user already has permission to access.

### Networking and VMs

#### How can I limit container internet access?

Docker Desktop doesn't have a built-in mechanism for this, but you can use process-level firewalls on the host. Apply rules to the `com.docker.vpnkit` user-space process to control where it can connect (DNS allowlists, packet filters) and which ports/protocols it can use.

For enterprise environments, consider [Air-gapped containers](/manuals/enterprise/security/hardened-desktop/air-gapped-containers.md) which provide network access controls for containers.

#### Can I apply firewall rules to container network traffic?

Yes. Docker Desktop uses a user-space process (`com.docker.vpnkit`) for network connectivity, which inherits constraints like firewall rules, VPN settings, and HTTP proxy properties from the user that launched it.

#### Does Docker Desktop for Windows with Hyper-V allow users to create other VMs?

No. The `DockerDesktopVM` name is hard-coded in the service, so you cannot use Docker Desktop to create or manipulate other virtual machines.

#### How does Docker Desktop achieve network isolation with Hyper-V and WSL 2?

Docker Desktop uses the same VM processes for both WSL 2 (in the `docker-desktop` distribution) and Hyper-V (in `DockerDesktopVM`). Host/VM communication uses `AF_VSOCK` hypervisor sockets (shared memory) rather than network switches or interfaces. All host networking is performed using standard TCP/IP sockets from the `com.docker.vpnkit.exe` and `com.docker.backend.exe` processes.

For more information, see [How Docker Desktop networking works under the hood](https://www.docker.com/blog/how-docker-desktop-networking-works-under-the-hood/).
