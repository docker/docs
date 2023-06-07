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
> Special thanks to [Alex Ellis](https://twitter.com/alexellisuk){:target="blank" rel="noopener" class=""}
> for granting permission to use his blog post [Builder pattern vs. Multi-stage builds in Docker](https://blog.alexellis.io/mutli-stage-docker-builds/){:target="blank" rel="noopener" class=""}
> as the basis of the examples on this page.

## Before multi-stage builds

One problem you may face as you build and publish images, is that the size of
those images can sometimes grow quite large. Traditionally, before multi-stage
builds were a thing, keeping the size of images down would require you to
manually clean up resources from the image, so as to keep it small.

In the past, it was common practice to have one Dockerfile for development,
and another, slimmed-down one to use for production.
The development version contained everything needed to build your application.
The production version only contained your application
and the dependencies needed to run it.

To write a truly efficient Dockerfile, you had to come up with shell tricks and
arcane solutions to keep the layers as small as possible.
All to ensure that each layer contained only the artifacts it needed,
and nothing else.

This has been referred to as the _builder pattern_.
The following examples show two Dockerfiles that adhere to this pattern:

- `build.Dockerfile`, for development builds
- `Dockerfile`, for slimmed-down production builds

**`build.Dockerfile`**:

```dockerfile
# syntax=docker/dockerfile:1
FROM golang:1.16
WORKDIR /go/src/github.com/alexellis/href-counter/
COPY app.go ./
RUN go get -d -v golang.org/x/net/html \
  && CGO_ENABLED=0 go build -a -installsuffix cgo -o app .
```

Notice how this example artificially compresses two `RUN` commands together
using the Bash `&&` operator. This is done to avoid creating an additional
layer in the image. Writing Dockerfiles like this is failure-prone and hard to
maintain. It's easy to insert another command and forget to continue the line
using the `\` character, for example.

**`Dockerfile`**:

```dockerfile
# syntax=docker/dockerfile:1
FROM alpine:latest  
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY app ./
CMD ["./app"]
```

The following example is a utility script that:

1. Builds the first image.
2. Creates a container from it to copy the artifact out.
3. Builds the second image.

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

Both images take up room on your system and you still end up with the `app`
artifact on your local disk as well.

Multi-stage builds simplifies this situation!

## Use multi-stage builds

With multi-stage builds, you use multiple `FROM` statements in your Dockerfile.
Each `FROM` instruction can use a different base, and each of them begins a new
stage of the build. You can selectively copy artifacts from one stage to
another, leaving behind everything you don't want in the final image. To show
how this works, you can adapt the `Dockerfile` from the previous section to use
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

You only need the single Dockerfile.
No need for a separate build script.
Just run `docker build`.

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

By default, the stages aren't named, and you refer to them by their integer
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

A few scenarios where this might be useful are:

- Debugging a specific build stage
- Using a `debug` stage with all debugging symbols or tools enabled, and a
  lean `production` stage
- Using a `testing` stage in which your app gets populated with test data, but
  building for production using a different stage which uses real data

## Use an external image as a "stage"

When using multi-stage builds, you aren't limited to copying from stages you
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

## Differences between legacy builder and BuildKit

The legacy Docker Engine builder processes all stages of a Dockerfile leading
up to the selected `--target`. It will build a stage even if the selected
target doesn't depend on that stage.

[BuildKit](../buildkit/index.md) only builds the stages that the target stage
depends on.

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

With [BuildKit enabled](../buildkit/index.md#getting-started), building the
`stage2` target in this Dockerfile means only `base` and `stage2` are processed.
There is no dependency on `stage1`, so it's skipped.

```console
$ DOCKER_BUILDKIT=1 docker build --no-cache -f Dockerfile --target stage2 .
[+] Building 0.4s (7/7) FINISHED                                                                    
 => [internal] load build definition from Dockerfile                                            0.0s
 => => transferring dockerfile: 36B                                                             0.0s
 => [internal] load .dockerignore                                                               0.0s
 => => transferring context: 2B                                                                 0.0s
 => [internal] load metadata for docker.io/library/ubuntu:latest                                0.0s
 => CACHED [base 1/2] FROM docker.io/library/ubuntu                                             0.0s
 => [base 2/2] RUN echo "base"                                                                  0.1s
 => [stage2 1/1] RUN echo "stage2"                                                              0.2s
 => exporting to image                                                                          0.0s
 => => exporting layers                                                                         0.0s
 => => writing image sha256:f55003b607cef37614f607f0728e6fd4d113a4bf7ef12210da338c716f2cfd15    0.0s
```

On the other hand, building the same target without BuildKit results in all
stages being processed:

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

The legacy builder processes `stage1`,
even if `stage2` doesn't depend on it.
