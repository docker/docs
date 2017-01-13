---
datafolder: engine-cli
datafile: docker_container_export
title: docker container export
---

<!--
Sorry, but the contents of this page are automatically generated from
Docker's source code. If you want to suggest a change to the text that appears
here, you'll need to find the string by searching this repo:

https://www.github.com/docker/docker
-->

{% include cli.md %}

## Examples

Export the contents of the container called angry_bell to a tar file
called angry_bell.tar:

```bash
$ docker export angry_bell > angry_bell.tar

$ docker export --output=angry_bell-latest.tar angry_bell

$ ls -sh angry_bell.tar

321M angry_bell.tar

$ ls -sh angry_bell-latest.tar

321M angry_bell-latest.tar
```
