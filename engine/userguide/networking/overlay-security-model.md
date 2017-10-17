---
description: Docker swarm mode overlay network security model
keywords: network, docker, documentation, user guide, multihost, swarm mode, overlay
title: Docker swarm mode overlay network security model
---

Overlay networking for Docker Engine swarm mode comes secure out of the box. The
swarm nodes exchange overlay network information using a gossip protocol. By
default the nodes encrypt and authenticate information they exchange via gossip
using the [AES algorithm](https://en.wikipedia.org/wiki/Galois/Counter_Mode) in
GCM mode. Manager nodes in the swarm rotate the key used to encrypt gossip data
every 12 hours.

You can also encrypt data exchanged between containers on different nodes on the
overlay network. To enable encryption, when you create an overlay network pass
the `--opt encrypted` flag:

```bash
$ docker network create --opt encrypted --driver overlay my-multi-host-network

dt0zvqn0saezzinc8a5g4worx
```

When you enable overlay encryption, Docker creates IPSEC tunnels between all the
nodes where tasks are scheduled for services attached to the overlay network.
These tunnels also use the AES algorithm in GCM mode and manager nodes
automatically rotate the keys every 12 hours.

> **Do not attach Windows nodes to encrypted overlay networks.**
>
> Overlay network encryption is not supported on Windows. If a Windows node
> attempts to connect to an encrypted overlay network, no error is detected but
> the node will not be able to communicate.
{: .warning }

## Swarm mode overlay networks and unmanaged containers

It is possible to use the overlay network feature with both `--opt encrypted --attachable`, and attach unmanaged containers to that network:

```bash
$ docker network create --opt encrypted --driver overlay --attachable my-attachable-multi-host-network

9s1p1sfaqtvaibq6yp7e6jsrt
```

Just like services that are attached to an encrypted network, regular containers can also have the advantage of encrypted traffic when attached to a network created this way.
