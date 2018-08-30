---
description: Prepare the application
keywords: Python, prepare, application
redirect_from:
- /docker-cloud/getting-started/python/3_prepare_the_app/
- /docker-cloud/getting-started/golang/3_prepare_the_app/
title: Prepare the application
---

In this step, you prepare a simple application that can be deployed.

## Clone the sample app

Run the following command to clone the sample application. You can use
either the Python or the Go version of this application, but you don't need to
install Python or Go to follow the tutorial.

### Python quickstart

```bash
$ git clone https://github.com/docker/dockercloud-quickstart-python.git
$ cd dockercloud-quickstart-python
```

### Go quickstart

```bash
$ git clone https://github.com/docker/dockercloud-quickstart-go.git
$ cd dockercloud-quickstart-go
```

## Build the application

*Skip the following step if you don't have Docker Engine installed locally.*

Next, we build this application to create an image. Run the following command to build the app. This creates a Docker image and tags it with whatever follows the word `tag`. Tag the image either `quickstart-python` or `quickstart-go` depending on which quickstart you are using.

### Python quickstart

```bash
$ docker build --tag quickstart-python .
```

### Go quickstart

```bash
$ docker build --tag quickstart-go .
```

## What's next?

Next, we [Push the Docker image to Docker Cloud's Registry](4_push_to_cloud_registry.md).
