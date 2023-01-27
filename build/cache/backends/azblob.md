---
title: "Azure Blob Storage cache"
keywords: build, buildx, cache, backend, azblob, azure
redirect_from:
  - /build/building/cache/backends/azblob/
---

> **Warning**
>
> This cache backend is unreleased. You can use it today, by using the
> `moby/buildkit:master` image in your Buildx driver.

The `azblob` cache store uploads your resulting build cache to
[Azure's blob storage service](https://azure.microsoft.com/en-us/services/storage/blobs/){:target="blank" rel="noopener" class=""}.

> **Note**
>
> This cache storage backend requires using a different driver than the default
> `docker` driver - see more information on selecting a driver
> [here](../../drivers/index.md). To create a new driver (which can act as a
> simple drop-in replacement):
>
> ```console
> $ docker buildx create --use --driver=docker-container
> ```

## Synopsis

```console
$ docker buildx build --push -t <registry>/<image> \
  --cache-to type=azblob,name=<cache-image>[,parameters...] \
  --cache-from type=azblob,name=<cache-image>[,parameters...] .
```

The following table describes the available CSV parameters that you can pass to
`--cache-to` and `--cache-from`.

| Name                | Option                  | Type        | Default | Description                                        |
|---------------------|-------------------------|-------------|---------|----------------------------------------------------|
| `name`              | `cache-to`,`cache-from` | String      |         | Required. The name of the cache image.             |
| `account_url`       | `cache-to`,`cache-from` | String      |         | Base URL of the storage account.                   |
| `secret_access_key` | `cache-to`,`cache-from` | String      |         | Blob storage account key, see [authentication][1]. |
| `mode`              | `cache-to`              | `min`,`max` | `min`   | Cache layers to export, see [cache mode][2].       |

[1]: #authentication
[2]: index.md#cache-mode

## Authentication

The `secret_access_key`, if left unspecified, is read from environment variables
on the BuildKit server following the scheme for the
[Azure Go SDK](https://docs.microsoft.com/en-us/azure/developer/go/azure-sdk-authentication){:target="blank" rel="noopener" class=""}.
The environment variables are read from the server, not the Buildx client.

## Further reading

For an introduction to caching see [Optimizing builds with cache](../index.md).

For more information on the `azblob` cache backend, see the
[BuildKit README](https://github.com/moby/buildkit#azure-blob-storage-cache-experimental){:target="blank" rel="noopener" class=""}.
