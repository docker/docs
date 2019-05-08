---
title: UCP configuration file
description: Set up UCP deployments by using a configuration file.
keywords: Docker EE, UCP, configuration, config
---

There are two ways to configure UCP:
- through the web interface, or
- by importing and exporting the UCP config in a TOML file. For more information about TOML, see the [TOML README on GitHub](https://github.com/toml-lang/toml/blob/master/README.md).

You can customize the UCP installation by creating a configuration file at the
time of installation. During the installation, UCP detects and starts using the
configuration specified in this file.

## The UCP configuration file

You can use the configuration file in different ways to set up your UCP
cluster.

- Install one cluster and use the UCP web interface to configure it as desired,
  export the configuration file, edit it as needed, and then import the edited
  configuration file into multiple other clusters.
- Install a UCP cluster, export and edit the configuration file, and then use the
  API to import the new configuration into the same cluster.
- Run the `example-config` command, edit the example configuration file, and
  set the configuration at install time or import after installation.

Specify your configuration settings in a TOML file.

## Export and modify an existing configuration

Use the `config-toml` API to export the current settings and write them to a file. Within the directory of a UCP admin user's [client certificate bundle](../../user-access/cli.md), the following command exports the current configuration for the UCP hostname `UCP_HOST` to a file named `ucp-config.toml`:

### Get an authtoken

```
AUTHTOKEN=$(curl --silent --insecure --data '{"username":"<username>","password":"<password>"}' https://UCP_HOST/auth/login | jq --raw-output .auth_token)
```

### Download config file

```
curl -X GET "https://UCP_HOST/api/ucp/config-toml" -H  "accept: application/toml" -H  "Authorization: Bearer $AUTHTOKEN" > ucp-config.toml
```

### Upload config file

```
curl -X PUT -H  "accept: application/toml" -H "Authorization: Bearer $AUTHTOKEN" --upload-file 'path/to/ucp-config.toml' https://UCP_HOST/api/ucp/config-toml
```

## Apply an existing configuration file at install time
You can configure UCP to import an existing configuration file at install time. To do this using the **Configs** feature of Docker Swarm, follow these steps.

1. Create a **Docker Swarm Config** object with a name of `com.docker.ucp.config` and the TOML value of your UCP configuration file contents.
2. When installing UCP on that cluster, specify the `--existing-config` flag to have the installer use that object for its initial configuration.
3. After installation, delete the `com.docker.ucp.config` object.

## Example configuration file

You can see an example TOML config file that shows how to configure UCP
settings. From the command line, run UCP with the `example-config` option:

```bash
docker container run --rm {{ page.ucp_org }}/{{ page.ucp_repo }}:{{ page.ucp_version }} example-config
```

## Configuration options

### auth table

| Parameter               | Required | Description                                                                                                                                                                          |
|:------------------------|:---------|:-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `backend`               | no       | The name of the authorization backend to use, either `managed` or `ldap`. The default is `managed`.                                                                                  |
| `default_new_user_role` | no       | The role that new users get for their private resource sets. Values are `admin`, `viewonly`, `scheduler`, `restrictedcontrol`, or `fullcontrol`. The default is `restrictedcontrol`. |


### auth.sessions

| Parameter                   | Required | Description                                                                                                                                                                                                                                                                             |
|:----------------------------|:---------|:----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `lifetime_minutes`          | no       | The initial session lifetime, in minutes. The default is 4320, which is 72 hours.                                                                                                                                                                                                       |
| `renewal_threshold_minutes` | no       | The length of time, in minutes, before the expiration of a session where, if used, a session will be extended by the current configured lifetime from then. A zero value disables session extension. The default is 1440, which is 24 hours.                                            |
| `per_user_limit`            | no       | The maximum number of sessions that a user can have active simultaneously. If creating a new session would put a user over this limit, the least recently used session will be deleted. A value of zero disables limiting the number of sessions that users may have. The default is 5. |

### registries array (optional)

An array of tables that specifies the DTR instances that the current UCP instance manages.

| Parameter      | Required | Description                                                                                                                                                                                 |
|:---------------|:---------|:--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `host_address` | yes      | The address for connecting to the DTR instance tied to this UCP cluster.                                                                                                                    |
| `service_id`   | yes      | The DTR instance's OpenID Connect Client ID, as registered with the Docker authentication provider.                                                                                         |
| `ca_bundle`    | no       | If you're using a custom certificate authority (CA), `ca_bundle` specifies the root CA bundle for the DTR instance. The value is a string with the contents of a `ca.pem` file. |

### custom headers (optional)

Included when you need to set custom API headers. You can repeat this section multiple times to specify multiple separate headers. If you include custom headers, you must specify both `name` and `value`.

[[custom_api_server_headers]]

| Item | Description |
| ----------- | ----------- |
| `name`       | Set to specify the name of the custom header with `name` = "*X-Custom-Header-Name*". |
| `value` | Set to specify the value of the custom header with `value` = "*Custom Header Value*".  |


### audit_log_configuration table (optional)
Configures audit logging options for UCP components.

| Parameter      | Required | Description                                                                                                                                                                                 |
|:---------------|:---------|:--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `level`        | no       | Specify the audit logging level. Leave empty for disabling audit logs (default). Other legal values are `metadata` and `request`. |
| `support_dump_include_audit_logs` | no | When set to true, support dumps will include audit logs in the logs of the `ucp-controller` container of each manager node. The default is `false`. |


### scheduling_configuration table (optional)

Specifies scheduling options and the default orchestrator for new nodes.

> **Note**: If you run the `kubectl` command, such as `kubectl describe nodes`, to view scheduling rules on Kubernetes nodes, it does not reflect what is configured in UCP Admin settings. UCP uses taints to control container scheduling on nodes and is unrelated to kubectl's `Unschedulable` boolean flag.

| Parameter                     | Required | Description                                                                                                                                |
|:------------------------------|:---------|:-------------------------------------------------------------------------------------------------------------------------------------------|
| `enable_admin_ucp_scheduling` | no       | Set to `true` to allow admins to schedule on containers on manager nodes. The default is `false`.                                          |
| `default_node_orchestrator`   | no       | Sets the type of orchestrator to use for new nodes that are joined to the cluster. Can be `swarm` or `kubernetes`. The default is `swarm`. |

### tracking_configuration table (optional)

Specifies the analytics data that UCP collects.

| Parameter            | Required | Description                                                                             |
|:---------------------|:---------|:----------------------------------------------------------------------------------------|
| `disable_usageinfo`  | no       | Set to `true` to disable analytics of usage information. The default is `false`.        |
| `disable_tracking`   | no       | Set to `true` to disable analytics of API call information. The default is `false`.     |
| `anonymize_tracking` | no       | Anonymize analytic data. Set to `true` to hide your license ID. The default is `false`. |
| `cluster_label`      | no       | Set a label to be included with analytics/                                              |

### trust_configuration table (optional)

Specifies whether DTR images require signing.

| Parameter                | Required | Description                                                                         |
|:-------------------------|:---------|:------------------------------------------------------------------------------------|
| `require_content_trust`  | no       | Set to `true` to require images be signed by content trust. The default is `false`. |
| `require_signature_from` | no       | A string array that specifies users or teams which must sign images.                |

### log_configuration table (optional)

> Note: This feature has been deprecated. Refer to the [Deprecation notice](https://docs.docker.com/ee/ucp/release-notes/#deprecation-notice) for additional information.

Configures the logging options for UCP components.

| Parameter  | Required | Description                                                                                                                                                                                     |
|:-----------|:---------|:------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `protocol` | no       | The protocol to use for remote logging. Values are `tcp` and `udp`. The default is `tcp`.                                                                                                       |
| `host`     | no       | Specifies a remote syslog server to send UCP controller logs to. If omitted, controller logs are sent through the default docker daemon logging driver from the `ucp-controller` container.     |
| `level`    | no       | The logging level for UCP components. Values are [syslog priority  levels](https://linux.die.net/man/5/syslog.conf): `debug`, `info`, `notice`, `warning`, `err`, `crit`, `alert`, and `emerg`. |

### license_configuration table (optional)

Specifies whether the your UCP license is automatically renewed.

| Parameter      | Required | Description                                                                                                                                                                                   |
|:---------------|:---------|:----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `auto_refresh` | no       | Set to `true` to enable attempted automatic license renewal when the license nears expiration. If disabled, you must manually upload renewed license after expiration. The default is `true`. |

### cluster_config table (required)

Configures the cluster that the current UCP instance manages.

The `dns`, `dns_opt`, and `dns_search` settings configure the DNS settings for UCP
components. Assigning these values overrides the settings in a container's
`/etc/resolv.conf` file. For more info, see
[Configure container DNS](/engine/userguide/networking/default_network/configure-dns/).

| Parameter                              | Required | Description                                                                                                                                                                                      |
|:---------------------------------------|:---------|:-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `controller_port`                      | yes      | Configures the port that the `ucp-controller` listens to. The default is `443`.                                                                                                                  |
| `kube_apiserver_port`                  | yes      | Configures the port the Kubernetes API server listens to.                                                                                                                                        |
| `swarm_port`                           | yes      | Configures the port that the `ucp-swarm-manager` listens to. The default is `2376`.                                                                                                              |
| `swarm_strategy`                       | no       | Configures placement strategy for container scheduling. This doesn't affect swarm-mode services. Values are `spread`, `binpack`, and `random`.                                                   |
| `dns`                                  | yes      | Array of IP addresses to add as nameservers.                                                                                                                                                     |
| `dns_opt`                              | yes      | Array of options used by DNS resolvers.                                                                                                                                                          |
| `dns_search`                           | yes      | Array of domain names to search when a bare unqualified hostname is used inside of a container.                                                                                                  |
| `profiling_enabled`                    | no       | Set to `true` to enable specialized debugging endpoints for profiling UCP performance. The default is `false`.                                                                                   |
| `kv_timeout`                           | no       | Sets the key-value store timeout setting, in milliseconds. The default is `5000`.                                                                                                                |
| `kv_snapshot_count`                    | no       | Sets the key-value store snapshot count setting. The default is `20000`.                                                                                                                         |
| `external_service_lb`                  | no       | Specifies an optional external load balancer for default links to services with exposed ports in the web interface.                                                                              |
| `cni_installer_url`                    | no       | Specifies the URL of a Kubernetes YAML file to be used for installing a CNI plugin. Applies only during initial installation. If empty, the default CNI plugin is used.                          |
| `metrics_retention_time`               | no       | Adjusts the metrics retention time.                                                                                                                                                              |
| `metrics_scrape_interval`              | no       | Sets the interval for how frequently managers gather metrics from nodes in the cluster.                                                                                                          |
| `metrics_disk_usage_interval`          | no       | Sets the interval for how frequently storage metrics are gathered. This operation can be expensive when large volumes are present.                                                               |
| `rethinkdb_cache_size`                 | no       | Sets the size of the cache used by UCP's RethinkDB servers. The default is 1GB, but leaving this field empty or specifying `auto` instructs RethinkDB to determine a cache size automatically. |
| `cloud_provider`                       | no       | Set the cloud provider for the kubernetes cluster.                                                                                                                                               |
| `pod_cidr`                             | yes      | Sets the subnet pool from which the IP for the Pod should be allocated from the CNI ipam plugin. Default is `192.168.0.0/16`.                                                                    |
| `calico_mtu`                           | no       | Set the MTU (maximum transmission unit) size for the Calico plugin.                                                                                                                              |
| `ipip_mtu`                             | no       | Set the IPIP MTU size for the calico IPIP tunnel interface.                                                                                                                                      |
| `azure_ip_count`                       | no       | Set the IP count for azure allocator to allocate IPs per Azure virtual machine.                                                                                                                               |
| `service-cluster-ip-range`             | yes      | Sets the subnet pool from which the IP for Services should be allocated. Default is `10.96.0.0/16`. 
| `nodeport_range`                       | yes      | Set the port range that for Kubernetes services of type NodePort can be exposed in. Default is `32768-35535`.                                                                                    |
| `custom_kube_api_server_flags`         | no       | Set the configuration options for the Kubernetes API server. (dev)                                                                                                                                   |
| `custom_kube_controller_manager_flags` | no       | Set the configuration options for the Kubernetes controller manager. (dev)                                                                                                                            |
| `custom_kubelet_flags`                 | no       | Set the configuration options for Kubelets. (dev)                                                                                                                                                      |
| `custom_kube_scheduler_flags`          | no       | Set the configuration options for the Kubernetes scheduler. (dev)                                                                                                                                      |
| `local_volume_collection_mapping`      | no       | Store data about collections for volumes in UCP's local KV store instead of on the volume labels. This is used for enforcing access control on volumes.                                          |
| `manager_kube_reserved_resources`      | no       | Reserve resources for Docker UCP and Kubernetes components which are running on manager nodes.                                                                                                   |
| `worker_kube_reserved_resources`       | no       | Reserve resources for Docker UCP and Kubernetes components which are running on worker nodes.                                                                                                    |


*dev indicates that the functionality is only for development and testing. Arbitrary Kubernetes configuration parameters are not tested and supported under the Docker Enterprise Software Support Agreement.
