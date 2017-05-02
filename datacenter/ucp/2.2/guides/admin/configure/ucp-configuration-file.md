---
description: Configure UCP deployments.
keywords: docker enterprise edition, ucp, universal control plane, swarm, cluster configuration, deploy
title: UCP configuration file
---

Override the default UCP settings by providing a configuration file when you create 
UCP manager nodes. This is useful for scripted installations.

```bash
$ docker config create --name ...  <ucp.cfg>
```

Specify your configuration settings in a TOML file. For more info, see
[Tom's Obvious, Minimal Language](https://github.com/toml-lang/toml/blob/master/README.md). 

## Example configuration file

Here's an example TOML config file that shows how to configure UCP settings.

```
// ExampleConfig contains an example config to help users understand how to configure UCP.

[[registries]]

# The address used to connect to the DTR instance tied to this UCP cluster.
host_address="example.com:444"

# The DTR instance's OpenID Connect Client ID, as registered with our auth provider.
service_id=""

# The root CA bundle for the DTR instance (if using a custom CA).
ca_bundle="-----BEGIN CERTIFICATE-----\nMIIEyjCCArKgAwIBAgIRAJYDdNEtRX3njQ4JJVCuaScwDQYJKoZIhvcNAQELBQAw\n..."

[scheduling_configuration]

# Allow admins to schedule containers on managers
# Set to true to allow admins to schedule on manager
enable_admin_ucp_scheduling=true

# Allow non-admin users to schedule containers on managers
# Set to true to allow users to schedule on managers
enable_user_ucp_scheduling=true

[tracking_configuration]

# Disable analytics of usage information
# Set to true to disable analytics
disable_usageinfo=false

# Disable analytics of API call information
# Set to true to disable analytics
disable_tracking=false

# Anonymize analytic data
# Set to true to hide your license ID
anonymize_tracking=false

[trust_configuration]

# Require images be signed by content trust
require_content_trust=false

# Specify users or teams which must sign images
require_signature_from=["team1", "team2"]

[log_configuration]

# Specify the protocol to use for remote logging
protocol="tcp"

# Specify a remote syslog server to send UCP controller logs to
# if omitted, controller logs will be sent through the default
# docker daemon logging driver from the ucp-controller container
host="example.com"

# Set the logging level for UCP components - uses syslog levels
level="DEBUG"

[license_configuration]

# Enable attempted automatic license renewal when the license nears expiration
# If disabled, you must manually upload renewed licesnse after expiration.
auto_refresh=true

[cluster_config]

# Configures the port the ucp-controller listens to
controller_port=443

# Configures the port the ucp-swarm-manager listens to
swarm_port=2376

# Configures Swarm scheduler strategy for container scheduling
# This does not affect swarm-mode services
swarm_strategy="spread"

# Configures DNS settings for the UCP components
dns=[]
dns_opt=[]
dns_search=[]

# Turn on specialized debugging endpoints for profiling UCP performance
profiling_enabled=false

# Tune the KV store timeout and snapshot settings
kv_timeout=5000 # milliseconds
kv_snapshot_count=20000

# Specify an optional external LB for default links to services with expose ports in the UI
external_service_lb="example.com"

# Adjust the metrics retention time
metrics_retention_time="24h"

# Set the interval for how frequently managers gather metrics from nodes in the cluster
metrics_scrape_interval="1m"

# Set the interval for how frequently storage metrics are gathered
# this operation can be expensive when large volumes are present
metrics_disk_usage_interval="2h"
```

## Config file and web UI

Admin users can open the UCP web UI, navigate to **Admin Settings**, 
and change UCP settings there. In most cases, the web UI is a front end 
for modifying this config file.

## registries array (required)

An array of tables that specifies the DTR instances that the current UCP instance manages.

| Parameter      | Required | Description                                                                                                                                                                                 |
| -------------- | -------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `host_address` | yes      | The address for connecting to the DTR instance tied to this UCP cluster.                                                                                                                    |
| `service_id`   | yes      | The DTR instance's OpenID Connect Client ID, as registered with the Docker authentication provider.                                                                                         |
| `ca_bundle`    | no       | If you're using a custom certificate authority (CA), the `ca_bundle` setting specifies the root CA bundle for the DTR instance. The value is a string with the contents of a `ca.pem` file. |

## scheduling_configuration table (optional)

Specifies the users who can schedule containers on manager nodes. 

| Parameter                     | Required | Description                                                                                        |
| ----------------------------- | -------- | -------------------------------------------------------------------------------------------------- |
| `enable_admin_ucp_scheduling` | no       | Set to `true` to allow admins to schedule on containers on manager nodes. The default is `false`.  |
| `enable_user_ucp_scheduling`  | no       | Set to `true` to allow non-admin users to schedule containers on managers. The default is `false`. |

## tracking_configuration table (optional)

Specifies the analytics data that UCP collects. 

| Parameter            | Required | Description                                                                             |
| -------------------- | -------- | --------------------------------------------------------------------------------------- |
| `disable_usageinfo`  | no       | Set to `true` to disable analytics of usage information. The default is `false`.        |
| `disable_tracking`   | no       | Set to `true` to disable analytics of API call information. The default is `false`.     |
| `anonymize_tracking` | no       | Anonymize analytic data. Set to `true` to hide your license ID. The default is `false`. |

## trust_configuration table (optional)

Specifies whether DTR images require signing. 

| Parameter                | Required | Description                                                                         |
| ------------------------ | -------- | ----------------------------------------------------------------------------------- |
| `require_content_trust`  | no       | Set to `true` to require images be signed by content trust. The default is `false`. |
| `require_signature_from` | no       | A string array that specifies users or teams which must sign images.                |

## log_configuration table (optional)

Configures the logging options for UCP components. 

| Parameter  | Required | Description                                                                                                                                                                                     |
| ---------- | -------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `protocol` | no       | The protocol to use for remote logging. Values are those supported by the Go [Dial function](https://golang.org/pkg/net/#Dial). The default is `tcp`.                                                                                                                               |
| `host`     | no       | Specifies a remote syslog server to send UCP controller logs to. If omitted, controller logs are sent through the default docker daemon logging driver from the `ucp-controller` container.     |
| `level`    | no       | The logging level for UCP components. Values are [syslog priority  levels](https://linux.die.net/man/5/syslog.conf): `debug`, `info`, `notice`, `warning`, `err`, `crit`, `alert`, and `emerg`. |

## license_configuration table (optional)

Specifies whether the your UCP license is automatically renewed.   

| Parameter      | Required | Description                                                                                                                                                                                    |
| -------------- | -------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `auto_refresh` | no       | Set to `true` to enable attempted automatic license renewal when the license nears expiration. If disabled, you must manually upload renewed licesnse after expiration. The default is `true`. |

## cluster_config table (required)

Configures the swarm cluster that the current UCP instance manages.

The `dns`, `dns_opt`, and `dns_search` settings configure the DNS settings for UCP 
components. Assigning these values overrides the settings in a container's 
`/etc/resolv.conf` file. For more info, see
[Configure container DNS](/engine/userguide/networking/default_network/configure-dns/).

| Parameter                     | Required | Description                                                                                                                                    |
| ----------------------------- | -------- | ---------------------------------------------------------------------------------------------------------------------------------------------- |
| `controller_port`             | yes      | Configures the port that the `ucp-controller` listens to. The default is `443`.                                                                |
| `swarm_port`                  | yes      | Configures the port that the `ucp-swarm-manager` listens to. The default is `2376`.                                                            |
| `swarm_strategy`              | no       | Configures placement strategy for container scheduling. This doesn't affect swarm-mode services. Values are `spread`, `binpack`, and `random`. |
| `dns`                         | yes      | Array of IP addresses to add as nameservers.                                                                                                   |
| `dns_opt`                     | yes      | Array of options used by DNS resolvers.                                                                                                        |
| `dns_search`                  | yes      | Array of domain names to search when a bare unqualified hostname is used inside of a container.                                                |
| `profiling_enabled`           | no       | Set to `true` to enable specialized debugging endpoints for profiling UCP performance. The default is `false`.                                 |
| `kv_timeout`                  | no       | Sets the key-value store timeout setting, in milliseconds. The default is `5000`.                                                              |
| `kv_snapshot_count`           | no       | Sets the key-value store snapshot count setting. The default is `20000`.                                                                       |
| `external_service_lb`         | no       | Specifies an optional external load balancer for default links to services with exposed ports in the web UI.                                   |
| `metrics_retention_time`      | no       | Adjusts the metrics retention time.                                                                                                            |
| `metrics_scrape_interval`     | no       | Sets the interval for how frequently managers gather metrics from nodes in the cluster.                                                        |
| `metrics_disk_usage_interval` | no       | Sets the interval for how frequently storage metrics are gathered. This operation can be expensive when large volumes are present.             |