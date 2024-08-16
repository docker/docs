---
title: Docker for GitHub Copilot
description: Learn how to use the Docker Extension for GitHub Copilot to streamline Docker-related tasks.
keywords: Docker, GitHub Copilot, extension, Visual Studio Code, chat, ai, containerization
---

{{% restricted title="Early access" %}}

The Docker for GitHub Copilot extension is a part of the [Limited Public Beta](https://github.blog/changelog/2024-05-21-copilot-extensions-now-in-limited-public-beta/)
for GitHub Copilot Extensions. Only members of the beta program are eligible to use this extension.

{{% /restricted %}}

The Docker Extension for GitHub Copilot integrates Docker's capabilities with
GitHub Copilot, providing assistance with containerizing applications and
generating Docker assets. This extension helps you streamline Docker-related
tasks wherever GitHub Copilot Chat is available.

The key features are:

- Ask questions and receive responses about containerization in any context
  when GitHub Copilot Chat is available, such as on GitHub.com and in Visual Studio Code.
- Automatically generate Dockerfiles, Docker Compose files, and .dockerignore
  files for a project.
- Open pull requests with generated Docker assets directly from the chat
  interface.
- Get summaries of project vulnerabilities from Docker Scout and receive next
  steps via CLI.

## Supported languages

The Docker Extension for GitHub Copilot supports the following
programming languages:

- Dockerfile
- Java
- JavaScript
- Python
- TypeScript

## Get Started

Here's how to get started with the Docker Extension for GitHub Copilot.

### Prerequisites

Before you start, ensure that:

- You have been granted access to GitHub's [limited public beta program](https://github.blog/changelog/2024-05-21-copilot-extensions-now-in-limited-public-beta/).
- The [Docker for GitHub Copilot](https://github.com/marketplace/docker-for-github-copilot)
  app is installed in your GitHub organization.

### Setup

Select where you want to use the Docker Extension for GitHub Copilot:

{{< tabs >}}
{{< tab name="Editor or IDE" >}}

For instructions on how to use the Docker Extension for GitHub Copilot in
your editor, see:

- [Visual Studio Code](https://docs.github.com/en/copilot/github-copilot-chat/copilot-chat-in-ides/using-github-copilot-chat-in-your-ide?tool=vscode)
- [Visual Studio](https://docs.github.com/en/copilot/github-copilot-chat/copilot-chat-in-ides/using-github-copilot-chat-in-your-ide?tool=visualstudio)
- [Codespaces](https://docs.github.com/en/codespaces/reference/using-github-copilot-in-github-codespaces)

{{< /tab >}}
{{< tab name="GitHub.com" >}}

No setup is needed for the Docker Extension for GitHub Copilot on GitHub.com.
Simply go to any repository and start using the chat interface,
see [Using the extension](#using-the-extension).

{{< /tab >}}
{{< /tabs >}}

## Using the Extension

The Docker Extension for GitHub Copilot provides a chat interface that you can
use to interact with the Docker agent. You can ask questions and get help
Dockerizing your project.

The Docker agent is trained to understand Docker-related questions, and provide
guidance on Dockerfiles, Docker Compose files, and other Docker assets.

### Editor

1. Open your project in your editor.
2. Open the Copilot extension or plugin.
3. Interact with the Docker agent using the `/docker` slash prompt, followed by your question.

### GitHub.com

To use the Docker Extension in the GitHub web interface:

1. Go to [github.com](https://github.com/) and sign in to your account.
2. Go to any repository.
3. Select the Copilot logo in the site menu to open the chat interface.

   ![Copilot chat button](/copilot/copilot-button.png?w=400px)

4. Interact with the Docker agent by tagging `@docker`.

## Example prompts

Here are some examples of the types of questions you can ask the Docker agent:

### Ask general Docker questions

You can ask general question about Docker. For example:

- `@docker what is a Dockerfile`
- `@docker how do I build a Docker image`
- `@docker how do I run a Docker container`

### Ask questions about your project

You can ask questions about your project, such as:

- `@docker what is the best way to Dockerize this project`
- `@docker can you help me find vulnerabilities in my project?`

The Docker agent will analyze your project, generate the necessary files, and,
if applicable, offer to raise a pull request with the necessary Docker assets.

## Performing actions on your behalf

Before the agent performs any actions on your behalf, such as opening a pull
request for you, you're prompted to provide your consent to allow the
operation. You can always roll back or back out of the changes.

![Copilot action prompt](/copilot/copilot-action-prompt.png?w=400px)

In the event that the agent encounters an error, for example during PR
creation, it handles timeouts and lack of responses gracefully.

## Feedback

For issues or feedback, visit the [GitHub feedback repository](https://github.com/docker/copilot-issues).
