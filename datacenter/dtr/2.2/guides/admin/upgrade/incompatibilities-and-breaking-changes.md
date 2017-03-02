---
title: Incompatibilities and breaking changes
description: Learn about the incompatibilities and breaking changes introduced by Docker Trusted Registry version {{ page.dtr_version }}
keywords: docker, ucp, upgrade, incompatibilities
---

With Docker Trusted Registry {{ page.dtr_version }}, the `/load_balancer_status`
endpoint is deprecated and is going to be removed in future versions. Use the
`/health` endpoint instead.
