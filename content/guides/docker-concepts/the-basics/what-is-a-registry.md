---
title: What is a registry?
keywords: concepts, build, images, container, docker desktop
description: What is a registry? This Docker Concept will explain what a registry is, explore their interoperability, and have you interact with registries.
---

<iframe width="650" height="365" src="https://www.youtube.com/embed/nsWWQ1xoEy0?rel=0" title="YouTube video player" frameborder="0" allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture; web-share" allowfullscreen></iframe>

## Explanation

Now that you know what a container image is and how it works, you might wonder- Where do you store these images? 

Well, you can store your container images on your computer system, but what if you want to share them with your friends or use them on another machine? That’s where the image registry comes in.

An image registry is a centralized location for storing and sharing your container images. It can be either public or private. Remember we talked about [Docker Hub](https://hub.docker.com)? It is a public registry that anyone can use, and Docker looks for images on `Docker Hub` by default. While Docker Hub is a popular option, there's a whole world of public container registries available today, including [Amazon Elastic Container Registry(ECR)](https://aws.amazon.com/ecr/), [Azure Container Registry (ACR)](https://azure.microsoft.com/en-in/products/container-registry), [Google Container Registry (GCR)](https://cloud.google.com/artifact-registry) etc. You can even run your private registry in your local system or inside your organization. For example, [VMware Harbor](https://docs.vmware.com/en/VMware-Harbor-Registry/services/vmware-harbor-registry/GUID-index.html), [JFrog Artifactory](https://jfrog.com/artifactory/), [GitLab Container registry](https://docs.gitlab.com/ee/user/packages/container_registry/) etc.


## Registry vs repository

While you’re working with `Docker Hub`, you might hear users using terms like **registry** and **repository** like they’re interchangeable. Even though they’re related, they’re not quite the same thing.

A registry is a centralized location that stores and manages container images whereas a repository is a collection of related container images within a registry. Think of it as a folder where you organize your images based on projects. Each repository contains one or more container images.

> You can create one private repository and unlimited public repositories using the free version of Docker Hub. For more information, visit the Docker Hub subscription page.


## Try it now

In this hands-on, you will learn how to build and push a Docker image to the Docker Hub repository.


### Sign up for a free Docker account {#sign-up-for-a-free-docker-account}



1. [Click here](https://hub.docker.com/signup) to create a new Docker account, if not already done.

![Screenshot of the official Docker Hub page showing the Sign up page](images/dockerhub-signup.webp?border)


You can use your Google account or GitHub credentials in order to sign-up. If you already have a Docker account, you can directly click on **Sign-in** to enter your credentials.


### Create your first repository  {#create-your-first-repository}



1. Once you login to Docker Hub, click on the **Create your repository** button.
2. Choose your namespace (ie. your Docker account) and your preferred repository name. For this demonstration, let’s name it `docker-quickstart`. Enter a short description to identify your repository.


![Screenshot of the Docker Hub page that shows how to create a public repository](images/create-hub-repository.webp?border)


3. Set the visibility to **Public**. 
4. Select **Create**.

That’s it. You've successfully created your first repository.🍻

This repository is empty right now. Let’s push some images to it to see them appear on the Docker Hub. 


### Log into Docker Hub using the CLI



1. [Download and install](https://www.docker.com/products/docker-desktop/) Docker Desktop, if not already installed.
2. Open your CLI terminal and run the following command:

```console
 docker login
```


Alternatively, you can use the following command to pass your username to the command line directly.

```console
 docker login registry-1.docker.io -u <your-dockerhub-id>
```

At the password prompt, enter your personal access token(PAT). If you see the **Login Succeeded** message, it shows that you’re successfully logged in to Docker Hub.


### Cloning a sample Node.js code


Let's clone a sample Node.js project from a GitHub repository. This repository contains a pre-built Dockerfile necessary for building a Docker image.

Don't worry about the specifics of how the Dockerfile was constructed; we'll cover that in detail in later sections.



1. Clone the GitHub repository

```console
 git clone https://github.com/dockersamples/helloworld-demo-node
```


2. Build a Docker image

Change to the `helloworld-demo-node `directory and run the following command to build a Docker image.


```console
 docker build -t <your-dockerhub-id>/docker-quickstart .
```

3. Run the following command to list the newly created Docker image.

```console
 docker images
 REPOSITORY                 TAG       IMAGE ID       CREATED         SIZE
 <your-dockerhub-id>/docker-quickstart   latest    476de364f70e   2 minutes ago   170MB
```

4. You can execute the following command to run a Docker container

```console
 docker run -d -p 8080:8080 <your-dockerhub-id>/docker-quickstart 
```


You can quickly verify if the container is working fine or not by executing the following `curl `command. If the command is not available, you can directly access the app by visiting [http://localhost:8080](http://localhost:8080) on your browser.

```console
 curl localhost:8080
 
           ##         .
     ## ## ##        ==
  ## ## ## ## ##    ===
 /""""""""""""""""\___/ ===
 {                       /  ===-
 \______ O           __/
  \    \         __/
   \____\_______/
 
 
 Hello from Docker!
```

5. Next, you can use `docker tag` command to tag the Docker image. Docker tags are used for labeling and versioning Docker images, essentially serving as a convenient reference to a specific version of an image. 

```console 
 docker tag <your-dockerhub-id>/docker-quickstart <your-dockerhub-id>/docker-quickstart:1.0 
```



6. Finally, it’s time to push the newly built image to your Docker Hub repository

```console 
 docker push -u <your-dockerhub-id>/docker-quickstart:1.0
```


By now, you should be able to find the new image by navigating to “Tags” section of your Docker Hub repository.



![Screenshot of the Docker Hub page that displays the newly added image tag](images/dockerhub-tags.webp?border=true) 
In this walkthrough, you learned how to sign up for a Docker Hub account and create your first Docker Hub repository. You also learned how to build and tag a Docker image and push it to your Docker Hub repository.


### Additional resources

- [Docker Hub Quickstart](https://docs.docker.com/docker-hub/quickstart/)
- [Manage Docker Hub Repositories](https://docs.docker.com/docker-hub/repo)
- [How to tag an image](https://docs.docker.com/reference/cli/docker/image/tag/)

{{< button text="What is Docker Compose" url="what-is-Docker-Compose" >}}
