---
title: Configure layer 7 routing service
description: Learn how to configure the layer 7 routing solution for UCP.
keywords: routing, proxy, interlock, load balancing
redirect_from:
  - /ee/ucp/interlock/deploy/configure/
  - /ee/ucp/interlock/usage/default-service/
---

To further customize the layer 7 routing solution, you must update the
`ucp-interlock` service with a new Docker configuration.

1. Find out what configuration is currently being used for the `ucp-interlock`
service and save it to a file:

   {% raw %}
   ```bash
   CURRENT_CONFIG_NAME=$(docker service inspect --format '{{ (index .Spec.TaskTemplate.ContainerSpec.Configs 0).ConfigName }}' ucp-interlock)
   docker config inspect --format '{{ printf "%s" .Spec.Data }}' $CURRENT_CONFIG_NAME > config.toml
   ```
   {% endraw %}

2. Make the necessary changes to the `config.toml` file. See [TOML file configuration options](#toml-file-configuration-options) for more information.

3. Create a new Docker configuration object from the `config.toml` file:

   ```bash
   NEW_CONFIG_NAME="com.docker.ucp.interlock.conf-$(( $(cut -d '-' -f 2 <<< "$CURRENT_CONFIG_NAME") + 1 ))"
   docker config create $NEW_CONFIG_NAME config.toml
   ```

3. Update the `ucp-interlock` service to start using the new configuration:

   ```bash
   docker service update \
     --config-rm $CURRENT_CONFIG_NAME \
     --config-add source=$NEW_CONFIG_NAME,target=/config.toml \
     ucp-interlock
   ```

By default, the `ucp-interlock` service is configured to pause if you provide an
invalid configuration. The service will not restart without manual intervention.

If you want the service to automatically rollback to a previous stable
configuration, you can update it with the following command:

```bash
docker service update \
  --update-failure-action rollback \
  ucp-interlock
```

**Note**: Every time you enable the layer 7 routing
solution from the UCP UI, the `ucp-interlock` service is started using the
default configuration.

If you've customized the configuration used by the `ucp-interlock` service,
you must update it again to use the Docker configuration object
you've created.

## TOML file configuration options
The following sections describe how to configure the primary Interlock services:

- Core
- Extension
- Proxy

### Core configuration

The core configuraton handles the Interlock service itself. These are the configuration options for the `ucp-interlock` service:

| Option             | Type        | Description                                                                                    |
|:-------------------|:------------|:-----------------------------------------------------------------------------------------------|
| `ListenAddr`       | string      | Address to serve the Interlock GRPC API. Defaults to `8080`.                                   |
| `DockerURL`        | string      | Path to the socket or TCP address to the Docker API. Defaults to `unix:///var/run/docker.sock` |
| `TLSCACert`        | string      | Path to the CA certificate for connecting securely to the Docker API.                          |
| `TLSCert`          | string      | Path to the certificate for connecting securely to the Docker API.                             |
| `TLSKey`           | string      | Path to the key for connecting securely to the Docker API.                                     |
| `AllowInsecure`    | bool        | Skip TLS verification when connecting to the Docker API via TLS.                               |
| `PollInterval`     | string      | Interval to poll the Docker API for changes. Defaults to `3s`.                                 |
| `EndpointOverride` | string      | Override the default GRPC API endpoint for extensions. The default is  detected via Swarm.     |
| `Extensions`       | []Extension | Array of extensions as listed below.                                                           |

### Extension configuration

Interlock must contain at least one extension to service traffic. The following options are available to configure the extensions:

| Option             | Type        | Description                                      |
|:-------------------|:------------|:-----------------------------------------------------------|
| `Image` | string | Name of the Docker Image to use for the extension service |
| `Args` | []string | Arguments to be passed to the Docker extension service upon creation |
| `Labels` | map[string]string | Labels to add to the extension service |
| `ContainerLabels` | map[string]string | labels to be added to the extension service tasks |
| `Constraints` | []string | one or more [constraints](https://docs.docker.com/engine/reference/commandline/service_create/#specify-service-constraints-constraint) to use when scheduling the extension service |
| `PlacementPreferences` | []string | one or more [placement prefs](https://docs.docker.com/engine/reference/commandline/service_create/#specify-service-placement-preferences-placement-pref) to use when scheduling the extension service |
| `ServiceName` | string | Name of the extension service |
| `ProxyImage` | string | Name of the Docker Image to use for the proxy service |
| `ProxyArgs` | []string | Arguments to be passed to the Docker proxy service upon creation |
| `ProxyLabels` | map[string]string | Labels to add to the proxy service |
| `ProxyContainerLabels` | map[string]string | labels to be added to the proxy service tasks |
| `ProxyServiceName` | string | Name of the proxy service |
| `ProxyConfigPath` | string | Path in the service for the generated proxy config |
| `ProxyReplicas` | uint | number of proxy service replicas |
| `ProxyStopSignal` | string | stop signal for the proxy service (i.e. `SIGQUIT`) |
| `ProxyStopGracePeriod` | string | stop grace period for the proxy service (i.e. `5s`) |
| `ProxyConstraints` | []string | one or more [constraints](https://docs.docker.com/engine/reference/commandline/service_create/#specify-service-constraints-constraint) to use when scheduling the proxy service |
| `ProxyPlacementPreferences` | []string | one or more [placement prefs](https://docs.docker.com/engine/reference/commandline/service_create/#specify-service-placement-preferences-placement-pref) to use when scheduling the proxy service |
| `ProxyUpdateDelay` | string | delay between rolling proxy container updates  |
| `ServiceCluster` | string | Name of the cluster this extension services |
| `PublishMode` | string (`ingress` or `host`) | Publish mode that the proxy service uses |
| `PublishedPort` | int | Port on which the proxy service serves non-SSL traffic |
| `PublishedSSLPort` | int | Port on which the proxy service serves SSL traffic |
| `Template` | string | Docker configuration object that is used as the extension template |
| `Config` | Config | Proxy configuration used by the extensions as described in the following table |

### Proxy
Options are made available to the extensions, and the extensions utilize the options needed for proxy service configuration. This provides overrides to the extension configuration.

Because Interlock passes the extension configuration directly to the extension, each extension has
different configuration options available.  Refer to the documentation for each extension for supported options:

- [Nginx](nginx-config.md)

#### Customize the default proxy service
The default proxy service used by UCP to provide layer 7 routing is NGINX. If users try to access a route that hasn't been configured, they will see the default NGINX 404 page:

![Default NGINX page](../../images/interlock-default-service-1.png){: .with-border}

You can customize this by labelling a service with
`com.docker.lb.default_backend=true`. In this case, if users try to access a route that's
not configured, they are redirected to this service.

As an example, create a `docker-compose.yml` file with:

```yaml
version: "3.2"

services:
  demo:
    image: ehazlett/interlock-default-app
    deploy:
      replicas: 1
      labels:
        com.docker.lb.default_backend: "true"
        com.docker.lb.port: 80
    networks:
      - demo-network

networks:
  demo-network:
    driver: overlay
```

Set up your CLI client with a [UCP client bundle](../../user-access/cli.md),
and deploy the service:

```bash
docker stack deploy --compose-file docker-compose.yml demo
```

If users try to access a route that's not configured, they are directed
to this demo service.

![Custom default page](../../images/interlock-default-service-2.png){: .with-border}

### Example Configuration
The following is an example configuration to use with the Nginx extension.

```toml
ListenAddr = ":8080"
DockerURL = "unix:///var/run/docker.sock"
PollInterval = "3s"

[Extensions.default]
  Image = "{{ page.ucp_org }}/interlock-extension-nginx:{{ page.ucp_version }}"
  Args = ["-D"]
  ServiceName = "interlock-ext"
  ProxyImage = "{{ page.ucp_org }}/ucp-interlock-proxy:{{ page.ucp_version }}"
  ProxyArgs = []
  ProxyServiceName = "interlock-proxy"
  ProxyConfigPath = "/etc/nginx/nginx.conf"
  ProxyStopGracePeriod = "3s"
  PublishMode = "ingress"
  PublishedPort = 80
  ProxyReplicas = 1
  TargetPort = 80
  PublishedSSLPort = 443
  TargetSSLPort = 443
  [Extensions.default.Config]
    User = "nginx"
    PidPath = "/var/run/proxy.pid"
    WorkerProcesses = 1
    RlimitNoFile = 65535
    MaxConnections = 2048
```

## Next steps

- [Configure host mode networking](host-mode-networking.md)
- [Configure an nginx extension](nginx-config.md)
- [Use application service labels](service-labels.md)
- [Tune the proxy service](tuning.md)
- [Update Interlock services](updates.md)
