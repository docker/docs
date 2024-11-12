---
description: Learn how to manage repositories on Docker Hub
keywords: Docker Hub, Hub, repositories
title: Repositories
weight: 20
aliases:
- /engine/tutorials/dockerrepos/
- /docker-hub/repos/
- /docker-hub/repos/configure/
---

A Docker Hub repository is a collection of container images, enabling you to
store, manage, and share Docker images publicly or privately. Each repository
serves as a dedicated space where you can store images associated with a
particular application, microservice, or project. Content in repositories is
organized by tags, which represent different versions of the same application,
allowing users to pull the right version when needed.

## Key features and concepts

- [Repository information](./manage/information.md): Each repository can include a
  description, an overview, and categories to help users understand its purpose
  and usage. Adding clear repository information ensures that others can find
  your images and know how to use them.

- [Access management](./manage/access.md): Docker Hub repositories offer flexible
  access options. You can make repositories public or private and add
  collaborators. For organizations, teams, roles, and access tokens enable
  fine-grained access, providing security and control at scale.

- [Image management](./manage/hub-images/_index.md): Repositories support
  various types of content, including OCI artifacts, and provide robust version
  control with tagging. You can push new images and move existing content
  between repositories, ensuring flexibility in managing different versions and
  types of container images.

- [Webhooks](./manage/webhooks.md): Webhooks let you automate responses to repository
  events. For instance, you can set up notifications or trigger actions in
  external systems whenever an image is pushed or updated, helping to streamline
  workflows.

- [Automated builds](./manage/builds/_index.md): Docker Hub
  repositories integrate with GitHub or Bitbucket for automated builds, ensuring
  that every code change triggers an image rebuild. This feature supports
  continuous integration and delivery, keeping images up-to-date and
  streamlining development pipelines.

- [Image security insights](./manage/vulnerability-scanning.md): Docker Hub repositories
  provide powerful features and controls to help uncover, understand, and fix
  issues with container images. Options for image security insights include
  continuous Docker Scout image analysis and static vulnerability scanning.