---
title: Using lifecycle hooks with Compose
linkTitle: Use lifecycle hooks
weight: 20
description: Learn how to use Docker Compose lifecycle hooks like post_start and pre_stop to customize container behavior.
keywords: docker compose lifecycle hooks, post_start, pre_stop, docker compose entrypoint, docker container stop hooks, compose hook commands
---

{{< summary-bar feature_name="Compose lifecycle hooks" >}}

## Services lifecycle hooks

When Docker Compose runs a container, it uses two elements, 
[ENTRYPOINT and COMMAND](/manuals/engine/containers/run.md#default-command-and-options), 
to manage what happens when the container starts and stops.

However, it can sometimes be easier to handle these tasks separately with lifecycle hooks - 
commands that run right after the container starts or just before it stops.

Lifecycle hooks are particularly useful because they can have special privileges 
(like running as the root user), even when the container itself runs with lower privileges 
for security. This means that certain tasks requiring higher permissions can be done without 
compromising the overall security of the container.

### Post-start hooks

Post-start hooks are commands that run after the container has started, but there's no 
set time for when exactly they will execute. The hook runs in parallel with the 
container's `entrypoint`, so avoid using it for tasks the application depends on 
during startup.

Use a post-start hook for side tasks that happen once the container is up, such as 
registering the container with an external system or emitting an audit event. The 
hook can run with higher privileges than the application, which lets it reach 
resources the application user does not have access to.

In the following example, the hook registers the container with an internal service 
discovery endpoint. The application runs as user `1001`, and the hook runs as `root`:

```yaml
services:
  app:
    image: backend
    user: 1001
    post_start:
      - command: wget -q -O - "http://registry.internal/register?name=backend"
        user: root
```

### Pre-stop hooks

Pre-stop hooks are commands that run before the container is stopped by a specific 
command (like `docker compose down` or stopping it manually with `Ctrl+C`). 
These hooks won't run if the container stops by itself or gets killed suddenly.

Use a pre-stop hook for shutdown steps the application itself cannot perform, such as 
notifying an external system that the container is draining. Avoid using it to flush 
application state: a well-written application already handles that on a graceful 
stop, and pre-stop hooks don't run on abrupt termination.

In the following example, the hook de-registers the service from an external 
discovery endpoint before the container stops:

```yaml
services:
  app:
    image: backend
    pre_stop:
      - command: wget -q -O - "http://registry.internal/deregister?name=backend"
```

## Reference information

- [`post_start`](/reference/compose-file/services.md#post_start)
- [`pre_stop`](/reference/compose-file/services.md#pre_stop)
