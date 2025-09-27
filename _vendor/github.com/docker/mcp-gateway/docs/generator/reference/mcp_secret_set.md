# docker mcp secret set

<!---MARKER_GEN_START-->
Set a secret in Docker Desktop's secret store

### Options

| Name         | Type     | Default | Description                            |
|:-------------|:---------|:--------|:---------------------------------------|
| `--provider` | `string` |         | Supported: credstore, oauth/<provider> |


<!---MARKER_GEN_END-->

## Examples

### Use secrets for postgres password with default policy

```console
docker mcp secret set POSTGRES_PASSWORD=my-secret-password
docker run -d -l x-secret:POSTGRES_PASSWORD=/pwd.txt -e POSTGRES_PASSWORD_FILE=/pwd.txt -p 5432 postgres
```

### Pass the secret via STDIN

```console
echo my-secret-password > pwd.txt
cat pwd.txt | docker mcp secret set POSTGRES_PASSWORD
```