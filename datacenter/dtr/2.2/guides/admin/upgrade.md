---
description: Learn how to upgrade your Docker Trusted Registry
keywords: docker, dtr, upgrade, install
title: Upgrade DTR
---

DTR uses [semantic versioning](http://semver.org/) and we aim to achieve specific
guarantees while upgrading between versions. We never support downgrading. We
support upgrades according to the following rules:

* When upgrading from one patch version to another you can skip patch versions
  because no data migration is done for patch versions.
* When upgrading between minor versions, you can't skip versions, but you can
  upgrade from any patch versions of the previous minor version to any patch
  version of the current minor version.
* When upgrading between major versions you also need to upgrade one major
  version at a time, but you need to upgrade to the earliest available minor
  version. We also strongly recommend upgrading to the latest minor/patch
  version for your major version first.

| Description                          | From  | To        | Supported |
|:-------------------------------------|:------|:----------|:----------|
| patch upgrade                        | x.y.0 | x.y.1     | yes       |
| skip patch version                   | x.y.0 | x.y.2     | yes       |
| patch downgrade                      | x.y.2 | x.y.1     | no        |
| minor upgrade                        | x.y.* | x.y+1.*   | yes       |
| skip minor version                   | x.y.* | x.y+2.*   | no        |
| minor downgrade                      | x.y.* | x.y-1.*   | no        |
| skip major version                   | x.*.* | x+2.*.*   | no        |
| major downgrade                      | x.*.* | x-1.*.*   | no        |
| major upgrade                        | x.y.z | x+1.0.0   | yes       |
| major upgrade skipping minor version | x.y.z | x+1.y+1.z | no        |

There may be at most a few seconds of interruption during the upgrade of a
DTR cluster. Schedule the upgrade to take place outside business peak hours
to ensure the impact on your business is close to none.

## Minor upgrade

Before starting your upgrade planning, make sure that the version of UCP you are
using is supported by the version of DTR you are trying to upgrade to. This can be
checked using the [Compatibility Matrix](https://success.docker.com/Policies/Compatibility_Matrix).

> **Warning**
>
> Before performing any upgrade it’s important to backup. See
> [DTR backups and recovery](/datacenter/dtr/2.2/guides/admin/backups-and-disaster-recovery.md).
{: .warning}

### Step 1. Upgrade DTR to 2.1 if necessary

Make sure you're running DTR 2.1. If that's not the case, [upgrade your installation to the 2.1 version](/datacenter/dtr/2.1/guides/install/upgrade/.md).

### Step 2. Upgrade DTR

Then pull the latest version of DTR:

```none
$ docker pull {{ page.docker_image }}
```

If the node you're upgrading doesn't have access to the internet, you can
follow the [offline installation documentation](install/install-offline.md)
to get the images.

Once you have the latest image on your machine (and the images on the target
nodes if upgrading offline), run the upgrade command:

```none
$ docker run -it --rm \
  {{ page.docker_image }} upgrade \
  --ucp-insecure-tls
```

By default the upgrade command runs in interactive mode and prompts you for
any necessary information. You can also check the
[reference documentation](../../reference/cli/index.md) for other existing flags.

The upgrade command will start replacing every container in your DTR cluster,
one replica at a time. It will also perform certain data migrations. If anything
fails or the upgrade is interrupted for any reason, you can re-run the upgrade
command and it will resume from where it left off.

## Patch upgrade

A patch upgrade changes only the DTR containers and it's always safer than a minor
upgrade. The command is the same as for a minor upgrade.

## Where to go next

* [Release notes](release-notes.md)
