---
title: Image-building best practices
keywords: get started, setup, orientation, quickstart, intro, concepts, containers,
  docker desktop
description: Tips for building images for your application
aliases:
- /get-started/09_image_best/
---

## Image layering

Using the `docker image history` command, you can see the command that was used
to create each layer within an image.

1. Use the `docker image history` command to see the layers in the `getting-started` image you
   created.

    ```console
    $ docker image history getting-started
    ```

    You should get output that looks something like the following.

    ```plaintext
    IMAGE               CREATED             CREATED BY                                      SIZE                COMMENT
    a78a40cbf866        18 seconds ago      /bin/sh -c #(nop)  CMD ["node" "src/index.j…    0B                  
    f1d1808565d6        19 seconds ago      /bin/sh -c yarn install --production            85.4MB              
    a2c054d14948        36 seconds ago      /bin/sh -c #(nop) COPY dir:5dc710ad87c789593…   198kB               
    9577ae713121        37 seconds ago      /bin/sh -c #(nop) WORKDIR /app                  0B                  
    b95baba1cfdb        13 days ago         /bin/sh -c #(nop)  CMD ["node"]                 0B                  
    <missing>           13 days ago         /bin/sh -c #(nop)  ENTRYPOINT ["docker-entry…   0B                  
    <missing>           13 days ago         /bin/sh -c #(nop) COPY file:238737301d473041…   116B                
    <missing>           13 days ago         /bin/sh -c apk add --no-cache --virtual .bui…   5.35MB              
    <missing>           13 days ago         /bin/sh -c #(nop)  ENV YARN_VERSION=1.21.1      0B                  
    <missing>           13 days ago         /bin/sh -c addgroup -g 1000 node     && addu…   74.3MB              
    <missing>           13 days ago         /bin/sh -c #(nop)  ENV NODE_VERSION=12.14.1     0B                  
    <missing>           13 days ago         /bin/sh -c #(nop)  CMD ["/bin/sh"]              0B                  
    <missing>           13 days ago         /bin/sh -c #(nop) ADD file:e69d441d729412d24…   5.59MB   
    ```

    Each of the lines represents a layer in the image. The display here shows the base at the bottom with
    the newest layer at the top. Using this, you can also quickly see the size of each layer, helping 
    diagnose large images.

2. You'll notice that several of the lines are truncated. If you add the `--no-trunc` flag, you'll get the
   full output.

    ```console
    $ docker image history --no-trunc getting-started
    ```

## Layer caching

Now that you've seen the layering in action, there's an important lesson to learn to help decrease build
times for your container images. Once a layer changes, all downstream layers have to be recreated as well.

Look at the following Dockerfile you created for the getting started app.

```dockerfile
# syntax=docker/dockerfile:1
FROM node:18-alpine
WORKDIR /app
COPY . .
RUN yarn install --production
CMD ["node", "src/index.js"]
```

Going back to the image history output, you see that each command in the Dockerfile becomes a new layer in the image.
You might remember that when you made a change to the image, the yarn dependencies had to be reinstalled. It doesn't make much sense to ship around the same dependencies every time you build.

To fix it, you need to restructure your Dockerfile to help support the caching
of the dependencies. For Node-based applications, those dependencies are defined
in the `package.json` file. You can copy only that file in first, install the
dependencies, and then copy in everything else. Then, you only recreate the yarn
dependencies if there was a change to the `package.json`.

1. Update the Dockerfile to copy in the `package.json` first, install dependencies, and then copy everything else in.

   ```dockerfile
   # syntax=docker/dockerfile:1
   FROM node:18-alpine
   WORKDIR /app
   COPY package.json yarn.lock ./
   RUN yarn install --production
   COPY . .
   CMD ["node", "src/index.js"]
   ```

