---
title: "Play with Docker"
keywords: get started
description: Play with Docker tests
---

## Play with Docker Test 1

Here is a play with Docker terminal.

<div id="myTerm" style="width 500px; height: 500px;"></div>
<script src="https://rawgit.com/play-with-docker/sdk/master/dist/pwd.js"></script>
<script>
    pwd.newSession([{selector: '#myTerm'}]);
</script>


## Play with Docker Test 2

```.term1
docker container run hello-world
```
```
Unable to find image 'hello-world:latest' locally
latest: Pulling from library/hello-world
03f4658f8b78: Pull complete
a3ed95caeb02: Pull complete
Digest: sha256:8be990ef2aeb16dbcb9271ddfe2610fa6658d13f6dfb8bc72074cc1ca36966a7
Status: Downloaded newer image for hello-world:latest

Hello from Docker.
This message shows that your installation appears to be working correctly.
```
