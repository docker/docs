// Package authz handles authentication and authorization for Docker Trusted
// Registry.
//
// Authentication is separated from Authorization so that a variety
// of authentication methods can be supported:
//
// - Basic authentication with hashing of managed user passwords.
// - Basic authentication via an LDAP BIND.
// - HTTP Cookie verification for client sessions.
// - And many more to come.
//
// Each authentication mechanism returns a common user account object which is
// used throughout Docker Trusted Registry and allows for separation of
// concerns for authorization to various resources. The user objects will store
// the username and which global roles (if any) they have as well as be able to
// lookup what teams they are in for other organization accounts.
//
// Currently, this package only handles authorization for repositories and
// repository namespaces. The same logic is used no matter what authentication
// mechanism was used.
//
// Repository Namespace Authorization:
//
// To determine the level of access that a user has to a given repository
// namespace, the following flow is used:
//
// 	Is the user a global admin?
// 		yes? -> return "admin" level access
// 		no? -> continue
//
// 	If the user is in the global "read-write" or "read-only" role, remember
// 	that (highest) corresponding level of access.
//
// 	Is the namespace the user's own (same name)?
// 		yes? -> return "admin" level access
// 		no? -> continue
//
// 	Is the namespace owned by some other user account?
// 		yes? -> return the highest global access level (if any) or
// 			none.
// 		no? -> continue
//
// 	Then the namespace must be owned by an organization account. The teams
// 	of this organization that the user is a member of are then queried:
//
// 	Is the user a member of the "owners" team in this organization?
// 		yes? -> return "admin" level access
// 		no? -> continue
//
// 	Get the highest level of namespace access that the user's teams are
// 	granted to this namespace. If that team level of access is higher than
// 	the user's global access level (if any), return that, otherwise return
// 	the user's global access level (which may be none).
//
// Repository Authorization:
//
// To determine the level of access that a user has to a given repository, the
// following flow is used:
//
// 	Note: If the repository has "public" visibility, return at least "read"
// 	level access if no access would be returned otherwise.
//
// 	Get the user's access to the repository's namespace using the method
// 	descripbed above.
//
// 	Is the namespace access level "admin"?
// 		yes? -> return "admin" level access
// 		no? -> continue
// 	Note: this handles the case of the user being a global admin, the
// 	repository being owned by the user, the user being an owner of the
// 	organization (if an organization's repository), and other global roles.
//
// 	If the user is in the global "read-write" or "read-only" role, remember
// 	that (highest) corresponding level of access.
//
// 	Is the repository owned by some other user account?
// 		yes? -> return the higher of the namespace access (if any)
// 			and the owner's granted access to this repository for
// 			the user (if any).
// 		no? -> continue
//
// 	Then the repository must be owned by an organization account. The
// 	teams of this organization that the user is a member of are then
// 	queried:
//
// 	Get the highest level of repository access that the user's teams are
// 	granted to this repository. If that team level of access is higher
// 	than the user's namespace access level (if any), return that, otherwise
// 	return the user's namespace access (which may be none).
package authz