2. Build a new image using `docker build`.

    ```console
    $ docker build -t getting-started .
    ```

    You should see output like the following.

    ```plaintext
    [+] Building 16.1s (10/10) FINISHED
    => [internal] load build definition from Dockerfile
    => => transferring dockerfile: 175B
    => [internal] load .dockerignore
    => => transferring context: 2B
    => [internal] load metadata for docker.io/library/node:18-alpine
    => [internal] load build context
    => => transferring context: 53.37MB
    => [1/5] FROM docker.io/library/node:18-alpine
    => CACHED [2/5] WORKDIR /app
    => [3/5] COPY package.json yarn.lock ./
    => [4/5] RUN yarn install --production
    => [5/5] COPY . .
    => exporting to image
    => => exporting layers
    => => writing image     sha256:d6f819013566c54c50124ed94d5e66c452325327217f4f04399b45f94e37d25
    => => naming to docker.io/library/getting-started
    ```

3. Now, make a change to the `src/static/index.html` file. For example, change the `<title>` to "The Awesome Todo App".

4. Build the Docker image now using `docker build -t getting-started .` again. This time, your output should look a little different.

    ```plaintext
    [+] Building 1.2s (10/10) FINISHED
    => [internal] load build definition from Dockerfile
    => => transferring dockerfile: 37B
    => [internal] load .dockerignore
    => => transferring context: 2B
    => [internal] load metadata for docker.io/library/node:18-alpine
    => [internal] load build context
    => => transferring context: 450.43kB
    => [1/5] FROM docker.io/library/node:18-alpine
    => CACHED [2/5] WORKDIR /app
    => CACHED [3/5] COPY package.json yarn.lock ./
    => CACHED [4/5] RUN yarn install --production
    => [5/5] COPY . .
    => exporting to image
    => => exporting layers
    => => writing image     sha256:91790c87bcb096a83c2bd4eb512bc8b134c757cda0bdee4038187f98148e2eda
    => => naming to docker.io/library/getting-started
    ```

    First off, you should notice that the build was much faster. And, you'll see
    that several steps are using previously cached layers. Pushing and pulling
    this image and updates to it will be much faster as well. 

## Multi-stage builds

Multi-stage builds are an incredibly powerful
tool to help use multiple stages to create an image. There are several advantages for them:

- Separate build-time dependencies from runtime dependencies
- Reduce overall image size by shipping only what your app needs to run

### Maven/Tomcat example

When building Java-based applications, you need a JDK to compile the source code to Java bytecode. However,
that JDK isn't needed in production. Also, you might be using tools like Maven or Gradle to help build the app.
Those also aren't needed in your final image. Multi-stage builds help.

```dockerfile
# syntax=docker/dockerfile:1
FROM maven AS build
WORKDIR /app
COPY . .
RUN mvn package

FROM tomcat
COPY --from=build /app/target/file.war /usr/local/tomcat/webapps 
```

In this example, you use one stage (called `build`) to perform the actual Java build using Maven. In the second
stage (starting at `FROM tomcat`), you copy in files from the `build` stage. The final image is only the last stage
being created, which can be overridden using the `--target` flag.

### React example

When building React applications, you need a Node environment to compile the JS code (typically JSX), SASS stylesheets,
and more into static HTML, JS, and CSS. If you aren't doing server-side rendering, you don't even need a Node environment
for your production build. You can ship the static resources in a static nginx container.

```dockerfile
# syntax=docker/dockerfile:1
FROM node:18 AS build
WORKDIR /app
COPY package* yarn.lock ./
RUN yarn install
COPY public ./public
COPY src ./src
RUN yarn run build

FROM nginx:alpine
COPY --from=build /app/build /usr/share/nginx/html
```

In the previous Dockerfile example, it uses the `node:18` image to perform the build (maximizing layer caching) and then copies the output
into an nginx container.

## Summary

In this section, you learned a few image building best practices, including layer caching and multi-stage builds.

Related information:
 - [Dockerfile reference](/reference/dockerfile/)
 - [Build with Docker guide](../../build/guide/_index.md)
 - [Dockerfile best practices](../../build/building/best-practices.md)

## Next steps

In the next section, you'll learn about additional resources you can use to continue learning about containers.

{{< button text="What next" url="10_what_next.md" >}}
