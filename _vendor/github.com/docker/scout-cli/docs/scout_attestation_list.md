# docker scout attestation list

<!---MARKER_GEN_START-->
List attestations for image

### Aliases

`docker scout attestation list`, `docker scout attest list`

### Options

| Name               | Type     | Default | Description                                                                                                                         |
|:-------------------|:---------|:--------|:------------------------------------------------------------------------------------------------------------------------------------|
| `--format`         | `string` | `list`  | Output format:<br>- list: list of attestations of the image<br>- json: json representation of the attestation list (default "json") |
| `--org`            | `string` |         | Namespace of the Docker organization                                                                                                |
| `-o`, `--output`   | `string` |         | Write the report to a file                                                                                                          |
| `--platform`       | `string` |         | Platform of image to analyze                                                                                                        |
| `--predicate-type` | `string` |         | Predicate-type for attestations                                                                                                     |
| `--ref`            | `string` |         | Reference to use if the provided tarball contains multiple references.<br>Can only be used with archive                             |


<!---MARKER_GEN_END-->

