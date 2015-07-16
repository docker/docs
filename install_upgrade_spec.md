# Installation and Upgrade of Orca

![Components](orca_components.png)


## Open Questions

* Should we link all our containers, or wire them up based on the punched through IP/ports?
    * If linking works across hosts, can we rely on that for all communication between orca:swarm:proxy?  If so, that might eliminate the need for a CA in v1, but how do you "secure" the cross-host communication? Feels like it might be a fallacy/chicken-and-egg problem...
* DB Clustering/HA?
* We might want to fold install/upgrade into one script since there's a lot of overlap
* What KV store (swarm discovery backend) should we use?  Can we let the user tweak this?
* How far away is core orca from supporting multiple swarms?
* Should the install "script" actually be mostly implemented as a golang binary, perhaps with a thin shell script wrapper that downloads the right arch binary?
    * Could help us leverage common code between the server and these little external "scripts"
* Should we allow the swarm manager to run on non-standard ports on the engines?
* Possibly include another curl|bash style script for "env" like "docker-machine env <machine>"
    * eval "$(curl -u myorcauser https://myorca/env | bash)" – generate user cert (if not already generated) download it, store it in ~/.docker/orca and echo the eval goop to set up the environment to use it
    * Same as above with https://myorca/swarm/env for admins to get certs to talk directly to swarm (required for upgrade flows below)
* Does it make sense to append the Orca root CA certificate to the local system's trusted certs?
    * Docker CLI needs more, browsers can handle one-off acceptance, so maybe this is just a waste of energy...
* How should "redeploy a broken orca/swarm" work?  Should they redeploy the single node that has orca, then "upgrade" the cluster from there?  If the proxies are busted, they'll likely have to re-add the nodes


## Assumptions

* We wont use data volume containers, but instead host volume mounts
* Most customers do not have swarm (yet), so our primary focus should be on making the Orca+Swarm deployment as clean and simple as possible
* Swarm requires a common single CA “on both sides” (incoming client communication and outgoing engine communication)
* Swarm Managers must have visibility to all the engines (or proxies) and be secured with TLS.  All Engines/Proxies must trust the root CA who signed the swarm cert
* Swarm manager and docker proxy may fold into one component, but this shouldn't fundamentally change the flow
* We'll "own" our own root CAs (One for orca, and one for swarm)
    * Must be "easy" to replace at least Orca's public visible cert/ca (replacing all the swarm ca/certs is desirable too)
        * system must keep functioning, but some capabilities (easy add of hosts, generating user certs) might go away when we no longer own the CA
    * We'll store the certs in a host volume mount
    * The volume could be swapped out for a keywhiz volume mount in the future (unclear if we can write to it though...)
    * Laying the groundwork of a central CA for our managed swarm will enable keywhiz for secret management post v1
* Installation script should be idempotent, and not clobber any pertinent state unless the user asks us to


## Deploy Orca With Swarm

Description:  Deploy orca+swarm onto a single existing engine.  Once deployed, additional engines can be added to the swarm.

```bash
curl https://get.docker.com/orca | bash
```

(or download a bundle with the script and saved/exported images)

Modes/Flags:

* version: Specify an exact version to pull, default is "latest"
* clobber: destroy any existing state containers and deploy fresh

Steps:

1. Pre-flight checks of target engine (Version, available ports, etc.)
    * Detect if we're pointed at an orca, swarm, or individual engine
        * Swarm: See deployment flow for existing swarm below
        * Orca:
            * Orca managed swarm: TBD - could start upgrade flow (see below)
            * Externally managed swarm: Fail, tell user to point to swarm or engine
        * Engine: This flow
2. (conditional) clobber existing state if requested
3. Pull images
    * Detect if saved images are present at the same location as the script (file naming scheme TBD) and if detected load those instead of pulling
    * If we don't have them local, and they aren't already on the target system, do a docker login and search for them, and give a good error message if they aren't visible
4. Generate certs if not present in two host volume paths:
    * /etc/ssl/orca: root CA cert and private key pair; orca server key pair
        * This chain is used for the incoming client requests to the orca server - user can replace certs with their own if they want to use their own or well known CA
        * We can expose a mechanism for a user account to get a signed key pair using this CA to authenticate CLIs (or other tools) against orca, mapping to their user account
    * /etc/ssl/swarm0: root CA cert and private key pair (different from above), swarm server key pair
        * This chain is used for communication from orca to swarm, and from swarm to the engines/proxies.  User typically wouldn't replace this chain
5. Generate cert for proxy signed by /etc/ssl/swarm0
    * If DOCKER\_HOST set, use the hostname/ip from there
    * If localhost, use IP and attempt to get hostname right (ugh) -- or maybe we force the user to tell us how to reach their localhost?
6. Deploy proxy with random exposed port
    * **Q: should we try to use 2375, then fall back to random if unavailable to make firewall updates easier?**
