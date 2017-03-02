---
title: Use a cache
description: Learn how to configure your Docker Trusted Registry account to pull images from a cache for faster download times.
keywords: docker, registry, cache
---

Docker Trusted Registry can be configured to have one or more caches. This
allows you to choose from which cache to pull images from for faster
download times.

If an administrator has set up caches, you can choose which cache to use when
pulling images. In the **DTR web UI**, navigate to your **user profile**,
and check the **Content Cache** options.

![](../../images/use-a-cache-1.png){: .with-border}

Once you save, your images are pulled from the cache instead of the central DTR.
