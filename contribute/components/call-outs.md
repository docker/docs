---
description: components and formatting examples used in Docker's docs
title: Callouts
toc_max: 3
---

We support these broad categories of callouts:

- Notes (no Liquid tag required)
- Important, which use the `{: .important}` tag
- Warning , which use the `{: .warning}` tag

## Examples

  > **Note**
  >
  > Note the way the `get_hit_count` function is written. This basic retry
  > loop lets us attempt our request multiple times if the redis service is
  > not available. This is useful at startup while the application comes
  > online, but also makes our application more resilient if the Redis
  > service needs to be restarted anytime during the app's lifetime. In a
  > cluster, this also helps handling momentary connection drops between
  > nodes.

> **Important**
>
> Treat access tokens like your password and keep them secret. Store your
> tokens securely (for example, in a credential manager).
{: .important}


> **Warning**
>
> Removing Volumes
>
> By default, named volumes in your compose file are NOT removed when running
> `docker-compose down`. If you want to remove the volumes, you will need to add
> the `--volumes` flag.
>
> The Docker Dashboard does not remove volumes when you delete the app stack.
{: .warning}

## HTML

```html
> **Note**
>
> Note the way the `get_hit_count` function is written. This basic retry
> loop lets us attempt our request multiple times if the redis service is
> not available. This is useful at startup while the application comes
> online, but also makes our application more resilient if the Redis
> service needs to be restarted anytime during the app's lifetime. In a
> cluster, this also helps handling momentary connection drops between
> nodes.

> **Important**
>
> Treat access tokens like your password and keep them secret. Store your
> tokens securely (for example, in a credential manager).
{: .important} 

> **Warning**
>
> Removing Volumes
>
> By default, named volumes in your compose file are NOT removed when running
> `docker-compose down`. If you want to remove the volumes, you will need to add
> the `--volumes` flag.
>
> The Docker Dashboard does _not_ remove volumes when you delete the app stack.
{: .warning}
```