---
title: Docker CE Edge documentation
description: Information about current Docker Edge releases
keywords: engine, edge, installation
current_edge: 17.05
---

The current Docker CE Edge release is {{ page.current_edge }}. The Docker CE
Edge channel provides monthly releases which allow you to try new features of
Docker and verify bug fixes quickly. Edge releases are only supported for one
month, and a given Edge release will not receive any updates once a new edge
release is available.

Stable releases are not published to the Edge channel, so Linux repository users
still need to subscribe to the stable channel as well.

Commercial support is not available for Docker CE.

For information about all Docker release channels and expectations about
support, see [Docker channels](/engine/installation/#docker-channels).

This page lists features that are only available in Docker CE Edge releases.
Where applicable, the API and CLI reference documentation has been updated to
reflect these features, but **full documentation for a given feature may not be
available until a Docker CE Stable release incorporates the feature**.

### Docker CE Edge new features

<ul class="nav nav-tabs">
  <li class="active"><a data-toggle="tab" data-target="#1705">17.05</a></li>
  <li><a data-toggle="tab" data-target="#1704">17.04</a></li>
</ul>

<div markdown="1" class="tab-content">

<div markdown="1" id="1705" class="tab-pane fade in active">
#### Docker CE Edge 17.05

The following major features and changes are included in Docker CE Edge 17.05.
**Docker CE Edge 17.05 also includes the features from Docker CE Edge 17.04.**
Continue reading, or go straight to [API and CLI](#api-and-cli),
[Builder](#builder), [Daemon](#daemon), [Dockerfile](#dockerfile),
[Logging](#logging), [Networking](#networking),
[Operating system support](#operating-system-support), [Runtime](#runtime),
[Security](#security), [Services](#services), or [Stacks](#stacks).

[Read the full release notes](https://github.com/moby/moby/releases/tag/v17.05.0-ce){: target="_blank" class="_" }

##### API and CLI

- Add [version 1.29](/engine/api/v1.29/) of the Docker API.
- Add `--mount` flag to `docker run` and `docker create`
  {% include github-pr.md pr="32251" %}
- Add the ability to specify `--type=secret` in `docker inspect`
  {% include github-pr.md pr="32124" %}
- Add `--format` and `--filter` options for `docker secret ls`
  {% include github-pr.md pr="30810" %} and {% include github-pr.md pr="31552" %}
- Add the ability to only show local or swarm networks in `docker network ls`
  {% include github-pr.md pr="31529" %}
- Add the ability to update the `--cpus` value for a running container
  {% include github-pr.md pr="31148" %}
- Add label filtering to all `prune` commands
  {% include github-pr.md pr="30740" %}
- Add the ability to remove multiple stacks with a single `docker stack rm`
  invocation {% include github-pr.md pr="32110" %}
* Improve `docker version --format` option when the client has downgraded the
  API version {% include github-pr.md pr="31022" %}
* You are now prompted when using an encrypted client certificate to connect to
  a Docker daemon {% include github-pr.md pr="31364" %}

##### Builder

- Add support for multi-stage builds {% include github-pr.md pr="31257" %}
  (see [Use multi-stage builds](/engine/userguide/eng-image/multistage-build.md))
- Add support for named build stages when using multi-stage builds
  {% include github-pr.md pr="32063" %}
- Add support for using `ARG` variables in the `FROM` instruction in the
  Dockerfile {% include github-pr.md pr="31352" %}
- `docker build` now reads the Dockerfile from `STDIN` when using the `-f` flag
  {% include github-pr.md pr="31236" %}
- Default build arguments are no longer included in the image history
  {% include github-pr.md pr="31584" %}

##### Daemon

- The `--api-cors-header` is no longer ignored if `--api-enable-cors` is not set
  {% include github-pr.md pr="32174" %}
- The `--graph` flag is deprecated in favor or `--data-root`
  {% include github-pr.md pr="28696" %}
- Document deprecation of the `--api-enable-cors` daemon flag. This flag was
  marked as deprecated in Docker 1.6.0 but was inadvertently omitted from the
  list of deprecated features. {% include github-pr.md pr="32352" %}

##### Logging

- You can now show the logs for an individual service task in the
  `docker service logs` command or via the `/tasks/{id}/logs` REST endpoint
  {% include github-pr.md pr="32015" %}
- Add support for logging driver plugins {% include github-pr.md pr="28403" %}
- Add the ability to set `--log-opt env-regex` to match an environment variable
  using a regular expression {% include github-pr.md pr="27565" %}
- `docker service logs` is no longer experimental.
  {% include github-pr.md pr="32462" %}

##### Networking

- Add the ability to replace, and customize the ingress network
  {% include github-pr.md pr="31714" %}
- Files are now written to the correct directory if a custom data-root is set
  {% include github-pr.md pr="32505" %}

##### Operating system support

- Ubuntu 12.04 (Precise Pangolin) is no longer supported.
  {% include github-pr.md pr="32520" %}
- Ubuntu 17.04 (Zesty Zapus) is now supported.
  {% include github-pr.md pr="32435" %}

##### Runtime

- The health probe is now consistently stopped when a container exits
  {% include github-pr.md pr="32274" %}
- You can now specify "grace periods" on healthchecks using `--start-period` and `--health-start-period`
  to support services and containers with an initial startup delay {% include github-pr.md pr="28938" %}

##### Security

- Add the ability to set SELinux type or MCS labels when using
  `--ipc=container:` or `--ipc=host ` {% include github-pr.md pr="30652" %}

##### Services

- You can now show the logs for an individual service task in the
  `docker service logs` command or via the `/tasks/{id}/logs` REST endpoint
  {% include github-pr.md pr="32015" %}
- You can now specify the update or rollback order for services using
  `--update-order` and `--rollback-order` {% include github-pr.md pr="30261" %}
- Add support for synchronous service creation and update {% include github-pr.md pr="31144" %}
- You can now specify "grace periods" on healthchecks using `--start-period` and `--health-start-period`
  to support services and containers with an initial startup delay
  {% include github-pr.md pr="28938" %}
- `docker service create` now omits fields that are not shown by the user, where
  possible. In addition, `docker service inspect` now shows default values that
  were not specified by the user. {% include github-pr.md pr="32284" %}
- `docker service logs` is no longer experimental.
  {% include github-pr.md pr="32462" %}
- Add support for Credential Spec and SELinux to services
  {% include github-pr.md pr="32339" %}
- You can override the entrypoint using the `--entrypoint` flag with
  `docker service create` or `docker service update`
  {% include github-pr.md pr="29228" %}
- Add `--network-add` and `--network-rm` flags to `docker service update`
  {% include github-pr.md pr="32062" %}
- Add the ability to specify the credential specification to
  `docker service create` and `docker service update`
  {% include github-pr.md pr="32339" %}
- Add the ability to filter by service replication mode when listing services
  {% include github-pr.md pr="31538" %}
- Add `--format` option to `docker node ls`
  {% include github-pr.md pr="30424" %}
- Tasks are no longer re-deployed unnecessarily when environment-variables are
  used {% include github-pr.md pr="32364" %}

##### Stacks

- Add the ability to update stacks to prune services that are no longer defined
  in the stack file {% include github-pr.md pr="31302" %}
- `docker stack deploy` now supports `endpoint_mode` when deploying from a stack
  file {% include github-pr.md pr="32333" %}
</div> <!-- 17.05 -->

<div markdown="1" id="1704" class="tab-pane fade">

#### Docker CE Edge 17.04

The following major features and changes are included in Docker CE Edge 17.04.

[Read the full release notes](https://github.com/moby/moby/releases/tag/v17.04.0-ce){: target="_blank" class="_" }

##### API and CLI

- Add [version 1.28](/engine/api/v1.28/) of the Docker API.

- Add `--device-cgroup-rule` flag to give containers access to devices that appear
  after the container is started. {% include github-pr.md pr="22563" %}

- Allow swarm nodes to join with `--availability=drain` to prevent them from
  taking non-manager workloads. {% include github-pr.md pr="24993" %}

- Add `publish` and `expose` filters to `docker ps`, so that containers can be
  filtered by port or port range for TCP or UDP protocols {% include github-pr.md pr="27557" %}

- Add `--no-trunc` and `--format` flags to the `docker service ls` command, and
  as well as the ability to specify the default format for `docker service ls`
  using the `ServicesFormat` option to the Docker CLI. Also add a
  `docker stack services` command. {% include github-pr.md pr="28199" %}

- Add ability to filter plugins by whether they are enabled or disabled in
  `docker plugin ls` output. {% include github-pr.md pr="28627" %}

- Add `mode` option to `--log-opts` flag for both `docker` and `dockerd`. If set
  to `non-blocking`, and the log buffer fills up, log messages will be lost, but
  the container will not block. The `max-buffer-size` option controls the
  maximum size of the ring buffer. Defaults to `blocking`, which will cause the
  container to block if messages cannot be logged. See
  [Options for all drivers](/engine/admin/logging/overview.md#options-for-all-drivers).
  {% include github-pr.md pr="28762" %}

- It is no longer possible to inadvertently pull images on an architecture where
  they will not run. {% include github-pr.md pr="29001" %}

- It is now possible to create AWS log groups when using the AWS logging driver.
  See [`awslogs-create-group`](engine/admin/logging/awslogs.md#awslogs-create-group).
  {% include github-pr.md pr="29504" %}

- Add the ability to filter `docker network ls` output by creation time, using
  the `{% raw %}{{CreatedAt}}{% endraw %}` format specifier.
  {% include github-pr.md pr="29900" %}

- Named but untagged images are now removed if you run `docker image prune` if
  `--dangling-only` is set to `true`. {% include github-pr.md pr="30330" %}

- Add `--add-host` flag to `docker build`, which will add entries to the
  `/etc/hosts` file of a container created from that image. The `/etc/hosts`
  file is not saved within the image itself. {% include github-pr.md pr="30383" %}

- Prevent `docker network ls` from pulling all the endpoints, to reduce
  impact on the network. {% include github-pr.md pr="30673" %}

- Windows-specific commands and options no longer show in command help text on
  non-Windows clients. {% include github-pr.md pr="30780" %}

- When you specify an IP address when running `docker network connect`, the
  IP address is now checked for validity. {% include github-pr.md pr="30807" %}

- Add the ability to customize bind-mount consistency to be more appropriate
  for some platforms and workloads. Options are `consistent` (the default),
  `cached`, or `delegated`. {% include github-pr.md pr="31047" %}

##### Daemon

- Docker Daemon logging settings no longer affect the `docker build` command.
  {% include github-pr.md pr="29552" %}

- Add a `registry-mirrors` configuration option for the Docker daemon, which
  replaces the daemon's registry mirrors with a new set of registry mirrors.
  {% include github-pr.md pr="29650" %}

- Add the ability to specify the default shared memory size for the Docker
  daemon, using the `--default-shm-size` or the `default-shm-size` key in
  `daemon.json`. {% include github-pr.md pr="29692" %}

- Add a `no-new-privileges` configuration option for the Docker daemon, which
  prevents unprivileged containers from gaining new privileges.
  {% include github-pr.md pr="29984" %}

- If a Docker client communicates with an older daemon and attempts to perform
  an operation not supported by the daemon, an error is printed, which shows
  the API versions of both the client and daemon.
  {% include github-pr.md pr="30187" %}

- The Docker daemon no longer depends upon `sqlite`. This change means that it
  is not possible to upgrade the Docker daemon from version 1.9 directly to the
  latest version. It is recommended to upgrade from one major version to the
  next, in sequence. {% include github-pr.md pr="30208" %}

##### Dockerfile

- Using the pattern `**/` in a Dockerfile now (correctly) behaves the same as
  `**`. {% include github-pr.md pr="29043" %}

- Time values less than 1 second are no longer allowed in health-check options
  in the Dockerfile. {% include github-pr.md pr="31177" %}

##### Services

- When a service is updated with both `--secret-add` and `--secret-rm` in the
  same operation, the order of operations is now changed so that the
  `--secret-rm` always occurs first. {% include github-pr.md pr="29802" %}

- Add the ability to create or update a service to be read-only using the
  `--read-only` flag. {% include github-pr.md pr="30162" %}

- Docker now updates swarm nodes if the swarm configuration is updated.
  {% include github-pr.md pr="30259" %}

- Add topology-aware placement preferences for Swarm services. This feature
  allows services to be balanced over nodes based on a particular user-defined
  property, such as which datacenter or rack they are located in.
  See [Control service scale and placement](/engine/swarm/services.md#control-service-scale-and-placement).
  {% include github-pr.md pr="30725" %}

- Add the ability to customize the stop signal which will be sent to nodes, when
  creating or updating a service. {% include github-pr.md pr="30754" %}

- Add the ability to address a secret by name or prefix, as well as ID, when
  updating it. {% include github-pr.md pr="30856" %}

- Add the ability to roll back to a previous version of a service if an
  updated service fails to deploy. Several flags are available at service
  creation or update,to control the rollback action, failure threshold,
  monitoring delay, rollback delay, and parallelism.
  {% include github-pr.md pr="31108" %}

- Add the ability to specify the stream when using the Docker service logs API.
  {% include github-pr.md pr="31313" %}

- Add `--tail` and `--since` flags to `docker service logs` command, to filter
  the logs by time or to show the tail of the logs and show new content as it
  is logged. {% include github-pr.md pr="31500" %}

- Add a `--verbose` flag to the `docker inspect` command. For swarm networks,
  this flag shows all nodes and services attached to the network.
  {% include github-pr.md pr="31710" %}

##### Stacks

- Compose file version 3.2 is now supported. This includes support for different
  types of endpoints and expands the options you can use when specifying mounts.
  {% include github-pr.md pr="31795" %}

</div> <!-- 17.04 -->
</div> <!-- tab-content -->
