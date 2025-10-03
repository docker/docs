# docker scout attestation get

<!---MARKER_GEN_START-->
Get attestation for image

### Aliases

`docker scout attestation get`, `docker scout attest get`

### Options

| Name               | Type     | Default                                                    | Description                                                                                             |
|:-------------------|:---------|:-----------------------------------------------------------|:--------------------------------------------------------------------------------------------------------|
| `--key`            | `string` | `https://registry.scout.docker.com/keyring/dhi/latest.pub` | Signature key to use for verification                                                                   |
| `--org`            | `string` |                                                            | Namespace of the Docker organization                                                                    |
| `-o`, `--output`   | `string` |                                                            | Write the report to a file                                                                              |
| `--platform`       | `string` |                                                            | Platform of image to analyze                                                                            |
| `--predicate`      |          |                                                            | Get in-toto predicate only dropping the subject                                                         |
| `--predicate-type` | `string` |                                                            | Predicate-type for attestation                                                                          |
| `--ref`            | `string` |                                                            | Reference to use if the provided tarball contains multiple references.<br>Can only be used with archive |
| `--skip-tlog`      |          |                                                            | Skip signature verification against public transaction log                                              |
| `--verify`         |          |                                                            | Verify the signature on the attestation                                                                 |


<!---MARKER_GEN_END-->

