---
description: Administer swarm clusters with Docker EE for IBM Cloud
keywords: ibm, ibm cloud, logging, iaas, tutorial
title: Administer swarm clusters with Docker EE for IBM Cloud
---

Docker Enterprise Edition for IBM Cloud (Beta) comes with a variety of integrations that simplify the swarm administration process.

Use the Docker EE for IBM Cloud CLI plug-in (`bx d4ic`) to provision swarm mode clusters and resources. Manage your cluster with the `bx d4ic` plug-in and the Docker EE Universal Control Plane (UCP) web UI.

## Create swarms

Create a Docker EE swarm cluster in IBM Cloud.

> Your beta license allows you to provision up to 20 nodes
>
> During the beta, your cluster can have a maximum of 20 nodes, up to 14 of which can be worker nodes. If you need more nodes than this, work with your Docker representative to acquire an additional Docker EE license.

Before you begin:

* [Complete the setup requirements](/docker-for-ibm-cloud/index.md#prerequisites).
* Make sure that you have the appropriate [IBM Cloud infrastructure permissions](faqs.md).
* Log in to [IBM Cloud infrastructure](https://control.softlayer.com/), select your user profile, and under the **API Access Information** section retrieve your **API Username** and **Authentication Key**.
* [Add your SSH key to IBM Cloud infrastructure](https://knowledgelayer.softlayer.com/procedure/add-ssh-key), note its label, and locate the file path of the private SSH key on your machine.
* Retrieve the Docker EE installation URL that you received in your beta welcome email.

To create a Docker EE for IBM Cloud cluster from the CLI:

1. Log in to the IBM Cloud CLI. If you have a federated ID, use the `--sso` option.

   ```bash
   $ bx login [--sso]
   ```

2. Target the IBM Cloud org and space.

   ```bash
   $ bx target --cf
   ```

3. Review the `bx d4ic create` command parameters. Parameters that are marked `Required` must be provided during the create command. Optional parameters are set to the default.

   | Parameter | Description | Default Value | Required? |
   | ---- | ----------- | ------------- | --- |
   | `--sl-user`, `-u` | [Log in to IBM Cloud infrastructure](https://control.softlayer.com/), select your profile, and locate your **API Username** under the API Access Information section. | | Required |
   | `--sl-api-key`, `-k` | [Log in to IBM Cloud infrastructure](https://control.softlayer.com/), select your profile, and locate your **Authentication Key** under the API Access Information section. | | Required |
   | `--ssh-label`, `--label` | Your IBM Cloud infrastructure SSH key label for the manager node. To create a key, [log in to IBM Cloud infrastructure](https://control.softlayer.com/) and select **Devices > Manage > SSH Keys > Add**. Copy the key label and insert it here. | | Required |
   | `--ssh-key` | The path to the SSH key on your local client that matches the SSH key label in your IBM Cloud infrastructure account. | | Required |
   | `--swarm-name`, `--name` | The name for your swarm and prefix for the names of each node. | d4ic | Required |
   | `--docker-ee-url` | The Docker EE installation URL associated with your subscription. For beta, you received this in your welcome email. | | Required |
   | `--manager` | Deploy 1, 3, or 5 manager nodes. | 3 | Optional |
   | `--workers`, `-w` | Deploy a minimum of 1 and maximum of 10 worker nodes. | 3 | Optional |
   | `--datacenter`, `-d` | The location (data center) that you deploy the cluster to. Availabe locations are dal12, dal13, fra02, hkg02, lon04, par01, syd01, syd04, tor01, wdc06, wdc07. | wdc07 | Optional |
   | `--verbose`, `-v` | Enable verbose mode | | Optional |
   | `--hardware` | If "dedicated" then the nodes are created on hosts with compute instances in the same account. | Shared | Optional |
   | `--manager-machine-type` | The machine type of the manager nodes: u1c.1x2, u1c.2x4, b1c.4x16, b1c.16x64, b1c.32x128, or b1c.56x242. More powerful machine types cost more, but deliver better performance. For example, u1c.2x4 is 2 cores and 4 GB memory, and b1c.56x242 is 56 cores and 242 GB memory. | u1c.2x4 | Optional |
   | `--worker-machine-type` | The machine type of the worker nodes: u1c.1x2, u1c.2x4, b1c.4x16, b1c.16x64, b1c.32x128, or b1c.56x242. More powerful machine types cost more, but deliver better performance. For example, u1c.2x4 is 2 cores and 4 GB memory, and b1c.56x242 is 56 cores and 242 GB memory. | u1c.1x2 | Optional |

4. Create the cluster. Use the `--swarm-name` flag to name your cluster, and fill in the credentials, SSH, and Docker EE installation URL variables with the information that you retrieved before you began.

    ```bash
    $ bx d4ic create --swarm-name my_swarm \
    --sl-user user.name.1234567 \
    --sl-api-key api-key\
    --ssh-label my_ssh_label \
    --ssh-key filepath_to_my_ssh_key \
    --docker-ee-url my_docker-ee-url
    ```

   > Tip to set environment variables
   >
   > You can set your infrastructure API credentials and Docker EE installation URL as environment variables so that you do not have to include them as options when using `bx d4ic` commands. For example:
   >
   > export SOFTLAYER_USERNAME=user.name.1234567
   >
   > export SOFTLAYER_API_KEY=api-key
   >
   > export D4IC_DOCKER_EE_URL=my_docker-ee-url


5. Note the cluster **Name** and **ID**.

   > Swarm provisioning
   >
   > Your cluster is provisioned in two stages, and takes a few minutes to provision. Don't try to modify your cluster just yet!
   > First, the manager node is deployed. Then, the additional infrastructure resources are deployed, including the worker nodes, DTR nodes, load balancers, subnet, and NFS volume.

   * **Provisioning Stage 1**: Check the status of the manager node:
     {% raw %}
    ```bash
     $ docker logs cluster-name_ID

     Apply complete! Resources: 3 added, 0 changed, 0 destroyed.

     The state of your infrastructure has been saved to the path
     below. This state is required to modify and destroy your
     infrastructure, so keep it safe. To inspect the complete state
     use the `terraform show` command.

     State path:

     Outputs:

     manager_public_ip = 169.##.###.##
     swarm_d4ic_id = ID
     swarm_name = cluster-name
     ucp_password = UCP-password
     ```
     {% endraw %}

   * **Provisioning Stage 2**: Check the status of the cluster infrastructure:
     {% raw %}
     ```bash
      $ bx d4ic show --swarm-name cluster-name --sl-user user.name.1234567 --sl-api-key api_key
      Getting swarm information...
      Infrastructure Details

      Swarm
      ID           ID
      Name         cluster-name
      Created By   user.name.1234567

      Nodes
      ID         Name                           Public IP       Private IP      CPU   Memory   Datacenter   Infrakit Group
      46506407   cluster-name-mgr1              169.##.###.##   10.###.###.##   2     4096     wdc07        managers
      ...

      Load Balancers
      ID                                     Name                Address                                          Type
      ID-string                              cluster-name-mgr    cluster-name-mgr-1234567-wdc07.lb.bluemix.net    mgr
      ...

      Subnets
      ID        Gateway          Datacenter
      ID-number 10.###.###.##   wdc07

      NFS Volumes
      ID              ID-number
      Mount Address   fsf-wdc0701b-fz.adn.networklayer.com:/ID_number/data01
      Datacenter      wdc07
      Capacity        20
      Type            ENDURANCE_FILE_STORAGE
      Tier Level      10_IOPS_PER_GB

      OK
      ```
      {% endraw %}

After creating the cluster, [log in to Docker UCP and download the Docker UCP client certificate bundle](#use-the-universal-control-plane).

## Use the Universal Control Plane

Docker EE for IBM Cloud uses [Docker Universal Control Plane (UCP)](/datacenter/ucp/2.2/guides/) to provide integrated container management and security, from development to production.

### Access UCP

Before you begin, [create a cluster](#create-swarms). Note the its **Name** and **ID**.

1. Retrieve your UCP password by using the cluster **Name** and **ID** that you made when you [created the cluster](#create-swarms).

   ```bash
   $ docker logs cluster-name_ID

   ...
   ucp_password = UCP-password
   ...
   ```

2. Retrieve the load balancer IP address.

   ```bash
   $ bx d4ic list --sl-user user.name.1234567 --sl-api-key api_key
   ```

3. Copy the **UCP URL** for your cluster from the `bx d4ic list` command, and in your browser navigate to it.

4. Log in to UCP. Your credentials are `admin` and the UCP password from the `docker logs` command, or the credentials that your admin created for you.

### Download client certificates

[Download the client certificate bundle](/datacenter/ucp/2.2/guides/user/access-ucp/cli-based-access/#download-client-certificates) to create objects and deploy services from a local Docker client.

1. [Access UCP](#access-ucp).
2. Under your user name (for example, **admin**), click **My Profile**.
3. Click **Client Bundles** > **New Client Bundle**. A zip file is generated.
4. In the GUI, you are now shown a label and public key. You can edit the label by clicking the pencil icon and giving it a name, e.g., _d4ic-ucp_.
5. In a terminal, navigate and unzip the client bundle.

   ```bash
   $ cd Downloads && unzip ucp-bundle-admin.zip
   ```

   > Keep your client bundle handy
   >
   > Move the certificate environment variable directory to a safe and accessible location on your machine. It contains secret information. You'll use it a lot!

6. From the client bundle directory, update your `DOCKER_HOST` and `DOCKER_CERT_PATH` environment variables by loading the `env.sh` script contents into your environment.

   ```bash
   $ source env.sh
   ```

   > Set your environment to use Docker EE for IBM Cloud
   >
   > Repeat this to set your environment variables each time you enter a new terminal session, or after you unset your variables, to connect to the Docker EE for IBM Cloud swarm.

7. Verify that your certificates are being sent to Docker Engine. The command returns information on your swarm.

   ```bash
   $ docker info
   ```

## View swarm resources

### Cluster-level resources

To review resources used within a particular Docker EE cluster, use the CLI or UCP.

**CLI**: The `bx d4ic` CLI lists, modifies, and automates cluster infrastructure, as well as the URLs to access UCP, DTR, or exposed Docker services.

  * To review a list of your clusters and their UCP URLs: `bx d4ic list --sl-user user.name.1234567 --sl-api-key api_key`.
  * To review details about the cluster, such as the IP address of manager nodes or the status of the cluster load balancers: `bx d4ic show --swarm-name my_swarm --sl-user user.name.1234567 --sl-api-key api_key`.

**UCP**: The Docker EE Universal Control Plane provides a web UI to manage swarm users and deployed applications. You can view swarm-related stacks, services, containers, images, nodes, networks, volumes, and secrets.

### Account-level resources

For an account-level view of services and infrastructure that can be used in your swarm, log in to your [IBM Cloud](https://console.bluemix.net/) account.

* The IBM Cloud dashboard provides information on connected IBM Cloud services in the account, such as Watson and Internet of Things.
* The IBM Cloud infrastructure portal shows account infrastructure resources such as virtual devices, storage, and networking.

### Other resources

To gather logging and metric data from your swarm, first [enable logging for the cluster](logging.md), and then access the data in your IBM Cloud organization and space.

## UCP and CLIs

Docker EE for IBM Cloud employs a flexible architecture and integration with IBM Cloud that you can use to leverage IBM Cloud resources and customize your swarm environment. Docker EE UCP exposes the standard Docker API, and as such, includes certain functions that instead should be done by using Docker EE for IBM Cloud capabilities.

> Self-healing capabilities so you don't have to modify cluster infrastructure.
>
> Docker EE for IBM Cloud uses the InfraKit toolkit to support self-healing infrastructure. After you create the swarm, the cluster maintains that specified number of nodes. If a manager node fails, you do not need to promote a worker node to manager; the swarm self-recovers the manager node.
>
> Do not use UCP to modify a cluster's underlying infrastructure, such as adding or promoting worker nodes to managers.

The table outlines when to use UCP and when to use the `bx d4ic` CLI for various types of tasks.

| Task type | UCP or `bx d4ic` CLI | Description |
| --- | --- | --- |
| Swarm nodes | CLI | [Create](#create-swarms), update, modify, and [delete](#delete-swarms) swarm nodes. |
| Certificates | UCP and CLI | From UCP, [download client bundles](#download-client-certificates) after swarm is created in CLI, and every time certificates are changed. From the CLI, run the script from the client bundle downloaded from UCP. |
| Labels | UCP | [Add and modify labels](/engine/userguide/labels-custom-metadata/). If the swarm nodes are removed or modified (such as during a rolling update), the labels must be re-created. |
| Access | UCP | Control access and grant permissions by users, roles, and teams. |
| Secrets | UCP and CLI | In UCP, [manage](/datacenter/ucp/2.2/guides/user/secrets/) and [grant access](/datacenter/ucp/2.2/guides/user/secrets/grant-revoke-access/) to secrets for general usage. For IBM Cloud services, use the CLI `bx d4ic key-create` [command](cli-ref.md#bx-d4ic-key-create) to create secrets. |
| Docker services | UCP and CLI | In UCP, view and manage connected services. From the CLI, [bind IBM Cloud services](binding-services.md). |
| Apps | UCP and CLI | In UCP, view and manage connected apps. From the CLI, [deploy apps](deploy.md) and run containers. |
| Registry | UCP and CLI | For UCP, use [Docker Trusted Registry](/datacenter/dtr/2.4/guides/) installed on the manager node. From the CLI, [run Docker and IBM Cloud Container Registry commands](registry.md). |
| Networking | UCP and CLI | Each cluster has three load balancers that you can use to access and expose various services. You cannot configure these load balancer. In UCP, view networks and related resources. Do not deploy services on the same port that the HTTP Routing Mesh uses. From the CLI, [retrieve load balancer URLs and expose services](load-balancer.md). |
| Data Storage Volumes | CLI | From the CLI, [create and connect file storage volumes](persistent-data-volumes.md) to your swarm. Do NOT use UCP to create IBM Cloud infrastructure data storage volumes. |
| Logging | UCP and CLI | From UCP, can send logs to a remote syslog server. From the CLI, [enable logging and monitoring](logging.md) to IBM Cloud and access by using Grafana and Kibana GUIs. |

## Grant user access

For IBM Cloud account access management, consult the [IBM Cloud Identity and Access Management documentation](https://console.bluemix.net/docs/iam/quickstart.html#getstarted).

For Docker EE cluster access management, use the [UCP Access Control documentation](/datacenter/ucp/2.2/guides/access-control/).

## Delete swarms

Before you begin:

* Log in to [IBM Cloud infrastructure](https://control.softlayer.com/), select your user profile, and under the **API Access Information** section retrieve your **API Username** and **Authentication Key**.
* Retrieve the label of your IBM Cloud infrastructure SSH key, and locate the file path of the private SSH key on your machine.

To delete a swarm:

1. Log in to the IBM Cloud CLI. If you have a federated ID, use the `--sso` option.

   ```bash
   $ bx login [--sso]
   ```

2. Target the IBM Cloud org and space:

   ```bash
   $ bx target --cf
   ```

3. Delete the swarm:

   ```bash
   $ bx d4ic delete (--swarm-name my_swarm | --id swarm_ID )\
   --sl-user user.name.1234567 \
   --sl-api-key api_key \
   --ssh-key filepath_to_my_ssh_key \
   [--force]
   ```

4. Restore the default Docker client settings by running the commands shown in the CLI:

   ```none
   unset DOCKER_HOST
   unset DOCKER_TLS_VERIFY
   unset DOCKER_CERT_PATH
   ```
