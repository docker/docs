# Migration

In Docker EE UCP 2.X, the layer 7 routing of the platform was implemented using HRM. Starting with Docker EE UCP 3.X, the layer 7 routing is now implemented using Interlock. This article describes the process in which the migration from HRM to Interlock happens.

For context, let me recap on the HRM to interlock migration process:

0. HRM service is already enabled and configured by the customer.
1. Migration from 2.2.4 to 3.0 starts
2. Reconcilation process for interlock starts and HRM service is inspected.
3. An interlock config is created from the HRM service spec and deployed as a docker config
4. The HRM service is removed.
5. The interlock service is deployed.
6. The interlock extension and proxy service are deployed by the interlock service.

Currently, this process is not atomic, so any failure along the way will not result in the reversal of previous steps. In any case, if a failure occurs today the migration would need manual fixing to operate normally again.

In the following section, I cover different failure points, how to diagnostic, and to solve (when possible) from a customer/support standpoint.

### Before starting migration

**TODO**: Information [on interlock manual deploymenbt](https://beta.docs.docker.com/ee/ucp/interlock/install/manual-deployment/#deployment) is incorrect. Please update with procedure in this PR: https://github.com/docker/orca/issues/12227.

#### Migration check procedure

Record which of your HRM-enabled applications are routable through HRM by doing a curl request:

 - For HTTP apps: ```curl -vs http://${CLUSTER_ID}:${HRM_HTTP_PORT}/ -H "Host: ${HOSTNAME}"```
 - For HTTPS apps: ```curl -vs http://${HOSTNAME}:${HRM_HTTPS_PORT}"```

When the migration, check that all applications that where routable before the migration still are.

#### Failure happens during step 1 before step 2

Not related to HRM interlock migration

#### Failure happens during step 2 before step 3

Codepath: https://github.com/docker/orca/blob/master/agent/agent/components/interlockservice/interlockservice.go

 - **Diagnostic**: Only happens if the HRM service gets deleted somehow after being detected as present but before the migration to interlock starts.
 - **Resolution**: Follow the *manual hrm migration* procedure.

#### Failure happens during step 3 before step 4

 - **Diagnostic**: `ucp-hrm` is still present, but the config `com.docker.ucp.interlock.conf-1` could not be created. Could happen due to name conflicts.
 - **Resolution**: Follow the *manual hrm migration* procedure.

#### Failure happens during step 4 before step 5

 - **Diagnostic**: `ucp-hrm` failed to remove. The ucp interlock config is already created. Only happens if the service delete operation for the `ucp-hrm` service fails for some reason (unlikely).
 - **Resolution**: Follow the *manual hrm migration* procedure, but replace the interlock swarm config name from `com.docker.ucp.interlock.conf-1` to `com.docker.ucp.interlock.conf-2`.

#### Failure happens during step 5 before step 6

 - **Diagnostic**: The `ucp-interlock` service failed to deploy. Only happens if the service create operation for the `ucp-interlock` service fails for some reason (unlikely).
 - **Resolution**: Enable interlock through the UI, making sure that the `PublishedPort` and `PublishedSSLPort` match the http and https ports you were using with HRM.

#### Failure happens during step 6

 - **Diagnostic**: The `ucp-interlock-extension` or `ucp-interlock-proxy` services failed to deploy. May happen if there is a port conflict.
 - **Resolution**: Disable and then re-enable interlock through the UI, making sure that the `PublishedPort` and `PublishedSSLPort` match the http and https ports you were using with HRM.

#### Manual Migration procedure

 1. delete the `ucp-hrm` service.
 2. Enable interlock manually (https://beta.docs.docker.com/ee/ucp/interlock/install/manual-deployment/), but instead of creating a new overlay network called `ucp-interlock`, reuse the existing `ucp-hrm`.
 3. Wait until the service named `ucp-interlock-proxy` is deployed with all replicas and then try the *migration check procedure* again.
