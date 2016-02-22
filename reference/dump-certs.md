+++
title = "dump-certs"
keywords= ["dump-certs, ucp"]
description = "Dump out public certificates"
[menu.main]
parent ="ucp_ref"
+++

# dump-certs

Dump out the public certs for this UCP controller.

## Usage

```
docker run --rm \
     --name ucp \
     -v /var/run/docker.sock:/var/run/docker.sock \
     docker/ucp \
     dump-certs [command options]
```

## Description

Dumps out the public certificates for the UCP controller running on the local
engine. By default, this command dumps both the CA and certificate. You can use
the output of this command to populate local certificate trust stores as
desired.


## Options

| Option                | Description                                                                         |
|-----------------------|-------------------------------------------------------------------------------------|
| `--debug`, `-D`       | Enable debug                                                                        |
| `--jsonlog`           | Produce json formatted output for easier parsing                                    |
| `--interactive`, `-i` | Enable interactive mode.,You are prompted to enter all required information.        |
| `--ca`                | Dump only the contents of the `ca.pem` file (default is to dump both ca and cert).  |
| `--cluster`           | Dump the internal UCP Cluster Root CA and cert instead of the public server cert.   |
