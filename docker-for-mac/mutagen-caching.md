---
description: Mutagen-based caching
keywords: mac, mutagen, volumes
title: High-performance two-way file synchronisation for volume mounts
---

Docker Desktop for Mac on Edge has a new filesharing feature which performs
a continuous two-way sync of files between the host and containers using
[mutagen](https://mutagen.io/). This feature is ideal for app development where
- the source code tree is quite large
- the source is edited on the Mac
- the source is compiled and run interactively inside Linux containers.

The following example shows how the feature should be used.

### Example: a simple react app

This example will bootstrap a simple react app with `npx` and configure
Docker Desktop to synchronise the source code between the host and a
development container.

First create the directory which will contain the app:
```
$ mkdir ~/workspace/my-app
```
Next, enable the two-way sync feature:

- Open the Whale menu, click Preferences, Resources and then File sharing.
- Type in the new directory name in the bottom of the list and hit enter.
- Click on the slider next to the directory name.
- Click on "Apply & Restart".
- Wait for Docker Desktop to restart.

When Docker Desktop has restarted, the Preferences window should look like
this:

![Caching with mutagen is "Ready"](images/mac-mutagen-ready.png)

Run the following command to start a container and bootstrap the app with `npx`:
```
$ docker run -it -v ~/workspace/my-app:/my-app -w /my-app -p 3000:3000 node:lts bash
root@95441305251a:/my-app# npx create-react-app app
root@95441305251a:/my-app# cd app
root@95441305251a:/my-app# npm start
```

Once the development webserver has started, open https://localhost:3000/ in
your browser and observe the app is running.

Return to the Desktop Preferences UI and observe the status text next to the 
cache on/off switch next to the directory name. The status text will be
updated as file changes are detected and then synchronised between the host
and the containers.

Wait until the text says "Ready" and then open the source code in your IDE on
the host. Edit the file `src/App.js`, save the changes and observe the change
in the webserver.

As you edit code on the host, the changes are detected and transferred to the
container for testing. Changes inside the container (e.g. the creation of build
artifacts) are detected and transferred back to the host. 

### Example: avoiding synchronising a sub directory

Although two-way file sync is suitable for many types of files, sometimes we 
know the container is going to generate lots of data which doesn't require
copying to the host (e.g. debug logs).

If your project has a subdirectory whose contents don't need to be continuously
copied back to the host, then use a named docker volume to bypass the sync.

First create a volume using:
```
$ docker volume create donotsyncme
donotsyncme
```
Use the volume for the subdir you want to avoid syncing:
``` 
$ docker run -it -v ~/workspace/my-app:/my-app -v donotsyncme:/my-app/dontsyncme -w /my-app -p 3000:3000 node:lts bash
```

Docker Desktop will synchronise all changes written by the app to `/my-app` to
the host, except changes written to `/my-app/dontsyncme` which will be written
to the named volume instead.

### Best-practices for two-way file synchronisation

For maximum performance:
- Avoid wasting disk space and CPU by minimising the size of the synchronised
  directories. For example, synchronse a project directory like `~/my-app` but
  never sync a large directory like `/Users/` or `/Volumes`. Remember that the
  files will be copied inside the container and therefore must fit within the
  `Docker.qcow2` or `Docker.raw` file.
- For every volume you want to sync in `docker run -v` or in
  `docker-compose.yml`, ensure either the directory itself or a
  parent/grandparent/... directory is listed in Preferences -> Resources ->
  File sharing. Note in particular that if only a *child* directory is listed
  in the Preferences, then the whole `docker run -v` may bypass the two-way
  sync and be slower.
- Avoid changing the same files from both the host and containers. If changes
  are detectedon both sides, the host will "win" and the changes in the
  containers will be discarded.
