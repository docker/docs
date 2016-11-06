<!--[metadata]>
+++
title ="upgrade"
description="Upgrade DTR 2.0.0 or later cluster to this version"
keywords= ["docker, dtr, cli, upgrade"]
[menu.main]
parent="dtr_menu_reference"
identifier="dtr_reference_upgrade"
+++
<![end-metadata]-->

# docker/dtr upgrade

Upgrade DTR 2.0.0 or later cluster to this version

## Usage

```bash
docker run -it --rm docker/dtr \
    upgrade [command options]
```

## Description


This command upgrades DTR 2.0.0 or later to the current version of this image.


## Options

| Option                    | Description                |
|:--------------------------|:---------------------------|
|`--ucp-url`|The UCP URL including domain and port|
|`--ucp-username`|The UCP administrator username|
|`--ucp-password`|The UCP administrator password|
|`--debug`|Enable debug mode for additional logging|
|`--hub-username`|Username to use when pulling images|
|`--hub-password`|Password to use when pulling images|
|`--ucp-insecure-tls`|Disable TLS verification for UCP|
|`--ucp-ca`|Use a PEM-encoded TLS CA certificate for UCP|
|`--existing-replica-id`|The ID of an existing DTR replica|

