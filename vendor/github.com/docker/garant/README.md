Garant
======

A JSON Web Token Service for Everyone.


### Introduction

Along with the development of the V2 Docker Registry, we desired a way to
separate authentication and authorization from the registry core and so a
pluggable authorization interface was developed for the v2 registry to provide
a list of resources and actions to be performed by a request and simply return
whether or not those action on those resources are authorized. With this
interface the new registry could cleanly separate the concerns of access
control and authn/authz mechanism from the rest of the application logic.
Admins can then build and configure a variety of auth backends. One such
backend which shipped with the new registry was a JSON Web Token backend. This
allows the authorized actions and respective resources to be be described in a
JSON object known as a claim set along with the authenticated user (if
applicable), expiration time, and unique ID per token. These tokens are
verifiable with public key cryptography so that the registry server does not
require a shared secret to verify and authorize requests. While this token
verifying backend is built into the registry, Garant is a complementary JWT
generating service.

### How it Works

The Garant Token Signing App uses configurable backends which allow for
authenticating clients in a variety of ways such as:

- Authenticating a client using HTTP Basic Auth, and verifying the password
  against a database or LDAP server.
- Secret Tokens associated with an account or authorized 3rd party as
  controlled by a centralized identity management system.
- TLS Client Certificates also associated with an account allowing for the use
  of no shared secrets throughout your authentication system.

Each backend must also be able to specify granted authorization to perform
actions on given resources. The interface is generic enough that any service
which defines typed and/or named resources with a known set of actions may
implement a backend that filters requested resource actions to a list of
granted resource actions for an authenticated client.

The granted resource actions and authenticated account information is then
placed in a JWT claim set and signed using the token app's private signing key.
Any appropriately configured service may use the same verification method to
ensure that resource access has been authorized by their Garant token singing
app. As such, this usage pattern is not limited to the v2 Docker Registry.

### Use by Docker

The official Docker Hub registry delegates all authentication and authorization
to a specially developed backend for Garant. This backend authenticates clients
using their Docker Hub username and password, and authorizes `push` and `pull`
actions on named repositories based on repository visibility, the authenticated
account, organization and group membership, and user or group collaborators as
specified on a per-repository basis.

The current design makes it easy for new authentication mechanisms to be added
to Docker Hub such as login tokens or deploy keys (possibly used by 3rd
parties) each with varying scope granularity for access control.
