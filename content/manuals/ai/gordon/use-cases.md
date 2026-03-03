---
title: Gordon use cases and examples
linkTitle: Use cases
description: Example prompts for common Docker workflows
weight: 10
---

{{< summary-bar feature_name="Gordon" >}}

Gordon handles Docker workflows through natural conversation. This page shows
example prompts for the most common use cases.

## Debug and troubleshoot

Fix broken containers, diagnose build failures, and resolve issues.

```console
# Diagnose container crashes
$ docker ai "why did my postgres container crash?"

# Debug build failures
$ docker ai "my build is failing at the pip install step, what's wrong?"

# Fix networking issues
$ docker ai "my web container can't reach my database container"

# Investigate performance problems
$ docker ai "my container is using too much memory, help me investigate"
```

## Build and containerize

Create Docker assets for applications and migrate to hardened images.

```console
# Create Dockerfile from scratch
$ docker ai "create a Dockerfile for my Node.js application"

# Generate compose file
$ docker ai "create a docker-compose.yml for my application stack"

# Migrate to Docker Hardened Images
$ docker ai "migrate my Dockerfile to use Docker Hardened Images"
```

## Execute operations

Run Docker commands to manage containers, images, and resources.

```console
# Start containers with configuration
$ docker ai "run a redis container with persistence"

# Build and tag images
$ docker ai "build my Dockerfile and tag it for production"

# Clean up resources
$ docker ai "clean up all unused Docker resources"
```

## Develop and optimize

Improve Dockerfiles and configure secure, efficient development environments.

```console
# Optimize existing Dockerfile
$ docker ai "rate my Dockerfile and suggest improvements"

# Add security improvements
$ docker ai "make my Dockerfile more secure"

# Configure development workflow
$ docker ai "set up my container for development with hot reload"
```

## Manage resources

Inspect containers, images, and resource usage.

```console
# Check container status
$ docker ai "show me all my containers and their status"

# Analyze disk usage
$ docker ai "how much disk space is Docker using?"

# Review image details
$ docker ai "list my images sorted by size"
```

## Learn Docker

Understand concepts and commands in the context of your projects.

```console
# Explain Docker concepts
$ docker ai "explain how Docker networking works"

# Understand commands
$ docker ai "what's the difference between COPY and ADD in Dockerfile?"

# Get troubleshooting guidance
$ docker ai "how do I debug a container that exits immediately?"
```


## Writing effective prompts

Be specific:
- Include relevant context: "my postgres container" not "the database"
- State your goal: "make my build faster" not "optimize"
- Include error messages when debugging

Gordon works best when you describe what you want to achieve rather than how to
do it.

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
