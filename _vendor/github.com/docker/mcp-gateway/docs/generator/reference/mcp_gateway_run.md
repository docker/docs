# docker mcp gateway run

<!---MARKER_GEN_START-->
Run the gateway

### Options

| Name                        | Type          | Default             | Description                                                                                                                                   |
|:----------------------------|:--------------|:--------------------|:----------------------------------------------------------------------------------------------------------------------------------------------|
| `--additional-catalog`      | `stringSlice` |                     | Additional catalog paths to append to the default catalogs                                                                                    |
| `--additional-config`       | `stringSlice` |                     | Additional config paths to merge with the default config.yaml                                                                                 |
| `--additional-registry`     | `stringSlice` |                     | Additional registry paths to merge with the default registry.yaml                                                                             |
| `--additional-tools-config` | `stringSlice` |                     | Additional tools paths to merge with the default tools.yaml                                                                                   |
| `--block-network`           | `bool`        |                     | Block tools from accessing forbidden network resources                                                                                        |
| `--block-secrets`           | `bool`        | `true`              | Block secrets from being/received sent to/from tools                                                                                          |
| `--catalog`                 | `stringSlice` | `[docker-mcp.yaml]` | Paths to docker catalogs (absolute or relative to ~/.docker/mcp/catalogs/)                                                                    |
| `--config`                  | `stringSlice` | `[config.yaml]`     | Paths to the config files (absolute or relative to ~/.docker/mcp/)                                                                            |
| `--cpus`                    | `int`         | `1`                 | CPUs allocated to each MCP Server (default is 1)                                                                                              |
| `--debug-dns`               | `bool`        |                     | Debug DNS resolution                                                                                                                          |
| `--dry-run`                 | `bool`        |                     | Start the gateway but do not listen for connections (useful for testing the configuration)                                                    |
| `--enable-all-servers`      | `bool`        |                     | Enable all servers in the catalog (instead of using individual --servers options)                                                             |
| `--interceptor`             | `stringArray` |                     | List of interceptors to use (format: when:type:path, e.g. 'before:exec:/bin/path')                                                            |
| `--log-calls`               | `bool`        | `true`              | Log calls to the tools                                                                                                                        |
| `--long-lived`              | `bool`        |                     | Containers are long-lived and will not be removed until the gateway is stopped, useful for stateful servers                                   |
| `--mcp-registry`            | `stringSlice` |                     | MCP registry URLs to fetch servers from (can be repeated)                                                                                     |
| `--memory`                  | `string`      | `2Gb`               | Memory allocated to each MCP Server (default is 2Gb)                                                                                          |
| `--oci-ref`                 | `stringArray` |                     | OCI image references to use                                                                                                                   |
| `--port`                    | `int`         | `0`                 | TCP port to listen on (default is to listen on stdio)                                                                                         |
| `--registry`                | `stringSlice` | `[registry.yaml]`   | Paths to the registry files (absolute or relative to ~/.docker/mcp/)                                                                          |
| `--secrets`                 | `string`      | `docker-desktop`    | Colon separated paths to search for secrets. Can be `docker-desktop` or a path to a .env file (default to using Docker Desktop's secrets API) |
| `--servers`                 | `stringSlice` |                     | Names of the servers to enable (if non empty, ignore --registry flag)                                                                         |
| `--static`                  | `bool`        |                     | Enable static mode (aka pre-started servers)                                                                                                  |
| `--tools`                   | `stringSlice` |                     | List of tools to enable                                                                                                                       |
| `--tools-config`            | `stringSlice` | `[tools.yaml]`      | Paths to the tools files (absolute or relative to ~/.docker/mcp/)                                                                             |
| `--transport`               | `string`      | `stdio`             | stdio, sse or streaming (default is stdio)                                                                                                    |
| `--verbose`                 | `bool`        |                     | Verbose output                                                                                                                                |
| `--verify-signatures`       | `bool`        |                     | Verify signatures of the server images                                                                                                        |
| `--watch`                   | `bool`        | `true`              | Watch for changes and reconfigure the gateway                                                                                                 |


<!---MARKER_GEN_END-->

