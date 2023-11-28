# docker scout stream

<!---MARKER_GEN_START-->
Manage streams (experimental)

### Options

| Name             | Type     | Default | Description                          |
|:-----------------|:---------|:--------|:-------------------------------------|
| `--org`          | `string` |         | Namespace of the Docker organization |
| `-o`, `--output` | `string` |         | Write the report to a file.          |
| `--platform`     | `string` |         | Platform of image to record          |


<!---MARKER_GEN_END-->

## Description

The `docker scout stream` command lists the deployment streams and records an image to it.

Once recorded, streams can be referred to by their name, eg. in the `docker scout compare` command using `--to-stream`.

## Examples

### List existing streams

```console
$ %[1]s %[2]s
prod-cluster-123
stage-cluster-234
```

### List images of a stream

```console
$ %[1]s %[2]s prod-cluster-123
namespace/repo:tag@sha256:9a4df4fadc9bbd44c345e473e0688c2066a6583d4741679494ba9228cfd93e1b
namespace/other-repo:tag@sha256:0001d6ce124855b0a158569c584162097fe0ca8d72519067c2c8e3ce407c580f
```

### Record an image to a stream, for a specific platform

```console
$ %[1]s %[2]s stage-cluster-234 namespace/repo:stage-latest --platform linux/amd64
✓ Pulled
✓ Successfully recorded namespace/repo:stage-latest in stream stage-cluster-234
```
