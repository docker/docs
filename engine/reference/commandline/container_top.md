---
datafolder: engine-cli
datafile: docker_container_top
title: docker container top
---

<!--
Sorry, but the contents of this page are automatically generated from
Docker's source code. If you want to suggest a change to the text that appears
here, you'll need to find the string by searching this repo:

https://www.github.com/docker/docker
-->

{% include cli.md %}

## Examples

Run **docker container top** with the ps option of -x:

    $ docker container top 8601afda2b -x
    PID      TTY       STAT       TIME         COMMAND
    16623    ?         Ss         0:00         sleep 99999