7. Verify we can see the proxy we just deployed (if not warn user firewall settings may need to be opened for port XXX)
8. Deploy swarm manager pointed at proxy (punched through to engine's public IP) - **Fail if swarm official port is taken?**
9. Verify we can see the swarm manager we just deployed (if not warn user firewall settings may need to be opened for port XXX)
11. Deploy DB with host volume mount for data directory
12. Deploy Orca server (prefer 80/443, use random ports if unavailable)
13. Add orca as trusted CA cert on local system:
    * Tell user what we're doing before the sudo prompt, instruct them to ^C to skip it
    * Linux: Append Orca CA cert in  /usr/local/share/ca-certificates/orca.pem and run update-ca-certificates
        * Tell the user that they can copy .../orca.pem to other systems and run "update-ca-certificates" to add it as a trusted system
    * Mac: sudo security add-trusted-cert -d -r trustRoot -k "/Library/Keychains/System.keychain" "/private/tmp/certs/orca.cer"
    * (future) Windows: certmgr.exe -add MyCert.cer -s -r localMachine trustedpublisher
14. Verify the Orca server is up before reporting address to user
15. (bonus round!) Download license key based on the users hub account and license Orca accordingly


## Deploy Orca on Existing Swarm
Description: Deploy orca onto an existing swarm.  Engine membership on the swarm is managed externally to orca.

Flow is similar to above, and would likely be implemented as one script, but outlined here in full.

1. Pre-flight checks of target swarm (Version, available ports, etc.)
2. (conditional) clobber existing state if requested
3. Pull images
    * Detect if saved images are present at the same location as the script (file naming scheme TBD) and if detected load those instead of pulling
    * If we don't have them local, and they aren't already on the target system, do a docker login and search for them, and give a good error message if they aren't visible
4. Generate certs if not present in two host volume paths:
    * /etc/ssl/orca: root CA cert and private key pair; orca server key pair
        * This chain is used for the incoming client requests to the orca server - user can replace certs with their own if they want to use their own or well known CA
        * We can expose a mechanism for a user account to get a signed key pair using this CA to authenticate CLIs (or other tools) against orca, mapping to their user account
    * /etc/ssl/swarm0: Copy connection certs from running session -- **If the root CA private key is present, should we copy it too so we can sign certs?**
        * This chain is used for communication from orca to swarm, and from swarm to the engines/proxies.
        * Note: we probably need additional metadata that this is an externally managed swarm besides just omitting the private CA key
5. Deploy DB with host volume mount for data directory
6. Deploy Orca server (prefer 80/443, use random ports if unavailable)
7. Add orca as trusted CA cert on local system  (see install flow above for details)
8. Verify the Orca server is up before reporting address to user
9. (bonus round!) Download license key based on the users hub account and license Orca accordingly


## Upgrade/Patch an Existing Orca With Swarm

Description: Pointed at an existing deployment, upgrade all the orca and swarm related components while persisting the configuration state of the system

1. Pre-flight checks of target system (existing version, desired target version)
    * Swarm mode: Verify at least two managers
        * If only one manager, we could fail and tell user to re-run deploy script for upgrade of single node system, or if the scripts are the same, just start that flow
    * Orca mode:
        * Managed swarm: Verify at least two managers, get temporary certs to talk to swarm, make sure we can talk to two managers
        * External swarm: Error out, tell user to point to swarm directly to ugprade Orca
    * Engine Mode: Reject?  (We could try to figure out if the swarm is already multi-node and try to connect to it, but this is getting complicated...)
2. Pull images
    * Detect if saved images are present at the same location as the script (file naming scheme TBD) and if detected load those instead of pulling
    * If we don't have them local, and they aren't already on the target system, do a docker login and search for them, and give a good error message if they aren't visible
3. Connect to primary swarm manager
4. For each secondary node (not the primary)
    * Deploy new proxy, verify it can be reached
    * Shutdown manager on this node
    * Deploy new manager, pointed at new proxy, verify it can be reached
    * Shutdown old proxy
5. **Can we trigger a manager switch for swarm at this point?**
6. Stop and remove Orca server and db
7. Start Orca db and server
8. Stop primary swarm manager
9. Switch to communicating with secondary swarm manager
10. Remove old primary swarm manager
11. Start old primary swarm manager
12. Health check swarm/orca
13. Discard temporary swarm connection cert

## Upgrade/Patch an Existing Orca with Externally Managed Swarm

Same as Deploy Orca on Existing Swarm flow

* Pre-flight check: Verify pointed at swarm, not orca, fail if orca and instruct user to load credentials for swarm, then proceed


## Add Host To Orca Managed Swarm

```bash
curl https://myorca/addhost | bash   # Unauthed in GET mode
```

1. Fail fast if the swarm isn't managed by this orca (If we don't have a root ca private key in orca-swarm-root-certs)
2. Pre-flight checks of target engine (same as install flow)
3. (conditional) pull images matching the existing orca/swarm
4. Create empty orca-swarm-certs data volume container
5. Prompt user for cred's to orca
6. curl/wget POST to https://myorca/addhost with hostname/IP of the target engine, piping output into the data volume container
    * (server) Verify user has proper rights to add hosts (admin)
    * (server) Using orca-swarm-root-certs, generate key pair for the engine using the hostname/IP specified, and return tar bundle with root CA cert (public portion only), and new server key pair
    * (server) If orca-swarm-root-certs missing private root ca key, reject request, as swarm and its certs are being externally managed, give pointer to official swarm docs explaining how to grow a swarm
7. Determine port number for engine proxy, deploy with random port if default port taken
    * Note: Once proxy is deployed, the remaining steps could be performed server side
8. Deploy swarm manager pointed at proxy (punched through to engine's public IP) - Fail if swarm official port is taken?
    * If we already have enough managers, should we skip this step?
9. Use users credentials against orca, verify new host is present, report success

Questions:

* Would it make sense to just deploy a conditional beachhead in the script, and do the rest of the host deployment logic on the server side?
    * If the engine already trusts our CA chain, just tell the server the endpoint to talk to
    * If the engine is local or doesn't trust our cert chain, deploy a proxy in the script, then tell the server the proxy endpoint
