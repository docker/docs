---
description: Using your DockerID to log in to Docker Cloud
keywords:
- one, two, three
menu:
  main:
    parent: docker-cloud
    weight: -99
title: Docker ID and Settings
---

# Your Docker ID and Docker Cloud account

Docker Cloud uses your Docker ID for access and access control, and this allows
you to link your Hub and Cloud accounts.

If you already have an account on Docker Hub you can use the same credentials to
log in to Docker Cloud.

If you don't have an existing Docker ID, you can sign up for a new Docker ID
from the Cloud website, or using the `docker login` command in the Docker CLI.
The name you choose for your Docker ID becomes part of your account namespace.

## Manage cloud services and source providers

You can link to your own hosts, or to hosted nodes from a Cloud Services
Provider such as Amazon Web Services or Microsoft Azure from your Docker Cloud
account.

You can also link to source code repositories such as GitHub and
Bitbucket from your Docker Cloud account settings.

<!-- TODO:
## API keys
API keys are used for what?
-->

## Notifications

You can configure your account so that you receive email notifications for certain types of events in Docker Cloud.

You can also connect Slack to your Docker Cloud account so you can get notifications through your chat channels. To learn more, see [Docker Cloud notifications in Slack](slack-integration.md).

## Paid accounts

Like free Docker Hub accounts, free Docker Cloud accounts come with one free
private repository.

If you require more private repositories, visit your **Cloud settings** and
select **Plan** to see the available paid accounts.
