---
title: "Faster development and testing with container-supported development"
linkTitle: Container-supported development
summary: |
  Containers don't have to be just for your app. Learn how to run your app's dependent services and other debugging tools to enhance your development environment.
description: |
  Use containers in your local development loop to develop and test faster… even if your main app isn't running in containers.
levels: [beginner]
params:
  image: images/learning-paths/container-supported-development.png
  time: TBD
  resource_links: []
---

Containers provide a standardized ability to build, share, and run applications. While containers are typically used to containerize your application, they also make it incredibly easy to run essential services needed for development. Instead of installing or connecting to a remote database, you can easily launch your own database. But the possibilities don't stop there.

With container-supported development, you use containers to enhance your development environment by emulating or running your own instances of the services your app needs. This provides faster feedback loops, less coupling with remote services, and a greater ability to test error states.

And best of all, you can have these benefits regardless of whether the main app under development is running in containers.

## What you'll learn

- The meaning of container-supported development
- How to connect non-containerized applications to containerized services
- Several examples of using containers to emulate or run local instances of services
- How to use containers to add additional troubleshooting and debugging tools to your development environment

## Who's this for?

- Teams that want to reduce the coupling they have on shared or deployed infrastructure or remote API endpoints
- Teams that want to reduce the complexity and costs associated with using cloud services directly during development
- Developers that want to make it easier to visualize what's going on in their databases, queues, etc.
- Teams that want to reduce the complexity of setting up their development environment without impacting the development of the app itself


## Tools integration

Works well with Docker Compose and Testcontainers


<div id="lp-survey-anchor"></div>