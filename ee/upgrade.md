---
title: Upgrade Docker EE
description: Learn how to upgrade your Docker Enterprise Edition, to start using the latest features and security patches.
keywords: enterprise, upgrade
redirect_from:
  - /enterprise/upgrade/
---

## Engine 18.09 Upgrades

In Docker Engine 18.09, significant architectural improvements were made to the network architecture in Swarm to increase 
the performance and scale of the built-in load balancing functionality. 

***NOTE:*** These changes introduce new constraints to the upgrade process that, if not correctly followed, can have 
impact on the availability of applications running on the Swarm. These constraints impact any upgrades coming from any 
version before 18.09 to version 18.09 or greater.

## IP Address Consumption in 18.09+
In Swarm overlay networks, each task connected to a network consumes an IP address on that network. Swarm networks have a 
finite amount of IPs based on the `--subnet` configured when the network is created. If no subnet is specified then Swarm 
defaults to a `/24` network with 254 available IP addresses. When the IP space of a network is fully consumed, Swarm tasks 
can no longer be scheduled on that network.

Docker Enterprise Engine 18.09 and later, each Swarm node will consume an IP address from every Swarm network. This IP 
address is consumed by the Swarm internal load balancer on the network. Swarm networks running on Engine versions 18.09 
or greater must be configured to account for this increase in IP usage. Networks at or near consumption prior to engine version 18.09 may have a risk of reaching full utilization that will prevent tasks from being scheduled on to the network. 
Maximum IP consumption per network at any given moment follows the following formula:

```
Max IP Consumed per Network = Number of Tasks on a Swarm Network + Number of Nodes
```
To prevent this from happening, overlay networks should have enough capacity prior to an upgrade to 18.09, such that the network will have enough capacity after the upgrade. The below instructions offer tooling and steps to ensure capacity is measured before performing an upgrade. 

>The above following only applies to containers running on Swarm overlay networks. This does not impact bridge, macvlan, host, or 3rd party docker networks.

## Cluster Upgrade Best Practices
Docker Engine upgrades in Swarm clusters should follow these guidelines in order to avoid exhaustion application downtime. 

* Workloads should not be actively scheduled in the cluster during upgrades. Large version mismatches between managers and workers can cause unintended consequences
* Manager nodes should all be upgraded first before upgrading worker nodes. Upgrading manager nodes sequentially is recommended if live workloads in the cluster during the upgrade.
* Once manager nodes are upgraded worker nodes should be upgraded next and then the Swarm cluster upgrade is complete.
* If running UCP, the UCP upgrade should follow once all of the Swarm engines have been upgraded.


To upgrade Docker Enterprise Edition you need to individually upgrade each of the
following components:

1. Docker Engine.
2. Universal Control Plane (UCP).
3. Docker Trusted Registry (DTR).

While upgrading, some of these components become temporarily unavailable.
So you should schedule your upgrades to take place outside business peak hours
to make sure there's no impact to your business.

## Create a backup

Before upgrading Docker EE, you should make sure you [create a backup](backup.md).
This makes it possible to recover if anything goes wrong during the upgrade.

## Check the compatibility matrix

