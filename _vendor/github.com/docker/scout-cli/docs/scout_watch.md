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
$ docker scout watch --org my-org --repository registry.example.com/my-service --all-images
```

### Configure Artifactory integration

The following example creates a web hook endpoint for Artifactory to push new
image events into:

```console
$ export DOCKER_SCOUT_ARTIFACTORY_API_USER=user
$ export DOCKER_SCOUT_ARTIFACTORY_API_PASSWORD=password
$ export DOCKER_SCOUT_ARTIFACTORY_WEBHOOK_SECRET=foo

$ docker scout watch --registry "type=artifactory,registry=example.jfrog.io,api=https://example.jfrog.io/artifactory,include=*/frontend*,exclude=*/dta/*,repository=docker-local,port=9000,subdomain-mode=true" --refresh-registry
```

This will launch an HTTP server on port `9000` that will receive all `component` web
hook events, optionally validating the HMAC signature.

### Configure Harbor integration

The following example creates a web hook endpoint for Harbor to push new image
events into:

```console
$ export DOCKER_SCOUT_HARBOR_API_USER=admin
$ export DOCKER_SCOUT_HARBOR_API_PASSWORD=password
$ export DOCKER_SCOUT_HARBOR_WEBHOOK_AUTH="token foo"

$ docker scout watch --registry 'type=harbor,registry=demo.goharbor.io,api=https://demo.goharbor.io,include=*/foo/*,exclude=*/bar/*,port=9000' --refresh-registry
```

This will launch an HTTP server on port `9000` that will receive all `component` web
hook events, optionally validating the HMAC signature.

### Configure Nexus integration

The following example shows how to configure Sonartype Nexus integration:

```console
$ export DOCKER_SCOUT_NEXUS_API_USER=admin
$ export DOCKER_SCOUT_NEXUS_API_PASSWORD=admin124

$ docker scout watch --registry 'type=nexus,registry=localhost:8082,api=http://localhost:8081,include=*/foo/*,exclude=*/bar/*,"repository=docker-test1,docker-test2"' --refresh-registry
```

This ingests all images and tags in Nexus repositories called `docker-test1`
and `docker-test2` that match the `*/foo/*` include and `*/bar/*` exclude glob
pattern.

You can also create a web hook endpoint for Nexus to push new image events into:

```console
$ export DOCKER_SCOUT_NEXUS_API_USER=admin
$ export DOCKER_SCOUT_NEXUS_API_PASSWORD=admin124
$ export DOCKER_SCOUT_NEXUS_WEBHOOK_SECRET=mysecret

$ docker scout watch --registry 'type=nexus,registry=localhost:8082,api=http://localhost:8081,include=*/foo/*,exclude=*/bar/*,"repository=docker-test1,docker-test2",port=9000' --refresh-registry
```

This will launch an HTTP server on port `9000` that will receive all `component` web
hook events, optionally validating the HMAC signature.

## Configure integration for other OCI registries

The following example shows how to integrate an OCI registry that implements the
`_catalog` endpoint:

```console
$ docker scout watch --registry 'type=oci,registry=registry.example.com,include=*/scout-artifact-registry/*'
```