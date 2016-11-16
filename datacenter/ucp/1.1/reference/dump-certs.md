---
description: Dump out public certificates
keywords: dump-certs, ucp
title: docker/ucp dump-certs
---

Dump out the public certs for this UCP controller.

## Usage

```bash
docker run --rm \
  --name ucp \
  -v /var/run/docker.sock:/var/run/docker.sock \
  docker/ucp \
  dump-certs [command options]
```

## Description

This utility will dump out the public certificates for the UCP controller
running on the local engine.  This can then be used to populate local
certificate trust stores as desired.

When connecting UCP to DTR, use the output of '--cluster --ca' to
configure DTR.


## Options

| Option        | Description                                                                       |
|:--------------|:----------------------------------------------------------------------------------|
| `--debug, -D` | Enable debug                                                                      |
| `--jsonlog`   | Produce json formatted output for easier parsing                                  |
| `--ca`        | Dump only the contents of the `ca.pem` file (default is to dump both ca and cert) |
| `--cluster`   | Dump the internal UCP Cluster Root CA and cert instead of the public server cert  |