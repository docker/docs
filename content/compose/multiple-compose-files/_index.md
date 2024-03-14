---
description: General overview for the different ways you can work with multiple compose
  files in Docker Compose
keywords: compose, compose file, merge, extends, include, docker compose, -f flag
title: Overview
---

{{< include "compose-eol.md" >}}

This section contains information on the ways you can work with multiple Compose files. 

Using multiple Compose files lets you customize a Compose application for different environments or workflows. This is useful for large applications that may use dozens of containers, with ownership distributed across multiple teams. For example, if your organization or team uses a monorepo, each team may have their own “local” Compose file to run a subset of the application. They then need to rely on other teams to provide a reference Compose file that defines the expected way to run their own subset. Complexity moves from the code in to the infrastructure and the configuration file.

The quickest way to work with multiple Compose files is to [merge](merge.md) Compose files using the `-f` flag in the command line to list out your desired Compose files. However, [merging rules](merge.md#merging-rules) means this can soon get quite complicated.

Docker Compose provides two other options to manage this complexity when working with multiple Compose files. Depending on your project's needs, you can: 

- [Extend a Compose file](extends.md) by referring to another Compose file and selecting the bits you want to use in your own application, with the ability to override some attributes.
- [Include other Compose files](include.md) directly in your Compose file.

