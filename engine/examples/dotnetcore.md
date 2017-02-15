---
description: Create a Docker image by layering your ASP.NET Core app on debian for Linux Containers or with Windows Nano Server containers using a Dockerfile.
keywords: docker, dockerize, dockerizing, dotnet, .NET, Core, article, example, platform, installation, containers, images, image, dockerfile, build, asp.net, asp.net core
title: Dockerizing a .NET Application
---

## Introduction

This example demonstrates how to dockerize an ASP.NET Core application.

> **Note:** This example assumes you already have a published ASP.NET Core app on your machine. If you are new to ASP.NET you can follow a [simple tutorial](https://www.asp.net/get-started) to initialize a project.

## Why build ASP.NET Core?
- [Open-source](https://github.com/aspnet/home)
- Develop and run your ASP.NET Core apps cross-platform on Windows, Mac and Linux
- Great for modern cloud-based apps, such as web apps, IoT apps and mobile backends 
- ASP.NET Core apps can run on [.NET Core](https://www.microsoft.com/net/core/platform) or on the full [.NET Framework](https://www.microsoft.com/net/framework)
- Architected to provide an optimized development framework for apps that are deployed to the cloud or run on-premises
- Modular components with minimal overhead retain flexibility while constructing your solutions

## Create a Dockerfile for an ASP.NET Core Application
1. If you haven't already, publish your project locally by running `dotnet restore` and `dotnet publish -o published`.
2. Create a `Dockerfile` in your project folder next to your published folder.
3. Add the text below to your `Dockerfile` for either Linux or [Windows Containers](https://docs.microsoft.com/en-us/virtualization/windowscontainers/about/).
4. The `Dockerfile` assumes that your application is called `aspnetapp`. Change the `Dockerfile` to use the .dll file of your project.

  For Linux Containers:
  ```dockerfile
  FROM microsoft/aspnetcore:1.1
  WORKDIR /app
  COPY published ./
  ENTRYPOINT ["dotnet", "aspnetapp.dll"]
  ```
  For Windows Containers:
  ```dockerfile
  FROM microsoft/dotnet:1.1-runtime-nanoserver
  WORKDIR /app
  ENV ASPNETCORE_URLS http://+:80
  COPY published ./
  ENTRYPOINT ["dotnet", "aspnetapp.dll"]
  ```
5. To make your build context as small as possible add a [`.dockerignore` file](https://docs.docker.com/engine/reference/builder/#dockerignore-file) to your project folder and copy the following into it.
```dockerignore
*
!published
```
> **Note:** If you are using Windows Containers on [Docker for Windows](https://docs.docker.com/docker-for-windows/) be sure to check that you are properly switched to Windows Containers. Do this by opening the system tray up arrow and right clicking on the Docker whale icon for a popup menu. In the popup menu make sure you select 'Switch to Windows Containers'. 

## Build and run the Docker image
1. Open the command prompt and navigate to your project folder.
2. Use the following commands to build and run your Docker image:
```console
docker build -t aspnetapp .
docker run -d -p 80:80 aspnetapp
```
## View your web page running from your container
* If you are using a Linux container you can simply browse to http://localhost:80 to access your app in a web browser.
* If you are using the Nano [Windows Container](https://docs.docker.com/docker-for-windows/) there is currently a bug that affects how [Windows 10 talks to Containers via "NAT"](https://github.com/Microsoft/Virtualization-Documentation/issues/181#issuecomment-252671828) (Network Address Translation). Today you must hit the IP of the container directly. You can get the IP address of your container with the following steps:
  1. Run `docker ps` and copy the hash of your container ID.
  2. Run `docker inspect -f "{{ .NetworkSettings.Networks.nat.IPAddress }}" HASH` where `HASH` is replaced with your container ID.
  3. Copy the container ip address and paste into your browser with port 80. (Ex: 172.16.240.197:80)

## Further reading
  - [ASP.NET Core](https://docs.microsoft.com/en-us/aspnet/core/)
  - [Microsoft ASP.NET Core on Docker Hub](https://hub.docker.com/r/microsoft/aspnetcore/)
  - [ASP.NET Core with Docker Tools for Visual Studio](https://blogs.msdn.microsoft.com/webdev/2016/11/16/new-docker-tools-for-visual-studio/)