You should also check the [compatibility matrix](https://success.docker.com/Policies/Compatibility_Matrix),
to make sure all Docker EE components are certified to work with one another.
You may also want to check the
[Docker EE maintenance lifecycle](https://success.docker.com/Policies/Maintenance_Lifecycle),
to understand until when your version may be supported.

## Apply firewall rules

Before you upgrade, make sure:

- Your firewall rules are configured to allow traffic in the ports UCP uses
  for communication. Learn about [UCP port requirements](ucp/admin/install/system-requirements.md#ports-used).
- Make sure you don't have containers or services that are listening on ports
  used by UCP.
- Configure your load balancer to forward TCP traffic to the Kubernetes API
  server port (6443/TCP by default) running on manager nodes.

> Certificates
>
> Externally signed certificates are used by the Kubernetes API server and
> the UCP controller.
{: .important}

## Upgrade Docker Engine

To avoid application downtime, you should be running Docker in Swarm mode and
deploying your workloads as Docker services. That way you can
drain the nodes of any workloads before starting the upgrade.

If you have workloads running as containers as opposed to swarm services,
make sure they are configured with a [restart policy](/engine/admin/start-containers-automatically/).
This ensures that your containers are started automatically after the upgrade.

To ensure that workloads running as Swarm services have no downtime, you need to:

1. Determine if the network is in danger of exhaustion, then
   a. Triage and fix an upgrade that exhausted IP address space, or
   b. Upgrade a service network live to add IP addresses, or 
3. Drain the node you want to upgrade so that services get scheduled in another node.
4. Upgrade the Docker Engine on that node.
5. Make the node available again.

If you do this sequentially for every node, you can upgrade with no
application downtime.
When upgrading manager nodes, make sure the upgrade of a node finishes before
you start upgrading the next node. Upgrading multiple manager nodes at the same
time can lead to a loss of quorum, and possible data loss.

### Determine if the network is in danger of exhaustion

Starting with a cluster with one or more services configured, determine whether some networks 
may require update in order to function correctly after an 18.09 upgrade.

1. SSH into a manager node.

2. Fetch and deploy a service that would exhaust IP addresses in one of its overlay networks, such as (https://raw.githubusercontent.com/ctelfer/moby-lb-upgrade-test/master/low_addrs/docker-compose.yml)


3. Run the following:

```
$ docker run -it --rm -v /var/run/docker.sock:/var/run/docker.sock docker/ip-util-check
```

If the network is in danger of exhaustion, the output will show similar warnings or errors:

```
 Overlay IP Utilization Report
    ----
    Network ex_net1/XXXXXXXXXXXX has an IP address capacity of 29 and uses 28 addresses
            ERROR: network will be over capacity if upgrading Docker engine version 18.06
                   or later.
    ----
    Network ex_net2/YYYYYYYYYYYY has an IP address capacity of 29 and uses 24 addresses
            WARNING: network could exhaust IP addresses if the cluster scales to 5 or more nodes
    ----
    Network ex_net3/ZZZZZZZZZZZZ has an IP address capacity of 61 and uses 52 addresses
            WARNING: network could exhaust IP addresses if the cluster scales to 9 or more nodes
```

####  Triage and fix an upgrade that exhausted IP address space

Starting with a cluster with services that exhaust their overlay address space in 18.09, adjust the deployment to fix this issue.

1. SSH into a manager node.

2. Fetch and deploy a service that exhausts IP addresses in one of its overlay networks such as (https://raw.githubusercontent.com/ctelfer/moby-lb-upgrade-test/master/exhaust_addrs_3_nodes/docker-compose.yml).

3. Check the `docker service ls` output. It will diplay the service that is  unable to completely fill all its replicas such as: 

```
ID                  NAME                MODE                REPLICAS   IMAGE               PORTS
wn3x4lu9cnln        ex_service          replicated          19/24      nginx:latest
```

4. Use `docker service ps ex_service` to find a failed replica such as:

```
ID                  NAME                IMAGE               NODE                DESIRED STATE       CURRENT STATE              ERROR               PORTS
    ... 
i64lee19ia6s         \_ ex_service.11   nginx:latest        tk1706-ubuntu-1     Shutdown            Rejected 7 minutes ago     "node is missing network attac…"
    ... 
```

5. Examine the error using `docker inspect`. In this example, the  `docker inspect i64lee19ia6s` output shows the error in the `Status.Err` field:

```
... 
            "Status": {
                "Timestamp": "2018-08-24T21:03:37.885405884Z",
                "State": "rejected",
                "Message": "preparing",
                **"Err": "node is missing network attachments, ip addresses may be exhausted",**
                "ContainerStatus": {
                    "ContainerID": "",
                    "PID": 0,
                    "ExitCode": 0
                },
                "PortStatus": {}
            },
    ...
```

6. Adjust the `- subnet:` field in `docker-compose.yml` to have a larger subnet such as `- subnet: 10.1.1.0/22`.  

7. Remove the original service and re-deploy with the new compose file. Confirm the adjusted service deployed successfully.

## Manager Upgrades When Moving to 18.09+
The following is a constraint introduced by architectural changes to the Swarm overlay networking when upgrading to 18.09. It only applies to this one-time upgrade and to workloads that are using the Swarm overlay driver. Once upgraded to 18.09, this constraint does not impact future upgrades.

When upgrading to 18.09, manager nodes cannot reschedule new workloads on the managers until all managers have been upgraded to the 18.09 (or higher) version. During the upgrade of the managers, there is a possibility that any new workloads that are scheduled on the managers will fail to schedule until all of the managers have been upgraded. 

In order to avoid any impactful application downtime, it is advised to reschedule any critical workloads on to Swarm worker nodes during the upgrade of managers. Worker nodes and their network functionality will continue to operate independently during any upgrades or outages on the managers. Note that this restriction only applies to managers and not worker nodes.

### Perform a hit-less upgrade

To upgrade an entire Docker environment, use the following steps.


1. SSH into the manager node.


2.Promote two other nodes to manager: 

```
$ docker node promote manager1 
$ docker node promote manager2
```
    
3. Start a stack with clients connecting to services. For example:

```
$ curl "https://raw.githubusercontent.com/ctelfer/moby-lb-upgrade-test/master/upgrade_test_ct/docker-compose.yml" > docker-compose.yml
docker stack deploy --compose-file docker-compose.yml test
```

4. Upgrade all subsequent managers:

   a. SSH into each manager.

   b. Drain containers from the node: 

      ```
      $ docker node update --availability drain $(docker node ls | grep managerY | awk '{print $1}')
      ```

   c. Verify containers have been moved off: 
      
      ```
      $ docker container ls
      ```
 
   d. Upgrade docker to 18.09 on the system.

5. After upgrading all the managers, reactivate all the nodes:
   
   a. SSH into each manager.
   
   b. Run the following to update all the nodes:

    ```
    $ for m in "manager0 manager1 manager2" ; do \
        docker node update --availability active $(docker node ls | grep $m | awk '{print $1}') \
    done
    ```

6. Repeat the steps above for each worker but with two differences:
   a. You muset drain and activeate the workers from a manager.
   b. It is possible to reactivate each worker as soon as the upgrade for that worker is done.

### Troubleshooting the hit-less upgrade.

If you re-activate a manager immediately instead of waiting for upgrades to the other managers, do the following:

1. Run the following and observe rejected tasks:

```
$ docker service_test_ps
```

2. Run the following and look for `"Err": "node is missing network attachments, ip addresses may be exhausted",`. 
XXX is the ID of one of the service tasks. 
 
```
$ docker inspect XXX
```

3. Finish the upgrade and service will resume.

If you forgot to drain the managers first, do the following:

1. Run the following, and observe rejected tasks with an error `"cannot create a swarm scoped …"`:

```
$ docker service_test_ps
```

2. Run the following and look for `"Err": "cannot create a swarm scoped network when swarm is not active",`
XXX is the ID of one of the service tasks. 
 
```
$ docker inspect XXX
```

3. Finish the upgrade and service will resume.

### Drain the node

Start by draining the node so that services get scheduled in another node and
continue running without downtime.
For that, run this command on a manager node:

```
docker node update --availability drain <node>
```

### Perform the upgrade

To upgrade a node individually by operating system, please follow the instructions
listed below:

* [Windows Server](/install/windows/docker-ee.md#update-docker-ee)
* [Ubuntu](/install/linux/docker-ee/ubuntu.md#upgrade-docker-ee)
* [RHEL](/install/linux/docker-ee/rhel.md#upgrade-docker-ee)
* [CentOS](/install/linux/docker-ee/centos.md#upgrade-docker-ee)
* [Oracle Linux](/install/linux/docker-ee/oracle.md#upgrade-docker-ee)
* [SLES](/install/linux/docker-ee/suse.md#upgrade-docker-ee)

### Make the node active

Once you finish upgrading the node, make it available to run workloads. For
this, run:

```
docker node update --availability active <node>
```

## Upgrade UCP

Once you've upgraded the Docker Engine running on all the nodes, upgrade UCP.
You can do this from the UCP web UI.

![UCP update notification banner](images/upgrade-1.png){: .with-border}

Click on the banner, and choose the version you want to upgrade to.

![UCP upgrade page - version selection](images/upgrade-2.png){: .with-border}

Once you click **Upgrade UCP**, the upgrade starts. If you want you can upgrade
UCP from the CLI instead. [Learn more](/ee/ucp/admin/install/upgrade.md).

## Upgrade DTR

Log in into the DTR web UI to check if there's a new version available.

![DTR settings page](images/upgrade-3.png){: .with-border}

Then follow these [instructions to upgrade DTR](/ee/dtr/admin/upgrade.md).
When this is finished, your Docker EE has been upgraded.

## Where to go next

- [Backup Docker EE](backup.md)
