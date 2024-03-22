---
title: What is a container?
keywords: concepts, build, images, container, docker desktop
description: What is a container? This concept page will teach you about containers and provide a quick hands-on where you will run your first container.
---

<iframe width="650" height="365" src="https://www.youtube.com/embed/nsWWQ1xoEy0?rel=0" title="YouTube video player" frameborder="0" allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture; web-share" allowfullscreen></iframe>

## Explanation

Imagine you’re developing a killer web app. It has three main parts - A `React` frontend, a `Python` API, and a `PostgreSQL` database.

- **React app:** A user interface that lets users browse the frontend.
- **Python API:**  An engine that powers the application logic.
- **Database:**  Used for storing all the essential information for your application.

Traditionally, all these components are bundled together in one single system. It works but it’s difficult to manage and scale. If one component fails, it might affect the whole application and result in downtime. 

Containers are like magic boxes for your app's components! Each component - the frontend React app, the Python API engine, and the database - gets its own secure container. Think of it like providing each app on your phone its own isolated space.

Here's what makes them awesome:

- **Self-Contained:** Each container has everything it needs to function perfectly, like a phone app with all its features built-in.
- **Isolated:** Containers are like separate sandboxes on the server. They don't see or influence each other, so even if one part malfunctions, it won't affect the others (like the database).
- **Independent:** You can manage each container independently, just like installing, updating, or deleting apps on your phone. Deleting one container (like the ordering app) won't affect the others.
- **Portable:** Containers can run on a developer’s laptop, VMs, in a data center, and on a cloud provider.


> Containers make your life easier by keeping things organized and independent. No more compatibility headaches! You can develop, test, and deploy your web app like a pro, with each part running smoothly – just like having separate, well-functioning apps on your phone!

## Try it now

In this hands-on, you will see how to run a simple Docker container using Docker Desktop GUI.



{{< tabs >}}
{{< tab name="Using Docker Desktop" >}}

## Run a container

Use the following instructions to run a container.

1. [Download  and install](https://www.docker.com/products/docker-desktop/) Docker Desktop.

2. Open Docker Desktop and select the **Search** on the top navigation bar.

3. Specify `welcome-to-docker` in the search and then select **Pull**.

![A screenshot of the Docker Desktop dashboard showing the search result for welcome-to-docker Docker image ](images/search-the-docker-image.webp?border=true&w=1000&h=700)

4. Once the image is successfully pulled, select **Run**.

5. Expand the **Optional settings**.

6. In the **Container name**, specify `welcome-to-docker`.

7. In the **Host port**, specify `8080`.

![A screenshot of Docker Desktop dashboard showing the container run dialog with welcome-to-docker typed in as the container name and 8080 specified as the port number](images/run-a-new-container.webp?border=true&w=550&h=400)
8. Click **Run**.

## View your container

Congratulations!! You just ran your first container! 🍻
 
You can view it in the **Containers** tab of the Docker Desktop GUI.


![Screenshot of the container tab of the Docker Desktop GUI showing the welcome-to-docker container running on the host port 8080](images/view-your-containers.webp?border=true&w=750&h=600)

This container runs a web server that displays a simple website. When working with more complex projects, you'll run different parts in different containers. For example, a different container for the frontend, backend, and database. In this walkthrough, you only have a simple frontend container.

## Access the frontend

The frontend is accessible on port `8080` of your local host. Select the link in the **Port(s)** column of your container, or visit [http://localhost:8080](https://localhost:8080)  in your browser to view it.

![Screenshot of the landing page of the nginx web server, coming from the running container](images/access-the-frontend.webp?border)

## Explore your container

Docker Desktop lets you easily view and interact with different aspects of your container. Try it out yourself. Select your container and then click on **Files** to explore your container's isolated file system.


![Screenshot of the Docker Desktop GUI showing the files and directories inside a running container](images/explore-your-container.webp?border)


## Stop your container

The `docker/welcome-to-docker` container continues to run until you stop it. To stop the container, go to the **Containers** tab and select the **Stop** icon in the **Actions** column of your container.


![Screenshot of the Docker Desktop GUI with the nginx container selected and being prepared to stop](images/stop-your-container.webp?border)

{{< /tab >}}
{{< tab name="CLI" >}}

## Run a container

Follow the instructions to to run a container using your CLI terminal:

1. [Download and install](https://www.docker.com/products/docker-desktop/) Docker Desktop

2. Open your CLI terminal and run the following command:

```console
  docker run -d -p 8080:80 docker/welcome-to-docker
```

Congratulations! You just fired up your first container! 🍻
## View your container

You can verify if the container is up and running by running the following command:

```console
  docker ps
```

You will see output like the following:

```console
 CONTAINER ID   IMAGE         COMMAND                  CREATED          STATUS          PORTS                    NAMES
 a1f7a4bb3a27   nginx         "/docker-entrypoint.…"   11 seconds ago   Up 11 seconds   0.0.0.0:80->80/tcp       gracious_keldysh
```

This output is the result of running the `docker ps` command. It shows you information about the containers currently running on your Docker engine.

This container runs a web server that displays a simple website. When working with more complex projects, you'll run different parts in different containers. For example, a different container for the `frontend`, `backend`, and `database`. In this walkthrough, you only have a simple frontend container.


## Access the frontend

The `frontend` is accessible on port `8080` of your localhost. Visit [http://localhost:8080](http://localhost:8080) in your browser to view it.

![Screenshot of the landing page of the nginx web server, coming from the running container](images/access-the-frontend.webp?border)

## Stop your container

The `docker/welcome-to-docker` container continues to run until you stop it. You can stop a container using a single Docker command. 

```console
 docker stop <the-container-id>
```

In order to remove a container, you can use `docker rm <the-container-id>` command.

Now that you have learned the basics of a Docker container, it's time to learn about Docker image and its layered architecture.


{{< /tab >}}
{{< /tabs >}}

Now that you have learned the basics of a Docker container, it's time to learn about Docker image and it's layered architecture.

## Additional resources

- [Running a container](https://docs.docker.com/engine/reference/run/)
- [What is a container](https://docs.docker.com/guides/walkthroughs/what-is-a-container/)
- [Overview of container](https://www.docker.com/resources/what-container/)
- [Why Docker?](https://www.docker.com/why-docker/)

{{< button text="What is an image" url="what-is-an-image" >}}
