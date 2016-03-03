<!--[metadata]>
+++
title = "Notary Changelog"
description = "Notary release changelog"
keywords = ["docker, notary, changelog, notary changelog, notary releases, releases, notary versions, versions"]
[menu.main]
parent="mn_notary"
weight=99
+++
<![end-metadata]-->

# Changelog

## v0.2
#### 2/24/2016
Adds support for [delegation roles](https://github.com/theupdateframework/tuf/blob/develop/docs/tuf-spec.txt#L387) in TUF.
Delegations allow for easier key management amongst collaborators in a notary trusted collection, and fine-grained permissions on what content each delegate is allowed to modify and sign.
This version also supports managing the snapshot key on notary server, which should be used when enabling delegations on a trusted collection.
Moreover, this version also adds more key management functionality to the notary CLI, and changes the docker-compose development configuration to use the official MariaDB image.

> Detailed release notes can be found here: [v0.2 release notes](https://github.com/docker/notary/releases/tag/v0.2.0).

## v0.1
#### 11/15/2015
Initial notary non-alpha release.
Implements The Update Framework (TUF) with root, targets, snapshot, and timestamp roles to sign and verify content of a trusted collection.

> Detailed release notes can be found here: [v0.1 release notes](https://github.com/docker/notary/releases/tag/v0.1).
