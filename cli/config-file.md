---
title: Docker client configuration file
description: |
  Learn about the configuration file that you can use to configure your Docker
  client, and the Docker CLI
keywords: cli, configuration, config.json
---

By default, the Docker command line stores its configuration files in a
directory called `.docker` within your `$HOME` directory.

Docker manages most of the files in the configuration directory
and you shouldn't modify them. However, you can modify the
`config.json` file to control certain aspects of how the `docker`
command behaves.

You can modify the `docker` command behavior using environment
variables or command-line options. You can also use options within
`config.json` to modify some of the same behavior. If an environment variable
and the `--config` flag are set, the flag takes precedent over the environment
variable. Command line options override environment variables and environment
variables override properties you specify in a `config.json` file.

### Change the `.docker` directory

To specify a different directory, use the `DOCKER_CONFIG`
environment variable or the `--config` command line option. The `--config`
option takes precedence over the `DOCKER_CONFIG` environment variable if you
specify both. The following example overrides the `docker ps` command using a
`config.json` file located in the `~/testconfigs/` directory.

```console
$ docker --config ~/testconfigs/ ps
```

This flag only applies to whatever command is being ran. For persistent
configuration, you can set the `DOCKER_CONFIG` environment variable in your
shell (e.g. `~/.profile` or `~/.bashrc`). The following example sets the new
directory to be `HOME/newdir/.docker`.

```console
$ echo export DOCKER_CONFIG=$HOME/newdir/.docker > ~/.profile
```

## Docker CLI configuration file (`config.json`) properties

Use the Docker CLI configuration to customize settings for the `docker` CLI. The
configuration file uses JSON formatting, and properties:

