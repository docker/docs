---
description: Docker Hub API v1 (deprecated)
keywords: kitematic, deprecated
title: Docker Hub API v1 (deprecated)
---

> **Deprecated**
>
> Docker Hub API v1 has been deprecated. Please use Docker Hub API v1 instead.
{: .warning }

The following API routes within the v1 path will no longer work and will return a 410 status code:
* `/v1/repositories/{name}/images`
* `/v1/repositories/{name}/tags`
* `/v1/repositories/{name}/tags/{tag_name}`
* `/v1/repositories/{namespace}/{name}/images`
* `/v1/repositories/{namespace}/{name}/tags`
* `/v1/repositories/{namespace}/{name}/tags/{tag_name}`

If you want to continue using the Docker Hub API in your current applications, update your clients to use v2 endpoints.

| **OLD** | **NEW** |
| -------------- | ------------ |
| [/v1/repositories/{name}/images](https://github.com/moby/moby/blob/v1.3.0/docs/sources/reference/api/docker-io_api.md#list-user-repository-images) *(Docker Hub API V1)*| [/v2/namespaces/{namespace}/repositories/{repository}/images](https://docs.docker.com/docker-hub/api/latest/#tag/images/operation/GetNamespacesRepositoriesImages) *(Docker Hub API V2)* | 
|  [/v1/repositories/{namespace}/{name}/images](https://github.com/moby/moby/blob/v1.3.0/docs/sources/reference/api/docker-io_api.md#list-user-repository-images) *(Docker Hub API V1)*| [/v2/namespaces/{namespace}/repositories/{repository}/images](https://docs.docker.com/docker-hub/api/latest/#tag/images/operation/GetNamespacesRepositoriesImages) *(Docker Hub API V2)* |
| [/v1/repositories/{name}/tags](https://github.com/moby/moby/blob/v1.8.3/docs/reference/api/registry_api.md#list-repository-tags) *(Registry API V1)*| [/v2/{name}/tags/list (where {name} is {namespace}/{name})](https://github.com/opencontainers/distribution-spec/blob/v1.0.1/spec.md#content-discovery) *(OCI Distribution Spec)* |
| [/v1/repositories/{namespace}/{name}/tags](https://github.com/moby/moby/blob/v1.8.3/docs/reference/api/registry_api.md#list-repository-tags) *(Registry API V1)*| [/v2/{name}/tags/list (where {name} is {namespace}/{name})](https://github.com/opencontainers/distribution-spec/blob/v1.0.1/spec.md#content-discovery) *(OCI Distribution Spec)* |
| **GET** [/v1/repositories/{namespace}/{name}/tags](https://github.com/moby/moby/blob/v1.8.3/docs/reference/api/registry_api.md#get-image-id-for-a-particular-tag) *(Registry API V1)* | **HEAD** [/v2/{name}/manifests/{reference}/](https://github.com/opencontainers/distribution-spec/blob/v1.0.1/spec.md#pulling-manifests) *(OCI Distribution Spec)* |
| **GET** [/v1/repositories/{namespace}/{name}/tags/{tag_name}](https://github.com/moby/moby/blob/v1.8.3/docs/reference/api/registry_api.md#get-image-id-for-a-particular-tag) *(Registry API V1)* | **HEAD** [/v2/{name}/manifests/{reference}/](https://github.com/opencontainers/distribution-spec/blob/v1.0.1/spec.md#pulling-manifests) *(OCI Distribution Spec)* |
