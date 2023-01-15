---
title: High-level builds with Bake
keywords: build, buildx, bake, buildkit, hcl, json, compose
aliases:
- /build/customize/bake/
---

> This command is experimental.
>
> The design of bake is in early stages, and we are looking for [feedback from users](https://github.com/docker/buildx/issues){:target="blank" rel="noopener" class=""}.
{: .experimental }

Buildx also aims to provide support for high-level build concepts that go beyond
invoking a single build command. We want to support building all the images in
your application together and let the users define project specific reusable
build flows that can then be easily invoked by anyone.

[BuildKit](https://github.com/moby/buildkit){:target="blank" rel="noopener" class=""}
efficiently handles multiple concurrent build requests and de-duplicating work.
The build commands can be combined with general-purpose command runners
(for example, `make`). However, these tools generally invoke builds in sequence
and therefore cannot leverage the full potential of BuildKit parallelization,
or combine BuildKit's output for the user. For this use case, we have added a
command called [`docker buildx bake`](../../engine/reference/commandline/buildx_bake.md).

The `bake` command supports building images from HCL, JSON and Compose files.
This is similar to [`docker compose build`](../../compose/compose-file/build.md),
but allowing all the services to be built concurrently as part of a single
request. If multiple files are specified they are all read and configurations are
combined.

We recommend using HCL files as its experience is more aligned with buildx UX
and also allows better code reuse, different target groups and extended features.

## Next steps

* [File definition](file-definition.md)
* [Configuring builds](configuring-build.md)
* [User defined HCL functions](hcl-funcs.md)
* [Defining additional build contexts and linking targets](build-contexts.md)
* [Building from Compose file](compose-file.md)
