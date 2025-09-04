---
title: Using the build cache
keywords: concepts, build, images, container, docker desktop
description: This concept page will teach you about the build cache, what changes invalidate the cache and how to effectively use the build cache.
summary: |
  Using the build cache effectively allows you to achieve faster builds by
  reusing results from previous builds and skipping unnecessary steps. To
  maximize cache usage and avoid resource-intensive and time-consuming
  rebuilds, it's crucial to understand how cache invalidation works. In this
  guide, you’ll learn how to use the Docker build cache efficiently for
  streamlined Docker image development and continuous integration workflows.
weight: 4
aliases: 
 - /guides/docker-concepts/building-images/using-the-build-cache/
---

{{< youtube-embed Ri6jMknjprY >}}

## Explanation

Consider the following Dockerfile that you created for the [getting-started](./writing-a-dockerfile/) app.


```dockerfile
FROM node:22-alpine
WORKDIR /app
COPY . .
RUN yarn install --production
CMD ["node", "./src/index.js"]
```

When you run the `docker build` command to create a new image, Docker executes each instruction in your Dockerfile, creating a layer for each command and in the order specified. For each instruction, Docker checks whether it can reuse the instruction from a previous build. If it finds that you've already executed a similar instruction before, Docker doesn't need to redo it. Instead, it’ll use the cached result. This way, your build process becomes faster and more efficient, saving you valuable time and resources.

Using the build cache effectively lets you achieve faster builds by reusing results from previous builds and skipping unnecessary work.
In order to maximize cache usage and avoid resource-intensive and time-consuming rebuilds, it's important to understand how cache invalidation works.
Here are a few examples of situations that can cause cache to be invalidated:

- Any changes to the command of a `RUN` instruction invalidates that layer. Docker detects the change and invalidates the build cache if there's any modification to a `RUN` command in your Dockerfile.

- Any changes to files copied into the image with the `COPY` or `ADD` instructions. Docker keeps an eye on any alterations to files within your project directory. Whether it's a change in content or properties like permissions, Docker considers these modifications as triggers to invalidate the cache.

- Once one layer is invalidated, all following layers are also invalidated. If any previous layer, including the base image or intermediary layers, has been invalidated due to changes, Docker ensures that subsequent layers relying on it are also invalidated. This keeps the build process synchronized and prevents inconsistencies.

When you're writing or editing a Dockerfile, keep an eye out for unnecessary cache misses to ensure that builds run as fast and efficiently as possible.

## Try it out

In this hands-on guide, you will learn how to use the Docker build cache effectively for a Node.js application.

### Build the application

