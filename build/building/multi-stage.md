---
title: Multi-stage builds
description: Keeping your images small with multi-stage builds
keywords: build, best practices
redirect_from:
- /engine/userguide/eng-image/multistage-build/
- /develop/develop-images/multistage-build/
---

Multi-stage builds are useful to anyone who has struggled to optimize
Dockerfiles while keeping them easy to read and maintain.

> **Acknowledgment**
>
> Special thanks to [Alex Ellis](https://twitter.com/alexellisuk) for granting
> permission to use his blog post [Builder pattern vs. Multi-stage builds in Docker](https://blog.alexellis.io/mutli-stage-docker-builds/)
> as the basis of the examples below.

## Before multi-stage builds

One of the most challenging things about building images is keeping the image
size down. Each `RUN`, `COPY`, and `ADD` instruction in the Dockerfile adds a layer to the image, and you
need to remember to clean up any artifacts you don't need before moving on to
the next layer. To write a really efficient Dockerfile, you have traditionally
needed to employ shell tricks and other logic to keep the layers as small as
possible and to ensure that each layer has the artifacts it needs from the
previous layer and nothing else.

It was actually very common to have one Dockerfile to use for development (which
contained everything needed to build your application), and a slimmed-down one
to use for production, which only contained your application and exactly what
was needed to run it. This has been referred to as the "builder
pattern". Maintaining two Dockerfiles is not ideal.

Here's an example of a `build.Dockerfile` and `Dockerfile` which adhere to the
builder pattern above:

**`build.Dockerfile`**:

```dockerfile
# syntax=docker/dockerfile:1
FROM golang:1.16
WORKDIR /go/src/github.com/alexellis/href-counter/
COPY app.go ./
RUN go get -d -v golang.org/x/net/html \
  && CGO_ENABLED=0 go build -a -installsuffix cgo -o app .
```

Notice that this example also artificially compresses two `RUN` commands together
using the Bash `&&` operator, to avoid creating an additional layer in the image.
This is failure-prone and hard to maintain. It's easy to insert another command
and forget to continue the line using the `\` character, for example.

**`Dockerfile`**:

```dockerfile
# syntax=docker/dockerfile:1
FROM alpine:latest  
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY app ./
CMD ["./app"]
```

**`build.sh`**:

```bash
#!/bin/sh
echo Building alexellis2/href-counter:build
docker build -t alexellis2/href-counter:build . -f build.Dockerfile

docker container create --name extract alexellis2/href-counter:build  
docker container cp extract:/go/src/github.com/alexellis/href-counter/app ./app  
docker container rm -f extract

echo Building alexellis2/href-counter:latest
docker build --no-cache -t alexellis2/href-counter:latest .
rm ./app
```

When you run the `build.sh` script, it needs to build the first image, create
a container from it to copy the artifact out, then build the second
image. Both images take up room on your system and you still have the `app`
artifact on your local disk as well.

Multi-stage builds vastly simplify this situation!

## Use multi-stage builds

With multi-stage builds, you use multiple `FROM` statements in your Dockerfile.
Each `FROM` instruction can use a different base, and each of them begins a new
stage of the build. You can selectively copy artifacts from one stage to
another, leaving behind everything you don't want in the final image. To show
how this works, let's adapt the `Dockerfile` from the previous section to use
multi-stage builds.

```dockerfile
# syntax=docker/dockerfile:1

FROM golang:1.16
WORKDIR /go/src/github.com/alexellis/href-counter/
RUN go get -d -v golang.org/x/net/html  
COPY app.go ./
RUN CGO_ENABLED=0 go build -a -installsuffix cgo -o app .

FROM alpine:latest  
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=0 /go/src/github.com/alexellis/href-counter/app ./
CMD ["./app"]
```

You only need the single Dockerfile. You don't need a separate build script,
either. Just run `docker build`.

```console
$ docker build -t alexellis2/href-counter:latest .
```

The end result is the same tiny production image as before, with a
significant reduction in complexity. You don't need to create any intermediate
images, and you don't need to extract any artifacts to your local system at all.

How does it work? The second `FROM` instruction starts a new build stage with
the `alpine:latest` image as its base. The `COPY --from=0` line copies just the
built artifact from the previous stage into this new stage. The Go SDK and any
intermediate artifacts are left behind, and not saved in the final image.

## Name your build stages

By default, the stages are not named, and you refer to them by their integer
number, starting with 0 for the first `FROM` instruction. However, you can
name your stages, by adding an `AS <NAME>` to the `FROM` instruction. This
example improves the previous one by naming the stages and using the name in
the `COPY` instruction. This means that even if the instructions in your
Dockerfile are re-ordered later, the `COPY` doesn't break.

```dockerfile
# syntax=docker/dockerfile:1

FROM golang:1.16 AS builder
WORKDIR /go/src/github.com/alexellis/href-counter/
RUN go get -d -v golang.org/x/net/html  
COPY app.go ./
RUN CGO_ENABLED=0 go build -a -installsuffix cgo -o app .

FROM alpine:latest  
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /go/src/github.com/alexellis/href-counter/app ./
CMD ["./app"]  
```

## Stop at a specific build stage

When you build your image, you don't necessarily need to build the entire
Dockerfile including every stage. You can specify a target build stage. The
following command assumes you are using the previous `Dockerfile` but stops at
the stage named `builder`:

```console
$ docker build --target builder -t alexellis2/href-counter:latest .
```

A few scenarios where this might be very powerful are:

- Debugging a specific build stage
- Using a `debug` stage with all debugging symbols or tools enabled, and a
  lean `production` stage
- Using a `testing` stage in which your app gets populated with test data, but
  building for production using a different stage which uses real data

## Use an external image as a "stage"

When using multi-stage builds, you are not limited to copying from stages you
created earlier in your Dockerfile. You can use the `COPY --from` instruction to
copy from a separate image, either using the local image name, a tag available
locally or on a Docker registry, or a tag ID. The Docker client pulls the image
if necessary and copies the artifact from there. The syntax is:

```dockerfile
COPY --from=nginx:latest /etc/nginx/nginx.conf /nginx.conf
```

## Use a previous stage as a new stage

You can pick up where a previous stage left off by referring to it when using
the `FROM` directive. For example:

```dockerfile
# syntax=docker/dockerfile:1

FROM alpine:latest AS builder
RUN apk --no-cache add build-base

FROM builder AS build1
COPY source1.cpp source.cpp
RUN g++ -o /binary source.cpp

FROM builder AS build2
COPY source2.cpp source.cpp
RUN g++ -o /binary source.cpp
```

## Version compatibility

Multi-stage build syntax was introduced in Docker Engine 17.05.

## Differences between legacy build and BuildKit

The legacy Docker Engine builder processes all stages of a Dockerfile leading up to the selected `--target`. It will build a stage even if the selected target doesn't depend on that stage.

BuildKit only builds the stages that the target stage depends on.

For example, given the following Dockerfile:

```dockerfile
# syntax=docker/dockerfile:1
FROM ubuntu AS base
RUN echo "base"

FROM base AS stage1
RUN echo "stage1"

FROM base AS stage2
RUN echo "stage2"
```

With BuildKit enabled, building the `stage2` target in this Dockerfile means only `base` and `stage2` are processed. There is no dependency on `stage1`, so it's skipped.

```console
$ DOCKER_BUILDKIT=1 docker build --no-cache -f Dockerfile --target stage2 .
[+] Building 0.4s (7/7) FINISHED                                                                                                                                                                                
 => [internal] load build definition from Dockerfile                                                                                                                                                       0.0s
 => => transferring dockerfile: 36B                                                                                                                                                                        0.0s
 => [internal] load .dockerignore                                                                                                                                                                          0.0s
 => => transferring context: 2B                                                                                                                                                                            0.0s
 => [internal] load metadata for docker.io/library/ubuntu:latest                                                                                                                                           0.0s
 => CACHED [base 1/2] FROM docker.io/library/ubuntu                                                                                                                                                        0.0s
 => [base 2/2] RUN echo "base"                                                                                                                                                                             0.1s
 => [stage2 1/1] RUN echo "stage2"                                                                                                                                                                         0.2s
 => exporting to image                                                                                                                                                                                     0.0s
 => => exporting layers                                                                                                                                                                                    0.0s
 => => writing image sha256:f55003b607cef37614f607f0728e6fd4d113a4bf7ef12210da338c716f2cfd15                                                                                                               0.0s
```

On the other hand, building the same target without BuildKit results in all stages being processed:

```console
$ DOCKER_BUILDKIT=0 docker build --no-cache -f Dockerfile --target stage2 .
Sending build context to Docker daemon  219.1kB
Step 1/6 : FROM ubuntu AS base
 ---> a7870fd478f4
Step 2/6 : RUN echo "base"
 ---> Running in e850d0e42eca
base
Removing intermediate container e850d0e42eca
 ---> d9f69f23cac8
Step 3/6 : FROM base AS stage1
 ---> d9f69f23cac8
Step 4/6 : RUN echo "stage1"
 ---> Running in 758ba6c1a9a3
stage1
Removing intermediate container 758ba6c1a9a3
 ---> 396baa55b8c3
Step 5/6 : FROM base AS stage2
 ---> d9f69f23cac8
Step 6/6 : RUN echo "stage2"
 ---> Running in bbc025b93175
stage2
Removing intermediate container bbc025b93175
 ---> 09fc3770a9c4
Successfully built 09fc3770a9c4
```

`stage1` gets executed when BuildKit is disabled, even if `stage2` does not depend on it.

BuildKit is enabled by default if you use Docker Desktop.

Always run multi-stage builds with
[BuildKit enabled](/develop/develop-images/build_enhancements/)
for better performance.
