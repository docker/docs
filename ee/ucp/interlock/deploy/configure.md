---
title: Configure the layer 7 routing service
description: Learn about Interlock, an application routing and load balancing system
  for Docker Swarm.
keywords: ucp, interlock, load balancing
ui_tabs:
- version: ucp-3.0
  orhigher: false
---

{% if include.version=="ucp-3.0" %}

[When enabling the layer 7 routing solution](index.md) from the UCP web UI,
you can configure the ports for incoming traffic. If you want to further
customize the layer 7 routing solution, you can do it by updating the
`ucp-interlock` service with a new Docker configuration object.

Here's how it works:

1. Find out what configuration is currently being used for the `ucp-interlock`
service and save it to a file:

{% raw %}
```bash
CURRENT_CONFIG_NAME=$(docker service inspect --format '{{ (index .Spec.TaskTemplate.ContainerSpec.Configs 0).ConfigName }}' ucp-interlock)
docker config inspect --format '{{ printf "%s" .Spec.Data }}' $CURRENT_CONFIG_NAME > config.toml
```
{% endraw %}

2. Make the necessary changes to the `config.toml` file.
[Learn about the configuration options available](configuration-reference.md).
3. Create a new Docker configuration object from the file you've edited:

```
NEW_CONFIG_NAME="com.docker.ucp.interlock.conf-$(( $(cut -d '-' -f 2 <<< "$CURRENT_CONFIG_NAME") + 1 ))"
docker config create $NEW_CONFIG_NAME config.toml
```

3. Update the `ucp-interlock` service to start using the new configuration:

```
docker service update \
  --config-rm $CURRENT_CONFIG_NAME \
  --config-add source=$NEW_CONFIG_NAME,target=/config.toml \
  ucp-interlock
```

By default the `ucp-interlock` service is configured to pause if something
goes wrong with the configuration update. The service won't do any updates
without manual intervention.

If you want the service to automatically rollback to a previous stable
configuration, you can update the service with:

```
docker service update \
  --update-failure-action rollback \
  ucp-interlock
```

{% endif %}
