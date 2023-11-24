# docker scout watch

<!---MARKER_GEN_START-->
Watch repositories in a registry and push images and indexes to Docker Scout (experimental)

### Options

| Name                 | Type          | Default | Description                                                                         |
|:---------------------|:--------------|:--------|:------------------------------------------------------------------------------------|
| `--all-images`       |               |         | Push all images instead of only the ones pushed during the watch command is running |
| `--dry-run`          |               |         | Watch images and prepare them, but do not push them                                 |
| `--interval`         | `int64`       | `60`    | Interval in seconds between checks                                                  |
| `--org`              | `string`      |         | Namespace of the Docker organization to which image will be pushed                  |
| `--refresh-registry` |               |         | Refresh the list of repositories of a registry at every run. Only with --registry.  |
| `--registry`         | `string`      |         | Registry to watch                                                                   |
| `--repository`       | `stringSlice` |         | Repository to watch                                                                 |
| `--sbom`             |               |         | Create and upload SBOMs                                                             |
| `--tag`              | `stringSlice` |         | Regular expression to match tags to watch                                           |
| `--workers`          | `int`         | `3`     | Number of concurrent workers                                                        |


<!---MARKER_GEN_END-->

## Description

The `docker scout watch` command watches repositories in a registry
and pushes images or analysis results to Docker Scout.

## Examples

### Watch for new images from two repositories and push them

```console
$ docker scout watch --org my-org --repository registry-1.example.com/repo-1 --repository registry-2.example.com/repo-2
```

### Only push images with a specific tag

```console
$ docker scout watch --org my-org --repository registry.example.com/my-service --tag latest
```

### Watch all repositories of a registry

```console
$ docker scout watch --org my-org --registry registry.example.com
```

### Push all images and not just the new ones

```console
$ docker scout watch--org my-org --repository registry.example.com/my-service --all-images
```
