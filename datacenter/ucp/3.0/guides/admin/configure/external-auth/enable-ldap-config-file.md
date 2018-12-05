---
title: Integrate with LDAP by using a configuration file
description: Set up LDAP authentication by using a configuration file.
keywords: UCP, LDAP, config
---

Docker UCP integrates with LDAP directory services, so that you can manage
users and groups from your organization's directory and automatically
propagate this information to UCP and DTR. You can set up your swarm's LDAP
configuration by using the UCP web UI, or you can use a
[UCP configuration file](../ucp-configuration-file.md).

To see an example TOML config file that shows how to configure UCP settings,
run UCP with the `example-config` option.
[Learn about UCP configuration files](../ucp-configuration-file.md).

```bash
$ docker container run --rm {{ page.ucp_org }}/{{ page.ucp_repo }}:{{ page.ucp_version }} example-config
```

## Set up LDAP by using a configuration file

1.  Use the following command to extract the name of the currently active
    configuration from the `ucp-agent` service.

    {% raw %}
    ```bash
    $ CURRENT_CONFIG_NAME=$(docker service inspect --format '{{ range $config := .Spec.TaskTemplate.ContainerSpec.Configs }}{{ $config.ConfigName }}{{ "\n" }}{{ end }}' ucp-agent | grep 'com.docker.ucp.config-')
    ```
    {% endraw %}

2.  Get the current configuration and save it to a TOML file.

    {% raw %}
    ```bash
    $ docker config inspect --format '{{ printf "%s" .Spec.Data }}' $CURRENT_CONFIG_NAME > config.toml
    ```
    {% endraw %}

3.  Use the output of the `example-config` command as a guide to edit your
    `config.toml` file. Under the `[auth]` sections, set `backend = "ldap"`
    and `[auth.ldap]` to configure LDAP integration the way you want.

4.  Once you've finished editing your `config.toml` file, create a new Docker
    Config object by using the following command.

    ```bash
    $ NEW_CONFIG_NAME="com.docker.ucp.config-$(( $(cut -d '-' -f 2 <<< "$CURRENT_CONFIG_NAME") + 1 ))"
    docker config create $NEW_CONFIG_NAME config.toml
    ```

5.  Update the `ucp-agent` service to remove the reference to the old config
    and add a reference to the new config.

    ```bash
    $ docker service update --config-rm "$CURRENT_CONFIG_NAME" --config-add "source=${NEW_CONFIG_NAME},target=/etc/ucp/ucp.toml" ucp-agent
    ```

6.  Wait a few moments for the `ucp-agent` service tasks to update across
    your swarm. If you set `jit_user_provisioning = true` in the LDAP
    configuration, users matching any of your specified search queries will
    have their accounts created when they log in with their username and LDAP
    password.

## Where to go next

-  [Create and manage users](../../../access-control/create-and-manage-users.md)
-  [Create and manage teams](../../../access-control/create-and-manage-teams.md)