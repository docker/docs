---
description: Page for deprecation announcements.
keywords: registry, manifest, images, signatures, repository, distribution, digest
title: Deprecation Notice
---

This document outlines the functionalities or components within DTR that will be deprecated.

### Enable Manifest List via the API

Since `v2.5`, it has been possible for repository admins to enable manifest lists when [creating a repository via the API](./reference/dtr/2.5/api/). You accomplish this by setting `enableManifestLists` to `true` when sending a POST request to the `/api/v0/repositories/{namespace}` endpoint. When enabled for a repository, any image that you push to an existing tag will be added to the list of manifests for that tag.

The above behavior and the field `enableManifestLists` will be removed in `v2.7`. Starting in `v2.7`, you can use the CLI command, `docker manifest` to [create and push a manifest list to any repository](./edge/engine/reference/commandline/manifest/). 

