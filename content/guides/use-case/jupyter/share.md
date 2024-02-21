---
description: Learn how to share your JupyterLab environment and notebooks in a Docker image.
keywords: jupyter, notebook, python, jupyterlab, data science
title: Share your JupyterLab environment and notebooks
toc_min: 1
toc_max: 2
---

## Overview

This section guides you on incorporating your Jupyter notebooks into a Docker
image and sharing it with the broader data science community through Docker Hub.

By adding your notebooks to an image, you create a portable and replicable
research environment that can be easily accessed and used by other data
scientists. This process not only facilitates collaboration but also ensures
that your work is preserved in an environment where it can be run without
compatibility issues.

You'll learn how to include your notebooks in a Docker image, how to publish
this image on Docker Hub, and how you or others can download and run your image.
This enables others to download, use, and build upon your work, fostering open
science and collaborative development.

## Prerequisites

- You have installed the latest version of [Docker Desktop](../../../get-docker.md).
- You have completed [Customize your JupyterLab environment](customize.md) and
  saved a notebook to your local machine.

## Add a notebook to a Docker image

This example builds upon the image and notebook created in [Customize your JupyterLab environment](customize.md). You'll now add that saved notebook from your local machine to the image.

Open the `Dockerfile` in and IDE or code editor and add a `COPY` instruction. The following is the contents of the `Dockerfile`.

```dockerfile {hl_lines=5}
# syntax=docker/dockerfile:1

FROM jupyter/base-notebook
RUN pip install --no-cache-dir matplotlib scikit-learn
COPY ./*.ipynb /home/jovyan/work/

```

Use the `COPY` command in a Dockerfile to copy files and directories from the
source file system (where the Dockerfile resides) into the filesystem of the
newly built Docker image. The specific command `COPY ./*.ipynb
/home/jovyan/work/` targets all IPython Notebook files (.ipynb) in the
Dockerfile's directory, copying them into the `/home/jovyan/work/` directory
inside the Docker container, making these notebooks available for use within the
JupyterLab environment.

Save the changes to your `Dockerfile`.

Build your image using the same command from [Customize your JupyterLab environment](customize.md).

```console
$ docker build -t my-jupyter-image .
```

Finally, run your image without a bind mount to verify that the image contains your notebook.

```console
$ docker run --rm -p 8888:8888 my-jupyter-image start-notebook.py --NotebookApp.token='my-token'
```

> **Note**
>
> Specifying a bind mount in the `docker run` command for the
> `/home/jovyan/work/` directory , or using Docker Compose with a bind mount
> will overwrite the files copied into the image. For this example, make sure
> you aren't using a bind mount.

To access the container, in a web browser navigate to
[localhost:8888/lab?token=my-token](http://localhost:8888/lab?token=my-token). You should find the notebook under the `work` directory.

## Publish your image to Docker Hub

Now that you have your notebook in an image, share it with the world. You'll need to create a repository on Docker Hub.

### Create a repository on Docker Hub

1. [Sign up](https://www.docker.com/pricing?utm_source=docker&utm_medium=webreferral&utm_campaign=docs_driven_upgrade) or sign in to [Docker Hub](https://hub.docker.com).

2. In Docker Hub, select **Repositories**, and then select the **Create
   repository** button.

3. For the repository name, use `my-jupyter-image`. Make sure the **Visibility**
   is **Public**.

4. Select **Create**.

Now that you have a repository, you can push your local image to the remote repository.

### Push your image to Docker Hub

1. Rename your image so that Docker knows which repository to push it to. Open a
   terminal and run the following `docker tag` command. Replace `YOUR-USER-NAME`
   with your Docker ID.

   ```console
   $ docker tag my-jupyter-image YOUR-USER-NAME/my-jupyter-image
   ```

2. Verify that you are signed in to Docker Desktop. In the
   Docker Dashboard , select **Sign in** if you aren't signed in yet.

3. Run the following `docker push` command to push the image to Docker Hub.
   Replace `YOUR-USER-NAME` with your Docker ID.

   ```console
   $ docker push YOUR-USER-NAME/my-jupyter-image
   ```

4. Verify that you pushed the image to Docker Hub.
   1. Go to [Docker Hub](https://hub.docker.com).
   2. Select **Repositories**.
   3. View the **Last pushed** time for your repository.

## Download and run your image

Others can now download and run your image. With Docker Desktop installed, they
can run the following command. Replace `YOUR-USER-NAME` with your Docker ID.

```console
$ docker run --rm -p 8888:8888 YOUR-USER-NAME/my-jupyer-image start-notebook.py --NotebookApp.token='my-token'
```

They can then access Jupyter and your notebook on their local machine at [localhost:8888/lab?token=my-token](http://localhost:8888/lab?token=my-token).

## Summary

In this section, you learned how to package your JupyterLab notebooks within a
Docker image and share them on Docker Hub, creating a portable and collaborative
research environment. You discovered how to include notebooks in a Docker image,
publish the image, and let others easily access and build upon your work.

Related information:

- [Docker Hub manual](../../../docker-hub/_index.md)
- [Dockerfile reference](/reference/dockerfile/)
- [Docker CLI reference](/reference/cli/docker/)