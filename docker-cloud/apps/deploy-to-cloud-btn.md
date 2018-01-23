---
description: Deploy to Docker Cloud
keywords: deploy, docker, cloud
redirect_from:
- /docker-cloud/feature-reference/deploy-to-cloud/
- /docker-cloud/tutorials/deploy-to-cloud/
title: Add a "Deploy to Docker Cloud" button
---

The **Deploy to Docker Cloud** button allows developers to deploy stacks with
one click in Docker Cloud as long as they are logged in. The button is intended
to be added to `README.md` files in public GitHub repositories, although it can
be used anywhere else.

> **Note**: You must be _logged in_ to Docker Cloud for the button to work
> Otherwise, the link results in a 404 error.

This is an example button to deploy our [python quickstart](https://github.com/docker/dockercloud-quickstart-python){: target="_blank" class="_"}:

<a href="https://cloud.docker.com/stack/deploy/?repo=https://github.com/docker/dockercloud-quickstart-python" target="_blank" class="_"><img src="https://files.cloud.docker.com/images/deploy-to-dockercloud.svg"></a>

The button redirects the user to the **Launch new Stack** wizard, with the stack
definition already filled with the contents of any of the following files (which
are fetched in the order shown) from the repository (taking into account branch
and relative path):

* `docker-cloud.yml`
* `docker-compose.yml`
* `fig.yml`

The user can still modify the stack definition before deployment.

## Add the 'Deploy to Docker Cloud' button in GitHub

You can simply add the following snippet to your `README.md` file:

```md
[![Deploy to Docker Cloud](https://files.cloud.docker.com/images/deploy-to-dockercloud.svg)](https://cloud.docker.com/stack/deploy/)
```

Docker Cloud detects the HTTP referer header and deploy the stack file found in the repository, branch and relative path where the source `README.md` file is stored.


## Add the 'Deploy to Docker Cloud' button in Docker Hub

If the button is displayed on the Docker Hub, Docker Cloud cannot automatically detect the source GitHub repository, branch and path. In this case, edit the repository description and add the following code:

```md
[![Deploy to Docker Cloud](https://files.cloud.docker.com/images/deploy-to-dockercloud.svg)](https://cloud.docker.com/stack/deploy/?repo=<repo_url>)
```

where `<repo_url>` is the path to your GitHub repository (see below).


## Add the 'Deploy to Docker Cloud' button anywhere else

If you want to use the button somewhere else, such as from external documentation or a landing site, you just need to create a link to the following URL:

```html
https://cloud.docker.com/stack/deploy/?repo=<repo_url>
```

where `<repo_url>` is the path to your GitHub repository. For example:

* `https://github.com/docker/dockercloud-quickstart-python`
* `https://github.com/docker/dockercloud-quickstart-python/tree/staging` to use branch `staging` instead of the default branch
* `https://github.com/docker/dockercloud-quickstart-python/tree/master/example` to use branch `master` and the relative path `/example` inside the repository

You can use your own image for the link (or no image). Our **Deploy to Docker Cloud** image is available at the following URL:

* `https://files.cloud.docker.com/images/deploy-to-dockercloud.svg`
