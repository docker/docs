---
description: Regenerate and update TLS certificates
keywords: machine, regenerate-certs, subcommand
title: docker-machine regenerate-certs
hide_from_sitemap: true
---

```none
Usage: docker-machine regenerate-certs [OPTIONS] [arg...]

Regenerate TLS Certificates for a machine

Description:
   Argument(s) are one or more machine names.

Options:

   --force, -f		Force rebuild and do not prompt
   --client-certs	Also regenerate client certificates and CA.
```

Regenerate TLS certificates and update the machine with new certs.

For example:

```bash
$ docker-machine regenerate-certs dev

Regenerate TLS machine certs?  Warning: this is irreversible. (y/n): y
Regenerating TLS certificates
```

If your certificates have expired, you'll need to regenerate the client certs
as well using the `--client-certs` option:

```bash
$ docker-machine regenerate-certs --client-certs dev

Regenerate TLS machine certs?  Warning: this is irreversible. (y/n): y
Regenerating TLS certificates
Regenerating local certificates
...
```
