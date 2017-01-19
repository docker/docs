---
datafolder: engine-cli
datafile: docker_container_exec
title: docker container exec
---

<!--
Sorry, but the contents of this page are automatically generated from
Docker's source code. If you want to suggest a change to the text that appears
here, you'll need to find the string by searching this repo:

https://www.github.com/docker/docker
-->

{% include cli.md %}

## Examples

    $ docker run --name ubuntu_bash --rm -i -t ubuntu bash

This will create a container named `ubuntu_bash` and start a Bash session.

    $ docker exec -d ubuntu_bash touch /tmp/execWorks

This will create a new file `/tmp/execWorks` inside the running container
`ubuntu_bash`, in the background.

    $ docker exec -it ubuntu_bash bash

This will create a new Bash session in the container `ubuntu_bash`.
