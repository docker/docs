---
description: Learn how to upgrade your Docker Trusted Registry
keywords: docker, dtr, upgrade, install
title: Upgrade DTR
---

DTR uses [semantic versioning](http://semver.org/) and we aim to achieve specific
guarantees while upgrading between versions. We never support downgrading. We
support upgrades according to the following rules:

* When upgrading from one patch version to another you can skip patch versions
  because no data migraiton is done for patch versions.
* When upgrading between minor versions, you can't skip versions, but you can
  upgrade from any patch versions of the previous minor version to any patch
  version of the current minor version.
* When upgrading between major versions you also have to upgrade one major
  version at a time, but you have to upgrade to the earliest available minor
  version. We also strongly recommend upgrading to the latest minor/patch
  version for your major version first.

|From| To| Description| Supported|
|:----|:---|:------------|----------|
| 2.2.0 | 2.2.1 | patch upgrade | yes |
| 2.2.0 | 2.2.2 | skip patch version | yes |
| 2.2.2 | 2.2.1 | patch downgrade | no |
| 2.1.0 | 2.2.0 | minor upgrade | yes |
| 2.1.1 | 2.2.0 | minor upgrade | yes |
| 2.1.2 | 2.2.2 | minor upgrade | yes |
| 2.0.1 | 2.2.0 | skip minor version | no |
| 2.2.0 | 2.1.0 | minor downgrade | no |
| 1.4.3 | 2.0.0 | major upgrade | yes |
| 1.4.3 | 2.0.3 | major upgrade | yes |
| 1.4.3 | 3.0.0 | skip major version | no |
| 1.4.1 | 2.0.3 | major upgrade from an old version | no |
| 1.4.3 | 2.1.0 | major upgrade skipping minor version | no |
| 2.0.0 | 1.4.3 | major downgrade | no |

There may be at most a few seconds of interruption during the upgrade of a
DTR cluster. Schedule the upgrade to take place outside business peak hours
to ensure the impact on your business is close to none.

## Minor upgrade

Before starting your upgrade planning, make sure that the version of UCP you are
using is supported by the version of DTR you are trying to upgrade to. <!--(TODO:
link to the compatibility matrix)-->

### Step 1. Upgrade DTR to 2.1 if necessary

Make sure you're running DTR 2.1. If that's not the case, [upgrade your installation to the 2.1 version](/datacenter/dtr/2.1/guides/install/upgrade/.md).

### Step 2. Upgrade DTR

Then pull the latest version of DTR:

```none
$ docker pull {{ page.docker_image }}
```

If the node you're upgrading doesn't have access to the internet, you can
follow the [offline installation documentation](../install/install-offline.md)
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
[reference documentation](../../../reference/cli/index.md) for other existing flags.

The upgrade command will start replacing every container in your DTR cluster,
one replica at a time. It will also perform certain data migrations. If anything
fails or the upgrade is interrupted for any reason, you can re-run the upgrade
command and it will resume from where it left off.

## Patch upgrade

A patch upgrade changes only the DTR containers and it's always safer than a minor
upgrade. The command is the same as for a minor upgrade.

## Where to go next

* [Release notes](release-notes.md)
