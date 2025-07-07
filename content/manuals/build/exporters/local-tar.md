---
title: Local and tar exporters
keywords: build, buildx, buildkit, exporter, local, tar
description: >
  The local and tar exporters save the build result to the local filesystem
aliases:
  - /build/building/exporters/local-tar/
---

The `local` and `tar` exporters output the root filesystem of the build result
into a local directory. They're useful for producing artifacts that aren't
container images.

- `local` exports files and directories.
- `tar` exports the same, but bundles the export into a tarball.

## Synopsis

Build a container image using the `local` exporter:

```console
$ docker buildx build --output type=local[,parameters] .
$ docker buildx build --output type=tar[,parameters] .
```

The following table describes the available parameters:

| Parameter        | Type    | Default | Description                                                                                                                                                                                                                            |
|------------------|---------|---------|----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `dest`           | String  |         | Path to copy files to                                                                                                                                                                                                                  |
| `platform-split` | Boolean | `true`  | When using the local exporter with a multi-platform build, by default, a subfolder matching each target platform is created in the destination directory. Set it to `false` to merge files from all platforms into the same directory. |

## Further reading

For more information on the `local` or `tar` exporters, see the
[BuildKit README](https://github.com/moby/buildkit/blob/master/README.md#local-directory).
