---
description: Use Docker images stored in IBM Cloud Container Registry
keywords: ibm, ibm cloud, registry, iaas, tutorial
title: Use images stored in IBM Cloud Container Registry
---

# Use IBM Cloud Container Registry to securely store your Docker images
[IBM Cloud Container Registry](https://www.ibm.com/cloud/container-registry) works with Docker EE for IBM Cloud to provide a secure image registry to use when creating swarm services and containers.

## Install the CLI and set up a namespace
Follow the [Getting Started with IBM Cloud Container Registry](https://console.bluemix.net/docs/services/Registry/index.html) instructions to install the registry CLI and set up a namespace.

## Log in to Docker with private registry credentials

### Log in with IBM Cloud Container Registry
You can place the credentials of your IBM Cloud account into Docker by running the registry login command:

```bash
$ bx cr login
```

The `bx cr login` command periodically expires, so you might need to run the command again when working with Docker or use a token.

### Log in with IBM Cloud Container Registry tokens
To prevent repeatedly logging in with `bx cr login`, you can create a non-expiring registry token with read-write permissions to use with Docker. Each token that you create is unique to the registry region. Repeat the steps for each registry region that you want to use with Docker.

1. [Create a registry token](https://console.bluemix.net/docs/services/Registry/registry_tokens.html#registry_tokens_create).
2. Instead of using the `bx cr login` command, log in to Docker with the registry token. Target the region for which you are using the registry, such as `ng` for US South:
    ```bash
    $ docker login -u token -p my_registry_token registry.ng.bluemix.net
    ```

### Log in with other private registries
You can also log in to Docker with other private registries. View that registry's documentation for the appropriate authentication methods. To use Docker Trusted Registry, [configure external IBM Cloud Object Storage](dtr-ibm-cos.md).

### Log in before running certain Docker commands
Log in to Docker (whether by the registry log-in, token, or other method) before running the following Docker commands:

- `docker pull` to download an image from IBM Cloud Container Registry.
- `docker push` to upload an image to IBM Cloud Container Registry.
- `docker service create` to create a service that uses an image that is stored in IBM Cloud Container Registry.
- Any Docker command that has the `--with-registry-auth` parameter.

## Create a container using an IBM Cloud Container Registry image
You can create a container using a registry image. You might want to run the image locally to test it before [creating a swarm service](#create-a-swarm-service-using-an-ibm-cloud-container-registry-image) based on the image.

Before you begin:

- [Install the registry CLI and set up a namespace](#install-the-cli-and-set-up-a-namespace).
- [Add an image in your registry namespace](https://console.bluemix.net/docs/services/Registry/registry_images_.html#registry_images_) to use to create the swarm service.
- [Log in to Docker](#log-in-to-docker-with-private-registry-credentials) with the appropriate registry credentials.

To create a local container that uses an IBM Cloud Container Registry image:

1. Get the name and tag of the image you want to use to create the service:

    ```bash
    $ bx cr images
    ```

2. Run the image locally:

    ```bash
    $ docker run --name my_container my_image:tag
    ```

**Tip**: If you no longer need the container, use the `docker kill` [command](/engine/reference/commandline/kill/) to remove it.

## Create a swarm service using an IBM Cloud Container Registry image
You can create a service that schedules tasks to spawn containers that are based on an image in your IBM Cloud Container Registry.

Before you begin:

- Install the Docker for IBM Cloud CLI.
- [Install the registry CLI and set up a namespace](#install-the-cli-and-set-up-a-namespace).
- [Add an image in your registry namespace](https://console.bluemix.net/docs/services/Registry/registry_images_.html#registry_images_) to use to create the service.
- [Log in to Docker](#log-in-to-docker-with-private-registry-credentials) with the appropriate registry credentials.
- [Create a Docker swarm](/engine/swarm/swarm-mode/#create-a-swarm).

To create a Docker swarm service that uses an IBM Cloud Container Registry image:

1. Get the name and tag of the image you want to use to create the service:

    ```bash
    $ bx cr images
    ```

2. Connect to your Docker for IBM Cloud swarm. Navigate to the directory where you [downloaded the UCP credentials](administering-swarms.md#download-client-certificates) and run the script. For example:

   ```bash
   $ cd filepath/to/certificate/repo && source env.sh
   ```

3. Send the registry authentication details when creating the Docker service for the image `<your-image:tag>`:

    ```bash
    $ docker service create --name my_service --with-registry-auth my_image:tag
    ```

4. Verify that your service was created:

    ```bash
    $ docker service ls --filter name=my_service
    ```

> More about creating Docker services
>
> You can also use the `docker service create` [command options](/engine/reference/commandline/service_create/) to set additional features such as replicas, global mode, or secrets. Visit the [Docker swarm services](/engine/swarm/services/) documentation to learn more.
