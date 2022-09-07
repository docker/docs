---
description: Configuring repositories on Docker Hub
keywords: Docker, docker, trusted, registry, accounts, plans, Dockerfile, Docker Hub, webhooks, docs, documentation, creating, deleting, consolidating
title: Configuring repositories
redirect_from:
---
## Private repositories

Private repositories let you keep container images private, either to your
own account or within an organization or team.

To create a private repository, select **Private** when creating a repository:

![Create Private Repo](/docker-hub/images/repo-create-private.png){: style="max-width: 60%"}

You can also make an existing repository private by going to its **Settings** tab:

![Convert Repo to Private](/docker-hub/images/repo-make-private.png){: style="max-width: 60%"}

You get one private repository for free with your Docker Hub user account (not
usable for organizations you're a member of). If you need more private
repositories for your user account, upgrade your Docker Hub plan from your
[Billing Information](https://hub.docker.com/billing/plan){: target="_blank" rel="noopener" class="_"} page.

Once the private repository is created, you can `push` and `pull` images to and
from it using Docker.

> **Note**: You need to be signed in and have access to work with a
> private repository.

> **Note**: Private repositories are not currently available to search through
> the top-level search or `docker search`.

You can designate collaborators and manage their access to a private
repository from that repository's **Settings** page. You can also toggle the
repository's status between public and private, if you have an available
repository slot open. Otherwise, you can upgrade your
[Docker Hub](https://hub.docker.com/account/billing-plans/){: target="_blank" rel="noopener" class="_"} plan.

### Permissions reference

Permissions are cumulative. For example, if you have Read & Write permissions,
you automatically have Read-only permissions:

- `Read-only` access allows users to view, search, and pull a private repository in the same way as they can a public repository.
- `Read & Write` access allows users to pull, push, and view a repository Docker
  Hub. In addition, it allows users to view, cancel, retry or trigger builds
- `Admin` access allows users to Pull, push, view, edit, and delete a
  repository; edit build settings; update the repository description modify the
  repositories "Description", "Collaborators" rights, "Public/Private"
  visibility, and "Delete".

> **Note**
>
> A User who has not yet verified their email address only has
> `Read-only` access to the repository, regardless of the rights their team
> membership has given them.


