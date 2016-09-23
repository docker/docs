---
aliases:
- /mackit/networking/
description: Networking
keywords:
- mac, networking
menu:
  main:
    identifier: mac-networking
    parent: pinata_mac_menu
    weight: 4
title: 'Networking '
---

# Networking

Docker for Mac provides several networking features to make it easier to use.

## Features

### VPN Passthrough

Docker for Mac's networking can work when attached to a VPN.
To do this, Docker for Mac intercepts traffic from the `HyperKit` and injects it into OSX as if it originated from the Docker application.

### Port Mapping

When you run a container with the `-p` argument, for example:
```
$ docker run -p 80:80 -d nginx
```
Docker for Mac will make the container port available at `localhost`.

### HTTP/HTTPS Proxy Support

Docker for Mac will detect HTTP/HTTPS Proxy Settings from OSX and automatically propagate these to Docker and to your containers.
For example, if you set your proxy settings to `http://proxy.example.com` in OSX, Docker will use this proxy when pulling containers.

![OSX Proxy Settings](images/proxy-settings.png)

When you start a container, you will see that your proxy settings propagate into the containers. For example:

```
$ docker run -it alpine env
PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin
HOSTNAME=b7edf988b2b5
TERM=xterm
HOME=/root
HTTP_PROXY=http://proxy.example.com:3128
http_proxy=http://proxy.example.com:3128
no_proxy=*.local, 169.254/16
```

You can see from the above output that the `HTTP_PROXY`, `http_proxy` and `no_proxy` environment variables are set.
When your proxy configuration changes, Docker restarts automatically to pick up the new settings.
If you have containers that you wish to keep running across restarts, you should consider using [restart policies](https://docs.docker.com/engine/reference/run/#restart-policies-restart)

## Known Limitations

### There is no docker0 bridge on OSX

Because of the way networking is implemented in Docker for Mac, you cannot see a `docker0` interface in OSX.
This interface is actually within `HyperKit`.

### I can't ping my containers

Unfortunately, due to limitations in OSX, we're unable to route traffic to containers, and from containers back to the host.

<p style="margin-bottom:300px">&nbsp;</p>
