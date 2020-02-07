---
title: Run only the images you trust
description: Configure a Docker UCP cluster to only allow running applications that use images you trust.
keywords: ucp, dtr, security, trust
---

>{% include enterprise_label_shortform.md %}

With Docker Universal Control Plane (UCP), you can enforce applications to only use Docker images signed by UCP users you trust. Each time a user attempts to deploy an application to the cluster, UCP checks whether the application is using a trusted Docker image (and will halt the deployment if that is not the case).

![Enforce image signing](../../images/run-only-the-images-you-trust-1.svg)

By signing and verifying Docker images, you ensure that:

* The images used in your cluster are ones that you trust.
* The images have not been altered either in the image registry or on their way from the image registry to your UCP cluster.

## Example workflow

1. A developer makes changes to a service and pushes the changes to a version control system.
2. A continuous integration (CI) system creates a build, runs tests, and pushes an image to Docker Trusted Registry (DTR) with the new changes.
3. The quality engineering team pulls the image and runs more tests. If the tests are successful, the team signs and then pushes the image.
4. The IT operations team deploys a service. If the image used for the service was signed by the QA team, UCP deploys it. Otherwise, UCP refuses to deploy the image.

## Configure UCP

To configure UCP to only allow running services that use Docker trusted images:

1. Access the UCP UI and browse to the **Admin Settings** page.
2. In the left navigation pane, click **Docker Content Trust**.
3. Select the **Run only signed images** option.

    ![UCP settings](../../images/run-only-the-images-you-trust-2.png){: .with-border}

    With this setting, UCP allows deploying any image as long as the image has
    been signed. 

    To enforce the requirement that the image be signed by specific teams, click **Add Team** and select the pertinent teams from the list.

    ![UCP settings](../../images/run-only-the-images-you-trust-3.png){: .with-border}

    If you specify multiple teams, the image needs to be signed by a member of each
    team, or by someone that is a member of all of those teams.

4. Click **Save.** 

    At this point, UCP starts enforcing the policy. Existing services will continue running and can be restarted if needed, however UCP only allows the deployment of new services that use a trusted image.

## Where to go next

- [Sign and push images to DTR](/ee/dtr/user/manage-images/sign-images.md)
