+++
draft = "true"
+++

# UCP Access Control
Access control in UCP is handled by request middleware.  This doc goes into
technical detail how it works.

UCP provides the Docker Remote API transparently to the user.  When a user
performs a Docker action (create container, start container, pull image, etc)
either by the user interface, the Docker client or another Docker API compatible
tool, it is handled by UCP and then proxied to Swarm.  UCP has middleware that
performs various operations (auth, access control, audit, etc).  We try to
limit the amount of extra state/metadata that we have to store in our datastore
for these but some is required.  For access control, we store the access
control lists in the datastore.

Access control lists are comprised of a team, label and access role.  An
account is added to this team and that account has the level of access defined
by the role based upon the resource label.  Here is an example:

ACL: Team dev has view access (read only) to label "dev"

- Team: `dev`
- Label: `dev`
- Role: `View`

This team has "view" access meaning any user part of this team can only view
the resources.  An example would be a container that has the `dev` UCP access
label means that any member of the dev team could inspect the container but
would not be able to restart, kill, remove, etc.

If a user is a member of multiple teams with varying access levels, the
highest level is used.  For example, if a user is a member of a team has a 
label with `View` access and another team that has `Full Control` to the
same label, the user will have `Full Control` privilges.

# Roles
The following are the roles supported by UCP.

## None
The `None` role does not provide any access.  A member in a team with a label
with this role will have no access to any resource.

## View Only
The `View` role provides read only access.  This means that a member in a team
with a label using the `View` role will not be able to edit the resource.
Examples of this include being able to view a container but not restart,
kill, remove, etc.  

## Restricted Control
The `Restricted Control` role provides restricted edit functionality.
For containers, this means a member in a team with a label using this role
an create, restart, kill, remove, etc.  However, `Restricted Control`
prevents `exec` access.

The `Restricted Control` role prevents the user from specifying options
when creating a container that are deemed dangerous to the system.  This
is not an exhaustive list and further policy enforcement will come in the
future.  However, this is an initial first step to prevent unwanted
operations on the cluster.  A `Restricted Control` role cannot specify the
following operations when creating a container; doing so will result in an
error:

- `privileged`
- `cap-add`
- host mounted volumes
- ipc mode
- pid mode

## Full Control
The `Full Control` role provides full access (create, kill, restart, exec) to
containers, however, cannot edit system configuration (settings, users, etc).

# User Default Permissions
Each user has a default permission.  This is used to decide what level of access
to provide to a user when no label is specified.  For admins, this is
ignored and the admin will always have access.  For non-admins, the same
privileges are used only applied at the user level.  This also means that
users can have "private" containers.  If a user has a default role of
`Restricted Control` or higher, the user can launch containers with no labels.
Only that user will be able to view and edit the container.

The default user permission is also used for non-container resources.  Since
we do not yet have label support for images, networks or volumes, the
default user role is used for accessing these resources.  Once we have
label support for these the behaviour will be the same as containers.


# Technical Design
The following are technical notes on the implementation.

## Middleware Details
This is a more technical explanation of how the access control middleware works.

A Route is a definition of an API Path and the set of functions that need to be 
executed for that API Path. These functions include middleware, such as access
control and tracking, and the actual handler. A Middleware Pipeline object is used
to simplify route definitions in `controller/api/routes.go` by having all routes use
the following middleware ordering:

`Authentication -> Route-specific Parser -> Access Control -> Tracking -> Auditing -> Route-specific Handler`

The Authentication middleware generates an OrcaRequestObject (defined in 
`controller/ctx/resource_context.go`) which is used as the sole argument to the 
followup middleware functions. 

The Route-specific Parsers are responsible for generating a set of Resource objects
for a user request, which are stored in the OrcaRequestObject. Each Resource object 
implements a `HasAccess(auth.Context) bool` method which is invoked by the Access Control
middleware to determine whether the current user has access to all requested resources. 
For example, a container inspect operation at `/containers/{id}/json` uses a containerViewParser,
which generates a ContainerResource for the container with the provided ID and appends that resource
to the OrcaRequestContext. The access control layer invokes HasAccess for that resource to determine
whether the current user has view-only access on the specified container.

The vision behind this design is based on the logic isolation of each individual stage.
- A first look on the table of Route definitions should be able to clearly convey what the
  access control restrictions are for a route, by observing the path and parser name.
- Parsers should break up a request to as many independent Resource objects as possible and
  handle any special route-specific logic themselves in the process. One exception to this is `create` operations, which might require special handling within the resource object itself.
- Resources should be as agnostic as possible to what the actual operation is. The goal is to move as much logic as possible away from the `HasAccess` methods.
- If there's access control logic that cannot be placed in a new resource type, it should be added to the parser rather than bloating up existing resources.
