---
title: Using the Docker for GitHub Copilot extension
linkTitle: Usage
description: |
  Learn how to use the Docker for GitHub Copilot extension to interact with the
  Docker agent, get help Dockerizing projects, and ask Docker-related questions
  directly from your IDE or GitHub.com.
weight: 20
---

{{< summary-bar feature_name="Docker GitHub Copilot" >}}

The Docker Extension for GitHub Copilot provides a chat interface that you can
use to interact with the Docker agent. You can ask questions and get help
Dockerizing your project.

The Docker agent is trained to understand Docker-related questions, and provide
guidance on Dockerfiles, Docker Compose files, and other Docker assets.

## Setup

Before you can start interacting with the Docker agent, make sure you've
[installed](./install.md) the extension for your organization.

### Enable GitHub Copilot chat in your editor or IDE

For instructions on how to use the Docker Extension for GitHub Copilot in
your editor, see:

- [Visual Studio Code](https://docs.github.com/en/copilot/github-copilot-chat/copilot-chat-in-ides/using-github-copilot-chat-in-your-ide?tool=vscode)
- [Visual Studio](https://docs.github.com/en/copilot/github-copilot-chat/copilot-chat-in-ides/using-github-copilot-chat-in-your-ide?tool=visualstudio)
- [Codespaces](https://docs.github.com/en/codespaces/reference/using-github-copilot-in-github-codespaces)

### Verify the setup

You can verify that the extension has been properly installed by typing
`@docker` in the Copilot Chat window. As you type, you should see the Docker
agent appear in the chat interface.

![Docker agent in chat](images/docker-agent-copilot.png)

The first time you interact with the agent, you're prompted to sign in and
authorize the Copilot extension with your Docker account.

## Asking Docker questions in your editor

To interact with the Docker agent from within your editor or IDE:

1. Open your project in your editor.
2. Open the Copilot chat interface.
3. Interact with the Docker agent by tagging `@docker`, followed by your question.

## Asking Docker questions on GitHub.com

To interact with the Docker agent from the GitHub web interface:

1. Go to [github.com](https://github.com/) and sign in to your account.
2. Go to any repository.
3. Select the Copilot logo in the site menu, or select the floating Copilot widget, to open the chat interface.

   ![Copilot chat button](images/copilot-button.png?w=400px)

4. Interact with the Docker agent by tagging `@docker`, followed by your question.
