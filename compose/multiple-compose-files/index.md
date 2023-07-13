---
description: General overview for the different ways you can work with multiple compose files in Docker Compose
keywords: compose, compose file, merge, extends, include, docker compose
title: Overview
---
{% include compose-eol.md %}

This section contains information on the three key ways you can use multiple Compose files in your Compose application. 

Using multiple Compose files lets you customize a Compose application for different environments or different workflows.

This is useful for large applications, using dozens, maybe hundreds of containers, with ownership distributed across multiple teams. For example, if your organization or team uses a monorepo it is common for teams to have their own “local” compose file to just run a subset of the application, but then they need to rely on other teams to provide some reference compose file defining the expected way to run their own subset.

With microservices and monorepo, it becomes common for an application to be split into dozens of services, and complexity is moved from code into infrastructure and configuration file. Docker Compose fits well with simple applications but is harder to use in such a context. At least it was, until now.

There are three different options when working with multiple Compose files depending on your needs. You can:
- [Extend a Compose file](extends.md) by referring to another compose file and selecting a service you want to also use in your own application, with the ability to override some attributes for your own needs.
- [Merge a set of Compose files](merge.md) together to create a composite Compose file.
- [Include other Compose files](include.md) directly in to your Compose file. 
