---
title: Create a repository
description: Learn how to create new repositories in Docker Trusted Registry.
keywords: registry, repository
redirect_from:
  - /datacenter/dtr/2.3/guides/user/manage-images/
---

Since DTR is secure by default, you need to create the image repository before
being able to push the image to DTR.

In this example, we'll create the 'golang' repository in DTR.

## Create a repository

To create a new repository, navigate to the **DTR web application**, and click
the **New repository** button.

![](../../images/create-repository-1.png){: .with-border}

Add a **name and description** for the repository, and choose whether your
repository is public or private:

  * Public repositories are visible to all users, but can only be changed by
  users granted with permission to write them.
  * Private repositories can only be seen by users that have been granted
  permissions to that repository.

![](../../images/create-repository-2.png){: .with-border}

Click **Save** to create the repository.

When creating a repository in DTR, the full name of the repository becomes
`<dtr-domain-name>/<user-or-org>/<repository-name>`. In this example, the full
name of our repository will be `dtr.example.org/dave.lauper/golang`.

> Image name size for DTR
>
> When creating an image name for use with DTR ensure that the organization and repository name has less than 56 characters and that the entire image name which includes domain, organization and repository name does not exceed 255 characters.
>
> The 56 character `<user-or-org/repository-name>` limit in DTR is due to an underlying limitation in how the image name information is stored within DTR metadata in RethinkDB.  RethinkDB currently has a Primary Key length limit of 127 characters.
>
> When DTR stores the above data it appends a sha256sum comprised of 72 characters to the end of the value to ensure uniqueness within the database.  If the `<user-or-org/repository-name>` exceeds 56 characters it will then exceed the 127 character limit in RethinkDB (72+56=128).
{: .important}

## Where to go next

- [Pull and push images](pull-and-push-images.md)
