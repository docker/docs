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

    Alternatively, install the [Docker ECS Integration for Linux](#install-the-docker-ecs-integration-cli-on-linux).

2. Ensure you have an AWS account.

  > **Note**
  >
  > If you had previously installed a Docker Desktop Stable release and now switched to Edge, ensure you turn on the experimental features flag. 
  >
  > From the Docker Desktop menu, click **Settings** (Preferences on macOS) > **Command Line** and then turn on the **Enable experimental features** toggle. Click **Apply & Restart** for the changes to take effect.

Docker not only runs multi-container applications locally, but also enables 
developers to seamlessly deploy Docker containers on Amazon ECS using a 
Compose file with the `docker compose up` command. The following sections 
contain instructions on how to deploy your Compose application on Amazon ECS.

### Create AWS context

Run the `docker context create ecs myecscontext` command to create an Amazon ECS docker 
context named `myecscontext`. If you have already installed and configured the AWS CLI, 
the setup command lets you select an existing AWS profile to connect to Amazon. 
Otherwise, you can create a new profile by passing an 
[AWS access key ID and a secret access key](https://docs.aws.amazon.com/general/latest/gr/aws-sec-cred-types.html#access-keys-and-secret-access-keys){: target="_blank" class="_"}.

After you have created an AWS context, you can list your Docker contexts by running the `docker context ls` command:

```console
NAME   DESCRIPTION  DOCKER ENDPOINT  KUBERNETES ENDPOINT ORCHESTRATOR
myecscontext *
default  Current DOCKER_HOST based configuration   unix:///var/run/docker.sock     swarm
```

## Run Compose applications

You can deploy and manage multi-container applications defined in Compose files
to Amazon ECS using the `docker compose` command. To do this:

- Ensure you are using your ECS context. You can do this either by specifying 
the `--context myecscontext` flag with your command, or by setting the 
current context using the command `docker context use myecscontext`.

- Run `docker compose up` and `docker compose down` to start and then 
stop a full Compose application.

  By default, `docker compose up` uses the `docker-compose.yaml` file in 
  the current folder. You can specify the Compose file directly using the 
  `--file` flag.

  You can also specify a name for the Compose application using the `--project-name` flag during deployment. If no name is specified, a name will be derived from the working directory.

- You can view services created for the Compose application on Amazon ECS and 
their state using the `docker compose ps` command.

- You can view logs from containers that are part of the Compose application 
using the `docker compose logs` command.

## Private Docker images

The Docker ECS integration automatically configures authorization so you can pull private images from the Amazon ECR registry on the same AWS account. To pull private images from another registry, including Docker Hub, you’ll have to create a Username + Password (or a Username + Token) secret on the [Amazon SSM service](https://docs.aws.amazon.com/systems-manager/latest/userguide/systems-manager-parameter-store.html){: target="_blank" class="_"}.

For your convenience, Docker ECS integration offers the `docker secret` command, so you can manage secrets created on AWS SMS without having to install the AWS CLI.

```console
docker secret create dockerhubAccessToken --username <dockerhubuser>  --password <dockerhubtoken>
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

### Secrets

You can pass secrets to your ECS services using Docker model to bind sensitive 
data as files under `/run/secrets`. If your Compose file declares a secret as 
file, such a secret will be created as part of your application deployment on 
ECS. If you use an existing secret as `external: true` reference in your 
Compose file, use the ECS Secrets Manager full ARN as the secret name:
```yaml
services:
  webapp:
    image: ...
    secrets:
      - foo

secrets:
  foo:
    name: "arn:aws:secretsmanager:eu-west-3:1234:secret:foo-ABC123"
```

Secrets will be available at runtime for your service as a plain text file `/run/secrets/foo`.

The AWS Secrets Manager allows you to store sensitive data either as a plain 
text (like Docker secret does), or as a hierarchical JSON document. You can 
use the latter with ECS integration by using custom field `x-asw-keys` to 
define which entries in the JSON document to bind as a secret in your service 
container.

```yaml
services:
  webapp:
    image: ...
    secrets:
      - foo

secrets:
  foo:
    name: "arn:aws:secretsmanager:eu-west-3:1234:secret:foo-ABC123"
    keys: 
      - "bar"
```

By doing this, the secret for `bar` key will be available at runtime for your 
service as a plain text file `/run/secrets/foo/bar`. You can use the special 
value `*` to get all keys bound in your container. 

### Logging

The ECS integration configures AWS CloudWatch Logs service for your containers. 
A log group is created for the application as `docker-compose/<application_name>`, 
and log streams are created for each service and container in your application 
as `<application_name>/<service_name>/<container_ID>`.

You can fine tune AWS CloudWatch Logs using extension field `x-aws-logs_retention` 
in your Compose file to set the number of retention days for log events. The 
default behaviour is to keep logs forever.

You can also pass `awslogs` driver parameters to your container as standard 
Compose file `logging.driver_opts` elements.

### Dependent service startup time and DNS resolution

Services get concurrently scheduled on ECS when a Compose file is deployed. AWS Cloud Map introduces an initial delay for DNS service to be able to resolve your services domain names. As a result, your code needs to be adjusted to support this delay by waiting for dependent services to be ready, or by adding a wait-script as the entrypoint to your Docker image, as documented in [Control startup order](https://docs.docker.com/compose/startup-order/).

Alternatively, you can use the [depends_on](https://github.com/compose-spec/compose-spec/blob/master/spec.md#depends_on){: target="_blank" class="_"} feature of the Compose file format. By doing this, dependent service will be created first, and application deployment will wait for it to be up and running before starting the creation of the dependent services.

### Rolling update

Your ECS services are created with rolling update configuration. As you run 
`docker compose up` with a modified Compose file, the stack will be 
updated to reflect changes, and if required, some services will be replaced. 
This replacement process will follow the rolling-update configuration set by 
your services [`deploy.update_config`](https://docs.docker.com/compose/compose-file/#update_config) 
configuration. 

AWS ECS uses a percent-based model to define the number of containers to be 
run or shut down during a rolling update. The ECS integration computes 
rolling update configuration according to the `parallelism` and `replicas` 
fields. However, you might prefer to directly configure a rolling update 
using the extension fields `x-aws-min_percent` and `x-aws-max_percent`. 
The former sets the minimum percent of containers to run for service, and the 
latter sets the maximum percent of additional containers to start before 
previous versions are removed.

By default, the ECS rolling update is set to run twice the number of 
containers for a service (200%), and has the ability to shut down 100% 
containers during the update.


### IAM roles

Your ECS Tasks are executed with a dedicated IAM role, granting access 
to AWS Managed policies[`AmazonECSTaskExecutionRolePolicy`](https://docs.aws.amazon.com/AmazonECS/latest/developerguide/task_execution_IAM_role.html) 
and [`AmazonEC2ContainerRegistryReadOnly`](https://docs.aws.amazon.com/AmazonECS/latest/developerguide/ecr_managed_policies.html). 
In addition, if your service uses secrets, IAM Role gets additional 
permissions to read and decrypt secrets from the AWS Secret Manager.

You can grant additional managed policies to your service execution 
by using `x-aws-policies` inside a service definition:

```yaml
services:
  foo:
    x-aws-policies:
      - "arn:aws:iam::aws:policy/AmazonS3FullAccess"
```

You can also write your own [IAM Policy Document](https://docs.aws.amazon.com/AmazonECS/latest/developerguide/task-iam-roles.html) 
to fine tune the IAM role to be applied to your ECS service, and use 
`x-aws-role` inside a service definition to pass the 
yaml-formatted policy document.

```yaml
services:
  foo:
    x-aws-role:
      Version: "2012-10-17"
      Statement: 
        - Effect: "Allow"
          Action: 
            - "some_aws_service"
          Resource": 
            - "*"
```

## Tuning the CloudFormation template

The Docker ECS integration relies on [Amazon CloudFormation](https://docs.aws.amazon.com/cloudformation/){: target="_blank" class="_"} to manage the application deployment. To get more control on the created resources, you can use `docker compose convert` to generate a CloudFormation stack file from your Compose file. This allows you to inspect resources it defines, or customize the template for your needs, and then apply the template to AWS using the AWS CLI, or the AWS web console.

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

### Install script

You can install the new CLI using the install script:

```console
curl -L https://raw.githubusercontent.com/docker/aci-integration-beta/main/scripts/install_linux.sh | sh
```


## FAQ

**What does the error `this tool requires the "new ARN resource ID format"` mean?**

This error message means that your integration requires the new ARN resource ID format for ECS. To learn more, see [Migrating your Amazon ECS deployment to the new ARN and resource ID format](https://aws.amazon.com/blogs/compute/migrating-your-amazon-ecs-deployment-to-the-new-arn-and-resource-id-format-2/){: target="_blank" class="_"}.

## Feedback

Thank you for trying out the Docker ECS Integration beta release. Your feedback is very important to us. Let us know your feedback by creating an issue in the [ecs-plugin](https://github.com/docker/ecs-plugin){: target="_blank" class="_"} GitHub repository.
