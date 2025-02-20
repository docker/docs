---
title: Docker Projects FAQs
linkTitle: FAQs
description: Find common FAQs for Docker Projects
keywords: faqs, docker projects, local, remote
---

## Why is a Compose file required?

A Compose file (`compose.yml`) defines how your application's containers should run together, including:

 - Services (e.g., web, database, API)
 - Networks for inter-container communication
 - Volumes for persistent data storage
 - Environment variables and configurations

Without a Compose file, Docker Projects doesn't have a way to understand how your application should be structured or executed.

## What if my project doesn’t have a Compose file?

If your project doesn't include a `compose.yml` file, you need to create one before opening it in Docker Projects. 