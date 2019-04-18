---
title: Build a Spring Boot project
description: Building a Spring Boot project using Docker Assemble
keywords: Assemble, Docker Enterprise, Spring Boot, container image
---

Ensure you are running the `backend` before you build any projects using Docker Assemble. For instructions on running the backend, see [Install Docker Assemble](/assemble/install).

Clone the git repository you would like to use. The following example uses the `docker-springfamework` repository.

```
~$ git clone https://github.com/anokun7/docker-springframework
Cloning into 'docker-springframework'...
«…»
```
When you build a Spring Boot project, Docker Assemble automatically detects the information it requires from the `pom.xml` project file.

Build the project using the `docker assemble build` command by passing it the path to the source repository:

```
~$ docker assemble build docker-springframework
«…»
Successfully built: docker.io/library/hello-boot:1
```
The resulting image is exported to the local Docker image store using a name and a tag which are automatically determined by the project metadata.

```
~$ docker image ls | head -n 2
REPOSITORY      TAG             IMAGE ID            CREATED           SIZE
hello-boot      1               00b0fbcf3c40        About a minute ago   97.4MB
```

An image name consists of `«namespace»/«name»:«tag»`. Where, `«namespace»/` is optional and defaults to `none`. If the project metadata does not contain a ‘tag’ (or a version), then `latest` is used. If the project metadata does not contain a ‘name’ and it was not provided on the command line, a fatal error occurs.

Use the `--namespace`, `--name` and `--tag` options to override each element of the image name:

```
~$ docker assemble build --name testing --tag latest docker-springframework/
«…»
INFO[0007] Successfully built "testing:latest"
~$ docker image ls
REPOSITORY       TAG            IMAGE ID            CREATED             SIZE
testing          latest         d7f41384814f        32 seconds ago      97.4MB
hello-boot       1              0dbc2c425cff        5 minutes ago       97.4MB
```

Run the container:

```
~$ docker run -d --rm -p 8080:8080 hello-boot:1
b2c88bdc35761ba2b99f85ce1f3e3ce9ed98931767b139a0429865cadb46ce13
~$ docker ps
CONTAINER ID    IMAGE           COMMAND                  «…» PORTS                    NAMES
b2c88bdc3576    hello-boot:1    "java -Djava.securit…"   «…» 0.0.0.0:8080->8080/tcp   silly_villani
~$ docker logs b2c88bdc3576

  .   ____          _            __ _ _
 /\\ / ___'_ __ _ _(_)_ __  __ _ \ \ \ \
( ( )\___ | '_ | '_| | '_ \/ _` | \ \ \ \
 \\/  ___)| |_)| | | | | || (_| |  ) ) ) )
  '  |____| .__|_| |_|_| |_\__, | / / / /
 =========|_|==============|___/=/_/_/_/
 :: Spring Boot ::        (v1.5.2.RELEASE)

«…» : Starting Application v1 on b2c88bdc3576 with PID 1 (/hello-boot-1.jar started by root in /)
«…»
~$ curl -s localhost:8080
Hello from b2c88bdc3576
~$ docker rm -f b2c88bdc3576
```
