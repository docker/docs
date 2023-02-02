---
description: Controlling and configuring Docker using systemd
keywords: docker, daemon, systemd, configuration
redirect_from:
  - /articles/host_integration/
  - /articles/systemd/
  - /engine/admin/systemd/
  - /engine/articles/systemd/
title: Configure the daemon with systemd
---

This page describes how to customize daemon settings when using systemd.

## Custom Docker daemon options

Most configuration options for the Docker daemon are set using the `daemon.json`
configuration file. See [Docker daemon configuration overview](./index.md) for
more information.

## Manually create the systemd unit files

When installing the binary without a package manager, you may want to integrate
Docker with systemd. For this, install the two unit files (`service` and
`socket`) from
[the github repository](https://github.com/moby/moby/tree/master/contrib/init/systemd)
to `/etc/systemd/system`.

## HTTP/HTTPS proxy

The Docker daemon uses the `HTTP_PROXY`, `HTTPS_PROXY`, and `NO_PROXY`
environmental variables in its start-up environment to configure HTTP or HTTPS
proxy behavior. You can't configure these environment variables using the
`daemon.json` file.

This example overrides the default `docker.service` file.

If you are behind an HTTP or HTTPS proxy server, for example in corporate
settings, you need to add this configuration in the Docker systemd service file.

> **Note for rootless mode**
>
> The location of systemd configuration files are different when running Docker
> in [rootless mode](../../engine/security/rootless.md). When running in
> rootless mode, Docker is started as a user-mode systemd service, and uses
> files stored in each users' home directory in
> `~/.config/systemd/user/docker.service.d/`. In addition, `systemctl` must be
> executed without `sudo` and with the `--user` flag. Select the _"rootless
> mode"_ tab below if you are running Docker in rootless mode.

<ul class="nav nav-tabs">
  <li class="active"><a data-toggle="tab" data-target="#rootful">regular install</a></li>
  <li><a data-toggle="tab" data-target="#rootless">rootless mode</a></li>
</ul>
<div class="tab-content">
<div id="rootful" class="tab-pane fade in active" markdown="1">

1. Create a systemd drop-in directory for the `docker` service:

   ```console
   $ sudo mkdir -p /etc/systemd/system/docker.service.d
   ```

2. Create a file named `/etc/systemd/system/docker.service.d/http-proxy.conf`
   that adds the `HTTP_PROXY` environment variable:

   ```systemd
   [Service]
   Environment="HTTP_PROXY=http://proxy.example.com:80"
   ```

   If you are behind an HTTPS proxy server, set the `HTTPS_PROXY` environment
   variable:

   ```systemd
   [Service]
   Environment="HTTPS_PROXY=https://proxy.example.com:443"
   ```

   Multiple environment variables can be set; to set both a non-HTTPS and a
   HTTPs proxy;

   ```systemd
   [Service]
   Environment="HTTP_PROXY=http://proxy.example.com:80"
   Environment="HTTPS_PROXY=https://proxy.example.com:443"
   ```

   > **Note**
   >
   > Special characters in the proxy value, such as `#?!()[]{}`, must be double
   > escaped using `%%`. For example:
   >
   > ```
   > [Service]
   > Environment="HTTP_PROXY=http://domain%%5Cuser:complex%%23pass@proxy.example.com:8080/"
   > ```

3. If you have internal Docker registries that you need to contact without
   proxying, you can specify them via the `NO_PROXY` environment variable.

   The `NO_PROXY` variable specifies a string that contains comma-separated
   values for hosts that should be excluded from proxying. These are the options
   you can specify to exclude hosts:

   - IP address prefix (`1.2.3.4`)
   - Domain name, or a special DNS label (`*`)
   - A domain name matches that name and all subdomains. A domain name with a
     leading "." matches subdomains only. For example, given the domains
     `foo.example.com` and `example.com`:
     - `example.com` matches `example.com` and `foo.example.com`, and
     - `.example.com` matches only `foo.example.com`
   - A single asterisk (`*`) indicates that no proxying should be done
   - Literal port numbers are accepted by IP address prefixes (`1.2.3.4:80`) and
     domain names (`foo.example.com:80`)

   Config example:

   ```systemd
   [Service]
   Environment="HTTP_PROXY=http://proxy.example.com:80"
   Environment="HTTPS_PROXY=https://proxy.example.com:443"
   Environment="NO_PROXY=localhost,127.0.0.1,docker-registry.example.com,.corp"
   ```

4. Flush changes and restart Docker

   ```console
   $ sudo systemctl daemon-reload
   $ sudo systemctl restart docker
   ```

5. Verify that the configuration has been loaded and matches the changes you
   made, for example:

   ```console
   $ sudo systemctl show --property=Environment docker

   Environment=HTTP_PROXY=http://proxy.example.com:80 HTTPS_PROXY=https://proxy.example.com:443 NO_PROXY=localhost,127.0.0.1,docker-registry.example.com,.corp
   ```

</div>
<div id="rootless" class="tab-pane fade in" markdown="1">

1. Create a systemd drop-in directory for the `docker` service:

   ```console
   $ mkdir -p ~/.config/systemd/user/docker.service.d
   ```

2. Create a file named `~/.config/systemd/user/docker.service.d/http-proxy.conf`
   that adds the `HTTP_PROXY` environment variable:

   ```systemd
   [Service]
   Environment="HTTP_PROXY=http://proxy.example.com:80"
   ```

   If you are behind an HTTPS proxy server, set the `HTTPS_PROXY` environment
   variable:

   ```systemd
   [Service]
   Environment="HTTPS_PROXY=https://proxy.example.com:443"
   ```

   Multiple environment variables can be set; to set both a non-HTTPS and a
   HTTPs proxy;

   ```systemd
   [Service]
   Environment="HTTP_PROXY=http://proxy.example.com:80"
   Environment="HTTPS_PROXY=https://proxy.example.com:443"
   ```

   > **Note**
   >
   > Special characters in the proxy value, such as `#?!()[]{}`, must be double
   > escaped using `%%`. For example:
   >
   > ```
   > [Service]
   > Environment="HTTP_PROXY=http://domain%%5Cuser:complex%%23pass@proxy.example.com:8080/"
   > ```

3. If you have internal Docker registries that you need to contact without
   proxying, you can specify them via the `NO_PROXY` environment variable.

   The `NO_PROXY` variable specifies a string that contains comma-separated
   values for hosts that should be excluded from proxying. These are the options
   you can specify to exclude hosts:

   - IP address prefix (`1.2.3.4`)
   - Domain name, or a special DNS label (`*`)
   - A domain name matches that name and all subdomains. A domain name with a
     leading "." matches subdomains only. For example, given the domains
     `foo.example.com` and `example.com`:
     - `example.com` matches `example.com` and `foo.example.com`, and
     - `.example.com` matches only `foo.example.com`
   - A single asterisk (`*`) indicates that no proxying should be done
   - Literal port numbers are accepted by IP address prefixes (`1.2.3.4:80`) and
     domain names (`foo.example.com:80`)

   Config example:

   ```systemd
   [Service]
   Environment="HTTP_PROXY=http://proxy.example.com:80"
   Environment="HTTPS_PROXY=https://proxy.example.com:443"
   Environment="NO_PROXY=localhost,127.0.0.1,docker-registry.example.com,.corp"
   ```

4. Flush changes and restart Docker

   ```console
   $ systemctl --user daemon-reload
   $ systemctl --user restart docker
   ```

5. Verify that the configuration has been loaded and matches the changes you
   made, for example:

   ```console
   $ systemctl --user show --property=Environment docker

   Environment=HTTP_PROXY=http://proxy.example.com:80 HTTPS_PROXY=https://proxy.example.com:443 NO_PROXY=localhost,127.0.0.1,docker-registry.example.com,.corp
   ```

</div>
</div> <!-- tab-content -->
