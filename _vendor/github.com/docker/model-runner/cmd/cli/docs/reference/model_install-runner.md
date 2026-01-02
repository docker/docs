# docker model install-runner

<!---MARKER_GEN_START-->
Install Docker Model Runner (Docker Engine only)

### Options

| Name             | Type     | Default     | Description                                                                                            |
|:-----------------|:---------|:------------|:-------------------------------------------------------------------------------------------------------|
| `--backend`      | `string` |             | Specify backend (llama.cpp\|vllm). Default: llama.cpp                                                  |
| `--debug`        | `bool`   |             | Enable debug logging                                                                                   |
| `--do-not-track` | `bool`   |             | Do not track models usage in Docker Model Runner                                                       |
| `--gpu`          | `string` | `auto`      | Specify GPU support (none\|auto\|cuda\|rocm\|musa\|cann)                                               |
| `--host`         | `string` | `127.0.0.1` | Host address to bind Docker Model Runner                                                               |
| `--port`         | `uint16` | `0`         | Docker container port for Docker Model Runner (default: 12434 for Docker Engine, 12435 for Cloud mode) |


<!---MARKER_GEN_END-->

## Description

 This command runs implicitly when a docker model command is executed. You can run this command explicitly to add a new configuration.
