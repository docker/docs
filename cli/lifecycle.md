---
title: Docker CLI release lifecycle
description: The product release lifecycle for the Docker CLI
keywords: cli, lifecycle, deprecation, removal, experimental features
---

This page describes the product lifecycle management of features in the Docker
CLI.

## CLI release cadence

The Docker CLI releases at the same cadence as the Docker Engine. The Docker
CLI is distributed as a separate package, but it's release always coincides
with the release of Docker Engine.

## Deprecation and removal

When a CLI command or flag becomes obsolete, replaced, or discontinued, it
first gets flagged as deprecated. When a command or flag is deprecated, it
means that you're advised to refrain from using it.

CLI features that are flagged as deprecated may come to be removed in the next
major release following their deprecation.

## Experimental features

Experimental CLI features provide early access to future product functionality.
Such features are intended for testing and feedback, and they may change
between releases, or be removed entirely, without warning.

Starting with Docker CLI version 20.10, experimental CLI features are enabled
by default. No configuration is required to enable them.
