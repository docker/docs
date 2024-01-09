---
description: Learn how to run a personal Jupyter Server with the JupyterLab using Docker.
keywords: jupyter, notebook, python, jupyterlab, data science
title: Run JupyterLab
toc_min: 1
toc_max: 2
---

## Overview

This section walks you through how to run personal Jupyter Server with
JupyterLab in a container. Containers ensure that your JupyterLab environment
remains the same across different deployments, regardless of the operating
system or setup. With a single Docker command, you can consistently run
JuptyerLab on different machines.

## Prerequisites

You have installed the latest version of [Docker Desktop](../../../get-docker.md).

## Run a JupyterLab container

In a terminal, run the following command to run your JupyterLab container.

```console
$ docker run --rm -p 8888:8888 jupyter/base-notebook start-notebook.py --NotebookApp.token='my-token'
```

The following is a breakdown of the command:

- The `docker run` command tells Docker to pull and run an image. In this case,
the image is the `jupyter/base-notebook` image.
- The `-p` option maps port `8888` on the host to port `8888` in the container.
- `--rm` tells Docker to remove the image after it's stopped.
- `start-notebook.py --NotebookApp.token='my-token'` is the command that the
  server will run after starting. In this case, it's a script that sets the
  access token for convenience. For more details about `start-notebook.py`, see
  [Jupyter Server Options](https://jupyter-docker-stacks.readthedocs.io/en/latest/using/common.html#jupyter-server-options).

To access the container, in a web browser navigate to [localhost:8888/lab?token=my-token](http://localhost:8888/lab?token=my-token).

After accessing the URL, you should see the JupterLab interface. Since you specified `--rm` when running the container, the container and any data you create is deleted when the container stops.

To stop the container, press `ctrl`+`c` in the terminal. If you see a prompt confirming shutdown in the terminal, specify `Y` and press `Enter`.

## Save and access notebooks locally

When you remove a container, all data in that container is deleted. To save
notebooks or access notebooks on local machine, you can use a
[bind mount](../../../storage/bind-mounts.md).

### Run a JupterLab container with a bind mount

Open a terminal, and change directory to a directory where you would like to save the notebook's data. Then, run the following command based on your operating system and shell.

{{< tabs >}}
{{< tab name="Windows (Command Prompt)" >}}

```console
$ docker run --rm -p 8888:8888 -v "%cd%":/home/jovyan/work jupyter/base-notebook start-notebook.py --NotebookApp.token='my-token'
```
{{< /tab >}}
{{< tab name="Windows (PowerShell)" >}}

```console
$ docker run --rm -p 8888:8888 -v "$(pwd):/home/jovyan/work" jupyter/base-notebook start-notebook.py --NotebookApp.token='my-token'
```

{{< /tab >}}
{{< tab name="Windows (Git Bash)" >}}

```console
$ docker run --rm -p 8888:8888 -v "/$(pwd):/home/jovyan/work" jupyter/base-notebook start-notebook.py --NotebookApp.token='my-token'
```

{{< /tab >}}
{{< tab name="Mac / Linux" >}}

```console
$ docker run --rm -p 8888:8888 -v "$(pwd):/home/jovyan/work" jupyter/base-notebook start-notebook.py --NotebookApp.token='my-token'
```

{{< /tab >}}
{{< /tabs >}}

The `-v` option tells Docker to provide access to your machine's present working directory (PWD) from the container at `/home/jovyan/work`.

To access the container, in a web browser navigate to
[localhost:8888/lab?token=my-token](http://localhost:8888/lab?token=my-token).
Notebooks will now be saved to your local computer and will accessible even when
the container is deleted.

### Create and save a notebook

For this example, you'll use the [Iris Dataset](https://scikit-learn.org/stable/auto_examples/datasets/plot_iris_dataset.html) example from sickit-learn.

1. In the JupterLab interface, double click on the **work** folder in the left
   sidebar. This lets you save your work to the bind mount you created.

2. In the **Launcher**, under **Notebook**, select **Python 3**.

3. In the notebook, specify the following to install the necessary package.

   ```console
   !pip install matplotlib scikit-learn
   ```

4. Select the play button to run the code.

5. In the notebook, specify the following code.
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
6. Select the play button to run the code. You should see a scatter plot of the
   Iris dataset.

7. In the top menu, select **File** and then **Save Notebook As...**.

8. Specify a name in the `work` directory. For example, `work/mynotebook.ipynb`.

9. Select **Save**.

The notebook is now saved on your host machine in the directory that you bind
mounted.

In the terminal, press `ctrl`+ `c` to stop the container.

### Access notebooks on your local machine from a container

1. After saving a notebook, run the container with a bind mount again. Follow
   the instructions from
   [Run a JupterLab container with a bind mount](#run-a-jupyterlab-container).

2. In a web browser navigate to
   [localhost:8888/lab?token=my-token](http://localhost:8888/lab?token=my-token)
   to access Jupyter.

3. In the JupterLab interface, double click on the **work** folder in the left
   sidebar. You should see the notebook that you saved.

When you run the data plot code again, it'll need to run `!pip install
matplotlib scikit-learn` and download the packages. You can avoid reinstalling
packages every time you run a new container by creating your own image with the
packages already installed.

## Summary

In this topic, you learned how to run JupyterLab in a Docker container, which
provides a consistent and isolated environment for data analysis and
visualization. You also learned how to use a bind mount to persist notebook data
on your host machine, and how to create and save a notebook using the Iris
Dataset example from scikit-learn.

Related information:
- [docker run CLI reference](/reference/cli/docker/container/run/)
- [Bind mounts](../../../storage/bind-mounts.md)
- [Jupyter Docker Stacks documentation](https://jupyter-docker-stacks.readthedocs.io/en/latest/)

## Next steps

Next, you'll learn how you can avoid reinstalling dependencies every time you
run a new container. You'll do this by building your own image with the
dependencies already installed in the environment.

{{< button text="Build your own environment" url="customize.md" >}}