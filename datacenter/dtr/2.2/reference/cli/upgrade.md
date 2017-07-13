---
title: docker/dtr upgrade
keywords: docker, dtr, cli, upgrade
description: Upgrade DTR 2.0.0 or later cluster to this version
---

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
|`--debug`|Enable debug mode for additional logging|
|`--existing-replica-id`|The ID of an existing DTR replica|
|`--hub-password`|Password to use when pulling images|
|`--hub-username`|Username to use when pulling images|
|`--ucp-ca`|Use a PEM-encoded TLS CA certificate for UCP|
|`--ucp-insecure-tls`|Disable TLS verification for UCP|
|`--ucp-password`|The UCP administrator password|
|`--ucp-url`|The UCP URL including domain and port|
|`--ucp-username`|The UCP administrator username|
|`--unsafe-upgrade`|Perform the upgrade ignoring version checks|

