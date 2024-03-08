---
title: Troubleshooting failing builds
keywords: concepts, build, images, docker desktop
description: Troubleshooting failing builds
---
<iframe width="650" height="365" src="https://www.youtube.com/embed/nsWWQ1xoEy0?rel=0" title="YouTube video player" frameborder="0" allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture; web-share" allowfullscreen></iframe>

## Explanation

In this concept, you will learn the following:
- how to troubleshooting common errors during the image building process

Building Docker images is a powerful way to containerize your applications, but sometimes things can go wrong, and builds can fail. This tutorial guides you through the essential steps to identify and resolve common issues encountered during Docker builds.

###  1. Error: "docker buildx build" requires exactly 1 argument

If you encounter the following error message while building the Docker image using `Dockerfile`, the easiest fix is adding `.` at the end of your `docker build` command. 

The `docker buildx build` command requires at least one argument, which is typically the path to the Dockerfile you want to use for the build

```console
 docker buildx build -t node-app .
```

This command assumes your `Dockerfile` is in the current directory. If it's located elsewhere, replace `.` with the actual path to your `Dockerfile`.

### 2. Error: Insufficient Scope, Authorization Failed


In case you encounter the following error message while configuring CI/CD to build Docker image

```console
 buildx failed with: ERROR: failed to solve: failed to push getting-started-todo-app:1.0: push access denied, repository does not exist or may require authorization: server message: insufficient_scope: authorization failed
```

This error message indicates that the CI/CD process was unable to push the built Docker image to the intended repository due to insufficient permissions. Here are the potential reasons:

- **Incorrect Credentials**: The DOCKER_USER and DOCKER_PAT secrets you provided in your GitHub repository might be incorrect, misspelled, or have insufficient access to the target repository.
- **Missing Repository**: The repository named `getting-started-todo-app` might not exist on the Docker Hub or the registry you're targeting.
- **Scope Limitations**: If you're using a personal access token (PAT) on Docker Hub, it might lack the necessary "repo:push" permission to push images to the repository.


#### How to fix

- **Verify Credentials**:

     - Double-check the DOCKER_USER and DOCKER_PAT values in your GitHub repository settings. Ensure they're accurate, including capitalization and special characters.
     - If you're using a personal access token, make sure it has the "repo:push" scope enabled for the specific repository you're targeting. You can create a new PAT with the appropriate scope on Docker Hub's "Your Account" page.

- **Confirm Repository Existence**:

     - If the repository doesn't exist, you'll need to create it on your Docker Hub account or the intended registry.
Check for Typos and Case Sensitivity:
     - Ensure there are no typos in the repository name (getting-started-todo-app) used in your workflow file or Dockerfile tags. Double-check for case sensitivity as well.

- **Review Workflow Configuration**:

     - Inspect your CI/CD workflow configuration (e.g., GitHub Actions) to verify that the repository name and credentials are set correctly.


### 3. Error: Docker no space left on device 

```console
  Docker: write /var/lib/docker/tmp/GetImageBlob104202287: no space left on device.
```

This error indicates that Docker has run out of storage space within the location it uses to store images, containers, and data. This often happens because Docker isn't effectively cleaning up unused or old resources.

#### How to fix

- **Check Disk Space**: Use commands like df -h to check the availability of space on your disk partitions. Look for the partition where Docker's data directory is typically located (/var/lib/docker) and see if there's enough free space.

- **Clean Up Unused Resources**: The most effective way to resolve this is to clean up unused Docker resources:

    - **Prune Images**: Use `docker image prune -a` to remove all unused and dangling images.
    - **Prune Containers**: Use `docker container prune` to remove stopped containers.
    - **Prune Volumes**: Use `docker volume prune` to remove unused volumes.
    - **Prune System**: Use  `docker system prune -a` to clean up images, containers, volumes, and networks. Exercise caution with this, as it's a more aggressive cleanup command.

- **Increase Disk Space**:

    - **Extend the Partition**: If your Docker partition is full, consider extending the partition where /var/lib/docker is located. This might involve resizing partitions on your hard disk.
    - **Change Docker Storage Location**: Move the Docker data directory to another partition with more space, using the -g option when starting the Docker daemon.

- **Prevention Strategies**:

    - **Regular Cleanup**:  Schedule regular executions of the `docker system prune` command to automate the removal of unused resources.
    - **Build Smaller Images**: Optimize your Dockerfiles to create smaller images, reducing their impact on disk space. Consider using multi-stage builds and compact base images.
    - **Monitor Disk Usage**: Regularly monitor disk space usage on the partition where Docker's data directory is located. Set up alerts to notify you when free space falls below a certain threshold.
