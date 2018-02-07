---
description: CLI reference for Docker for IBM Cloud
keywords: ibm, ibm cloud, cli, iaas, reference
title: CLI reference for Docker EE for IBM Cloud
---

With the Docker EE for IBM Cloud (beta) plug-in for the IBM Cloud CLI, you can manage your Docker swarms alongside other IBM Cloud operations.

## Docker for IBM Cloud plug-in commands

Refer to these commands to manage your Docker EE for IBM Cloud clusters.

* To view a list of commands, run the `bx d4ic help` command.
* For help with a specific command, run `bx d4ic help [command_name]`.
* To view the version of your Docker for IBM Cloud plug-in, run the `bx d4ic version` command.


| Commands | | |
|---|---|---|
| [bx d4ic create](#bx-d4ic-create) | [bx d4ic delete](#bx-d4ic-delete) | [bx d4ic key-create](#bx-d4ic-key-create) |
| [bx d4ic list](#bx-d4ic-list) | [bx d4ic logmet](#bx-d4ic-logmet) | [bx d4ic show](#bx-d4ic-show) |

## bx d4ic create

Create a Docker EE swarm cluster.

### Usage

```bash
$ bx d4ic create --sl-user SOFTLAYER_USERNAME --sl-api-key SOFTLAYER_API_KEY --ssh-label SSH_KEY_LABEL --ssh-key SSH_KEY_PATH --docker-ee-url DOCKER_EE_URL --swarm-name SWARM_NAME [--datacenter DATACENTER] [--workers NUMBER] [--managers NUMBER] [--hardware SHARED|DEDICATED] [--manager-machine-type MANAGER_MACHINE_TYPE] [--worker-machine-type WORKER_MACHINE_TYPE] [--disable-dtr-storage] [-v] [--version VERSION]
```

### Options

| Name, shorthand | Description | Default | Required? |
|---|---|---|---|
| `--sl-user`, `-u` | [Log in to IBM Cloud infrastructure](https://control.softlayer.com/), select your profile, and locate your **API Username** under the API Access Information section. | | Required |
| `--sl-api-key`, `-k` | [Log in to IBM Cloud infrastructure](https://control.softlayer.com/), select your profile, and locate your **Authentication Key** under the API Access Information section. | | Required |
| `--ssh-label`, `--label` | Your IBM Cloud infrastructure SSH key label for the manager node. To create a key, [log in to IBM Cloud infrastructure](https://control.softlayer.com/) and select **Devices > Manage > SSH Keys > Add**. Copy the key label and insert it here. | | Required |
| `--ssh-key` | The path to the SSH key on your local client that matches the SSH key label in your IBM Cloud infrastructure account. | | Required |
| `--swarm-name`, `--name` | The name for your swarm and prefix for the names of each node. | | Required |
| `--docker-ee-url` | The Docker EE installation URL associated with your subscription. [Email IBM](mailto:sealbou@us.ibm.com) to get a trial subscription during the beta. | | Required |
| `--manager` | Deploy 1, 3, or 5 manager nodes. | 3 | Optional |
| `--workers`, `-w` | Deploy a minimum of 1 and maximum of 10 worker nodes. | 3 | Optional |
| `--datacenter`, `-d` | The location (data center) that you deploy the cluster to. Availabe locations are dal12, dal13, fra02, hkg02, lon04, par01, syd01, syd04, tor01, wdc06, wdc07. | wdc07 | Optional |
| `--verbose`, `-v` | Enable verbose mode | | Optional |
| `--hardware` | If "dedicated" then the nodes are created on hosts with compute instances in the same account. | Shared | Optional |
| `--manager-machine-type` | The machine type of the manager nodes: u1c.1x2, u1c.2x4, b1c.4x16, b1c.16x64, b1c.32x128, or b1c.56x242. Higher machine types cost more, but deliver better performance: for example, u1c.2x4 is 2 cores and 4 GB memory, and b1c.56x242 is 56 cores and 242 GB memory. | b1c.4x16 | Optional |
| `--worker-machine-type` | The machine type of the worker nodes: u1c.1x2, u1c.2x4, b1c.4x16, b1c.16x64, b1c.32x128, or b1c.56x242. Higher machine types cost more, but deliver better performance: for example, u1c.2x4 is 2 cores and 4 GB memory, and b1c.56x242 is 56 cores and 242 GB memory. | u1c.1x2 | Optional |
| `--disable-dtr-storage` | By default, the `bx d4ic create` command orders an IBM Cloud Swift API Object Storage account and creates a container named `dtr-container`. If you want to prevent this, include the `--disable-dtr-storage`. Then, [set up IBM Cloud Object Storage](dtr-ibm-cos.md) yourself so that DTR works with your cluster. | Enabled by default. | Optional |
| `--version` | The Docker EE version of the created cluster. For the beta, only the default version is available. | Default version | Optional |


## bx d4ic delete

Delete a Docker EE swarm cluster.

### Usage

```bash
$ bx d4ic delete (--swarm-name SWARM_NAME | --id ID) --sl-user SOFTLAYER_USERNAME --sl-api-key SOFTLAYER_API_KEY --ssh-label SSH_KEY_LABEL --ssh-key SSH_KEY_PATH [--insecure] [--force]
```

### Options

| Name, shorthand | Description | Default | Required? |
|---|---|---|---|
| `--sl-user`, `-u` | [Log in to IBM Cloud infrastructure](https://control.softlayer.com/), select your profile, and locate your **API Username** under the API Access Information section. | | Required |
| `--sl-api-key`, `-k` | [Log in to IBM Cloud infrastructure](https://control.softlayer.com/), select your profile, and locate your **Authentication Key** under the API Access Information section. | | Required |
| `--ssh-label`, `--label` | Your IBM Cloud infrastructure SSH key label for the manager node. To create a key, [log in to IBM Cloud infrastructure](https://control.softlayer.com/) and select **Devices > Manage > SSH Keys > Add**. Copy the key label and insert it here. | | Required |
| `--ssh-key` | The path to the SSH key on your local client that matches the SSH key label in your IBM Cloud infrastructure account. | | Required |
| `--swarm-name`, `--name` | The name of your cluster. If the name is not provided, you must provide the ID. | | Required |
| `--id` | The ID of your cluster. If the ID is not provided, you must provide the name. | | Required |
| `--verbose`, `-v`| Enable verbose mode | | Optional |
| `--insecure` | Do not verify the identity of the remote host and accept any host key. This is not recommended. | | Optional |
| `--force`, `-f` | Force deletion without confirmation. | | Optional |

## bx d4ic key-create

Create a key for a service instance. Before you can create a key, create an IBM Cloud service.

### Usage

```bash
$ bx d4ic key-create (--swarm-name SWARM_NAME | --id ID) --cert-path CERT_PATH --service-name SERVICE_NAME --service-key SERVICE_KEY --sl-user SOFTLAYER_USERNAME --sl-api-key SOFTLAYER_API_KEY
```

### Options

| Name, shorthand | Description | Default | Required? |
|---|---|---|---|
| `--cert-path`, `--cp` | The directory containing the [Docker UCP client certificate bundle](administering-swarms.md#download-client-certificates). | | Required |
| `--swarm-name`, `--name` | The name of your cluster. If the name is not provided, you must provide the ID. | | Required |
| `--id` | The ID of your cluster. If the ID is not provided, you must provide the name. | | Required |
| `--service-name`, `--name` |  Name of an IBM Cloud service. | | Required |
| `--service-key`, `--key` |  Key of an IBM Cloud service. | | Required |
| `--sl-user`, `-u` | [Log in to IBM Cloud infrastructure](https://control.softlayer.com/), select your profile, and locate your **API Username** under the API Access Information section. | | Required |
| `--sl-api-key`, `-k` | [Log in to IBM Cloud infrastructure](https://control.softlayer.com/), select your profile, and locate your **Authentication Key** under the API Access Information section. | | Required |

## bx d4ic list

List the clusters in your Docker EE for IBM Cloud account.

### Usage

```bash
$ bx d4ic list --sl-user SOFTLAYER_USERNAME --sl-api-key SOFTLAYER_API_KEY [--json]
```

### Options

| Name, shorthand | Description | Default | Required? |
|---|---|---|---|
| `--sl-user`, `-u` | [Log in to IBM Cloud infrastructure](https://control.softlayer.com/), select your profile, and locate your **API Username** under the API Access Information section. | | Required |
| `--sl-api-key`, `-k` | [Log in to IBM Cloud infrastructure](https://control.softlayer.com/), select your profile, and locate your **Authentication Key** under the API Access Information section. | | Required |
| `--json` | Prints the output as JSON. | | Optional |


## bx d4ic logmet

Enable or disable transmission of container log and metric data to IBM Cloud [Log Analysis](https://console.bluemix.net/docs/services/CloudLogAnalysis/log_analysis_ov.html#log_analysis_ov) and [Monitoring](https://console.bluemix.net/docs/services/cloud-monitoring/monitoring_ov.html#monitoring_ov) services.

### Usage

```bash
$ bx d4ic logmet (--swarm-name SWARM_NAME | --id ID) --cert-path CERT_PATH --sl-user SOFTLAYER_USERNAME --sl-api-key SOFTLAYER_API_KEY [--enable | --disable]
```

### Options

| Name, shorthand | Description | Default | Required? |
|---|---|---|---|
| `--swarm-name`, `--name` | The name of your cluster. If the name is not provided, you must provide the ID. | | Required |
| `--id` | The ID of your cluster. If the ID is not provided, you must provide the name. | | Required |
| `--cert-path`, `--cp` | The directory containing the [Docker UCP client certificate bundle](administering-swarms.md#download-client-certificates). | | Required |
| `--sl-user`, `-u` | [Log in to IBM Cloud infrastructure](https://control.softlayer.com/), select your profile, and locate your **API Username** under the API Access Information section. | | Required |
| `--sl-api-key`, `-k` | [Log in to IBM Cloud infrastructure](https://control.softlayer.com/), select your profile, and locate your **Authentication Key** under the API Access Information section. | | Required |
| `--enable` | Send log activity to IBM Cloud [Log Analysis](https://console.bluemix.net/docs/services/CloudLogAnalysis/log_analysis_ov.html#log_analysis_ov) and [Monitoring](https://console.bluemix.net/docs/services/cloud-monitoring/monitoring_ov.html#monitoring_ov) services to the ORG and SPACE that you're currently logged in to. You must include either `--enable` or `--disable` in the command. | | Optional |
| `--disable` | Disable sending log activity to IBM Cloud Log Analysis and Monitoring services. You must include either `--enable` or `--disable` in the command. | | Optional |

## bx d4ic show

Show information about the IBM Cloud infrastructure components, such as load balancer URLs, of a specific cluster.

### Usage

```bash
$ bx d4ic show (--swarm-name SWARM_NAME | --id ID) --sl-user SOFTLAYER_USERNAME --sl-api-key SOFTLAYER_API_KEY [--json]
```

### Options

| Name, shorthand | Description | Default | Required? |
|---|---|---|---|
| `--sl-user`, `-u` | [Log in to IBM Cloud infrastructure](https://control.softlayer.com/), select your profile, and locate your **API Username** under the API Access Information section. | | Required |
| `--sl-api-key`, `-k` | [Log in to IBM Cloud infrastructure](https://control.softlayer.com/), select your profile, and locate your **Authentication Key** under the API Access Information section. | | Required |
| `--id` | The ID of the cluster. You must provide either the ID or the swarm name. | | Required |
| `--swarm-name`, `--name` | The name of your cluster. You must provide either the name or the ID.| | Required |
| `--json` | Prints the output as JSON. | | Optional |
