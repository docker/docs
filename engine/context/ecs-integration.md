---
title: Deploying Docker containers on ECS
description: Deploying Docker containers on ECS
keywords: Docker, AWS, ECS, Integration, context, Compose, cli, deploy, containers, cloud
toc_min: 1
toc_max: 2
---

## Overview

The Docker ECS Integration enables developers to use native Docker commands to run applications in Amazon EC2 Container Service (ECS) when building cloud-native applications.

The integration between Docker and Amazon ECS allow developers to use the Docker CLI to:

Set up an AWS context in one Docker command, allowing you to switch from a local context to a cloud context and run applications quickly and easily
Simplify multi-container application development on Amazon ECS using the Compose specification

>**Note**
>
> Docker ECS Integration is currently a beta release. The commands and flags are subject to change in subsequent releases.
{:.important}

## Prerequisites

To deploy Docker containers on ECS, you must meet the following requirements:

1. Download and install Docker Desktop Edge version 2.3.3.0 or later.

    - [Download for Mac](https://desktop.docker.com/mac/edge/Docker.dmg){: target="_blank" class="_"}
    - [Download for Windows](https://desktop.docker.com/win/edge/Docker%20Desktop%20Installer.exe){: target="_blank" class="_"}

    Alternatively, install [Docker ECS Integration for Linux](#install-the-docker-ecs-integration-cli-on-linux).

2. Ensure you have an AWS account.

  > **Note**
  >
  > If you had previously installed a Docker Desktop Stable release and now switched to Edge, ensure you turn on the experimental features flag. 
  >
  > From the Docker Desktop menu, click **Settings** (Preferences on macOS) > **Command Line** and then turn on the **Enable experimental features** toggle. Click **Apply & Restart** for the changes to take effect.

Check your installation by running the command `docker ecs version`.

Docker not only runs multi-container applications locally, but also enables developers to seamlessly deploy Docker containers on Amazon ECS using a Compose file with the `docker ecs compose up` command. The following sections contain instructions on how to deploy your Compose application on Amazon ECS.

### Create AWS context

Run the `docker ecs setup` command to create an AWS docker context. If you have already installed and configured the AWS CLI, the setup command lets you select an existing AWS profile to connect to Amazon. Otherwise, you can create a new profile by passing an [AWS access key ID and a secret access key](https://docs.aws.amazon.com/general/latest/gr/aws-sec-cred-types.html#access-keys-and-secret-access-keys){: target="_blank" class="_"}.

The `docker ecs setup` command will let you select an existing AWS configuration, or create one with provided secrets and tokens.

After you have created an AWS context, you can list your Docker contexts by running the `docker context ls` command:

```console
NAME   DESCRIPTION  DOCKER ENDPOINT  KUBERNETES ENDPOINT ORCHESTRATOR
aws *
default  Current DOCKER_HOST based configuration   unix:///var/run/docker.sock     swarm
```

## Run Compose applications

You can deploy and manage multi-container applications defined in Compose files to Amazon ECS using the `docker ecs compose` command. To do this:

- Ensure you are using your AWS context. You can do this either by specifying the `--context aws` flag with your command, or by setting the current context using the command `docker context use aws`.

- Run `docker ecs compose up` and `docker ecs compose down` to start and then stop a full Compose application.

  By default, `docker ecs compose up` uses the `docker-compose.yaml` file in the current folder. You can specify the Compose file directly using the `--file` flag.

  You can also specify a name for the Compose application using the `--project-name` flag during deployment. If no name is specified, a name will be derived from the working directory.

- You can view services created for the Compose application on Amazon ECS and their state using the `docker ecs compose ps` command.

- You can view logs from containers that are part of the Compose application using the `docker ecs compose logs` command.

## Private Docker images

The Docker ECS integration automatically configures authorization so you can pull private images from the Amazon ECR registry on the same AWS account. To pull private images from another registry, including Docker Hub, you’ll have to create a Username + Password (or a Username + Token) secret on the [Amazon SSM service](https://docs.aws.amazon.com/systems-manager/latest/userguide/systems-manager-parameter-store.html){: target="_blank" class="_"}.

For your convenience, Docker ECS integration offers the `docker ecs secret` command, so you can manage secrets created on AWS SMS without having to install the AWS CLI.

```console
docker ecs secret create dockerhubAccessToken --username <dockerhubuser>  --password <dockerhubtoken>
arn:aws:secretsmanager:eu-west-3:12345:secret:DockerHubAccessToken
```

Once created, you can use this ARN in you Compose file using using `x-aws-pull_credentials` custom extension with the Docker image URI for your service.

```console
version: 3.8
services:
  worker:
    image: mycompany/privateimage
    x-aws-pull_credentials: "arn:aws:secretsmanager:eu-west-3:12345:secret:DockerHubAccessToken"
```

>**Note**
>
> If you set the Compose file version to 3.8 or later, you can use the same Compose file for local deployment using `docker-compose`. Custom extensions will be ignored in this case.

## Service discovery

Service-to-service communication is implemented by the [Security Groups](https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/ec2-security-groups.html){: target="_blank" class="_"} rules, allowing services sharing a common Compose file “network” to communicate together. This allows individual services to run with distinct constraints (memory, cpu) and replication rules. However, it comes with a constraint that Docker images have to handle service discovery and wait for dependent services to be available.

### Service names

Services are registered by the Docker ECS integration on [AWS Cloud Map](https://docs.aws.amazon.com/cloud-map/latest/dg/what-is-cloud-map.html){: target="_blank" class="_"} during application deployment. They are declared as fully qualified domain names of the form: `<service>.<compose_project_name>.local`. Services can retrieve their dependencies using this fully qualified name, or can just use a short service name (as they do with docker-compose) as Docker ECS integration automatically injects the `LOCALDOMAIN` variable. This works out of the box if your Docker image fully implements domain name resolution standards, otherwise (typically, when using Alpine-based Docker images), you’ll have to include an [entrypoint script](https://docs.docker.com/develop/develop-images/dockerfile_best-practices/#entrypoint) in your Docker image to force this option:

```console
#! /bin/sh

if [ "${LOCALDOMAIN}x" != "x" ]; then echo "search ${LOCALDOMAIN}" >> /etc/resolv.conf; fi
exec "$@"
```

### Dependent service startup time and DNS resolution

Services get concurrently scheduled on ECS when a Compose file is deployed. AWS Cloud Map introduces an initial delay for DNS service to be able to resolve your services domain names. As a result, your code needs to be adjusted to support this delay by waiting for dependent services to be ready, or by adding a wait-script as the entrypoint to your Docker image, as documented in [Control startup order](https://docs.docker.com/compose/startup-order/).

Alternatively, you can use the [depends_on](https://github.com/compose-spec/compose-spec/blob/master/spec.md#depends_on){: target="_blank" class="_"} feature of the Compose file format. By doing this, dependent service will be created first, and application deployment will wait for it to be up and running before starting the creation of the dependent services.

## Tuning the CloudFormation template

The Docker ECS integration relies on [Amazon CloudFormation](https://docs.aws.amazon.com/cloudformation/){: target="_blank" class="_"} to manage the application deployment. To get more control on the created resources, you can use `docker ecs compose convert` to generate a CloudFormation stack file from your Compose file. This allows you to inspect resources it defines, or customize the template for your needs, and then apply the template to AWS using the AWS CLI, or the AWS web console.

By default, the Docker ECS integration creates an ECS cluster for your Compose application, a Security Group per network in your Compose file on your AWS account’s default VPC, and a LoadBalancer to route traffic to your services. If your AWS account does not have [permissions](https://github.com/docker/ecs-plugin/blob/master/docs/requirements.md#permissions){: target="_blank" class="_"} to create such resources, or you want to manage these yourself, you can use the following custom Compose extensions:

- Use `x-aws-cluster` as a top-level element in your Compose file to set the ARN
of an ECS cluster when deploying a Compose application. Otherwise, a 
cluster will be created for the Compose project.

- Use `x-aws-vpc` as a top-level element in your Compose file to set the ARN 
of a VPC when deploying a Compose application.

- Use `x-aws-loadbalancer` as a top-level element in your Compose file to set
the ARN of an existing LoadBalancer.

- Use `x-aws-securitygroup` inside a network definition in your Compose file to
set the ARN of an existing SecurityGroup used to implement network connectivity
between services.

## Install the Docker ECS Integration CLI on Linux

The Docker ECS Integration CLI adds support for running and managing containers on ECS.

>**Note**
>
> Docker ECS Integration is a beta release. The installation process, commands, and flags will change in future releases.
{:.important}

### Prerequisites

[Docker 19.03 or later](https://docs.docker.com/get-docker/)

### Download the plugin

You can download the Docker ECS plugin from the [docker/ecs-plugin](https://github.com/docker/ecs-plugin){: target="_blank" class="_"} GitHub repository using the following command:

```console
$ curl -LO https://github.com/docker/ecs-plugin/releases/latest/download/docker-ecs-linux-amd64
```

You will then need to make it an executable:

```console
$ chmod +x docker-ecs-linux-amd64
```

### Install the plugin

Move the plugin you’ve downloaded to the right place so the Docker CLI can use it:

```console
$ mkdir -p /usr/local/lib/docker/cli-plugins

$ mv docker-ecs-linux-amd64 /usr/local/lib/docker/cli-plugins/docker-ecs
```

You can move the CLI plugin into any of the following directories:

- `/usr/local/lib/docker/cli-plugins`
- `/usr/local/libexec/docker/cli-plugins`
- `/usr/lib/docker/cli-plugins`
- `/usr/libexec/docker/cli-plugins`

Finally, you must enable the experimental features on the CLI. You can do this by setting the environment variable `DOCKER_CLI_EXPERIMENTAL=enabled`, or by setting experimental to `enabled` in your Docker config file located at `~/.docker/config.json`:

```console
$ export DOCKER_CLI_EXPERIMENTAL=enabled

$ DOCKER_CLI_EXPERIMENTAL=enabled docker help

$ cat ~/.docker/config.json
{
  "experimental" : "enabled",
  "auths" : {
    "https://index.docker.io/v1/" : {

    }
  }
}
```

You can verify whether the CLI plugin installation is successful by checking whether it appears in the CLI help output, or by checking the plugin version. For example:

```console
$ docker help | grep ecs
  ecs*        Docker ECS (Docker Inc., 0.0.1)

$ docker ecs version
Docker ECS plugin 0.0.1
```

## FAQ

**What does the error `this tool requires the "new ARN resource ID format"` mean?**

This error message means that your integration requires the new ARN resource ID format for ECS. To learn more, see [Migrating your Amazon ECS deployment to the new ARN and resource ID format](https://aws.amazon.com/blogs/compute/migrating-your-amazon-ecs-deployment-to-the-new-arn-and-resource-id-format-2/){: target="_blank" class="_"}.

## Feedback

Thank you for trying out the Docker ECS Integration beta release. Your feedback is very important to us. Let us know your feedback by creating an issue in the [ecs-plugin](https://github.com/docker/ecs-plugin){: target="_blank" class="_"} GitHub repository.
