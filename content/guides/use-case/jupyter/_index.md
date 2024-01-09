---
description: Run, develop, and share data science projects using JupyterLab and Docker
keywords: getting started, jupyter, notebook, python, jupyterlab, data science
title: Data science with JupyterLab
toc_min: 1
toc_max: 2
---

Docker and JupyterLab are two powerful tools that can enhance your data science
workflow. In this guide, you will learn how to use them together to create and
run reproducible data science environments. This guide is based on
[Supercharging AI/ML Development with JupyterLab and
Docker](https://www.docker.com/blog/supercharging-ai-ml-development-with-jupyterlab-and-docker/).

In this guide, you'll learn how to:

- Run a personal Jupyter Server with JupyterLab on your local machine
- Customize your JupyterLab environment
- Share your JupyterLab notebook and environment with other data scientists

## What is JupyterLab?

[JupyterLab](https://jupyterlab.readthedocs.io/en/stable/) is an open source application built around the concept of a computational notebook document. It enables sharing and executing code, data processing, visualization, and offers a range of interactive features for creating graphs.

## Why use Docker and JupyterLab together?

By combining Docker and JupyterLab, you can benefit from the advantages of both tools, such as:

- Containerization ensures a consistent JupyterLab environment across all
  deployments, eliminating compatibility issues.
- Containerized JupyterLab simplifies sharing and collaboration by removing the
  need for manual environment setup.
- Containers offer scalability for JupyterLab, supporting workload distribution
  and efficient resource management with platforms like Kubernetes.

Start by running a personal Jupyter Server with JupyterLab in a container.

{{< button text="Run JupyterLab in a container" url="run-jupyter.md" >}}