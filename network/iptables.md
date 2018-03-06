---
title: Docker and iptables
description: The basics of how Docker works with iptables
keywords: network, iptables
---

On Linux, Docker manipulates `iptables` rules to provide network isolation.
This is an implementation detail, and you should not modify the rules Docker
inserts into your `iptables` policies.

## Add iptables policies before Docker's rules

All of Docker's `iptables` rules are added to the `DOCKER` chain. Do not
manipulate this table manually. If you need to add rules which load before
Docker's rules, add them to the `DOCKER-USER` chain. These rules are loaded
before any rules Docker creates automatically.

### Restrict connections to the Docker daemon

By default, all external source IPs are allowed to connect to the Docker daemon.
To allow only a specific IP or network to access the containers, insert a
negated rule at the top of the DOCKER filter chain. For example, the following
rule restricts external access to all IP addresses except 192.168.1.1:

```bash
$ iptables -I DOCKER-USER -i ext_if ! -s 192.168.1.1 -j DROP
```

You could instead allow connections from a source subnet. The following rule
only allows access from the subnet 192.168.1.0/24:

```bash
$ iptables -I DOCKER-USER -i ext_if ! -s 192.168.1.0/24 -j DROP
```

Finally, you can specify a range of IP addresses to accept using `--src-range`
(Remember to also add `-m iprange` when using `--src-range` or `--dst-range`):

```bash
$ iptables -I DOCKER-USER -m iprange -i ext_if ! --src-range 192.168.1.1-192.168.1.3 -j DROP
```

You can combine `-s` or `--src-range` with `-d` or `--dst-range` to control both
the source and destination. For instance, if the Docker daemon listens on both
192.168.1.99 and 10.1.2.3, you can make rules specific to `10.1.2.3` and leave
`192.168.1.99` open.

`iptables` is complicated and more complicated rule are out of scope for this
topic. See the [Netfilter.org HOWTO](https://www.netfilter.org/documentation/HOWTO/NAT-HOWTO.html)
for a lot more information.


## Prevent Docker from manipulating iptables

To prevent Docker from manipulating the `iptables` policies at all, set the
`iptables` key to `false` in `/etc/docker/daemon.json`. This is inappropriate
for most users, because the `iptables` policies then need to be managed by hand.

## Next steps

- Read [Docker Reference Architecture: Designing Scalable, Portable Docker Container Networks](https://success.docker.com/Architecture/Docker_Reference_Architecture%3A_Designing_Scalable%2C_Portable_Docker_Container_Networks)
