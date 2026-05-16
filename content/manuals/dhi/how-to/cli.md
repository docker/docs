---
title: Use the DHI CLI
linkTitle: Use the CLI
weight: 50
keywords: docker dhi, CLI, command line, docker hardened images
description: Learn how to install and use docker dhi, the command-line interface for managing Docker Hardened Images.
---

The `docker dhi` command-line interface (CLI) is a tool for managing Docker Hardened Images:
- Browse the catalog of available DHI images and their metadata
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

```bash
docker dhi --help
docker dhi catalog list --help
```

### Browse the DHI catalog

List all available DHI images:

```bash
docker dhi catalog list
```

Filter by type, name, or compliance:

```bash
docker dhi catalog list --type image
docker dhi catalog list --filter golang
docker dhi catalog list --fips
docker dhi catalog list --stig
```

Get details of a specific image, including available tags and CVE counts:

```bash
docker dhi catalog get <image-name>
```

### Mirror DHI images

{{< summary-bar feature_name="Docker Hardened Images" >}}

Start mirroring one or more DHI images to your Docker Hub organization:

```bash
docker dhi mirror start --org my-org \
  -r dhi/golang,my-org/dhi-golang \
  -r dhi/nginx,my-org/dhi-nginx \
  -r dhi/prometheus-chart,my-org/dhi-prometheus-chart
```

Mirror with dependencies:

```bash
docker dhi mirror start --org my-org -r dhi/golang,my-org/dhi-golang --dependencies
```

List mirrored images in your organization:

```bash
docker dhi mirror list --org my-org
```

Filter mirrored images by name or type:

```bash
docker dhi mirror list --org my-org --filter python
docker dhi mirror list --org my-org --type image
docker dhi mirror list --org my-org --type helm-chart
```

Stop mirroring one or more images:

```bash
docker dhi mirror stop dhi-golang --org my-org
docker dhi mirror stop dhi-python dhi-golang --org my-org
```

Stop mirroring and delete the repositories:

```bash
docker dhi mirror stop dhi-golang --org my-org --delete
docker dhi mirror stop dhi-golang --org my-org --delete --force
```

### Customize DHI images

{{< summary-bar feature_name="Docker Hardened Images" >}}

The CLI can be used to create and manage DHI image customizations. For detailed
instructions on creating customizations using the GUI, see [Customize a Docker
Hardened Image](./customize.md).

The following is a quick reference for CLI commands. For complete details on all
options and flags, see the
[CLI reference](/reference/cli/docker/dhi/).

```bash
# Prepare a customization scaffold
docker dhi customization prepare golang 1.25 \
  --org my-org \
  --destination my-org/dhi-golang \
  --name "golang with git" \
  --output my-customization.yaml

# Create a customization
docker dhi customization create my-customization.yaml --org my-org

# List customizations
docker dhi customization list --org my-org

# Filter customizations by name, repository, or source
docker dhi customization list --org my-org --filter git
docker dhi customization list --org my-org --repo dhi-golang
docker dhi customization list --org my-org --source golang

# Get a customization
docker dhi customization get my-org/dhi-golang "golang with git" --org my-org --output my-customization.yaml

# Update a customization
# The YAML file must include the 'id' field to identify the customization to update
docker dhi customization edit my-customization.yaml --org my-org

# Delete a customization
docker dhi customization delete my-org/dhi-golang "golang with git" --org my-org

# Delete without confirmation prompt
docker dhi customization delete my-org/dhi-golang "golang with git" --org my-org --yes
```

### Enterprise package authentication

{{< summary-bar feature_name="Docker Hardened Images Enterprise" >}}

Generate authentication credentials for accessing the enterprise hardened
package repository. This is used when configuring your package manager to
install compliance-specific packages in your own images. For detailed
instructions, see [Enterprise
repository](./hardened-packages.md#enterprise-repository).

```bash
docker dhi auth apk
```

### Monitor customization builds

{{< summary-bar feature_name="Docker Hardened Images" >}}

List builds for a customization:

```bash
docker dhi customization build list my-org/dhi-golang "golang with git" --org my-org
docker dhi customization build list my-org/dhi-golang "golang with git" --org my-org --json
```

Get details of a specific build:

```bash
docker dhi customization build get my-org/dhi-golang "golang with git" <build-id> --org my-org
docker dhi customization build get my-org/dhi-golang "golang with git" <build-id> --org my-org --json
```

View build logs:

```bash
docker dhi customization build logs my-org/dhi-golang "golang with git" <build-id> --org my-org
docker dhi customization build logs my-org/dhi-golang "golang with git" <build-id> --org my-org --json
```

### JSON output

Most list and get commands support a `--json` flag for machine-readable output:

```bash
docker dhi catalog list --json
docker dhi catalog get golang --json
docker dhi mirror list --org my-org --json
docker dhi mirror start --org my-org -r golang --json
docker dhi customization list --org my-org --json
docker dhi customization build list my-org/dhi-golang "golang with git" --org my-org --json
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
