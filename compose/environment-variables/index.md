---
title: Overview
description: Explainer on the ways to set, use and manage environment variables in Compose
keywords: compose, orchestration, environment, env file
---

This section covers:
- The various ways you can [use environment variables in Compose](set-environment-variables.md).
- How to set environment variables with the correct syntax depending on your chosen method.
- Explains and illustrates how environment variable precedence works. 

## Why use environmental variables?

If you need to define various configuration values, environment variables are your best friend. Like many tools, Docker and, more specifically, Docker Compose can both define and read environment variables to help keep your configurations clean and modular.

Environment variables keep your app secure, flexible and organized. The main advantages:

keep confidential information confidential: prevent baking your passwords and secrets into your image
Organize confidential information in one file; we can .gitignore all of it
Reuse env files by injecting them in multiple containers; all configurations in one place


the first is for ease of cross-environment deployment, and the second is to avoid leaking sensitive information, notably authentication credentials.