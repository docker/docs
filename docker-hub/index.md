---
description: Docker Hub Quickstart
keywords: Docker, docker, registry, accounts, plans, Dockerfile, Docker Hub, docs, documentation, accounts, organizations, repositories, groups, teams
title: Docker Hub Quickstart
redirect_from:
- /docker-hub/overview/
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

[Docker Hub](https://hub.docker.com) is a service provided by Docker for
finding and sharing container images with your team. It provides the following
major features:
* [Repositories](repos.md): Push and pull container images.
* [Teams & Organizations](orgs.md): Manage access to private
repositories of container images.
* [Official Images](official_images.md): Pull and use high-quality
container images provided by Docker.
* [Publisher Images](publish/customer_faq.md): Pull and use high-
quality container images provided by external vendors.
* [Builds](builds/index.md): Automatically build container images from
GitHub and Bitbucket and push them to Docker Hub.
* [Webhooks](webhooks.md): Trigger actions after a successful push
  to a repository to integrate Docker Hub with other services.



### Step 1: Sign up for Docker Hub

Start by [creating an account](https://hub.docker.com/signup).

### Step 2: Create your first repository

To create a repo:

1. Sign in to [Docker Hub](https://hub.docker.com).

2. Click on **Create a Repository** on the Docker Hub welcome page:

    ![Welcome](images/index-welcome.png)

3. Name it **&lt;your-username&gt;/my-first-repo** as shown below. Select
   **Private**:

    ![Create Repository](images/index-create-repo.png)

    You've created your first repo. You should see:

    ![Repository Created](images/index-repo-created.png)

### Step 3: Download and install Docker Desktop

We'll need to download Docker Desktop to build and push a container image to
Docker Hub.

1. Download and install [Docker Desktop](https://docker.com/get-started). If on
Linux, download [Docker Engine - Community](https://hub.docker.com/search?type=edition&offering=community).

2. Open the terminal and sign in to Docker Hub on your computer by running
   `docker login`.

### Step 4: Build and push a container image to Docker Hub from your computer

1. Start by creating a [Dockerfile](https://docs.docker.com/engine/reference/builder/)
to specify your application as shown below:
```shell
cat > Dockerfile <<EOF
FROM busybox
CMD echo "Hello world! This is my first Docker image."
EOF
```

2. Run `docker build -t <your_username>/my-first-repo .` to build your Docker
   image.

3. Test your docker image locally by running `docker run <your_username>/my-first-repo`.

4. Run `docker push <your_username>/my-first-repo` to push your Docker image to
Docker Hub.

    You should see output similar to:

    ![Terminal](images/index-terminal-2019.png)

    And in Docker Hub, your repository should have a new `latest` tag available
    under **Tags**:

    ![Tag Created](images/index-tag.png)

Congratulations! You've successfully:
- Signed up for Docker Hub
- Created your first repository
- Built a Docker container image on your computer
- Pushed it to Docker Hub

### Next steps

- Create an [organization](orgs.md) to use Docker Hub with your team.
- Automatically build container images from code through [builds](builds/index.md).
- [Explore](https://hub.docker.com/explore) official & publisher images.
- [Upgrade your plan](upgrade.md) to push additional private Docker images to
Docker Hub.
