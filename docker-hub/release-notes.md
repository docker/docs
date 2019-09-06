---
title: Docker Hub release notes
description: Learn about the new features, bug fixes, and breaking changes for Docker Hub
keywords: docker hub, whats new, release notes
toc_max: 2
---

Here you can learn about the latest changes, new features, bug fixes, and known issues for each Docker Hub release.

## 2019-09-05

### Enhancements

* The `Tags` tab on an image page now provides additional information for each tag:
  * A list of digests associated with the tag
  * The architecture it was built on
  * The OS
  * Information about who edited the tag last
* The vulnerability scanning viewer for [official images](https://docs.docker.com/docker-hub/official_images/) has been updated. You can reach it by clicking a tag's digest link. If each layer has either a green checkmark or a red X next to it, that image has been scanned. Click on any layer to see more information.

### Known Issues

* Scan results don't appear for some official images.
