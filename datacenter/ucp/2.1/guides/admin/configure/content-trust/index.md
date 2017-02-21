---
description: Configure a Docker UCP cluster to only allow running applications that use images you trust.
keywords: docker, ucp, backup, restore, recovery
title: Run only the images you trust
redirect_from:
  - /datacenter/ucp/2.1/guides/admin/configure/only-allow-running-signed-images/
---

## About trusted images

When transferring data among networked systems, _trust_ is a central concern. In
particular, when communicating over an untrusted medium such as the internet, it
is critical to ensure the integrity and the publisher of all the data a system
operates on. Docker allows you to push images to, and pull images from, public
and private registries.

Docker provides a mechanism called
[content trust](/engine/security/trust/content_trust.md), which you can use to
verify that the contents of the image have been approved by people you trust,
and to prevent untrusted images from being used in your UCP instance.

### Example workflow for using trusted images

An example workflow that takes advantage of content trust might look like this:

1.  Developers push code into source control.
2.  A CI system performs automated tests. If the tests pass, the CI system
    builds and cryptographically signs an image containing the code.
3.  A quality engineering team pulls the image signed by the CI system and
    performs quality tests on it. When the image is approved for production,
    part of the approval process is to cryptographically sign the image again.
4.  If any image is not signed both by the CI group and the QA group, UCP
    refuses to deploy it.

## Configuration overview

First, an administrator performs the following configuration tasks, which are
detailed in [Server-side tasks for content trust in UCP](admin_tasks.md).

1.  Configure UCP.

2.  Configure the Notary client on the administrator's system.

3.  Initialize the trusted repository in DTR.

4.  Delegate image signing to users in the correct groups.


Afterward, members of approved teams perform the following tasks, which are
detailed in [Configure the Docker client to sign images](client_configuration.md):

1.  Set up the Docker CLI to use the signing certificates from the UCP client
    bundle and to require images to be signed when pulling them from
    repositories.

2.  Sign and push an image to a repository.

## Next steps

- [Server-side tasks for content trust in UCP](admin_tasks.md)
- [Configure the Docker client to sign images](client_configuration.md)
