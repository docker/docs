---
title: docker/ucp dump-certs
description: Print the public certificates used by this UCP web server
keywords: docker, ucp, cli, dump-certs
---

Print the public certificates used by this UCP web server

## Description

This command outputs the public certificates for the UCP web server running on
this node. By default it prints the contents of the ca.pem and cert.pem files.

When integrating UCP and DTR, use this command with the `--cluster --ca` flags
to configure DTR.


## Options

| Option                    | Description                |
|:--------------------------|:---------------------------|
|`--debug, D`|Enable debug mode|
|`--jsonlog`|Produce json formatted output for easier parsing|
|`--ca`|Only print the contents of the ca.pem file|
|`--cluster`|Print the internal UCP swarm root CA and cert instead of the public server cert|
