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

<table border=0>
   <tr>
    <td><b>Option</b></td>
    <td><b>Description</b></td>
   </tr>
   <tr>
    <td><code>--debug, D</code></td>
    <td>Enable debug mode</td>
   </tr>
   <tr>
    <td><code>--jsonlog</code></td>
    <td>Produce json formatted output for easier parsing</td>
   </tr>
   <tr>
    <td><code>--jsonlog</code></td>
    <td>Produce json formatted output for easier parsing</td>
   </tr>
   <tr>
    <td><code>--interactive, i</code></td>
    <td> Run in interactive mode and prompt for configuration values </td>
   </tr>
   <tr>
    <td><code>--admin-username</code></td>
    <td> The UCP administrator username</td>
   </tr>
   <tr>
    <td><code>--admin-password</code></td>
    <td>The UCP administrator password</td>
   </tr>
   <tr>
    <td><code>--san</code></td>
    <td>Add subject alternative names to certificates (e.g. <code>--san www1.acme.com --san www2.acme.com<code>)</td>
   </tr>
   <tr>
    <td><code>--unmanaged-cni</codde></td>
    <td>This determines who manages the CNI plugin, using <code>true</code> or <code>false</code>. The default is <code>false</code> The <code>true</code> value installs UCP without a managed CNI plugin. UCP and the Kubernetes components will be running but pod to pod networking will not function until a CNI plugin is manually installed. This will impact some functionality of UCP until a CNI plugin is running.</td>
   </tr>
   <tr>
    <td><code>--host-address</code></td>
    <td>The network address to advertise to other nodes. Format: IP address or network interface name</td>
   </tr>
   <tr>
    <td><code>--data-path-addr</code></td>
    <td>Address or interface to use for data path traffic. Format: IP address or network interface name</td>
   </tr>
   <tr>
    <td><code>--controller-port</code></td>
    <td>Port for the web UI and API</td>
   </tr>
   <tr>
    <td><code>--kube-apiserver-port</code></td>
    <td>Port for the Kubernetes API server (default: 6443)</td>
   </tr>
   <tr>
    <td><code>--swarm-port</code></td>
    <td>Port for the Docker Swarm manager. Used for backwards compatibility</td>
   </tr>
   <tr>
    <td><code>--swarm-grpc-port</code></td>
    <td>Port for communication between nodes</td>
   </tr>
   <tr>
    <td><code>--cni-installer-url</code></td>
    <td> Deprecated feature. A URL pointing to a Kubernetes YAML file to be used as an installer for the CNI plugin of the cluster. If specified, the default CNI plugin is not installed. If the URL uses the HTTPS scheme, no certificate verification is performed.</td>
   </tr>
   <tr>
    <td><code>--pod-cidr</code></td>
    <td>Kubernetes cluster IP pool for the pods to allocated IPs from (Default: `192.168.0.0/16`)</td>
   </tr>
   <tr>
    <td><code>--cloud-provider</code></td>
    <td>The cloud provider for the cluster</td>
   </tr>
   <tr>
    <td><code>--dns</code></td>
    <td>Set custom DNS servers for the UCP containers</td>
   </tr>
   <tr>
    <td><code>--dns-opt</code></td>
    <td>Set DNS options for the UCP containers</td>
   </tr>
   <tr>
    <td><code>--dns-search</code></td>
    <td>Set custom DNS search domains for the UCP containers</td>
   </tr>
   <tr>
    <td><code>--unlock-key</code></td>
    <td>The unlock key for this swarm-mode cluster, if one exists.</td>
   </tr>
   <tr>
    <td><code>--existing-config</code></td>
    <td>Use the latest existing UCP config during this installation. The install fails if a configuration is not found.</td>
   </tr>
   <tr>
    <td><code>--force-minimums</code></td>
    <td>Force the install/upgrade even if the system doesn't meet the minimum requirements.</td>
   </tr>
   <tr>
    <td><code>--pull</code></td>
    <td>Pull UCP images: <code>always</code> when <code>missing</code> or <code>never</code>.</td>
   </tr>
   <tr>
    <td><code>--registry-username</code></td>
    <td>Username to use when pulling images</td>
   </tr>
   <tr>
    <td><code>--registry-password</code></td>
    <td>Password to use when pulling images</td>
   </tr>
   <tr>
    <td><code>--kv-timeout</code></td>
    <td>Timeout in milliseconds for the key-value store</td>
   </tr>
   <tr>
    <td><code>--kv-snapshot-count</code></td>
    <td>Number of changes between key-value store snapshots</td>
   </tr>
   <tr>
    <td><code>--swarm-experimental</code></td>
    <td>Enable Docker Swarm experimental features. Used for backwards compatibility.</td>
   </tr>
   <tr>
    <td><code>--disable-tracking</code></td>
    <td>Disable anonymous tracking and analytics</td>
   </tr>
   <tr>
    <td><code>--disable-usage</code></td>
    <td>Disable anonymous usage reporting</td>
   </tr>
   <tr>
    <td><code>--external-server-cert</code></td>
    <td>Use the certificates in the <code>ucp-controller-server-certs</code> volume instead of generating self-signed certs during installation.</td>
   </tr>
   <tr>
    <td><code>--preserve-certs</code></td>
    <td>Do not generate certificates if they already exist.</td>
   </tr>
   <tr>
    <td><code>--binpack</code></td>
    <td>Set the Docker Swarm scheduler to binpack mode. Used for backwards compatibility.</td>
   </tr>
   <tr>
    <td><code>--random</code></td>
    <td>Set the Docker Swarm scheduler to random mode. Used for backwards compatibility.</td>
   </tr>
   <tr>
    <td><code>--external-service-lb</code></td>
    <td>Set the external service load balancer reported in the UI</td>
   </tr>
   <tr>
    <td><code>--enable-profiling</code></td>
    <td>Enable performance profiling</td>
   </tr>
   <tr>
    <td><code>--license</code></td>
    <td>Add a license: e.g. <code>--license "$(cat license.lic)"</code>
   </tr>
   <tr>
    <td><code>--force-insecure-tcp</code></td>
    <td>Force install to continue even with unauthenticated Docker Engine ports</td>                   
   </tr>
</table>
