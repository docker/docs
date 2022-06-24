---
description: Dev Environments
keywords: Dev Environments, share, collaborate, local
title: Development Environments (Beta)
---
## Share your Dev Environment

{% include upgrade-cta.html
  body="Docker Pro, Team, and Business users can now share Dev Environments with their team members."
  header-text="This feature requires a paid Docker subscription"
  target-url="https://www.docker.com/pricing?utm_source=docker&utm_medium=webreferral&utm_campaign=docs_driven_upgrade"
%}

When you are ready to share your environment, just click the **Share** button and specify the Docker Hub namespace where you’d like to push your Dev Environment to.

![Share a Dev environment](../images/dev-env-share.png){:width="700px"}

This creates a Docker image of your dev environment, uploads it to the Docker Hub namespace you have specified in the previous step, and provides a tiny URL which you can use to share your work with your team members.

![Dev environment shared](../images/dev-env-shared.png){:width="700px"}

Your team members need to open the **Create** dialog, select the **Existing Dev Environment** tab, and then paste the URL. Your Dev Environment now starts in the exact same state as you shared it.

Using this shared Dev Environment, your team members can access the code, any dependencies, and the current Git branch you are working on. They can also review your changes and provide feedback even before you create a pull request!
