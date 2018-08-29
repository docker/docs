---
description: Using your DockerID to log in to Docker Cloud
keywords: one, two, three
title: Docker Cloud settings and Docker ID
---

Docker Cloud uses your Docker ID for access and access control, and this allows
you to link your Hub and Cloud accounts.

If you already have an account on Docker Hub, you can use the same credentials to
log in to Docker Cloud.

If you don't have a [Docker ID](../docker-id/) yet, you can sign up for one from
the Cloud website, or using the `docker login` command in the Docker CLI. The
name you choose for your Docker ID becomes part of your account namespace.

## Manage cloud services and source providers

You can link to your own hosts, or to hosted nodes from a Cloud Services
Provider such as Amazon Web Services or Microsoft Azure from your Docker Cloud
account.

You can also link to source code repositories such as GitHub and
Bitbucket from your Docker Cloud account settings.

## Email addresses

You can associate multiple email addresses with your Docker ID, and one of these
becomes the primary address for the account. The primary address is used by
Docker to send password reset notifications and other important information, so
be sure to keep it updated.

To add another email address to your Docker ID:

1. In Docker Cloud, click the user icon menu at top right, and click **Account Settings**.
2. In the **Emails** section, enter a new email address for the account.
3. Click the **plus sign** icon (**+**) to add the address and send a verification email.

The new email address is not added to the account until you confirm it by
clicking the link in the verification email. This link is only good for a
limited time. To send a new verification email, click the envelope icon next to
the email address that you want to verify.

If you have multiple verified email addresses associated with the account, you can click **Set as primary** to change the primary email address.

## Notifications

You can configure your account so that you receive email notifications for certain types of events in Docker Cloud.

You can also connect Slack to your Docker Cloud account so you can get notifications through your chat channels. To learn more, see [Docker Cloud notifications in Slack](slack-integration.md).

## Paid accounts

Like free Docker Hub accounts, free Docker Cloud accounts come with one free
private repository.

If you require more private repositories, visit your **Cloud settings** and
select **Plan** to see the available paid accounts.
