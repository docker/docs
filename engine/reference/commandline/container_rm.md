---
datafolder: engine-cli
datafile: docker_container_rm
title: docker container rm
---

<!--
Sorry, but the contents of this page are automatically generated from
Docker's source code. If you want to suggest a change to the text that appears
here, you'll need to find the string by searching this repo:

https://www.github.com/docker/docker
-->

{% include cli.md %}

## Examples

### Removing a container using its ID

To remove a container using its ID, find either from a **docker ps -a**
command, or use the ID returned from the **docker run** command, or retrieve
it from a file used to store it using the **docker run --cidfile**:

    $ docker container rm abebf7571666

### Removing a container using the container name

The name of the container can be found using the **docker ps -a**
command. The use that name as follows:

    $ docker container rm hopeful_morse

### Removing a container and all associated volumes

    $ docker container rm -v redis
    redis

This command will remove the container and any volumes associated with it.
Note that if a volume was specified with a name, it will not be removed.

    $ docker create -v awesome:/foo -v /bar --name hello redis

    hello
    
    $ docker container rm -v hello

In this example, the volume for `/foo` will remain in tact, but the volume for
`/bar` will be removed. The same behavior holds for volumes inherited with
`--volumes-from`.
