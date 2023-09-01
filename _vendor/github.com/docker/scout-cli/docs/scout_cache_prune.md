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

Docker Scout stores temporary files to generate SBOMs of images and cache those SBOMs to not generate or fetch them again.

The `docker scout cache prune` command will remove all the temporary files used while generating the SBOMs.

By default the cached SBOMs will not be deleted are they can be used by the different `docker scout` commands. But the `--sboms`
flag can be used to delete them.

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
