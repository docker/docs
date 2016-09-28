# UCP-Seattle Proxy Spec

* The `ucp-proxy` container can be evolved into a beachhead, as it already has the
`docker.sock` mounted.  The set of existing bootstrapper operations that could
be performed from the proxy are:
	* join: spawn all non-proxy containers depending on whether the node is a
	  replica or not
	* dump-certs
	* images
	* fingerprint - most probably deprecated in favor of `docker swarm join`
	* uninstall - potentially deprecated in favor of `docker node rm`. The
	  uninstall operation would need to remain in the bootstrapper for host
	  cleanups , if that is the case.

* Even though there might be functional overlap between the controller and proxy
containers, let's try to keep them as separate binaries for the UCP-Seattle
release to minimize risk. If the actual code reuse between the controller and
proxy ends up being substantial, we can consider merging them in the same binary
for a future release.

# To Be Determined

* It is still unknown how the role of a node will be determined (manager vs
non-manager). It's possible that this information will be either conveyed 
during the `docker cluster join` operation or that any join node could
become a manager at a later step.

* The final installation process for a UCP controller will be based on the concrete
  designs for the setup of a swarm manager, as we will need to design the
  process by which controller-specific installation parameters are defined
	by the admin, such as external certs or additional SANs. As of now, it seems
	like the manager directive would be passed as `docker swarm join --manager`.

# Proxy Installation
* Use a one-time secret (maybe a short-lived JWT) for phoning home to the
  controller to do all the "join" logic (CSR dance, etc.)
* Preferred strategy: Use a `mode:fill` service definition to make sure our
  proxy runs on every node as soon as they join the cluster.
	* Note: we can't feed the secret in via env, as that would need to
	  change, and re-create the proxy every time
	* To workaround this limitation, we could side-band the secret into the
	  `ucp-node-certs` volume, and have the proxy loop until it appears.
	  Using channels, we can concurrently block on the appearance of the
	  secret, a timeout and any error returns from the controller.
* Fallback strategy: Integrate into the node lifecycle operations (see node
  spec) and potentially event stream so the existing controller(s) are
notified on node join to the cluster.
	* Secrets (or maybe all the certs?) side-band loaded onto the node named
	  volume(s) first, then start the proxy
	* This strategy will be throw-away code - in the future we'll use service
	  definitions, we just might not be able to for this release

# Proxy Initialization
* The proxy identifies whether the node should be initialized as a
  controller/manager or a node/agent. 
* In the event of a controller, the proxy identifies and locates the expected 
configuration of this specific controller (SANs, host address, external
server certs). Details TBD
* Automatically clean up any cruft that's there from any prior joins - the
  active manager and cluster state trumps any local state that was there.
* Different setup for managers vs agents

* Agents:
	* Agents get non-replica join bits (just the v1 swarm-join for now)
	* Check for local certs in the volumes
	* If not found, use a short-lived secret from the volume to call CSR/join routines 
	on an existing controller
		* If for some reason the secret is invalid, abort, and let the
		  existing controller try again
		* Get the certs, stuff them in the volumes
	* Spawn the other required containers (only the swarm v1 join in this case)
	* Get the advertise address from how swarmv2 wired the node up (this is
	  autodetected for non-manager nodes)
* Managers:
	* Same initialization process as agents, with more containers to manage
	* Use a random secret to transfer the root key material over securely
	* Once we have a valid root replica, update the cluster config


## Other notable aspects

* Like etcd, if quorum is lost, the system becomes wedged to prevent split-brain.
Also like etcd, there isn't a super simple recovery model.  You must shut down
all managers, then bounce one of them with a `--force` flag and it will reset
the cluster.

* It's still a little murky if you can join a manager that's not an agent to a
cluster.  If we can not run containers against these non-agent manager nodes, we
may need to tweak the CA workflow so we can detect this and reject the join
attempt and tell the user they must also make the node an agent, then allow
prevention of user workloads using scheduler hacks like we do for swarm v1.
