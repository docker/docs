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

## Ask general Docker questions

You can ask general question about Docker. For example:

- `@docker what is a Dockerfile?`
- `@docker how do I build a Docker image?`
- `@docker how do I run a Docker container?`

## Get help containerizing your project

You can ask the agent to help you containerize your existing project:

- `@docker can you help create a compose file for this project?`
- `@docker can you create a Dockerfile for this project?`

The Docker agent will analyze your project, generate the necessary files, and,
if applicable, offer to [raise a pull request](#performing-actions-on-your-behalf)
with the necessary Docker assets.

## Analyze a project for vulnerabilities

The agent can help you improve your security posture with [Docker
Scout](/manuals/scout/_index.md):

- `@docker can you help me find vulnerabilities in my project?`
- `@docker does my project contain any insecure dependencies?`

The agent will run use Docker Scout to analyze your project's dependencies, and
report whether you're vulnerable to any [known CVEs](/manuals/scout/deep-dive/advisory-db-sources.md).

![Copilot vulnerabilities report](images/copilot-vuln-report.png?w=500px&border=1)

## Performing actions on your behalf

Before the agent performs any actions on your behalf, such as opening a pull
request for you, you're prompted to provide your consent to allow the
operation. You can always roll back or back out of the changes.

![Copilot action prompt](images/copilot-action-prompt.png?w=400px)

## Feedback

For issues or feedback, visit the [GitHub feedback repository](https://github.com/docker/copilot-issues).
