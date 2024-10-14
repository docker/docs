---
title: Example prompts for the Docker agent
linkTitle: Example prompts
description: |
  Discover example prompts to interact with the Docker agent and learn how to
  automate tasks like Dockerizing projects or opening pull requests.
weight: 30
---

{{% restricted title="Early Access" %}}
The Docker for GitHub Copilot extension is an [early access](/release-lifecycle#early-access-ea) product.
{{% /restricted %}}

Here are some examples of the types of questions you can ask the Docker agent:

### Ask general Docker questions

You can ask general question about Docker. For example:

- `@docker what is a Dockerfile?`
- `@docker how do I build a Docker image?`
- `@docker how do I run a Docker container?`

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

![Copilot action prompt](images/copilot-action-prompt.png?w=400px)

In the event that the agent encounters an error, for example during PR
creation, it handles timeouts and lack of responses gracefully.

## Feedback

For issues or feedback, visit the [GitHub feedback repository](https://github.com/docker/copilot-issues).
