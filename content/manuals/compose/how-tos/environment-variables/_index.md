---
title: Environment variables in Compose
linkTitle: Use environment variables
weight: 40
description: Explainer on the ways to set, use and manage environment variables in
  Compose
keywords: compose, orchestration, environment, env file
aliases:
- /compose/environment-variables/
---

By leveraging environment variables and interpolation in Docker Compose, you can create versatile and reusable configurations, making your Dockerized applications easier to manage and deploy across different environments. 

> [!TIP]
>
> Before using environment variables, read through all of the information first to get a full picture of environment variables in Docker Compose.

This section covers:

- [How to set environment variables within your container's environment](set-environment-variables.md).
- [How environment variable precedence works within your container's environment](envvars-precedence.md).
- [Pre-defined environment variables](envvars.md).

It also covers: 
- How [interpolation](variable-interpolation.md) can be used to set variables within your Compose file and how it relates to a container's environment.
- Some [best practices](best-practices.md).
