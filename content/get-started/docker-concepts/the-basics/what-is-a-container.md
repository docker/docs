---
title: What is a container?
weight: 10
keywords: concepts, build, images, container, docker desktop
description: What is a container? This concept page will teach you about containers and provide a quick hands-on where you will run your first container.
aliases:
- /guides/walkthroughs/what-is-a-container/
- /guides/walkthroughs/run-a-container/
- /guides/walkthroughs/
- /get-started/run-your-own-container/
- /get-started/what-is-a-container/
- /guides/docker-concepts/the-basics/what-is-a-container/
---

{{< youtube-embed W1kWqFkiu7k >}}

## Explanation

Imagine you're developing a killer web app that has three main components - a React frontend, a Python API, and a PostgreSQL database. If you wanted to work on this project, you'd have to install Node, Python, and PostgreSQL. 

How do you make sure you have the same versions as the other developers on your team? Or your CI/CD system? Or what's used in production?

How do you ensure the version of Python (or Node or the database) your app needs isn't affected by what's already on your machine? How do you manage potential conflicts?

Enter containers!

What is a container? Simply put, containers are isolated processes for each of your app's components. Each component - the frontend React app, the Python API engine, and the database - runs in its own isolated environment, completely isolated from everything else on your machine. 

Here's what makes them awesome. Containers are:

- Self-contained. Each container has everything it needs to function with no reliance on any pre-installed dependencies on the host machine.
- Isolated. Since containers are run in isolation, they have minimal influence on the host and other containers, increasing the security of your applications.
- Independent. Each container is independently managed. Deleting one container won't affect any others.
- Portable. Containers can run anywhere! The container that runs on your development machine will work the same way in a data center or anywhere in the cloud!

### Containers versus virtual machines (VMs)

Without getting too deep, a VM is an entire operating system with its own kernel, hardware drivers, programs, and applications. Spinning up a VM only to isolate a single application is a lot of overhead.

A container is simply an isolated process with all of the files it needs to run. If you run multiple containers, they all share the same kernel, allowing you to run more applications on less infrastructure.

> **Using VMs and containers together**
>
> Quite often, you will see containers and VMs used together. As an example, in a cloud environment, the provisioned machines are typically VMs. However, instead of provisioning one machine to run one application, a VM with a container runtime can run multiple containerized applications, increasing resource utilization and reducing costs.


## Try it out

In this hands-on, you will see how to run a Docker container using the Docker Desktop GUI.

{{< tabs group=concept-usage persist=true >}}
{{< tab name="Using the GUI" >}}

Use the following instructions to run a container.

1. Open Docker Desktop and select the **Search** field on the top navigation bar.

2. Specify `welcome-to-docker` in the search input and then select the **Pull** button.

    ![A screenshot of the Docker Desktop Dashboard showing the search result for welcome-to-docker Docker image ](images/search-the-docker-image.webp?border=true&w=1000&h=700)

3. Once the image is successfully pulled, select the **Run** button.

4. Expand the **Optional settings**.

5. In the **Container name**, specify `welcome-to-docker`.

6. In the **Host port**, specify `8080`.

    ![A screenshot of Docker Desktop Dashboard showing the container run dialog with welcome-to-docker typed in as the container name and 8080 specified as the port number](images/run-a-new-container.webp?border=true&w=550&h=400)

7. Select **Run** to start your container.

Congratulations! You just ran your first container! 🎉
 
### View your container

You can view all of your containers by going to the **Containers** view of the Docker Desktop Dashboard.

![Screenshot of the container view of the Docker Desktop GUI showing the welcome-to-docker container running on the host port 8080](images/view-your-containers.webp?border=true&w=750&h=600)

This container runs a web server that displays a simple website. When working with more complex projects, you'll run different parts in different containers. For example, you might run a different container for the frontend, backend, and database.

### Access the frontend

When you launched the container, you exposed one of the container's ports onto your machine. Think of this as creating configuration to let you to connect through the isolated environment of the container. 

For this container, the frontend is accessible on port `8080`. To open the website, select the link in the **Port(s)** column of your container or visit [http://localhost:8080](http://localhost:8080) in your browser.

![Screenshot of the landing page coming from the running container](images/access-the-frontend.webp?border)

### Explore your container

Docker Desktop lets you explore and interact with different aspects of your container. Try it out yourself. 

1. Go to the **Containers** view in the Docker Desktop Dashboard.

2. Select your container.

3. Select the **Files** tab to explore your container's isolated file system.

    ![Screenshot of the Docker Desktop Dashboard showing the files and directories inside a running container](images/explore-your-container.webp?border)

### Stop your container

The `docker/welcome-to-docker` container continues to run until you stop it. 

1. Go to the **Containers** view in the Docker Desktop Dashboard.

2. Locate the container you'd like to stop.

3. Select the **Stop** action in the **Actions** column.

    ![Screenshot of the Docker Desktop Dashboard with the welcome container selected and being prepared to stop](images/stop-your-container.webp?border)

{{< /tab >}}
{{< tab name="Using the CLI" >}}

Follow the instructions to run a container using the CLI:

1. Open your CLI terminal and start a container by using the [`docker run`](/reference/cli/docker/container/run/) command:

    ```console
    $ docker run -d -p 8080:80 docker/welcome-to-docker
    ```

    The output from this command is the full container ID. 

Congratulations! You just fired up your first container! 🎉

### View your running containers

You can verify if the container is up and running by using the [`docker ps`](/reference/cli/docker/container/ls/) command:

```console
docker ps
```

You will see output like the following:

```console
 CONTAINER ID   IMAGE                      COMMAND                  CREATED          STATUS          PORTS                      NAMES
 a1f7a4bb3a27   docker/welcome-to-docker   "/docker-entrypoint.…"   11 seconds ago   Up 11 seconds   0.0.0.0:8080->80/tcp       gracious_keldysh
```

This container runs a web server that displays a simple website. When working with more complex projects, you'll run different parts in different containers. For example, a different container for the `frontend`, `backend`, and `database`.

> [!TIP]
>
> The `docker ps` command will show you _only_ running containers. To view stopped containers, add the `-a` flag to list all containers: `docker ps -a`


### Access the frontend

When you launched the container, you exposed one of the container's ports onto your machine. Think of this as creating configuration to let you to connect through the isolated environment of the container. 

For this container, the frontend is accessible on port `8080`. To open the website, select the link in the **Port(s)** column of your container or visit [http://localhost:8080](http://localhost:8080) in your browser.

![Screenshot of the landing page of the Nginx web server, coming from the running container](images/access-the-frontend.webp?border)

### Stop your container

The `docker/welcome-to-docker` container continues to run until you stop it. You can stop a container using the `docker stop` command.

1. Run `docker ps` to get the ID of the container

2. Provide the container ID or name to the [`docker stop`](/reference/cli/docker/container/stop/) command:

    ```console
    docker stop <the-container-id>
    ```

> [!TIP]
>
> When referencing containers by ID, you don't need to provide the full ID. You only need to provide enough of the ID to make it unique. As an example, the previous container could be stopped by running the following command:
>
> ```console
> docker stop a1f
> ```

{{< /tab >}}
{{< /tabs >}}

## Additional resources

The following links provide additional guidance into containers:

- [Running a container](/engine/containers/run/)
- [Overview of container](https://www.docker.com/resources/what-container/)
- [Why Docker?](https://www.docker.com/why-docker/)

## Next steps

Now that you have learned the basics of a Docker container, it's time to learn about Docker images.

{{< button text="What is an image?" url="what-is-an-image" >}}
