---
title: Best practices for working with environment variables in Docker Compose
linkTitle: Best practices
description: Explainer on the best ways to set, use, and manage environment variables in
  Compose
keywords: compose, orchestration, environment, env file, environment variables
tags: [Best practices]
weight: 50
aliases:
- /compose/environment-variables/best-practices/
---

#### Handle sensitive information securely

Be cautious about including sensitive data in environment variables. Consider using [Secrets](../use-secrets.md) for managing sensitive information.

#### Understand environment variable precedence

Be aware of how Docker Compose handles the [precedence of environment variables](envvars-precedence.md) from different sources (`.env` files, shell variables, Dockerfiles).

#### Use specific environment files

Consider how your application adapts to different environments. For example development, testing, production, and use different `.env` files as needed.

#### Know interpolation
   
Understand how [interpolation](variable-interpolation.md) works within compose files for dynamic configurations.

#### Command line overrides
    
Be aware that you can [override environment variables](set-environment-variables.md#cli) from the command line when starting containers. This is useful for testing or when you have temporary changes.

