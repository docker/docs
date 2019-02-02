---
description: Page for deprecation announcements.
keywords: dtr, manifest list, api, repository, digest
title: Deprecation Notice
---

This document outlines the functionalities or components within DTR that will be deprecated.

### Enable Manifest List via the API

#### 2.5 and 2.6

Since `v2.5`, it has been possible for repository admins to autogenerate manifest lists when [creating a repository via the API](/datacenter/dtr/2.5/reference/api/). You accomplish this by setting `enableManifestLists` to `true` when sending a POST request to the `/api/v0/repositories/{namespace}` endpoint. When enabled for a repository, any image that you push to an existing tag will be appended to the list of manifests for that tag. `enableManifestLists` is set to false by default, which means pushing a new image to an existing tag will overwrite the manifest entry for that tag.

#### 2.7

The above behavior and the `enableManifestLists` field will be removed in `v2.7`. Starting in `v2.7`, you can use the DTR CLI to create and push a manifest list to any repository. 

