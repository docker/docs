---
title: Example prompts for the Docker agent
linkTitle: Example prompts
description: |
  Discover example prompts to interact with the Docker agent and learn how to
  automate tasks like Dockerizing projects or opening pull requests.
weight: 30
---

{{< summary-bar feature_name="Docker GitHub Copilot" >}}

## Use cases

Here are some examples of the types of questions you can ask the Docker agent:

### Ask general Docker questions

You can ask general question about Docker. For example:

- `@docker what is a Dockerfile?`
- `@docker how do I build a Docker image?`
- `@docker how do I run a Docker container?`
- `@docker what does 'docker buildx imagetools inspect' do?`

### Get help containerizing your project

You can ask the agent to help you containerize your existing project:

- `@docker can you help create a compose file for this project?`
- `@docker can you create a Dockerfile for this project?`

#### Opening pull requests

The Docker agent will analyze your project, generate the necessary files, and,
if applicable, offer to raise a pull request with the necessary Docker assets.

Automatically opening pull requests against your repositories is only available
when the agent generates new Docker assets.

### Analyze a project for vulnerabilities

The agent can help you improve your security posture with [Docker
Scout](/manuals/scout/_index.md):

- `@docker can you help me find vulnerabilities in my project?`
- `@docker does my project contain any insecure dependencies?`

The agent will run use Docker Scout to analyze your project's dependencies, and
report whether you're vulnerable to any [known CVEs](/manuals/scout/deep-dive/advisory-db-sources.md).

![Copilot vulnerabilities report](images/copilot-vuln-report.png?w=500px&border=1)

## Limitations

- The agent is currently not able to access specific files in your repository,
  such as the currently-opened file in your editor, or if you pass a file
  reference with your message in the chat message.

## Feedback

For issues or feedback, visit the [GitHub feedback repository](https://github.com/docker/copilot-issues).
