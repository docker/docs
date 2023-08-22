# docker scout cves

```
docker scout cves [OPTIONS] [IMAGE|DIRECTORY|ARCHIVE]
```

<!---MARKER_GEN_START-->
Display CVEs identified in a software artifact

### Options

| Name                  | Type          | Default    | Description                                                                                                                                                                                              |
|:----------------------|:--------------|:-----------|:---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `--details`           |               |            | Print details on default text output                                                                                                                                                                     |
| `-e`, `--exit-code`   |               |            | Return exit code '2' if vulnerabilities are detected                                                                                                                                                     |
| `--format`            | `string`      | `packages` | Output format of the generated vulnerability report:<br>- packages: default output, plain text with vulnerabilities grouped by packages<br>- sarif: json Sarif output<br>- markdown: markdown output<br> |
| `--ignore-base`       |               |            | Filter out CVEs introduced from base image                                                                                                                                                               |
| `--locations`         |               |            | Print package locations including file paths and layer diff_id                                                                                                                                           |
| `--multi-stage`       |               |            | Show packages from multi-stage Docker builds                                                                                                                                                             |
| `--only-cve-id`       | `stringSlice` |            | Comma separated list of CVE ids (like CVE-2021-45105) to search for                                                                                                                                      |
| `--only-fixed`        |               |            | Filter to fixable CVEs                                                                                                                                                                                   |
| `--only-package-type` | `stringSlice` |            | Comma separated list of package types (like apk, deb, rpm, npm, pypi, golang, etc)                                                                                                                       |
| `--only-severity`     | `stringSlice` |            | Comma separated list of severities (critical, high, medium, low, unspecified) to filter CVEs by                                                                                                          |
| `--only-stage`        | `stringSlice` |            | Comma separated list of multi-stage Docker build stage names                                                                                                                                             |
| `--only-unfixed`      |               |            | Filter to unfixed CVEs                                                                                                                                                                                   |
| `-o`, `--output`      | `string`      |            | Write the report to a file.                                                                                                                                                                              |
| `--platform`          | `string`      |            | Platform of image to analyze                                                                                                                                                                             |
| `--ref`               | `string`      |            | Reference to use if the provided tarball contains multiple references.<br>Can only be used with --type archive.                                                                                          |
| `--type`              | `string`      | `image`    | Type of the image to analyze. Can be one of:<br>- image<br>- oci-dir<br>- archive (docker save tarball)<br>                                                                                              |


<!---MARKER_GEN_END-->

## Description

The `docker scout cves` command analyzes a software artifact for vulnerabilities.

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

If the artifact you want to analyze is an OCI directory or a tarball archive, you must use the `--type` flag.

## Examples

### Display vulnerabilities grouped by package

```console
$ docker scout cves alpine
Analyzing image alpine
    ✓ Image stored for indexing
    ✓ Indexed 18 packages
    ✓ No vulnerable package detected
```

### Display vulnerabilities from a `docker save` tarball

```console
$ docker save alpine > alpine.tar

$ docker scout cves --type archive alpine.tar
Analyzing archive alpine.tar
    ✓ Archive read
    ✓ SBOM of image already cached, 18 packages indexed
    ✓ No vulnerable package detected
```

### Display vulnerabilities from an OCI directory

```console
$ skopeo copy --override-os linux docker://alpine oci:alpine

$ docker scout cves --type oci-dir alpine
Analyzing OCI directory alpine
    ✓ OCI directory read
    ✓ Image stored for indexing
    ✓ Indexed 19 packages
    ✓ No vulnerable package detected
```

### Export vulnerabilities to a SARIF JSON file

```console
$ docker scout cves --format sarif --output alpine.sarif.json alpine
Analyzing image alpine
    ✓ SBOM of image already cached, 18 packages indexed
    ✓ No vulnerable package detected
    ✓ Report written to alpine.sarif.json
```
