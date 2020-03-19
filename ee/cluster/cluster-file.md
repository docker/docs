---
description: Cluster file reference and guidelines
keywords: documentation, docs, docker, cluster, infrastructure, automation
title: Cluster file version 1 reference
toc_max: 5
toc_min: 1
---

>{% include enterprise_label_shortform.md %}

This topic describes version 1 of the Cluster file format.

## Cluster file structure and examples

<div class="panel panel-default">
    <div class="panel-heading collapsed" data-toggle="collapse" data-target="#collapseSample1" style="cursor: pointer">
    Example Cluster file version 1
    <i class="chevron fa fa-fw"></i></div>
    <div class="collapse block" id="collapseSample1">
<pre><code>
variable:
  domain: "YOUR DOMAIN, e.g. docker.com"
  subdomain: "A SUBDOMAIN, e.g. cluster"
  region: "THE AWS REGION TO DEPLOY, e.g. us-east-1"
  email: "YOUR.EMAIL@COMPANY.COM"
  ucp_password:
    type: prompt
provider:
  acme:
    email: ${email}
    server_url: https://acme-staging-v02.api.letsencrypt.org/directory
  aws:
    region: ${region}
cluster:
  dtr:
    version: docker/dtr:2.6.5
  engine:
    version: ee-stable-18.09.5
  ucp:
    username: admin
    password: ${ucp_password}
    version: docker/ucp:3.1.6
resource:
  aws_instance:
    managers:
      instance_type: t2.xlarge
      os: Ubuntu 16.04
      quantity: 3
    registry:
      instance_type: t2.xlarge
      os: Ubuntu 16.04
      quantity: 3
    workers:
      instance_type: t2.xlarge
      os: Ubuntu 16.04
      quantity: 3
  aws_lb:
    apps:
      domain: ${subdomain}.${domain}
      instances:
      - workers
      ports:
      - 80:8080
      - 443:8443
    dtr:
      domain: ${subdomain}.${domain}
      instances:
      - registry
      ports:
      - 443:443
    ucp:
      domain: ${subdomain}.${domain}
      instances:
      - managers
      ports:
      - 443:443
      - 6443:6443
  aws_route53_zone:
    dns:
      domain: ${domain}
      subdomain: ${subdomain}
</code></pre>
    </div>
</div>

The topics on this reference page are organized alphabetically by top-level keys
to reflect the structure of the Cluster file. Top-level keys that define
a section in the configuration file, such as `cluster`, `provider`, and `resource`,
are listed with the options that support them as sub-topics. This information
maps to the indent structure of the Cluster file.

### cluster
Specifies components to install and configure for a cluster.

The following components are available:

- `subscription`: (Optional) Configuration options for Docker Enterprise
  Subscriptions.
- `cloudstor`: (Optional) Configuration options for Docker Cloudstor.
- `dtr`: (Optional) Configuration options for Docker Trusted Registry.
- `engine`: (Optional) Configuration options for Docker Engine.
- `ucp`: (Optional) Configuration options for Docker Universal Control Plane.
- `registry`: (Optional) Configuration options for authenticating nodes with a registry to pull Docker images.

#### subscription
Provide Docker Enterprise subscription information
-  `id`: (Optional) The subscription UUID for this Docker Enterprise
   installation `sub-xxxx`.
-  `license`: (Optional) The absolute path to a local Docker Enterprise license
   file `/path/to/docker_subscription.lic`.
-  `trial`: (Optional) Specify if this is a trial subscription. Default is
   `false`

#### cloudstor
Docker Cloudstor is a Docker Swarm Plugin that provides persistent storage to
Docker Swarm Clusters deployed on to AWS or Azure. By default Docker Cloudstor
is not installed on Docker Enterprise environments created with Docker Cluster.

```yaml
cluster:
  cloudstor:
    version: '1.0'
```

The following optional elements can be specified:

- `version`: (Required) The version of Docker Cloudstor to install. The default
  is `disabled`. The only released version of Docker Cloudstor at this time is
  `1.0`.
