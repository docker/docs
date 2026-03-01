# docker model restart-runner

<!---MARKER_GEN_START-->
Restart Docker Model Runner (Docker Engine only)

### Options

| Name             | Type     | Default     | Description                                                                                            |
|:-----------------|:---------|:------------|:-------------------------------------------------------------------------------------------------------|
| `--debug`        | `bool`   |             | Enable debug logging                                                                                   |
| `--do-not-track` | `bool`   |             | Do not track models usage in Docker Model Runner                                                       |
| `--gpu`          | `string` | `auto`      | Specify GPU support (none\|auto\|cuda\|rocm\|musa\|cann)                                               |
| `--host`         | `string` | `127.0.0.1` | Host address to bind Docker Model Runner                                                               |
| `--port`         | `uint16` | `0`         | Docker container port for Docker Model Runner (default: 12434 for Docker Engine, 12435 for Cloud mode) |
| `--proxy-cert`   | `string` |             | Path to a CA certificate file for proxy SSL inspection                                                 |


<!---MARKER_GEN_END-->

## Description

This command restarts the Docker Model Runner without pulling container images. Use this command to restart the runner when you already have the required images locally.

For the first-time setup or to ensure you have the latest images, use `docker model install-runner` instead.
