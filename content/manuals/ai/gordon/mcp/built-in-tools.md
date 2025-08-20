---
title: Built-in tools in Gordon
description: Use and configure Gordon's built-in tools for Docker, Kubernetes, security, and development workflows
keywords: ai, mcp, gordon, docker, kubernetes, security, developer tools, toolbox, configuration, usage
aliases:
 - /desktop/features/gordon/mcp/built-in-tools/
---

Gordon includes an integrated toolbox that gives you access to system tools and
capabilities. These tools extend Gordon's functionality so you can interact with
the Docker Engine, Kubernetes, Docker Scout security scanning, and other
developer utilities. This article describes the available tools, how to
configure them, and usage patterns.

## Configure tools

Configure tools globally in the toolbox to make them available throughout
Gordon, including Docker Desktop and the CLI.

To configure tools:

1. In the **Ask Gordon** view in Docker Desktop, select the **Toolbox** button at the bottom left of the input area.

   ![Screenshot showing Gordon page with the toolbox button.](../images/gordon.png)

1. To enable or disable a tool, select it in the left menu and select the toggle.

   ![Screenshot showing Gordon's Toolbox.](../images/toolbox.png)

   For more information about Docker tools, see [Reference](#reference).

## Usage examples

This section shows common tasks you can perform with Gordon tools.

### Manage Docker containers

#### List and monitor containers

```console
# List all running containers
$ docker ai "Show me all running containers"

# List containers using specific resources
$ docker ai "List all containers using more than 1GB of memory"

# View logs from a specific container
$ docker ai "Show me logs from my running api-container from the last hour"
```

#### Manage container lifecycle

```console
# Run a new container
$ docker ai "Run a nginx container with port 80 exposed to localhost"

# Stop a specific container
$ docker ai "Stop my database container"

# Clean up unused containers
$ docker ai "Remove all stopped containers"
```

### Work with Docker images

```console
# List available images
$ docker ai "Show me all my local Docker images"

# Pull a specific image
$ docker ai "Pull the latest Ubuntu image"

# Build an image from a Dockerfile
$ docker ai "Build an image from my current directory and tag it as myapp:latest"

# Clean up unused images
$ docker ai "Remove all my unused images"
```

### Manage Docker volumes

```console
# List volumes
$ docker ai "List all my Docker volumes"

# Create a new volume
$ docker ai "Create a new volume called postgres-data"

# Back up data from a container to a volume
$ docker ai "Create a backup of my postgres container data to a new volume"
```

### Perform Kubernetes operations

```console
# Create a deployment
$ docker ai "Create an nginx deployment and make sure it's exposed locally"

# List resources
$ docker ai "Show me all deployments in the default namespace"

# Get logs
$ docker ai "Show me logs from the auth-service pod"
```

### Run security analysis

```console
# Scan for CVEs
$ docker ai "Scan my application for security vulnerabilities"

# Get security recommendations
$ docker ai "Give me recommendations for improving the security of my nodejs-app image"
```

### Use development workflows

```console
# Analyze and commit changes
$ docker ai "Look at my local changes, create multiple commits with sensible commit messages"

# Review branch status
$ docker ai "Show me the status of my current branch compared to main"
```

## Reference

This section lists the built-in tools in Gordon's toolbox.

### Docker tools

Interact with Docker containers, images, and volumes.

#### Container management

<!-- vale off -->

| Name          | Description                      |
|---------------|----------------------------------|
| `docker`      | Access the Docker CLI            |
| `list_builds` | List builds in the Docker daemon |
| `build_logs`  | Show build logs                  |

#### Volume management

| Tool           | Description                |
|----------------|---------------------------|
| `list_volumes` | List all Docker volumes   |
| `remove_volume`| Remove a Docker volume    |
| `create_volume`| Create a new Docker volume|

#### Image management

| Tool           | Description                    |
|----------------|-------------------------------|
| `list_images`  | List all Docker images         |
| `remove_images`| Remove Docker images           |
| `pull_image`   | Pull an image from a registry  |
| `push_image`   | Push an image to a registry    |
| `build_image`  | Build a Docker image           |
| `tag_image`    | Tag a Docker image             |
| `inspect`      | Inspect a Docker object        |

### Kubernetes tools

Interact with your Kubernetes cluster.

#### Pod management

| Tool           | Description                        |
|----------------|------------------------------------|
| `list_pods`    | List all pods in the cluster       |
| `get_pod_logs` | Get logs from a specific pod       |

#### Deployment management

| Tool               | Description                        |
|--------------------|------------------------------------|
| `list_deployments` | List all deployments               |
| `create_deployment`| Create a new deployment            |
| `expose_deployment`| Expose a deployment as a service   |
| `remove_deployment`| Remove a deployment                |

#### Service management

| Tool           | Description                |
|----------------|---------------------------|
| `list_services`| List all services         |
| `remove_service`| Remove a service         |

#### Cluster information

| Tool             | Description                  |
|------------------|-----------------------------|
| `list_namespaces`| List all namespaces         |
| `list_nodes`     | List all nodes in the cluster|

### Docker Scout tools

Security analysis powered by Docker Scout.

| Tool                           | Description                                                                                                             |
|--------------------------------|-------------------------------------------------------------------------------------------------------------------------|
| `search_for_cves`              | Analyze a Docker image, project directory, or other artifacts for vulnerabilities using Docker Scout CVEs.              |
| `get_security_recommendations` | Analyze a Docker image, project directory, or other artifacts for base image update recommendations using Docker Scout. |

### Developer tools

General-purpose development utilities.

| Tool              | Description                      |
|-------------------|----------------------------------|
| `fetch`           | Retrieve content from a URL      |
| `get_command_help`| Get help for CLI commands        |
| `run_command`     | Execute shell commands           |
| `filesystem`      | Perform filesystem operations    |
| `git`             | Execute git commands             |

### AI model tools

| Tool           | Description                        |
|----------------|------------------------------------|
| `list_models`  | List all available Docker models   |
| `pull_model`   | Download a Docker model            |
| `run_model`    | Query a model with a prompt        |
| `remove_model` | Remove a Docker model              |

### Docker MCP Catalog

If you have enabled the [MCP Toolkit feature](../../mcp-catalog-and-toolkit/_index.md),
all the tools you have enabled and configured are available for Gordon to use.
