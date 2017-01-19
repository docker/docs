---
datafolder: engine-cli
datafile: docker_container_wait
title: docker container wait
---

<!--
Sorry, but the contents of this page are automatically generated from
Docker's source code. If you want to suggest a change to the text that appears
here, you'll need to find the string by searching this repo:

https://www.github.com/docker/docker
-->

{% include cli.md %}

## Examples

# EXAMPLES

```bash
$ docker run -d fedora sleep 99

079b83f558a2bc52ecad6b2a5de13622d584e6bb1aea058c11b36511e85e7622

$ docker container wait 079b83f558a2bc

0
```
