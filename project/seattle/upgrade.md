# Upgrade Spec


Strawman upgrade flow:

1. Upgrade engine on one controller to 1.12
    * Must be one with a valid CA (detect in the bootstrapper and reject)
2. Upgrade UCP on that controller using the bootstrapper
    * Almost same upgrade logic as today...
    * Stop all containers, preserving data
    * Fixup any cmd parameters that need adjustment
    * Run any migration/upgrade auxilary containers
    * Init new swarm-mode cluster, wiring up to the external CA(s)
    * Define the new proxy/beachhead service, and wait for it to finish standing everything up
3. Upgrade engine on second controller
4. Join second engine to the cluster
    * This will trigger the proxy/beachhead to start on the node
    * proxy/beachhead detects existing containers running the wrong image
    * Stop all containers, preserving data
    * Detect if the CA has valid material.  If not, reach back to the controller to get it (just like join flow)
    * Start new containers with new image
5. User repeats for remaining controller engines (upgrade engine, join cluster)
6. USer repeats for all worker nodes (upgrade engine, join cluster)



Refinement 1:

* Detect 1.12 nodes that are part of classic swarm inventory, but not in the new swarm-mode inventory.  When detected, expose a button "Complete upgrade on 1.12 nodes"
* Button would call the "node join" API on the node, and accept the join on the new cluster


Refinement 2:

* Controller runs goroutine that periodically (or via events) detects 1.12 nodes that are in classic swarm and not in new swarm-mode.
* When detected, automatically joins them.
