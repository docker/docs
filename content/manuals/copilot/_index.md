---
title: Docker for GitHub Copilot
params:
  sidebar:
    badge:
      color: violet
      text: EA
weight: 100
description: |
  Learn how to streamline Docker-related tasks with the Docker for GitHub
  Copilot extension. This integration helps you generate Docker assets, analyze
  vulnerabilities, and automate containerization through GitHub Copilot Chat in
  various development environments.
keywords: Docker, GitHub Copilot, extension, Visual Studio Code, chat, ai, containerization
---

{{% restricted title="Early Access" %}}
The Docker for GitHub Copilot extension is an [early access](/release-lifecycle#early-access-ea) product.
{{% /restricted %}}

The [Docker for GitHub Copilot](https://github.com/marketplace/docker-for-github-copilot)
extension integrates Docker's capabilities with GitHub Copilot, providing
assistance with containerizing applications, generating Docker assets, and
analyzing project vulnerabilities. This extension helps you streamline
Docker-related tasks wherever GitHub Copilot Chat is available.

## Key features

Key features of the Docker for GitHub Copilot extension include:

- Ask questions and receive responses about containerization in any context
  where GitHub Copilot Chat is available, such as on GitHub.com and in Visual Studio Code.
- Automatically generate Dockerfiles, Docker Compose files, and `.dockerignore`
  files for a project.
- Open pull requests with generated Docker assets directly from the chat
  interface.
- Get summaries of project vulnerabilities from [Docker
  Scout](/manuals/scout/_index.md) and receive next steps via the CLI.

## Data Privacy

The Docker agent is trained exclusively on Docker's documentation and tools to
assist with containerization and related tasks. It does not have access to your
project's data outside the context of the questions you ask.

When using the Docker Extension for GitHub Copilot, GitHub Copilot may include
a reference to the currently open file in its request if authorized by the
user. The Docker agent can read the file to provide context-aware responses.

If the agent is requested to check for vulnerabilities or generate
Docker-related assets, it will clone the referenced repository into in-memory
storage to perform the necessary actions.

Source code, questions, and responses are not persistently stored. The Docker
agent processes data only to provide responses and discards it immediately
afterward.

## Supported languages

The Docker Extension for GitHub Copilot supports the following
programming languages:

- Go
- Java
- JavaScript
- Python
- Rust
- TypeScript
