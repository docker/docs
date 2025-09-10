# docker model install-runner

<!---MARKER_GEN_START-->
Install Docker Model Runner (Docker Engine only)

### Options

| Name             | Type     | Default | Description                                      |
|:-----------------|:---------|:--------|:-------------------------------------------------|
| `--do-not-track` | `bool`   |         | Do not track models usage in Docker Model Runner |
| `--gpu`          | `string` | `auto`  | Specify GPU support (none\|auto\|cuda)           |
| `--port`         | `uint16` | `12434` | Docker container port for Docker Model Runner    |


<!---MARKER_GEN_END-->

## Description

 This command runs implicitly when a docker model command is executed. You can run this command explicitly to add a new configuration.
