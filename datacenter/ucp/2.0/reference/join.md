+++
title = "join"
keywords= ["join, ucp"]
description = "Joins a node to an existing Docker Universal Control Plane cluster."
[menu.main]
parent = "ucp_ref"
identifier = "ucp_ref_join"
+++

# docker/ucp join

Join this engine to an existing UCP

## Usage

```
docker run --rm -it \
           --name ucp \
           -v /var/run/docker.sock:/var/run/docker.sock \
           docker/ucp \
           join [command options]
```

## Description

The join command is no longer used.  To join a node to UCP, simply run `docker swarm join ...`
