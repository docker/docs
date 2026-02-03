# docker model stop-runner

<!---MARKER_GEN_START-->
Stop Docker Model Runner (Docker Engine only)

### Options

| Name       | Type   | Default | Description                 |
|:-----------|:-------|:--------|:----------------------------|
| `--models` | `bool` |         | Remove model storage volume |


<!---MARKER_GEN_END-->

## Description

This command stops the Docker Model Runner by removing the running containers, but preserves the container images on disk. Use this command when you want to temporarily stop the runner but plan to start it again later.

To completely remove the runner including images, use `docker model uninstall-runner --images` instead.