The default location of the configuration file is `~/.docker/config.json`.
Refer to the
[change the `.docker` directory](#change-the-docker-directory) section to use a
different location.

> **Warning**
> 
> The configuration file and other files inside the `~/.docker` configuration
> directory may contain sensitive information, such as authentication information
> for proxies or, depending on your credential store, credentials for your image
> registries. Review your configuration file's content before sharing with others,
> and prevent committing the file to version control.

### Customize the default output format for commands

These fields allow you to customize the default output format for some commands
if you don't use the `--format` flag.

| Property               | Description                                                                                                                                                                                                              |
|:-----------------------|:-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `configFormat`         | Custom default format for `docker config ls` output. Refer to the [format the output section in the `docker config ls` documentation](../engine/reference/commandline/config_ls.md#format) for a list of supported formatting directives.            |
| `imagesFormat`         | Custom default format for `docker images` / `docker image ls` output. Refer to the [format the output section in the `docker images` documentation](../engine/reference/commandline/images.md#format) for a list of supported formatting directives. |
| `nodesFormat`          | Custom default format for `docker node ls` output. Refer to the [formatting section in the `docker node ls` documentation](../engine/reference/commandline/node_ls.md#format) for a list of supported formatting directives.                         |
| `pluginsFormat`        | Custom default format for `docker plugin ls` output. Refer to the [formatting section in the `docker plugin ls` documentation](../engine/reference/commandline/plugin_ls.md#format) for a list of supported formatting directives.                   |
| `psFormat`             | Custom default format for `docker ps` / `docker container ps` output. Refer to the [formatting section in the `docker ps` documentation](../engine/reference/commandline/ps.md#format) for a list of supported formatting directives.                |
| `secretFormat`         | Custom default format for `docker secret ls` output. Refer to the [format the output section in the `docker secret ls` documentation](../engine/reference/commandline/secret_ls.md#format) for a list of supported formatting directives.            |
| `serviceInspectFormat` | Custom default format for `docker service inspect` output. Refer to the [formatting section in the `docker service inspect` documentation](../engine/reference/commandline/service_inspect.md#format) for a list of supported formatting directives. |
| `servicesFormat`       | Custom default format for `docker service ls` output. Refer to the [formatting section in the `docker service ls` documentation](../engine/reference/commandline/service_ls.md#format) for a list of supported formatting directives.                |
| `statsFormat`          | Custom default format for `docker stats` output. Refer to the [formatting section in the `docker stats` documentation](../engine/reference/commandline/stats.md#format) for a list of supported formatting directives.                               |

### Custom HTTP headers

The property `HttpHeaders` specifies a set of headers to include in all messages
sent from the Docker client to the daemon. Docker doesn't try to interpret or
understand these headers; it simply puts them into the messages. Docker does
not allow these headers to change any headers it sets for itself.

### Credential store options

The property `credsStore` specifies an external binary to serve as the default
credential store. When this property is set, `docker login` will attempt to
store credentials in the binary specified by `docker-credential-<value>` which
is visible on `$PATH`. If this property isn't set, credentials will be stored
in the `auths` property of the config. For more information, see the
[Credentials store section in the `docker login` documentation](../engine/reference/commandline/login.md#credentials-store)

The property `credHelpers` specifies a set of credential helpers to use
preferentially over `credsStore` or `auths` when storing and retrieving
credentials for specific registries. If this property is set, the binary
`docker-credential-<value>` will be used when storing or retrieving credentials
for a specific registry. For more information, see the
[Credential helpers section in the `docker login` documentation](../engine/reference/commandline/login.md#credential-helpers)

### Automatic proxy configuration for containers

The property `proxies` specifies proxy environment variables to be automatically
set on containers, and set as `--build-arg` on containers used during `docker build`.
A `"default"` set of proxies can be configured, and will be used for any Docker
daemon that the client connects to, or a configuration per host (Docker daemon),
for example, `https://docker-daemon1.example.com`. The following properties can
be set for each environment:

| Property       | Description                                                                                             |
|:---------------|:--------------------------------------------------------------------------------------------------------|
| `httpProxy`    | Default value of `HTTP_PROXY` and `http_proxy` for containers, and as `--build-arg` on `docker build`   |
| `httpsProxy`   | Default value of `HTTPS_PROXY` and `https_proxy` for containers, and as `--build-arg` on `docker build` |
| `ftpProxy`     | Default value of `FTP_PROXY` and `ftp_proxy` for containers, and as `--build-arg` on `docker build`     |
| `noProxy`      | Default value of `NO_PROXY` and `no_proxy` for containers, and as `--build-arg` on `docker build`       |
| `allProxy`     | Default value of `ALL_PROXY` and `all_proxy` for containers, and as `--build-arg` on `docker build`     |

These settings are used to configure proxy settings for containers only, and not
used as proxy settings for the `docker` CLI or the `dockerd` daemon. Refer to the
[environment variables](./env-vars.md) and [HTTP/HTTPS proxy](../config/daemon/systemd.md#httphttps-proxy)
sections for configuring proxy settings for the daemon.

> **Warning**
> 
> Proxy settings may contain sensitive information (for example, if the proxy
> requires authentication). Environment variables are stored as plain text in
> the container's configuration, and as such can be inspected through the remote
> API or committed to an image when using `docker commit`.

### Default key-sequence to detach from containers

Once attached to a container, users detach from it and leave it running using
the using `CTRL-p CTRL-q` key sequence. This detach key sequence is customizable
using the `detachKeys` property. Specify a `<sequence>` value for the
property. The format of the `<sequence>` is a comma-separated list of either
a letter [a-Z], or the `ctrl-` combined with any of the following:

* `a-z` (a single lowercase alpha character )
* `@` (at sign)
* `[` (left bracket)
* `\\` (two backward slashes)
*  `_` (underscore)
* `^` (caret)

Your customization applies to all containers started in with your Docker client.
Users can override your custom or the default key sequence on a per-container
basis. To do this, the user specifies the `--detach-keys` flag with the `docker
attach`, `docker exec`, `docker run` or `docker start` command.

### CLI Plugin options

The property `plugins` contains settings specific to CLI plugins. The
key is the plugin name, while the value is a further map of options,
which are specific to that plugin.


### Sample configuration file

Following is a sample `config.json` file to illustrate the format used for
various fields:

```json
{% raw %}
{
  "HttpHeaders": {
    "MyHeader": "MyValue"
  },
  "psFormat": "table {{.ID}}\\t{{.Image}}\\t{{.Command}}\\t{{.Labels}}",
  "imagesFormat": "table {{.ID}}\\t{{.Repository}}\\t{{.Tag}}\\t{{.CreatedAt}}",
  "pluginsFormat": "table {{.ID}}\t{{.Name}}\t{{.Enabled}}",
  "statsFormat": "table {{.Container}}\t{{.CPUPerc}}\t{{.MemUsage}}",
  "servicesFormat": "table {{.ID}}\t{{.Name}}\t{{.Mode}}",
  "secretFormat": "table {{.ID}}\t{{.Name}}\t{{.CreatedAt}}\t{{.UpdatedAt}}",
  "configFormat": "table {{.ID}}\t{{.Name}}\t{{.CreatedAt}}\t{{.UpdatedAt}}",
  "serviceInspectFormat": "pretty",
  "nodesFormat": "table {{.ID}}\t{{.Hostname}}\t{{.Availability}}",
  "detachKeys": "ctrl-e,e",
  "credsStore": "secretservice",
  "credHelpers": {
    "awesomereg.example.org": "hip-star",
    "unicorn.example.com": "vcbait"
  },
  "plugins": {
    "plugin1": {
      "option": "value"
    },
    "plugin2": {
      "anotheroption": "anothervalue",
      "athirdoption": "athirdvalue"
    }
  },
  "proxies": {
    "default": {
      "httpProxy":  "http://user:pass@example.com:3128",
      "httpsProxy": "https://my-proxy.example.com:3129",
      "noProxy":    "intra.mycorp.example.com",
      "ftpProxy":   "http://user:pass@example.com:3128",
      "allProxy":   "socks://example.com:1234"
    },
    "https://manager1.mycorp.example.com:2377": {
      "httpProxy":  "http://user:pass@example.com:3128",
      "httpsProxy": "https://my-proxy.example.com:3129"
    }
  }
}
{% endraw %}
```
