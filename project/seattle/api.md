# UCP-Seattle API and Applications Spec

# API Routes Compatibility and Integration
* The existing Docker remote REST API paths map to swarm v1 manager.  This
  should require no change to the wiring that is already present

* The New APIs map to the engine (swarm v2 manager) endpoint on the local
  controller node. These APIs are expected to correspond to the `cluster`, CLI
verb.

* The majority of the new APIs are expected to be hijacked by UCP for Access
  Control and Node Lifecycle management. 

* Since the Swarm v2 API is not fully specced out yet, the specific API
  endpoints will be fleshed out as they become stable. Our source of truth will
be https://github.com/docker/engine-api-private/tree/1.12-integration

# Multi-resource Access Control requests

* The declarative API of engine 1.12 will allow users to access multiple
resources in the context of a single API request, motivating the extension of
our RBAC system.  As a first step, the existing logic for containers can be
extended to Networks and Volumes. This resolves the highly-requested feature
where non-admin users can bind to the Root CA volumes. Moreover, the implicit
creation of Volumes during container creation makes for a great first iteration
on how declarative requests of multiple resources can be handled.

* The proposed flow for this feature is the following:
	* A complex request is disected into individual resource objects by a
	  `parser` method in the pipeline. The parser is route-specific and it is
allowed to "peek" inside the body of the request. The parser identifies whether
each of the requested resources is a new resource or an existing resource in
the system and generates an appropriate ResourceRequest object.
	* The Access Control layer ensures that the current user has access to each
	  of the ResourceRequest objects.
	* The actual handlers responsible for resource creation (e.g.
	  applicationCreate) would ensure that all newly-created resources are
tagged with the appropriate owner label. If there is no native API support for
this, the top-level handler would need to invoke the handlers of other
resources with the appropriate ResourceRequest, such as volumeCreate or
containerCreate.
	* The parser of the "edit" operation would need to consider the diff of the
	  declarative API request before generating the corresponding
ResourceRequest objects. This would allow low-privilege users to edit parts of
multi-tier apps without full control.


# Service Accounts / Privileged System Access

* In the current iteration of the designs, the proxy beachhead will need an
  authorized way of calling back to UCP (see the ./proxy.md spec). This
  motivates the creation of a service accounts feature in the long term. 
  However, given the scope of this release, it might be possible to get away
  with short-lived secrets.

# Outside the scope of this release

* RBAC for Images requires us to reconsider our entire Access Control model, as
  the existing model on containers is not sufficient for images. Images are
currently expected to be globally available within a UCP cluster, and changing
this model would break existing workflows. For this reason, implementing RBAC
for images requires a deeper iteration over our Access Control model,
potentially also requiring registry integration.

## Misc

* Expose config for cluster itself
	* https://github.com/docker/docker-1.12-integration/issues/1
