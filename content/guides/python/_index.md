---
title: Python language-specific guide
linkTitle: Python
description: Containerize Python apps using Docker
keywords: Docker, getting started, Python, language
summary: |
  This guide explains how to containerize Python applications using Docker.
aliases:
  - /language/python/
  - /guides/language/python/
  - /language/python/build-images/
  - /language/python/run-containers/
  - /language/python/containerize/
  - /language/python/develop/
  - /language/python/configure-ci-cd/
  - /guides/language/python/configure-ci-cd/
  - /language/python/deploy/
  - /guides/python/configure-github-actions/
  - /guides/python/containerize/
  - /guides/python/deploy/
  - /guides/python/develop/
  - /guides/python/lint-format-typing/
  - /guides/python/secure-supply-chain/
params:
  tags: [cicd]
  time: 20 minutes
---


> **Acknowledgment**
>
> This guide is a community contribution. Docker would like to thank
> [Esteban Maya](https://www.linkedin.com/in/esteban-x64/) and [Igor Aleksandrov](https://www.linkedin.com/in/igor-aleksandrov/) for their contribution
> to this guide.

The Python language-specific guide teaches you how to containerize a Python application using Docker. In this guide, you’ll learn how to:

- Containerize and run a Python application
- Set up a local environment to develop a Python application using containers
- Lint, format, typing and best practices
- Configure a CI/CD pipeline for a containerized Python application using GitHub Actions
- Deploy your containerized Python application locally to Kubernetes to test and debug your deployment

Start by containerizing an existing Python application.

## Containerize a Python application

### Prerequisites

- You have installed the latest version of [Docker Desktop](/get-started/get-docker.md).

### Overview

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

### Create the application

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
uvicorn==0.34.3
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

# Secrets
db/password.txt
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

### Create the Docker assets

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
uvicorn==0.34.3
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

# Secrets
db/password.txt
```
{{< /file >}}

{{< /files >}}

To learn more about each file, see the following:

- [Dockerfile](/reference/dockerfile.md)
- [.dockerignore](/reference/dockerfile.md#dockerignore-file)
- [compose.yaml](/reference/compose-file/_index.md)

### Run the application

Inside the `python-docker-example` directory, run the following command in a
terminal.

```console
$ docker compose up --build
```

Open a browser and view the application at [http://localhost:8000](http://localhost:8000). You should see a simple FastAPI application.

In the terminal, press `ctrl`+`c` to stop the application.

#### Run the application in the background

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

### Summary

In this section, you learned how you can containerize and run your Python
application using Docker.

Related information:

- [Docker Hardened Images](/dhi/)
- [Dockerfile reference](/reference/dockerfile.md)
- [Multi-stage builds](/manuals/build/building/multi-stage.md)
- [Docker Compose overview](/manuals/compose/_index.md)

### Next steps

In the next section, you'll take a look at how to set up a local development environment using Docker containers.

## Use containers for Python development

### Prerequisites

Complete [Containerize a Python application](containerize.md).

### Overview

Once your application runs in a container, the next step is making the
container loop part of your everyday development workflow. Code changes should
show up quickly, and services your app depends on, like databases, should run
right alongside it.

In this section, you'll extend the project from the previous topic by adding a
PostgreSQL database service to your `compose.yaml`, persisting the database
data in a named volume, and enabling Compose Watch so that changes you save in
your editor are picked up by the running container without a manual rebuild.

### Update the application

You'll update your application to connect to a PostgreSQL database. Continue
working in your `python-docker-example` directory.

Replace `app.py` and `requirements.txt`, and add a new `config.py` file with the
following contents.

> [!NOTE]
>
> The application won't run yet after this step. It tries to connect to a
> PostgreSQL database that doesn't exist. The next two sections add the
> database service and the Docker configuration needed to run everything
> together.

{{< files name="python-docker-example" >}}

{{< file path="app.py" status="modified" >}}
```python
# FastAPI application backed by a PostgreSQL database via SQLModel.
# The FastAPI lifespan handler creates database tables at startup.
# Endpoints: GET / (greeting), POST /heroes/ (create), GET /heroes/ (list).
# See https://fastapi.tiangolo.com/ and https://sqlmodel.tiangolo.com/

from collections.abc import AsyncGenerator, Sequence
from contextlib import asynccontextmanager

from fastapi import FastAPI
from sqlmodel import Field, Session, SQLModel, create_engine, select

from config import settings


class Hero(SQLModel, table=True):
    id: int | None = Field(default=None, primary_key=True)
    name: str = Field(index=True)
    secret_name: str
    age: int | None = Field(default=None, index=True)


engine = create_engine(str(settings.SQLALCHEMY_DATABASE_URI))


def create_db_and_tables() -> None:
    SQLModel.metadata.create_all(engine)


@asynccontextmanager
async def lifespan(_app: FastAPI) -> AsyncGenerator[None, None]:
    create_db_and_tables()
    yield


app = FastAPI(lifespan=lifespan)


@app.get("/")
def hello() -> str:
    return "Hello, Docker!"


@app.post("/heroes/")
def create_hero(hero: Hero) -> Hero:
    with Session(engine) as session:
        session.add(hero)
        session.commit()
        session.refresh(hero)
        return hero


@app.get("/heroes/")
def read_heroes() -> Sequence[Hero]:
    with Session(engine) as session:
        heroes = session.exec(select(Hero)).all()
        return heroes
```
{{< /file >}}

{{< file path="config.py" status="new" >}}
```python
# Pydantic settings that read PostgreSQL connection details from the
# environment. Supports a password file (Docker secrets) via
# POSTGRES_PASSWORD_FILE in addition to POSTGRES_PASSWORD.
# See https://docs.pydantic.dev/latest/concepts/pydantic_settings/

import os
from typing import Any

from pydantic import (
    PostgresDsn,
    computed_field,
    field_validator,
    model_validator,
)
from pydantic_core import MultiHostUrl
from pydantic_settings import BaseSettings


class Settings(BaseSettings):
    POSTGRES_SERVER: str
    POSTGRES_PORT: int = 5432
    POSTGRES_USER: str
    POSTGRES_PASSWORD: str | None = None
    POSTGRES_PASSWORD_FILE: str | None = None
    POSTGRES_DB: str

    @model_validator(mode="before")
    @classmethod
    def check_postgres_password(cls, data: Any) -> Any:
        """Validate that either POSTGRES_PASSWORD or POSTGRES_PASSWORD_FILE is set."""
        if isinstance(data, dict):
            password_file: str | None = data.get("POSTGRES_PASSWORD_FILE")  # type: ignore
            password: str | None = data.get("POSTGRES_PASSWORD")  # type: ignore
            if password_file is None and password is None:
                raise ValueError(
                    "At least one of POSTGRES_PASSWORD_FILE and POSTGRES_PASSWORD must be set."
                )
        return data  # type: ignore

    @field_validator("POSTGRES_PASSWORD_FILE", mode="before")
    @classmethod
    def read_password_from_file(cls, v: str | None) -> str | None:
        if v is not None:
            file_path = v
            if os.path.exists(file_path):
                with open(file_path) as file:
                    return file.read().strip()
            raise ValueError(f"Password file {file_path} does not exist.")
        return v

    @computed_field
    @property
    def SQLALCHEMY_DATABASE_URI(self) -> PostgresDsn:
        url = MultiHostUrl.build(
            scheme="postgresql+psycopg",
            username=self.POSTGRES_USER,
            password=self.POSTGRES_PASSWORD
            if self.POSTGRES_PASSWORD
            else self.POSTGRES_PASSWORD_FILE,
            host=self.POSTGRES_SERVER,
            port=self.POSTGRES_PORT,
            path=self.POSTGRES_DB,
        )
        return PostgresDsn(url)


settings = Settings()  # type: ignore
```
{{< /file >}}

{{< file path="requirements.txt" status="modified" hl_lines="5-7" >}}
```text
# Python package dependencies for the application, pinned for reproducible builds.
# See https://pip.pypa.io/en/stable/reference/requirements-file-format/

fastapi==0.115.12
sqlmodel==0.0.24
psycopg[binary]==3.2.9
pydantic-settings==2.9.1
uvicorn==0.34.3
```
{{< /file >}}

{{< file path="Dockerfile" >}}
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

{{< file path="compose.yaml" >}}
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

{{< file path=".dockerignore" >}}
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

# Secrets
db/password.txt
```
{{< /file >}}

{{< /files >}}

### Update Docker assets

Replace `Dockerfile` and `compose.yaml` with the following.

{{< files name="python-docker-example" >}}

{{< file path="app.py" >}}
```python
# FastAPI application backed by a PostgreSQL database via SQLModel.
# The FastAPI lifespan handler creates database tables at startup.
# Endpoints: GET / (greeting), POST /heroes/ (create), GET /heroes/ (list).
# See https://fastapi.tiangolo.com/ and https://sqlmodel.tiangolo.com/

from collections.abc import AsyncGenerator, Sequence
from contextlib import asynccontextmanager

from fastapi import FastAPI
from sqlmodel import Field, Session, SQLModel, create_engine, select

from config import settings


class Hero(SQLModel, table=True):
    id: int | None = Field(default=None, primary_key=True)
    name: str = Field(index=True)
    secret_name: str
    age: int | None = Field(default=None, index=True)


engine = create_engine(str(settings.SQLALCHEMY_DATABASE_URI))


def create_db_and_tables() -> None:
    SQLModel.metadata.create_all(engine)


@asynccontextmanager
async def lifespan(_app: FastAPI) -> AsyncGenerator[None, None]:
    create_db_and_tables()
    yield


app = FastAPI(lifespan=lifespan)


@app.get("/")
def hello() -> str:
    return "Hello, Docker!"


@app.post("/heroes/")
def create_hero(hero: Hero) -> Hero:
    with Session(engine) as session:
        session.add(hero)
        session.commit()
        session.refresh(hero)
        return hero


@app.get("/heroes/")
def read_heroes() -> Sequence[Hero]:
    with Session(engine) as session:
        heroes = session.exec(select(Hero)).all()
        return heroes
```
{{< /file >}}

{{< file path="config.py" >}}
```python
# Pydantic settings that read PostgreSQL connection details from the
# environment. Supports a password file (Docker secrets) via
# POSTGRES_PASSWORD_FILE in addition to POSTGRES_PASSWORD.
# See https://docs.pydantic.dev/latest/concepts/pydantic_settings/

import os
from typing import Any

from pydantic import (
    PostgresDsn,
    computed_field,
    field_validator,
    model_validator,
)
from pydantic_core import MultiHostUrl
from pydantic_settings import BaseSettings


class Settings(BaseSettings):
    POSTGRES_SERVER: str
    POSTGRES_PORT: int = 5432
    POSTGRES_USER: str
    POSTGRES_PASSWORD: str | None = None
    POSTGRES_PASSWORD_FILE: str | None = None
    POSTGRES_DB: str

    @model_validator(mode="before")
    @classmethod
    def check_postgres_password(cls, data: Any) -> Any:
        """Validate that either POSTGRES_PASSWORD or POSTGRES_PASSWORD_FILE is set."""
        if isinstance(data, dict):
            password_file: str | None = data.get("POSTGRES_PASSWORD_FILE")  # type: ignore
            password: str | None = data.get("POSTGRES_PASSWORD")  # type: ignore
            if password_file is None and password is None:
                raise ValueError(
                    "At least one of POSTGRES_PASSWORD_FILE and POSTGRES_PASSWORD must be set."
                )
        return data  # type: ignore

    @field_validator("POSTGRES_PASSWORD_FILE", mode="before")
    @classmethod
    def read_password_from_file(cls, v: str | None) -> str | None:
        if v is not None:
            file_path = v
            if os.path.exists(file_path):
                with open(file_path) as file:
                    return file.read().strip()
            raise ValueError(f"Password file {file_path} does not exist.")
        return v

    @computed_field
    @property
    def SQLALCHEMY_DATABASE_URI(self) -> PostgresDsn:
        url = MultiHostUrl.build(
            scheme="postgresql+psycopg",
            username=self.POSTGRES_USER,
            password=self.POSTGRES_PASSWORD
            if self.POSTGRES_PASSWORD
            else self.POSTGRES_PASSWORD_FILE,
            host=self.POSTGRES_SERVER,
            port=self.POSTGRES_PORT,
            path=self.POSTGRES_DB,
        )
        return PostgresDsn(url)


settings = Settings()  # type: ignore
```
{{< /file >}}

{{< file path="requirements.txt" >}}
```text
# Python package dependencies for the application, pinned for reproducible builds.
# See https://pip.pypa.io/en/stable/reference/requirements-file-format/

fastapi==0.115.12
sqlmodel==0.0.24
psycopg[binary]==3.2.9
pydantic-settings==2.9.1
uvicorn==0.34.3
```
{{< /file >}}

{{< file path="Dockerfile" status="modified" hl_lines="11,27-34,37,45" >}}
```dockerfile
# syntax=docker/dockerfile:1

# Comments are provided throughout this file to help you get started.
# If you need more help, visit the Dockerfile reference guide at
# https://docs.docker.com/go/dockerfile-reference/

# This Dockerfile uses Docker Hardened Images (DHI) for enhanced security.
# For more information, see https://docs.docker.com/dhi/

# Use the dev image to build and install dependencies.
# The builder stage is also used directly in development (see compose.yaml).
FROM dhi.io/python:3.12-dev AS builder

WORKDIR /app

RUN python3 -m venv /venv
ENV PATH="/venv/bin:$PATH"

# Download dependencies as a separate step to take advantage of Docker's caching.
# Leverage a cache mount to /root/.cache/pip to speed up subsequent builds.
# Leverage a bind mount to requirements.txt to avoid having to copy them
# into this layer.
RUN --mount=type=cache,target=/root/.cache/pip \
    --mount=type=bind,source=requirements.txt,target=requirements.txt \
    pip install -r requirements.txt

# Copy the source code into the container.
COPY . .

# Expose the port that the application listens on.
EXPOSE 8000

# Run the application.
CMD ["/venv/bin/python3", "-m", "uvicorn", "app:app", "--host=0.0.0.0", "--port=8000"]


# Use the minimal runtime image for production. It runs as nonroot by default.
FROM dhi.io/python:3.12

WORKDIR /app

COPY --from=builder /venv /venv
ENV PATH="/venv/bin:$PATH"

COPY --from=builder /app .

EXPOSE 8000

CMD ["/venv/bin/python3", "-m", "uvicorn", "app:app", "--host=0.0.0.0", "--port=8000"]
```
{{< /file >}}

{{< file path="compose.yaml" status="modified" hl_lines="8" >}}
```yaml
services:
  # Application service. The `target: builder` line builds the development
  # image (includes a shell and tools); the production stage of the
  # Dockerfile is unused in development.
  server:
    build:
      context: .
      target: builder
    ports:
      - 8000:8000
```
{{< /file >}}

{{< file path=".dockerignore" >}}
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

# Secrets
db/password.txt
```
{{< /file >}}

{{< /files >}}

#### About these changes

The `Dockerfile` builder stage now includes `COPY . .` and a `CMD`
instruction, which makes it directly runnable. This lets Compose target the
builder stage during development without rebuilding the production stage. The
production stage at the bottom is unchanged and still produces a minimal,
nonroot runtime image for shipping.

In `compose.yaml`, the new `target: builder` line tells Compose to build and
run the builder stage of the Dockerfile during development. Unlike the minimal
production image, the development image includes a shell and additional tools
that make debugging easier. If you need a shell in a running production
container, use [Docker Debug](/reference/cli/docker/debug/) instead.

### Add a local database and persist data

You can use containers to set up local services, like a database. In this
section, you'll update the `compose.yaml` file to define a database service
and a volume to persist data, and add a `db/password.txt` file that holds the
database password.

{{< files name="python-docker-example" >}}

{{< file path="app.py" >}}
```python
# FastAPI application backed by a PostgreSQL database via SQLModel.
# The FastAPI lifespan handler creates database tables at startup.
# Endpoints: GET / (greeting), POST /heroes/ (create), GET /heroes/ (list).
# See https://fastapi.tiangolo.com/ and https://sqlmodel.tiangolo.com/

from collections.abc import AsyncGenerator, Sequence
from contextlib import asynccontextmanager

from fastapi import FastAPI
from sqlmodel import Field, Session, SQLModel, create_engine, select

from config import settings


class Hero(SQLModel, table=True):
    id: int | None = Field(default=None, primary_key=True)
    name: str = Field(index=True)
    secret_name: str
    age: int | None = Field(default=None, index=True)


engine = create_engine(str(settings.SQLALCHEMY_DATABASE_URI))


def create_db_and_tables() -> None:
    SQLModel.metadata.create_all(engine)


@asynccontextmanager
async def lifespan(_app: FastAPI) -> AsyncGenerator[None, None]:
    create_db_and_tables()
    yield


app = FastAPI(lifespan=lifespan)


@app.get("/")
def hello() -> str:
    return "Hello, Docker!"


@app.post("/heroes/")
def create_hero(hero: Hero) -> Hero:
    with Session(engine) as session:
        session.add(hero)
        session.commit()
        session.refresh(hero)
        return hero


@app.get("/heroes/")
def read_heroes() -> Sequence[Hero]:
    with Session(engine) as session:
        heroes = session.exec(select(Hero)).all()
        return heroes
```
{{< /file >}}

{{< file path="config.py" >}}
```python
# Pydantic settings that read PostgreSQL connection details from the
# environment. Supports a password file (Docker secrets) via
# POSTGRES_PASSWORD_FILE in addition to POSTGRES_PASSWORD.
# See https://docs.pydantic.dev/latest/concepts/pydantic_settings/

import os
from typing import Any

from pydantic import (
    PostgresDsn,
    computed_field,
    field_validator,
    model_validator,
)
from pydantic_core import MultiHostUrl
from pydantic_settings import BaseSettings


class Settings(BaseSettings):
    POSTGRES_SERVER: str
    POSTGRES_PORT: int = 5432
    POSTGRES_USER: str
    POSTGRES_PASSWORD: str | None = None
    POSTGRES_PASSWORD_FILE: str | None = None
    POSTGRES_DB: str

    @model_validator(mode="before")
    @classmethod
    def check_postgres_password(cls, data: Any) -> Any:
        """Validate that either POSTGRES_PASSWORD or POSTGRES_PASSWORD_FILE is set."""
        if isinstance(data, dict):
            password_file: str | None = data.get("POSTGRES_PASSWORD_FILE")  # type: ignore
            password: str | None = data.get("POSTGRES_PASSWORD")  # type: ignore
            if password_file is None and password is None:
                raise ValueError(
                    "At least one of POSTGRES_PASSWORD_FILE and POSTGRES_PASSWORD must be set."
                )
        return data  # type: ignore

    @field_validator("POSTGRES_PASSWORD_FILE", mode="before")
    @classmethod
    def read_password_from_file(cls, v: str | None) -> str | None:
        if v is not None:
            file_path = v
            if os.path.exists(file_path):
                with open(file_path) as file:
                    return file.read().strip()
            raise ValueError(f"Password file {file_path} does not exist.")
        return v

    @computed_field
    @property
    def SQLALCHEMY_DATABASE_URI(self) -> PostgresDsn:
        url = MultiHostUrl.build(
            scheme="postgresql+psycopg",
            username=self.POSTGRES_USER,
            password=self.POSTGRES_PASSWORD
            if self.POSTGRES_PASSWORD
            else self.POSTGRES_PASSWORD_FILE,
            host=self.POSTGRES_SERVER,
            port=self.POSTGRES_PORT,
            path=self.POSTGRES_DB,
        )
        return PostgresDsn(url)


settings = Settings()  # type: ignore
```
{{< /file >}}

{{< file path="requirements.txt" >}}
```text
# Python package dependencies for the application, pinned for reproducible builds.
# See https://pip.pypa.io/en/stable/reference/requirements-file-format/

fastapi==0.115.12
sqlmodel==0.0.24
psycopg[binary]==3.2.9
pydantic-settings==2.9.1
uvicorn==0.34.3
```
{{< /file >}}

{{< file path="Dockerfile" >}}
```dockerfile
# syntax=docker/dockerfile:1

# Comments are provided throughout this file to help you get started.
# If you need more help, visit the Dockerfile reference guide at
# https://docs.docker.com/go/dockerfile-reference/

# This Dockerfile uses Docker Hardened Images (DHI) for enhanced security.
# For more information, see https://docs.docker.com/dhi/

# Use the dev image to build and install dependencies.
# The builder stage is also used directly in development (see compose.yaml).
FROM dhi.io/python:3.12-dev AS builder

WORKDIR /app

RUN python3 -m venv /venv
ENV PATH="/venv/bin:$PATH"

# Download dependencies as a separate step to take advantage of Docker's caching.
# Leverage a cache mount to /root/.cache/pip to speed up subsequent builds.
# Leverage a bind mount to requirements.txt to avoid having to copy them
# into this layer.
RUN --mount=type=cache,target=/root/.cache/pip \
    --mount=type=bind,source=requirements.txt,target=requirements.txt \
    pip install -r requirements.txt

# Copy the source code into the container.
COPY . .

# Expose the port that the application listens on.
EXPOSE 8000

# Run the application.
CMD ["/venv/bin/python3", "-m", "uvicorn", "app:app", "--host=0.0.0.0", "--port=8000"]


# Use the minimal runtime image for production. It runs as nonroot by default.
FROM dhi.io/python:3.12

WORKDIR /app

COPY --from=builder /venv /venv
ENV PATH="/venv/bin:$PATH"

COPY --from=builder /app .

EXPOSE 8000

CMD ["/venv/bin/python3", "-m", "uvicorn", "app:app", "--host=0.0.0.0", "--port=8000"]
```
{{< /file >}}

{{< file path="compose.yaml" status="modified" hl_lines="11-46" >}}
```yaml
services:
  # Application service. The `target: builder` line builds the development
  # image (includes a shell and tools); the production stage of the
  # Dockerfile is unused in development.
  server:
    build:
      context: .
      target: builder
    ports:
      - 8000:8000
    environment:
      - POSTGRES_SERVER=db
      - POSTGRES_USER=postgres
      - POSTGRES_DB=example
      - POSTGRES_PASSWORD_FILE=/run/secrets/db-password
    depends_on:
      db:
        condition: service_healthy
    secrets:
      - db-password
  # Database service. Reads the password from a Docker secret mounted at
  # /run/secrets/db-password. Compose waits for the healthcheck to pass
  # before starting the server, via the server's depends_on.
  db:
    image: dhi.io/postgres:18
    restart: always
    user: postgres
    secrets:
      - db-password
    volumes:
      - db-data:/var/lib/postgresql
    environment:
      - POSTGRES_DB=example
      - POSTGRES_PASSWORD_FILE=/run/secrets/db-password
    expose:
      - 5432
    healthcheck:
      test: ["CMD", "pg_isready"]
      interval: 10s
      timeout: 5s
      retries: 5
volumes:
  db-data:
secrets:
  db-password:
    file: db/password.txt
```
{{< /file >}}

{{< file path="db/password.txt" status="new" >}}
```text
mysecretpassword
```
{{< /file >}}

{{< file path=".dockerignore" >}}
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

# Secrets
db/password.txt
```
{{< /file >}}

{{< /files >}}

> [!NOTE]
>
> To learn more about the instructions in the Compose file, see [Compose file
> reference](/reference/compose-file/).

Now, run the following `docker compose up` command to start your application.

```console
$ docker compose up --build
```

Now test your API endpoint. Open a new terminal then make a request to the server using the curl commands:

Create an object with a POST request:

```console
$ curl -X 'POST' \
  'http://localhost:8000/heroes/' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '{
  "id": 1,
  "name": "my hero",
  "secret_name": "austing",
  "age": 12
}'
```

You should receive the following response:

```json
{
  "age": 12,
  "id": 1,
  "name": "my hero",
  "secret_name": "austing"
}
```

Now make a GET request:

```console
$ curl -X 'GET' \
  'http://localhost:8000/heroes/' \
  -H 'accept: application/json'
```

You should receive the same response as above because it's the only object in the database.

```json
{
  "age": 12,
  "id": 1,
  "name": "my hero",
  "secret_name": "austing"
}
```

Press `ctrl+c` in the terminal to stop your application.

### Automatically update services

Use Compose Watch to automatically update your running Compose services as you
edit and save your code. For more details about Compose Watch, see [Use Compose
Watch](/manuals/compose/how-tos/file-watch.md).

Open your `compose.yaml` file in an IDE or text editor and add the highlighted
Compose Watch instructions.

{{< files name="python-docker-example" >}}

{{< file path="app.py" >}}
```python
# FastAPI application backed by a PostgreSQL database via SQLModel.
# The FastAPI lifespan handler creates database tables at startup.
# Endpoints: GET / (greeting), POST /heroes/ (create), GET /heroes/ (list).
# See https://fastapi.tiangolo.com/ and https://sqlmodel.tiangolo.com/

from collections.abc import AsyncGenerator, Sequence
from contextlib import asynccontextmanager

from fastapi import FastAPI
from sqlmodel import Field, Session, SQLModel, create_engine, select

from config import settings


class Hero(SQLModel, table=True):
    id: int | None = Field(default=None, primary_key=True)
    name: str = Field(index=True)
    secret_name: str
    age: int | None = Field(default=None, index=True)


engine = create_engine(str(settings.SQLALCHEMY_DATABASE_URI))


def create_db_and_tables() -> None:
    SQLModel.metadata.create_all(engine)


@asynccontextmanager
async def lifespan(_app: FastAPI) -> AsyncGenerator[None, None]:
    create_db_and_tables()
    yield


app = FastAPI(lifespan=lifespan)


@app.get("/")
def hello() -> str:
    return "Hello, Docker!"


@app.post("/heroes/")
def create_hero(hero: Hero) -> Hero:
    with Session(engine) as session:
        session.add(hero)
        session.commit()
        session.refresh(hero)
        return hero


@app.get("/heroes/")
def read_heroes() -> Sequence[Hero]:
    with Session(engine) as session:
        heroes = session.exec(select(Hero)).all()
        return heroes
```
{{< /file >}}

{{< file path="config.py" >}}
```python
# Pydantic settings that read PostgreSQL connection details from the
# environment. Supports a password file (Docker secrets) via
# POSTGRES_PASSWORD_FILE in addition to POSTGRES_PASSWORD.
# See https://docs.pydantic.dev/latest/concepts/pydantic_settings/

import os
from typing import Any

from pydantic import (
    PostgresDsn,
    computed_field,
    field_validator,
    model_validator,
)
from pydantic_core import MultiHostUrl
from pydantic_settings import BaseSettings


class Settings(BaseSettings):
    POSTGRES_SERVER: str
    POSTGRES_PORT: int = 5432
    POSTGRES_USER: str
    POSTGRES_PASSWORD: str | None = None
    POSTGRES_PASSWORD_FILE: str | None = None
    POSTGRES_DB: str

    @model_validator(mode="before")
    @classmethod
    def check_postgres_password(cls, data: Any) -> Any:
        """Validate that either POSTGRES_PASSWORD or POSTGRES_PASSWORD_FILE is set."""
        if isinstance(data, dict):
            password_file: str | None = data.get("POSTGRES_PASSWORD_FILE")  # type: ignore
            password: str | None = data.get("POSTGRES_PASSWORD")  # type: ignore
            if password_file is None and password is None:
                raise ValueError(
                    "At least one of POSTGRES_PASSWORD_FILE and POSTGRES_PASSWORD must be set."
                )
        return data  # type: ignore

    @field_validator("POSTGRES_PASSWORD_FILE", mode="before")
    @classmethod
    def read_password_from_file(cls, v: str | None) -> str | None:
        if v is not None:
            file_path = v
            if os.path.exists(file_path):
                with open(file_path) as file:
                    return file.read().strip()
            raise ValueError(f"Password file {file_path} does not exist.")
        return v

    @computed_field
    @property
    def SQLALCHEMY_DATABASE_URI(self) -> PostgresDsn:
        url = MultiHostUrl.build(
            scheme="postgresql+psycopg",
            username=self.POSTGRES_USER,
            password=self.POSTGRES_PASSWORD
            if self.POSTGRES_PASSWORD
            else self.POSTGRES_PASSWORD_FILE,
            host=self.POSTGRES_SERVER,
            port=self.POSTGRES_PORT,
            path=self.POSTGRES_DB,
        )
        return PostgresDsn(url)


settings = Settings()  # type: ignore
```
{{< /file >}}

{{< file path="requirements.txt" >}}
```text
# Python package dependencies for the application, pinned for reproducible builds.
# See https://pip.pypa.io/en/stable/reference/requirements-file-format/

fastapi==0.115.12
sqlmodel==0.0.24
psycopg[binary]==3.2.9
pydantic-settings==2.9.1
uvicorn==0.34.3
```
{{< /file >}}

{{< file path="Dockerfile" >}}
```dockerfile
# syntax=docker/dockerfile:1

# Comments are provided throughout this file to help you get started.
# If you need more help, visit the Dockerfile reference guide at
# https://docs.docker.com/go/dockerfile-reference/

# This Dockerfile uses Docker Hardened Images (DHI) for enhanced security.
# For more information, see https://docs.docker.com/dhi/

# Use the dev image to build and install dependencies.
# The builder stage is also used directly in development (see compose.yaml).
FROM dhi.io/python:3.12-dev AS builder

WORKDIR /app

RUN python3 -m venv /venv
ENV PATH="/venv/bin:$PATH"

# Download dependencies as a separate step to take advantage of Docker's caching.
# Leverage a cache mount to /root/.cache/pip to speed up subsequent builds.
# Leverage a bind mount to requirements.txt to avoid having to copy them
# into this layer.
RUN --mount=type=cache,target=/root/.cache/pip \
    --mount=type=bind,source=requirements.txt,target=requirements.txt \
    pip install -r requirements.txt

# Copy the source code into the container.
COPY . .

# Expose the port that the application listens on.
EXPOSE 8000

# Run the application.
CMD ["/venv/bin/python3", "-m", "uvicorn", "app:app", "--host=0.0.0.0", "--port=8000"]


# Use the minimal runtime image for production. It runs as nonroot by default.
FROM dhi.io/python:3.12

WORKDIR /app

COPY --from=builder /venv /venv
ENV PATH="/venv/bin:$PATH"

COPY --from=builder /app .

EXPOSE 8000

CMD ["/venv/bin/python3", "-m", "uvicorn", "app:app", "--host=0.0.0.0", "--port=8000"]
```
{{< /file >}}

{{< file path="compose.yaml" status="modified" hl_lines="21-24" >}}
```yaml
services:
  # Application service. The `target: builder` line builds the development
  # image (includes a shell and tools); the production stage of the
  # Dockerfile is unused in development.
  server:
    build:
      context: .
      target: builder
    ports:
      - 8000:8000
    environment:
      - POSTGRES_SERVER=db
      - POSTGRES_USER=postgres
      - POSTGRES_DB=example
      - POSTGRES_PASSWORD_FILE=/run/secrets/db-password
    depends_on:
      db:
        condition: service_healthy
    secrets:
      - db-password
    develop:
      watch:
        - action: rebuild
          path: .
  db:
    image: dhi.io/postgres:18
    restart: always
    user: postgres
    secrets:
      - db-password
    volumes:
      - db-data:/var/lib/postgresql
    environment:
      - POSTGRES_DB=example
      - POSTGRES_PASSWORD_FILE=/run/secrets/db-password
    expose:
      - 5432
    healthcheck:
      test: ["CMD", "pg_isready"]
      interval: 10s
      timeout: 5s
      retries: 5
volumes:
  db-data:
secrets:
  db-password:
    file: db/password.txt
```
{{< /file >}}

{{< file path="db/password.txt" >}}
```text
mysecretpassword
```
{{< /file >}}

{{< file path=".dockerignore" >}}
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

# Secrets
db/password.txt
```
{{< /file >}}

{{< /files >}}

Run the following command to run your application with Compose Watch.

```console
$ docker compose watch
```

In a terminal, curl the application to get a response.

```console
$ curl http://localhost:8000
Hello, Docker!
```

Any changes to the application's source files on your local machine will now be immediately reflected in the running container.

Open `python-docker-example/app.py` in an IDE or text editor and update the `Hello, Docker!` string by adding a few more exclamation marks.

```diff
-    return 'Hello, Docker!'
+    return 'Hello, Docker!!!'
```

Save the changes to `app.py` and then wait a few seconds for the application to rebuild. Curl the application again and verify that the updated text appears.

```console
$ curl http://localhost:8000
Hello, Docker!!!
```

Press `ctrl+c` in the terminal to stop your application.

### Summary

In this section, you took a look at setting up your Compose file to add a local
database and persist data. You also learned how to use Compose Watch to automatically rebuild and run your container when you update your code.

Related information:

- [Compose file reference](/reference/compose-file/)
- [Compose secrets](/reference/compose-file/secrets.md)
- [Compose Watch](/manuals/compose/how-tos/file-watch.md)
- [Multi-stage builds](/manuals/build/building/multi-stage.md)

### Next steps

In the next section, you'll learn how you can set up linting, formatting, and type checking to follow the best practices in Python apps.

## Linting, formatting, and type checking for Python

### Prerequisites

Complete [Develop your app](develop.md). This topic requires a local Python
installation because the tools and Git hooks introduced here run on your
host. If you don't want to install Python locally, skip this topic. The same
checks run in CI in the [next topic](configure-github-actions.md).

### Overview

Linting, formatting, and type checking are automated ways to catch bugs,
enforce style, and spot type errors before code runs. Running them on every
commit, in CI, and in your editor catches problems early when they're cheap
to fix.

In this section, you'll configure three tools for your Python application.
Ruff handles linting and formatting in a single fast pass. Pyright statically
checks your code for type errors. Pre-commit hooks run both of these
automatically before each Git commit so problems are caught locally before
they're committed.

### Linting and formatting with Ruff

Ruff is an extremely fast Python linter and formatter written in Rust. It replaces multiple tools like flake8, isort, and black with a single unified tool.

Create a `pyproject.toml` file in your `python-docker-example` directory:

{{< files name="python-docker-example" >}}

{{< file path="pyproject.toml" status="new" >}}
```toml
# Configuration for code-quality tools.
# - [tool.ruff]: linting and formatting (https://docs.astral.sh/ruff/)
# - [tool.pyright]: static type checking (https://microsoft.github.io/pyright/)

[tool.ruff]
target-version = "py312"

[tool.ruff.lint]
select = [
    "E",  # pycodestyle errors
    "W",  # pycodestyle warnings
    "F",  # pyflakes
    "I",  # isort
    "B",  # flake8-bugbear
    "C4",  # flake8-comprehensions
    "UP",  # pyupgrade
    "ARG001", # unused arguments in functions
]
ignore = [
    "E501",  # line too long, handled by black
    "B008",  # do not perform function calls in argument defaults
    "W191",  # indentation contains tabs
    "B904",  # Allow raising exceptions without from e, for HTTPException
]
```
{{< /file >}}

{{< /files >}}

Install Ruff:

```console
$ pip install ruff
```

If you're using a virtual environment, make sure it is activated so the `ruff`
command is available.

Run these commands to check and format your code:

```console
# Check for errors
$ ruff check .

# Automatically fix fixable errors
$ ruff check --fix .

# Format code
$ ruff format .
```

### Type checking with Pyright

Pyright is a fast static type checker for Python that works well with modern Python features.

Update `pyproject.toml` to add the Pyright configuration at the bottom.

{{< files name="python-docker-example" >}}

{{< file path="pyproject.toml" status="modified" hl_lines="25-29" >}}
```toml
# Configuration for code-quality tools.
# - [tool.ruff]: linting and formatting (https://docs.astral.sh/ruff/)
# - [tool.pyright]: static type checking (https://microsoft.github.io/pyright/)

[tool.ruff]
target-version = "py312"

[tool.ruff.lint]
select = [
    "E",  # pycodestyle errors
    "W",  # pycodestyle warnings
    "F",  # pyflakes
    "I",  # isort
    "B",  # flake8-bugbear
    "C4",  # flake8-comprehensions
    "UP",  # pyupgrade
    "ARG001", # unused arguments in functions
]
ignore = [
    "E501",  # line too long, handled by black
    "B008",  # do not perform function calls in argument defaults
    "W191",  # indentation contains tabs
    "B904",  # Allow raising exceptions without from e, for HTTPException
]

[tool.pyright]
typeCheckingMode = "strict"
pythonVersion = "3.12"
exclude = [".venv"]
```
{{< /file >}}

{{< /files >}}

Install Pyright and run it:

```console
$ pip install pyright
$ pyright
```

### Setting up pre-commit hooks

Pre-commit hooks run checks automatically before each commit on your local
machine. Create a `.pre-commit-config.yaml` file in your `python-docker-example`
directory to set up Ruff hooks:

{{< files name="python-docker-example" >}}

{{< file path=".pre-commit-config.yaml" status="new" >}}
```yaml
# Pre-commit hook configuration. Runs Ruff (lint + format) on every
# `git commit`. See https://pre-commit.com/

repos:
  - repo: https://github.com/astral-sh/ruff-pre-commit
    rev: v0.15.15
    hooks:
      - id: ruff
        args: [--fix]
      - id: ruff-format
```
{{< /file >}}

{{< /files >}}

To install and use:

```console
$ pip install pre-commit
$ pre-commit install
$ git commit -m "Test commit"  # Automatically runs checks
```

### Summary

In this section, you learned how to:

- Configure and use Ruff for linting and formatting
- Set up Pyright for static type checking
- Automate checks with pre-commit hooks

These tools help maintain code quality and catch errors early in development.

Related information:

- [Ruff documentation](https://docs.astral.sh/ruff/)
- [Pyright documentation](https://microsoft.github.io/pyright/)
- [pre-commit framework](https://pre-commit.com/)

### Next steps

- [Configure GitHub Actions](configure-github-actions.md) to run these checks automatically
- Customize linting rules to match your team's style preferences
- Explore advanced type checking features

## Automate your builds with GitHub Actions

### Prerequisites

Complete all the previous sections of this guide, starting with [Containerize a Python application](containerize.md). You must have a [GitHub](https://github.com/signup) account and a verified [Docker](https://hub.docker.com/signup) account to complete this section.

If you didn't create a [GitHub repository](https://github.com/new) for your project yet, it is time to do it. After creating the repository, don't forget to [add a remote](https://docs.github.com/en/get-started/getting-started-with-git/managing-remote-repositories) and ensure you can commit and [push your code](https://docs.github.com/en/get-started/using-git/pushing-commits-to-a-remote-repository#about-git-push) to GitHub.

1. In your project's GitHub repository, open **Settings**, and go to **Secrets and variables** > **Actions**.

2. Under the **Variables** tab, create a new **Repository variable** named `DOCKER_USERNAME` and your Docker ID as a value.

3. Create a new [Personal Access Token (PAT)](/manuals/security/access-tokens.md#create-an-access-token) for Docker Hub. You can name this token `docker-tutorial`. Make sure access permissions include Read and Write.

4. Add the PAT as a **Repository secret** in your GitHub repository, with the name
   `DOCKERHUB_TOKEN`.

### Overview

GitHub Actions is a CI/CD automation tool built into GitHub. A workflow is a
YAML file that tells GitHub which jobs to run when something happens in your
repository, like a push to a branch or a pull request opening. Workflows live
in the `.github/workflows/` directory of your repository.

In this section, you'll add a workflow that runs your linting, formatting, and
type checks on every push to the main branch, then builds your Docker image
and pushes it to Docker Hub.

### 1. Define the GitHub Actions workflow

You can create a GitHub Actions workflow by creating a YAML file in the `.github/workflows/` directory of your repository. To do this use your favorite text editor or the GitHub web interface. The following steps show you how to create a workflow file using the GitHub web interface.

If you prefer to use the GitHub web interface, follow these steps:

1. Go to your repository on GitHub and then select the **Actions** tab.

2. Select **set up a workflow yourself**.

   This takes you to a page for creating a new GitHub Actions workflow file in
   your repository. By default, the file is created under `.github/workflows/main.yml`. Change the file name to `build.yml`.

If you prefer to use your text editor, create a new file named `build.yml` in the `.github/workflows/` directory of your repository.

Add the following content to the file:

{{< files name="python-docker-example" >}}

{{< file path=".github/workflows/build.yml" status="new" >}}
```yaml
# GitHub Actions workflow that runs on every push to main.
# - lint-test: runs pre-commit hooks (Ruff) and Pyright type checks.
# - build_and_push: signs in to Docker Hub and the DHI registry, then
#   builds and pushes the image (with SBOM and provenance attestations).
name: Build and push Docker image

on:
  push:
    branches:
      - main

jobs:
  lint-test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@{{% param "checkout_action_version" %}}

      - name: Set up Python
        uses: actions/setup-python@v6
        with:
          python-version: '3.12'

      - name: Install dependencies
        run: |
          python -m pip install --upgrade pip
          pip install -r requirements.txt
          pip install pre-commit pyright

      - name: Run pre-commit hooks
        run: pre-commit run --all-files

      - name: Run pyright
        run: pyright

  build_and_push:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@{{% param "checkout_action_version" %}}

      - name: Login to Docker Hub
        uses: docker/login-action@{{% param "login_action_version" %}}
        with:
          username: ${{ vars.DOCKER_USERNAME }}
          password: ${{ secrets.DOCKERHUB_TOKEN }}

      - name: Login to Docker Hardened Images
        uses: docker/login-action@{{% param "login_action_version" %}}
        with:
          registry: dhi.io
          username: ${{ vars.DOCKER_USERNAME }}
          password: ${{ secrets.DOCKERHUB_TOKEN }}

      - name: Set up Docker Buildx
        uses: docker/setup-buildx-action@{{% param "setup_buildx_action_version" %}}

      - name: Build and push
        uses: docker/build-push-action@{{% param "build_push_action_version" %}}
        with:
          push: true
          tags: ${{ vars.DOCKER_USERNAME }}/${{ github.event.repository.name }}:latest
```
{{< /file >}}

{{< /files >}}

Each GitHub Actions workflow includes one or several jobs. Each job consists of steps. Each step can either run a set of commands or use already [existing actions](https://github.com/marketplace?type=actions). The action above has three steps:

1. [Login to Docker Hub](https://github.com/docker/login-action): Action logs in to Docker Hub using the Docker ID and Personal Access Token (PAT) you created earlier.

2. [Set up Docker Buildx](https://github.com/docker/setup-buildx-action): Action sets up Docker [Buildx](https://github.com/docker/buildx), a CLI plugin that extends the capabilities of the Docker CLI.

3. [Build and push](https://github.com/docker/build-push-action): Action builds and pushes the Docker image to Docker Hub. The `tags` parameter specifies the image name and tag. The `latest` tag is used in this example.

### 2. Run the workflow

Commit the changes and push them to the `main` branch. This workflow is runs every time you push changes to the `main` branch. You can find more information about workflow triggers [in the GitHub documentation](https://docs.github.com/en/actions/writing-workflows/choosing-when-your-workflow-runs/events-that-trigger-workflows).

Go to the **Actions** tab of you GitHub repository. It displays the workflow. Selecting the workflow shows you the breakdown of all the steps.

When the workflow is complete, go to your [repositories on Docker Hub](https://hub.docker.com/repositories). If you see the new repository in that list, it means the GitHub Actions workflow successfully pushed the image to Docker Hub.

### Summary

In this section, you learned how to set up a GitHub Actions workflow for your Python application that includes:

- Running pre-commit hooks for linting and formatting
- Static type checking with Pyright
- Building and pushing Docker images

Related information:

- [Introduction to GitHub Actions](/guides/gha.md)
- [Docker Build GitHub Actions](/manuals/build/ci/github-actions/_index.md)
- [docker/login-action](https://github.com/docker/login-action)
- [docker/build-push-action](https://github.com/docker/build-push-action)
- [Create a Docker Hub access token](/manuals/security/access-tokens.md#create-an-access-token)

### Next steps

In the next section, you'll learn how to inspect and generate supply chain
attestations for your image. See [Secure your supply chain](secure-supply-chain.md).

## Secure your Python image supply chain

### Prerequisites

Complete [Configure CI/CD for your Python application](configure-github-actions.md).

### Overview

When you ship a container image, what's inside it and where it came from
matters. Supply chain attestations are signed records that answer questions
like which packages are in the image, what vulnerabilities affect them, how
the image was built, and what security checks it passed.

In this section, you'll inspect the attestations that ship with your Docker
Hardened Image base, generate your own SBOM and provenance attestations
during CI, and pin the base image by digest so your builds are reproducible.

The inspection commands in this topic are shown manually so you can see what
each one returns. In a real workflow you'd automate these checks with
[Docker Scout](/scout/), which runs the same scans on every push,
enforces policies in CI, and surfaces results in your registry and pull
requests.

### Inspect the base image attestations

Docker Hardened Images are built to SLSA Build Level 3 and ship with a set of
signed attestations covering bill-of-materials, vulnerabilities, build
provenance, and security scans. See
[DHI attestations](/manuals/dhi/core-concepts/attestations.md) for the full
list of types and how to verify their signatures with Cosign.

List all the attestations available on the Python DHI:

```console
$ docker scout attest list registry://dhi.io/python:3.12
```

View the SBOM:

```console
$ docker scout sbom registry://dhi.io/python:3.12
```

Check known vulnerabilities:

```console
$ docker scout cves registry://dhi.io/python:3.12
```

> [!NOTE]
>
> The `registry://` prefix forces `docker scout` to fetch the image and its
> attestations from the registry instead of reading a locally pulled copy. If
> you've already pulled or built against the base image, the local copy
> doesn't have the attached attestations, so the prefix is required to see
> them.

When you base your own image on a DHI image, these attestations stay attached to the base layer in the registry. Tools that inspect your image can follow the chain back to the DHI source.

### Generate attestations for your image

Update your GitHub Actions workflow to attach SBOM and provenance attestations to the image you push.

Edit `.github/workflows/build.yml` and update the build-and-push step:

```yaml {hl_lines="6-7"}
- name: Build and push Docker image
  uses: docker/build-push-action@v6
  with:
    context: .
    push: true
    sbom: true
    provenance: mode=max
    tags: ${{ steps.meta.outputs.tags }}
```

- `sbom: true` tells BuildKit to scan the built image and attach an SBOM attestation.
- `provenance: mode=max` records detailed build provenance, including the source repository, commit, and build parameters.

The next time your workflow runs, the pushed image will carry these attestations alongside the image manifest in the registry.

### Inspect your pushed image's attestations

After your workflow pushes the image, inspect it the same way you inspected the base image:

```console
$ docker scout attest list registry://DOCKER_USERNAME/REPO_NAME:latest
$ docker scout sbom registry://DOCKER_USERNAME/REPO_NAME:latest
```

The SBOM includes packages from every layer, including those inherited from `dhi.io/python:3.12`. The provenance record references the DHI base image by digest, so consumers of your image can trace the build chain back to the DHI source.

### Pin the base image by digest

Image tags like `dhi.io/python:3.12` move over time as new patches land. For reproducible builds, pin to an immutable digest.

The Dockerfile uses two tags, `dhi.io/python:3.12-dev` in the builder stage
and `dhi.io/python:3.12` in the runtime stage. Each tag has its own digest,
so look up both:

```console
$ docker buildx imagetools inspect dhi.io/python:3.12-dev --format "{{ .Manifest.Digest }}"
sha256:4f53cda18c2baa0c0354bb5f9a3ecbe5ed12ab4d8e11ba873c2f11161202b945
$ docker buildx imagetools inspect dhi.io/python:3.12 --format "{{ .Manifest.Digest }}"
sha256:e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
```

Each digest is a 64-character hex string. Update your `Dockerfile` to reference
each digest on the matching `FROM` line:

```dockerfile
FROM dhi.io/python:3.12-dev@sha256:4f53cda18c2baa0c0354bb5f9a3ecbe5ed12ab4d8e11ba873c2f11161202b945 AS builder
# ...
FROM dhi.io/python:3.12@sha256:e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
```

> [!TIP]
>
> Pinning by digest also pins you to that image's vulnerabilities. Use [Dependabot](https://docs.github.com/en/code-security/dependabot) or [Renovate](https://docs.renovatebot.com/) to automate digest updates so you get a PR when a new patched image is available, with a changelog to review before merging.

### Summary

In this section, you learned how to:

- Inspect the supply chain attestations that ship with the DHI base image, including SBOMs, CVE reports, VEX statements, and scan results
- Generate SBOM and provenance attestations for your own image in CI
- Pin base images by digest for reproducible builds

Related information:

- [DHI attestations](/manuals/dhi/core-concepts/attestations.md)
- [Verify a Docker Hardened Image](/manuals/dhi/how-to/verify.md)
- [Docker Scout](/scout/)
- [Build attestations](/manuals/build/metadata/attestations/_index.md)

### Next steps

In the next section, you'll deploy your application to Kubernetes.

## Test your Python deployment

### Prerequisites

- Complete all the previous sections of this guide, starting with [Use containers for Python development](develop.md).
- [Turn on Kubernetes](/manuals/desktop/use-desktop/kubernetes.md#enable-kubernetes) in Docker Desktop.

### Overview

[Kubernetes](https://kubernetes.io/) is an open source platform that runs and
orchestrates container workloads across one or more machines. You describe
what you want to run, like which container images, how many replicas, and
which network ports to expose, in YAML manifest files. Kubernetes reads the
manifests and makes the cluster match that description.

In this section, you'll use the Kubernetes environment built into Docker
Desktop to deploy your application locally. You'll write two manifest files,
one for the PostgreSQL database and one for the FastAPI application, apply
them with `kubectl`, and verify the deployment by hitting your application
from a terminal.

### Registry authentication

The Docker Hardened Images used in this guide are hosted on `dhi.io`. Docker
Desktop's Kubernetes shares credentials with Docker Desktop, so the `docker login dhi.io`
you completed earlier is all that's needed. No additional image pull secret is required.

> [!NOTE]
>
> If you're deploying to a Kubernetes cluster outside of Docker Desktop, you'll
> need to create an image pull secret and reference it in your pod specs. See
> [Use a Docker Hardened Image](/dhi/how-to/use/#use-with-kubernetes) for instructions.

### Create a Kubernetes YAML file

Create the following two Kubernetes manifest files in your
`python-docker-example` directory. Before applying
`docker-python-kubernetes.yaml`, replace `DOCKER_USERNAME/REPO_NAME` with your
Docker username and the repository name that you created in [Configure CI/CD for
your Python application](./configure-github-actions.md).

{{< files name="python-docker-example" >}}

{{< file path="docker-postgres-kubernetes.yaml" status="new" >}}
```yaml
# Kubernetes manifests for the PostgreSQL database used by the FastAPI app.
# Contains a Deployment, Service, PersistentVolumeClaim, and Secret.

# Deployment: runs one PostgreSQL pod. The image, port, env vars, and the
# persistent volume mount are all defined here.
apiVersion: apps/v1
kind: Deployment
metadata:
  name: postgres
  namespace: default
spec:
  replicas: 1
  selector:
    matchLabels:
      app: postgres
  template:
    metadata:
      labels:
        app: postgres
    spec:
      containers:
        - name: postgres
          image: dhi.io/postgres:18
          ports:
            - containerPort: 5432
          env:
            - name: POSTGRES_DB
              value: example
            - name: POSTGRES_USER
              value: postgres
            - name: POSTGRES_PASSWORD
              valueFrom:
                secretKeyRef:
                  name: postgres-secret
                  key: POSTGRES_PASSWORD
          volumeMounts:
            - name: postgres-data
              mountPath: /var/lib/postgresql
      volumes:
        - name: postgres-data
          persistentVolumeClaim:
            claimName: postgres-pvc
---
# Service: exposes PostgreSQL inside the cluster on port 5432 so the
# application pod can reach it by the DNS name `postgres`.
apiVersion: v1
kind: Service
metadata:
  name: postgres
  namespace: default
spec:
  ports:
    - port: 5432
  selector:
    app: postgres
---
# PersistentVolumeClaim: storage that survives pod restarts.
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: postgres-pvc
  namespace: default
spec:
  accessModes:
    - ReadWriteOnce
  resources:
    requests:
      storage: 1Gi
---
# Secret: holds the database password (base64-encoded). Referenced by both
# the postgres Deployment and the application Deployment.
apiVersion: v1
kind: Secret
metadata:
  name: postgres-secret
  namespace: default
type: Opaque
data:
  POSTGRES_PASSWORD: cG9zdGdyZXNfcGFzc3dvcmQ= # Base64 encoded password (e.g., 'postgres_password')
```
{{< /file >}}

{{< file path="docker-python-kubernetes.yaml" status="new" >}}
```yaml
# Kubernetes manifests for the FastAPI application.
# Contains a Deployment and a NodePort Service.

# Deployment: runs the FastAPI app. Connection details to the postgres
# service are passed in via environment variables, and the database
# password comes from the shared postgres-secret.
apiVersion: apps/v1
kind: Deployment
metadata:
  name: docker-python-demo
  namespace: default
spec:
  replicas: 1
  selector:
    matchLabels:
      service: fastapi
  template:
    metadata:
      labels:
        service: fastapi
    spec:
      containers:
        - name: fastapi-service
          image: DOCKER_USERNAME/REPO_NAME
          imagePullPolicy: Always
          env:
            - name: POSTGRES_PASSWORD
              valueFrom:
                secretKeyRef:
                  name: postgres-secret
                  key: POSTGRES_PASSWORD
            - name: POSTGRES_USER
              value: postgres
            - name: POSTGRES_DB
              value: example
            - name: POSTGRES_SERVER
              value: postgres
            - name: POSTGRES_PORT
              value: "5432"
          ports:
            - containerPort: 8000
---
# Service: exposes the FastAPI app on port 30001 of the cluster node so
# you can reach it from your host with `curl http://localhost:30001/`.
apiVersion: v1
kind: Service
metadata:
  name: service-entrypoint
  namespace: default
spec:
  type: NodePort
  selector:
    service: fastapi
  ports:
    - port: 8000
      targetPort: 8000
      nodePort: 30001
```
{{< /file >}}

{{< /files >}}

In these Kubernetes YAML files, there are various objects, separated by the `---`:

- A Deployment, describing a scalable group of identical pods. In this case,
  you'll get just one replica, or copy of your pod. That pod, which is
  described under `template`, has just one container in it. The
  container is created from the image built by GitHub Actions in [Configure CI/CD for
  your Python application](configure-github-actions.md).
- A Service, which will define how the ports are mapped in the containers.
- A PersistentVolumeClaim, to define a storage that will be persistent through restarts for the database.
- A Secret, which stores the database password as a Kubernetes Secret resource.
- A NodePort service, which will route traffic from port 30001 on your host to
  port 8000 inside the pods it routes to, so you can reach your app
  from the network.

To learn more about Kubernetes objects, see the [Kubernetes documentation](https://kubernetes.io/docs/home/).

> [!NOTE]
>
> The `NodePort` service is good for development and testing. For production, implement an [ingress controller](https://kubernetes.io/docs/concepts/services-networking/ingress-controllers/) instead.

### Deploy and check your application

1. In a terminal, navigate to `python-docker-example` and deploy your database to
   Kubernetes.

   ```console
   $ kubectl apply -f docker-postgres-kubernetes.yaml
   ```

   You should see output that looks like the following, indicating your Kubernetes objects were created successfully.

   ```console
   deployment.apps/postgres created
   service/postgres created
   persistentvolumeclaim/postgres-pvc created
   secret/postgres-secret created
   ```

   Now, deploy your Python application.

   ```console
   $ kubectl apply -f docker-python-kubernetes.yaml
   ```

   You should see output that looks like the following, indicating your Kubernetes objects were created successfully.

   ```console
   deployment.apps/docker-python-demo created
   service/service-entrypoint created
   ```

2. Make sure everything worked by listing your deployments.

   ```console
   $ kubectl get deployments
   ```

   Your deployment should be listed as follows:

   ```console
   NAME                 READY   UP-TO-DATE   AVAILABLE   AGE
   docker-python-demo   1/1     1            1           48s
   postgres             1/1     1            1           2m39s
   ```

   This indicates all one of the pods you asked for in your YAML are up and running. Do the same check for your services.

   ```console
   $ kubectl get services
   ```

   You should get output like the following.

   ```console
   NAME                 TYPE        CLUSTER-IP     EXTERNAL-IP   PORT(S)          AGE
   kubernetes           ClusterIP   10.43.0.1      <none>        443/TCP          13h
   postgres             ClusterIP   10.43.209.25   <none>        5432/TCP         3m10s
   service-entrypoint   NodePort    10.43.67.120   <none>        8000:30001/TCP   79s
   ```

   In addition to the default `kubernetes` service, you can see your `service-entrypoint` service, accepting traffic on port 30001/TCP and the internal `ClusterIP` `postgres` with the port `5432` open to accept connections from your Python app.

3. In a terminal, curl the root endpoint to verify the application is running.

   ```console
   $ curl http://localhost:30001/
   Hello, Docker!
   ```

4. Exercise the database by creating a hero with a POST request:

   ```console
   $ curl -X 'POST' \
     'http://localhost:30001/heroes/' \
     -H 'accept: application/json' \
     -H 'Content-Type: application/json' \
     -d '{
     "id": 1,
     "name": "my hero",
     "secret_name": "austing",
     "age": 12
   }'
   ```

   You should receive the following response:

   ```json
   {
     "age": 12,
     "id": 1,
     "name": "my hero",
     "secret_name": "austing"
   }
   ```

   Then read it back with a GET request:

   ```console
   $ curl http://localhost:30001/heroes/
   ```

   You should receive an array containing the hero you just created. This
   confirms the application can read from and write to the PostgreSQL database
   running in the cluster.

5. Run the following commands to tear down your application.

   ```console
   $ kubectl delete -f docker-python-kubernetes.yaml
   $ kubectl delete -f docker-postgres-kubernetes.yaml
   ```

### Summary

In this section, you learned how to use Docker Desktop to deploy your application to a fully-featured Kubernetes environment on your development machine.

Related information:

- [Kubernetes documentation](https://kubernetes.io/docs/home/)
- [Deploy on Kubernetes with Docker Desktop](/manuals/desktop/use-desktop/kubernetes.md)
- [Use a Docker Hardened Image with Kubernetes](/dhi/how-to/use/#use-with-kubernetes)
