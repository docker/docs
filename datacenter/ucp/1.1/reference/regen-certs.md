---
description: Regenerate certificates for Docker Universal Control Plane.
keywords: install, ucp
title: docker/ucp regen-certs
---

Regenerate keys and certificates for a UCP controller

## Usage

```bash
$ docker run --rm -it \
     --name ucp \
     -v /var/run/docker.sock:/var/run/docker.sock \
     docker/ucp \
     regen-certs [command options]
```

## Description

This utility will generate new private keys and certs for UCP controllers.

By default it will leave the Root CA keys and certs intact and only
regenerate server and client certs on the controller.  This can be used
to change the list of SANs within the certs after install and refresh
the expiration of the certificates.

You may regenerate the Root CAs with this tool using "--root-ca-only"
then follow a multi-step procedure to regenerate all certs in the cluster.

WARNING: REGENERATING THE ROOT CAs IS A DISRUPTIVE OPERATION!

First run "regen-certs --root-ca-only" on one controller.  If this is an
HA cluster, then perform a "backup --root-ca-only" on this controller,
and "restore --root-ca-only" on all other controllers.  Then on all of
the controllers run "regen-certs" during which the cluster will become
unavailable until 1/2+1 of the controllers are running with new certs.
Once all controllers have new certs, restart all the docker daemons on
the controller nodes.  Once the cluster controllers have recovered, run
"join --fresh-install" on all non-controller nodes to re-join them to
the cluster.  After completing the process, all user bundles will be
invalid and new bundles must be downloaded.

## Options

| Option                                | Description                                                                                 |
|:--------------------------------------|:--------------------------------------------------------------------------------------------|
| `--debug, -D`                         | Enable debug mode                                                                           |
| `--jsonlog`                           | Produce json formatted output for easier parsing                                            |
| `--interactive, -i`                   | Enable interactive mode.  You will be prompted to enter all required information            |
| `--root-ca-only`                      | Regenerate the Root CAs on this node (Do only once in an HA cluster!)                       |
| `--id`                                | The ID of the UCP instance to regenerate certificates for                                   |
| `--san` `[--san option --san option]` | Additional Subject Alternative Names for certs. For example, `--san foo1.bar.com --san foo2.bar.com` |
| `--external-server-cert`              | Omit regenerating the UCP Controller web server certificate signed with an external CA      |