1. [Download and install](https://www.docker.com/products/docker-desktop/) Docker Desktop.

2. Open a terminal and [clone this sample application](https://github.com/dockersamples/todo-list-app).


    ```console
    $ git clone https://github.com/dockersamples/todo-list-app
    ```

3. Navigate into the `todo-list-app` directory:


    ```console
    $ cd todo-list-app
    ```

    Inside this directory, you'll find a file named `Dockerfile` with the following content:


    ```dockerfile
    FROM node:22-alpine
    WORKDIR /app
    COPY . .
    RUN yarn install --production
    EXPOSE 3000
    CMD ["node", "./src/index.js"]
    ```

4. Execute the following command to build the Docker image:

    ```console
    $ docker build .
    ```

    Here’s the result of the build process:

    ```console
    [+] Building 20.0s (10/10) FINISHED
    ```

    The first line indicates that the entire build process took *20.0 seconds*. The first build may take some time as it installs dependencies.

5. Rebuild without making changes.

   Now, re-run the `docker build` command without making any change in the source code or Dockerfile as shown:

    ```console
    $ docker build .
    ```

   Subsequent builds after the initial are faster due to the caching mechanism, as long as the commands and context remain unchanged. Docker caches the intermediate layers generated during the build process. When you rebuild the image without making any changes to the Dockerfile or the source code, Docker can reuse the cached layers, significantly speeding up the build process.

    ```console
    [+] Building 1.0s (9/9) FINISHED                                                                            docker:desktop-linux
     => [internal] load build definition from Dockerfile                                                                        0.0s
     => => transferring dockerfile: 187B                                                                                        0.0s
     ...
     => [internal] load build context                                                                                           0.0s
     => => transferring context: 8.16kB                                                                                         0.0s
     => CACHED [2/4] WORKDIR /app                                                                                               0.0s
     => CACHED [3/4] COPY . .                                                                                                   0.0s
     => CACHED [4/4] RUN yarn install --production                                                                              0.0s
     => exporting to image                                                                                                      0.0s
     => => exporting layers                                                                                                     0.0s
     => => exporting manifest
   ```


   The subsequent build was completed in just 1.0 second by leveraging the cached layers. No need to repeat time-consuming steps like installing dependencies.


    <table>
      <thead>
        <tr>
          <th>Steps
          </th>
          <th>Description
          </th>
          <th>Time Taken (1st Run)
          </th>
          <th>Time Taken (2nd Run)
          </th>
        </tr>
      </thead>
      <tbody>
      <tr>
       <td>1
       </td>
       <td><code>Load build definition from Dockerfile</code>
       </td>
       <td>0.0 seconds
       </td>
       <td>0.0 seconds
       </td>
      </tr>
      <tr>
       <td>2
       </td>
       <td><code>Load metadata for docker.io/library/node:22-alpine</code>
       </td>
       <td>2.7 seconds
       </td>
       <td>0.9 seconds
       </td>
      </tr>
      <tr>
       <td>3
       </td>
       <td><code>Load .dockerignore</code>
       </td>
       <td>0.0 seconds
       </td>
       <td>0.0 seconds
       </td>
      </tr>
      <tr>
       <td>4
       </td>
       <td><code>Load build context</code>
    <p>
    (Context size: 4.60MB)
       </td>
       <td>0.1 seconds
       </td>
       <td>0.0 seconds
       </td>
      </tr>
      <tr>
       <td>5
       </td>
       <td><code>Set the working directory (WORKDIR)</code>
       </td>
       <td>0.1 seconds
       </td>
       <td>0.0 seconds
       </td>
      </tr>
      <tr>
       <td>6
       </td>
       <td><code>Copy the local code into the container</code>
       </td>
       <td>0.0 seconds
       </td>
       <td>0.0 seconds
       </td>
      </tr>
      <tr>
       <td>7
       </td>
       <td><code>Run yarn install --production</code>
       </td>
       <td>10.0 seconds
       </td>
       <td>0.0 seconds
       </td>
      </tr>
      <tr>
       <td>8
       </td>
       <td><code>Exporting layers</code>
       </td>
       <td>2.2 seconds
       </td>
       <td>0.0 seconds
       </td>
      </tr>
      <tr>
       <td>9
       </td>
       <td><code>Exporting the final image</code>
       </td>
       <td>3.0 seconds
       </td>
       <td>0.0 seconds
       </td>
     </tr>
     </tbody>
    </table>


    Going back to the `docker image history` output, you see that each command in the Dockerfile becomes a new layer in the image. You might remember that when you made a change to the image, the `yarn` dependencies had to be reinstalled. Is there a way to fix this? It doesn't make much sense to reinstall the same dependencies every time you build, right?

    To fix this, restructure your Dockerfile so that the dependency cache remains valid unless it really needs to be invalidated. For Node-based applications, dependencies are defined in the `package.json` file. You'll want to reinstall the dependencies if that file changes, but use cached dependencies if the file is unchanged. So, start by copying only that file first, then install the dependencies, and finally copy everything else. Then, you only need to recreate the yarn dependencies if there was a change to the `package.json` file.

6. Update the Dockerfile to copy in the `package.json` file first, install dependencies, and then copy everything else in.

     ```dockerfile
     FROM node:22-alpine
     WORKDIR /app
     COPY package.json yarn.lock ./
     RUN yarn install --production 
     COPY . . 
     EXPOSE 3000
     CMD ["node", "src/index.js"]
     ```

7. Create a file named `.dockerignore` in the same folder as the Dockerfile with the following contents.

     ```plaintext
     node_modules
     ```

8. Build the new image:

    ```console
    $ docker build .
    ```

    You'll then see output similar to the following:

    ```console
    [+] Building 16.1s (10/10) FINISHED
    => [internal] load build definition from Dockerfile                                               0.0s
    => => transferring dockerfile: 175B                                                               0.0s
    => [internal] load .dockerignore                                                                  0.0s
    => => transferring context: 2B                                                                    0.0s
    => [internal] load metadata for docker.io/library/node:22-alpine                                  0.0s
    => [internal] load build context                                                                  0.8s
    => => transferring context: 53.37MB                                                               0.8s
    => [1/5] FROM docker.io/library/node:22-alpine                                                    0.0s
    => CACHED [2/5] WORKDIR /app                                                                      0.0s
    => [3/5] COPY package.json yarn.lock ./                                                           0.2s
    => [4/5] RUN yarn install --production                                                           14.0s
    => [5/5] COPY . .                                                                                 0.5s
    => exporting to image                                                                             0.6s
    => => exporting layers                                                                            0.6s
    => => writing image     
    sha256:d6f819013566c54c50124ed94d5e66c452325327217f4f04399b45f94e37d25        0.0s
    => => naming to docker.io/library/node-app:2.0                                                 0.0s
    ```

    You'll see that all layers were rebuilt. Perfectly fine since you changed the Dockerfile quite a bit.

9. Now, make a change to the `src/static/index.html` file (like change the title to say "The Awesome Todo App").

10. Build the Docker image. This time, your output should look a little different.

    ```console
    $ docker build -t node-app:3.0 .
    ```

    You'll then see output similar to the following:

    ```console
    [+] Building 1.2s (10/10) FINISHED 
    => [internal] load build definition from Dockerfile                                               0.0s
    => => transferring dockerfile: 37B                                                                0.0s
    => [internal] load .dockerignore                                                                  0.0s
    => => transferring context: 2B                                                                    0.0s
    => [internal] load metadata for docker.io/library/node:22-alpine                                  0.0s 
    => [internal] load build context                                                                  0.2s
    => => transferring context: 450.43kB                                                              0.2s
    => [1/5] FROM docker.io/library/node:22-alpine                                                    0.0s
    => CACHED [2/5] WORKDIR /app                                                                      0.0s
    => CACHED [3/5] COPY package.json yarn.lock ./                                                    0.0s
    => CACHED [4/5] RUN yarn install --production                                                     0.0s
    => [5/5] COPY . .                                                                                 0.5s 
    => exporting to image                                                                             0.3s
    => => exporting layers                                                                            0.3s
    => => writing image     
    sha256:91790c87bcb096a83c2bd4eb512bc8b134c757cda0bdee4038187f98148e2eda       0.0s
    => => naming to docker.io/library/node-app:3.0                                                 0.0s
    ```

    First off, you should notice that the build was much faster. You'll see that several steps are using previously cached layers. That's good news; you're using the build cache. Pushing and pulling this image and updates to it will be much faster as well.

By following these optimization techniques, you can make your Docker builds faster and more efficient, leading to quicker iteration cycles and improved development productivity.

## Additional resources

* [Optimizing builds with cache management](/build/cache/)
* [Cache Storage Backend](/build/cache/backends/)
* [Build cache invalidation](/build/cache/invalidation/)


## Next steps

Now that you understand how to use the Docker build cache effectively, you're ready to learn about Multi-stage builds.

{{< button text="Multi-stage builds" url="multi-stage-builds" >}}
