---
title: Containerize a Python application
keywords: python, flask, containerize, initialize
description: Learn how to containerize a Python application.
aliases:
  - /language/python/build-images/
  - /language/python/run-containers/
---

## Prerequisites

* You have installed the latest version of [Docker Desktop](../../get-docker.md).
* You have a [git client](https://git-scm.com/downloads). The examples in this section use a command-line based git client, but you can use any client.

## Overview

This section walks you through containerizing and running a Python application.

## Get the sample application

The sample application uses the popular [Flask](https://flask.palletsprojects.com/) framework.

Clone the sample application to use with this guide. Open a terminal, change directory to a directory that you want to work in, and run the following command to clone the repository:

```console
$ git clone https://github.com/docker/python-docker
```

## Test the application without Docker (optional)

You can test the application locally without Docker before you continue building and running the application with Docker. This section requires you to have Python 3.11 or later installed on your machine. Download and install [Python](https://www.python.org/downloads/).

Open your terminal and navigate to the working directory you created. Create an environment, install the dependencies, and start the application to make sure it’s running.

```console
$ cd /path/to/python-docker
$ python3 -m venv .venv
$ source .venv/bin/activate
(.venv) $ python3 -m pip install -r requirements.txt
(.venv) $ python3 -m flask run
```

To test that the application is working, open a new browser and navigate to `http://localhost:5000`.

Switch back to the terminal where the server is running and you should see the following requests in the server logs. The data and timestamp will be different on your machine.

```shell
127.0.0.1 - - [22/Sep/2020 11:07:41] "GET / HTTP/1.1" 200 -
```

## Initialize Docker assets

Now that you have an application, you can use `docker init` to create the necessary Docker assets to containerize your application. Inside the `python-docker` directory, run the `docker init` command. Refer to the following example to answer the prompts from `docker init`.

```console
$ docker init
Welcome to the Docker Init CLI!

This utility will walk you through creating the following files with sensible defaults for your project:
  - .dockerignore
  - Dockerfile
  - compose.yaml

Let's get started!

? What application platform does your project use? Python
? What version of Python do you want to use? 3.11.4
? What port do you want your app to listen on? 5000
? What is the command to run your app? python3 -m flask run --host=0.0.0.0
```

You should now have the following contents in your `python-docker`
directory.

   ```
   ├── python-docker/
   │ ├── app.py
   │ ├── requirements.txt
   │ ├── .dockerignore
   │ ├── compose.yaml
   │ ├── Dockerfile
   │ └── README.md
   ```

To learn more about the files that `docker init` added, see the following:
 - [Dockerfile](../../engine/reference/builder.md)
 - [.dockerignore](../../engine/reference/builder.md#dockerignore-file)
 - [compose.yaml](../../compose/compose-file/_index.md)

## Run the application

Inside the `python-docker` directory, run the following command in a
terminal.

```console
$ docker compose up --build
```

Open a browser and view the application at [http://localhost:5000](http://localhost:5000). You should see a simple Flask application.

In the terminal, press `ctrl`+`c` to stop the application.

### Run the application in the background

You can run the application detached from the terminal by adding the `-d`
option. Inside the `python-docker` directory, run the following command
in a terminal.

```console
$ docker compose up --build -d
```

Open a browser and view the application at [http://localhost:5000](http://localhost:5000).

You should see a simple Flask application.

In the terminal, run the following command to stop the application.

```console
$ docker compose down
```

For more information about Compose commands, see the [Compose CLI
reference](../../compose/reference/_index.md).

## Summary

In this section, you learned how you can containerize and run your Python
application using Docker.

Related information:
 - [Build with Docker guide](../../build/guide/index.md)
 - [Docker Compose overview](../../compose/_index.md)

## Next steps

In the next section, you'll learn how you can develop your application using
containers.

{{< button text="Develop your application" url="develop.md" >}}
