---
description: Deploy Apps on Docker EE for IBM Cloud
keywords: ibm cloud, ibm, iaas, deploy
title: Deploy your app on Docker EE for IBM Cloud
---

## Deploy apps
Before you begin:

* Ensure that you [completed the account prerequisites](/docker-for-ibm-cloud/index.md).
* [Create a swarm](administering-swarms.md#create-swarms).
* [Set the environment variables to your swarm](administering-swarms.md#download-client-certificates).
* Review [best practices for creating a Dockerfile](/engine/userguide/eng-image/dockerfile_best-practices/) for your app's image.
* Review [Docker Compose file reference](/compose/compose-file/) for creating a YAML to define services, networks, or volumes that your app uses.

Steps:

* To deploy services and applications, see the [UCP documentation](/datacenter/ucp/2.2/guides/user/services/deploy-a-service/).
* To deploy IBM Cloud services such as Watson, see [Binding IBM Cloud services to swarms](binding-services.md).

## Run apps
After deploying an app to the cluster, you can create containers and services by using Docker commands such as `docker run`.

You can run websites too. Ports exposed with `-p` are automatically exposed
through the platform load balancer. For example:

  ```bash
  $ docker service create nginx \
  --name nginx \
  -p 80:80
  ```

Learn more about [load balancing in Docker EE for IBM Cloud](load-balancer.md).

### Execute Docker commands in all swarm nodes

You might need to execute a Docker command in all the nodes across the swarm, such as when installing a volume plug-in. Use the `swarm-exec` tool:

  ```bash
  $ swarm-exec {Docker command}
  ```

The following example installs a test plug-in in all the nodes in the swarm:

  ```bash
  $ swarm-exec docker plugin install --grant-all-permissions mavenugo/test-docker-netplugin
  ```

The `swarm-exec` tool internally uses the Docker global-mode service that runs a task on
each node in the cluster. The task in turn executes the `docker`
command. The global-mode service also guarantees that when a new node is added
to the cluster or during upgrades, a new task is executed on that node and hence
the `docker` command is automatically executed.

### Distributed Application Bundles

To deploy complex multi-container apps, you can use [distributed application
bundles](/compose/bundles.md). You can either run `docker deploy` to deploy a
bundle on your machine over an SSH tunnel, or copy the bundle (for example using
`scp`) to a manager node, SSH into the manager, and then run `docker deploy` (if
you have multiple managers, you have to ensure that your session is on one that
has the bundle file).

> SSH into manager
>
> Remember, the port to access your manager is 56422. For example: `ssh -A docker@managerIP -p 56422`.

A good sample app to test application bundles is the [Docker voting
app](https://github.com/docker/example-voting-app).

By default, apps deployed with bundles do not have ports publicly exposed.
When you change port mappings for services, Docker automatically updates the
underlying platform load balancer:

  ```bash
  $ docker service update --publish-add 80:80 my_service
  ```

> Publishing services on ports
>
> Your cluster's service load balancer can expose up to 10 ports. [Learn more](load-balancer.md#service-load-balancer).

### Images in private repos

To create swarm services using images in private repos, first make sure that you're
authenticated and have access to the private repo, and then create the service with
the `--with-registry-auth` flag. The following example assumes you're using Docker
Hub:

  ```bash
  $ docker login
  ...
  $ docker service create --with-registry-auth user/private-repo
  ...
  ```

The swarm caches and uses the cached registry credentials when creating containers for the service.

See [Using images with Docker for IBM Cloud](registry.md) for more information about using IBM Cloud Container Registry and Docker Trusted Registry.
