---
title: Upgrade DTR
description: Learn how to upgrade your Docker Trusted Registry
keywords: dtr, upgrade, install
---

{% assign previous_version="2.4" %}

DTR uses [semantic versioning](http://semver.org/) and we aim to achieve specific
guarantees while upgrading between versions. We never support downgrading. We
support upgrades according to the following rules:

* When upgrading from one patch version to another you can skip patch versions
  because no data migration is done for patch versions.
* When upgrading between minor versions, you can't skip versions, but you can
  upgrade from any patch versions of the previous minor version to any patch
  version of the current minor version.
* When upgrading between major versions you also have to upgrade one major
  version at a time, but you have to upgrade to the earliest available minor
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

Before starting your upgrade, make sure that:
* The version of UCP you are using is supported by the version of DTR you
are trying to upgrade to. [Check the compatibility matrix](https://success.docker.com/Policies/Compatibility_Matrix).
* You have a recent [DTR backup](disaster-recovery/create-a-backup.md).
* You [disable Docker content trust in UCP](/datacenter/ucp/2.2/guides/admin/configure/run-only-the-images-you-trust.md).

### Step 1. Upgrade DTR to {{ previous_version }} if necessary

Make sure you're running DTR {{ previous_version }}. If that's not the case,
[upgrade your installation to the {{ previous_version }} version](/datacenter/dtr/{{ previous_version }}/guides/admin/upgrade.md).

### Step 2. Upgrade DTR

Then pull the latest version of DTR:

```bash
docker pull {{ page.dtr_org }}/{{ page.dtr_repo }}:{{ page.dtr_version }}
```

If the node you're upgrading doesn't have access to the internet, you can
follow the [offline installation documentation](install/install-offline.md)
to get the images.

Once you have the latest image on your machine (and the images on the target
nodes if upgrading offline), run the upgrade command:

```bash
docker run -it --rm \
  {{ page.dtr_org }}/{{ page.dtr_repo }}:{{ page.dtr_version }} upgrade \
  --ucp-insecure-tls
```

By default the upgrade command runs in interactive mode and prompts you for
any necessary information. You can also check the
[reference documentation](/reference/dtr/2.5/cli/index.md) for other existing flags.

The upgrade command will start replacing every container in your DTR cluster,
one replica at a time. It will also perform certain data migrations. If anything
fails or the upgrade is interrupted for any reason, you can re-run the upgrade
command and it will resume from where it left off.

## Patch upgrade

A patch upgrade changes only the DTR containers and it's always safer than a minor
upgrade. The command is the same as for a minor upgrade.

## Download the vulnerability database

After upgrading DTR, you need to re-download the vulnerability database.
[Learn how to update your vulnerability database](configure/set-up-vulnerability-scans.md#update-the-cve-scanning-database).

## Where to go next

- [Release notes](../release-notes.md)
