---
title: Access control design with Docker EE Standard
description: Learn how to architect multitenancy by using Docker Enterprise Edition Advanced.
keywords: authorize, authentication, users, teams, groups, sync, UCP, role, access control
---

[Collections and grants](index.md) are strong tools that can be used to control
access and visibility to resources in UCP. This tutorial describes a fictitious
company named OrcaBank that is designing the access architecture for two
application teams that they have, Payments and Mobile.

This tutorial introduces many concepts include collections, grants, centralized
LDAP/AD, and also the ability for resources to be shared between different teams
and across collections.

## Team access requirements

OrcaBank has organized their application teams to specialize more and provide
shared services to other applications. A `db` team was created just to manage
the databases that other applications will utilize. Additionally, OrcaBank
recently read a book about DevOps. They have decided that developers should be
able to deploy and manage the lifecycle of their own applications.

- `security` should have visibility-only access across all applications in the
  swarm.
- `db` should have the full set of capabilities against all database
  applications and their respective resources.
- `payments` should have the full set of capabilities to deploy Payments apps
  and also access some of the shared services provided by the `db` team.
- `mobile` has the same rights as the `payments` team, with respect to the
  Mobile applications.

## Role composition

OrcaBank will use a combination of default and custom roles, roles which they
have created specifically for their use case. They are using the default
`View Only` role to provide security access to only see but not edit resources.
There is an `ops` role that they created which can do almost all operations
against all types of resources.  They also created the
`View & Use Networks + Secrets` role. This type of role will enable application
DevOps teams to use shared resources provided by other teams. It will enable
applications to connect to networks and use secrets that will also be used by
`db` containers, but not the ability to see or impact the `db` applications
themselves.

![image](../images/design-access-control-adv-0.png){: .with-border}

## Collection architecture

OrcaBank will also create some collections that fit the organizational structure
of the company. Since all applications will share the same physical resources,
all nodes and applications are built in to collections underneath the `/Shared`
built-in collection.

- `/Shared/payments` hosts all applications and resources for the Payments
  applications.
- `/Shared/mobile` hosts all applications and resources for the Mobile
  applications.

Some other collections will be created to enable the shared `db` applications.

- `/Shared/db` will be a top-level collection for all `db` resources.
- `/Shared/db/payments` will be specifically for `db` resources providing
  service to the Payments applications.
- `/Shared/db/mobile` will do the same for the Mobile applications.

The following grant composition will show that this collection architecture
allows an app team to access shared `db` resources without providing access
to _all_ `db` resources. At the same time _all_ `db` resources will be managed
by a single `db` team.

## LDAP/AD integration

OrcaBank has standardized on LDAP for centralized authentication to help their
identity team scale across all the platforms they manage. As a result LDAP
groups will be mapped directly to UCP teams using UCP's native LDAP/AD
integration. As a result users can be added to or removed from UCP teams via
LDAP which can be managed centrally by OrcaBank's identity team. The following
grant composition shows how LDAP groups are mapped to UCP teams .

## Grant composition

Two grants are applied for each application team, allowing each team to fully
manage their own apps in their collection, but also have limited access against
networks and secrets within the `db` collection. This kind of grant composition
provides flexibility to have different roles against different groups of
resources.

![image](../images/design-access-control-adv-1.png){: .with-border}

## OrcaBank access architecture

The resulting access architecture shows applications connecting across
collection boundaries. Multiple grants per team allow Mobile applications and
Databases to connect to the same networks and use the same secrets so they can
securely connect with each other but through a secure and controlled interface.
These resources are still deployed across the same group of UCP
worker nodes. Node segmentation is discussed in the [next tutorial](#).

![image](../images/design-access-control-adv-2.png){: .with-border}

### DB team

The `db` team is responsible for deploying and managing the full lifecycle
of the databases used by the application teams. They have the full set of
operations against all database resources.

![image](../images/design-access-control-adv-3.png){: .with-border}

### Mobile team

The `mobile` team is responsible for deploying their own application stack,
minus the database tier which is managed by the `db` team.

![image](../images/design-access-control-adv-4.png){: .with-border}

## Where to go next

- [Access control design with Docker EE Advanced](access-control-design-ee-advanced.md)


