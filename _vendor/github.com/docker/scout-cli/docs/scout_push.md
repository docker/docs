# docker scout push

<!---MARKER_GEN_START-->
Push an image or image index to Docker Scout (experimental)

### Options

| Name             | Type     | Default | Description                                                        |
|:-----------------|:---------|:--------|:-------------------------------------------------------------------|
| `--author`       | `string` |         | Name of the author of the image                                    |
| `--org`          | `string` |         | Namespace of the Docker organization to which image will be pushed |
| `-o`, `--output` | `string` |         | Write the report to a file.                                        |
| `--sbom`         |          |         | Create and upload SBOMs                                            |
| `--timestamp`    | `string` |         | Timestamp of image or tag creation                                 |


<!---MARKER_GEN_END-->

## Description

The `docker scout push` command lets you push an image or analysis result to Docker Scout.

## Examples

### Push an image to Docker Scout

```console
$ docker scout push --org my-org registry.example.com/repo:tag
```
