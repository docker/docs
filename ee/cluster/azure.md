---
description: Get started with Docker Cluster on Azure
keywords: documentation, docs, docker, cluster, infrastructure, automation, Azure
title: Get started with Docker Cluster on Azure
---

>{% include enterprise_label_shortform.md %}

## Prerequisites

- Completed installation of [Docker Desktop Enterprise](../ee/desktop/) on Windows or Mac, or the [Docker Enterprise Engine](https://docs.docker.com/ee/supported-platforms/) on Linux.
- Sign up for the following items for your Azure account:
  - Service Principal UUID
  - Service Principal App Secret
  - Subscription UUID
  - Tenant UUID
- Organizations wishing to provision roles with explicit permissions should refer to [custom roles](https://docs.microsoft.com/en-us/azure/role-based-access-control/custom-roles) 
and [Azure Permissions](https://github.com/kubernetes/cloud-provider-azure/blob/master/docs/azure-permissions.md) for more information.

More information can be found on obtaining these with either the [Azure CLI](https://docs.microsoft.com/en-us/cli/azure/create-an-azure-service-principal-azure-cli?view=azure-cli-latest) or through the [Azure Portal](https://docs.microsoft.com/en-us/azure/active-directory/develop/howto-create-service-principal-portal).

To securely utilize this Azure credential information, we will create a cluster secrets
file which will inject this data into the environment at runtime. For example, create
a file named `my-azure-creds.sh` similar to the following containing your credentials:

```bash
export ARM_CLIENT_ID='aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee'
export ARM_CLIENT_SECRET='ABCDEFGHIJKLMNOPQRSTUVWXYZ01234567890abcdef='
export ARM_SUBSCRIPTION_ID='ffffffff-gggg-hhhh-iiii-jjjjjjjjjjjj'
export ARM_TENANT_ID='kkkkkkkk-llll-mmmm-nnnn-oooooooooooo'
```

This file should be treated as sensitive data with file permissions set appropriately.
To use this file, we _do not_ source or run this file directly in the shell. Instead,
we reference this file via the CLUSTER_SECRETS_FILE variable in our environment before
running cluster:

```bash
$ export CLUSTER_SECRETS_FILE=~/.my-secrets/my-azure-creds.sh
$ docker cluster create ....
```

Docker cluster will bindmount this file into its container runtime to inject the
credential data as needed.

## Create a cluster

When you create a docker cluster in Azure, the cluster created has:
 - 3 UCP Managers
 - 3 Workers
 - 3 DTR Replicas

Create a file called `cluster.yml` in your directory and paste this in:

```yaml
variable:
  region: "Azure region to deploy"
  ucp_password: 
    type: prompt

provider:
  azurerm:
    region: ${region}

cluster:
  engine:
    version: ee-stable-19.03
  ucp:
    version: docker/ucp:3.2.0
    username: admin
    password: ${ucp_password}
  dtr:
    version: docker/dtr:2.7.1

resource:
  azurerm_virtual_machine:
    managers:
      quantity: 3
    registry:
      quantity: 3
    workers:
      quantity: 3

  azurerm_lb:
    ucp:
      instances:
      - managers
      ports:
      - "443:443"
      - "6443:6443"
```

Provide values for the variable section. For instance:

    region: "centralus"

The values will be substituted in the cluster definition. This makes it
easy to define a reusable cluster definition and then change the variables
to create multiple instances of a cluster.

Run `docker cluster create --file cluster.yml --name quickstart`

    $ docker cluster create --file cluster.yml --name quickstart
    Please provide a value for ucp_password:
    Checking for licenses on Docker Hub
    Docker Enterprise Platform 3.0
    Planning cluster on azurerm                                                  OK
    Creating: [===========>                                              ]  19%  [ ]

After about 5-10 minutes, depending on the number of resources requested, the cluster will be provisioned in the cloud and Docker Enterprise Platform installation will begin:

    $ docker cluster create --file cluster.yml --name quickstart
    Please provide a value for ucp_password:
    Checking for licenses on Docker Hub
    Docker Enterprise Platform 3.0
    Planning cluster on azurerm                                                   OK
    Creating: [==========================================================] 100%   OK
    Installing Docker Enterprise Platform                                         OK

After about 15-20 minutes, Docker Enterprise installation will complete:

    $ docker cluster create --file cluster.yml --name quickstart
    Please provide a value for ucp_password:
    Checking for licenses on Docker Hub
    Docker Enterprise Platform 3.0
    Planning cluster on azurerm                                                   OK
    Creating: [==========================================================] 100%   OK
    Installing Docker Enterprise Platform                                         OK
    Installing Docker Enterprise Engine                                           OK
    Installing Docker Universal Control Plane                                     OK
    Installing Docker Trusted Registry                                            OK

    quickstart
    Successfully created context "quickstart"
    Connect to quickstart at:

    https://ucp-e58dd2a77567-y4pl.centralus.cloudapp.azure.com

    e58dd2a77567

After all operations complete succesfully, the cluster ID will be the last statement
to print. You can login to the URL and begin interacting with the cluster.

## View cluster information

To see an inventory of the current clusters you've created, run `docker cluster ls`

    $ docker cluster ls
    ID             NAME         PROVIDER        ENDPOINT                                                     STATE
    e58dd2a77567   quickstart   azurerm         https://ucp-e58dd2a77567-y4pl.centralus.cloudapp.azure.com   running

To see detailed information about an individual cluster, run `docker cluster inspect quickstart`

$ docker cluster inspect quickstart
```yaml
name: quickstart
shortid: e58dd2a77567
variable:
  region: centralus
  ucp_password: xxxxxxxxxx
provider:
  azurerm:
    environment: public
    region: centralus
    version: ~> 1.32.1
cluster:
  dtr:
    version: docker/dtr:2.7.1
  engine:
    storage_volume: /dev/disk/azure/scsi1/lun0
    url: https://storebits.docker.com/ee/ubuntu/sub-xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
    version: ee-stable-19.03
  kubernetes:
    cloud_provider: true
    load_balancer: false
    nfs_storage: false
  subscription:
    id: sub-xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
    license: /data/license/docker-ee.lic
    trial: "True"
  ucp:
    azure_ip_count: "128"
    pod_cidr: 172.31.0.0/16
    username: admin
    version: docker/ucp:3.2.0
resource:
  azurerm_lb:
    ucp:
      _running:
        dns_name: ucp-e58dd2a77567-y4pl.centralus.cloudapp.azure.com
      path: /data/ssl-certs/
      ports:
      - 443:443
      - 6443:6443
  azurerm_virtual_machine:
    managers:
      data_disk_size: "40"
      enable_public_ips: "true"
      instance_type: Standard_DS3_v2
      os: Ubuntu 18.04
      quantity: 3
      role: manager
    registry:
      data_disk_size: "40"
      enable_public_ips: "true"
      instance_type: Standard_DS3_v2
      os: Ubuntu 18.04
      quantity: 3
      role: dtr
    workers:
      data_disk_size: "40"
      enable_public_ips: "true"
      instance_type: Standard_DS3_v2
      os: Ubuntu 18.04
      quantity: 3
      role: worker
```

The information displayed by `docker cluster inspect` can be used as a cluster definition to clone the cluster.

## Use context

Docker cluster creates a context on your local machine. To use this context, and interact with the cluster, run `docker context use quickstart`

    $ docker context use quickstart
    quickstart
    Current context is now "quickstart"

To verify that the client is connected to the cluster, run `docker version`

```bash
$ docker version
Client: Docker Engine - Enterprise
 Version:           19.03.1
 API version:       1.40
 Go version:        go1.12.5
 Git commit:        f660560
 Built:             Thu Jul 25 20:56:44 2019
 OS/Arch:           darwin/amd64
 Experimental:      false

Server: Docker Enterprise 3.0
 Engine:
  Version:          19.03.1
  API version:      1.40 (minimum version 1.12)
  Go version:       go1.12.5
  Git commit:       f660560
  Built:            Thu Jul 25 20:57:45 2019
  OS/Arch:          linux/amd64
  Experimental:     false
 containerd:
  Version:          1.2.6
  GitCommit:        894b81a4b802e4eb2a91d1ce216b8817763c29fb
 runc:
  Version:          1.0.0-rc8
  GitCommit:        425e105d5a03fabd737a126ad93d62a9eeede87f
 docker-init:
  Version:          0.18.0
  GitCommit:        fec3683
 Universal Control Plane:
  Version:          3.2.0
  ApiVersion:       1.40
  Arch:             amd64
  BuildTime:        Wed Jul 17 23:27:40 UTC 2019
  GitCommit:        586d782
  GoVersion:        go1.12.7
  MinApiVersion:    1.20
  Os:               linux
 Kubernetes:
  Version:          1.14+
  buildDate:        2019-06-06T16:18:13Z
  compiler:         gc
  gitCommit:        7cfcb52617bf94c36953159ee9a2bf14c7fcc7ba
  gitTreeState:     clean
  gitVersion:       v1.14.3-docker-2
  goVersion:        go1.12.5
  major:            1
  minor:            14+
  platform:         linux/amd64
 Calico:
  Version:          v3.5.7
  cni:              v3.5.7
  kube-controllers: v3.5.7
  node:             v3.5.7

$ docker context use default
default
Current context is now "default"
```

## Scale a cluster

Open `cluster.yml`. Change the number of workers to 6:

```yaml
resource:
  azurerm_virtual_machine:
    managers:
      quantity: 3
    registry:
      quantity: 3
    workers:
      quantity: 6
```

Since the cluster is already created, the next step is to `update` the cluster's
desired state.  Run  `docker cluster update quickstart --file cluster.yml`

    $ docker cluster update quickstart --file cluster.yml
    Docker Enterprise Platform 3.0
    Preparing quickstart                                                       [OK]
    Planning cluster on azure                                                  [OK]
    Updating: [==================                                            ] 30%

After about 10 minutes the update operation adds the new nodes and joins them to the cluster:

    $ docker cluster update quickstart --file examples/docs.yml
    Please provide a value for ucp_password
    Docker Enterprise Platform 3.0
    Preparing quickstart                                                       [OK]
    Planning cluster on azure                                                  [OK]
    Updating: [==============================================================] 100%
    Installing Docker Enterprise Platform Requirements                         [OK]
    Installing Docker Enterprise Engine                                        [OK]
    Installing Docker Universal Control Plane                                  [OK]
    Installing Docker Trusted Registry                                         [OK]

    e58dd2a77567

A quick `docker cluster inspect e58dd2a77567` will show the worker count increased:

```yaml
...
    workers:
      data_disk_size: "40"
      enable_public_ips: "true"
      instance_type: Standard_DS3_v2
      os: Ubuntu 18.04
      quantity: 6
      role: worker
```

## Backup a cluster

Before we proceed with more operations on the cluster, let's take a backup of the running cluster. To create a full backup of the cluster, run `docker cluster backup quickstart --file "backup-$(date '+%Y-%m-%d').tar.gz" `

Provide a passphrase to encrypt the UCP backup.

    $ docker cluster backup quickstart --file "backup-$(date '+%Y-%m-%d').tar.gz"
    Passphrase for UCP backup:
    Docker Enterprise Platform 3.0
    Create archive file.                                                       [OK]

    Backup of e58dd2a77567 saved to backup-2019-05-07.tar.gz

Save the backups on external storage for disaster recovery.

To restore a cluster, run `docker cluster restore quickstart --file backup-2019-05-07.tar.gz`

Provide the passphrase from the backup step to decrypt the UCP backup.

## Upgrade a cluster

Open `cluster.yml`.  Change the cluster versions:

```yaml
cluster:
  dtr:
    version: docker/dtr:2.7.0
  engine:
    version: ee-stable-19.03.01
  ucp:
    version: docker/ucp:3.2.0
```

Run  `docker cluster update quickstart --file cluster.yml `

    $ docker cluster update quickstart --file examples/docs.yml
    Please provide a value for ucp_password
    Docker Enterprise Platform 3.0
    Preparing quickstart                                                       [OK]
    Planning cluster on azure                                                  [OK]
    Updating: [==============================================================] 100%
    Installing Docker Enterprise Platform Requirements                         [OK]
    Upgrading Docker Enterprise Engine                                         [OK]
    Upgrading Docker Universal Control Plane                                   [OK]
    Upgrading Docker Trusted Registry                                          [OK]

    e58dd2a77567

## Destroy a cluster

When the cluster has reached end-of-life, run `docker cluster rm quickstart`

    $ docker cluster rm quickstart
    Removing quickstart
    Removing: [==========================================================] 100%   OK

    quickstart
    e58dd2a77567


## Where to go next

- [Explore the full list of Cluster commands](./reference/index.md)
- [Cluster configuration file reference](/ee/cluster-file/index.md)