- `use_efs`: (Optional) Specifies whether an Elastic File System should be
  provisioned. By default Docker Cloudstor on AWS uses Elastic Block Store,
  therefore this value defaults to `false`.

#### dtr
Customizes the installation of Docker Trusted Registry.
```yaml
cluster:
  dtr:
    version: "docker/dtr:2.6.5"
    install_options:
    - "--debug"
    - "--enable-pprof"
```

The following optional elements can be specified:

- `version`: (Optional) The version of DTR to install.  Defaults to `docker/dtr:2.6.5`.
- `ca`: (Optional) The path to a root CA public certificate.
- `key`: (Optional) The path to a TLS private key.
- `cert`: (Optional) The path to a public key certificate.
- `install_options`: (Optional) Additional DTR install options.

#### engine
Customizes the installation of Docker Enterprise Engine.
```yaml
cluster:
  engine:
    channel: "stable"
    edition: "ee"
    version: "19.03"
```

The following optional elements can be specified:
- `version`: (Optional) The version of the Docker Engine to install.  Defaults to `19.03`.
- `edition`: (Optional) The family of Docker Engine to install.  Defaults to `ee` for Enterprise edition.
- `channel`: (Optional) The channel on the repository to pull updated packages.  Defaults to `stable`.
- `url`: (Optional) Defaults to "https://storebits.docker.com/ee".
- `storage_driver`: (Optional)  The storage driver to use for the storage volume.  Default
value is dependent on the operating system.
    - Amazon Linux 2 is `overlay2`.
    - Centos is `overlay2`.
    - Oracle Linux is `overlay2`.
    - RedHat is `overlay2`.
    - SLES is `btrfs`.
    - Ubuntu is `overlay2`.
- `storage_fstype`: (Optional) File system to use for storage volume.  Default value is dependent on the operating system.
    - Amazon Linux 2 is `xfs`.
    - Centos is `xfs`.
    - Oracle Linux is `xfs`.
    - RedHat is `xfs`.
    - SLES is `btrfs`.
    - Ubuntu is `ext4`.
- `storage_volume`: (Optional) Docker storage volume path for `/var/lib/docker` Default value is provider dependent.
    - AWS
        - non-NVME is `/dev/xvdb`.
        - NVME disks are one of `/dev/nvme[0-26]n1`.
    - Azure is `/dev/disk/azure/scsi1/lun0`.
- `daemon`: (Optional) Provides docker daemon options. Defaults to "".
- `ca`:  (dev) Defaults to "".
- `key`: (dev) Defaults to "".
- `enable_remote_tcp`: (dev) Enables direct access to docker engine. Defaults to `false`.

*dev indicates that the functionality is only for development and testing.

#### kubernetes
Enables provider-specific options for Kubernetes support.

##### AWS Kubernetes options

- `cloud_provider`: (Optional)Enable cloud provider support for Kubernetes.  Defaults to `false`.
- `ebs_persistent_volumes`: (Optional) Enable persistent volume support with EBS volumes. Defaults to `false`.
- `efs_persistent_volumes`:  (Optional) Enable persistent volume support with EFS.  Defaults to `false`.
- `load_balancer`: (Optional) Enable Kubernetes pods to instantiate a load-balancer.  Defaults to `false`.
- `nfs_storage`: (Optional) Install additional packages on node for NFS support.  Defaults to `false`.
- `lifecycle`: (Optional) Defaults to `owned`.

#### registry
Customizes the registry from which the installation should pull images.  By default, Docker Hub and credentials to access Docker Hub are used.

```yaml
cluster:
  registry:
    password: ${base64decode("TVJYeTNDQWpTSk5HTW1ZRzJQcE1kM0tVRlQ=")}
    url: https://index.docker.io/v1/
    username: user
```

