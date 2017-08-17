---
redirect_from:
- /engine/misc/breaking/
description: Breaking changes
keywords: docker, documentation, about, technology, breaking, incompatibilities
title: Breaking changes and incompatibilities
---

Every Engine release strives to be backward compatible with its predecessors,
and interface stability is always a priority at Docker.

In all cases, feature removal is communicated three releases
in advance and documented as part of the [deprecated features](deprecated.md)
page.
  
The following list compiles any updates to Docker Engine that created
backwards-incompatibility for old versions of Docker tools.

> **Note**: In the case of your local environment, you should be updating your
  Docker Engine using [Docker for Mac](/docker-for-mac),
  [Docker for Windows](/docker-for-windows). That way all your tools stay
  in sync with Docker Engine.

## Engine 1.10

There were two breaking changes in the 1.10 release that affected
Registry and Docker Content Trust:

**Registry**

Registry 2.3 includes improvements to the image manifest that caused a
breaking change. Images pushed by Engine 1.10 to a Registry 2.3 cannot be
pulled by digest by older Engine versions. A `docker pull` that encounters this
situation returns the following error:

```none
 Error response from daemon: unsupported schema version 2 for tag TAGNAME
```

Docker Content Trust heavily relies on pull by digest. As a result, images
pushed from the Engine 1.10 CLI to a 2.3 Registry cannot be pulled by older
Engine CLIs (< 1.10) with Docker Content Trust enabled.

If you are using an older Registry version (< 2.3), this problem does not occur
with any version of the Engine CLI; push, pull, with and without content trust
work as you would expect.

**Docker Content Trust**

Engine older than the current 1.10 cannot pull images from repositories that
have enabled key delegation. Key delegation is a feature which requires a
manual action to enable.
