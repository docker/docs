# docker scout policy

<!---MARKER_GEN_START-->
Display the policy results of an image (experimental)

### Options

| Name                | Type     | Default | Description                                    |
|:--------------------|:---------|:--------|:-----------------------------------------------|
| `--env`             | `string` |         | Name of the environment to compare to.         |
| `-e`, `--exit-code` |          |         | Return exit code '2' if policies are not met.  |
| `--org`             | `string` |         | Namespace of the Docker organization           |
| `-o`, `--output`    | `string` |         | Write the report to a file.                    |
| `--platform`        | `string` |         | Platform of image to pull policy results from. |


<!---MARKER_GEN_END-->

## Description

The `docker scout policy` command displays the policy results of an image if there are any.

## Examples

### Display the policy results of an image

```console
$ docker scout policy dockerscoutpolicy/customers-api-service:0.0.1
```

### Compare policy results for a repository in a specific environment

```console
$ docker scout policy dockerscoutpolicy/customers-api-service --env production
```
