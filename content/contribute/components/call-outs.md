---
description: components and formatting examples used in Docker's docs
title: Callouts
toc_max: 3
---

We support these broad categories of callouts:

- Version callouts
- Notes (no HTML attribute required)
- Experimental, which use the `{ .experimental }` class
- Restricted, which use the `{ .restricted }` class

## Examples

{{< introduced buildx 0.10.4 >}}

> [!NOTE]
>
> Note the way the `get_hit_count` function is written. This basic retry
> loop lets us attempt our request multiple times if the redis service is
> not available. This is useful at startup while the application comes
> online, but also makes our application more resilient if the Redis
> service needs to be restarted anytime during the app's lifetime. In a
> cluster, this also helps handling momentary connection drops between
> nodes.

> [!TIP]
>
> For a smaller base image, use `alpine`.


> [!IMPORTANT]
>
> Treat access tokens like your password and keep them secret. Store your
> tokens securely (for example, in a credential manager).


> [!WARNING]
>
> Removing Volumes
>
> By default, named volumes in your compose file are NOT removed when running
> `docker compose down`. If you want to remove the volumes, you will need to add
> the `--volumes` flag.
>
> The Docker Dashboard does not remove volumes when you delete the app stack.

For both of the following callouts, consult [the Docker release lifecycle](/release-lifecycle) for more information on when to use them.

> **Beta feature**
>
> The Builds view is currently in Beta. This feature may change or be removed from future releases.
{ .experimental }

> **Restricted**
>
> Docker Scout is an [early access](/release-lifecycle/#early-access-ea)
> product.
{ .restricted}

## Formatting 

```go
{{</* introduced buildx 0.10.4 "../../release-notes.md#0104" */>}}
```

```html
> [!NOTE]
>
> Note the way the `get_hit_count` function is written. This basic retry
> loop lets us attempt our request multiple times if the redis service is
> not available. This is useful at startup while the application comes
> online, but also makes our application more resilient if the Redis
> service needs to be restarted anytime during the app's lifetime. In a
> cluster, this also helps handling momentary connection drops between
> nodes.

> [!TIP]
>
> For a smaller base image, use `alpine`.

> [!IMPORTANT]
>
> Treat access tokens like your password and keep them secret. Store your
> tokens securely (for example, in a credential manager).

> [!WARNING]
>
> Removing Volumes
>
> By default, named volumes in your compose file are NOT removed when running
> `docker compose down`. If you want to remove the volumes, you will need to add
> the `--volumes` flag.
>
> The Docker Dashboard does _not_ remove volumes when you delete the app stack.

> **Beta feature**
>
> The Builds view is currently in Beta. This feature may change or be removed from future releases.
{ .experimental }

> **Restricted**
>
> Docker Scout is an [early access](/release-lifecycle/#early-access-ea)
> product.
{ .restricted }
```
