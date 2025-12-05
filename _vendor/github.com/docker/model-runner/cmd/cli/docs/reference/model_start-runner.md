# docker model start-runner

<!---MARKER_GEN_START-->
Start Docker Model Runner (Docker Engine only)

### Options

| Name             | Type     | Default | Description                                                                                            |
|:-----------------|:---------|:--------|:-------------------------------------------------------------------------------------------------------|
| `--backend`      | `string` |         | Specify backend (llama.cpp\|vllm). Default: llama.cpp                                                  |
| `--debug`        | `bool`   |         | Enable debug logging                                                                                   |
| `--do-not-track` | `bool`   |         | Do not track models usage in Docker Model Runner                                                       |
| `--gpu`          | `string` | `auto`  | Specify GPU support (none\|auto\|cuda\|rocm\|musa\|cann)                                               |
| `--port`         | `uint16` | `0`     | Docker container port for Docker Model Runner (default: 12434 for Docker Engine, 12435 for Cloud mode) |


<!---MARKER_GEN_END-->

## Description

This command starts the Docker Model Runner without pulling container images. Use this command to start the runner when you already have the required images locally.

For the first-time setup or to ensure you have the latest images, use `docker model install-runner` instead.
