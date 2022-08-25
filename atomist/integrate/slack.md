---
title: Slack
description:
keywords:
---

{% include atomist/disclaimer.md %}

Atomist has a powerful Slack integration to help your team leverage the power of
ChatOps. After installing the bot, you'll need to link some channels to GitHub
repositories. The bot will only create notifications in channels that have
subscribed to them.

## Connecting Slack

Start by
[installing the GitHub Notification skill](https://go.atomist.com/catalog/skills/atomist/github-notifications-skill).

![start-enable](img/slack/0-skill-start-enable.png)

Click the **Continue** button on the right side.

1.  Connect to Slack.

    In the next step, you'll be redirected to Slack to install the Atomist
    application into your workspace. This step can only be performed by a Slack
    user that has permission to install new applications in the Slack workspace.

2.  You will be redirected to Slack.

    ![select Slack](img/slack/1-skill-choose-slack.png)

3.  Authorize the application.

    ![redirect to Slack](img/slack/2-redirect-to-slack.png)

4.  You will be redirected back to Atomist.

    After being redirected back to Atomist, you'll be asked to link a channel to
    one of your GitHub repositories.

    ![link channels](img/slack/3-link-channels.png)

5.  Select a **Slack Channel** on the left, and one or more **Repositories** on
    the right, then click the **Link** button.

    - You can always add, edit and remove channel links in the \*\*Manage
      > Integrations > Slack\*\* page.
    - The `@atomist` bot will automatically be invited to linked channels.

6.  Complete Installation

    ![set parameters](img/slack/4-set-params.png)

## Disconnecting Slack

You might want to disconnect Slack when:

- You want to change the Slack workspace that is connected to your Atomist
  workspace. To do so, disconnect the old Slack workspace first, then follow the
  instructions above for Connecting to Slack to connect the new workspace.
- You want to remove Atomist from a Slack workspace where you no longer need
  Atomist.

To disconnect your Atomist workspace from your Slack workspace:

1.  Visit **Manage > Integrations > Slack** and click the **Disconnect** button.
    This removes the Slack connection from your Atomist workspace.
2.  Go to the Slack [Manage Apps page](https://slack.com/apps/manage){:
    target="blank" rel="noopener" class=""}, find the Atomist app and select
    **Remove App**. This removes the Atomist Slack App from your Slack
    workspace.
