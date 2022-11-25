---
description: Integrate Docker Hub with Slack
keywords: Slack, integrate, notifications
redirect_from:
- /docker-cloud/tutorials/slack-integration/
- /docker-cloud/slack-integration/
title: Set up Docker Hub notifications in Slack
sitemap: false
---

Docker Hub can integrate with your **Slack** team to provide notifications about builds.

## Set up a Slack integration

Before you begin, make sure that you are signed into the Slack team that you want to show notifications in.

1. Log in to the Docker account that owns the builds that you want to receive notifications about.

    > **Note**: If you are setting up notifications for an organization, log in as a member of the organization's `Owners` team, then switch to the organization account to change the settings.

2. Click **Account Settings** in the left hand navigation, and scroll down to the **Notifications** section.

3. Click the plug icon next to **Slack**.

    The Docker Hub page refreshes to show a Slack authorization screen.

4. On the page that appears, double check that you're signed in to the correct Slack team. (If necessary sign in to the correct one.)
5. Select the channel that should receive notifications.
6. Click **Authorize**.

    Once you click **Authorize**, you should see a message in the Slack channel notifying you of the new integration.

Once configured, choose a notification level:

* **Off** Do not receive any notifications.
* **Only failures** Only receive notifications about failed builds.
* **Everything** Receive notifications for both failed and successful builds.
  ![slack notifications](images/slack-notification-updates.png)

Enjoy your new Slack channel integration!

## Edit a Slack integration

* Click **Account Settings** in the lower left, scroll down to **Notifications**, and locate the **Slack** section. From here you can choose a new notification level, or remove the integration.

* From the Slack **Notifications** section you can also change the channel that the integration posts to. Click the reload icon (two arrows) next to the Slack integration to reopen the OAuth channel selector.

* Alternately, go to the <a href="https://slack.com/apps/manage" target="_blank" rel="noopener">Slack App Management page</a> and search for "Docker Hub". Click the result to see all of the Docker Hub notification channels set for the Slack team.
