---
description: Build your own JupyterLab environment into a Docker image
keywords: jupyter, notebook, python, jupyterlab, data science
title: Customize your JupyterLab environment
toc_min: 1
toc_max: 2
---

## Overview

This section walks you through how to create your own JupyterLab environment and
build it into an image using Docker. By building your own image, you can
customize your JupyterLab environment with the packages and tools you need, and
ensure that it's consistent and reproducible across different deployments.
Building your own image also makes it easier to share your JupyterLab
environment with others, or to use it as a base for further development.

## Prerequisites

- You have installed the latest version of [Docker Desktop](../../../get-docker.md).
- Optionally, you have completed [Run JupyterLab](run-jupyter.md) for more
  context.

## Identify your environment

For this example, you'll use the Iris Dataset example from
[Run JupyterLab](run-jupyter.md). When you run a new container and then run that
notebook, the dependencies, `matplotlib` and `scikit-learn`, need to be
installed every time. While the dependencies in that small example download and
install quickly, it may become a problem as your list of dependencies grow.
There may also be other tools, packages, or files that you always want in your
environment.

In this case, you can install the dependencies as part of the environment in the
image. Then, every time you run your container, the dependencies will always be
installed.

Once you have identified your environment, you can start defining it in a
Dockerfile. A Dockerfile is a text file that instructs Docker how to create an
image of your JupyterLab environment. An image contains everything you want and
need when running JupyterLab, such as files, packages, and tools.

## Define your environment in a Dockerfile

In a directory of your choice, create a new text file named `Dockerfile`. Open the `Dockerfile` in an IDE or text editor and then add the following contents.

```dockerfile
# syntax=docker/dockerfile:1

FROM jupyter/base-notebook
RUN pip install --no-cache-dir matplotlib scikit-learn
```

Before you proceed, save your changes to the `Dockerfile`.

The following is a line-by-line breakdown of the `Dockerfile`.


```dockerfile
# syntax=docker/dockerfile:1
```

This optional instruction is a parser directive that tells Docker to use a
specific Dockerfile syntax for building the image. The syntax is defined by an
external image that contains the instructions and features for the Dockerfile.
Different versions of the Dockerfile syntax may have different features and
instructions that you can use to customize your image. By using a parser
directive, you can ensure that your Dockerfile is interpreted correctly by the
Docker Engine, and avoid potential errors or inconsistencies.

```dockerfile
FROM jupyter/base-notebook
```

This `FROM` instruction specifies the base image to use for building your own image. The base image is an existing image that has some of the files, packages, and tools you need for your JupyterLab environment. By using a base image, you can save time and avoid repeating common steps in your Dockerfile.

```dockerfile
RUN pip install --no-cache-dir matplotlib scikit-learn
```

The `RUN` instruction executes commands on top of the current image. In this
case, it's running `pip install --no-cache-dir matplotlib scikit-learn` to install dependencies. You can also use the `RUN` instruction to install JupyterLab extensions.

