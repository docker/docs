---
title: Hardened Desktop
description: Overview of what Hardened Desktop is
keywords: security, hardened desktop, enhanced container isolation,
---

Hardened Desktop is Docker's ongoing effort to increase Docker Desktop security without impacting the developer experience

These new features from Docker follow a number of security related acquisitions (Nestybox, Atomist) that sees the company establishing itself as the undoubted market leader in providing an Enterprise-ready offering for containerised development.

This configuration is designed for organizations that don't give root/admin access to their developers on their machines, and wish to configure Docker Desktop to also be within the organization's centralized control.

In order to use this model, the application and VM image need to be installed as root/admin so that the user cannot modify them. All containers are run unprivileged in the VM, in user namespaces. Root access to the VM is removed, and privileged containers cannot be run, and there is no access to the host namespaces. The ownership boundary of system code in the VM moves to the organization. The user owns the (unprivileged) containers that they run, the equivalent of being able to run unprivileged applications on the host but not being able to modify the host configuration.


We have introduced features such as registry access management, that controls which registries a user can pull from on Docker Desktop, as organizations want to only allow users to pull from their central repository, but again this cannot actually be enforced if the user can modify the VM freely and disable controls, which this model prevents.


We have some longer term roadmap items around secure boot and code verification to increase trust in the code on the VM, as well as supporting trusted logging and audit.

 a new security model for Docker Desktop. The Hardened Desktop security model is designed to provide Enterprise admins with a simple and powerful way to increase their security posture for containerised development.


 As part of the Hardened Desktop model, Docker announced the release of two initial features. The first is Enhanced Container Isolation, a setting that helps admins to instantly enhance security by preventing containers from running as root in Docker Desktop’s Linux VM. The second is Admin Controls, which helps Enterprise admins to confidently manage and control usage of Docker Desktop. With just a few lines of JSON, admins will be able to enforce preferences like HTTP proxies, Network settings and the Docker Engine configuration, saving them significant time and cost in securing their developer workflows.

With the Hardened Desktop security model, and our new Enhanced Container Isolation and Admin Controls features we’re moving the ownership boundary for containers to the organization, meaning that any security controls admins set cannot be altered by the user.

 Docker will be adding more security enhancements to their Hardened Desktop model over the coming months.
