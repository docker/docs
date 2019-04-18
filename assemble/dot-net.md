---
title: Build a C# ASP.NET Core project
description: Building a C# ASP.NET Core project using Docker Assemble
keywords: Assemble, Docker Enterprise, Spring Boot, container image
---

Ensure you are running the `backend` before you build any projects using Docker Assemble. For instructions on running the backend, see [Install Docker Assemble](/assemble/install).

Clone the git repository you would like to use. The following example uses the `dotnetdemo` repository.

```
~$ git clone https://github.com/mbentley/dotnetdemo
Cloning into 'dotnetdemo'...
«…»
```

Build the project using the `docker assemble build` command by passing it the path to the source repository (or a subdirectory in the following example):

```
~$ docker assemble build dotnetdemo/dotnetdemo
«…»
Successfully built: docker.io/library/dotnetdemo:latest
```
The resulting image is exported to the local Docker image store using a name and a tag which are automatically determined by the project metadata.

```
~$ docker image ls
REPOSITORY      TAG             IMAGE ID            CREATED             SIZE
dotnetdemo      latest          a055e61e3a9e        24 seconds ago      349MB
```

An image name consists of `«namespace»/«name»:«tag»`. Where, `«namespace»/` is optional and defaults to `none`. If the project metadata does not contain a ‘tag’ (or a version), then latest is used. If the project metadata does not contain a ‘name’ and it was not provided on the command line, then a fatal error occurs.

Use the `--namespace`, `--name` and `--tag` options to override each element of the image name:

```
~$ docker assemble build --name testing --tag latest dotnetdemo/
«…»
INFO[0007] Successfully built "testing:latest"
~$ docker image ls
REPOSITORY       TAG            IMAGE ID            CREATED             SIZE
testing          latest         d7f41384814f        32 seconds ago      97.4MB
hello-boot       1              0dbc2c425cff        5 minutes ago       97.4MB
```

Run the container:

```
~$ docker run -d --rm -p 8080:80 dotnetdemo:latest
e1c54291e96967dad402a81c4217978a544e4d7b0fdd3c0a2e2cca384c3b4adb
~$ docker ps
CONTAINER ID    IMAGE             COMMAND                  «…» PORTS                    NAMES
e1c54291e969    dotnetdemo:latest "dotnet dotnetdemo.d…"   «…» 0.0.0.0:8080->80/tcp     lucid_murdock
~$ docker logs e1c54291e969
warn: Microsoft.AspNetCore.DataProtection.KeyManagement.XmlKeyManager[35]
      No XML encryptor configured. Key {11bba23a-71ad-4191-b583-4f974e296033} may be persisted to storage in unencrypted form.
Hosting environment: Production
Content root path: /app
Now listening on: http://[::]:80
Application started. Press Ctrl+C to shut down.
~$ curl -s localhost:8080 | grep '<h4>'
<h4>This environment is </h4>
<h4>served from e1c54291e969 at 11/22/2018 16:00:23</h4>
~$ docker rm -f e1c54291e969
```
