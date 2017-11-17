---
title:  Node access control in Docker EE Advanced
description: Learn how to architect node access with Docker Enterprise Edition Standard.
keywords: authorize, authentication, node, UCP, role, access control
redirect_from:
- /ucp/
ui_tabs:
- version: ucp-3.0
  orhigher: true
- version: ucp-2.2
  orlower: true
---

{% if include.ui %}
{% if include.version=="ucp-3.0" %}

*Node access control* lets you segment scheduling and visibility by node. Node access control is available in [Docker EE Advanced](https://www.docker.com/pricing) with Swarm orchestration only.

{% elsif include.version=="ucp-2.2" %}

*Node access control* lets you segment scheduling and visibility by node. Node access control is available in [Docker EE Advanced](https://www.docker.com/pricing).

{% endif %}

By default, non-infrastructure nodes (non-UCP & DTR nodes) belong to a built-in
collection called `/Shared`. All application workloads in the cluster are
scheduled on nodes in the `/Shared` collection. This includes those deployed in
private collections (`/Shared/Private/`) or any other collection under
`/Shared`.

This setting is enabled by a built-in grant that assigns every UCP user the
`scheduler` capability against the `/Shared` collection.

Node Access Control works by placing nodes in custom collections (outside of
`/Shared`). If a user or team is granted a role with the `scheduler` capability
against a collection, then they can schedule containers and services on these
nodes.

In the following example, users with `scheduler` capability against
`/collection1` can schedule applications on those nodes.

Again, these collections lie outside of the `/Shared` collection so users
without grants do not have access to these collections unless it is explicitly
granted. These users can only deploy applications on the built-in `/Shared`
collection nodes.

![image](../images/design-access-control-adv-custom-grant.png){: .with-border}

The tree representation of this collection structure looks like this:

```
/
├── Shared
├── System
├── collection1
└── collection2
    ├── sub-collection1
    └── sub-collection2
```

With the use of default collections, users, teams, and organizations can be
constrained to what nodes and physical infrastructure they are capable of
deploying on.


## Where to go next

- [Isolate swarm nodes to a specific team](isolate-nodes-between-teams.md)
{% endif %}
