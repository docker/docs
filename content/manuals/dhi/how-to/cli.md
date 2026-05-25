---
title: Use the DHI CLI
linkTitle: Use the CLI
weight: 50
keywords: docker dhi, CLI, command line, docker hardened images
description: Learn how to install and use docker dhi, the command-line interface for managing Docker Hardened Images.
---

The `docker dhi` command-line interface (CLI) is a tool for managing Docker Hardened Images:
- Browse the catalog of available DHI images and their metadata
- View attestations for DHI images, including SBOMs and provenance
- Mirror DHI images to your Docker Hub organization
- Create and manage customizations of DHI images
- Generate authentication for enterprise package repositories
- Monitor customization builds

## Installation

The `docker dhi` CLI is available in [Docker Desktop](https://docs.docker.com/desktop/) version 4.65 and later.
You can also install the standalone `dhictl` binary.

### Docker Desktop

The `docker dhi` command is included in Docker Desktop 4.65 and later. No additional installation is required.

### Standalone binary

1. Download the `dhictl` binary for your platform from the
   [releases](https://github.com/docker-hardened-images/dhictl/releases) page.
2. Move it to a directory in your `PATH`:
    - `mv dhictl /usr/local/bin/` on _Linux_ and _macOS_
    - Move `dhictl.exe` to a directory in your `PATH` on _Windows_

## Usage

Every command has built-in help accessible with the `--help` flag:

```console
$ docker dhi --help
$ docker dhi catalog list --help
```

### Browse the DHI catalog

List all available DHI images:

```console
$ docker dhi catalog list
```

Filter by type, name, or compliance:

```console
$ docker dhi catalog list --type image
$ docker dhi catalog list --filter golang
$ docker dhi catalog list --fips
$ docker dhi catalog list --stig
```

Get details of a specific image, including available tags and CVE counts:

```console
$ docker dhi catalog get <image-name>
```

### View attestations

List all attestations attached to a DHI image:

```console
$ docker dhi attestation list dhi/nginx:1.27
$ docker dhi attestation list dhi/nginx:1.27 --platform linux/amd64
$ docker dhi attestation list dhi/nginx:1.27 --predicate-type https://slsa.dev/provenance/v1
$ docker dhi attestation list dhi/nginx:1.27 --json
```

Get a specific attestation by its referrer digest:

```console
$ docker dhi attestation get dhi/nginx:1.27 sha256:<digest>
$ docker dhi attestation get dhi/nginx:1.27 sha256:<digest> -o provenance.json
```

Display the SPDX SBOM for an image:

```console
$ docker dhi attestation sbom dhi/nginx:1.27
$ docker dhi attestation sbom dhi/nginx:1.27 --platform linux/amd64
```

### Mirror DHI images

{{< summary-bar feature_name="Docker Hardened Images" >}}

Start mirroring one or more DHI images to your Docker Hub organization:

```console
$ docker dhi mirror start --org my-org \
  dhi/golang,my-org/dhi-golang \
  dhi/nginx,my-org/dhi-nginx \
  dhi/prometheus-chart,my-org/dhi-prometheus-chart
```

Mirror with dependencies:

```console
$ docker dhi mirror start --org my-org dhi/golang,my-org/dhi-golang --dependencies
```

List mirrored images in your organization:

```console
$ docker dhi mirror list --org my-org
```

Filter mirrored images by name or type:

```console
$ docker dhi mirror list --org my-org --filter python
$ docker dhi mirror list --org my-org --type image
$ docker dhi mirror list --org my-org --type helm-chart
```

Stop mirroring one or more images:

```console
$ docker dhi mirror stop dhi-golang --org my-org
$ docker dhi mirror stop dhi-python dhi-golang --org my-org
```

Stop mirroring and delete the repositories:

```console
$ docker dhi mirror stop dhi-golang --org my-org --delete
$ docker dhi mirror stop dhi-golang --org my-org --delete --force
```

### Customize DHI images

{{< summary-bar feature_name="Docker Hardened Images" >}}

The CLI can be used to create and manage DHI image customizations. For detailed
instructions on creating customizations using the GUI, see [Customize a Docker
Hardened Image](./customize.md).

The following is a quick reference for CLI commands. For complete details on all
options and flags, see the
[CLI reference](/reference/cli/docker/dhi/).

```console
# Prepare a single customization scaffold
$ docker dhi customization prepare golang 1.25 \
  --org my-org \
  --destination my-org/dhi-golang \
  --name "golang with git" \
  > my-customization.yaml

# Prepare a bulk customization scaffold (pipe JSON array via stdin)
$ echo '[{"destination":"my-org/dhi-golang","tag-definition-id":"golang/alpine-3.23/1.24-dev"}]' \
  | docker dhi customization prepare --name "golang with git" --org my-org \
  > my-customization.yaml

# Create a customization
$ docker dhi customization create my-customization.yaml --org my-org

# Create with flag overrides (flags take precedence over the YAML file)
$ docker dhi customization create my-customization.yaml --org my-org \
  --destination my-org/dhi-golang \
  --name "golang with git"

# List customizations
$ docker dhi customization list --org my-org

# Filter customizations by name, repository, or source
$ docker dhi customization list --org my-org --filter git
$ docker dhi customization list --org my-org --repo dhi-golang
$ docker dhi customization list --org my-org --source golang

# Get a customization by ID
$ docker dhi customization get <id> --org my-org

# Update a customization
# The YAML file must include the 'id' field to identify the customization to update
$ docker dhi customization edit my-customization.yaml --org my-org

# Delete a customization by ID
$ docker dhi customization delete <id> --org my-org

# Delete multiple customizations
$ docker dhi customization delete <id1> <id2> --org my-org

# Delete without confirmation prompt
$ docker dhi customization delete <id> --org my-org --force
```

For a complete reference of all YAML fields, see
[Image customization YAML file](/dhi/how-to/customize/#image-customization-yaml-file).

### Enterprise package authentication

{{< summary-bar feature_name="Docker Hardened Images Enterprise" >}}

Generate authentication credentials for accessing the enterprise hardened
package repository. These credentials are used when configuring your package
manager to install compliance and security-patched packages in your own images. For detailed
instructions, see [Enterprise
repository](./hardened-packages.md#enterprise-repository).

For Alpine-based images:

```console
$ docker dhi auth apk
```

For Debian-based images:

```console
$ docker dhi auth deb
```

### Monitor customization builds

{{< summary-bar feature_name="Docker Hardened Images" >}}

List builds for a customization:

```console
$ docker dhi customization build list <customization-id> --org my-org
$ docker dhi customization build list <customization-id> --org my-org --json
```

Get details of a specific build:

```console
$ docker dhi customization build get <customization-id> <build-id> --org my-org
$ docker dhi customization build get <customization-id> <build-id> --org my-org --json
```

View build logs:

```console
$ docker dhi customization build logs <customization-id> <build-id> --org my-org
$ docker dhi customization build logs <customization-id> <build-id> --org my-org --json
```

### JSON output

Most list and get commands support a `--json` flag for machine-readable output:

```console
$ docker dhi catalog list --json
$ docker dhi catalog get golang --json
$ docker dhi attestation list dhi/nginx:1.27 --json
$ docker dhi mirror list --org my-org --json
$ docker dhi mirror start --org my-org dhi/golang,my-org/dhi-golang --json
$ docker dhi customization list --org my-org --json
$ docker dhi customization build list <customization-id> --org my-org --json
```

## Configuration

The `docker dhi` CLI can be configured with a YAML file located at:
- `$HOME/.config/dhictl/config.yaml` on _Linux_ and _macOS_
- `%USERPROFILE%\.config\dhictl\config.yaml` on _Windows_

If `$XDG_CONFIG_HOME` is set, the configuration file is located at `$XDG_CONFIG_HOME/dhictl/config.yaml` (see the [XDG Base Directory Specification](https://specifications.freedesktop.org/basedir/spec/latest/)).

Available configuration options:

| Option      | Environment Variable | Description                                                                                                               |
|-------------|----------------------|---------------------------------------------------------------------------------------------------------------------------|
| `org`       | `DHI_ORG`            | Default Docker Hub organization for mirror and customization commands.                                                    |
| `api_token` | `DHI_API_TOKEN`      | Docker token for authentication. You can generate a token in your [Docker Hub account settings](https://hub.docker.com/). |

Environment variables take precedence over configuration file values.
