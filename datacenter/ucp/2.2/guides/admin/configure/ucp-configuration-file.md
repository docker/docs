---
title: UCP configuration file
description: Configure UCP deployments.
keywords: docker enterprise edition, ucp, universal control plane, swarm, configuration, deploy
---

Override the default UCP settings by providing a configuration file when you
create UCP manager nodes. This is useful for scripted installations.

## UCP configuration file

The `ucp-agent` service uses a configuration file to set up UCP.
You can use the configuration file in different ways to set up your UCP
swarms.

- Install one swarm and use the UCP web UI to configure it as desired,
  extract the configuration file, edit it as needed, and use the edited
  config file to make copies to multiple other swarms.
- Install a UCP swarm, extract and edit the configuration file, and use the
  CLI to apply the new configuration to the same swarm.
- Run the `example-config` command, edit the example configuration file, and
  apply the file at install time or after installation.

Specify your configuration settings in a TOML file.
[Learn about Tom's Obvious, Minimal Language](https://github.com/toml-lang/toml/blob/master/README.md).

The configuration has a versioned naming convention, with a trailing decimal
number that increases with each version, like `com.docker.ucp.config-1`. The
`ucp-agent` service maps the configuration to the file at `/etc/ucp/ucp.toml`.

## Inspect and modify existing configuration

Use the `docker config inspect` command to view the current settings and emit
them to a file.

{% raw %}
```bash
# CURRENT_CONFIG_NAME will be the name of the currently active UCP configuration
CURRENT_CONFIG_NAME=$(docker service inspect ucp-agent --format '{{range .Spec.TaskTemplate.ContainerSpec.Configs}}{{if eq "/etc/ucp/ucp.toml" .File.Name}}{{.ConfigName}}{{end}}{{end}}')
# Collect the current config with `docker config inspect`
docker config inspect --format '{{ printf "%s" .Spec.Data }}' $CURRENT_CONFIG_NAME > ucp-config.toml
```
{% endraw %}

Edit the file, then use the `docker config create` and `docker service update`
commands to create and apply the configuration from the file.


```bash
# NEXT_CONFIG_NAME will be the name of the new UCP configuration
NEXT_CONFIG_NAME=${CURRENT_CONFIG_NAME%%-*}-$((${CURRENT_CONFIG_NAME##*-}+1))
# Create the new swarm configuration from the file ucp-config.toml
docker config create $NEXT_CONFIG_NAME  ucp-config.toml
# Use the `docker service update` command to remove the current configuration
# and apply the new configuration to the `ucp-agent` service.
docker service update --config-rm $CURRENT_CONFIG_NAME --config-add source=$NEXT_CONFIG_NAME,target=/etc/ucp/ucp.toml ucp-agent
```

## Example configuration file

You can see an example TOML config file that shows how to configure UCP
settings. From the command line, run UCP with the `example-config` option:

```bash
$ docker container run --rm {{ page.ucp_org }}/{{ page.ucp_repo }}:{{ page.ucp_version }} example-config
```


## Configuration file and web UI

Admin users can open the UCP web UI, navigate to **Admin Settings**,
and change UCP settings there. In most cases, the web UI is a front end
for modifying this config file.

## auth table

| Parameter               | Required | Description                                                                                                                                                                        |
|:------------------------|:---------|:-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `backend`               | no       | The name of the authorization backend to use, either `managed` or `ldap`. The default is `managed`.                                                                                |
| `default_new_user_role` | no       | The role that new users get for their private collections. Values are `admin`, `viewonly`, `scheduler`, `restrictedcontrol`, or `fullcontrol`. The default is `restrictedcontrol`. |


## auth.sessions

| Parameter                   | Required | Description                                                                                                                                                                                                                                                                             |
|:----------------------------|:---------|:----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `lifetime_minutes`          | no       | The initial session lifetime, in minutes. The default is 4320, which is 72 hours.                                                                                                                                                                                                       |
| `renewal_threshold_minutes` | no       | The length of time, in minutes, before the expiration of a session where, if used, a session will be extended by the current configured lifetime from then. A zero value disables session extension. The default is 1440, which is 24 hours.                                            |
| `per_user_limit`            | no       | The maximum number of sessions that a user can have active simultaneously. If creating a new session would put a user over this limit, the least recently used session will be deleted. A value of zero disables limiting the number of sessions that users may have. The default is 5. |

## auth.ldap (optional)

| Parameter               | Required | Description                                                                                                                                                                      |
|:------------------------|:---------|:---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `server_url`            | no       | The URL of the LDAP server.                                                                                                                                                      |
| `no_simple_pagination`  | no       | Set to `true` if the LDAP server doesn't support the Simple Paged Results control extension (RFC 2696). The default is `false`.                                                  |
| `start_tls`             | no       | Set to `true` to use StartTLS to secure the connection to the server, ignored if the server URL scheme is 'ldaps://'. The default is `false`.                                    |
| `root_certs`            | no       | A root certificate PEM bundle to use when establishing a TLS connection to the server.                                                                                           |
| `tls_skip_verify`       | no       | Set to `true` to skip verifying the server's certificate when establishing a TLS connection, which isn't recommended unless testing on a secure network. The default is `false`. |
| `reader_dn`             | no       | The distinguished name the system uses to bind to the LDAP server when performing searches.                                                                                      |
| `reader_password`       | no       | The password that the system uses to bind to the LDAP server when performing searches.                                                                                           |
| `sync_schedule`         | no       | The scheduled time for automatic LDAP sync jobs, in CRON format. Needs to have the seconds field set to zero. The default is @hourly if empty or omitted.                        |
| `jit_user_provisioning` | no       | Whether to only create user accounts upon first login (recommended). The default is `true`.                                                                                      |


## auth.ldap.additional_domains array (optional)

A list of additional LDAP domains and corresponding server configs from which
to sync users and team members. This is an advanced feature which most
environments don't need.

| Parameter              | Required | Description                                                                                                                                                                                                                                                                 |
|:-----------------------|:---------|:----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `domain`               | no       | The root domain component of this server, for example, `dc=example,dc=com`. A longest-suffix match of the base DN for LDAP searches is used to select which LDAP server to use for search requests. If no matching domain is found, the default LDAP server config is used. |
| `server_url`           | no       | The URL of the LDAP server for the current additional domain.                                                                                                                                                                                                               |
| `no_simple_pagination` | no       | Set to true if the LDAP server for this additional domain does not support the Simple Paged Results control extension (RFC 2696). The default is `false`.                                                                                                                   |
| `server_url`           | no       | The URL of the LDAP server.                                                                                                                                                                                                                                                 |
| `start_tls`            | no       | Whether to use StartTLS to secure the connection to the server, ignored if the server URL scheme is 'ldaps://'.                                                                                                                                                             |
| `root_certs`           | no       | A root certificate PEM bundle to use when establishing a TLS connection to the server for the current additional domain.                                                                                                                                                    |
| `tls_skip_verify`      | no       | Whether to skip verifying the additional domain server's certificate when establishing a TLS connection, not recommended unless testing on a secure network. The default is `true`.                                                                                         |
| `reader_dn`            | no       | The distinguished name the system uses to bind to the LDAP server when performing searches under the additional domain.                                                                                                                                                     |
| `reader_password`      | no       | The password that the system uses to bind to the LDAP server when performing searches under the additional domain.                                                                                                                                                          |

## auth.ldap.user_search_configs array (optional)

Settings for syncing users.

| Parameter                 | Required | Description                                                                                                                                                                                                                                                                              |
|:--------------------------|:---------|:-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `base_dn`                 | no       | The distinguished name of the element from which the LDAP server will search for users, for example, `ou=people,dc=example,dc=com`.                                                                                                                                                      |
| `scope_subtree`           | no       | Set to `true` to search for users in the entire subtree of the base DN. Set to `false` to search only one level under the base DN. The default is `false`.                                                                                                                               |
| `username_attr`           | no       | The name of the attribute of the LDAP user element which should be selected as the username. The default is `uid`.                                                                                                                                                                       |
| `full_name_attr`          | no       | The name of the attribute of the LDAP user element which should be selected as the full name of the user. The default is `cn`.                                                                                                                                                           |
| `filter`                  | no       | The LDAP search filter used to select user elements, for example, `(&(objectClass=person)(objectClass=user))`. May be left blank.                                                                                                                                                        |
| `match_group`             | no       | Whether to additionally filter users to those who are direct members of a group. The default is `true`.                                                                                                                                                                                  |
| `match_group_dn`          | no       | The distinguished name of the LDAP group, for example, `cn=ddc-users,ou=groups,dc=example,dc=com`. Required if `matchGroup` is `true`.                                                                                                                                                   |
| `match_group_member_attr` | no       | The name of the LDAP group entry attribute which corresponds to distinguished names of members. Required if `matchGroup` is `true`. The default is `member`.                                                                                                                             |
| `match_group_iterate`     | no       | Set to `true` to get all of the user attributes by iterating through the group members and performing a lookup for each one separately. Use this instead of searching users first, then applying the group selection filter. Ignored if `matchGroup` is `false`. The default is `false`. |

## auth.ldap.admin_sync_opts (optional)

Settings for syncing system admininistrator users.

| Parameter              | Required | Description                                                                                                                                                                                               |
|:-----------------------|:---------|:----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `enable_sync`          | no       | Set to `true` to enable syncing admins. If `false`, all other fields in this table are ignored. The default is `true`.                                                                                    |
| `select_group_members` | no       | Set to `true` to sync using a group DN and member attribute selection. Set to `false` to use a search filter. The default is `true`.                                                                      |
| `group_dn`             | no       | The distinguished name of the LDAP group, for example, `cn=ddc-admins,ou=groups,dc=example,dc=com`. Required if `select_group_members` is `true`.                                                         |
| `group_member_attr`    | no       | The name of the LDAP group entry attribute which corresponds to distinguished names of members. Required if `select_group_members` is `true`. The default is `member`.                                    |
| `search_base_dn`       | no       | The distinguished name of the element from which the LDAP server will search for users, for example, `ou=people,dc=example,dc=com`. Required if `select_group_members` is `false`.                        |
| `search_scope_subtree` | no       | Set to `true` to search for users in the entire subtree of the base DN. Set to `false` to search only one level under the base DN. The default is `false`. Required if `select_group_members` is `false`. |
| `search_filter`        | no       | The LDAP search filter used to select users if `select_group_members` is `false`, for example, `(memberOf=cn=ddc-admins,ou=groups,dc=example,dc=com)`. May be left blank.                                 |


## registries array (required)

An array of tables that specifies the DTR instances that the current UCP instance manages.

| Parameter      | Required | Description                                                                                                                                                                                 |
|:---------------|:---------|:--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `host_address` | yes      | The address for connecting to the DTR instance tied to this UCP cluster.                                                                                                                    |
| `service_id`   | yes      | The DTR instance's OpenID Connect Client ID, as registered with the Docker authentication provider.                                                                                         |
| `ca_bundle`    | no       | If you're using a custom certificate authority (CA), the `ca_bundle` setting specifies the root CA bundle for the DTR instance. The value is a string with the contents of a `ca.pem` file. |

## scheduling_configuration table (optional)

Specifies the users who can schedule containers on manager nodes.

| Parameter                     | Required | Description                                                                                        |
|:------------------------------|:---------|:---------------------------------------------------------------------------------------------------|
| `enable_admin_ucp_scheduling` | no       | Set to `true` to allow admins to schedule on containers on manager nodes. The default is `false`.  |
| `enable_user_ucp_scheduling`  | no       | Set to `true` to allow non-admin users to schedule containers on managers. The default is `false`. |

## tracking_configuration table (optional)

Specifies the analytics data that UCP collects.

| Parameter            | Required | Description                                                                             |
|:---------------------|:---------|:----------------------------------------------------------------------------------------|
| `disable_usageinfo`  | no       | Set to `true` to disable analytics of usage information. The default is `false`.        |
| `disable_tracking`   | no       | Set to `true` to disable analytics of API call information. The default is `false`.     |
| `anonymize_tracking` | no       | Anonymize analytic data. Set to `true` to hide your license ID. The default is `false`. |

## trust_configuration table (optional)

Specifies whether DTR images require signing.

| Parameter                | Required | Description                                                                         |
|:-------------------------|:---------|:------------------------------------------------------------------------------------|
| `require_content_trust`  | no       | Set to `true` to require images be signed by content trust. The default is `false`. |
| `require_signature_from` | no       | A string array that specifies users or teams which must sign images.                |

## log_configuration table (optional)

Configures the logging options for UCP components.

| Parameter  | Required | Description                                                                                                                                                                                     |
|:-----------|:---------|:------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `protocol` | no       | The protocol to use for remote logging. Values are `tcp` and `udp`. The default is `tcp`.                                                                                                       |
| `host`     | no       | Specifies a remote syslog server to send UCP controller logs to. If omitted, controller logs are sent through the default docker daemon logging driver from the `ucp-controller` container.     |
| `level`    | no       | The logging level for UCP components. Values are [syslog priority  levels](https://linux.die.net/man/5/syslog.conf): `debug`, `info`, `notice`, `warning`, `err`, `crit`, `alert`, and `emerg`. |

## license_configuration table (optional)

Specifies whether the your UCP license is automatically renewed.   

| Parameter      | Required | Description                                                                                                                                                                                   |
|:---------------|:---------|:----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `auto_refresh` | no       | Set to `true` to enable attempted automatic license renewal when the license nears expiration. If disabled, you must manually upload renewed license after expiration. The default is `true`. |

## cluster_config table (required)

Configures the swarm cluster that the current UCP instance manages.

The `dns`, `dns_opt`, and `dns_search` settings configure the DNS settings for UCP
components. Assigning these values overrides the settings in a container's
`/etc/resolv.conf` file. For more info, see
[Configure container DNS](/engine/userguide/networking/default_network/configure-dns/).

| Parameter                         | Required | Description                                                                                                                                                                                                           |
|:----------------------------------|:---------|:----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `controller_port`                 | yes      | Configures the port that the `ucp-controller` listens to. The default is `443`.                                                                                                                                       |
| `swarm_port`                      | yes      | Configures the port that the `ucp-swarm-manager` listens to. The default is `2376`.                                                                                                                                   |
| `swarm_strategy`                  | no       | Configures placement strategy for container scheduling. This doesn't affect swarm-mode services. Values are `spread`, `binpack`, and `random`.                                                                        |
| `dns`                             | yes      | Array of IP addresses to add as nameservers.                                                                                                                                                                          |
| `dns_opt`                         | yes      | Array of options used by DNS resolvers.                                                                                                                                                                               |
| `dns_search`                      | yes      | Array of domain names to search when a bare unqualified hostname is used inside of a container.                                                                                                                       |
| `profiling_enabled`               | no       | Set to `true` to enable specialized debugging endpoints for profiling UCP performance. The default is `false`.                                                                                                        |
| `kv_timeout`                      | no       | Sets the key-value store timeout setting, in milliseconds. The default is `5000`.                                                                                                                                     |
| `kv_snapshot_count`               | no       | Sets the key-value store snapshot count setting. The default is `20000`.                                                                                                                                              |
| `external_service_lb`             | no       | Specifies an optional external load balancer for default links to services with exposed ports in the web UI.                                                                                                          |
| `metrics_retention_time`          | no       | Adjusts the metrics retention time. Units can be `s/m/h` (`12h` for rexample).                                                                                                                                        |
| `metrics_scrape_interval`         | no       | Sets the interval for how frequently managers gather metrics from nodes in the cluster. Units can be `s/m/h` (`12h` for rexample).                                                                                    |
| `metrics_disk_usage_interval`     | no       | Sets the interval for how frequently storage metrics are gathered. This operation can be expensive when large volumes are present. Units can be `s/m/h` (`12h` for rexample).                                         |
| `rethinkdb_cache_size`            | no       | Set the size of the cache used by UCP's RethinkDB servers. The default is 512MB, but leaving this field empty or specifying the special value "auto" will instruct RethinkDB to determine a cache size automatically. |
| `min_tls_version`                 | no       | Set the minimum TLS version for the controller to serve. Valid options are tlsv1, tlsv1.0, tlsv1.1, and tlsv1.2.                                                                                                      |
| `local_volume_collection_mapping` | no       | Store data about collections for volumes in UCP's local KV store instead of on the volume labels. This is used for enforcing access control on volumes.                                                               |
