# Certificate Authority Spec

In the OSS 1.12 release, Swarm V2 has a built-in CA, derived from CFSSL,
but wrapped in a customized gRPC service.  This will have a plugin model
for authorization.  We will swap out our ~pure CFSSL based CA service
with an instance of this CA service, running a UCP plugin.

* The CA tightly controls the CNs - they're randomly generated IDs
* The blacklist tracked at the manager level is based on CNs
    * Removal flows should blacklist the CNs for the "thing" (user, service, node, etc.)
    * When an admin user's admin flag is set to false, purge all their certs, and add their cert CN to the cluster blacklist
    * The controller middleware should leverage the same blacklist to block connections from revoked CNs
* The CA will now be gRPC based instead of CFSSL REST API based, so UCP will have to switch to using a gRPC client lib to talk to it
    * For bootstrapping scenarios we can continue to use vendored CFSSL library code and use the "raw" cert material and configuration
* Mutual-TLS
    * In the prior release, we tightly controlled mTLS to only allow controllers to connect.  All other clients were blocked at the TLS handshake level.
    * Now anyone should be able to connect and make a request, but if the client cert of the requester isn't trusted, then the CSR will be queued until an admin accepts it  (This is already implemented upstream)
* Plugin:
    * If the cert request is a renewal, allow (this should be existing behavior)
    * If it's a new cert request, notify UCP (see below for more node details)
    * After Admin accepts the request (or automation within UCP does so - future most likely) then unblock the CSR



## Misc notes:

* Wire up https://github.com/docker/orca/issues/1513
* gRPC can coexist with REST, so if there's a compelling reason, we could consider supporting both in our CA implementation
    * This might be interesting for upgrade, but breaking secondary controllers during the upgrade seems fine, and the added complexity of supporting both just for upgrade seems unwarranted (since this is a "private" service within UCP)
