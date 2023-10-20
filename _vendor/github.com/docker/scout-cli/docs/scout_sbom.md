# docker scout sbom

<!---MARKER_GEN_START-->
Generate or display SBOM of an image

### Options

| Name                  | Type          | Default | Description                                                                                                                                   |
|:----------------------|:--------------|:--------|:----------------------------------------------------------------------------------------------------------------------------------------------|
| `--format`            | `string`      | `json`  | Output format:<br>- list: list of packages of the image<br>- json: json representation of the SBOM<br>- spdx: spdx representation of the SBOM |
| `--only-package-type` | `stringSlice` |         | Comma separated list of package types (like apk, deb, rpm, npm, pypi, golang, etc)<br>Can only be used with --format list                     |
| `-o`, `--output`      | `string`      |         | Write the report to a file.                                                                                                                   |
| `--platform`          | `string`      |         | Platform of image to analyze                                                                                                                  |
| `--ref`               | `string`      |         | Reference to use if the provided tarball contains multiple references.<br>Can only be used with archive.                                      |


<!---MARKER_GEN_END-->

## Description

The `docker scout sbom` command analyzes a software artifact to generate the corresponding Software Bill Of Materials (SBOM).

The SBOM can be used to list all packages, or the ones from a specific type (as dep, maven, etc).

If no image is specified, the most recently built image will be used.

The following artifact types are supported:

- Images
- OCI layout directories
- Tarball archives, as created by `docker save`
- Local directory or file

The tool analyzes the provided software artifact, and generates a vulnerability report.

By default, the tool expects an image reference, such as:

- `redis`
- `curlimages/curl:7.87.0`
- `mcr.microsoft.com/dotnet/runtime:7.0`

If the artifact you want to analyze is an OCI directory, a tarball archive, a local file or directory,
or if you want to control from where the image will be resolved, you must prefix the reference with one of the following:

- `image://` (default) use a local image, or fall back to a registry lookup
- `local://` use an image from the local image store (don't do a registry lookup)
- `registry://` use an image from a registry (don't use a local image)
- `oci-dir://` use an OCI layout directory
- `archive://` use a tarball archive, as created by docker save
- `fs://` use a local directory or file

## Examples

### Display the list of packages

```console
$ docker scout sbom --format list alpine
```

### Only display packages of a specific type

```console
 $ docker scout sbom --format list --only-package-type apk alpine
```

### Display the full SBOM as json

```console
$ docker scout sbom alpine
```

### Display the full SBOM of the most recently buitl image

```console
$ docker scout sbom
```

### Write SBOM to a file

```console
$ docker scout sbom --output alpine.sbom alpine
```
