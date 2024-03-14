# docker scout config

<!---MARKER_GEN_START-->
Manage Docker Scout configuration


<!---MARKER_GEN_END-->

## Description

`docker scout config` allows you to list, get and set Docker Scout configuration.

Available configuration key:

- `organization`: Namespace of the Docker organization to be used by default.

## Examples

### List existing configuration

```console
$ docker scout config
organization=my-org-namespace
```

### Print configuration value

```console
$ docker scout config organization
my-org-namespace
```

### Set configuration value

```console
$ docker scout config organization my-org-namespace
    ✓ Successfully set organization to my-org-namespace
```
