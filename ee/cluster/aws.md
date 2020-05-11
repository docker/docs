---
description: Get started with Docker Cluster on AWS
keywords: documentation, docs, docker, cluster, infrastructure, automation, AWS
title: Get started with Docker Cluster on AWS
---

>{% include enterprise_label_shortform.md %}

## Prerequisites

- Completed installation of [Docker Desktop Enterprise](/desktop/enterprise/) on Windows or Mac, or the [Docker Enterprise Engine](https://docs.docker.com/ee/supported-platforms/) on Linux.
- [Access keys](https://docs.aws.amazon.com/general/latest/gr/aws-sec-cred-types.html#access-keys-and-secret-access-keys) to an AWS subscription. You can provide these credentials in many ways, but the recommended way is to create an `~/.aws/credentials` file. Refer to [AWS CLI configuration](https://docs.aws.amazon.com/cli/latest/userguide/cli-chap-configure.html) for details on creating one.

## Create a cluster
When you create a docker cluster in AWS, the created cluster has:
 - 3 UCP Managers
 - 3 Workers
 - 3 DTR Replicas

Create a `cluster.yml` file with the following information:
```yaml
    variable:
      domain: "YOUR DOMAIN, e.g. docker.com"
      subdomain: "A SUBDOMAIN, e.g. cluster"
      region: "THE AWS REGION TO DEPLOY, e.g. us-east-1"
      email: "YOUR.EMAIL@COMPANY.COM"
      ucp_password:
        type: prompt

    provider:
      acme:
        email: ${email}
        server_url: https://acme-staging-v02.api.letsencrypt.org/directory
      aws:
        region: ${region}
    cluster:
      dtr:
        version: docker/dtr:2.6.5
      engine:
        version: ee-stable-18.09.5
      ucp:
        username: admin
        password: ${ucp_password}
        version: docker/ucp:3.1.6
    resource:
      aws_instance:
        managers:
          instance_type: t2.xlarge
          os: Ubuntu 16.04
          quantity: 3
        registry:
          instance_type: t2.xlarge
          os: Ubuntu 16.04
          quantity: 3
        workers:
          instance_type: t2.xlarge
          os: Ubuntu 16.04
          quantity: 3
      aws_lb:
        apps:
          domain: ${subdomain}.${domain}
          instances:
          - workers
          ports:
          - 80:8080
          - 443:8443
        dtr:
          domain: ${subdomain}.${domain}
          instances:
          - registry
          ports:
          - 443:443
        ucp:
          domain: ${subdomain}.${domain}
          instances:
          - managers
          ports:
          - 443:443
          - 6443:6443
      aws_route53_zone:
        dns:
          domain: ${domain}
          subdomain: ${subdomain}
```
In this example, the cluster takes on the following topology:

![Docker Cluster Topology](./images/docker_cluster_aws.png)

Provide values for the variable section.  For example:

    domain: "docker.notreal"
    subdomain: "quickstart"
    region: "us-east-1"
    email: "cluster@docker.com"

The values are substituted in the cluster definition, which makes it
easy to define a re-usable cluster definition and then change the variables
to create multiple instances of a cluster.

Run `docker cluster create --file cluster.yml --name quickstart`.

    $ docker cluster create --file cluster.yml --name quickstart
    Please provide a value for ucp_password
    Docker Enterprise Platform 3.0
    Preparing quickstart                                                       [OK]
    Planning cluster on aws                                                    [OK]
    Creating: [===========================                                   ] 44%

After approximately 10 minutes, resources are provisioned, and Docker Enterprise installation is started:

    $ docker cluster create --file cluster.yml --name quickstart
    Please provide a value for ucp_password
    Docker Enterprise Platform 3.0
    Preparing quickstart                                                       [OK]
    Planning cluster on aws                                                    [OK]
    Creating: [==============================================================] 100%
    Installing Docker Enterprise Platform Requirements                         [OK]
    docker-ee : Ensure old versions of Docker are not installed.               [-]

After approximately 20 minutes, Docker Enterprise installation completes:

    $ docker cluster create -f examples/docs.yml -n quickstart
    Please provide a value for ucp_password
    Docker Enterprise Platform 3.0
    Preparing quickstart                                                       [OK]
    Planning cluster on aws                                                    [OK]
    Creating: [==============================================================] 100%
    Installing Docker Enterprise Platform Requirements                         [OK]
    Installing Docker Enterprise Engine                                        [OK]
    Installing Docker Universal Control Plane                                  [OK]
    Installing Docker Trusted Registry                                         [OK]

    Successfully created context "quickstart"
    Connect to quickstart at:

     https://ucp.quickstart.docker.notreal

    911c882340b2

After all operations complete succesfully, the cluster ID is the last statement
to print. You can now log in to the URL and begin interacting with the cluster.

## View cluster information

To view an inventory of the clusters you created, run `docker cluster ls`:

    $ docker cluster ls
    ID             NAME         PROVIDER    ENGINE              UCP                DTR                STATE
    911c882340b2   quickstart   acme, aws   ee-stable-18.09.5   docker/ucp:3.1.6   docker/dtr:2.6.5   running

For detailed information about the cluster, run `docker cluster inspect quickstart`.

    $ docker cluster inspect quickstart
```yaml
name: quickstart
shortid: 911c882340b2
variable:
  domain: docker.notreal
  email: cluster@docker.com
  region: us-east-1
  subdomain: quickstart
provider:
  acme:
    server_url: https://acme-staging-v02.api.letsencrypt.org/directory
  aws:
    region: us-east-1
    version: ~> 1.0
cluster:
  dtr:
    version: docker/dtr:2.6.5
  engine:
    storage_volume: /dev/xvdb
    version: ee-stable-18.09.5
  registry:
    url: https://index.docker.io/v1/
    username: user
  ucp:
    username: admin
    version: docker/ucp:3.1.6
resource:
  aws_instance:
    managers:
      instance_type: t2.xlarge
      os: Ubuntu 16.04
      quantity: 3
      role: manager
    registry:
      instance_type: t2.xlarge
      os: Ubuntu 16.04
      quantity: 3
      role: dtr
    workers:
      instance_type: t2.xlarge
      os: Ubuntu 16.04
      quantity: 3
      role: worker
  aws_lb:
    apps:
      domain: quickstart.docker.notreal
      path: /data/ssl-certs/
      ports:
      - 80:8080
      - 443:8443
    dtr:
      domain: quickstart.docker.notreal
      path: /data/ssl-certs/
      ports:
      - 443:443
    ucp:
      domain: quickstart.docker.notreal
      path: /data/ssl-certs/
      ports:
      - 443:443
      - 6443:6443
  aws_route53_zone:
    dns:
      domain: docker.notreal
      subdomain: quickstart
```
The information displayed by `docker cluster inspect` can be used as a cluster definition to clone the cluster.

## Use context

`docker cluster` creates a context on your local machine.  To use this context and interact with the cluster, run `docker context use quickstart`:

    $ docker context use quickstart
    quickstart
    Current context is now "quickstart"

To verify that the client is connected to the cluster, run `docker version`:

    $ docker version

    Client: Docker Engine - Enterprise
    Version:           19.03.0-beta1
    API version:       1.39 (downgraded from 1.40)
    Go version:        go1.12.1
    Git commit:        90dbc83
    Built:             Fri Apr  5 23:35:58 2019
    OS/Arch:           darwin/amd64
    Experimental:      false

    Server: Docker Enterprise 2.1
    Engine:
      Version:          18.09.5
      API version:      1.39 (minimum version 1.12)
      Go version:       go1.10.8
      Git commit:       be4553c
      Built:            Thu Apr 11 06:19:48 2019
      OS/Arch:          linux/amd64
      Experimental:     false
    Universal Control Plane:
      Version:          3.1.6
      ApiVersion:       1.39
      Arch:             amd64
      BuildTime:        Wed Apr 10 22:35:22 UTC 2019
      GitCommit:        944388b
      GoVersion:        go1.10.6
      MinApiVersion:    1.20
      Os:               linux
    Kubernetes:
      Version:          1.11+
      buildDate:        2019-03-26T02:54:43Z
      compiler:         gc
      gitCommit:        2d582ce995b1ff65b89ad851e8b09b6bc1a84c85
      gitTreeState:     clean
      gitVersion:       v1.11.9-docker-1
      goVersion:        go1.10.8
      major:            1
      minor:            11+
      platform:         linux/amd64
    Calico:
      Version:          v3.5.3
      cni:              v3.5.3
      kube-controllers: v3.5.3
      node:             v3.5.3

To change the context back to your local machine, run `docker context use default`:

    $ docker context use default
    default
    Current context is now "default"

## Scale a cluster
Open `cluster.yml`.  Change the number of workers to 6:
```yaml
      workers:
        instance_type: t2.xlarge
        os: Ubuntu 16.04
        quantity: 6
```
Since the cluster is already created, the next step is to `update` the cluster's
desired state. Run  `docker cluster update quickstart --file cluster.yml`:

    $ docker cluster update quickstart --file cluster.yml
    Docker Enterprise Platform 3.0
    Preparing quickstart                                                       [OK]
    Planning cluster on aws                                                    [OK]
    Updating: [==================                                            ] 30%

After approximately 10 minutes, use the `update` operation to add the new nodes and join them to the cluster:

    $ docker cluster update quickstart --file examples/docs.yml
    Please provide a value for ucp_password
    Docker Enterprise Platform 3.0
    Preparing quickstart                                                       [OK]
    Planning cluster on aws                                                    [OK]
    Updating: [==============================================================] 100%
    Installing Docker Enterprise Platform Requirements                         [OK]
    Installing Docker Enterprise Engine                                        [OK]
    Installing Docker Universal Control Plane                                  [OK]
    Installing Docker Trusted Registry                                         [OK]

    911c882340b2

To view the new nodes in the cluster:

    $ docker --context quickstart node ls
    ID                            HOSTNAME                                      STATUS              AVAILABILITY        MANAGER STATUS      ENGINE VERSION
    mpyk5jxkvgnh75cqmfdzddp7g     ip-172-31-0-116.us-east-2.compute.internal    Ready               Active                                  18.09.5
    s0pd7kqjg8ufelwa9ndkbf1k5     ip-172-31-6-9.us-east-2.compute.internal      Ready               Active              Leader              18.09.5
    ddnvnasq8wibtz9kedlvnxru0     ip-172-31-7-9.us-east-2.compute.internal      Ready               Active                                  18.09.5
    vzta920dhpke9nf4vipqtkuuw     ip-172-31-15-210.us-east-2.compute.internal   Ready               Active                                  18.09.5
    tk98g0tfsb9kzri4slqdh2d2x     ip-172-31-18-95.us-east-2.compute.internal    Ready               Active                                  18.09.5
    g1kwut63oule9v0x245ms7wsw     ip-172-31-21-212.us-east-2.compute.internal   Ready               Active                                  18.09.5
    04jgx94jwscgnac2afdzcd9hp *   ip-172-31-25-45.us-east-2.compute.internal    Ready               Active              Reachable           18.09.5
    5ubqk4mojz198sr72m9zegeew     ip-172-31-29-201.us-east-2.compute.internal   Ready               Active                                  18.09.5
    32rthfhjpm9gaz7n5608k5coj     ip-172-31-33-183.us-east-2.compute.internal   Ready               Active                                  18.09.5
    zqg81yv81auy7eot3a1kson2g     ip-172-31-42-49.us-east-2.compute.internal    Ready               Active                                  18.09.5
    qu84bv2zytv5nubcuntkzwbu5     ip-172-31-43-6.us-east-2.compute.internal     Ready               Active                                  18.09.5
    j6kzzog8a2yv4ragpx826juyv     ip-172-31-43-108.us-east-2.compute.internal   Ready               Active              Reachable           18.09.5

## Back up a cluster

Before performing operations on the cluster, perform a full backup of the running cluster by running `docker cluster backup quickstart --file "backup-$(date '+%Y-%m-%d').tar.gz" `.

Provide a passphrase to encrypt the UCP backup.

    $ docker cluster backup quickstart --file "backup-$(date '+%Y-%m-%d').tar.gz"
    Passphrase for UCP backup:
    Docker Enterprise Platform 3.0
    Create archive file.                                                       [OK]

    Backup of 911c882340b2 saved to backup-2019-05-07.tar.gz

Save the backup on external storage for disaster recovery.

To restore a cluster, run `docker cluster restore quickstart --file backup-2019-05-07.tar.gz`.

Provide the passphrase from the backup step to decrypt the UCP backup.

## Upgrade a cluster
Open `cluster.yml`.  Change the cluster versions:
```yaml
cluster:
  dtr:
    version: docker/dtr:2.7.0
  engine:
    version: ee-stable-19.03
  ucp:
    version: docker/ucp:3.2.0
```
Run  `docker cluster update quickstart --file cluster.yml `:

    $ docker cluster update quickstart --file examples/docs.yml
    Please provide a value for ucp_password
    Docker Enterprise Platform 3.0
    Preparing quickstart                                                       [OK]
    Planning cluster on aws                                                    [OK]
    Updating: [==============================================================] 100%
    Installing Docker Enterprise Platform Requirements                         [OK]
    Upgrading Docker Enterprise Engine                                         [OK]
    Upgrading Docker Universal Control Plane                                   [OK]
    Upgrading Docker Trusted Registry                                          [OK]

    911c882340b2

## Destroy a cluster
When the cluster has reached end-of-life, run `docker cluster rm quickstart`:

    $ docker cluster rm quickstart
    Removing quickstart                                                        [OK]
    Removing: [==============================================================] 100%

    quickstart
    911c882340b2

All provisioned resources are destroyed and the context for the cluster is removed.

## Where to go next

- View the quick start guide for [Azure](azure.md)
- [Explore the full list of Cluster commands](/engine/reference/commandline/cluster/)
- [Cluster configuration file reference](cluster-file.md)
