<!--[metadata]>
+++
aliases = ["/docker-trusted-registry/userguide/"]
title = "Create a repository"
description = "Learn how to manage your repositories on Docker Trusted Registry."
keywords = ["docker, registry, management, repository"]
[menu.main]
parent="dtr_menu_repos_and_images"
weight=0
+++
<![end-metadata]-->

## Create a repository

Before you can push images to your Docker Trusted Registry, you need to
create a repository for them.

To create a new repository:

1. In your browser navigate to the **Docker Trusted Registry web application**.

2. Navigate to the **Repositories** page. <!-- TODO: add sreenshot -->

3. Click **New repository**. <!-- TODO: add sreenshot -->


4. Add a **name and description** for the repository.
<!-- TODO: add sreenshot -->

5. Choose whether your repository is public or private:

  * Private repositories are visible to all users, but can only be changed by
  users granted with permission to write them.
  * Private repositories can only be seen by users that have been granted
  permissions to that repository.

6. Click **Create** to create the repository.

Now you can push your images to this repository.
