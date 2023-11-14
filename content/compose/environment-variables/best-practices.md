---
title: Best practices for working with environment variables in Docker Compose
description: Explainer on the ways to set, use and manage environment variables in
  Compose
keywords: compose, orchestration, environment, env file
---

1. Handle Sensitive Information Securely:
   Be cautious about including sensitive data in environment variables. Consider using [Secrets](../use-secrets.md) for managing sensitive information.

2. Understand Environment variable precedence:
    Be aware of how Docker Compose handles the precedence of environment variables from different sources (`.env` files, shell variables, Dockerfiles).

3. Use specific environment files:
   Consider how your application adapts to different environments. for example development, testing, production, and use different `.env` files as needed.

4. Know interpolation:
   Understand how interplation works within Docker Compose files for dynamic configurations.

5. Command line overrides:
    Be aware that you can override environment variables from the command line when starting containers, useful for testing or temporary changes.

