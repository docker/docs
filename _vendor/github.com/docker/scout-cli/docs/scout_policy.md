# docker scout policy

<!---MARKER_GEN_START-->
Evaluate policies against an image and display the policy evaluation results (experimental)

### Options

| Name                | Type     | Default | Description                                                 |
|:--------------------|:---------|:--------|:------------------------------------------------------------|
| `--env`             | `string` |         | Name of the environment to compare to                       |
| `-e`, `--exit-code` |          |         | Return exit code '2' if policies are not met, '0' otherwise |
| `--org`             | `string` |         | Namespace of the Docker organization                        |
| `-o`, `--output`    | `string` |         | Write the report to a file                                  |
| `--platform`        | `string` |         | Platform of image to pull policy results from               |


<!---MARKER_GEN_END-->

## Description

The `docker scout policy` command evaluates policies against an image. The image is indexed into the Scout platform - if it wasn't already - and policies will run against it. The policy results may take a few minutes to become available.

## Examples

### Evaluate policies against an image and display the results

```console
$ docker scout policy dockerscoutpolicy/customers-api-service:0.0.1
```

### Evaluate policies against an image for a specific organization

```console
$ docker scout policy dockerscoutpolicy/customers-api-service:0.0.1 --org dockerscoutpolicy
```

### Evaluate policies against an image with a specific platform

```console
$ docker scout policy dockerscoutpolicy/customers-api-service:0.0.1 --platform linux/amd64
```

### Compare policy results for a repository in a specific environment

```console
$ docker scout policy dockerscoutpolicy/customers-api-service --env production
```
