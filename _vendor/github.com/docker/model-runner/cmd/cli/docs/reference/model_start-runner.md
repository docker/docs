# docker model start-runner

<!---MARKER_GEN_START-->
Start Docker Model Runner (Docker Engine only)

### Options

| Name             | Type     | Default     | Description                                                                                            |
|:-----------------|:---------|:------------|:-------------------------------------------------------------------------------------------------------|
| `--backend`      | `string` |             | Specify backend (llama.cpp\|vllm\|diffusers). Default: llama.cpp                                       |
| `--debug`        | `bool`   |             | Enable debug logging                                                                                   |
| `--do-not-track` | `bool`   |             | Do not track models usage in Docker Model Runner                                                       |
| `--gpu`          | `string` | `auto`      | Specify GPU support (none\|auto\|cuda\|rocm\|musa\|cann)                                               |
| `--host`         | `string` | `127.0.0.1` | Host address to bind Docker Model Runner                                                               |
| `--port`         | `uint16` | `0`         | Docker container port for Docker Model Runner (default: 12434 for Docker Engine, 12435 for Cloud mode) |
| `--proxy-cert`   | `string` |             | Path to a CA certificate file for proxy SSL inspection                                                 |
| `--tls`          | `bool`   |             | Enable TLS/HTTPS for Docker Model Runner API                                                           |
| `--tls-cert`     | `string` |             | Path to TLS certificate file (auto-generated if not provided)                                          |
| `--tls-key`      | `string` |             | Path to TLS private key file (auto-generated if not provided)                                          |
| `--tls-port`     | `uint16` | `0`         | TLS port for Docker Model Runner (default: 12444 for Docker Engine, 12445 for Cloud mode)              |


<!---MARKER_GEN_END-->

## Description

This command starts the Docker Model Runner without pulling container images. Use this command to start the runner when you already have the required images locally.

For the first-time setup or to ensure you have the latest images, use `docker model install-runner` instead.
