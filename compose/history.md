---
title: Evolution of Compose
description: History of Compose V1 and Compose YAML schema versioning
keywords: compose, compose yaml, swarm, migration, compatibility
redirect_from:
 - /compose/cli-command-compatibility/
 - /compose/compose-v2/
---
{% include compose-eol.md %}

The Compose project has existed for many years now, so this page details a brief history of the major versions and file format revisions.

The currently supported version of Compose is V2, which uses YAML files as defined by the [Compose specification](https://github.com/compose-spec/compose-spec){:
target="_blank" rel="noopener" class="_"}.

For users supporting projects originally targeting older versions of Compose, this can serve as a guide to understanding compatibility and the evolution of changes.

## Docker Compose CLI versioning
There are two major versions of Docker Compose, i.e. the command-line binary.

Compose V1 was first released in 2014, is written in Python, and is invoked as `docker-compose`.
Typically, Compose V1 projects include a `version` field in YAML with values ranging from `2.0` to `3.8`.

Compose V2 was announced in 2020, is written in Go, and is invoked as `docker compose`.
Compose V2 ignores the `version` field in YAML.

## Compose file format versioning
With Compose V1, projects declared a Compose file format version in YAML.

Three major versions of the Compose file format for Compose V1 were released:
- [Compose file format 1](/compose/compose-file/compose-versioning/#version-1-to-2x) with Compose 1.0.0 in 2014
- [Compose file format 2.x](../compose-file/compose-file-v2/) with Compose 1.6.0 in 2016
- [Compose file format 3.x](../compose-file/compose-file-v3/) with Compose 1.10.0 in 2017

Compose file format 1 was substantially different than all following formats, lacking a top-level `services` key.
Its usage is historical and files written in this format will not run with Compose V2.

Compose file format 2.x and 3.x were very similar to each other, but the latter introduced many new options targeted at Swarm deployments.

To address confusion around Compose CLI versioning, Compose file format versioning, and feature parity depending on whether Swarm mode was in use, file format 2.x and 3.x were merged into the [Compose specification](https://github.com/compose-spec/compose-spec){:target="_blank" rel="noopener" class="_"}.
Unlike the prior file formats, the Compose specification is rolling and eliminates the `version` field in YAML.

Compose V2 uses the Compose specification for project definition.
The `version` field should be omitted and will be ignored if present.

To make migration easier, Compose V2 has backwards compatibility for certain elements that have been deprecated or changed between Compose file format 2.x/3.x and the Compose specification.
In these cases, a warning will be logged when running Compose V2 commands, and you should update your YAML accordingly.
Future versions of Compose might begin to return an error in these cases.
