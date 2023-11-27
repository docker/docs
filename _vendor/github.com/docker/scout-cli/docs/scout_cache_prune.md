# docker scout cache prune

<!---MARKER_GEN_START-->
Remove temporary or cached data

### Options

| Name            | Type | Default | Description                    |
|:----------------|:-----|:--------|:-------------------------------|
| `-f`, `--force` |      |         | Do not prompt for confirmation |
| `--sboms`       |      |         | Prune cached SBOMs             |


<!---MARKER_GEN_END-->

## Description

The `docker scout cache prune` command removes temporary data and SBOM cache.

By default, `docker scout cache prune` only deletes temporary data.
To delete temporary data and clear the SBOM cache, use the `--sboms` flag.

## Examples

### Delete temporary data

```console
$ docker scout cache prune
? Are you sure to delete all temporary data? Yes
    ✓ temporary data deleted
```

### Delete temporary _and_ cache data

```console
$ docker scout cache prune --sboms
? Are you sure to delete all temporary data and all cached SBOMs? Yes
    ✓ temporary data deleted
    ✓ cached SBOMs deleted
```
