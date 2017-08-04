---
title: "Play with Docker"
keywords: get started
description: Play with Docker tests
---
{% include pwd.html %}

## Play with Docker Test 1

Here is a play with Docker terminal.
<div id="term1" class="pwdterm" style="height: 300px; width: 400px; position: relative; right: 20px; top: 75px; z-index:100;"></div>





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

<script>
        pwd.newSession([{selector: '#term1'}]);
        console.log(pwd);
</script>