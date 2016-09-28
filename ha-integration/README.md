# DTR Multi-Node Integration Test Framework


Running the Tests
=================

* Populate your environment with the appropriate flags (self-documented in Makefile)
* `make ha-integration` from ../

## Running on AWS EC2
* Make a security group that allows all inbound and outbound TCP traffic and inbound traffic on UDP ports 4789 and 7946 (for overlay networking to work properly)
* Run m3.medium AWS instance with ubuntu and the security group you just made in step 1
* You'll need the following to pass to the amazonec2 driver that docker-machine uses: secret-key, access-key, security-group, vpc-id, zone, region
* Setup agent forwarding for your ssh then login into the instance you set up
* Install docker, add ubuntu user to the docker group (curl https://get.docker.com | sudo sh; sudo usermod -aG docker ubuntu)
* Install docker-machine
* Log out and in again
* Login to docker and do the following:
```
sudo apt-get install -y git make
sudo apt-get install build-essential
docker pull busybox
git clone git@github.com:docker/dhe-deploy.git
cd dhe-deploy/
```
* If you don't have Go installed, install it as you'll need it when you run `make ha-integration`
* The HA integration tests use the following environment variables so export/modify them as you wish:
```
NUM_MACHINES=3                          # Number of HA nodes and instances to provision
TEST_ARGS="-trace -notify -v ./ha-integration"
MACHINE_DRIVER=amazonec2
MACHINE_CREATE_FLAGS=--amazonec2-secret-key  --amazonec2-access-key  --amazonec2-security-group --amazonec2-vpc-id  --amazonec2-zone --amazonec2-region  --amazonec2-instance-type t2.small --amazonec2-ami ami-9abea4fb
MACHINE_PREFIX=         # Set to differentiate the DTRTest machine names
UCP_PULL_IMAGES=               # if non-empty will pull using registry env vars
DTR_PULL_IMAGES=
PURGE_MACHINES=             # if non-empty will purge any lingering test machines at end of run
REGISTRY_USERNAME=    # required for pulling
REGISTRY_PASSWORD=
REGISTRY_EMAIL=
UCP_REPO=dockerorcadev
UCP_TAG=1.1.0-rc1
DTR_REPO=dockerhubenterprise
DTR_TAG=2.0.0-rc1-004459_gd48039b
```
* If you didn't set the `PULL_IMAGES` variables to 1, then:
  * pull the UCP and DTR images so they can be transferred by the test framework to the cluster nodes
  * or make them from source and update your env vars to refer to the right tags
  * Note: This setup is configured to work with the latest version of UCP which has eNZi so make sure you're pulling the right UCP images
  * Note 2: If the official UCP images don't have eNZi built in, refer to "Use Josh's custom build of UCP 1.1" section in DTR Hacks

* Use -focus 'installs with constraints' in TEST_ARGS to not run any tests, just deploy the cluster

* Use this for setting constraints for various containers:

```
export DTR_LB_CONSTRAINTS="constraint:container!=ucp-controller"
export DTR_CONSTRAINTS="constraint:container!=ucp-controller|constraint:container!=dtr-haproxy"
```

* Cluster setup
  * A cluster will be created for you with the given name prefix or reused if it exists

* From dhe-deply root: `make ha-integration`
  * The initial cluster set up (assuming 4 nodes: 1 UCP controller and 3 DTR replicas) takes about 12 minutes
  * Skipping cluster setup for the same number of machines should take less than 4 minutes

* Example with the generic driver:
```
export MACHINE_DRIVER=generic
export GENERIC_MACHINE_LIST=52.36.173.142,52.40.229.120,52.39.99.139
export MACHINE_CREATE_FLAGS="--generic-ssh-user ubuntu --generic-ssh-key /home/v/.ssh/id_rsa"
export MACHINE_PREFIX=my-awesome-cluster
```

### Setting up the HAProxy load balancer on top of the UCP controller
* Once the cluster is setup, we want to set up a load balancer to direct traffic during tests to all the DTR replicas. We can set up HAProxy on the UCP controller node as follows:
  * SSH into the UCP controller node
  * Write the following Dockerfile and use it to build the HAProxy image:
```
FROM haproxy:1.6-alpine
COPY haproxy.cfg /usr/local/etc/haproxy/haproxy.cfg
```
  * Then put the following in your `haproxy.cfg` file while replacing xxx with external IP addresses of your DTR replica machines:
```
frontend localnodes
	bind *:443
	mode tcp
	default_backend nodes
	timeout client 1m

backend nodes
	mode tcp
	balance roundrobin
	server clone1 xxx:443 check
	server clone2 xxx:443 check
	server clone3 xxx:443 check
	timeout connect 10s
	timeout server 1m
```
  * Then run: `docker build -t haproxy .`
  * Run the following to make sure the config file is valid: `docker run -it --rm --name haproxy-syntax-check haproxy haproxy -c -f /usr/local/etc/haproxy/haproxy.cfg`
  * Then run the HAProxy: `docker run -d -p 443:443 --name haproxy haproxy`
    * Check to make sure it's actually running. If you don't see any errors after running the following command, move to the next step: `docker logs haproxy`
  * Then hit `https://<ip-address-of-ucp-controller>` (you can get this IP either via the AWS dashboard or by doing `docker-machine ls` from the shell you ran `make ha-integration`) and you should be redirected to a DTR replica and possibly redirected again to the primary DTR replica if the first replica wasn't the primary one.

### Setting up NFS for HA storage
#### NFS Server setup
* SSH into the UCP controller and run the following:
```
apt-get update
apt-get install nfs-kernel-server
```
* Create a directory to serve as the central point of shared content, such as `mkdir /var/nfsshare`
* Change the ownership of the folder: `chown nobody:nogroup /var/nfsshare`
* Next we'll share the NFS directory over the networks by changing the file `/etc/exports/`. The nodes refer to all client nodes that are running a DTR replica:
```
/var/nfsshare <node_1_IP>(rw,sync,no_root_squash,no_subtree_check) <node_2_IP>(rw,sync,no_root_squash,no_subtree_check) <node_3_IP>(rw,sync,no_root_squash,no_subtree_check)
```
* Update the NFS table with the new sharing points: `exportfs -a`
* Start the NFS service: `service nfs-kernel-server start`

#### NFS Clients setup
* Repeat the following steps for each client node that's running a DTR replica
* Install the required package: `sudo apt-get update; sudo apt-get install nfs-common`
* Change directories to `/var/lib/docker/volumes/dtr-registry-<replica_id>/_data`
* Mount `/var/nfsshare` of the NFS server to the client:
```
mount -t nfs4 -o actimeo=0 <server_ip_address>:/var/nfsshare `pwd`
```
* Check to make sure the mount worked properly: `mount -t nfs4`
* Do a test by going into the mounted diretory and creating an empty file: `touch test`
  * If after doing `ls` in the server's `/var/nfsshare` directory the file shows up, it's working correctly.

Writing Tests
=================

* DTR is initialized and torn down for each test suite
* The DTRHATestSuite holds an HAFramework, a util object, a load-balanced apiclient and individual apiclient instances for all DTR replicas
* Individual machines can be accessed with the `HAFramework.Machines []Machine`. Useful methods are `machine.MachineSSH`, `machine.GetIP()`, `machine.GetClient()` 


Future nice-to-haves
====================

* Bring-your-own-cluster logic in BuildFramework. Some of the UCP and eNZi validation logic could be still used.