The following optional elements can be specified:
- `username`: The username for logging in to the registry on each node.  Default value is the current docker user.
- `url`: The registry to use for pulling Docker images.  Defaults to "https://index.docker.io/v1/".
- `password`: The password for logging in to the registry on each node.  Default value is the current docker user's password base64 encoded and wrapped in a call to base64decode.

#### ucp

- `version`: Specifies the version of UCP to install.  Defaults to `docker/ucp:3.1.6`.
- `username`: Specifies the username of the first user to create in UCP.  Defaults to `admin`.
- `password`: Specifies the password of the first user to create in UCP.  Defaults to `dockerdocker`.
- `ca`: Specifies a path to a root CA public certificate.
- `key`: Specifies a path to a TLS private key.
- `cert`: Specifies a path to a public key certificate.
- `install_options`: Lists additional UCP install options.

##### Additional UCP configuration options
Docker Cluster also accepts all UCP configuration options and creates the initial UCP config on
installation. The following list provides supported options:
- `anonymize_tracking`: Anonymizes analytic data. Specify 'true' to hide the license ID. Defaults to 'false'.
- `audit_level`: Specifies the audit logging level. Leave empty for disabling audit logs (default).
Other valid values are 'metadata' and 'request'.
- `auto_refresh`: Specify 'true' to enable attempted automatic license renewal when the license
nears expiration. If disabled, you must manually upload renewed license after expiration. Defaults to 'true'.
- `azure_ip_count`: Sets the IP count for azure allocator to allocate IPs per Azure virtual machine.
- `backend`: Specifie the name of the authorization backend to use, either 'managed' or 'ldap'. Defaults to 'managed'.
- `calico_mtu`: Specifies the MTU (maximum transmission unit) size for the Calico plugin. Defaults to '1480'.
- `cloud_provider`: Specifies the cloud provider for the kubernetes cluster.
- `cluster_label`: Specifies a label to be included with analytics/.
- `cni_installer_url`: Specifies the URL of a Kubernetes YAML file to be used for installing a CNI plugin.
Only applies during initial installation. If empty, the default CNI plugin is used.
- `controller_port`: Configures the port that the 'ucp-controller' listens to. Defaults to '443'.
- `custom_header_name`: Specifies the name of the custom header with 'name' = '*X-Custom-Header-Name*'.
- `custom_header_value`: Specifies the value of the custom header with 'value' = '*Custom Header Value*'.
- `default_new_user_role`: Specifies the role that new users get for their private resource sets.
Values are 'admin', 'viewonly', 'scheduler', 'restrictedcontrol', or 'fullcontrol'. Defaults to 'restrictedcontrol'.
- `default_node_orchestrator`: Specifies the type of orchestrator to use for new nodes that are
joined to the cluster. Can be 'swarm' or 'kubernetes'. Defaults to 'swarm'.
- `disable_tracking`: Specify 'true' to disable analytics of API call information. Defaults to 'false'.
- `disable_usageinfo`: Specify 'true' to disable analytics of usage information. Defaults to 'false'.
- `dns`: Specifies a CSV list of IP addresses to add as nameservers.
- `dns_opt`: Specifies a CSV list of options used by DNS resolvers.
- `dns_search`: Specifies a CSV list of domain names to search when a bare unqualified hostname is
used inside of a container.
- `enable_admin_ucp_scheduling`: Specify 'true' to allow admins to schedule on containers on manager nodes.
Defaults to 'false'.
- `external_service_lb`: Specifies an optional external load balancer for default links to services with
exposed ports in the web interface.
- `host_address`: Specifies the address for connecting to the DTR instance tied to this UCP cluster.
- `log_host`: Specifies a remote syslog server to send UCP controller logs to. If omitted, controller
logs are sent through the default docker daemon logging driver from the 'ucp-controller' container.
- `idpMetadataURL`: Specifies the Identity Provider Metadata URL.
- `image_repository`: Specifies the repository to use for UCP images.
- `install_args`: Specifies additional arguments to pass to the UCP installer.
- `ipip_mtu`: Specifies the IPIP MTU size for the calico IPIP tunnel interface.
- `kube_apiserver_port`: Configures the port to which the Kubernetes API server listens.
- `kv_snapshot_count`: Sets the key-value store snapshot count setting. Defaults to '20000'.
- `kv_timeout`: Sets the key-value store timeout setting, in milliseconds. Defaults to '5000'.
- `lifetime_minutes`: Specifies the initial session lifetime, in minutes. Defaults to `4320`, which is 72 hours.
- `local_volume_collection_mapping`: Stores data about collections for volumes in UCP's local KV store
instead of on the volume labels. This is used for enforcing access control on volumes.
- `log_level`: Specifies the logging level for UCP components. Values are syslog priority
levels (https://linux.die.net/man/5/syslog.conf): 'debug', 'info', 'notice', 'warning', 'err', 'crit', 'alert',
and 'emerg'.
- `managedPasswordDisabled`: Indicates if managed password is disabled. Defaults to false.
- `managedPasswordFallbackUser`: The fallback user when the managed password authentication is disabled. Defaults to "".
- `manager_kube_reserved_resources`: Specifies reserve resources for Docker UCP and Kubernetes components
that are running on manager nodes.
- `metrics_disk_usage_interval`: Specifies the interval for how frequently storage metrics are gathered.
This operation can impact performance when large volumes are present.
- `metrics_retention_time`: Adjusts the metrics retention time.
- `metrics_scrape_interval`: Specifies the interval for how frequently managers gather metrics from nodes in the cluster.
- `nodeport_range`: Specifies the port range that for Kubernetes services of type NodePort can be exposed in.
Defaults to '32768-35535'.
- `per_user_limit`: Specifies the maximum number of sessions that a user can have active simultaneously. If
the creation of a new session would put a user over this limit, the least recently used session is deleted.
A value of zero disables limiting the number of sessions that users can have. Defaults to `5`.
- `pod_cidr`: Specifies the subnet pool from which the IP for the Pod should be allocated from the CNI ipam plugin.
- `profiling_enabled`: Specify 'true' to enable specialized debugging endpoints for profiling UCP performance.
Defaults to 'false'.
- `log_protocol`: Specifies the protocol to use for remote logging. Values are 'tcp' and 'udp'. Defaults to 'tcp'.
- `renewal_threshold_minutes`: Specifies the length of time, in minutes, before the expiration of a
session. When used, a session is extended by the current configured lifetime from that point in time. A zero value disables session extension. Defaults to `1440`, which is 24 hours.
- `require_content_trust`: Specify 'true' to require images be signed by content trust. Defaults to 'false'.
- `require_signature_from`: Specifies a csv list of users or teams required to sign images.
- `rethinkdb_cache_size`: Sets the size of the cache used by UCP's RethinkDB servers. TDefaults to 1GB,
but leaving this field empty or specifying `auto` instructs RethinkDB to determine a cache size automatically.
- `rootCerts`: Defaults to empty.
- `samlEnabled`: Indicates if saml is used.
- `samlLoginText`: Specifies the customized SAML login button text.
- `service_id`: Specifies the DTR instance's OpenID Connect Client ID, as registered with the Docker
authentication provider.
- `spHost`: Specifies the Service Provider Host.
- `storage_driver`: Specifies the UCP storage driver to install.
- `support_dump_include_audit_logs`: When set to `true`, support dumps include audit logs in the logs
of the 'ucp-controller' container of each manager node. Defaults to 'false'.
- `swarm_port`: Configures the port that the 'ucp-swarm-manager' listens to. Defaults to '2376'.
- `swarm_strategy`: Configures placement strategy for container scheduling.
This doesn't affect swarm-mode services. Values are 'spread', 'binpack', and 'random'.
- `tlsSkipVerify`: Specifies TLS Skip verify for IdP Metadata.
- `unmanaged_cni`: Defaults to 'false'.
- `worker_kube_reserved_resources`: Reserves resources for Docker UCP and Kubernetes components
that are running on worker nodes.
- `custom_kube_api_server_flags`: Specifies the configuration options for the Kubernetes API server. (dev)
- `custom_kube_controller_manager_flags`: Specifies the configuration options for the Kubernetes controller manager. (dev)
- `custom_kube_scheduler_flags`: Specifies the configuration options for the Kubernetes scheduler. (dev)
- `custom_kubelet_flags`: Specifies the configuration options for Kubelets. (dev)

*dev indicates that the functionality is only for development and testing. Arbitrary Kubernetes configuration parameters are not tested and supported under the Docker Enterprise Software Support Agreement.

#### vpc

If you are deploying on to AWS, by default Docker Cluster will create a new AWS
VPC (Virtual Private Cloud) for the Docker Enterprise resources. To specify an
existing VPC, a user can specify a VPC ID in the Cluster File.

```yaml
cluster:
  vpc:
    id: vpc-existing-vpc-id
```

Docker Cluster assumes the VPC CIDR is `172.31.0.0/16`, so will therefore
attempt to create AWS subnets from this range. Docker Cluster can not utilise
existing AWS subnets. To instruct Docker Cluster to provision subnets for an
alternative CIDR you can pass a new CIDR into the Cluster File.

```yaml
cluster:
  vpc:
    id: vpc-existing-vpc-id
    cidr: "192.168.0.0/16"
```

The following elements can be specified:

- `id` - (Required) The existing AWS VPC ID `vpc-xxx`
- `cidr` - If the VPC's CIDR is not the default `172.31.0.0/16` an alternative
  CIDR can be specified here.

### provider
Defines where the cluster's resources are provisioned, as well as provider-specific configuration such as tags.

{% raw %}
```yaml
provider:
  acme:
    email: ${email}
    server_url: https://acme-staging-v02.api.letsencrypt.org/directory
  aws:
    region: ${region}
```
{% endraw %}

#### acme
The Automated Certificate Management Environment (ACME) is an evolving standard for the automation of a domain-validated certificate authority. Docker Cluster uses the ACME provider to create SSL certificates that are signed by [Let's Encrypt](https://letsencrypt.org/).

The ACME provider Configuration for the ACME provider supports arguments that closely align with the [Terraform ACME provider](https://www.terraform.io/docs/providers/acme/index.html):

The following elements can be specified:
- `email`: (Required) The email to associate the certificates with.
- `server_url`: (Optional) The URL to the ACME endpoint's directory.  Default is "https://acme-v02.api.letsencrypt.org/directory"

#### aws
Configuration for the AWS provider supports arguments that closely align with the [Terraform AWS provider](https://www.terraform.io/docs/providers/aws/index.html).

```yaml
aws:
  region: "us-east-1"
  tags:
    Owner: "Infra"
    Environment: "Test"
```
The following elements can be specified:
- `region` - (Required) This is the AWS region. It can be sourced from the `AWS_DEFAULT_REGION` environment variables, or
  via a shared credentials file if `profile` is specified.
- `tags` - (Optional) Additional name value pairs to assign to every resource (which
 supports tagging) in the cluster.
- `access_key` - (Required) This is the AWS access key. It can be sourced from
the `AWS_ACCESS_KEY_ID` environment variable, or via
  a shared credentials file if `profile` is specified.
- `secret_key` - (Required) This is the AWS secret key. It can be sourced from
the `AWS_SECRET_ACCESS_KEY` environment variable, or
  via a shared credentials file if `profile` is specified.
- `profile` - (Optional) This is the AWS profile name as set in the shared credentials
  file.
- `assume_role` - (Optional) An `assume_role` block (documented below). Only one
  `assume_role` block can be in the configuration.
- `endpoints` - (Optional) Configuration block for customizing service endpoints. See the
[Custom Service Endpoints Guide](/docs/providers/aws/guides/custom-service-endpoints.html)
for more information about connecting to alternate AWS endpoints or AWS compatible solutions.
- `shared_credentials_file` = (Optional) This is the path to the shared
 credentials file.  If this is not set and a profile is specified,
 `~/.aws/credentials` is used.
- `token` - (Optional) Session token for validating temporary credentials.
Typically provided after successful identity federation or Multi-Factor
Authentication (MFA) login. With MFA login, this is the session token
provided afterwards, not the 6 digit MFA code used to get temporary
credentials.  It can also be sourced from the `AWS_SESSION_TOKEN`
environment variable.
- `max_retries` - (Optional) This is the maximum number of times an API
  call is retried, in the case where requests are being throttled or
  experiencing transient failures. The delay between the subsequent API
  calls increases exponentially.
- `allowed_account_ids` - (Optional) List of allowed, white listed, AWS
  account IDs to prevent you from mistakenly using an incorrect one (and
  potentially end up destroying a live environment). Conflicts with
  `forbidden_account_ids`.
- `forbidden_account_ids` - (Optional) List of forbidden, blacklisted,
  AWS account IDs to prevent you mistakenly using a wrong one (and
  potentially end up destroying a live environment). Conflicts with
  `allowed_account_ids`.
- `insecure` - (Optional) Explicitly allows the provider to
  perform "insecure" SSL requests. If omitted, defaults to `false`.
- `skip_credentials_validation` - (Optional) Skips the credentials
  validation via the STS API. Useful for AWS API implementations that do
  not have STS available or implemented.
- `skip_get_ec2_platforms` - (Optional) Skips getting the supported EC2
  platforms. Used by users that don't have `ec2:DescribeAccountAttributes`
  permissions.
- `skip_region_validation` - (Optional) Skips validation of provided region name.
  Useful for AWS-like implementations that use their own region names
  or to bypass the validation for regions that aren't publicly available yet.

### resource
Resources to provision for a cluster.  Resources are organized as shown in the following example:

```yaml
resource:
  type:
    name:
      parameters
```
For a given `type`, there may be more one or more named resources to provision.

For a given `name`, a resource may have one or more parameters.

#### aws_instance

```yaml
resource:
  aws_instance:
    workers:
      instance_type: t2.xlarge
      price: 0.25
      os: Ubuntu 16.04
```
- `quantity`: (Required) The number of instances to create.
- `os`: An alias that is expanded by `docker cluster` to the AMI owner and AMI name to install.
The following aliases are supported by `docker cluster`:
    - `CentOS 7`
    - `RHEL 7.1`
    - `RHEL 7.2`
    - `RHEL 7.3`
    - `RHEL 7.4`
    - `RHEL 7.5`
    - `RHEL 7.6`
    - `Oracle Linux 7.3`
    - `Oracle Linux 7.4`
    - `Oracle Linux 7.5`
    - `SLES 12.2`
    - `SLES 12.3`
    - `SLES 15`
    - `Ubuntu 14.04`
    - `Ubuntu 16.04`
    - `Ubuntu 18.04`
    - `Windows Server 2016`
    - `Windows Server 1709`
    - `Windows Server 1803`
    - `Windows Server 2019`
    > Note
    > 
    > Make sure the OS you select is [compatible](https://success.docker.com/article/compatibility-matrix)
    with the product you're installing. Docker Cluster validates the support during installation.
- `instance_type`: Specifies the [AWS instance type](https://aws.amazon.com/ec2/instance-types/) to provision.
- `key_name`: By default, Docker Cluster creates an [AWS EC2 Key Pair](https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/ec2-key-pairs.html) and registers it with AWS for the cluster.
To use an existing AWS EC2 Key Pair, set this value to the name of the AWS EC2 Key Pair.
- `ssh_private_key`: By default, Docker Cluster creates an [AWS EC2 Key Pair](https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/ec2-key-pairs.html) and registers it with AWS for the cluster. To use an existing AWS EC2 Key Pair, set this value to the path of the private SSH key.
- `username`: Specifies the username for the node with Administrative privileges. By default, the `os` option
sets this to the well-known username for the AMIs (which can change by distribution):
    - Amazon Linux 2 is `ec2-user`.
    - Centos is `centos`.
    - Oracle Linux is `ec2-user`.
    - RedHat is `ec2-user`.
    - SLES is `ec2-user`.
    - Ubuntu is `ubuntu`.
    - Windows is `Administrator`.
- `password`: This value is only used by Windows nodes. By default, Windows nodes have a random password generated.
- `ami`: Specifies a custom AMI, or one that's not currently available as an OS. Specify either the id or
the owner/name to query for the latest.
    - `id`: Specifies the ID of the AMI. For example, `ami-0510c89f1a2691cf2`.
    - `owner`: Specifies the AWS account ID of the image owner. For example, `099720109477`.
    - `name`: Specifies the name of the AMI that was provided during image creation. For example, `ubuntu/images/hvm-ssd/ubuntu-xenial-16.04-amd64-server-*`.
    - `platform`: Specify `windows` for Windows instances.
- `tags`: (Optional) Specifies additional name value pairs to assign to every instance.
- `swarm_labels`: (Optional) Specifies additional key value pairs that represent swarm labels to apply to every node.

#### aws_spot_instance_request

Provisions a spot instance request in AWS to dramatically reduce the cost of instances. Spot instance
availability is not guaranteed. Therefore, it is recommended to use `aws_spot_instance_request` for
additional worker nodes and not for mission-critical nodes like managers and registry.

```yaml
resource:
  aws_spot_instance_request:
    workers:
      instance_type: t2.xlarge
      price: 0.25
      os: Ubuntu 16.04
      quantity: 3
```

Supports the same set of parameters as [aws_instance](index.md#aws_instance), with the addition of an optional price to limit the max bid for a spot instance.
- `price`: (Optional) Specifies a maximum price to bid on the spot instance.

#### aws_lb
Provisions an AWS Load Balancer.
```yaml
resource:
  aws_lb:
    ucp:
      domain: "example.com"
      instances:
      - managers
      ports:
      - 443:443
      - 6443:6443
```
The following options are supported:

- `instances`: (Required) Specifies a list of `aws_instance` and `aws_spot_instance_request` names to
attach to the load balancer.
- `ports`: (Required) Specifies a list of `listening port[/protocol]:target port[/protocol]` mappings
to define how the load balancer should route traffic. By default, the protocol is `tcp`.
- `domain`: Specifies the domain in which to create DNS records for this load balancer. The record is named the
same as this resource, appended by the domain. For example, if the resource is `ucp` and the domain is `example.com`,
the `A` record is `ucp.example.com`.
- `internal`: (Optional) Defaults to `false`.
- `type`: (Optional) Defaults to `network`.
- `enable_cross_zone_load_balancing`: (Optional) Defaults to `false`.

#### aws_route53_zone
Creates a subdomain in an AWS route53 zone. The following example creates a public zone for `testing.example.com`:

```yaml
resource:
  aws_route53_zone:
    dns:
      domain: example.com
      subdomain: testing
```
The following elements are required:
- `domain`: (Required) Specifies the name of the hosted zone.
- `subdomain`: (Required) Specifies the subdomain to create in the `domain` hosted zone.

### variable
Docker cluster supports basic parameterization. The variable section defines a make of keys and values. A key can have a sub-key named `type`, which changes the behavior of the variable.

```yaml
variable:
  region: "us-east-1"
  password:
    prompt: true
  instance_type:
    env: AWS_INSTANCE_TYPE
```

Variables are referenced in the cluster definition as `${variable_name}`. For
example, `${region}` is substituted as `us-east-2` through the cluster
definition.

In addition to providing a literal value for variables, you can specify values
by:

 - `prompt: true` - Request the value from the user and do not echo characters
   as the value is entered.
 - `env`: Obtain the value from an environment variable

### Where to go next
 * [UCP CLI reference](https://docs.docker.com/reference/)
 * [DTR CLI reference](https://docs.docker.com/reference/)
 
