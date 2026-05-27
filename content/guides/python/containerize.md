---
title: Containerize a Python application
linkTitle: Containerize your app
weight: 10
keywords: python, flask, containerize, initialize
description: Learn how to containerize a Python application.
aliases:
  - /language/python/build-images/
  - /language/python/run-containers/
  - /language/python/containerize/
  - /guides/language/python/containerize/
---

## Prerequisites

- You have installed the latest version of [Docker Desktop](/get-started/get-docker.md).

## Overview

Containerizing your application means packaging it together with its
dependencies, configuration, and runtime into a single portable unit called a
container image. Running that image creates a container, an isolated process
that behaves the same on any machine, whether it's your laptop, a CI runner, or
a production server.

In this section, you'll containerize a simple
[FastAPI](https://fastapi.tiangolo.com) web application. You'll write a
`Dockerfile` that describes how to build the image, add a `compose.yaml` file
that defines how Docker runs your container, and then build and start the
application with one command.

You'll use [Docker Hardened Images](/dhi/) as the base. These are minimal,
secure Python images maintained by Docker.

## Create the application

The sample application is a minimal FastAPI service with a single endpoint
that returns a JSON greeting. Create the following files in a new
`python-docker-example` directory. To create all the files at once, switch to
the **Scaffold script** tab in the file browser and copy the shell command.

{{< files name="python-docker-example" >}}

{{< file path="app.py" status="new" >}}
```python
# A minimal FastAPI application.
# The root endpoint (GET /) returns a JSON "Hello World" response.
# See https://fastapi.tiangolo.com/ for the framework reference.

from fastapi import FastAPI

app = FastAPI()


@app.get("/")
async def root():
    return {"message": "Hello World"}
```
{{< /file >}}

{{< file path="requirements.txt" status="new" >}}
```text
# Python package dependencies for the application, pinned for reproducible builds.
# See https://pip.pypa.io/en/stable/reference/requirements-file-format/

fastapi==0.115.12
```
{{< /file >}}

{{< file path=".gitignore" status="new" >}}
```text
# Files and directories that Git should ignore. This is the standard Python
# template covering bytecode, build artifacts, virtual environments, and IDE
# settings. See https://git-scm.com/docs/gitignore for syntax reference.

# Byte-compiled / optimized / DLL files
__pycache__/
*.py[cod]
*$py.class

# C extensions
*.so

# Distribution / packaging
.Python
build/
develop-eggs/
dist/
downloads/
eggs/
.eggs/
lib/
lib64/
parts/
sdist/
var/
wheels/
share/python-wheels/
*.egg-info/
.installed.cfg
*.egg
MANIFEST

# Unit test / coverage reports
htmlcov/
.tox/
.nox/
.coverage
.coverage.*
.cache
nosetests.xml
coverage.xml
*.cover
*.py,cover
.hypothesis/
.pytest_cache/
cover/

# PEP 582; used by e.g. github.com/David-OConnor/pyflow and github.com/pdm-project/pdm
__pypackages__/

# Environments
.env
.venv
env/
venv/
ENV/
env.bak/
venv.bak/
```
{{< /file >}}

{{< /files >}}

If you already have Python installed and want to verify the app works before
containerizing it, you can run it locally:

```console
$ python3 -m venv .venv
$ source .venv/bin/activate
$ pip install -r requirements.txt
$ uvicorn app:app --reload
```

> [!NOTE]
>
> On Windows, activate the virtual environment with `.venv\Scripts\activate`
> instead of `source .venv/bin/activate`.

If you don't have Python installed, skip ahead to the next section. The
remaining steps run the application in a container, with no local Python
required.

## Create the Docker assets

Sign in to the DHI registry so Docker can pull the Python base images during
the build. The available Python images are listed in the
[catalog](https://hub.docker.com/hardened-images/catalog/dhi/python).

```console
$ docker login dhi.io
```

Add the following three files to your `python-docker-example` directory. The
`Dockerfile` describes how to build the image, `compose.yaml` defines how
Docker runs the container, and `.dockerignore` keeps unwanted files out of the
build context.

> [!TIP]
>
> [Gordon](/ai/gordon/), Docker's AI assistant, can generate Docker assets for
> your project. Ask Gordon to create a Dockerfile, Compose file, and
> `.dockerignore` tailored to your application.

{{< files name="python-docker-example" >}}

{{< file path="app.py" >}}
```python
# A minimal FastAPI application.
# The root endpoint (GET /) returns a JSON "Hello World" response.
# See https://fastapi.tiangolo.com/ for the framework reference.

from fastapi import FastAPI

app = FastAPI()


@app.get("/")
async def root():
    return {"message": "Hello World"}
```
{{< /file >}}

{{< file path="requirements.txt" >}}
```text
# Python package dependencies for the application, pinned for reproducible builds.
# See https://pip.pypa.io/en/stable/reference/requirements-file-format/

fastapi==0.115.12
```
{{< /file >}}

{{< file path="Dockerfile" status="new" >}}
```dockerfile
# syntax=docker/dockerfile:1

# Comments are provided throughout this file to help you get started.
# If you need more help, visit the Dockerfile reference guide at
# https://docs.docker.com/go/dockerfile-reference/

# This Dockerfile uses Docker Hardened Images (DHI) for enhanced security.
# For more information, see https://docs.docker.com/dhi/

# Use the dev image to build and install dependencies.
FROM dhi.io/python:3.12-dev AS builder

WORKDIR /app

RUN python3 -m venv /venv
ENV PATH="/venv/bin:$PATH"

# Download dependencies as a separate step to take advantage of Docker's caching.
# Leverage a cache mount to /root/.cache/pip to speed up subsequent builds.
# Leverage a bind mount to requirements.txt to avoid having to copy them into
# this layer.
RUN --mount=type=cache,target=/root/.cache/pip \
    --mount=type=bind,source=requirements.txt,target=requirements.txt \
    pip install -r requirements.txt

# Use the minimal runtime image. It runs as nonroot by default.
FROM dhi.io/python:3.12

WORKDIR /app

COPY --from=builder /venv /venv
ENV PATH="/venv/bin:$PATH"

# Copy the source code into the container.
COPY . .

# Expose the port that the application listens on.
EXPOSE 8000

# Run the application.
CMD ["/venv/bin/python3", "-m", "uvicorn", "app:app", "--host=0.0.0.0", "--port=8000"]
```
{{< /file >}}

{{< file path="compose.yaml" status="new" >}}
```yaml
# Comments are provided throughout this file to help you get started.
# If you need more help, visit the Docker Compose reference guide at
# https://docs.docker.com/go/compose-spec-reference/

# Here the instructions define your application as a service called "server".
# This service is built from the Dockerfile in the current directory.
# You can add other services your application may depend on here, such as a
# database or a cache. For examples, see the Awesome Compose repository:
# https://github.com/docker/awesome-compose
services:
  server:
    build:
      context: .
    ports:
      - 8000:8000
```
{{< /file >}}

{{< file path=".dockerignore" status="new" >}}
```text
# Include any files or directories that you don't want to be copied to your
# container here (e.g., local build artifacts, temporary files, etc.).
#
# For more help, visit the .dockerignore file reference guide at
# https://docs.docker.com/go/build-context-dockerignore/

**/.DS_Store
**/__pycache__
**/.venv
**/.classpath
**/.dockerignore
**/.env
**/.git
**/.gitignore
**/.project
**/.settings
**/.toolstarget
**/.vs
**/.vscode
**/*.*proj.user
**/*.dbmdl
**/*.jfm
**/bin
**/charts
**/docker-compose*
**/compose.y*ml
**/Dockerfile*
**/node_modules
**/npm-debug.log
**/obj
**/secrets.dev.yaml
**/values.dev.yaml
LICENSE
README.md
```
{{< /file >}}

{{< file path=".gitignore" >}}
```text
# Files and directories that Git should ignore. This is the standard Python
# template covering bytecode, build artifacts, virtual environments, and IDE
# settings. See https://git-scm.com/docs/gitignore for syntax reference.

# Byte-compiled / optimized / DLL files
__pycache__/
*.py[cod]
*$py.class

# C extensions
*.so

# Distribution / packaging
.Python
build/
develop-eggs/
dist/
downloads/
eggs/
.eggs/
lib/
lib64/
parts/
sdist/
var/
wheels/
share/python-wheels/
*.egg-info/
.installed.cfg
*.egg
MANIFEST

# Unit test / coverage reports
htmlcov/
.tox/
.nox/
.coverage
.coverage.*
.cache
nosetests.xml
coverage.xml
*.cover
*.py,cover
.hypothesis/
.pytest_cache/
cover/

# PEP 582; used by e.g. github.com/David-OConnor/pyflow and github.com/pdm-project/pdm
__pypackages__/

# Environments
.env
.venv
env/
venv/
ENV/
env.bak/
venv.bak/
```
{{< /file >}}

{{< /files >}}

To learn more about each file, see the following:

- [Dockerfile](/reference/dockerfile.md)
- [.dockerignore](/reference/dockerfile.md#dockerignore-file)
- [compose.yaml](/reference/compose-file/_index.md)

## Run the application

Inside the `python-docker-example` directory, run the following command in a
terminal.

```console
$ docker compose up --build
```

Open a browser and view the application at [http://localhost:8000](http://localhost:8000). You should see a simple FastAPI application.

In the terminal, press `ctrl`+`c` to stop the application.

### Run the application in the background

You can run the application detached from the terminal by adding the `-d`
option. Inside the `python-docker-example` directory, run the following command
in a terminal.

```console
$ docker compose up --build -d
```

Open a browser and view the application at [http://localhost:8000](http://localhost:8000).

To see the OpenAPI docs you can go to [http://localhost:8000/docs](http://localhost:8000/docs).

You should see a simple FastAPI application.

In the terminal, run the following command to stop the application.

```console
$ docker compose down
```

For more information about Compose commands, see the [Compose CLI
reference](/reference/cli/docker/compose/).

## Summary

In this section, you learned how you can containerize and run your Python
application using Docker.

Related information:

- [Docker Compose overview](/manuals/compose/_index.md)

## Next steps

In the next section, you'll take a look at how to set up a local development environment using Docker containers.
