# docker scout compare

<!---MARKER_GEN_START-->
Compare two images and display differences (experimental)

### Aliases

`docker scout compare`, `docker scout diff`

### Options

| Name                  | Type          | Default | Description                                                                                                                                                                    |
|:----------------------|:--------------|:--------|:-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `-e`, `--exit-code`   |               |         | Return exit code '2' if vulnerability changes are detected                                                                                                                     |
| `--format`            | `string`      | `text`  | Output format of the generated vulnerability report:<br>- text: default output, plain text with or without colors depending on the terminal<br>- markdown: Markdown output<br> |
| `--ignore-base`       |               |         | Filter out CVEs introduced from base image                                                                                                                                     |
| `--ignore-unchanged`  |               |         | Filter out unchanged packages                                                                                                                                                  |
| `--multi-stage`       |               |         | Show packages from multi-stage Docker builds                                                                                                                                   |
| `--only-fixed`        |               |         | Filter to fixable CVEs                                                                                                                                                         |
| `--only-package-type` | `stringSlice` |         | Comma separated list of package types (like apk, deb, rpm, npm, pypi, golang, etc)                                                                                             |
| `--only-severity`     | `stringSlice` |         | Comma separated list of severities (critical, high, medium, low, unspecified) to filter CVEs by                                                                                |
| `--only-stage`        | `stringSlice` |         | Comma separated list of multi-stage Docker build stage names                                                                                                                   |
| `--only-unfixed`      |               |         | Filter to unfixed CVEs                                                                                                                                                         |
| `--org`               | `string`      |         | Namespace of the Docker organization                                                                                                                                           |
| `-o`, `--output`      | `string`      |         | Write the report to a file.                                                                                                                                                    |
| `--platform`          | `string`      |         | Platform of image to analyze                                                                                                                                                   |
| `--ref`               | `string`      |         | Reference to use if the provided tarball contains multiple references.<br>Can only be used with --type archive.                                                                |
| `--to`                | `string`      |         | Image, directory, or archive to compare to                                                                                                                                     |
| `--to-env`            | `string`      |         | Name of environment to compare to                                                                                                                                              |
| `--to-latest`         |               |         | Latest image processed to compare to                                                                                                                                           |
| `--to-ref`            | `string`      |         | Reference to use if the provided tarball contains multiple references.<br>Can only be used with --type archive.                                                                |
| `--to-type`           | `string`      | `image` | Image type to analyze. Can be one of:<br>- image<br>- oci-dir<br>- archive (docker save tarball)<br>                                                                           |
| `--type`              | `string`      | `image` | Type of the image to analyze. Can be one of:<br>- image<br>- oci-dir<br>- archive (docker save tarball)<br>                                                                    |


<!---MARKER_GEN_END-->

## Description

The `docker scout compare` command analyzes two images and displays a comparison of both.

> This command is **experimental** and its behaviour might change in the future

The main usage is to compare two versions of the same image.
For instance when a new image is built and compared to the version running in production.

If no image is specified, the most recently built image will be used.

The following artifact types are supported:

- Images
- OCI layout directories
- Tarball archives, as created by `docker save`

The tool analyzes the provided software artifact, and generates a vulnerability report.

By default, the tool expects an image reference, such as:

- `redis`
- `curlimages/curl:7.87.0`
- `mcr.microsoft.com/dotnet/runtime:7.0`

If the artifact you want to analyze is an OCI directory or a tarball archive, you must use the `--type` or `--to-type` flag.

## Examples

### Compare the most recently built image to the latest tag

```console
$ docker scout compare --to namespace/repo:latest
```

### Ignore base images

```console
$ docker scout compare --ignore-base --to namespace/repo:latest namespace/repo:v1.2.3-pre
```

### Generate a markdown output

```console
$ docker scout compare --format markdown --to namespace/repo:latest namespace/repo:v1.2.3-pre
```

### Only compare maven packages and only display critical vulnerabilities for maven packages

```console
$ docker scout compare --only-package-type maven --only-severity critical --to namespace/repo:latest namespace/repo:v1.2.3-pre
```
