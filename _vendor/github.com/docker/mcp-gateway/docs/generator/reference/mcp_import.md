# docker mcp import

<!---MARKER_GEN_START-->
Import and parse a server definition from an official MCP registry URL.

This command fetches the server definition from the provided URL, parses it as a ServerDetail,
converts it to the internal Server format, and displays the results.

Example:
  docker mcp officialregistry import https://registry.example.com/servers/my-server

### Options

| Name             | Type     | Default | Description                     |
|:-----------------|:---------|:--------|:--------------------------------|
| `--catalog`      | `string` |         | import to local catalog         |
| `--mcp-registry` | `string` |         | import from MCP registry format |
| `--push`         | `bool`   |         | push the new server artifact    |


<!---MARKER_GEN_END-->

