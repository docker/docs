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

## Options

| Option                   | Description                                                                                                                                                                                                                               |
|:-------------------------|:------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `--admin-password`       | The UCP administrator password |
| `--admin-username`       | The UCP administrator username                             |
| `--binpack`              | Set the Docker Swarm scheduler to binpack mode. Used for backwards compatibility       |
| `--cloud-provider`       | The cloud provider for the cluster 
| `--cni-installer-url`    | Deprecated feature. A URL pointing to a Kubernetes YAML file to be used as an installer for the CNI plugin of the cluster. If specified, the default CNI plugin is not installed. If the URL uses the HTTPS scheme, no certificate verification is performed.       |
| `--controller-port`      | Port for the web UI and API 
| `--data-path-addr`       | Address or interface to use for data path traffic. Format: IP address or network interface name
| `--debug, D`             | Enable debug mode  |
| `--disable-tracking`     | Disable anonymous tracking and analytics                                               |
| `--disable-usage`        | Disable anonymous usage reporting                                                      |
| `--dns`                  | Set custom DNS servers for the UCP containers                                          |                                                                                                                                                   
| `--dns-opt`              | Set DNS options for the UCP containers                                                 |                                                                                                                                                   
| `--dns-search`           | Set custom DNS search domains for the UCP containers                                   |
| `--enable-profiling`     | Enable performance profiling                                                           |
| `--existing-config`      | Use the latest existing UCP config during this installation. The install fails if a config is not found.          |
| `--external-server-cert` | Use the certificates in the `ucp-controller-server-certs` volume instead of generating self-signed certs during installation                                                                                           |
| `--external-service-lb`  | Set the external service load balancer reported in the UI                              |
| `--force-insecure-tcp`   | Force install to continue even with unauthenticated Docker Engine ports                |
| `--force-minimums`       | Force the install/upgrade even if the system doesn't meet the minimum requirements.    |
| `--host-address`         | The network address to advertise to other nodes. Format: IP address or network interface name |
| `--interactive, i`       | Run in interactive mode and prompt for configuration values |
| `--jsonlog`              | Produce json formatted output for easier parsing |
| `--kube-apiserver-port`  | Port for the Kubernetes API server (default: 6443)                                     |
| `--kv-snapshot-count`    | Number of changes between key-value store snapshots                                    |
| `--kv-timeout`           | Timeout in milliseconds for the key-value store                                        |
| `--license`              | Add a license: e.g.` --license "$(cat license.lic)" `                                  |
| `--pod-cidr`             | Kubernetes cluster IP pool for the pods to allocated IPs from (Default: `192.168.0.0/16`) |
| `--preserve-certs`       | Don't generate certificates if they already exist                                      |
| `--pull`                 | Pull UCP images: `always`, when `missing`, or `never`                                  |
| `--random`               | Set the Docker Swarm scheduler to random mode. Used for backwards compatibility        |
| `--registry-username`    | Username to use when pulling images                                                    |
| `--registry-password`    | Password to use when pulling images                                                    |
| `--san`                  | Add subject alternative names to certificates (e.g. --san www1.acme.com --san www2.acme.com) |
| `--skip-cloud-provider`  | Disables checks that rely on detecting the cloud provider (if any) on which the cluster is currently running. | 
| `--swarm-experimental`   | Enable Docker Swarm experimental features. Used for backwards compatibility            |
| `--swarm-port`           | Port for the Docker Swarm manager. Used for backwards compatibility                    | 
| `--swarm-grpc-port`      | Port for communication between nodes                                                   | 
| `--unlock-key`           | The unlock key for this swarm-mode cluster, if one exists.                             |  
| `--unmanaged-cni`        |The default value of `false` indicates that Kubernetes networking is managed by UCP with its default managed CNI plugin, Calico. When set to `true`, UCP does not deploy or manage the lifecycle of the default CNI plugin - the CNI plugin is deployed and managed independently of UCP. Note that when `unmanaged-cni=true`, networking in the cluster will not function for Kubernetes until a CNI plugin is deployed.    |                                                                                                                                                                                                                                                                         
