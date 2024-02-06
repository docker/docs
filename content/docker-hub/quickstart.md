---
description: Learn how to get started using Docker Hub
keywords: Docker, docker, registry, accounts, plans, Dockerfile, Docker Hub, accounts,
  organizations, repositories, groups, teams
title: Docker Hub quickstart
aliases:
- /apidocs/docker-cloud/
- /docker-cloud/
- /docker-cloud/migration/
- /docker-cloud/migration/cloud-to-swarm/
- /docker-cloud/migration/cloud-to-kube-aks/
- /docker-cloud/migration/cloud-to-kube-gke/
- /docker-cloud/migration/cloud-to-aws-ecs/
- /docker-cloud/migration/deregister-swarms/
- /docker-cloud/migration/kube-primer/
- /docker-cloud/cloud-swarm/
- /docker-cloud/cloud-swarm/using-swarm-mode/
- /docker-cloud/cloud-swarm/register-swarms/
- /docker-cloud/cloud-swarm/register-swarms/
- /docker-cloud/cloud-swarm/create-cloud-swarm-aws/
- /docker-cloud/cloud-swarm/create-cloud-swarm-azure/
- /docker-cloud/cloud-swarm/connect-to-swarm/
- /docker-cloud/cloud-swarm/link-aws-swarm/
- /docker-cloud/cloud-swarm/link-azure-swarm/
- /docker-cloud/cloud-swarm/ssh-key-setup/
- /docker-cloud/infrastructure/
- /docker-cloud/infrastructure/deployment-strategies/
- /docker-cloud/infrastructure/link-aws/
- /docker-cloud/infrastructure/link-do/
- /docker-cloud/infrastructure/link-azure/
- /docker-cloud/infrastructure/link-packet/
- /docker-cloud/infrastructure/link-softlayer/
- /docker-cloud/infrastructure/ssh-into-a-node/
- /docker-cloud/infrastructure/docker-upgrade/
- /docker-cloud/infrastructure/byoh/
- /docker-cloud/infrastructure/cloud-on-packet.net-faq/
- /docker-cloud/infrastructure/cloud-on-aws-faq/
- /docker-cloud/standard/
- /docker-cloud/getting-started/
- /docker-cloud/getting-started/intro_cloud/
- /docker-cloud/getting-started/connect-infra/
- /docker-cloud/getting-started/your_first_node/
- /docker-cloud/getting-started/your_first_service/
- /docker-cloud/getting-started/deploy-app/1_introduction/
- /docker-cloud/getting-started/deploy-app/2_set_up/
- /docker-cloud/getting-started/deploy-app/3_prepare_the_app/
- /docker-cloud/getting-started/deploy-app/4_push_to_cloud_registry/
- /docker-cloud/getting-started/deploy-app/5_deploy_the_app_as_a_service/
- /docker-cloud/getting-started/deploy-app/6_define_environment_variables/
- /docker-cloud/getting-started/deploy-app/7_scale_the_service/
- /docker-cloud/getting-started/deploy-app/8_view_logs/
- /docker-cloud/getting-started/deploy-app/9_load-balance_the_service/
- /docker-cloud/getting-started/deploy-app/10_provision_a_data_backend_for_your_service/
- /docker-cloud/getting-started/deploy-app/11_service_stacks/
- /docker-cloud/getting-started/deploy-app/12_data_management_with_volumes/
- /docker-cloud/apps/
- /docker-cloud/apps/deploy-to-cloud-btn/
- /docker-cloud/apps/auto-destroy/
- /docker-cloud/apps/autorestart/
- /docker-cloud/apps/auto-redeploy/
- /docker-cloud/apps/load-balance-hello-world/
- /docker-cloud/apps/deploy-tags/
- /docker-cloud/apps/stacks/
- /docker-cloud/apps/ports/
- /docker-cloud/apps/service-redeploy/
- /docker-cloud/apps/service-scaling/
- /docker-cloud/apps/api-roles/
- /docker-cloud/apps/service-links/
- /docker-cloud/apps/triggers/
- /docker-cloud/apps/volumes/
- /docker-cloud/apps/stack-yaml-reference/
- /docker-cloud/installing-cli/
- /docker-cloud/docker-errors-faq/
- /docker-cloud/release-notes/
- /docker-store/
- /apidocs/
- /apidocs/overview/
---

The following section contains step-by-step instructions on how to get started with Docker Hub.

### Step 1: Sign up for a free Docker account

