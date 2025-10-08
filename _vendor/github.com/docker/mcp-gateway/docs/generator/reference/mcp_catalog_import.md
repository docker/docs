# docker mcp catalog import

<!---MARKER_GEN_START-->
Import an MCP server catalog from a URL or local file. The catalog will be downloaded 
and stored locally for use with the MCP gateway.

When --mcp-registry flag is used, the argument must be an existing catalog name, and the
command will import servers from the MCP registry URL into that catalog.

### Options

| Name             | Type     | Default | Description                                               |
|:-----------------|:---------|:--------|:----------------------------------------------------------|
| `--dry-run`      | `bool`   |         | Show Imported Data but do not update the Catalog          |
| `--mcp-registry` | `string` |         | Import server from MCP registry URL into existing catalog |


<!---MARKER_GEN_END-->

