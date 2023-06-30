---
description: components and formatting examples used in Docker's docs
title: Callouts
toc_max: 3
---

We support these broad categories of callouts:

- Notes (no Liquid tag required)
- Tips, which use the `{: .tip}` tag
- Important, which use the `{: .important}` tag
- Warning , which use the `{: .warning}` tag
- Experimental, which use the `{: .experimental}` tag
- Restricted, which use the `{: .restricted}` tag

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

> **Tip**
>
> For a smaller base image, use `alpine`.
{: .tip }


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

> **Beta feature**
>
> The Builds view is currently in Beta. This feature may change or be removed from future releases.
{: .experimental}

> **Restricted**
>
> Docker Scout is an [early access](/release-lifecycle/#early-access-ea)
> product.
{: .restricted}

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

> **Tip**
>
> For a smaller base image, use `alpine`.
{: .tip }

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

> **Beta feature**
>
> The Builds view is currently in Beta. This feature may change or be removed from future releases.
{: .experimental}

> **Restricted**
>
> Docker Scout is an [early access](/release-lifecycle/#early-access-ea)
> product.
{: .restricted}
```
