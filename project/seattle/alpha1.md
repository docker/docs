# UCP Seattle Alpha 1

UCP Seattle is the code name for the upcoming release with swarm-mode
support.  The UCP team is proud to announce the first internal Alpha for
Docker employees to start exploring the new release.  Please be aware
there are many known limitations in this alpha (see below)


## Highlights

* Swarm-mode support
    * UCP Seattle will always run in swarm mode
* Backwards Compatible "Classic Swarm"
    * UCP cert bundles will expose both classic swarm and swarm-mode capability to remote clients
    * `docker run ...` based containers will work and be scheduled across your same cluster
    * `docker logs` `docker exec` and friends can be used against tasks anywhere in the cluster
* No more "engine discovery" mode
    * The source of numerous perf/reliability issues in the field
* Service based deployment
    * Once UCP is installed, you manage the cluster membership using the same commands as OSS (`docker swarm join` etc.)
    * UCP can be installed on top of an existing OSS swarm-mode cluster and retain all the cluster nodes, and service definitions
* Service centric UI workflows
    * Deploy and manage your services from the UCP GUI
    * Includes "experimental" bundle support

## Installation Instructions

### AWS

If you have an AWS account, try out the Docker AWS Edition with UCP/DDC support [<img src=https://s3.amazonaws.com/cloudformation-examples/cloudformation-launch-stack.png>](https://console.aws.amazon.com/cloudformation/home?#/stacks/new?stackName=Docker&templateURL=https://s3.amazonaws.com/docker-for-aws/aws/alpha/aws-v1.12.0-rc3-beta1-ddc.json)



### Local nodes

You can install UCP on an existing swarm or on a stand-alone 1.12 engine

1. Pre-load UCP images on **all** nodes (yes, repeat this for each and every node!)

    ```bash
    for i in $(docker run --rm dockerorcadev/ucp:1.2.0-alpha1 images --list --image-version dev: ) ; do docker pull $i; done
    ```

2. Install UCP on one of the controller nodes

    ```bash
    docker run --rm -it --name ucp \
        -v /var/run/docker.sock:/var/run/docker.sock \
        dockerorcadev/ucp:1.2.0-alpha1 \
        install --image-version dev: [other args...]
    ```


## Known issues/limitations

* Incompatible with DTR
    * Overlay networking limitations (see below)
* Private images are klunky
    * Don't forget to pre-load all required UCP images on **all** nodes in your cluster (or any nodes you plan to join)
* HA Stability issues
    * Various stability glitches when joining HA nodes concurrently
* Promote and Demote are not yet wired up
* Removal of nodes not yet supported
* Upgrade not yet supported
    * You must run uninstall, then re-install to update to a new build (no service definitions will be lost, but all UCP state will be (users, teams, permissions, etc.)
* Overlay networks not supported for classic `docker run` containers
    * Core team is working hard on a resolution

For more details on the open issues check out
https://github.com/docker/orca/milestones and look at the upcoming
Seattle sprints.

