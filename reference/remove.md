<!--[metadata]>
+++
title ="remove"
description="Remove a DTR replica"
keywords= ["docker, dtr, cli, remove"]
[menu.main]
parent="dtr_menu_reference"
identifier="dtr_reference_remove"
+++
<![end-metadata]-->

# docker/dtr remove

Remove a DTR replica

## Usage

```bash
docker run -it --rm docker/dtr \
    remove [command options]
```

## Description


This command removes a replica from a DTR deployment. All DTR containers and
volumes are removed from the node.


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
|`--force-remove`|Force a DTR replica to be removed|
|`--replica-id`|Assign an ID to the DTR replica. By default the ID is random|
|`--existing-replica-id`|The ID of an existing DTR replica|

