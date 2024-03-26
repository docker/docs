---
title: What is Docker Compose?
keywords: concepts, build, images, container, docker desktop
description: What is Docker Compose?
---

{{< youtube-embed xhcUIK4fGtY >}}

## Explanation

Up to this point, you’ve been working with a single container application, but now you’re thinking, “Hey, I want to add another container such as a `PostgreSQL` database or a frontend `React` app. So, where do these new containers fit in? Do you need to install them in the same container or run it separately? In the world of containerisation, it’s generally a good idea for each container to do one thing and do it well. Here’s why:

* You might want to scale your Python APIs and React frontend app differently than your PostgreSQL databases. Having separate containers gives you the flexibility to tweak and adjust each part of your application independently.
* Imagine you want to update a specific version of your React frontend without messing with your PostgreSQL database. With separate containers, you can isolate these changes and update versions without causing a whole fuss.
* While you might use a PostgreSQL container for your database during development, you might opt for a managed PostgreSQL service in production. You don't want to carry around your database engine with your app wherever it goes, right?

You could use multiple `docker run` commands and link those containers but you might feel overwhelmed with everything you need to do to start up this application. Creating a network, starting containers in the right sequence, specifying all the environment variables, exposing ports and more! That’s a lot to remember and it’s certainly making things harder to pass along to someone else.

Thanks to Docker Compose. With `Docker Compose`, you can define all these containers and their configurations in a single file. Instead of running individual containers separately, it allows you to run multi-container applications using a single command.

One of the significant advantages of using Compose is developers can define their application stack in a file, keeping it at the root of your project repository, and easily enable someone else to contribute to your project. All someone would need to do is clone your repository and start the compose app. It's important to understand that Compose is a declarative tool - you simply define it and go. You don't always need to recreate everything from scratch. Just run `docker-compose up` again, and Compose will reconcile the changes in your file and apply them intelligently.

>**Tip**
>
>  Dockerfile provides instructions to build a specific container image while the Compose file defines your containers. The Compose file often references a Dockerfile to build an image to use for a particular service.


## Try it now 

In this hands-on, you will learn how to use a Docker compose to run multi-container applications. We’ll use a simple to-do list app built with Node.js and MySQL as a database server.

Follow the instructions to run the to-do list app on your system.

1. [Download  and install ](https://www.docker.com/products/docker-desktop/) Docker Desktop.
2. Open a terminal and [clone this sample application](https://github.com/dockersamples/todo-list-app).

>**Tip**
>
> If you have Git installed in your system, you can clone the repository for the sample application. Otherwise, you can download the sample application. You can choose your preferred options.

```console
 git clone https://github.com/dockersamples/todo-list-app 
```

3. Execute the following command to bring up the application

Change to the `todo-list-app` directory and start up the application stack by running the following command. 

```console
 docker compose up -d --build
```

When you run this command, you should see an output like this:

```console
 [+] Running 4/4
 ✔ app 3 layers [⣿⣿⣿]      0B/0B      Pulled                                                                   7.1s
   ✔ e6f4e57cc59e Download complete                                                                            0.9s
   ✔ df998480d81d Download complete                                                                            1.0s
   ✔ 31e174fedd23 Download complete                                                                            2.5s
 [+] Running 2/4
  ⠸ Network todo-list-app_default           Created                                                             0.3s
  ⠸ Volume "todo-list-app_todo-mysql-data"  Created                                                             0.3s
  ✔ Container todo-list-app-app-1           Started                                                             0.3s
  ✔ Container todo-list-app-mysql-1         Started                                                             0.3s
```

The Docker Compose configuration defines two services: a `Node.js` application for the to-do list and a `MySQL` database server. The application connects to the MySQL server using the environment variables to find the server and authenticate. The volume for the MySQL data ensures that the database data is persisted on the host machine. 

Don’t worry about the details of the services yet. 

4. Open up Docker Dashboard to view the containers.

![A screenshot of Docker Desktop dashboard showing the list of containers running todo-list app](images/todo-list-containers.webp?border=true&w=950&h=400)

You'll see two container images get downloaded from Docker Hub and, after a moment, the application will be up and running! No need to install or configure anything on your machine! Open [https://localhost:3000](https://localhost:3000) or simply click on Port `3000` to access the to-do list application. 

![A screenshot of a webpage showing the todo-list application running on port 3000](images/todo-list-app.webp?border=true&w=950&h=400)

### Tear it down

To remove the containers, you can select the application stack and select **Delete**.

![A screenshot of Docker Desktop Dashboard showing the containers tab with an arrow pointing to the “Delete” button](images/todo-list-delete.webp?w=930&h=400)

Alternatively, you can run the following CLI command to tear down the application stack.

```console
 docker compose down
```

In this walkthrough, you learned how to run a multi-container application using Docker Compose.

## Additional resources {#additional-resources}

* [Overview of Docker Compose CLI](https://docs.docker.com/compose/reference/)
* [Overview of Docker Compose](https://docs.docker.com/compose/)
* [How Compose works](https://docs.docker.com/compose/compose-application-model/)
