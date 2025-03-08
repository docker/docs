---
title: Configure your Notary client
description: Learn how to configure your Notary client to push and pull images from Docker Trusted Registry.
keywords: docker, registry, notary, trust
---

The Docker CLI client makes it easy to sign images but to streamline that
process it generates a set of private and public keys that are not tied
to your UCP account. This means that you can push and sign images to
DTR, but UCP doesn't trust those images since it doesn't know anything about
the keys you're using.

So before signing and pushing images to DTR you should:

* Configure the Notary CLI client
* Import your UCP private keys to the Notary client

This allows you to  start signing images with the private keys in your UCP
client bundle, that UCP can trace back to your user account.

## Download the Notary CLI client

If you're using Docker for Mac or Docker for Windows, you already have the
`notary` command installed.

If you're running Docker on a Linux distribution, you can [download the
latest version](https://github.com/docker/notary/releases). As an example:

```bash
# Get the latest binary
curl -L <download-url> -o notary

# Make it executable
chmod +x notary

# Move it to a location in your path
sudo mv notary /usr/bin/
```

## Configure the Notary CLI client

Before you use the Notary CLI client, you need to configure it to make it
talk with the Notary server that's part of DTR.

There's two ways to do this, either by passing flags to the notary command,
or using a configuration file.

### With flags

Run the Notary command with:

```bash
notary --server https://<dtr-url> --trustDir ~/.docker/trust --tlscacert <dtr-ca.pem>
```

Here's what the flags mean:

| Flag          | Purpose                                                                                                                           |
|:--------------|:----------------------------------------------------------------------------------------------------------------------------------|
| `--server`    | The Notary server to query                                                                                                        |
| `--trustDir`  | Path to the local directory where trust metadata will be stored                                                                   |
| `--tlscacert` | Path to the DTR CA certificate. If you've configured your system to trust the DTR CA certificate, you don't need to use this flag |

To avoid having to type all the flags when using the command, you can set an
alias:

```none
# Bash
alias notary="notary --server https://<dtr-url> --trustDir ~/.docker/trust --tlscacert <dtr-ca.pem>"

# PowerShell
set-alias notary "notary --server https://<dtr-url> --trustDir ~/.docker/trust --tlscacert <dtr-ca.pem>"
```

### With a configuration file

You can also configure Notary by creating a `~/.notary/config.json` file with
the following content:

```json
{
  "trust_dir" : "~/.docker/trust",
  "remote_server": {
    "url": "<dtr-url>",
    "root_ca": "<dtr-ca.pem>"
  }
}
```

To validate your configuration, try running the `notary list` command on a
DTR repository that already has signed images:

```none
# Assumes you've configured notary
notary list <dtr-repository>
```

The command should print a list of digests for each signed image on the
repository.

## Import your UCP key

The last step in configuring the Notary CLI client is to import the private
key of your UCP client bundle.
[Get a new client bundle if you don't have one yet](/datacenter/ucp/2.1/guides/user/access-ucp/cli-based-access.md).

Import the private key in your UCP bundle into the Notary CLI client:

```none
# Assumes you've configured notary
notary key import <path-to-key.pem>
```

The private key is copied to `~/.docker/trust`, and you are prompted for a
password to encrypt it.

You can validate what keys Notary knows about by running:

```none
notary key list
```

The key you've imported should be listed with the role `delegation`.
