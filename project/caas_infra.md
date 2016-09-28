# CaaS Infrastructure
This doc describes the infrastructure used by the CaaS team.

Note: once the AWS edition is stable we will migrate everything to that.

# Swarm
CaaS maintains a small Docker 1.12 Swarm for internal tooling and testing.
This resides in AWS in the `5468-4868-6991` account.

## Manager
Currently there is a single Swarm manager.  We may grow this as needed if we
happen to need HA.

## Workers
There is currently a single worker.  We will scale this as needed for demand.

# NFS
There is a small instance that has an attached EBS volume that provides NFS
storage for applications that need it. This is mounted on `/mnt/nfs` on each
node.

# ELB
We use an elastic load balancer to provide tier 1 load balancing for
applications and services.  The current LB is `caas-lb`.  It has SSL support
for the `*.caas.docker.io` domain.

# Route53
Domain management is handled with Route53.  We use the zone `caas.docker.io`.

# Applications and Services
The following is a list of long running services that we maintain on the Swarm.

## Universal Control Plane
There is an instance of UCP that is used for Swarm management.  All actions
should be performed using UCP to help test and find use cases.  It is
available at https://ucp.caas.docker.io:8443.  Ping `ehazlett` or `dhiltgen`
in the `#caas` Slack channel to get an account.

## Jenkins
This is the Jenkins instance that is responsible for our automated testing of
Universal Control Plane.  It consists of a Jenkins master service and a
scalable Jenkins "slave" service.

To create the manager service (pinned to the manager):

```
docker service create \
    --name jenkins \
    --publish 8080:8080 \
    --publish 8081:8081 \
    --network jenkins \
    --replicas 1 \
    --restart-condition any \
    --constraint node=ip-10-0-0-10 \
    --mount type=bind,source=/mnt/nfs/jenkins,target=/var/lib/jenkins,writable=true \
    --mount type=bind,source=/var/run/docker.sock,target=/var/run/docker.sock,writable=true \
    ehazlett/jenkins
```

To create the slave service:

```
docker service create \
    --name jenkins-slave \
    --network jenkins \
    --restart-condition any \
    --constraint node.name!=ip-10-0-0-10 \
    --mount type=bind,source=/mnt/nfs,target=/mnt/nfs,writable=true \
    --mount type=bind,source=/var/run/docker.sock,target=/var/run/docker.sock,writable=true \
    --env JENKINS_MANAGER=http://jenkins:8080 \
    --env JENKINS_USERNAME=api \
    --env JENKINS_PASSWORD=zju1zme5mg \
    --env JENKINS_ROOT_PATH=/mnt/nfs/jenkins-slave \
    ehazlett/jenkins slave
```

The slave service is configured to use the `api` user in Jenkins for automating
the setup.  It uses a generic Jenkins image located
https://github.com/ehazlett/dockerfiles/tree/master/jenkins

The slave service uses an NFS mount so it can be re-scheduled throughout the
cluster and maintain state.
