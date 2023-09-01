# docker scout recommendations

<!---MARKER_GEN_START-->
Display available base image updates and remediation recommendations

### Options

| Name             | Type     | Default | Description                                                                                                     |
|:-----------------|:---------|:--------|:----------------------------------------------------------------------------------------------------------------|
| `--only-refresh` |          |         | Only display base image refresh recommendations                                                                 |
| `--only-update`  |          |         | Only display base image update recommendations                                                                  |
| `--org`          | `string` |         | Namespace of the Docker organization                                                                            |
| `-o`, `--output` | `string` |         | Write the report to a file.                                                                                     |
| `--platform`     | `string` |         | Platform of image to analyze                                                                                    |
| `--ref`          | `string` |         | Reference to use if the provided tarball contains multiple references.<br>Can only be used with --type archive. |
| `--tag`          | `string` |         | Specify tag                                                                                                     |
| `--type`         | `string` | `image` | Type of the image to analyze. Can be one of:<br>- image<br>- oci-dir<br>- archive (docker save tarball)<br>     |


<!---MARKER_GEN_END-->

## Description

The `docker scout recommendations` command display recommendations for base images updates.
It analyzes the image and display recommendations to refresh or update the base image.
For each recommendation it shows a list of benefits like less vulnerabilities, smaller image, etc.

If no image is specified, the most recently built image will be used.

The following artifact types are supported:

- Images
- OCI layout directories
- Tarball archives, as created by `docker save`

The tool analyzes the provided software artifact, and generates base image updates and remediation recommendations.

By default, the tool expects an image reference, such as:

- `redis`
- `curlimages/curl:7.87.0`
- `mcr.microsoft.com/dotnet/runtime:7.0`

If the artifact you want to analyze is an OCI directory or a tarball archive, you must use the `--type` flag.

## Examples

### Display base image update recommendations

```console
$ docker scout recommendations golang:1.19.4
```

### Display base image refresh only recommendations

```console
$ docker scout recommendations --only-refresh golang:1.19.4
```

### Display base image update only recommendations

```console
$ docker scout recommendations --only-update golang:1.19.4
```
