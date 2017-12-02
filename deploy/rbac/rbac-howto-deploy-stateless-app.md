---
title: Deploy a simple stateless app with RBAC
description: Learn how to deploy a simple application and customize access to resources.
keywords: rbac, authorize, authentication, users, teams, UCP, Docker
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

This tutorial explains how to deploy a nginx web server and limit access to one
team with role-based access control (RBAC).

## Scenario

You are the Docker EE admin at Acme Company and need to configure permissions to
company resources. The best way to do this is to:

- Build the organization with teams and users
- Create collections or namespaces for storing resources.
- Create grants that specify which team can do what operations on which
  collection or namespace.
- Give the `ops` team the all-clear to deploy nginx.

## Build the organization

Add the organization, `acme-datacenter`, and create three teams according to the
following structure:

```
acme-datacenter
├── dba
│   └── Alex Alutin
├── dev
│   └── Bett Bhatia
└── ops
    └── Chad Chavez
```

> Easy username / passwords:
> - alex / alexalutin
> - bett / bettbhatia
> - chad / chadchavez

See: [Create and configure users and teams](./usermgmt-create-subjects.md).

## Kubernetes deployment

In this section, we deploy `nginx` with Kubernetes. See [Swarm stack](#swarm-stack)
for the same exercise with Swarm.

### Create namespace

Create a namespace to logically store the nginx application:

1. Click **Kubernetes** > **Namespaces**.
2. Paste the following manifest in the terminal window:

```
apiVersion: v1
kind: Namespace
metadata:
  name: nginx-namespace
```

2. Click **Create**.


### Grant roles

Grant the ops team (and only the ops team) access to nginx-namespace with the
built-in role, **Full Control**.

```
acme-datacenter/ops + Full Control + nginx-namespace
```

### Deploy Nginx

You've configured Docker EE. The `ops` team can now deploy `nginx`.

1. Log on to UCP as chad (on the `ops`team).
2. Click **Kubernetes** > **Namespaces**.
3. Paste the following manifest in the terminal window and click **Create**.

```
apiVersion: apps/v1beta2 # for versions before 1.8.0 use apps/v1beta1
kind: Deployment
metadata:
  name: nginx-deployment
spec:
  selector:
    matchLabels:
      app: nginx
  replicas: 2
  template:
    metadata:
      labels:
        app: nginx
    spec:
      containers:
      - name: nginx
        image: nginx:1.7.9
        ports:
        - containerPort: 80
```

4. Log on to UCP as each user and ensure that:
- `dba` (alex) cannot see `nginx-namespace`.
- `dev` (bett) cannot see `nginx-namespace`.


## Swarm Stack

In this section, we deploy `nginx` as a Swarm service. See [Kubernetes Deployment](#kubernetes-deployment)
for the same exercise with Swarm.

### Create collection paths

Create a collection for nginx resources, nested under the `/Shared` collection:

```
/
├── System
└── Shared
    └── nginx-collection
```

> **Tip**: To drill into a collection, click **View Children**.

See: [Group and isolate cluster resources](./resources-group-resources.md).

### Grant roles

Grant the ops team (and only the ops team) access to nginx-collection with the
built-in role, **Full Control**.

```
acme-datacenter/ops + Full Control + /Shared/nginx-collection
```

See: [Grant access to cluster resources](./usermgmt-grant-permissions.md).

### Deploy Wordpress and MySQL with Swarm

You've configured Docker EE. The `ops` team can now deploy an `nginx` Swarm
service.

1. Log on to UCP as chad (on the `ops`team).
2. Click **Swarm** > **Services**.
3. Click **Create Stack**.
4. On the Details tab, enter:
   - Name: `nginx-service`
   - Image: nginx:latest
4. On the Collections tab:
   - Click `/Shared` in the breadcrumbs.
   - Select `nginx-collection`.
5. Click **Create**.
6. Log on to UCP as each user and ensure that:
   - `dba` (alex) cannot see `nginx-collection`.
   - `dev` (bett) cannot see `nginx-collection`.


{% elsif include.version=="ucp-2.2" %}

## Swarm Stack

In this section, we deploy `nginx` as a Swarm service. See [Kubernetes Deployment](#kubernetes-deployment)
for the same exercise with Swarm.

### Create collection paths

Create a collection for nginx resources, nested under the `/Shared` collection:

```
/
├── System
└── Shared
    └── nginx-collection
```

> **Tip**: To drill into a collection, click **View Children**.

See: [Group and isolate cluster resources](./resources-group-resources.md).

### Grant roles

Grant the ops team (and only the ops team) access to nginx-collection with the
built-in role, **Full Control**.

```
acme-datacenter/ops + Full Control + /Shared/nginx-collection
```

See: [Grant access to cluster resources](./usermgmt-grant-permissions.md).

### Deploy Wordpress and MySQL with Swarm

You've configured Docker EE. The `ops` team can now deploy an `nginx` Swarm
service.

1. Log on to UCP as chad (on the `ops`team).
2. Click **Swarm** > **Services**.
3. Click **Create Stack**.
4. On the Details tab, enter:
   - Name: `nginx-service`
   - Image: nginx:latest
4. On the Collections tab:
   - Click `/Shared` in the breadcrumbs.
   - Select `nginx-collection`.
5. Click **Create**.
6. Log on to UCP as each user and ensure that:
   - `dba` (alex) cannot see `nginx-collection`.
   - `dev` (bett) cannot see `nginx-collection`.

{% endif %}
{% endif %}
