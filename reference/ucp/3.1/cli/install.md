---
title: docker/ucp install
description: Install UCP on a node
keywords: ucp, cli, install
---

Install UCP on a node

## Usage

```bash
docker container run --rm -it \
    --name ucp \
    -v /var/run/docker.sock:/var/run/docker.sock \
    docker/ucp \
    install [command options]
```

## Description

This command initializes a new swarm, turns anode into a manager, and installs
Docker Universal Control Plane (UCP).

When installing UCP you can customize:

  * The UCP web server certificates. Create a volume named `ucp-controller-server-certs` and copy the `ca.pem`, `cert.pem`, and `key.pem` files to the root directory. Then run the install command with the `--external-server-cert` flag.
  * The license used by UCP, which you can accomplish by bind-mounting the file at `/config/docker_subscription.lic` in the tool. For example, `-v /path/to/my/config/docker_subscription.lic:/config/docker_subscription.lic` or by specifying the `--license $(cat license.lic)` option.

If you're joining more nodes to this swarm, open the following ports in your
firewall:

  * 443 or the `--controller-port`
  * 2376 or the `--swarm-port`
  * 12376, 12379, 12380, 12381, 12382, 12383, 12384, 12385, 12386, 12387
  * 4789 (udp) and 7946 (tcp/udp) for overlay networking

If you have SELinux policies enabled for your Docker install, you will need to
use `docker container run --rm -it --security-opt label=disable ...` when running this
command.

If you are installing on Azure, see [Install UCP on Azure](/ee/ucp/admin/install/install-on-azure/).

## Options

| Option                          | Description                                                                                                                                                                                                                                            |
|:--------------------------------|:-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `--debug, -D`                   | Enable debug mode                                                                                                                                                                                                                                      |
| `--jsonlog`                     | Produce json formatted output for easier parsing                                                                                                                                                                                                       |
| `--interactive, -i`             | Run in interactive mode and prompt for configuration values                                                                                                                                                                                            |
| `--admin-password` *value*      | The UCP administrator password [$UCP_ADMIN_PASSWORD]                                                                                                                                                                                                   |
| `--admin-username` *value*      | The UCP administrator username [$UCP_ADMIN_USER]                                                                                                                                                                                                       |
| `--binpack`                     | Set the Docker Swarm scheduler to binpack mode. Used for backwards compatibility                                                                                                                                                                       |
| `--cloud-provider` *value*      | The cloud provider for the cluster                                                                                                                                                                                                                     |
| `--cni-installer-url` *value*   | A URL pointing to a kubernetes YAML file to be used as an installer for the CNI plugin of the cluster. If specified, the default CNI plugin will not be installed. If the URL is using the HTTPS scheme, no certificate verification will be performed |
| `--controller-port` *value*     | Port for the web UI and API (default: 443)                                                                                                                                                                                                             |
| `--data-path-addr` *value*      | Address or interface to use for data path traffic. Format: IP address or network interface name [$UCP_DATA_PATH_ADDR]                                                                                                                                  |
| `--disable-tracking`            | Disable anonymous tracking and analytics                                                                                                                                                                                                               |
| `--disable-usage`               | Disable anonymous usage reporting                                                                                                                                                                                                                      |
| `--dns-opt` *value*             | Set DNS options for the UCP containers [$DNS_OPT]                                                                                                                                                                                                      |
| `--dns-search` *value*          | Set custom DNS search domains for the UCP containers [$DNS_SEARCH]                                                                                                                                                                                     |
| `--dns` *value*                 | Set custom DNS servers for the UCP containers [$DNS]                                                                                                                                                                                                   |
| `--enable-profiling`            | Enable performance profiling                                                                                                                                                                                                                           |
| `--existing-config`             | Use the latest existing UCP config during this installation. The install will fail if a config is not found                                                                                                                                            |
| `--external-server-cert`        | Customize the certificates used by the UCP web server                                                                                                                                                                                                  |
| `--external-service-lb` *value* | Set the IP address of the load balancer that published services are expected to be reachable on                                                                                                                                                        |
| `--force-insecure-tcp`          | Force install to continue even with unauthenticated Docker Engine ports.                                                                                                                                                                               |
| `--force-minimums`              | Force the install/upgrade even if the system does not meet the minimum requirements                                                                                                                                                                    |
| `--host-address` *value*        | The network address to advertise to other nodes. Format: IP address or network interface name [$UCP_HOST_ADDRESS]                                                                                                                                      |
| `--kube-apiserver-port` *value* | Port for the Kubernetes API server (default: 6443)                                                                                                                                                                                                     |
| `--kv-snapshot-count` *value*   | Number of changes between key-value store snapshots (default: 20000) [$KV_SNAPSHOT_COUNT]                                                                                                                                                              |
| `--kv-timeout` *value*          | Timeout in milliseconds for the key-value store (default: 5000) [$KV_TIMEOUT]                                                                                                                                                                          |
| `--license` *value*             | Add a license: e.g. --license "$(cat license.lic)" [$UCP_LICENSE]                                                                                                                                                                                      |
| `--nodeport-range` *value*      | Allowed port range for Kubernetes services of type NodePort (Default: 32768-35535) (default: "32768-35535")                                                                                                                                            |
| `--pod-cidr` *value*            | Kubernetes cluster IP pool for the pods to allocated IP from (Default: 192.168.0.0/16) (default: "192.168.0.0/16")                                                                                                                                     |
| `--preserve-certs`              | Don't generate certificates if they already exist                                                                                                                                                                                                      |
| `--pull` *value*                | Pull UCP images: 'always', when 'missing', or 'never' (default: "missing")                                                                                                                                                                             |
| `--random`                      | Set the Docker Swarm scheduler to random mode. Used for backwards compatibility                                                                                                                                                                        |
| `--registry-password` *value*   | Password to use when pulling images [$REGISTRY_PASSWORD]                                                                                                                                                                                               |
| `--registry-username` *value*   | Username to use when pulling images [$REGISTRY_USERNAME]                                                                                                                                                                                               |
| `--san` *value*                 | Add subject alternative names to certificates (e.g. --san www1.acme.com --san www2.acme.com) [$UCP_HOSTNAMES]                                                                                                                                          |
| `--skip-cloud-provider-check`   | Disables checks which rely on detecting which (if any) cloud provider the cluster is currently running on                                                                                                                                              |
| `--swarm-experimental`          | Enable Docker Swarm experimental features. Used for backwards compatibility                                                                                                                                                                            |
| `--swarm-grpc-port` *value*     | Port for communication between nodes (default: 2377)                                                                                                                                                                                                   |
| `--swarm-port` *value*          | Port for the Docker Swarm manager. Used for backwards compatibility (default: 2376)                                                                                                                                                                    |
| `--unlock-key` *value*          | The unlock key for this swarm-mode cluster, if one exists. [$UNLOCK_KEY]                                                                                                                                                                               |
| `--unmanaged-cni`               | Flag to indicate if cni provider is calico and managed by UCP (calico is the default CNI provider)                                                                                                                                                     |                                                                                                                                                   |

