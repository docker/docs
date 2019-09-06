---
title: Install UCP on Azure
description: Learn how to install Docker Universal Control Plane in a Microsoft Azure environment.
keywords: Universal Control Plane, UCP, install, Docker EE, Azure, Kubernetes
---

Docker UCP closely integrates into Microsoft Azure for its Kubernetes Networking 
and Persistent Storage feature set. UCP deploys the Calico CNI provider, in Azure
the Calico CNI leverages the Azure networking infrastructure for data path 
networking and the Azure IPAM for IP address management. There are 
infrastructure prerequisites that are required prior to UCP installation for the 
Calico / Azure integration.

## Docker UCP Networking

Docker UCP configures the Azure IPAM module for Kubernetes to allocate
IP addresses to Kubernetes pods.  The Azure IPAM module requires each Azure
VM that's part of the Kubernetes cluster to be configured with a pool of
IP addresses.

You have two options for deploying the VMs for the Kubernetes cluster on Azure:
- Install the cluster on Azure stand-alone virtual machines. Docker UCP provides
  an [automated mechanism](#configure-ip-pools-for-azure-stand-alone-vms)
  to configure and maintain IP pools for stand-alone Azure VMs. 
- Install the cluster on an Azure virtual machine scale set. Configure the
  IP pools by using an ARM template like 
  [this one](#set-up-ip-configurations-on-an-azure-virtual-machine-scale-set).

The steps for setting up IP address management are different in the two
environments. If you're using a scale set, you set up `ipConfigurations`
in an ARM template. If you're using stand-alone VMs, you set up IP pools
for each VM by using a utility container that's configured to run as a
global Swarm service, which Docker provides.

## Azure Prerequisites 

The following list of infrastructure prerequisites need to be met in order 
to successfully deploy Docker UCP on Azure.

- All UCP Nodes (Managers and Workers) need to be deployed into the same 
Azure Resource Group. The Azure Networking (Vnets, Subnets, Security Groups) 
components could be deployed in a second Azure Resource Group.
- All UCP Nodes (Managers and Workers) need to be attached to the same 
Azure Subnet.
- All UCP (Managers and Workers) need to be tagged in Azure with the 
`Orchestrator` tag. Note the value for this tag is the Kubernetes version number
in the format `Orchestrator=Kubernetes:x.y.z`. This value may change in each 
UCP release. To find the relevant version please see the UCP 
[Release Notes](../../release-notes). For example for UCP 3.0.6 the tag 
would be `Orchestrator=Kubernetes:1.8.15`. 
- The Azure Computer Name needs to match the Node Operating System's Hostname. 
Note this applies to the FQDN of the host including domain names. 
- An Azure Service Principal with `Contributor` access to the Azure Resource 
Group hosting the UCP Nodes. Note, if using a separate networking Resource 
Group the same Service Principal will need `Network Contributor` access to this 
Resource Group.

The following information will be required for the installation:

- `subscriptionId` - The Azure Subscription ID in which the UCP 
objects are being deployed. 
- `tenantId` - The Azure Active Directory Tenant ID in which the UCP 
objects are being deployed. 
- `aadClientId` - The Azure Service Principal ID
- `aadClientSecret` - The Azure Service Principal Secret Key

### Azure Configuration File

For Docker UCP to integrate in to Microsoft Azure, an Azure configuration file 
will need to be placed within each UCP node in your cluster. This file 
will need to be placed at `/etc/kubernetes/azure.json`. 

See the template below. Note entries that do not contain `****` should not be 
changed.

```
{
    "cloud":"AzurePublicCloud", 
    "tenantId": "***",
    "subscriptionId": "***",
    "aadClientId": "***",
    "aadClientSecret": "***",
    "resourceGroup": "***",
    "location": "****",
    "subnetName": "/****",
    "securityGroupName": "****",
    "vnetName": "****",
    "cloudProviderBackoff": false,
    "cloudProviderBackoffRetries": 0,
    "cloudProviderBackoffExponent": 0,
    "cloudProviderBackoffDuration": 0,
    "cloudProviderBackoffJitter": 0,
    "cloudProviderRatelimit": false,
    "cloudProviderRateLimitQPS": 0,
    "cloudProviderRateLimitBucket": 0,
    "useManagedIdentityExtension": false,
    "useInstanceMetadata": true
}
```

There are some optional values for Azure deployments:

- `"primaryAvailabilitySetName": "****",` - The Worker Nodes availability set.
- `"vnetResourceGroup": "****",` - If your Azure Network objects live in a 
seperate resource group.
- `"routeTableName": "****",` - If you have defined multiple Route tables within
an Azure subnet.

More details on this configuration file can be found 
[here](https://github.com/kubernetes/kubernetes/blob/master/pkg/cloudprovider/providers/azure/azure.go).



## Considerations for IPAM Configuration

The subnet and the virtual network associated with the primary interface of
the Azure VMs need to be configured with a large enough address prefix/range. 
The number of required IP addresses depends on the number of pods running
on each node and the number of nodes in the cluster.

For example, in a cluster of 256 nodes, to run a maximum of 128 pods
concurrently on a node, make sure that the address space of the subnet and the
virtual network can allocate at least 128 * 256 IP addresses, _in addition to_
initial IP allocations to VM NICs during Azure resource creation.

Accounting for IP addresses that are allocated to NICs during VM bring-up, set
the address space of the subnet and virtual network to 10.0.0.0/16. This
ensures that the network can dynamically allocate at least 32768 addresses,
plus a buffer for initial allocations for primary IP addresses.

> Azure IPAM, UCP, and Kubernetes
> 
> The Azure IPAM module queries an Azure virtual machine's metadata to obtain
> a list of IP addresses that are assigned to the virtual machine's NICs. The
> IPAM module allocates these IP addresses to Kubernetes pods. You configure the
> IP addresses as `ipConfigurations` in the NICs associated with a virtual machine
> or scale set member, so that Azure IPAM can provide them to Kubernetes when
> requested.
{: .important}

#### Additional Notes

- The `IP_COUNT` variable defines the subnet size for each node's pod IPs. This subnet size is the same for all hosts.
-  The Kubernetes `pod-cidr` must match the Azure Vnet of the hosts. 

## Configure IP pools for Azure stand-alone VMs

Follow these steps once the underlying infrastructure has been provisioned.

### Configure multiple IP addresses per VM NIC

Follow the steps below to configure multiple IP addresses per VM NIC.

1.  Initialize a swarm cluster comprising the virtual machines you created
    earlier. On one of the nodes of the cluster, run:

    ```bash
    docker swarm init
    ```

2.  Note the tokens for managers and workers. You may retrieve the join tokens 
    at any time by running `$ docker swarm join-token manager` or `$ docker swarm 
    join-token worker` on the manager node.
3.  Join two other nodes on the cluster as manager (recommended for HA) by running:

    ```bash
    docker swarm join --token <manager-token>
    ```

4.  Join remaining nodes on the cluster as workers: 

    ```bash
    docker swarm join --token <worker-token>
    ```

5.  Create a file named "azure_ucp_admin.toml" that contains contents from
    creating the Service Principal.

    ```
    $ cat > azure_ucp_admin.toml <<EOF
    AZURE_CLIENT_ID = "<AD App ID field from Step 1>"
    AZURE_TENANT_ID = "<AD Tenant ID field from Step 1>"
    AZURE_SUBSCRIPTION_ID = "<Azure subscription ID>"
    AZURE_CLIENT_SECRET = "<AD App Secret field from Step 1>"
    EOF
    ```

6.  Create a Docker Swarm secret based on the "azure_ucp_admin.toml" file. 

    ```bash
    docker secret create azure_ucp_admin.toml azure_ucp_admin.toml
    ```

7.  Create a global swarm service using the [docker4x/az-nic-ips](https://hub.docker.com/r/docker4x/az-nic-ips/)
    image on Docker Hub. Use the Swarm secret to prepopulate the virtual machines
    with the desired number of IP addresses per VM from the VNET pool. Set the
    number of IPs to allocate to each VM through the IP_COUNT environment variable.
    For example, to configure 128 IP addresses per VM, run the following command: 

    ```bash
    docker service create \
      --mode=global \
      --secret=azure_ucp_admin.toml \
      --log-driver json-file \
      --log-opt max-size=1m \
      --env IP_COUNT=128 \
      --name ipallocator \
      --constraint "node.platform.os == linux" \
      docker4x/az-nic-ips
    ```

## Set up IP configurations on an Azure virtual machine scale set

Configure IP Pools for each member of the VM scale set during provisioning by
associating multiple `ipConfigurations` with the scale set’s
`networkInterfaceConfigurations`. Here's an example `networkProfile`
configuration for an ARM template that configures pools of 32 IP addresses
for each VM in the VM scale set.

```json
"networkProfile": {
  "networkInterfaceConfigurations": [
    {
      "name": "[variables('nicName')]",
      "properties": {
        "ipConfigurations": [
          {
            "name": "[variables('ipConfigName1')]",
            "properties": {
              "primary": "true",
              "subnet": {
                "id": "[concat('/subscriptions/', subscription().subscriptionId,'/resourceGroups/', resourceGroup().name, '/providers/Microsoft.Network/virtualNetworks/', variables('virtualNetworkName'), '/subnets/', variables('subnetName'))]"
              },
              "loadBalancerBackendAddressPools": [
                {
                  "id": "[concat('/subscriptions/', subscription().subscriptionId,'/resourceGroups/', resourceGroup().name, '/providers/Microsoft.Network/loadBalancers/', variables('loadBalancerName'), '/backendAddressPools/', variables('bePoolName'))]"
                }
              ],
              "loadBalancerInboundNatPools": [
                {
                  "id": "[concat('/subscriptions/', subscription().subscriptionId,'/resourceGroups/', resourceGroup().name, '/providers/Microsoft.Network/loadBalancers/', variables('loadBalancerName'), '/inboundNatPools/', variables('natPoolName'))]"
                }
              ]
            }
          },
          {
            "name": "[variables('ipConfigName2')]",
            "properties": {
              "subnet": {
                "id": "[concat('/subscriptions/', subscription().subscriptionId,'/resourceGroups/', resourceGroup().name, '/providers/Microsoft.Network/virtualNetworks/', variables('virtualNetworkName'), '/subnets/', variables('subnetName'))]"
              }
            }
          }
          .
          .
          .
          {
            "name": "[variables('ipConfigName32')]",
            "properties": {
              "subnet": {
                "id": "[concat('/subscriptions/', subscription().subscriptionId,'/resourceGroups/', resourceGroup().name, '/providers/Microsoft.Network/virtualNetworks/', variables('virtualNetworkName'), '/subnets/', variables('subnetName'))]"
              }
            }
          }
        ],
        "primary": "true"
      }
    }
  ]
}
```

## Install UCP 

Run the following command to install UCP on a manager node. The `--pod-cidr`
option maps to the IP address range that you have configured for the Azure
subnet, and the `--host-address` maps to the private IP address of the master node

> Note: The `pod-cidr` range must match the Azure Virtual Network's Subnet host.
> attached the hosts. For example, if the Azure Virtual Network had the range
> `172.0.0.0/16` with Virtual Machines provisioned on an Azure Subnet of
> `172.0.1.0/24`, then the Pod CIDR should also be `172.0.1.0/24`.

```bash
docker container run --rm -it \
  --name ucp \
  -v /var/run/docker.sock:/var/run/docker.sock \
  {{ page.ucp_org }}/{{ page.ucp_repo }}:{{ page.ucp_version }} install \
  --host-address <ucp-ip> \
  --pod-cidr <ip-address-range> \
  --cloud-provider Azure \
  --interactive
```