---
datafolder: engine-cli
datafile: docker_container_logs
title: docker container logs
---

<!--
Sorry, but the contents of this page are automatically generated from
Docker's source code. If you want to suggest a change to the text that appears
here, you'll need to find the string by searching this repo:

https://www.github.com/docker/docker
-->

{% include cli.md %}

## Examples

### Display all containers, including non-running

```bash
$ docker container ls -a

CONTAINER ID        IMAGE                 COMMAND                CREATED             STATUS      PORTS    NAMES
a87ecb4f327c        fedora:20             /bin/sh -c #(nop) MA   20 minutes ago      Exit 0               desperate_brattain
01946d9d34d8        vpavlin/rhel7:latest  /bin/sh -c #(nop) MA   33 minutes ago      Exit 0               thirsty_bell
c1d3b0166030        acffc0358b9e          /bin/sh -c yum -y up   2 weeks ago         Exit 1               determined_torvalds
41d50ecd2f57        fedora:20             /bin/sh -c #(nop) MA   2 weeks ago         Exit 0               drunk_pike
```

### Display only IDs of all containers, including non-running

```bash
$ docker container ls -a -q

a87ecb4f327c
01946d9d34d8
c1d3b0166030
41d50ecd2f57
```

# Display only IDs of all containers that have the name `determined_torvalds`

```bash
$ docker container ls -a -q --filter=name=determined_torvalds

c1d3b0166030
```

### Display containers with their commands

```bash
{% raw %}
$ docker container ls --format "{{.ID}}: {{.Command}}"

a87ecb4f327c: /bin/sh -c #(nop) MA
01946d9d34d8: /bin/sh -c #(nop) MA
c1d3b0166030: /bin/sh -c yum -y up
41d50ecd2f57: /bin/sh -c #(nop) MA
{% endraw %}
```

### Display containers with their labels in a table

```bash
{% raw %}
$ docker container ls --format "table {{.ID}}\t{{.Labels}}"

CONTAINER ID        LABELS
a87ecb4f327c        com.docker.swarm.node=ubuntu,com.docker.swarm.storage=ssd
01946d9d34d8
c1d3b0166030        com.docker.swarm.node=debian,com.docker.swarm.cpu=6
41d50ecd2f57        com.docker.swarm.node=fedora,com.docker.swarm.cpu=3,com.docker.swarm.storage=ssd
{% endraw %}
```

### Display containers with their node label in a table

```bash
{% raw %}
$ docker container ls --format 'table {{.ID}}\t{{(.Label "com.docker.swarm.node")}}'

CONTAINER ID        NODE
a87ecb4f327c        ubuntu
01946d9d34d8
c1d3b0166030        debian
41d50ecd2f57        fedora
{% endraw %}
```

### Display containers with `remote-volume` mounted

```bash
{% raw %}
$ docker container ls --filter volume=remote-volume --format "table {{.ID}}\t{{.Mounts}}"

CONTAINER ID        MOUNTS
9c3527ed70ce        remote-volume
{% endraw %}
```

### Display containers with a volume mounted in `/data`

```bash
{% raw %}
$ docker container ls --filter volume=/data --format "table {{.ID}}\t{{.Mounts}}"

CONTAINER ID        MOUNTS
9c3527ed70ce        remote-volume
{% endraw %}
```
