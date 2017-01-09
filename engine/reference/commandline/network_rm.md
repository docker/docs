---
datafolder: engine-cli
datafile: docker_network_rm
title: docker network rm
---
<!--
Sorry, but the contents of this page are automatically generated from
Docker's source code. If you want to suggest a change to the text that appears
here, you'll need to find the string by searching this repo:

https://www.github.com/docker/docker
-->
{% include cli.md %}

## Examples

```bash
  $ docker network rm my-network
```

To delete multiple networks in a single `docker network rm` command, provide
multiple network names or ids. The following example deletes a network with id
`3695c422697f` and a network named `my-network`:

```bash
  $ docker network rm 3695c422697f my-network
```

When you specify multiple networks, the command attempts to delete each in turn.
If the deletion of one network fails, the command continues to the next on the
list and tries to delete that. The command reports success or failure for each
deletion.
