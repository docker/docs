---
title: Docker+Wasm (Beta)
description: How to use the Wasm integration in Docker Desktop
keywords: Docker, WebAssembly, wasm, containerd, engine
toc_max: 3
---

This page provides information about the ongoing integration of `containerd` for image and file system management in the Docker Engine.

> **Beta**
>
> The Docker+Wasm feature is currently in [Beta](../../release-lifecycle.md/#beta). We recommend that you do not use this feature in production environments as this feature may change or be removed from future releases.


## Enabling the Docker+Wasm integration

The Docker+Wasm integration currently requires a special build of Docker Desktop.

- **Important note #1:** This is a technical preview build of Docker Desktop and things might not work as expected. Be sure to back up your containers and images before proceeding.
- **Important note #2:** This preview has the containerd image store enabled and cannot be disabled. If you’re not currently using the containerd image store, then pre-existing images and containers will be inaccessible.

You can download the technical preview build of Docker Desktop here:

- [macOS Apple Silicon](https://dockr.ly/3sf56vH)
- [macOS Intel](https://dockr.ly/3VF6uFB)
- [Windows AMD64](https://dockr.ly/3ShlsP0)
- Linux Arm64 ([deb](https://dockr.ly/3TDcjRV))
- Linux AMD64 ([deb](https://dockr.ly/3TgpWH8), [rpm](https://dockr.ly/3eG6Mvp), [tar](https://dockr.ly/3yUhdCk))


## Usage examples

### Running a Wasm application with docker run

```
$ docker run -dp 8080:8080 \
  --name=wasm-example \
  --runtime=io.containerd.wasmedge.v1 \
  --platform=wasi/wasm32 \
  michaelirwin244/wasm-example
```

Note the addition of two additional flags to the run command:

- **--runtime=io.containerd.wasmedge.v1** - This informs the Docker engine that we want to use the Wasm containerd shim instead of the standard Linux container runtime
- **--platform=wasi/wasm32** - This specifies the architecture of the image we want to use. By leveraging a Wasm architecture, we don’t need to build separate images for the different machine architectures. The Wasm runtime will do the final step of converting the Wasm binary to machine instructions.

### Running a Wasm application with Docker Compose

The same application can be run using the following Docker Compose file:

```yaml
services:
  app:
    image: michaelirwin244/wasm-example
    platform: wasi/wasm32
    runtime: io.container.wasmedge.v1
    ports:
      - 8080:8080
```

Then start the application using the normal Docker Compose commands:

```
docker compose up
```


### Running a multi-service application with Wasm

Networking works the same as you expect with Linux containers, giving you the flexibility to combine Wasm applications with other containerized workloads (such as a database) in a single application stack.

In this example, the Wasm application will leverage a MariaDB database running in a container.

1. Start by cloning the repository.

    ```
    $ git clone https://github.com/second-state/microservice-rust-mysql.git
    Cloning into 'microservice-rust-mysql'...
    remote: Enumerating objects: 75, done.
    remote: Counting objects: 100% (75/75), done.
    remote: Compressing objects: 100% (42/42), done.
    remote: Total 75 (delta 29), reused 48 (delta 14), pack-reused 0
    Receiving objects: 100% (75/75), 19.09 KiB | 1.74 MiB/s, done.
    Resolving deltas: 100% (29/29), done.
    ```

2. Navigate into the cloned project and start the project using Docker Compose.

    ```
    $ cd microservice-rust-mysql
    $ docker compose up
    [+] Running 0/1
    ⠿ server Warning                                                                                                  0.4s
    [+] Building 4.8s (13/15)
    ...
    microservice-rust-mysql-db-1      | 2022-10-19 19:54:45 0 [Note] mariadbd: ready for connections.
    microservice-rust-mysql-db-1      | Version: '10.9.3-MariaDB-1:10.9.3+maria~ubu2204'  socket: '/run/mysqld/mysqld.sock'  port: 3306  mariadb.org binary distribution
    ```

3. In another terminal, we can see the Wasm image that was created.

    ```
    $ docker images
    REPOSITORY   TAG       IMAGE ID       CREATED         SIZE
    server       latest    2c798ddecfa1   2 minutes ago   3MB
    ```

4. Inspecting the image will show the image has a `wasi/wasm32` platform (combination of Os and Architecture).

    ```
    $ docker image inspect server | grep -A 3 "Architecture"
            "Architecture": "wasm32",
            "Os": "wasi",
            "Size": 3001146,
            "VirtualSize": 3001146,
    ```

5. Open the website at http://localhost:8090 and create a few sample orders. All of these are interacting with the Wasm server.

6. When you're all done, tear everything down by hitting Ctrl+C in the terminal you launched the application.


### Building and pushing a Wasm module

1. Create a Dockerfile that will build your Wasm application. This will vary depending on the language you are using.

2. In a separate stage in your `Dockerfile`, extract the module and set it as the `ENTRYPOINT`.

    ```
    FROM scratch
    COPY --from=build /build/hello_world.wasm /hello_world.wasm
    ENTRYPOINT [ "hello_world.wasm" ]
    ```

3. Build and push the image specifying the `wasi/wasm32` architecture. Buildx makes this easy to do in a single command.

    ```
    $ docker buildx build --platform wasi/wasm32 -t username/hello-world .
    ...
    => exporting to image                                                                             0.0s
    => => exporting layers                                                                            0.0s
    => => exporting manifest sha256:2ca02b5be86607511da8dc688234a5a00ab4d58294ab9f6beaba48ab3ba8de56  0.0s
    => => exporting config sha256:a45b465c3b6760a1a9fd2eda9112bc7e3169c9722bf9e77cf8c20b37295f954b    0.0s
    => => naming to docker.io/username/hello-world:latest                                            0.0s
    => => unpacking to docker.io/username/hello-world:latest                                         0.0s
    $ docker push username/hello-world
    ```


## Docker+Wasm Release Notes

(2022-10-24)  
Initial release

### New
- Initial implementation of Wasm integration

### Known issues
- Docker Compose may not exit cleanly when interrupted
    - Workaround: Clean up `docker-compose` processes by sending them a SIGKILL (`killall -9 docker-compose`).
- Pushes to Hub might give an error stating `server message: insufficient_scope: authorization failed`, even after logging in using Docker Desktop
    - Workaround: Run `docker login` in the CLI

## Feedback

Thanks for trying the new Docker+Wasm integration. We’d love to hear from you! Please feel free to give feedback or report any bugs you may find through the issues tracker on the [public roadmap item](https://github.com/docker/roadmap/issues/426){: target="_blank" rel="noopener" class="_"}.
