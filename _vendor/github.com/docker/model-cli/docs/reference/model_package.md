# docker model package

<!---MARKER_GEN_START-->
Package a GGUF file into a Docker model OCI artifact, with optional licenses, and pushes it to the specified registry

### Options

| Name              | Type          | Default | Description                           |
|:------------------|:--------------|:--------|:--------------------------------------|
| `--gguf`          | `string`      |         | absolute path to gguf file (required) |
| `-l`, `--license` | `stringArray` |         | absolute path to a license file       |
| `--push`          | `bool`        |         | push to registry (required)           |


<!---MARKER_GEN_END-->

