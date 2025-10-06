# docker model package

<!---MARKER_GEN_START-->
Package a GGUF file into a Docker model OCI artifact, with optional licenses. The package is sent to the model-runner, unless --push is specified.
When packaging a sharded model --gguf should point to the first shard. All shard files should be siblings and should include the index in the file name (e.g. model-00001-of-00015.gguf).

### Options

| Name              | Type          | Default | Description                                                                            |
|:------------------|:--------------|:--------|:---------------------------------------------------------------------------------------|
| `--chat-template` | `string`      |         | absolute path to chat template file (must be Jinja format)                             |
| `--context-size`  | `uint64`      | `0`     | context size in tokens                                                                 |
| `--gguf`          | `string`      |         | absolute path to gguf file (required)                                                  |
| `-l`, `--license` | `stringArray` |         | absolute path to a license file                                                        |
| `--push`          | `bool`        |         | push to registry (if not set, the model is loaded into the Model Runner content store) |


<!---MARKER_GEN_END-->