Start by creating a [Docker ID](https://hub.docker.com/signup).

A [Docker ID](../docker-id/_index.md) grants you access to Docker Hub repositories and lets you explore available images from the community and verified publishers. You also need a Docker ID to share images on Docker Hub.

> **Tip**
>
> Explore [Docker's core subscriptions](https://www.docker.com/pricing/) to see what else Docker can offer you. 
{ .tip }

### Step 2: Create your first repository

To create a repository:

1. Sign in to [Docker Hub](https://hub.docker.com).
2. On the Repositories page, select **Create repository**.
3. Name it **&lt;your-username&gt;/my-private-repo**.
4. Set the visibility to **Private**.
5. Select **Create**.

You've created your first repository.

### Step 3: Download and install Docker Desktop

You need to download Docker Desktop to build, push, and pull container images.

1. Download and install [Docker Desktop](../desktop/index.md).

2. Sign in to Docker Desktop using the Docker ID you created in step one.

### Step 4: Pull and run a container image from Docker Hub

1. In your terminal, run `docker pull hello-world` to pull the image from Docker Hub. You should see output similar to:

   ```console
   $ docker pull hello-world
   Using default tag: latest
   latest: Pulling from library/hello-world
   2db29710123e: Pull complete
   Digest:   sha256:7d246653d0511db2a6b2e0436cfd0e52ac8c066000264b3ce63331ac66dca625
   Status: Downloaded newer image for hello-world:latest
   docker.io/library/hello-world:latest
   ```

2. Run `docker run hello-world` to run the image locally. You should see output similar to:

    ```console
    $ docker run hello-world
    Hello from Docker!
    This message shows that your installation appears to be working correctly.

    To generate this message, Docker took the following steps:
     1. The Docker client contacted the Docker daemon.
     2. The Docker daemon pulled the "hello-world" image from the Docker Hub.
     (amd64)
     3. The Docker daemon created a new container from that image which runs the
     executable that produces the output you are currently reading.
     4. The Docker daemon streamed that output to the Docker client, which sent
     it to your terminal.

    To try something more ambitious, you can run an Ubuntu container with:
     $ docker run -it ubuntu bash

    Share images, automate workflows, and more with a free Docker ID:
     https://hub.docker.com/

    For more examples and ideas, visit:
     https://docs.docker.com/get-started/
    ```

### Step 5: Build and push a container image to Docker Hub from your computer

1. Start by creating a [Dockerfile](../engine/reference/builder.md) to specify your application as shown below:

   ```dockerfile
   # syntax=docker/dockerfile:1
   FROM busybox
   CMD echo "Hello world! This is my first Docker image."
   ```

2. Run `docker build -t <your_username>/my-private-repo .` to build your Docker
   image.

3. Run `docker run <your_username>/my-private-repo` to test your
Docker image locally.

4. Run `docker push <your_username>/my-private-repo` to push your Docker image to Docker Hub. You should see output similar to:

   ```console
   $ cat > Dockerfile <<EOF
   FROM busybox
   CMD echo "Hello world! This is my first Docker image."
   EOF
   $ docker build -t mobythewhale/my-private-repo .
   [+] Building 1.2s (5/5) FINISHED
   => [internal] load build definition from Dockerfile
   => => transferring dockerfile: 110B
   => [internal] load .dockerignore
   => => transferring context: 2B
   => [internal] load metadata for docker.io/library/busybox:latest
   => CACHED [1/1] FROM docker.io/library/busybox@sha256:a9286defaba7n3a519
   => exporting to image
   => => exporting layers
   => => writing image sha256:dcdb1fd928bf257bfc0122ea47accd911a3a386ce618
   => => naming to docker.io/mobythewhale/my-private-repo
   $ docker run mobythewhale/my-private-repo
   Hello world! This is my first Docker image.
   $ docker push mobythewhale/my-private-repo
   The push refers to repository [docker.io/mobythewhale/my-private-repo]
   d2421964bad1: Layer already exists
   latest: digest: sha256:7604fbf8eeb03d866fd005fa95cdbb802274bf9fa51f7dafba6658294
   efa9baa size: 526
   ```

    >**Note**
    >
    > You must be signed in to Docker Hub through Docker Desktop or the command line, and you must also name your images correctly, as per the above steps.

5. Your repository in Docker Hub should now display a new `latest` tag under **Tags**:

    ![Tag created](images/index-tag.webp)

You've successfully:

- Signed up for a Docker account
- Created your first repository
- Pulled an existing container image from Docker Hub
- Built your own container image on your computer
- Pushed it successfully to Docker Hub

### Next steps

- Create an [organization](orgs.md) to use Docker Hub with your team.
- Automatically build container images from code through [builds](builds/index.md).
- [Explore](https://hub.docker.com/explore) official & publisher images.
- [Upgrade your subscription](https://www.docker.com/pricing) to push additional private Docker images to
Docker Hub.