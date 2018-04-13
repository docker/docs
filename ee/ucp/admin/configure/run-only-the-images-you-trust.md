---
title: Run only the images you trust
description: Configure a Docker UCP cluster to only allow running applications that use images you trust.
keywords: ucp, dtr, security, trust
---

With Docker Universal Control Plane you can enforce applications to only use
Docker images signed by UCP users you trust. When a user tries to deploy an
application to the cluster, UCP checks if the application uses a Docker image
that is not trusted, and won’t continue with the deployment if that’s the case.

![Enforce image signing](../../images/run-only-the-images-you-trust-1.svg)

By signing and verifying the Docker images, you ensure that the images being
used in your cluster are the ones you trust and haven’t been altered either in
the image registry or on their way from the image registry to your UCP cluster.

## Example workflow

Here's an example of a typical workflow:

1. A developer makes changes to a service and pushes their changes to a version
   control system.
2. A CI system creates a build, runs tests, and pushes an image to DTR with the
   new changes.
3. The quality engineering team pulls the image and runs more tests. If
   everything looks good they sign and push the image.
4. The IT operations team deploys a service. If the image used for the service
   was signed by the QA team, UCP deploys it. Otherwise UCP refuses to deploy.

## Configure UCP

To configure UCP to only allow running services that use Docker images you
trust, go to the UCP web UI, navigate to the **Admin Settings** page, and in
the left pane, click **Docker Content Trust**.

Select the **Run Only Signed Images** option to only allow deploying
applications if they use images you trust.

![UCP settings](../../images/run-only-the-images-you-trust-2.png){: .with-border}

With this setting, UCP allows deploying any image as long as the image has
been signed. It doesn't matter who signed the image.

To enforce that the image needs to be signed by specific teams, click **Add Team**
and select those teams from the list.

![UCP settings](../../images/run-only-the-images-you-trust-3.png){: .with-border}

If you specify multiple teams, the image needs to be signed by a member of each
team, or someone that is a member of all those teams, again which must be a part
of the `docker-datacenter` organization.

> Signing with teams
>
> Teams used for signing policy enforcement must be in the `docker-datacenter`
> organization.

Click **Save** for UCP to start enforcing the policy. From now on, existing
services will continue running and can be restarted if needed, but UCP will only
allow deploying new services that use a trusted image.

## Where to go next

- [Sign and push images to DTR](/ee/dtr/user/manage-images/sign-images.md)
