---
title: Gordon use cases and examples
linkTitle: Use cases
description: Example prompts for common Docker workflows
weight: 10
---

{{< summary-bar feature_name="Gordon" >}}

Gordon handles Docker workflows through natural conversation. In Docker
Desktop, Gordon is available from the sidebar for open-ended sessions and from
contextual entry points in views like Containers, Images, Builds, and Volumes.
Selecting Gordon from one of these views opens a conversation pre-loaded with
context about the item you're looking at. You can ask the same questions from
the CLI with `docker ai`.

## Debug a failing container

You're in the Containers view and a container has crashed or behaves
unexpectedly. Open Gordon from the container row to ask about that container's
state and configuration:

- "Why did this container exit?"
- "What environment variables are set in this container?"
- "How long did this container run?"
- "What security settings are applied to this container?"

From the CLI:

```console
$ docker ai "why is my postgres container crashing on startup?"
```

## Debug a failed build

You're in the Builds view looking at a build that failed or is slower than
expected. Open Gordon from the build to inspect the Dockerfile, build
arguments, and cache behavior:

- "Why did this build fail?"
- "How can I improve cache usage for this build?"
- "What Dockerfile instructions were used?"
- "What build arguments were used?"

From the CLI:

```console
$ docker ai "my build is failing at the pip install step, what's wrong?"
```

## Inspect an image

You're in the Images view and want to understand what's in an image before
running it, or you want to size up a base image:

- "How do I run this image in the CLI?"
- "What environment variables are configured?"
- "What entrypoint is configured?"
- "What's the base architecture of this image?"
- "Is there a lighter version of this image?"

From the CLI:

```console
$ docker ai "compare my python:3.12 image to python:3.12-slim"
```

## Manage volumes and resources

From the Volumes view, ask Gordon about what's stored, which containers use a
volume, or how to clean up. From any view, use the Gordon sidebar to inspect
your wider environment:

- "Which containers are using this volume?"
- "Show me all my containers and their status"
- "How much disk space is Docker using?"
- "List my images sorted by size"

From the CLI:

```console
$ docker ai "clean up all unused Docker resources"
```

## Build and containerize

For new projects, start a conversation in the Gordon sidebar or via `docker
ai` from your project directory. Gordon reads your working directory and
proposes the right files:

- "Containerize my Node.js app"
- "Create a docker-compose for my stack"
- "Set up a dev environment with Postgres and Redis"

From the CLI:

```console
$ cd ~/my-project
$ docker ai "create a Dockerfile for this application"
```

## Develop and optimize

Ask Gordon to review and improve existing Dockerfiles or service definitions.
You can start from the Images view (for an image you've already built) or from
the Gordon sidebar with your project context:

- "Optimize this Dockerfile"
- "Add a health check to my service"
- "Make my Dockerfile more secure"

From the CLI:

```console
$ docker ai "rate my Dockerfile and suggest improvements"
```

## Learn Docker

For conceptual questions, use the Gordon sidebar or CLI. Gordon explains
concepts grounded in your environment, not generic answers:

- "What is a Docker volume?"
- "Explain multi-stage builds"
- "How does networking work in Docker?"

From the CLI:

```console
$ docker ai "what's the difference between COPY and ADD in a Dockerfile?"
```

## Writing effective prompts

Be specific:

- Include relevant context: "my postgres container" not "the database"
- State your goal: "make my build faster" not "optimize"
- Include error messages when debugging

Gordon works best when you describe what you want to achieve rather than how
to do it. Gordon maintains context across a conversation, so you can follow up
with clarifications or ask related questions without repeating yourself.

### Working directory context

When using `docker ai` in the CLI, Gordon uses your current working directory
as the default context for file operations. Change to your project directory
before starting Gordon to ensure it has access to the right files:

```console
$ cd ~/my-project
$ docker ai "review my Dockerfile"
```

You can also override the working directory with the `-C` flag. See [Using
Gordon via CLI](./how-to/cli.md#working-directory) for details.
