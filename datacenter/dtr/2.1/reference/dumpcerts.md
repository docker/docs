<!--[metadata]>
+++
title ="dumpcerts"
description="Print the TLS certificates used by DTR"
keywords= ["docker, dtr, cli, dumpcerts"]
[menu.main]
parent="dtr_menu_reference"
identifier="dtr_reference_dumpcerts"
+++
<![end-metadata]-->

# docker/dtr dumpcerts

Print the TLS certificates used by DTR

## Usage

```bash
docker run -i --rm docker/dtr \
    dumpcerts [command options] > backup.tar
```

## Description


This command creates a backup of the certificates used by DTR for
communicating across replicas with TLS.


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

