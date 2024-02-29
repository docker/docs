---
title: Overview
description: Explainer on the ways to set, use and manage environment variables in
  Compose
keywords: compose, orchestration, environment, env file
---

{{< include "compose-eol.md" >}}

Use environment variables to pass configuration information to containers at runtime. 

Environment variables are key-value pairs that contain data that can be used by processes running inside a Docker container. They are often used to configure application settings and other parameters that may vary between different environments, such as development, testing, and production.  

> **Tip**
>
> Before using environment variables, read through all of the information first to get a full picture of environment variables in Docker Compose.
{ .tip }

This section covers key content such as:
- How to set, use, and manage variables in a Compose file
- [Set environment variables within your container's environment](set-environment-variables.md).
- Changing pre-defined [environment variables](envvars.md).

It also covers: 
- [How environment variable precedence works within your container's environment](envvars-precedence.md).
- The correct syntax for an [environment file](env-file.md).
- Some [best practices](best-practices.md).


