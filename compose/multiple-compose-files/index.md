---
description: General overview for the different ways you can work with multiple compose files in Docker Compose
keywords: compose, compose file, merge, extends, include, docker compose
title: Overview
---
{% include compose-eol.md %}

This section contains information on the three key ways you can use multiple Compose files in your Compose application. 

Using multiple Compose files lets you customize a Compose application
for different environments or different workflows.


This is useful for large applications, using dozens, maybe hundreds of containers, with ownership distributed across multiple teams. For example, if your organization or team uses a monorepo it is common for teams to have their own “local” compose file to just run a subset of the application, but then they need to rely on other teams to provide some reference compose file defining the expected way to run their own subset.

