---
title: Configure Docker to use a proxy server
description: How to configure the Docker client to use a proxy server
keywords: network, networking, proxy, client
---

> **Note**
>
> This page describes how to configure the Docker CLI to configure proxies via environment variables in containers.
> For information on configuring Docker Desktop to use HTTP/HTTPS proxies, see [proxies on Mac](../desktop/mac/index.md#proxies) and [proxies on Windows](../desktop/windows/index.md#proxies).
> If you are not running Docker Desktop, and have installed the Docker Engine in
> other ways, refer to the "HTTP/HTTPS proxy" section in
> [configuring the Docker daemon with systemd](../config/daemon/systemd.md#httphttps-proxy).

If your container needs to use an HTTP, HTTPS, or FTP proxy server, you can
configure it in different ways:

- In Docker 17.07 and higher, you can
  [configure the Docker client](#configure-the-docker-client) to pass
  proxy information to containers automatically.

- In Docker 17.06 and earlier versions, you must set the appropriate
  [environment variables](#use-environment-variables)
  within the container. You can do this when you build the image (which makes
  the image less portable) or when you create or run the container.

## Configure the Docker client

1.  On the Docker client, create or edit the file `~/.docker/config.json` in the
    home directory of the user that starts containers. Add JSON similar to the
    following example. Substitute the type of proxy with `httpsProxy` or `ftpProxy` if necessary, and substitute the address and port of the proxy server. You can also configure multiple proxy servers simultaneously.

    You can optionally exclude hosts or ranges from going through the proxy
    server by setting a `noProxy` key to one or more comma-separated IP
    addresses or hosts. Using the `*` character as a wildcard for hosts and using CIDR notation for IP addresses is supported as
    shown in this example:

    ```json
    {
     "proxies":
     {
       "default":
       {
         "httpProxy": "http://192.168.1.12:3128",
         "httpsProxy": "http://192.168.1.12:3128",
         "noProxy": "*.test.example.com,.example2.com,127.0.0.0/8"
       }
     }
    }
    ```

    Save the file.

 2. When you create or start new containers, the environment variables are
    set automatically within the container.

## Use environment variables

### Set the environment variables manually

When you build the image, or using the `--env` flag when you create or run the
container, you can set one or more of the following variables to the appropriate
value. This method makes the image less portable, so if you have Docker 17.07
or higher, you should [configure the Docker client](#configure-the-docker-client)
instead.

| Variable      | Dockerfile example                                | `docker run` example                                |
|:--------------|:--------------------------------------------------|:----------------------------------------------------|
| `HTTP_PROXY`  | `ENV HTTP_PROXY="http://192.168.1.12:3128"`          | `--env HTTP_PROXY="http://192.168.1.12:3128"`          |
| `HTTPS_PROXY` | `ENV HTTPS_PROXY="https://192.168.1.12:3128"`        | `--env HTTPS_PROXY="https://192.168.1.12:3128"`        |
| `FTP_PROXY`   | `ENV FTP_PROXY="ftp://192.168.1.12:3128"`            | `--env FTP_PROXY="ftp://192.168.1.12:3128"`            |
| `NO_PROXY`    | `ENV NO_PROXY="*.test.example.com,.example2.com"` | `--env NO_PROXY="*.test.example.com,.example2.com"` |