For more information about Dockerfile instructions, see the [Dockerfile reference](/reference/dockerfile/). For more recipes, see the [Docker Stacks documentation](https://jupyter-docker-stacks.readthedocs.io/en/latest/using/recipes.html).

## Build your environment into an image

After you have a `Dockerfile` to define your environment, you can use `docker
build` to build an image using your `Dockerfile`.

Open a terminal, change directory to the directory where your `Dockerfile` is
located, and then run the following command.

```console
$ docker build -t my-jupyter-image .
```

The command `docker build -t my-jupyter-image .` builds a Docker image from your
`Dockerfile` and a context. The `-t` option specifies the name and tag of the
image, in this case `my-jupyter-image`. The `.` indicates that the current
directory is the context, which means that the files in that directory can be
used in the image creation process.

You can verify that the image was built by viewing the `Images` view in Docker Desktop, or by running the `docker image ls` command in a terminal. You should see an image named `my-jupyter-image`.

## Run your image as a container

To run your image as a container, you use the same `docker run` command from the
[Run JupyterLab](run-jupyter.md) topic, but replace the image name with your
own. Make sure you stop any previously ran containers from this guide to prevent
port conflicts.

```console
$ docker run --rm -p 8888:8888 base-notebook start-notebook.py --NotebookApp.token='my-token'
```

To access the container, in a web browser navigate to
[localhost:8888/lab?token=my-token](http://localhost:8888/lab?token=my-token).

You can now use the packages without having to install them in your notebook.

1. In the **Launcher**, under **Notebook**, select **Python 3**.

2. In the notebook, specify the following code.

   ```python
   from sklearn import datasets

   iris = datasets.load_iris()
   import matplotlib.pyplot as plt

   _, ax = plt.subplots()
   scatter = ax.scatter(iris.data[:, 0], iris.data[:, 1], c=iris.target)
   ax.set(xlabel=iris.feature_names[0], ylabel=iris.feature_names[1])
   _ = ax.legend(
      scatter.legend_elements()[0], iris.target_names, loc="lower right", title="Classes"
   )
   ```

3. Select the play button to run the code. You should see a scatter plot of the Iris dataset.

In the terminal, press `ctrl`+ `c` to stop the container.

## Use Compose to run your container

Docker Compose is a tool for defining and running multi-container applications.
In this case, the application isn't a multi-container application, but Docker
Compose can make it easier to run by defining all the `docker run` options in a
file.

### Create a Compose file

To use Compose, you need a `compose.yaml` file. In the same directory as your
`Dockerfile`, create a new file named `compose.yaml`.

Open the `compose.yaml` file in an IDE or text editor and add the following
contents.

```yaml
services:
  jupyter:
    build:
      context: .
    ports:
      - 8888:8888
    volumes:
      - .:/home/jovyan/work
    command: start-notebook.py --NotebookApp.token='my-token'
```

The following is a line-by-line breakdown of the `compose.yaml` file.

```yaml
services:
```

This is the root key in the Docker Compose file, which you use to define the
services that make up your application.

```yaml
  jupyter:
```

This is the name of the service you're defining. You can name your services
whatever you like, but it's good practice to choose names that reflect the
purpose of the service.

```yaml
    build:
      context: .
```

`build` indicates that the service needs to be built from a Dockerfile and `context: .` specifies that the build context is the current directory (.).

```yaml
    ports:
      - 8888:8888
```
You use `ports` to map ports between the host and the container. This maps port `8888` on the host to port `8888` on the container.

```yaml
    volumes:
      - .:/home/jovyan/work
```

You use `volumes` to define mount points for volumes, which are directories that
persist data across container restarts and allow for data sharing between the
host and the container. In this case, `.:/home/jovyan/work` mounts the current
directory `.` (where the Compose file is located) to `/home/jovyan/work` in the
container.

```yaml
    command: start-notebook.py --NotebookApp.token='my-token'
```

You use `command` to run the `start-notebook.py --NotebookApp.token='my-token'`
command when the container starts. This is for your convenience. Otherwise, the
server starts with a random token that you must obtain from the logs.

For more details about the Compose instructions, see the
[Compose file reference](../../../compose/compose-file/_index.md).

### Run your container using Compose

Open a terminal, change directory to where your `compose.yaml` file is located, and then run the following command.

```console
$ docker compose up --build
```

This command builds your image and runs it as a container using the instructions
specified in the `compose.yaml` file. The `--build` option ensures that your
image is rebuilt, which is necessary if you made changes to your `Dockefile`.

To access the container, in a web browser navigate to
[localhost:8888/lab?token=my-token](http://localhost:8888/lab?token=my-token).

In the terminal, press `ctrl`+ `c` to stop the container.

## Summary

In this section, you learned how to construct a customized JupyterLab
environment within a Docker container. This involves creating a Dockerfile that
encapsulates your preferred packages and tools, allowing for a consistent and
reproducible development setup. You also explored how to use Docker Compose to
build and run containers with a simple command.

Related information:

- [Compose file reference](../../../compose/compose-file/_index.md)
- [Docker CLI reference](/reference/cli/docker/)
- [Dockerfile reference](/reference/dockerfile/)
- [Jupyter Docker Stacks documentation](https://jupyter-docker-stacks.readthedocs.io/en/latest/)

## Next steps

Next, you'll learn how you can add your notebooks to the image and share them.
Sharing the image with other data scientists facilitates
collaboration, reproducibility, and verification of results.

{{< button text="Share your JuptyerLab environment and notebook" url="share.md" >}}