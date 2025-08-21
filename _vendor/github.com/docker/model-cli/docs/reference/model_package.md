# docker model package

<!---MARKER_GEN_START-->
Package a GGUF file into a Docker model OCI artifact, with optional licenses. The package is sent to the model-runner, unless --push is specified

### Options

| Name              | Type          | Default | Description                                                                            |
|:------------------|:--------------|:--------|:---------------------------------------------------------------------------------------|
| `--context-size`  | `uint64`      | `0`     | context size in tokens                                                                 |
| `--gguf`          | `string`      |         | absolute path to gguf file (required)                                                  |
| `-l`, `--license` | `stringArray` |         | absolute path to a license file                                                        |
| `--push`          | `bool`        |         | push to registry (if not set, the model is loaded into the Model Runner content store. |


<!---MARKER_GEN_END-->

